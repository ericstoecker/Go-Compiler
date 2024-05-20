package vm

import (
	"compiler/code"
	"compiler/compiler"
	"compiler/object"
	"fmt"
)

const StackSize = 2048

var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,
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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			vm.executeBinaryOperation(op)
		case code.OpPop:
			vm.pop()
		}
	}

	return nil
}

func (vm *VM) executeBinaryOperation(op code.Opcode) {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	switch {
	case rightType == object.INT && leftType == object.INT:
		leftValue := left.(*object.Integer).Value
		rightValue := right.(*object.Integer).Value
		vm.executeBinaryIntegerOperation(op, leftValue, rightValue)
	}
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left int64, right int64) {
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
	}

	vm.push(result)
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
