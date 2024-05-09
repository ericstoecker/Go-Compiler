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
		{
			`let arr = [1,2,3]
            arr[0]
            `,
			1,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult, ok := output.(*object.Integer)
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
			"true && false",
			false,
		},
		{
			"false || true",
			true,
		},
		{
			"false || false",
			false,
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
		{
			`let x = fn(l) { if (l > 3) { return true } return false }
            x(5)
            `,
			true,
		},
		{
			`"ab" == "ab"`,
			true,
		},
		{
			`"ab" == "ba"`,
			false,
		},
		{
			`"ab" != "ba"`,
			true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator := New()

		output := evaluator.Evaluate(program)

		intResult, ok := output.(*object.Boolean)
		if !ok {
			t.Fatalf("Expected BooleanObject. Got %T", output)
		}

		if intResult.Value != tt.expected {
			t.Fatalf("Expected %t. Got %t", tt.expected, intResult.Value)
		}
	}
}

func TestStringEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`"test"`,
			"test",
		},
		{
			`"str" + "ing"`,
			"string",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator := New()

		output := evaluator.Evaluate(program)

		stringObj, ok := output.(*object.String)
		if !ok {
			t.Fatalf("Expected StringObject. Got %T", output)
		}

		if stringObj.Value != tt.expected {
			t.Fatalf("Expected %s. Got %s", tt.expected, stringObj.Value)
		}
	}

}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser had %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
