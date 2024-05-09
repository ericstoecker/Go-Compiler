package evaluator

import (
	"compiler/ast"
	"compiler/object"
	"compiler/token"
	"fmt"
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
		result = evaluateStatement(stmt, eval.environment)
	}

	return result
}

func evaluateStatement(stmt ast.Statement, env *Environment) object.Object {
	switch v := stmt.(type) {
	case *ast.ExpressionStatement:
		return evaluateExpression(v.Expression, env)
	case *ast.LetStatement:
		value := evaluateExpression(v.Value, env)
		env.put(v.Name.Value, value)
		return nil
	case *ast.ReturnStatement:
		return &object.Return{ReturnValue: evaluateExpression(v.ReturnValue, env)}
	default:
		return nil
	}
}

func evaluateExpression(expression ast.Expression, env *Environment) object.Object {
	switch v := expression.(type) {
	case *ast.PrefixExpression:
		return evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return evaluateInfixExpression(v, env)
	case *ast.IntegerExpression:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanExpression:
		return &object.Boolean{Value: v.Value}
	case *ast.StringExpression:
		return &object.String{Value: v.Value}
	case *ast.Identifier:
		return env.get(v.Value)
	case *ast.FunctionLiteral:
		params := make([]string, 0)
		for _, param := range v.Parameters {
			params = append(params, param.Value)
		}
		return &object.Function{Parameters: params, Body: v.Body}
	case *ast.CallExpression:
		functionObj, ok := env.get(v.TokenLiteral()).(*object.Function)
		if !ok {
			return nil
		}

		newEnv := FromEnvironment(env)

		for i, param := range functionObj.Parameters {
			newEnv.put(param, evaluateExpression(v.Arguments[i], env))
		}

		result := evaluateBlockStatement(functionObj.Body, newEnv)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt.ReturnValue
		}
		return result
	case *ast.IfExpression:
		evaluatedCondition := evaluateExpression(v.Condition, env)

		booleanEvalResult, ok := evaluatedCondition.(*object.Boolean)
		if !ok {
			return nil
		}

		if booleanEvalResult.Value {
			return evaluateBlockStatement(v.Consequence, env)
		} else if v.Alternative != nil {
			return evaluateBlockStatement(v.Alternative, env)
		}

		// TODO return null object
		return nil
	case *ast.ArrayExpression:
		elements := make([]object.Object, len(v.Elements))
		for i, e := range v.Elements {
			elements[i] = evaluateExpression(e, env)
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		arrObj, ok := env.get(v.TokenLiteral()).(*object.Array)
		if !ok {
			fmt.Print("NOT ARRAY")
			return nil
		}

		index := v.Index.Value
		if 0 <= index && int(index) <= len(arrObj.Elements) {
			return arrObj.Elements[v.Index.Value]
		}
		return nil
	default:
		return nil
	}
}

func evaluateBlockStatement(blockStmt *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object

	for _, stmt := range blockStmt.Statements {
		result = evaluateStatement(stmt, env)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt
		}
	}

	return result
}

func evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) object.Object {
	expr := evaluateExpression(prefixExpr.Right, env)
	switch prefixExpr.Operator {
	case token.MINUS:
		intExpr, ok := expr.(*object.Integer)
		if !ok {
			return nil
		}

		return &object.Integer{Value: -intExpr.Value}
	case token.BANG:
		boolExpr, ok := expr.(*object.Boolean)
		if !ok {
			return nil
		}

		return &object.Boolean{Value: !boolExpr.Value}
	default:
		return nil
	}
}

func evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) object.Object {
	left := evaluateExpression(infixExpr.Left, env)
	right := evaluateExpression(infixExpr.Right, env)

	switch true {
	case left.Type() == object.INT && right.Type() == object.INT:
		leftInt := left.(*object.Integer)
		rightInt := right.(*object.Integer)

		return evaluateIntegerInfixExpression(infixExpr.Operator, leftInt, rightInt)
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		leftBool := left.(*object.Boolean)
		rightBool := right.(*object.Boolean)

		return evaluateBooleanInfixExpression(infixExpr.Operator, leftBool, rightBool)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		leftStr := left.(*object.String)
		rightStr := right.(*object.String)

		return evaluateStringInfixExpression(infixExpr.Operator, leftStr, rightStr)
	default:
		return nil
	}
}

func evaluateIntegerInfixExpression(operator token.TokenType, left *object.Integer, right *object.Integer) object.Object {
	switch operator {
	case token.PLUS:
		return &object.Integer{Value: left.Value + right.Value}
	case token.MINUS:
		return &object.Integer{Value: left.Value - right.Value}
	case token.ASTERIK:
		return &object.Integer{Value: left.Value * right.Value}
	case token.EQUALS:
		return &object.Boolean{Value: left.Value == right.Value}
	case token.NOT_EQUALS:
		return &object.Boolean{Value: left.Value != right.Value}
	case token.GT:
		return &object.Boolean{Value: left.Value > right.Value}
	case token.LT:
		return &object.Boolean{Value: left.Value < right.Value}
	case token.GREATER_EQUAL:
		return &object.Boolean{Value: left.Value >= right.Value}
	case token.LESS_EQUAL:
		return &object.Boolean{Value: left.Value <= right.Value}
	default:
		return nil
	}
}

func evaluateBooleanInfixExpression(operator token.TokenType, left *object.Boolean, right *object.Boolean) object.Object {
	switch operator {
	case token.EQUALS:
		return &object.Boolean{Value: left.Value == right.Value}
	case token.NOT_EQUALS:
		return &object.Boolean{Value: left.Value != right.Value}
	case token.AND:
		return &object.Boolean{Value: left.Value && right.Value}
	case token.OR:
		return &object.Boolean{Value: left.Value || right.Value}
	default:
		return nil
	}
}

func evaluateStringInfixExpression(operator token.TokenType, left *object.String, right *object.String) object.Object {
	switch operator {
	case token.EQUALS:
		return &object.Boolean{Value: left.Value == right.Value}
	case token.NOT_EQUALS:
		return &object.Boolean{Value: left.Value != right.Value}
	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}
	default:
		return nil
	}
}
