package evaluator

import (
	"fmt"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObjects(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	}
	return NULL
}

func isError(obj object.Object) bool { return obj != nil && obj.Type() == object.ERROR_OBJ }

func evalProgram(p *ast.Program) object.Object {
	var ans object.Object
	for _, sm := range p.Statements {
		ans = Eval(sm)

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
		return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
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

func evalIfExpression(ie *ast.IfExpression) object.Object {
	con := Eval(ie.Condition)
	if isError(con) {
		return con
	}
	if isTruthy(con) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
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

func evalBlockStatement(b *ast.BlockStatement) object.Object {
	var ans object.Object

	for _, s := range b.Statements {
		ans = Eval(s)
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
