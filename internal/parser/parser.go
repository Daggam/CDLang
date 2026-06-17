package parser

import (
	"fmt"

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
	default:
		return nil
	}
}

func (p *Parser) parseOfferStatement() *ast.OfferStatement {
	stmt := &ast.OfferStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Collectable = &ast.Collectable{Token: p.curToken, Value: p.curToken.Literal}

	//Saltea las expresiones hasta que encuentra un punto y coma
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
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
