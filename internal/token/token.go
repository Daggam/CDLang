package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identificadores y literales
	IDENT = "IDENT"
	INT   = "INT"

	//Operadores
	ASSIGN = "="
	PLUS   = "+"

	//Delimitadores
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"

	//Keywords
	OFFER   = "OFFER"
	GET     = "GET"
	VIEW    = "VIEW"
	SEND    = "SEND"
	FOR     = "FOR"
	WHERE   = "WHERE"
	ACCEPT  = "ACCEPT"
	DECLINE = "DECLINE"
	TRADE   = "TRADE"
)
