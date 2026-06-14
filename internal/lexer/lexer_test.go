package lexer

import (
	"testing"

	"github.com/Daggam/CDL/internal/token"
)

func TestNextToken(t *testing.T) {
	input := `SEND OFFER messi FOR messi IN USER pepe;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SEND, "SEND"},
		{token.OFFER, "OFFER"},
		{token.IDENT, "messi"},
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
