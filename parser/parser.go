package parser

import (
	"fmt"
	"strconv"

	"github.com/clg0803/circus/ast"
	"github.com/clg0803/circus/lexer"
	"github.com/clg0803/circus/token"
)

// PRIORITY
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         //+
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // func call

)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // 中缀表达式接收已解析的左侧的内容
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerIdentifier)

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
		program.Statements = append(program.Statements, s)
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
		return p.parseExpressionStatement()
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

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	s := &ast.ExpressionStatement{Token: p.curToken}
	s.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return s
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

// ---------------------------------------
// binded func TOKEN.TYPE refers to
// ---------------------------------------

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerIdentifier() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as Integer \n", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = v

	return lit
}
