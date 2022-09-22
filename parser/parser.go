package parser

import (
	"fmt"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/lexer"
	"github.com/clg0803/circus/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekErrors(t token.TokenType) {
	m := fmt.Sprintf("expected next token to be %s, got %s instead\n",
		t, p.curToken.Type)
	p.errors = append(p.errors, m)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		s := p.parseStatement()
		if s != nil {
			program.Statements = append(program.Statements, s)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{Token: p.curToken}

	if !p.exceptPeek(token.IDENT) {
		return nil
	}

	s.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.exceptPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 跳过对等号后面表达式的处理直至 `;`
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return s
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// peek at eat
func (p *Parser) exceptPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekErrors(t)
		return false
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: parse expression

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return s
}
