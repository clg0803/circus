package parser

import (
	"testing"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 991 332;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram return nil")
	}
	
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got =  %d ",
			len(program.Statements))
	}

	// tests := []struct {
	// 	expectedIdentifier string
	// }{
	// 	{"x"},
	// 	{"y"},
	// 	{"foobar"},
	// }

	// for i, tt := range tests {
	// 	s := program.Statements[i]
	// 	if !testLetStatement(t, s, tt.expectedIdentifier) {
	// 		return
	// 	}
	// }

	for _, s := range program.Statements {
		rs, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("s not *ast.ReturnStatement, got %T", s)
			continue
		}
		if rs.TokenLiteral() != "return" {
			t.Errorf("returnStatement.Tokenliteral not 'return', got %q", 
			rs.TokenLiteral())
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

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let', got =%q", s.TokenLiteral())
		return false
	}

	ls, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got = %T", s)
		return false
	}
	if ls.Name.Value != name {
		t.Errorf("ls.Name.Value not '%s', got = %s", name, ls.Name.Value)
		return false
	}
	if ls.Name.TokenLiteral() != name {
		t.Errorf("ls.Name.TokenLiteral() not '%s', got = %s", name, ls.Name.TokenLiteral())
		return false
	}

	return true
}
