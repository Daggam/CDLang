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

	program := createProgram(t, input)

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

func TestGetOfferStatement(t *testing.T) {
	input := `
	GET OFFER messi;
	GET OFFER ronaldo;
	GET OFFER mbappe;
	`
	program := createProgram(t, input)

	tests := []struct {
		expectedIdentifier string
	}{
		{"messi"},
		{"ronaldo"},
		{"mbappe"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.TokenLiteral() != "GET" {
			t.Errorf("stmt.TokenLiteral no es un GET sino un %q", stmt.TokenLiteral())
		}
		getOfferStmt, ok := stmt.(*ast.GetOfferStatement)

		if !ok {
			t.Errorf("stmt no es del tipo GetOfferStatement, sino que es del tipo %T", stmt)
		}

		identifier := getOfferStmt.Identifier
		if identifier.Value != tt.expectedIdentifier {
			t.Errorf("El valor de los identificadores no son iguales.")
		}
	}
}

func createProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() retorno nulo.")
	}
	return program
}

func TestSendOffer(t *testing.T) {
	input := `
	SEND OFFER ronaldo FOR messi IN USER pepe;
	SEND OFFER ronaldo(20) FOR messi IN USER pepe;
	SEND OFFER ronaldo FOR wanchope(20) IN USER pepe;
	SEND OFFER mbappe, ronaldo FOR messi IN USER pepe;
	SEND OFFER mbappe, ronaldo(20) FOR messi IN USER pepe;
	SEND OFFER mbappe, ronaldo, wanchope FOR messi, neymar IN USER pepe;
	SEND OFFER mbappe, ronaldo, wanchope FOR messi, neymar(20) IN USER pepe;
	`

	program := createProgram(t, input)

	type Collectable struct {
		Value  string
		Amount int
	}

	tests := []struct {
		expectedLCollectables []Collectable
		expectedRCollectables []Collectable
		userValue             string
	}{
		{
			[]Collectable{{Value: "ronaldo", Amount: 1}},
			[]Collectable{{Value: "messi", Amount: 1}},
			"pepe",
		},
		{
			[]Collectable{{Value: "ronaldo", Amount: 20}},
			[]Collectable{{Value: "messi", Amount: 1}},
			"pepe",
		},
		{
			[]Collectable{{Value: "ronaldo", Amount: 1}},
			[]Collectable{{Value: "wanchope", Amount: 20}},
			"pepe",
		},
		{
			[]Collectable{{Value: "mbappe", Amount: 1}, {Value: "ronaldo", Amount: 1}},
			[]Collectable{{Value: "messi", Amount: 1}},
			"pepe",
		},
		{
			[]Collectable{{Value: "mbappe", Amount: 1}, {Value: "ronaldo", Amount: 20}},
			[]Collectable{{Value: "messi", Amount: 1}},
			"pepe",
		},
		{
			[]Collectable{{Value: "mbappe", Amount: 1}, {Value: "ronaldo", Amount: 1}, {Value: "wanchope", Amount: 1}},
			[]Collectable{{Value: "messi", Amount: 1}, {Value: "neymar", Amount: 1}},
			"pepe",
		},
		{
			[]Collectable{{Value: "mbappe", Amount: 1}, {Value: "ronaldo", Amount: 1}, {Value: "wanchope", Amount: 1}},
			[]Collectable{{Value: "messi", Amount: 1}, {Value: "neymar", Amount: 20}},
			"pepe",
		},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.TokenLiteral() != "SEND" {
			t.Errorf("stmt.TokenLiteral no es SEND sino %q", stmt.TokenLiteral())
		}

		sendOfferStatement, ok := stmt.(*ast.SendOfferStatement)
		if !ok {
			t.Errorf("stmt no es del tipo SendOfferStatement, sino que es del tipo %T", stmt)
		}

		//Checamos los LValues y RValues
		for i, LExpectedCollection := range tt.expectedLCollectables {

			if sendOfferStatement.LCollectables[i].Value != LExpectedCollection.Value {
				t.Errorf("Se esperaba como coleccionable a %q, pero se obtuvo %q", LExpectedCollection.Value, sendOfferStatement.LCollectables[i].Value)
			}
			if sendOfferStatement.LCollectables[i].Amount != LExpectedCollection.Amount {
				t.Errorf("Se esperaba como cantidad del coleccionable a %d, pero se obtuvo %d", LExpectedCollection.Amount, sendOfferStatement.LCollectables[i].Amount)
			}
		}
		for i, RExpectedCollection := range tt.expectedRCollectables {
			if sendOfferStatement.RCollectables[i].Value != RExpectedCollection.Value {
				t.Errorf("Se esperaba como coleccionable a %q, pero se obtuvo %q", RExpectedCollection.Value, sendOfferStatement.RCollectables[i].Value)
			}
			if sendOfferStatement.RCollectables[i].Amount != RExpectedCollection.Amount {
				t.Errorf("Se esperaba como cantidad del coleccionable a %d, pero se obtuvo %d", RExpectedCollection.Amount, sendOfferStatement.RCollectables[i].Amount)
			}
		}

		//Checamos al usuario
		if sendOfferStatement.Username.Value != tt.userValue {
			t.Errorf("Se esperaba a %q como valor de username, pero se obtuvo %q", tt.userValue, sendOfferStatement.Username.Value)
		}
	}
}

func TestViewOffer(t *testing.T) {
	input := `
	VIEW OFFER;
	`
	program := createProgram(t, input)

	if len(program.Statements) != 1 {
		t.Error("Se esperaba una unica sentencia.")
	}
	stmt := program.Statements[0]

	if stmt.TokenLiteral() != "VIEW" {
		t.Errorf("stmt.TokenLiteral no es VIEW, sino que es %q", stmt.TokenLiteral())
	}

	_, ok := stmt.(*ast.ViewOfferStatement)

	if !ok {
		t.Errorf("stmt no es del tipo ast.ViewOfferStatement, sino que es del tipo %T", stmt)
	}
}
