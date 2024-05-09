package evaluator

import (
	"compiler/ast"
	"compiler/object"
	"compiler/token"
)

type Evaluator struct {
	environment *Environment
}

func New() *Evaluator {
	environment := NewEnvironment()
	return &Evaluator{environment: environment}
}

func (eval *Evaluator) Evaluate(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = eval.evaluateStatement(stmt, eval.environment)
	}

	return result
}

func (eval *Evaluator) evaluateStatement(stmt ast.Statement, env *Environment) object.Object {
	switch v := stmt.(type) {
	case *ast.ExpressionStatement:
		return eval.evaluateExpression(v.Expression, env)
	case *ast.LetStatement:
		value := eval.evaluateExpression(v.Value, env)
		eval.environment.put(v.Name.Value, value)
		return nil
	case *ast.ReturnStatement:
		return &object.ReturnObject{ReturnValue: eval.evaluateExpression(v.ReturnValue, env)}
	default:
		return nil
	}
}

func (eval *Evaluator) evaluateExpression(expression ast.Expression, env *Environment) object.Object {
	switch v := expression.(type) {
	case *ast.PrefixExpression:
		return eval.evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return eval.evaluateInfixExpression(v, env)
	case *ast.IntegerExpression:
		return &object.IntegerObject{Value: v.Value}
	case *ast.BooleanExpression:
		return &object.BooleanObject{Value: v.Value}
	case *ast.Identifier:
		return env.get(v.Value)
	case *ast.FunctionLiteral:
		params := make([]string, 0)
		for _, param := range v.Parameters {
			params = append(params, param.Value)
		}
		return &object.FunctionObject{Parameters: params, Body: v.Body}
	case *ast.CallExpression:
		functionObj, ok := env.get(v.TokenLiteral()).(*object.FunctionObject)
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

func (eval *Evaluator) evaluateBlockStatement(blockStmt *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object

	for _, stmt := range blockStmt.Statements {
		result = eval.evaluateStatement(stmt, env)

		returnStmt, ok := result.(*object.ReturnObject)
		if ok {
			return returnStmt.ReturnValue
		}
	}

	return result
}

func (eval *Evaluator) evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) object.Object {
	expr := eval.evaluateExpression(prefixExpr.Right, env)
	switch prefixExpr.Operator {
	case token.MINUS:
		intExpr, ok := expr.(*object.IntegerObject)
		if !ok {
			return nil
		}

		return &object.IntegerObject{Value: -intExpr.Value}
	case token.BANG:
		boolExpr, ok := expr.(*object.BooleanObject)
		if !ok {
			return nil
		}

		return &object.BooleanObject{Value: !boolExpr.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) object.Object {
	left := e.evaluateExpression(infixExpr.Left, env)
	right := e.evaluateExpression(infixExpr.Right, env)

	switch true {
	case left.Type() == "INT" && right.Type() == "INT":
		leftInt := left.(*object.IntegerObject)
		rightInt := right.(*object.IntegerObject)

		return e.evaluateIntegerInfixExpression(infixExpr.Operator, leftInt, rightInt)
	case left.Type() == "BOOLEAN" && right.Type() == "BOOLEAN":
		leftBool := left.(*object.BooleanObject)
		rightBool := right.(*object.BooleanObject)

		return e.evaluateBooleanInfixExpression(infixExpr.Operator, leftBool, rightBool)
	default:
		return nil
	}
}

func (e *Evaluator) evaluateIntegerInfixExpression(operator token.TokenType, left *object.IntegerObject, right *object.IntegerObject) object.Object {
	switch operator {
	case token.PLUS:
		return &object.IntegerObject{Value: left.Value + right.Value}
	case token.MINUS:
		return &object.IntegerObject{Value: left.Value - right.Value}
	case token.ASTERIK:
		return &object.IntegerObject{Value: left.Value * right.Value}
	case token.EQUALS:
		return &object.BooleanObject{Value: left.Value == right.Value}
	case token.NOT_EQUALS:
		return &object.BooleanObject{Value: left.Value != right.Value}
	case token.GT:
		return &object.BooleanObject{Value: left.Value > right.Value}
	case token.LT:
		return &object.BooleanObject{Value: left.Value < right.Value}
	case token.GREATER_EQUAL:
		return &object.BooleanObject{Value: left.Value >= right.Value}
	case token.LESS_EQUAL:
		return &object.BooleanObject{Value: left.Value <= right.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateBooleanInfixExpression(operator token.TokenType, left *object.BooleanObject, right *object.BooleanObject) object.Object {
	switch operator {
	case token.EQUALS:
		return &object.BooleanObject{Value: left.Value == right.Value}
	case token.NOT_EQUALS:
		return &object.BooleanObject{Value: left.Value != right.Value}
	default:
		return nil
	}
}

func (e *Evaluator) evaluateEquals(infixExpr *ast.InfixExpression, env *Environment) object.Object {
	left := e.evaluateExpression(infixExpr.Left, env)
	right := e.evaluateExpression(infixExpr.Right, env)

	intLeft, ok := left.(*object.BooleanObject)
	if !ok {
		return nil
	}

	intRight, ok := right.(*object.BooleanObject)
	if !ok {
		return nil
	}
	return &object.BooleanObject{Value: intLeft.Value == intRight.Value}
}
