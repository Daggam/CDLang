package lexer

import (
	"github.com/Daggam/CDL/internal/token"
)

type Lexer struct {
	input        string
	position     int  //Posicion actual del input (apunta al caracter actual)
	readPosition int  // Posicion actual de lectura en el input (despues del caracter actual)
	ch           byte // el caracter actual
}

// Creamos el lexer y lo retornamos.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Inicializamos nuestro l.ch, l.readPosition, l.position
	return l
}

// función helper avanzar al proximo caracter.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// Me retorna el token del caracter actual y pasa a leer el siguiente.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()
	l.skipComments()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			//Veamos que tipo de token es
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// Retorna una nueva token segun el tipo y el caracter.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Retorna el literal del identificador. (Lee el identificador completo)
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Funcion auxiliar que me ayuda a decidir el formato del identificador.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Lee el numero.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Funcion auxiliar para detección de números enteros.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// Funcion auxiliar que saltea los espacios en blanco u otros caracteres especiales.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	for l.ch == '/' && l.input[l.readPosition] == '/' {
		for l.ch != '\n' {
			l.readChar()
		}
		//fmt.Printf("%v", l.ch)
		l.readChar() // Pues se encuentra con \n y necesita saltearselo.
	}
}
