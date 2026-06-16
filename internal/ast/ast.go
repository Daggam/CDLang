package ast

import (
	"github.com/Daggam/CDL/internal/token"
)

// Una interfaz de nodo para nuestro AST
type Node interface {
	TokenLiteral() string // Utilizado unicamente para debugging y testing
}

// Un nodo de tipo Statement
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// El nodo raíz de cada AST que produzca nuestro parser.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// OFFER STATEMENT
type OfferStatement struct {
	Token       token.Token
	Collectable *Collectable //Esto hay que cambiarlo luego por *Collectables
}

func (ofs *OfferStatement) statementNode()       {}
func (ofs *OfferStatement) TokenLiteral() string { return ofs.Token.Literal }

type Collectable struct {
	Token token.Token // token IDENT
	Value string      // Nombre del coleccionable
	// Después agregamos la cantidad
}

func (c *Collectable) expressionNode()      {}
func (c *Collectable) TokenLiteral() string { return c.Token.Literal }
