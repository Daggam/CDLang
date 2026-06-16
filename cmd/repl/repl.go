package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Daggam/CDL/internal/lexer"
	"github.com/Daggam/CDL/internal/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)

	for {
		fmt.Printf(PROMPT)
		line, err := reader.ReadString('\n')
		if err != nil {
			panic("Hubo un error no esperado")
		}

		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
