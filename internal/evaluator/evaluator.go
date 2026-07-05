package evaluator

import (
	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Collectable:
		return &object.Collectable{Name: object.CollectableName(node.Value), Amount: node.Amount}
	}
	return nil
}
