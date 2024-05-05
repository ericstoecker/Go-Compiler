package evaluator

import "compiler/ast"

type Evaluator struct {
}

func New() *Evaluator {
	return nil
}

func (eval *Evaluator) Evaluate(program *ast.Program) Object {
	stmt := program.Statements[0]
	exprStmt := stmt.(*ast.ExpressionStatement)
	intStmt := exprStmt.Expression.(*ast.IntegerExpression)
	return &IntegerObject{Value: intStmt.Value}
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
