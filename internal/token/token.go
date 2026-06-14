package token

import (
	"strings"
)

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
	IN      = "IN"
	USER    = "USER"
)

var keywords = map[string]TokenType{
	"offer":   OFFER,
	"get":     GET,
	"view":    VIEW,
	"send":    SEND,
	"for":     FOR,
	"where":   WHERE,
	"accept":  ACCEPT,
	"decline": DECLINE,
	"trade":   TRADE,
	"in":      IN,
	"user":    USER,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[strings.ToLower(ident)]; ok {
		return tok
	}
	return IDENT
}
