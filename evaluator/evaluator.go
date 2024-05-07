package evaluator

import (
	"compiler/ast"
	"compiler/token"
)

type Evaluator struct {
	environment *Environment
}

func New() *Evaluator {
	environment := NewEnvironment()
	return &Evaluator{environment: environment}
}

func (eval *Evaluator) Evaluate(program *ast.Program) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = eval.evaluateStatement(stmt, eval.environment)
	}

	return result
}

func (eval *Evaluator) evaluateStatement(stmt ast.Statement, env *Environment) Object {
	switch v := stmt.(type) {
	case *ast.ExpressionStatement:
		return eval.evaluateExpression(v.Expression, env)
	case *ast.LetStatement:
		value := eval.evaluateExpression(v.Value, env)
		eval.environment.put(v.Name.Value, value)
		return nil
	case *ast.ReturnStatement:
		return &ReturnObject{ReturnValue: eval.evaluateExpression(v.ReturnValue, env)}
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateExpression(expression ast.Expression, env *Environment) Object {
	switch v := expression.(type) {
	case *ast.PrefixExpression:
		return eval.evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return eval.evaluateInfixExpression(v, env)
	case *ast.IntegerExpression:
		return &IntegerObject{Value: v.Value}
	case *ast.BooleanExpression:
		return &BooleanObject{Value: v.Value}
	case *ast.Identifier:
		return env.get(v.Value)
	case *ast.FunctionLiteral:
		params := make([]string, 0)
		for _, param := range v.Parameters {
			params = append(params, param.Value)
		}
		return &FunctionObject{Parameters: params, Body: v.Body}
	case *ast.CallExpression:
		functionObj, ok := env.get(v.TokenLiteral()).(*FunctionObject)
		if !ok {
			return nil
		}

		newEnv := FromEnvironment(env)

		for i, param := range functionObj.Parameters {
			newEnv.put(param, eval.evaluateExpression(v.Arguments[i], env))
		}

		return eval.evaluateBlockStatement(functionObj.Body, newEnv)
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateBlockStatement(blockStmt *ast.BlockStatement, env *Environment) Object {
	var result Object

	for _, stmt := range blockStmt.Statements {
		result = eval.evaluateStatement(stmt, env)

		returnStmt, ok := result.(*ReturnObject)
		if ok {
			return returnStmt.ReturnValue
		}
	}

	return result
}

func (eval *Evaluator) evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) Object {
	expr := eval.evaluateExpression(prefixExpr.Right, env)
	switch prefixExpr.Operator {
	case token.MINUS:
		intExpr, ok := expr.(*IntegerObject)
		if !ok {
			return nil
		}

		return &IntegerObject{Value: -intExpr.Value}
	case token.BANG:
		boolExpr, ok := expr.(*BooleanObject)
		if !ok {
			return nil
		}

		return &BooleanObject{Value: !boolExpr.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) Object {
	left := e.evaluateExpression(infixExpr.Left, env)
	right := e.evaluateExpression(infixExpr.Right, env)

	switch true {
	case left.Type() == "INT" && right.Type() == "INT":
		leftInt := left.(*IntegerObject)
		rightInt := right.(*IntegerObject)

		return e.evaluateIntegerInfixExpression(infixExpr.Operator, leftInt, rightInt)
	case left.Type() == "BOOLEAN" && right.Type() == "BOOLEAN":
		leftBool := left.(*BooleanObject)
		rightBool := right.(*BooleanObject)

		return e.evaluateBooleanInfixExpression(infixExpr.Operator, leftBool, rightBool)
	default:
		return nil
	}
}

func (e *Evaluator) evaluateIntegerInfixExpression(operator token.TokenType, left *IntegerObject, right *IntegerObject) Object {
	switch operator {
	case token.PLUS:
		return &IntegerObject{Value: left.Value + right.Value}
	case token.MINUS:
		return &IntegerObject{Value: left.Value - right.Value}
	case token.ASTERIK:
		return &IntegerObject{Value: left.Value * right.Value}
	case token.EQUALS:
		return &BooleanObject{Value: left.Value == right.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateBooleanInfixExpression(operator token.TokenType, left *BooleanObject, right *BooleanObject) Object {
	switch operator {
	case token.EQUALS:
		return &BooleanObject{Value: left.Value == right.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateEquals(infixExpr *ast.InfixExpression, env *Environment) Object {
	left := e.evaluateExpression(infixExpr.Left, env)
	right := e.evaluateExpression(infixExpr.Right, env)

	intLeft, ok := left.(*BooleanObject)
	if !ok {
		return nil
	}

	intRight, ok := right.(*BooleanObject)
	if !ok {
		return nil
	}
	return &BooleanObject{Value: intLeft.Value == intRight.Value}
}

type ObjectType string

type Object interface {
	Type() ObjectType
}

type IntegerObject struct {
	Value int64
}

func (intObj *IntegerObject) Type() ObjectType { return "INT" }

type FunctionObject struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (funcObj *FunctionObject) Type() ObjectType { return "FUNCTION" }

type ReturnObject struct {
	ReturnValue Object
}

func (returnObj *ReturnObject) Type() ObjectType { return "RETURN_OBJ" }

type BooleanObject struct {
	Value bool
}

func (boolObj *BooleanObject) Type() ObjectType { return "BOOLEAN" }
