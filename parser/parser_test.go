package parser

import (
	"testing"

	"github.com/clg0803/circus/lexer"
	"github.com/clg0803/circus/ast"
)

// parser/parser_test.go

// parser/parser_test.go

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d",
            len(program.Statements))
    }
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
            program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
    }
    if literal.Value != 5 {
        t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
    }
    if literal.TokenLiteral() != "5" {
        t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
            literal.TokenLiteral())
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
