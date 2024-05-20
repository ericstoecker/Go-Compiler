package vm

import (
	"compiler/code"
	"compiler/compiler"
	"compiler/object"
	"fmt"
)

const StackSize = 2048

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
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()

			rightInt := right.(*object.Integer)
			leftInt := left.(*object.Integer)

			vm.push(&object.Integer{Value: leftInt.Value + rightInt.Value})
		case code.OpSub:
			right := vm.pop()
			left := vm.pop()

			rightInt := right.(*object.Integer)
			leftInt := left.(*object.Integer)

			vm.push(&object.Integer{Value: leftInt.Value - rightInt.Value})
		case code.OpMul:
			right := vm.pop()
			left := vm.pop()

			rightInt := right.(*object.Integer)
			leftInt := left.(*object.Integer)

			vm.push(&object.Integer{Value: leftInt.Value * rightInt.Value})
		case code.OpDiv:
			right := vm.pop()
			left := vm.pop()

			rightInt := right.(*object.Integer)
			leftInt := left.(*object.Integer)

			vm.push(&object.Integer{Value: leftInt.Value / rightInt.Value})
		case code.OpPop:
			vm.pop()
		}
	}

	return nil
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
