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
	Token        token.Token
	Collectables []*Collectable //Esto hay que cambiarlo luego por *Collectables
}

func (ofs *OfferStatement) statementNode()       {}
func (ofs *OfferStatement) TokenLiteral() string { return ofs.Token.Literal }

type Collectable struct {
	Token  token.Token // token IDENT
	Value  string      // Nombre del coleccionable
	Amount int         // Después agregamos la cantidad
}

func (c *Collectable) expressionNode()      {}
func (c *Collectable) TokenLiteral() string { return c.Token.Literal }

// GET OFFER STATEMENT
type GetOfferStatement struct {
	Token      token.Token
	Identifier *Identifier
}

func (gofs *GetOfferStatement) statementNode()       {}
func (gofs *GetOfferStatement) TokenLiteral() string { return gofs.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// SEND OFFER STATEMENT
type SendOfferStatement struct {
	Token         token.Token
	LCollectables []*Collectable
	RCollectables []*Collectable
	Username      *Identifier
}

func (sos *SendOfferStatement) statementNode()       {}
func (sos *SendOfferStatement) TokenLiteral() string { return sos.Token.Literal }

//View Offer Statement

type ViewOfferStatement struct {
	Token token.Token
}

func (vos *ViewOfferStatement) statementNode()       {}
func (vos *ViewOfferStatement) TokenLiteral() string { return vos.Token.Literal }

// Accept Trade Statement
type AcceptTradeStatement struct {
	Token   token.Token
	OfferID []int
}

func (ats *AcceptTradeStatement) statementNode()       {}
func (ats *AcceptTradeStatement) TokenLiteral() string { return ats.Token.Literal }

// Decline Offer Statement

type DeclineTradeStatement struct {
	Token   token.Token
	OfferID []int
}

func (dos *DeclineTradeStatement) statementNode()       {}
func (dos *DeclineTradeStatement) TokenLiteral() string { return dos.Token.Literal }

// Delete Offer Statement

type DeleteOfferStatement struct {
	Token        token.Token
	Collectables []*Collectable
}

func (dos *DeleteOfferStatement) statementNode()       {}
func (dos *DeleteOfferStatement) TokenLiteral() string { return dos.Token.Literal }

// Explain Statement
type ExplainStatement struct {
	Token token.Token
	Inner Statement
}

func (es *ExplainStatement) statementNode()       {}
func (es *ExplainStatement) TokenLiteral() string { return es.Token.Literal }
