package evaluator

import (
	"compiler/lexer"
	"compiler/parser"
	"testing"
)

func TestIntegerExpression(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{
			"10",
			10,
		},
		{
			"5",
			5,
		},
		{
			"-3",
			-3,
		},
		{
			"10 + 3",
			13,
		},
		{
			"-2 - -10",
			8,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult, ok := output.(*IntegerObject)
		if !ok {
			t.Fatalf("Expected IntegerObject. Got %T", output)
		}

		if intResult.Value != tt.expected {
			t.Fatalf("Expected %d. Got %d", tt.expected, intResult.Value)
		}
	}
}
