package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Daggam/CDL/internal/lexer"
	"github.com/Daggam/CDL/internal/parser"
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
		p := parser.New(l)

		p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
		}

		//io.WriteString(out, program.String())
		//io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
