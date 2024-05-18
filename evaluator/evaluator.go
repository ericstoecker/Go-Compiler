package evaluator

import (
	"compiler/ast"
	"compiler/object"
	"compiler/token"
	"fmt"
)

type Builtin func(*ast.CallExpression, *Environment) object.Object

var NULL = &object.Null{}
var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}

type Evaluator struct {
	environment *Environment
	builtins    map[string]Builtin
}

func New() *Evaluator {
	environment := FromMap(constructBuiltins())

	e := &Evaluator{environment: environment}

	return e
}

func (e *Evaluator) Evaluate(node ast.Node) object.Object {
	return evaluate(node, e.environment)
}

func evaluate(node ast.Node, env *Environment) object.Object {
	switch v := node.(type) {
	case *ast.ExpressionStatement:
		return evaluate(v.Expression, env)
	case *ast.LetStatement:
		value := evaluate(v.Value, env)
		env.put(v.Name.Value, value)
		return NULL
	case *ast.ReturnStatement:
		return &object.Return{ReturnValue: evaluate(v.ReturnValue, env)}
	case *ast.Program:
		var result object.Object
		for _, stmt := range v.Statements {
			result = evaluate(stmt, env)

			switch r := result.(type) {
			case *object.Error:
				return r
			}
		}
		return result
	case *ast.PrefixExpression:
		return evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return evaluateInfixExpression(v, env)
	case *ast.IntegerExpression:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanExpression:
		return newBool(v.Value)
	case *ast.StringExpression:
		return &object.String{Value: v.Value}
	case *ast.Identifier:
		variable := env.get(v.Value)
		if variable == nil {
			return newError("undefined: %s", v.TokenLiteral())
		}
		return variable
	case *ast.FunctionLiteral:
		params := make([]string, 0)
		for _, param := range v.Parameters {
			params = append(params, param.Value)
		}
		return &object.Function{Parameters: params, Body: v.Body}
	case *ast.CallExpression:
		left := evaluate(v.Left, env)
		if left.Type() == object.ERROR {
			return left
		}

		builtin, ok := left.(*object.Builtin)
		if ok {
			args := make([]object.Object, len(v.Arguments))
			for i, arg := range v.Arguments {
				args[i] = evaluate(arg, env)
			}
			return builtin.Fn(args...)
		}

		function, ok := left.(*object.Function)
		if !ok {
			return newError("not a function. Got %s", left.Type())
		}

		if numArg, numParam := len(v.Arguments), len(function.Parameters); numArg != numParam {
			return newError("wrong number of arguments: expected %d. Got %d", numParam, numArg)
		}

		newEnv := FromEnvironment(env)

		for i, param := range function.Parameters {
			newEnv.put(param, evaluate(v.Arguments[i], env))
		}

		result := evaluateBlockStatement(function.Body, newEnv)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt.ReturnValue
		}
		return result
	case *ast.IfExpression:
		evaluatedCondition := evaluate(v.Condition, env)

		booleanEvalResult, ok := evaluatedCondition.(*object.Boolean)
		if !ok {
			return newError("non-boolean condition in if-expression")
		}

		if booleanEvalResult.Value {
			return evaluateBlockStatement(v.Consequence, env)
		} else if v.Alternative != nil {
			return evaluateBlockStatement(v.Alternative, env)
		}

		return NULL
	case *ast.ArrayExpression:
		elements := make([]object.Object, len(v.Elements))
		for i, el := range v.Elements {
			elements[i] = evaluate(el, env)
		}
		return &object.Array{Elements: elements}
	case *ast.MapExpression:
		entries := make(map[string]object.Object, len(v.Entries))
		for key, value := range v.Entries {
			keyObject := evaluate(key, env)
			valueObject := evaluate(value, env)

			entries[keyObject.String()] = valueObject
		}
		return &object.Map{Entries: entries}
	case *ast.IndexExpression:
		left := evaluate(v.Left, env)
		index := evaluate(v.Index, env)
		switch true {
		case left.Type() == object.ARRAY && index.Type() == object.INT:
			arr := left.(*object.Array)
			indexInt := index.(*object.Integer)
			return evaluateArrayIndexExpression(arr, indexInt)
		case left.Type() == object.MAP:
			mapObj := left.(*object.Map)
			return evaluateMapIndexExpression(mapObj, index)
		default:
			return newError("Operation not supported: %s (must be either ARRAY or MAP)", v.String())
		}
	default:
		return newError("Node of type %T unknown", v)
	}
}

func evaluateArrayIndexExpression(left *object.Array, index *object.Integer) object.Object {
	if 0 > index.Value || index.Value >= int64(len(left.Elements)) {
		return newError("index %d out of bounds for array of length %d", index.Value, len(left.Elements))
	}
	return left.Elements[index.Value]
}

func evaluateMapIndexExpression(left *object.Map, index object.Object) object.Object {
	key := index.String()

	value, ok := left.Entries[key]
	if !ok {
		return NULL
	}
	return value
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
			return newError("Operation not supported: -%s (type missmatch, expected INT. Got %s)", expr.String(), expr.Type())
		}

		return &object.Integer{Value: -intExpr.Value}
	case token.BANG:
		boolExpr, ok := expr.(*object.Boolean)
		if !ok {
			return newError("Operation not supported: !%s (type missmatch, expected BOOLEAN. Got %s)", expr.String(), expr.Type())
		}

		return newBool(!boolExpr.Value)
	default:
		return newError("Prefix unknown: %s", prefixExpr.Operator)
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
		return newError("Operation not supported %s %s %s", left.Type(), infixExpr.Operator, right.Type())
	}
}

func (e *Evaluator) evaluatePush(call *ast.CallExpression, env *Environment) object.Object {
	if numArg := len(call.Arguments); numArg != 2 {
		return newError("wrong number of arguments: expected 2. Got %d", numArg)
	}

	firstArg := evaluate(call.Arguments[0], env)
	arrArg, ok := firstArg.(*object.Array)
	if !ok {
		return newError("Operation not supported: %s (type missmatch, expected ARRAY. Got %s)", call.String(), firstArg.Type())
	}

	secondArg := evaluate(call.Arguments[1], env)

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
		return newBool(left.Value == right.Value)
	case token.NOT_EQUALS:
		return newBool(left.Value != right.Value)
	case token.GT:
		return newBool(left.Value > right.Value)
	case token.LT:
		return newBool(left.Value < right.Value)
	case token.GREATER_EQUAL:
		return newBool(left.Value >= right.Value)
	case token.LESS_EQUAL:
		return newBool(left.Value <= right.Value)
	default:
		return newError("Infix operator unknown: %s", operator)
	}
}

func evaluateBooleanInfixExpression(operator token.TokenType, left *object.Boolean, right *object.Boolean) object.Object {
	switch operator {
	case token.EQUALS:
		return newBool(left == right)
	case token.NOT_EQUALS:
		return newBool(left != right)
	case token.AND:
		return newBool(left == TRUE && right == TRUE)
	case token.OR:
		return newBool(left == TRUE || right == TRUE)
	default:
		return newError("Infix operator unknown: %s", operator)
	}
}

func evaluateStringInfixExpression(operator token.TokenType, left *object.String, right *object.String) object.Object {
	switch operator {
	case token.EQUALS:
		return newBool(left.Value == right.Value)
	case token.NOT_EQUALS:
		return newBool(left.Value != right.Value)
	case token.PLUS:
		return &object.String{Value: left.Value + right.Value}
	default:
		return newError("Infix operator unknown: %s", operator)
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func newBool(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}
