package parser

import (
	"testing"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/lexer"
)

func TestOfferStatement(t *testing.T) {
	input := `OFFER messi;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() retorno nulo.")
	}

	tests := []struct {
		expectedCollection string
	}{
		{"messi"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testOfferStatement(t, stmt, tt.expectedCollection) {
			return
		}
	}
}

func testOfferStatement(t *testing.T, s ast.Statement, collection string) bool {
	if s.TokenLiteral() != "OFFER" {
		t.Errorf("s.TokenLiteral no es 'OFFER'. se obtuvo=%q", s.TokenLiteral())
		return false
	}

	offerStmt, ok := s.(*ast.OfferStatement)

	if !ok {
		t.Errorf("s no es un *ast.OfferStatement, se obtuvo=%T", s)
	}

	if offerStmt.Collectable.Value != collection {
		t.Errorf("offerStmt.Collectable.Value no es %s, sino =%s", collection, offerStmt.Collectable.Value)
		return false
	}

	if offerStmt.Collectable.TokenLiteral() != collection {
		t.Errorf("offerStmt.Collectable.TokenLiteral() no es %s, sino =%s", collection, offerStmt.Collectable.TokenLiteral())
		return false
	}

	return true
}
