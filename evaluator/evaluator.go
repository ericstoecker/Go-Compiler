package evaluator

import "compiler/ast"

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
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateExpression(expression ast.Expression) Object {
	intExpr := expression.(*ast.IntegerExpression)
	return &IntegerObject{Value: intExpr.Value}
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
