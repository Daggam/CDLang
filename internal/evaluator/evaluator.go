package evaluator

import (
	"fmt"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.Collectable:
		return &object.Collectable{Name: object.CollectableName(node.Value), Amount: node.Amount}
	case *ast.OfferStatement:
		queryCollectables := []*object.Collectable{}
		for _, c := range node.Collectables {
			c_eval, ok := Eval(c, env).(*object.Collectable)
			if !ok {
				return nil
			}
			queryCollectables = append(queryCollectables, c_eval)
		}
		err := env.SetExchangeableCollection(queryCollectables)
		if err != nil {
			return newError("%s", err.Error())
		}
	}
	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
