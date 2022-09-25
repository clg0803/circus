package evaluator

import (
	"fmt"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObjects(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body}
	case *ast.CallExpression:
		f := Eval(node.Function, env)
		if isError(f) {
			return f
		}
		args := evalExpression(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(f, args)
	}
	return NULL
}

func isError(obj object.Object) bool { return obj != nil && obj.Type() == object.ERROR_OBJ }

func applyFunction(fn object.Object,
	args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		eEnv := extendFunctionEnv(fn, args)
		eva := Eval(fn.Body, eEnv)
		return unwrapReturnValue(eva)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function,
	args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvirnment(fn.Env)
	for i, p := range fn.Parameters {
		env.Set(p.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok { // unwrap at last
		return returnVal.Value
	}
	return obj
}

func evalExpression(args []ast.Expression,
	env *object.Environment) []object.Object {
	var ans []object.Object

	for _, e := range args {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		ans = append(ans, evaluated)
	}

	return ans
}

func evalIdentifier(node *ast.Identifier,
	env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalProgram(p *ast.Program, env *object.Environment) object.Object {
	var ans object.Object
	for _, sm := range p.Statements {
		ans = Eval(sm, env)

		switch ans := ans.(type) {
		case *object.ReturnValue:
			return ans.Value
		case *object.Error:
			return ans
		}
	}
	return ans
}

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	NULL = &object.Null{}
)

func nativeBoolToBooleanObjects(in bool) *object.Boolean {
	if in {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", op, right.Type())
	}
}

func evalInfixExpression(op string,
	left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ &&
		right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	case left.Type() == object.STRING_OBJ &&
		right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(op, left, right)
	case op == "==":
		return nativeBoolToBooleanObjects(left == right)
	case op == "!=":
		return nativeBoolToBooleanObjects(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s % s %s", left.Type(), op, right.Type())
	default:
		return newError("unknown operator: %s % s %s", left.Type(), op, right.Type())
	}
}

func evalIntegerInfixExpression(op string,
	left object.Object, right object.Object) object.Object {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	case "<":
		return nativeBoolToBooleanObjects(lv < rv)
	case ">":
		return nativeBoolToBooleanObjects(lv > rv)
	case "==":
		return nativeBoolToBooleanObjects(lv == rv)
	case "!=":
		return nativeBoolToBooleanObjects(lv != rv)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string,
	left object.Object, right object.Object) object.Object {
	if op != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), op, right.Type())
	}
	lv := left.(*object.String).Value
	rv := right.(*object.String).Value
	return &object.String{Value: lv + rv}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	// reverse
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	v := right.(*object.Integer).Value
	return &object.Integer{Value: -v}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	con := Eval(ie.Condition, env)
	if isError(con) {
		return con
	}
	if isTruthy(con) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(con object.Object) bool {
	switch con {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalBlockStatement(b *ast.BlockStatement, env *object.Environment) object.Object {
	var ans object.Object

	for _, s := range b.Statements {
		ans = Eval(s, env)
		if ans != nil {
			if t := ans.Type(); t == object.RETURN_VALUE_OBJ || t == object.ERROR_OBJ {
				return ans
			}
		}
	}

	return ans
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
