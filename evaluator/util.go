package evaluator

import "compiler/object"

func constructBuiltins() map[string]object.Object {
	builtins := make(map[string]object.Object)
	builtins["push"] = &object.Builtin{
		Name: "push",
		Fn: func(args ...object.Object) object.Object {
			if numArg := len(args); numArg != 2 {
				return newError("wrong number of arguments: expected 2. Got %d", numArg)
			}

			arrArg, ok := args[0].(*object.Array)
			if !ok {
				return newError(
					"Operation not supported: push(%s, %s) (type missmatch, expected ARRAY. Got %s)",
					args[0].Type(), args[1].Type(), args[0].Type(),
				)
			}

			return &object.Array{Elements: append(arrArg.Elements, args[1])}
		},
	}
	builtins["len"] = &object.Builtin{
		Name: "len",
		Fn: func(args ...object.Object) object.Object {
			if numArg := len(args); numArg != 1 {
				return newError("wrong number of arguments: expected 1. Got %d", numArg)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("Operation not supported: len(%s)", arg.Type())
			}
		},
	}
	builtins["isEmpty"] = &object.Builtin{
		Name: "isEmpty",
		Fn: func(args ...object.Object) object.Object {
			if numArg := len(args); numArg != 1 {
				return newError("wrong number of arguments: expected 1. Got %d", numArg)
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return newBool(len(arg.Elements) == 0)
			case *object.String:
				return newBool(len(arg.Value) == 0)
			case *object.Map:
				return newBool(len(arg.Entries) == 0)
			default:
				return newError("Operation not supported: isEmpty(%s)", arg.Type())
			}
		},
	}

	return builtins
}
