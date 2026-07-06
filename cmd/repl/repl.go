package main

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Daggam/CDL/internal/evaluator"
	"github.com/Daggam/CDL/internal/lexer"
	"github.com/Daggam/CDL/internal/object"
	"github.com/Daggam/CDL/internal/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	reader := bufio.NewReader(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)
		line, err := reader.ReadString('\n')
		if err != nil {
			panic("Hubo un error no esperado")
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if errObj, ok := evaluated.(*object.Error); ok {
			io.WriteString(out, "[EVALUATOR] "+errObj.Message)
			io.WriteString(out, "\n")
			continue
		}
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
		//io.WriteString(out, program.String())
		//io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "[LEXER] "+msg+"\n")
	}
}
