package evaluator

import (
	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObjects(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		Left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, Left, right)
	}
	return nil
}

func evalStatements(s []ast.Statement) object.Object {
	var ans object.Object
	for _, sm := range s {
		ans = Eval(sm)
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
		return NULL
	}
}

func evalInfixExpression(op string,
	left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ &&
		right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(op, left, right)
	default:
		return nil
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
		return &object.Boolean{Value: lv < rv}
	case ">":
		return &object.Boolean{Value: lv > rv}
	case "==":
		return &object.Boolean{Value: lv == rv}
	case "!=":
		return &object.Boolean{Value: lv != rv}
	default:
		return nil
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
		return NULL
	}
	v := right.(*object.Integer).Value
	return &object.Integer{Value: -v}
}
