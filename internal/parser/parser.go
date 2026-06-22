package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Daggam/CDL/internal/ast"
	"github.com/Daggam/CDL/internal/lexer"
	"github.com/Daggam/CDL/internal/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	//Leemos dos tokens, así curToken y peekToken sean seteados correctamente
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("se esperaba que la próxima token sea %s y en su lugar fue %s.",
		t, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}

// Avanzamos al próximo token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Se encargará de parsear el programa.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.OFFER:
		return p.parseOfferStatement()
	case token.GET:
		return p.parseGetOfferStatement()
	case token.SEND:
		return p.parseSendOfferStatement()
	case token.VIEW:
		return p.parseViewOfferStatement()
	case token.ACCEPT:
		return p.parseAcceptTradeStatement()
	case token.DELETE:
		return p.parseDeleteOfferStatement()
	default:
		return nil
	}
}

//REGION DE PARSERS

func (p *Parser) parseOfferStatement() *ast.OfferStatement {
	stmt := &ast.OfferStatement{Token: p.curToken}

	collectables, err := p.parseCollectables()

	if err != nil {
		return nil
	}

	stmt.Collectables = collectables

	//Saltea las expresiones hasta que encuentra un punto y coma
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseGetOfferStatement() *ast.GetOfferStatement {
	stmt := &ast.GetOfferStatement{Token: p.curToken}
	if !p.expectPeek(token.OFFER) {
		return nil
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Identifier = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseSendOfferStatement() *ast.SendOfferStatement {
	stmt := &ast.SendOfferStatement{Token: p.curToken}

	if !p.expectPeek(token.OFFER) {
		return nil
	}
	//Catcheamos los LValues
	LCollectables, err := p.parseCollectables()
	if err != nil {
		return nil
	}
	stmt.LCollectables = LCollectables

	if !p.expectPeek(token.FOR) {
		return nil
	}
	//Catcheamos los RValues
	RCollectables, err := p.parseCollectables()
	if err != nil {
		return nil
	}
	stmt.RCollectables = RCollectables

	//Catcheamos a los usuarios
	if !p.expectPeek(token.IN) {
		return nil
	}
	if !p.expectPeek(token.USER) {
		return nil
	}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Username = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseViewOfferStatement() *ast.ViewOfferStatement {
	stmt := &ast.ViewOfferStatement{Token: p.curToken}

	if !p.expectPeek(token.OFFER) {
		return nil
	}

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseAcceptTradeStatement() *ast.AcceptTradeStatement {
	stmt := &ast.AcceptTradeStatement{Token: p.curToken}
	if !p.expectPeek(token.TRADE) {
		return nil
	}

	offerId := []int{}

	//Opcional: Si detecta un número/s se los agrega.
	if p.peekTokenIs(token.INT) {
		p.nextToken()
		firstOID, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return nil
		}
		offerId = append(offerId, firstOID)
		//Opcional: Si detecta comas, puede haber más de un entero.
		for p.peekTokenIs(token.COMMA) {
			p.nextToken()

			if !p.expectPeek(token.INT) {
				return nil
			}

			if oID, err := strconv.Atoi(p.curToken.Literal); err == nil {
				offerId = append(offerId, oID)
			} else {
				return nil
			}
		}
	}

	stmt.OfferID = offerId

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseDeleteOfferStatement() *ast.DeleteOfferStatement {
	stmt := &ast.DeleteOfferStatement{Token: p.curToken}
	if !p.expectPeek(token.OFFER) {
		return nil
	}
	collectables, err := p.parseCollectables()

	if err != nil {
		return nil
	}
	stmt.Collectables = collectables

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}
	return stmt
}

// ¿El token actual es del tipo t?
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// ¿El siguiente token es del tipo t?
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Funcion que avanzará al siguiente token si este es del tipo t
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseCollectable() (*ast.Collectable, error) {
	if !p.expectPeek(token.IDENT) {
		return nil, errors.New("Se esperaba la token IDENT")
	}

	collectable := &ast.Collectable{Token: p.curToken, Value: p.curToken.Literal, Amount: 1}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		if !p.expectPeek(token.INT) {
			return collectable, errors.New("Se esperaba la token INT")
		}
		value, err := strconv.Atoi(p.curToken.Literal)
		if err != nil {
			return collectable, err
		}
		collectable.Amount = value
		if !p.expectPeek(token.RPAREN) {
			return collectable, errors.New("Se esperaba la token ')'")
		}
	}

	return collectable, nil
}

func (p *Parser) parseCollectables() ([]*ast.Collectable, error) {
	collectables := []*ast.Collectable{}

	firstCollectable, err := p.parseCollectable()
	if err != nil {
		return nil, errors.New("No pudo obtenerse el primer coleccionable.")
	}

	collectables = append(collectables, firstCollectable)

	//COMO OPCIONAL, NOS DEBERÍA PERMITIR AGREGAR VARIAS COLECCIONES (por eso es parseCollectables)
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		collectable, err := p.parseCollectable()

		if err != nil {
			return collectables, errors.New("Hubo un error en el parseo de los coleccionables.")
		}

		collectables = append(collectables, collectable)
	}

	return collectables, nil
}
