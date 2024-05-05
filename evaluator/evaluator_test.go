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
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult := output.(*IntegerObject)
		if intResult.Value != tt.expected {
			t.Fatalf("Expected %d. Got %d", tt.expected, intResult.Value)
		}
	}
}
