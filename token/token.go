package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 标识符 + 字面量
	IDENT = "IDENT" // add foo x y ...
	INT   = "INT"   // 1234

	// binary operator
	ASSIGN = "="
	PLUS   = "+"

	// split
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type Token struct {
	Type    TokenType // 上面定义的常量枚举
	Literal string
}
