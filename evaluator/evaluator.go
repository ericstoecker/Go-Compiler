package evaluator

import (
	"compiler/ast"
	"compiler/object"
	"compiler/token"
	"fmt"
)

type Builtin func(*ast.CallExpression, *Environment) object.Object

type Evaluator struct {
	environment *Environment
	builtins    map[string]Builtin
}

func New() *Evaluator {
	environment := NewEnvironment()

	e := &Evaluator{environment: environment}

	builtins := make(map[string]Builtin)
	builtins["push"] = e.evaluatePush
	e.builtins = builtins

	return e
}

func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	return e.evaluate(node, e.environment)
}

// refactorings todo:
// single true/false and null

func (e *Evaluator) evaluate(node ast.Node, env *Environment) object.Object {
	switch v := node.(type) {
	case *ast.ExpressionStatement:
		return e.evaluate(v.Expression, env)
	case *ast.LetStatement:
		value := e.evaluate(v.Value, env)
		env.put(v.Name.Value, value)
		return nil
	case *ast.ReturnStatement:
		return &object.Return{ReturnValue: e.evaluate(v.ReturnValue, env)}
	case *ast.Program:
		var result object.Object
		for _, stmt := range v.Statements {
			result = e.evaluate(stmt, env)
		}
		return result
	case *ast.PrefixExpression:
		return e.evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return e.evaluateInfixExpression(v, env)
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
		builtin := e.builtins[v.TokenLiteral()]
		if builtin != nil {
			return builtin(v, env)
		}

		functionObj, ok := env.get(v.TokenLiteral()).(*object.Function)
		if !ok {
			return nil
		}

		newEnv := FromEnvironment(env)

		for i, param := range functionObj.Parameters {
			newEnv.put(param, e.evaluate(v.Arguments[i], env))
		}

		result := e.evaluateBlockStatement(functionObj.Body, newEnv)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt.ReturnValue
		}
		return result
	case *ast.IfExpression:
		evaluatedCondition := e.evaluate(v.Condition, env)

		booleanEvalResult, ok := evaluatedCondition.(*object.Boolean)
		if !ok {
			return nil
		}

		if booleanEvalResult.Value {
			return e.evaluateBlockStatement(v.Consequence, env)
		} else if v.Alternative != nil {
			return e.evaluateBlockStatement(v.Alternative, env)
		}

		// TODO return null object
		return nil
	case *ast.ArrayExpression:
		elements := make([]object.Object, len(v.Elements))
		for i, el := range v.Elements {
			elements[i] = e.evaluate(el, env)
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

func (e *Evaluator) evaluateBlockStatement(blockStmt *ast.BlockStatement, env *Environment) object.Object {
	var result object.Object

	for _, stmt := range blockStmt.Statements {
		result = e.evaluate(stmt, env)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt
		}
	}

	return result
}

func (e *Evaluator) evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) object.Object {
	expr := e.evaluate(prefixExpr.Right, env)
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

func (e *Evaluator) evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) object.Object {
	left := e.evaluate(infixExpr.Left, env)
	right := e.evaluate(infixExpr.Right, env)

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
		return newError("Operation not supported %s %s %s", left.Type(), infixExpr.Operator, right.Type())
	}
}

func (e *Evaluator) evaluatePush(call *ast.CallExpression, env *Environment) object.Object {
	if len(call.Arguments) != 2 {
		return nil
	}

	firstArg := e.evaluate(call.Arguments[0], env)
	arrArg, ok := firstArg.(*object.Array)
	if !ok {
		return nil
	}

	secondArg := e.evaluate(call.Arguments[1], env)

	return &object.Array{Elements: append(arrArg.Elements, secondArg)}
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
