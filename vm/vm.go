package vm

import (
	"compiler/code"
	"compiler/compiler"
	"compiler/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 6048

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}
var NULL = &object.Null{}

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int

	globals []object.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),
	}
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(FALSE)
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(NULL)
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreater, code.OpGreaterEqual:
			err := vm.executeComparisonOperation(op)
			if err != nil {
				return err
			}
		case code.OpMinus, code.OpBang:
			err := vm.executePrefixOperation(op)
			if err != nil {
				return err
			}
		case code.OpJumpNotTrue:
			condition := vm.pop()
			booleanCondition, ok := condition.(*object.Boolean)
			if !ok {
				return fmt.Errorf("not able to execute OpJumpNotTrue if non boolean object on top of stack")
			}

			if !booleanCondition.Value {
				jumpPosition := code.ReadUint16(vm.instructions[ip+1:])
				ip = int(jumpPosition - 1)
			} else {
				ip += 2
			}
		case code.OpJump:
			jumpPosition := code.ReadUint16(vm.instructions[ip+1:])
			ip = int(jumpPosition - 1)
		case code.OpSetGlobal:
			globalsIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[globalsIndex] = vm.pop()
		case code.OpGetGlobal:
			globalsIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			obj := vm.globals[globalsIndex]
			if obj == nil {
				panic("trying to get global that does not exist")
			}

			err := vm.push(obj)
			if err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			vm.sp = vm.sp - numElements

			elements := make([]object.Object, numElements)
			for i := 0; i < numElements; i++ {
				elements[i] = vm.stack[vm.sp+i]
			}

			err := vm.push(&object.Array{Elements: elements})
			if err != nil {
				return err
			}
		case code.OpPop:
			vm.pop()
		}
	}

	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	switch {
	case rightType == object.INT && leftType == object.INT:
		leftValue := left.(*object.Integer).Value
		rightValue := right.(*object.Integer).Value
		return vm.executeBinaryIntegerOperation(op, leftValue, rightValue)
	case rightType == object.STRING && leftType == object.STRING:
		leftValue := left.(*object.String).Value
		rightValue := right.(*object.String).Value
		return vm.executeBinaryStringOperation(op, leftValue, rightValue)
	}

	return nil
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left int64, right int64) error {
	result := &object.Integer{}

	switch op {
	case code.OpAdd:
		result.Value = left + right
	case code.OpSub:
		result.Value = left - right
	case code.OpMul:
		result.Value = left * right
	case code.OpDiv:
		result.Value = left / right
	default:
		return fmt.Errorf("operator %d not known for type INT", op)
	}

	vm.push(result)
	return nil
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left string, right string) error {
	switch op {
	case code.OpAdd:
		vm.push(&object.String{Value: left + right})
	default:
		return fmt.Errorf("operator %d not known for type STRING", op)
	}

	return nil
}

func (vm *VM) executeComparisonOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	switch {
	case rightType == object.INT && leftType == object.INT:
		leftValue := left.(*object.Integer).Value
		rightValue := right.(*object.Integer).Value
		return vm.executeIntegerComparison(op, leftValue, rightValue)
	case rightType == object.BOOLEAN && leftType == object.BOOLEAN:
		leftValue := left.(*object.Boolean)
		rightValue := right.(*object.Boolean)
		return vm.executeBooleanComparison(op, leftValue, rightValue)
	case rightType == object.STRING && leftType == object.STRING:
		leftValue := left.(*object.String).Value
		rightValue := right.(*object.String).Value
		return vm.executeStringComparison(op, leftValue, rightValue)
	default:
		return fmt.Errorf("operator %d not known for type %s, %s", op, leftType, rightType)
	}
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left int64, right int64) error {
	switch op {
	case code.OpEqual:
		vm.push(booleanObjectFromBool(left == right))
	case code.OpNotEqual:
		vm.push(booleanObjectFromBool(left != right))
	case code.OpGreater:
		vm.push(booleanObjectFromBool(left > right))
	case code.OpGreaterEqual:
		vm.push(booleanObjectFromBool(left >= right))
	default:
		return fmt.Errorf("operator %d not known for type INT", op)
	}

	return nil
}

func (vm *VM) executeBooleanComparison(op code.Opcode, left *object.Boolean, right *object.Boolean) error {
	switch op {
	case code.OpEqual:
		vm.push(booleanObjectFromBool(left == right))
	case code.OpNotEqual:
		vm.push(booleanObjectFromBool(left != right))
	default:
		return fmt.Errorf("operator %d not known for type BOOLEAN", op)
	}

	return nil
}

func (vm *VM) executeStringComparison(op code.Opcode, left string, right string) error {
	switch op {
	case code.OpEqual:
		vm.push(booleanObjectFromBool(left == right))
	case code.OpNotEqual:
		vm.push(booleanObjectFromBool(left != right))
	default:
		return fmt.Errorf("operator %d not known for type STRING", op)
	}

	return nil
}

func (vm *VM) executePrefixOperation(op code.Opcode) error {
	right := vm.pop()

	switch op {
	case code.OpBang:
		if right.Type() != object.BOOLEAN {
			return fmt.Errorf("operator %d not known for type %s", op, right.Type())
		}

		rightValue := right.(*object.Boolean).Value
		vm.push(booleanObjectFromBool(!rightValue))
	case code.OpMinus:
		if right.Type() != object.INT {
			return fmt.Errorf("operator %d not known for type %s", op, right.Type())
		}

		rightValue := right.(*object.Integer).Value
		vm.push(&object.Integer{Value: -rightValue})
	default:
		return fmt.Errorf("operator %d not known as prefix operator", op)
	}

	return nil
}

func booleanObjectFromBool(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	vm.sp--
	return vm.stack[vm.sp]
}
