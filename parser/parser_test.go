package parser

import (
	"fmt"
	"testing"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/lexer"
)

// parser/parser_test.go

// parser/parser_test.go

// parser/parser_test.go

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input        string
        operator     string
        integerValue int64
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
    }

    for _, tt := range prefixTests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
        }
        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s",
                tt.operator, exp.Operator)
        }
        if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
            return
        }
    }
}

func checkParserErrors(t *testing.T, p *Parser) {
	e := p.Errors()
	if len(e) == 0 {
		return
	} else {
		t.Errorf("parser has %d errors", len(e))
		for _, msg := range e {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, v int64) bool {
    iv, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not *ast.IntegerLiteral, got = %T", il)
        return false
    }
    if iv.Value != v {
        t.Errorf("iv not %d, got = %d", v, iv.Value)
        return false
    }
    if iv.TokenLiteral() != fmt.Sprintf("%d", v) {
        t.Errorf("iv.TokenLiteral not %d, got = %s", v, iv.TokenLiteral())
        return false
    }

    return true
}