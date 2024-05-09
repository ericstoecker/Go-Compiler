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

func (eval *Evaluator) Evaluate(node ast.Node) object.Object {
	return evaluate(node, eval.environment)
}

func evaluate(node ast.Node, env *Environment) object.Object {
	switch v := node.(type) {
	case *ast.ExpressionStatement:
		return evaluate(v.Expression, env)
	case *ast.LetStatement:
		value := evaluate(v.Value, env)
		env.put(v.Name.Value, value)
		return nil
	case *ast.ReturnStatement:
		return &object.Return{ReturnValue: evaluate(v.ReturnValue, env)}
	case *ast.Program:
		var result object.Object
		for _, stmt := range v.Statements {
			result = evaluate(stmt, env)
		}
		return result
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
			newEnv.put(param, evaluate(v.Arguments[i], env))
		}

		result := evaluateBlockStatement(functionObj.Body, newEnv)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt.ReturnValue
		}
		return result
	case *ast.IfExpression:
		evaluatedCondition := evaluate(v.Condition, env)

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
			elements[i] = evaluate(e, env)
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
		result = evaluate(stmt, env)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt
		}
	}

	return result
}

func evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) object.Object {
	expr := evaluate(prefixExpr.Right, env)
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
	left := evaluate(infixExpr.Left, env)
	right := evaluate(infixExpr.Right, env)

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
