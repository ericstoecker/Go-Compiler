package repl

import (
	"bufio"
	"compiler/evaluator"
	"compiler/parser"
	scannergenerator "compiler/scanner-generator"
	"compiler/token"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	generator := scannergenerator.NewScannerGenerator()
	dfa := generator.GenerateScanner(token.TokenClassifications)

	e := evaluator.New()
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		s := scannergenerator.New(line, dfa)
		// l := lexer.New(line)
		p := parser.New(s)

		program := p.ParseProgram()
		if len(p.Errors) != 0 {
			printParserErrors(out, p.Errors)
		}
		evaluationResult := e.Evaluate(program)

		if evaluationResult != nil {
			io.WriteString(out, evaluationResult.String())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
