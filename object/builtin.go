package object

var Builtins = []*Builtin{
	{
		Name: "push",
		Fn: func(args ...Object) interface{} {
			if numArg := len(args); numArg != 2 {
				return NewError("wrong number of arguments: expected 2. Got %d", numArg)
			}

			arrArg, ok := args[0].(*Array)
			if !ok {
				return NewError("type missmatch: first argument of push must be %s. Got %s", ARRAY, args[0].Type())
			}

			return &Array{Elements: append(arrArg.Elements, args[1])}
		},
	},
	{
		Name: "len",
		Fn: func(args ...Object) interface{} {
			if numArg := len(args); numArg != 1 {
				return NewError("wrong number of arguments: expected 1. Got %d", numArg)
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return NewError("type missmatch: len(%s) not supported", arg.Type())
			}
		},
	},
	{
		Name: "isEmpty",
		Fn: func(args ...Object) interface{} {
			if numArg := len(args); numArg != 1 {
				return NewError("wrong number of arguments: expected 1. Got %d", numArg)
			}

			switch arg := args[0].(type) {
			case *Array:
				return len(arg.Elements) == 0
			case *String:
				return len(arg.Value) == 0
			case *Map:
				return len(arg.Entries) == 0
			default:
				return NewError("type missmatch: isEmpty(%s) not supported", arg.Type())
			}
		},
	},
}
