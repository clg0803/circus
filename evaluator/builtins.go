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
			default:
				return newError("arg to `len` not supported so far, got %s",
					args[0].Type())
			}
		},
	},
}
