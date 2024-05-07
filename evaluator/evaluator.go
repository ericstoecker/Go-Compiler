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
	switch prefixExpr.Operator {
	case token.MINUS:
		expr := eval.evaluateExpression(prefixExpr.Right, env)
		intExpr, ok := expr.(*IntegerObject)
		if !ok {
			return nil
		}

		return &IntegerObject{Value: -intExpr.Value}
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) Object {
	left := eval.evaluateExpression(infixExpr.Left, env)
	right := eval.evaluateExpression(infixExpr.Right, env)
	intLeft, ok := left.(*IntegerObject)
	if !ok {
		return nil
	}

	intRight, ok := right.(*IntegerObject)
	if !ok {
		return nil
	}

	switch infixExpr.Operator {
	case token.PLUS:
		return &IntegerObject{Value: intLeft.Value + intRight.Value}
	case token.MINUS:
		return &IntegerObject{Value: intLeft.Value - intRight.Value}
	case token.ASTERIK:
		return &IntegerObject{Value: intLeft.Value * intRight.Value}
	default:
		return nil
	}
}

type Object interface {
	object()
}

type IntegerObject struct {
	Value int64
}

func (intObj *IntegerObject) object() {}

type FunctionObject struct {
	Parameters []string
	Body       *ast.BlockStatement
}

func (funcObj *FunctionObject) object() {}

type ReturnObject struct {
	ReturnValue Object
}

func (returnObj *ReturnObject) object() {}
