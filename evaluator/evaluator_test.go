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
		{
			`len("abc")`,
			3,
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
			"true && true",
			true,
		},
		{
			`let x = true
            x && true`,
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
		{
			`isEmpty([])`,
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
		{
			`let x = ["a"]
		    let y = push(x, "b")
		    y[1]
		    `,
			"b",
		},
		{
			`let hashmap = { "a": "b" }
		    hashmap["a"]`,
			"b",
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		error string
	}{
		{
			`true - 10`,
			"Operation not supported BOOLEAN - INT",
		},
		{
			`10 - true
		    20
		    `,
			"Operation not supported INT - BOOLEAN",
		},
		{
			`a(10)`,
			"undefined: a",
		},
		{
			`if (300) { }`,
			"non-boolean condition in if-expression",
		},
		{
			`-false`,
			"Operation not supported: -false (type missmatch, expected INT. Got BOOLEAN)",
		},
		{
			`!30`,
			"Operation not supported: !30 (type missmatch, expected BOOLEAN. Got INT)",
		},
		{
			`let l = 10
            l[10]`,
			"Operation not supported: l[10] (must be either ARRAY or MAP)",
		},
		{
			`let l = [0, 1]
		    l[2]`,
			"index 2 out of bounds for array of length 2",
		},
		{
			`let x = true
		    push(x, false)`,
			"Operation not supported: push(BOOLEAN, BOOLEAN) (type missmatch, expected ARRAY. Got BOOLEAN)",
		},
		{
			`push(1, 2, 3)`,
			"wrong number of arguments: expected 2. Got 3",
		},
		{
			`let f = fn(a,b,c) { return a + b + c }
		    f(1,2)`,
			"wrong number of arguments: expected 3. Got 2",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		evaluator := New()

		output := evaluator.Evaluate(program)

		err, ok := output.(*object.Error)
		if !ok {
			t.Fatalf("Expected ErrorObject. Got %T", output)
		}

		if err.Message != tt.error {
			t.Fatalf("Expected error message to be: '%s'. Got '%s'", tt.error, err.Message)
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
