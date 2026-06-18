package parser

import (
	"testing"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/lexer"
)

func TestOfferStatement(t *testing.T) {
	input := `
	OFFER messi, cristiano(20), mbappe;
	OFFER cristiano;
	OFFER mbappe;`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() retorno nulo.")
	}

	tests := []struct {
		expectedCollections []string
	}{
		{[]string{"messi", "cristiano", "mbappe"}},
		{[]string{"cristiano"}},
		{[]string{"mbappe"}},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testOfferStatement(t, stmt, tt.expectedCollections) {
			return
		}
	}
}

func testOfferStatement(t *testing.T, s ast.Statement, collections []string) bool {
	if s.TokenLiteral() != "OFFER" {
		t.Errorf("s.TokenLiteral no es 'OFFER'. se obtuvo=%q", s.TokenLiteral())
		return false
	}

	offerStmt, ok := s.(*ast.OfferStatement)

	if !ok {
		t.Errorf("s no es un *ast.OfferStatement, se obtuvo=%T", s)
	}
	for i, collectable := range offerStmt.Collectables {
		if collectable.Value != collections[i] {
			t.Errorf("offerStmt.Collectable.Value no es %s, sino =%s", collections[i], collectable.Value)
			return false
		}

		if collectable.TokenLiteral() != collections[i] {
			t.Errorf("offerStmt.Collectable.TokenLiteral() no es %s, sino =%s", collections[i], collectable.TokenLiteral())
			return false
		}
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if len(errors) == 0 {
		return
	}

	t.Errorf("El parser tuvo %d errores.", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow()
}
