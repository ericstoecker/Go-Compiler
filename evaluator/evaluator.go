package evaluator

import (
	"compiler/ast"
	"compiler/token"
)

type Evaluator struct {
}

func New() *Evaluator {
	return nil
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
	case *ast.ReturnStatement:
		return nil
	case *ast.LetStatement:
		return nil
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
	default:
		return nil
	}
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
	String() string
}

type IntegerObject struct {
	Value int64
}

func (intObj *IntegerObject) String() string {
	return "10"
}
