// evaluator/evaluator_test.go

package evaluator

import (
	"testing"

	"github.com/clg0803/circus/lexer"
	"github.com/clg0803/circus/object"
	"github.com/clg0803/circus/parser"
)

// evaluator/evaluator_test.go

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"1", 1},
		{"-5", -5},
		{"-1", -1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		exp   bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!5", true},
		{"!!true", true},
		{"!!false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.exp)
	}
}
