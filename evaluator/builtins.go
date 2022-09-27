package evaluator

import (
	"github.com/clg0803/circus/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args, got %d, want = 1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("arg to `len` not supported so far, got %s",
					args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args, got %d, want = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("arguments to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args, got %d, want = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("arguments to `last` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			if l > 0 {
				return arr.Elements[l-1]
			}

			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of args, got %d, want = 1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("arguments to `rest` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			if l > 0 {
				ne := make([]object.Object, l-1)
				copy(ne, arr.Elements[1:l])
				return &object.Array{Elements: ne}
			}
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of args, got %d, want = 2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("arguments to `push` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			l := len(arr.Elements)
			ne := make([]object.Object, l + 1)
			copy(ne, arr.Elements)
			ne[l] = args[1]
			return &object.Array{Elements: ne}
		},
	},
}
