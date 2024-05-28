package vm

import (
	"compiler/code"
	"compiler/compiler"
	"compiler/object"
	"fmt"
)

const StackSize = 2048
const GlobalsSize = 6048
const FramesSize = 1024

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}
var NULL = &object.Null{}

type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames     []*Frame
	frameIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn)

	frames := make([]*Frame, FramesSize)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]object.Object, StackSize),
		sp:    0,

		globals: make([]object.Object, GlobalsSize),

		frames:     frames,
		frameIndex: 0,
	}
}

func (vm *VM) LastPopped() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.currentFrame().Instructions()); ip++ {
		op := code.Opcode(vm.currentFrame().Instructions()[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
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
				jumpPosition := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
				ip = int(jumpPosition - 1)
			} else {
				ip += 2
			}
		case code.OpJump:
			jumpPosition := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
			ip = int(jumpPosition - 1)
		case code.OpSetGlobal:
			globalsIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
			ip += 2

			vm.globals[globalsIndex] = vm.pop()
		case code.OpGetGlobal:
			globalsIndex := code.ReadUint16(vm.currentFrame().Instructions()[ip+1:])
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
			numElements := int(code.ReadUint16(vm.currentFrame().Instructions()[ip+1:]))
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
		case code.OpMap:
			numEntries := int(code.ReadUint16(vm.currentFrame().Instructions()[ip+1:]))
			ip += 2

			vm.sp = vm.sp - numEntries*2

			entries := make(map[string]object.Object, numEntries)
			for i := 0; i < numEntries; i++ {
				key := vm.stack[vm.sp+2*i]
				value := vm.stack[vm.sp+2*i+1]

				hashableKey, ok := key.(object.Hashable)
				if !ok {
					return fmt.Errorf("type missmatch: non-hashable object provided as hash key")
				}

				entries[hashableKey.Hash()] = value
			}

			err := vm.push(&object.Map{Entries: entries})
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			switch left.Type() {
			case object.ARRAY:
				arr := left.(*object.Array)

				intIndex, ok := index.(*object.Integer)
				if !ok {
					return fmt.Errorf("type missmatch: cannot use %s as array index", index.Type())
				}

				indexValue := intIndex.Value
				if indexValue < 0 || int(indexValue) >= len(arr.Elements) {
					return fmt.Errorf("index %d out of range for array of length %d", indexValue, len(arr.Elements))
				}

				value := arr.Elements[indexValue]
				err := vm.push(value)
				if err != nil {
					return err
				}
			case object.MAP:
				mapObj := left.(*object.Map)

				hashableIndex, ok := index.(object.Hashable)
				if !ok {
					return fmt.Errorf("type missmatch: cannot use %s as key for hashmap", hashableIndex.Type())
				}

				value := mapObj.Entries[hashableIndex.Hash()]
				if value == nil {
					value = NULL
				}
				err := vm.push(value)
				if err != nil {
					return err
				}
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
	if vm.sp <= 0 {
		panic("trying to pop from empty stack")
	}

	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) pushFrame() {}

func (vm *VM) popFrame() {}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex]
}
