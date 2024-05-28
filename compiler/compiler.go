package compiler

import (
	"compiler/ast"
	"compiler/code"
	"compiler/object"
	"fmt"
)

type Compiler struct {
	constants []object.Object

	scopes     []*CompilationScope
	scopeIndex int

	symbolTable *SymbolTable
}

type CompilationScope struct {
	instructions code.Instructions

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type EmittedInstruction struct {
	code     code.Opcode
	position int
}

func New() *Compiler {
	mainScope := &CompilationScope{
		instructions: code.Instructions{},
	}

	return &Compiler{
		scopes:     []*CompilationScope{mainScope},
		scopeIndex: 0,

		symbolTable: NewSymbolTable(),
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(node.Name.Value)

		c.emit(code.OpSetGlobal, symbol.Index)
	case *ast.Identifier:
		symbol, ok := c.symbolTable.RetrieveSymbol(node.Value)
		if !ok {
			return fmt.Errorf("undefined: %s", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)
	case *ast.ArrayLiteral:
		elements := node.Elements
		for _, e := range elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(elements))
	case *ast.MapLiteral:
		entries := node.Entries
		for key, value := range entries {
			err := c.Compile(key)
			if err != nil {
				return err
			}

			err = c.Compile(value)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpMap, len(entries))
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		opJumpNotTruePosition := c.emit(code.OpJumpNotTrue, 0)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastInstruction()
		}

		opJumpPosition := c.emit(code.OpJump, 0)
		c.replaceInstruction(opJumpNotTruePosition, code.Make(code.OpJumpNotTrue, len(c.currentInstructions())))

		if node.Alternative != nil {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastInstruction()
			}

		} else {
			c.emit(code.OpNull)
		}

		c.replaceInstruction(opJumpPosition, code.Make(code.OpJump, len(c.currentInstructions())))
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.FunctionLiteral:
		c.enterScope()

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastInstruction()
			c.emit(code.OpReturnValue)
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		functionInstructions := c.leaveScope()

		compiledFn := &object.CompiledFunction{Instructions: functionInstructions}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	case *ast.CallExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		c.emit(code.OpCall)
	case *ast.InfixExpression:
		switch node.Operator {
		case "<=", "<":
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
		default:
			err := c.Compile(node.Left)
			if err != nil {
				return err
			}

			err = c.Compile(node.Right)
			if err != nil {
				return err
			}
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">=", "<=":
			c.emit(code.OpGreaterEqual)
		case ">", "<":
			c.emit(code.OpGreater)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.StringLiteral:
		string := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(string))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.scopes[c.scopeIndex].previousInstruction = c.scopes[c.scopeIndex].lastInstruction
	c.scopes[c.scopeIndex].lastInstruction = EmittedInstruction{code: op, position: pos}
	return pos
}

func (c *Compiler) removeLastInstruction() {
	c.scopes[c.scopeIndex].instructions = c.currentInstructions()[:c.scopes[c.scopeIndex].lastInstruction.position]
	c.scopes[c.scopeIndex].lastInstruction = c.scopes[c.scopeIndex].previousInstruction
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	return c.scopes[c.scopeIndex].lastInstruction.code == op
}

func (c *Compiler) replaceInstruction(pos int, ins []byte) {
	for i, instruction := range ins {
		c.currentInstructions()[pos+i] = instruction
	}
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	c.scopes[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
	return posNewInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	newScope := &CompilationScope{
		instructions: code.Instructions{},
	}

	c.scopes = append(c.scopes, newScope)
	c.scopeIndex++
}

func (c *Compiler) leaveScope() code.Instructions {
	instructionsInCurrentScope := c.scopes[c.scopeIndex].instructions

	c.scopes = c.scopes[:c.scopeIndex]
	c.scopeIndex--

	return instructionsInCurrentScope
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
