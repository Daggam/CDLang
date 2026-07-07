package lexer

import (
	"testing"

	"github.com/Daggam/CDLang/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `SEND OFFER AR-LM10 FOR messi IN USER pepe;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SEND, "SEND"},
		{token.OFFER, "OFFER"},
		{token.IDENT, "AR-LM10"},
		{token.FOR, "FOR"},
		{token.IDENT, "messi"},
		{token.IN, "IN"},
		{token.USER, "USER"},
		{token.IDENT, "pepe"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestExplainToken(t *testing.T) {
	input := `EXPLAIN ACCEPT TRADE 58;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.EXPLAIN, "EXPLAIN"},
		{token.ACCEPT, "ACCEPT"},
		{token.TRADE, "TRADE"},
		{token.INT, "58"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
