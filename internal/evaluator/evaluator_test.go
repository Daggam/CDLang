package evaluator

import (
	"testing"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/lexer"
	"github.com/Daggam/CDL/internal/object"
	"github.com/Daggam/CDL/internal/parser"
	"github.com/Daggam/CDL/internal/token"
)

func TestEvalCollectables(t *testing.T) {
	token := token.Token{Type: token.IDENT, Literal: object.ES_AM7}
	tests := []struct {
		input    *ast.Collectable
		expected *object.Collectable
	}{
		{&ast.Collectable{Token: token, Value: token.Literal, Amount: 20}, &object.Collectable{Name: object.CollectableName(token.Literal), Amount: 20}},
	}

	env := object.NewEnvironment()
	for _, test := range tests {
		value := Eval(test.input, env)
		if value.Type() != object.COLLECTABLE_OBJ {
			t.Error("El objeto no es del tipo COLLECTABLE")
		}
		valueCollectable, ok := value.(*object.Collectable)
		if !ok {
			t.Errorf("Se esperaba que value sea del tipo Collectable pero es del tipo %T", value)
		}
		if valueCollectable.Name != test.expected.Name {
			t.Errorf("Se esperaba que valueCollectable tenga el nombre %s pero es %s", test.expected.Name, valueCollectable.Name)
		}
		if valueCollectable.Amount != test.expected.Amount {
			t.Errorf("Se esperaba que valueCollectable tenga el valor %d pero es %d", test.expected.Amount, valueCollectable.Amount)

		}
	}
}

func TestEvalOfferStatement(t *testing.T) {
	input := `OFFER AR-LM10;`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			t.Error(e)
			t.FailNow()
		}
	}

	env := object.NewEnvironment()

	Eval(program, env)

	//Corroboremos la base de datos. se tendría que haber agregado 1 messi y quitado otro messi.
	collectables := env.GetCollectables()
	for _, c := range collectables {
		if c.Name == object.AR_LM10 && c.Amount != 4 {
			t.Errorf("Hay %d AR-LM10 cuando se esperaban 4", c.Amount)
		}
	}
}

func TestErrorOfferStatements(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"OFFER AR-LM;", "unknown collectable: El coleccionable AR-LM no existe."},
		{"OFFER AR-LM10(29);", "no stock: No tienes suficiente coleccionables AR-LM10 para ofrecer. (Tienes 5)"},
	}
	for _, test := range tests {
		input := test.input
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			for _, e := range p.Errors() {
				t.Error(e)
				t.FailNow()
			}
		}

		env := object.NewEnvironment()
		evaluated := Eval(program, env)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No se obtuvo el error, sino %T(%+v)", evaluated, evaluated)
		}

		if errObj.Message != test.expectedMessage {
			t.Errorf("Se esperaba un mensaje de error: %q\nSin embargo se tiene: %q", test.expectedMessage, errObj.Message)
		}

	}
}
