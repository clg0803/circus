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
