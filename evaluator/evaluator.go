package evaluator

import (
	"compiler/ast"
	"compiler/token"
)

type Evaluator struct {
	environment map[string]Object
}

func New() *Evaluator {
	environment := make(map[string]Object)
	return &Evaluator{environment: environment}
}

func (eval *Evaluator) Evaluate(program *ast.Program) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = eval.evaluateStatement(stmt)
	}

	return result
}

func (eval *Evaluator) evaluateStatement(stmt ast.Statement) Object {
	switch v := stmt.(type) {
	case *ast.ExpressionStatement:
		return eval.evaluateExpression(v.Expression)
	case *ast.LetStatement:
		value := eval.evaluateExpression(v.Value)
		eval.environment[v.Name.Value] = value
		return nil
	case *ast.ReturnStatement:
		return &ReturnObject{ReturnValue: eval.evaluateExpression(v.ReturnValue)}
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateExpression(expression ast.Expression) Object {
	switch v := expression.(type) {
	case *ast.PrefixExpression:
		return eval.evaluatePrefixExpression(v)
	case *ast.InfixExpression:
		return eval.evaluateInfixExpression(v)
	case *ast.IntegerExpression:
		return &IntegerObject{Value: v.Value}
	case *ast.Identifier:
		return eval.environment[v.Value]
	case *ast.FunctionLiteral:
		params := make([]string, 0)
		for _, param := range v.Parameters {
			params = append(params, param.Value)
		}
		return &FunctionObject{Parameters: params, Body: v.Body}
	case *ast.CallExpression:
		functionObj, ok := eval.environment[v.TokenLiteral()].(*FunctionObject)
		if !ok {
			return nil
		}

		for i, param := range functionObj.Parameters {
			eval.environment[param] = eval.evaluateExpression(v.Arguments[i])
		}

		return eval.evaluateBlockStatement(functionObj.Body)
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateBlockStatement(blockStmt *ast.BlockStatement) Object {
	var result Object

	for _, stmt := range blockStmt.Statements {
		result = eval.evaluateStatement(stmt)

		returnStmt, ok := result.(*ReturnObject)
		if ok {
			return returnStmt.ReturnValue
		}
	}

	return result
}

func (eval *Evaluator) evaluatePrefixExpression(prefixExpr *ast.PrefixExpression) Object {
	switch prefixExpr.Operator {
	case token.MINUS:
		expr := eval.evaluateExpression(prefixExpr.Right)
		intExpr, ok := expr.(*IntegerObject)
		if !ok {
			return nil
		}

		return &IntegerObject{Value: -intExpr.Value}
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateInfixExpression(infixExpr *ast.InfixExpression) Object {
	left := eval.evaluateExpression(infixExpr.Left)
	right := eval.evaluateExpression(infixExpr.Right)
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
