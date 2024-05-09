package evaluator

import (
	"compiler/lexer"
	"compiler/object"
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
		{
			"(10 + 2) * 2",
			24,
		},
		{
			`let x = 10
		    x
		    `,
			10,
		},
		{
			`let x = -2
		    x + 2
		    `,
			0,
		},
		{
			`let x = fn (a) { return a + 1 }
		       x(2)
		       `,
			3,
		},
		{
			`let x = 10
		       let y = fn(x) { return x }
		       y(3)
		       `,
			3,
		},
		{
			`let x = 10
		       let y = fn(x) { return x }
		       y(4)
		       x
		       `,
			10,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult, ok := output.(*object.IntegerObject)
		if !ok {
			t.Fatalf("Expected IntegerObject. Got %T", output)
		}

		if intResult.Value != tt.expected {
			t.Fatalf("Expected %d. Got %d", tt.expected, intResult.Value)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true",
			true,
		},
		{
			"!true",
			false,
		},
		{
			"true == false",
			false,
		},
		{
			"10 + 2 == 12",
			true,
		},
		{
			"true != false",
			true,
		},
		{
			"10 != 10",
			false,
		},
		{
			"10 < 5",
			false,
		},
		{
			"10 > 5",
			true,
		},
		{
			"5 <= 5",
			true,
		},
		{
			"4 >= 5",
			false,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult, ok := output.(*object.BooleanObject)
		if !ok {
			t.Fatalf("Expected BooleanObject. Got %T", output)
		}

		if intResult.Value != tt.expected {
			t.Fatalf("Expected %t. Got %t", tt.expected, intResult.Value)
		}
	}
}
