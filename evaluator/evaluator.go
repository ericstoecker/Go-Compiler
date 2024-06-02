package evaluator

import (
	"compiler/ast"
	"compiler/object"
	"compiler/token"
)

type Builtin func(*ast.CallExpression, *Environment) object.Object

var NULL = &object.Null{}
var TRUE = &object.Boolean{Value: true}
var FALSE = &object.Boolean{Value: false}

type Evaluator struct {
	environment *Environment
}

func getBuiltinByName(name string) *object.Builtin {
	for _, builtin := range object.Builtins {
		if builtin.Name == name {
			return builtin
		}
	}

	return nil
}

func New() *Evaluator {
	environment := NewEnvironment()

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
		if isError(value) {
			return value
		}

		env.put(v.Name.Value, value)
		return NULL
	case *ast.ReturnStatement:
		value := evaluate(v.ReturnValue, env)
		if isError(value) {
			return value
		}

		return &object.Return{ReturnValue: value}
	case *ast.Program:
		var result object.Object
		for _, stmt := range v.Statements {
			result = evaluate(stmt, env)
			if isError(result) {
				return result
			}
		}
		return result
	case *ast.PrefixExpression:
		return evaluatePrefixExpression(v, env)
	case *ast.InfixExpression:
		return evaluateInfixExpression(v, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return newBool(v.Value)
	case *ast.StringLiteral:
		return &object.String{Value: v.Value}
	case *ast.Identifier:
		builtin := getBuiltinByName(v.Value)
		if builtin != nil {
			return builtin
		}

		variable := env.get(v.Value)
		if variable == nil {
			return object.NewError("undefined: %s", v.TokenLiteral())
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
		if isError(left) {
			return left
		}

		builtin, ok := left.(*object.Builtin)
		if ok {
			args := make([]object.Object, len(v.Arguments))
			for i, arg := range v.Arguments {
				currentArg := evaluate(arg, env)
				if isError(currentArg) {
					return currentArg
				}
				args[i] = currentArg
			}
			return wrapNativeValue(builtin.Fn(args...))
		}

		function, ok := left.(*object.Function)
		if !ok {
			return object.NewError("not a function. Got %s", left.Type())
		}

		if numArg, numParam := len(v.Arguments), len(function.Parameters); numArg != numParam {
			return object.NewError("wrong number of arguments: expected %d. Got %d", numParam, numArg)
		}

		newEnv := FromEnvironment(env)

		for i, param := range function.Parameters {
			currentParam := evaluate(v.Arguments[i], env)
			if isError(currentParam) {
				return currentParam
			}

			newEnv.put(param, currentParam)
		}

		result := evaluateBlockStatement(function.Body, newEnv)

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt.ReturnValue
		}
		return result
	case *ast.IfExpression:
		evaluatedCondition := evaluate(v.Condition, env)
		if isError(evaluatedCondition) {
			return evaluatedCondition
		}

		booleanEvalResult, ok := evaluatedCondition.(*object.Boolean)
		if !ok {
			return object.NewError("non-boolean condition in if-expression")
		}

		if booleanEvalResult.Value {
			return evaluateBlockStatement(v.Consequence, env)
		} else if v.Alternative != nil {
			return evaluateBlockStatement(v.Alternative, env)
		}

		return NULL
	case *ast.ArrayLiteral:
		elements := make([]object.Object, len(v.Elements))
		for i, element := range v.Elements {
			evaluatedElement := evaluate(element, env)
			if isError(evaluatedElement) {
				return evaluatedElement
			}
			elements[i] = evaluatedElement
		}
		return &object.Array{Elements: elements}
	case *ast.MapLiteral:
		entries := make(map[string]object.Object, len(v.Entries))
		for key, value := range v.Entries {
			keyObject := evaluate(key, env)
			if isError(keyObject) {
				return keyObject
			}
			hashableKey, ok := keyObject.(object.Hashable)
			if !ok {
				return object.NewError("not a valid hash key: %s", keyObject.Type())
			}

			valueObject := evaluate(value, env)
			if isError(valueObject) {
				return valueObject
			}

			entries[hashableKey.Hash()] = valueObject
		}
		return &object.Map{Entries: entries}
	case *ast.IndexExpression:
		left := evaluate(v.Left, env)
		if isError(left) {
			return left
		}

		index := evaluate(v.Index, env)
		if isError(index) {
			return index
		}

		switch true {
		case left.Type() == object.ARRAY && index.Type() == object.INT:
			arr := left.(*object.Array)
			indexInt := index.(*object.Integer)
			return evaluateArrayIndexExpression(arr, indexInt)
		case left.Type() == object.MAP:
			mapObj := left.(*object.Map)
			return evaluateMapIndexExpression(mapObj, index)
		default:
			return object.NewError("type missmatch: cannot index %s", left.Type())
		}
	default:
		return object.NewError("Node of type %T unknown", v)
	}
}

func evaluateArrayIndexExpression(left *object.Array, index *object.Integer) object.Object {
	if 0 > index.Value || index.Value >= int64(len(left.Elements)) {
		return object.NewError("index %d out of bounds for array of length %d", index.Value, len(left.Elements))
	}
	return left.Elements[index.Value]
}

func evaluateMapIndexExpression(left *object.Map, index object.Object) object.Object {
	hashableIndex, ok := index.(object.Hashable)
	if !ok {
		return object.NewError("not a valid hash key: %s", index.Type())
	}

	key := hashableIndex.Hash()

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
		if isError(result) {
			return result
		}

		returnStmt, ok := result.(*object.Return)
		if ok {
			return returnStmt
		}
	}

	return result
}

func evaluatePrefixExpression(prefixExpr *ast.PrefixExpression, env *Environment) object.Object {
	expr := evaluate(prefixExpr.Right, env)
	if isError(expr) {
		return expr
	}

	switch prefixExpr.Operator {
	case token.MINUS:
		intExpr, ok := expr.(*object.Integer)
		if !ok {
			return object.NewError("Operation not supported: -%s", expr.Type())
		}

		return &object.Integer{Value: -intExpr.Value}
	case token.BANG:
		boolExpr, ok := expr.(*object.Boolean)
		if !ok {
			return object.NewError("Operation not supported: !%s", expr.Type())
		}

		return newBool(!boolExpr.Value)
	default:
		return object.NewError("Prefix unknown: %s", prefixExpr.Operator)
	}
}

func evaluateInfixExpression(infixExpr *ast.InfixExpression, env *Environment) object.Object {
	left := evaluate(infixExpr.Left, env)
	if isError(left) {
		return left
	}

	right := evaluate(infixExpr.Right, env)
	if isError(right) {
		return right
	}

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
		return object.NewError("Operation not supported %s %s %s", left.Type(), infixExpr.Operator, right.Type())
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
		return object.NewError("Infix operator unknown: %s", operator)
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
		return object.NewError("Infix operator unknown: %s", operator)
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
		return object.NewError("Infix operator unknown: %s", operator)
	}
}

func newBool(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	switch obj.(type) {
	case *object.Error:
		return true
	default:
		return false
	}
}

func wrapNativeValue(value interface{}) object.Object {
	switch value := value.(type) {
	case bool:
		return newBool(value)
	case nil:
		return NULL
	case object.Object:
		return value
	default:
		panic("Trying to wrap native value that does not exist in evaluator")
	}
}
