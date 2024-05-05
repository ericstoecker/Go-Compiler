package evaluator

import (
	"compiler/lexer"
	"compiler/parser"
	"testing"
)

func TestSimpleExpressions(t *testing.T) {
	input := `
    10
    `

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	evaluator := New()

	output := evaluator.Evaluate(program)

	if output.String() != "10" {
		t.Fatalf("Expected 10. Got %s", output.String())
	}
}
