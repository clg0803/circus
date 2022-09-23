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

// parser/parser.go

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

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
	p.infixParseFns = make(map[token.TokenType]infixParseFn)

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerIdentifier)

	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)

	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.registerPrefix(token.LPAREN, p.parseGroupExpression)

	p.registerPrefix(token.IF, p.parseIfExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

	p.nextToken()
	s.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	s := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	s.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
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
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// KEY !!!
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
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

func (p *Parser) parsePrefixExpression() ast.Expression {
	e := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	e.Right = p.parseExpression(PREFIX)
	return e
}

func (p *Parser) parseInfixExpression(l ast.Expression) ast.Expression {
	e := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     l,
	}

	precedence := p.curPrecedence()
	p.nextToken() // eat op
	e.Right = p.parseExpression(precedence)

	return e
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.exceptPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	e := &ast.IfExpression{Token: p.curToken}
	if !p.exceptPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	e.Condition = p.parseExpression(LOWEST)

	if !p.exceptPeek(token.RPAREN) {
		return nil
	}
	if !p.exceptPeek(token.LBRACE) {
		return nil
	}

	e.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.exceptPeek(token.LBRACE) {
			return nil
		}

		e.Alternative = p.parseBlockStatement()
	}

	return e
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	b := &ast.BlockStatement{Token: p.curToken}
	b.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		s := p.parseStatement()
		b.Statements = append(b.Statements, s)
		p.nextToken()
	}

	return b
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	l := &ast.FunctionLiteral{Token: p.curToken}
	if !p.exceptPeek(token.LPAREN) {
		return nil
	}
	l.Parameters = p.parseFunctionParameters()

	if !p.exceptPeek(token.LBRACE) {
		return nil
	}
	l.Body = p.parseBlockStatement()
	return l
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	i := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return i
	}
	p.nextToken()

	ii := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	i = append(i, ii)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ii := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		i = append(i, ii)
	}

	if !p.exceptPeek(token.RPAREN) {
		return nil
	}

	return i
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	e := &ast.CallExpression{Token: p.curToken, Function: f}
	e.Arguments = p.parseCallArguments()
	return e
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}
	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.exceptPeek(token.RPAREN) {
		return nil
	}
	return args
}
