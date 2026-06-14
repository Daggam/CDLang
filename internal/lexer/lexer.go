package lexer

import "github.com/Daggam/CDL/internal/token"

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
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
