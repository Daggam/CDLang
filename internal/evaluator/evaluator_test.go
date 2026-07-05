package evaluator

import (
	"testing"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/object"
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
