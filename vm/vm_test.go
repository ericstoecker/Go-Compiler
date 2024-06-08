package vm

import (
	"compiler/ast"
	"compiler/compiler"
	"compiler/lexer"
	"compiler/object"
	"compiler/parser"
	"fmt"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"2 - 1", 1},
		{"3 * 2", 6},
		{"6 / 2", 3},
		{"-5", -5},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"true == false", false},
		{"true == true", true},
		{"true != true", false},
		{"true != false", true},
		{"10 == 5", false},
		{"10 == 10", true},
		{"6 != 6", false},
		{"6 != 5", true},
		{"10 <= 3", false},
		{"10 <= 11", true},
		{"10 <= 10", true},
		{"10 < 10", false},
		{"10 < 15", true},
		{"3 >= 4", false},
		{"3 >= 3", true},
		{"3 > 3", false},
		{"3 > 2", true},
		{`"ab" == "ba"`, false},
		{`"ab" == "ab"`, true},
		{`"ab" != "ab"`, false},
		{`"ab" != "cc"`, true},
		{"!true", false},
		{"!false", true},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"string"`, "string"},
		{`"str" + "ing"`, "string"},
	}

	runVmTests(t, tests)
}

func TestIfElseExpression(t *testing.T) {
	tests := []vmTestCase{
		{`if (false) { 10 } 20`, 20},
		{`if (true) { 10 }`, 10},
		{`if (true) { 10 } else { 20 }`, 10},
		{`if (false) { 10 } else { 20 }`, 20},
		{`if (false) { 10 }`, NULL},
	}

	runVmTests(t, tests)
}

func TestLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{`let x = 10; x`, 10},
		{`let x = "test"; let y = x; y`, "test"},
		{`let var = 20; if (true) { var + 1 }`, 21},
	}

	runVmTests(t, tests)
}

func TestArrayExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`[1, 2, 3]`, []int{1, 2, 3}},
		{`[1 * 1, 4 - 2, 6 / 2]`, []int{1, 2, 3}},
		{`let var = [1, 2, 3]; var`, []int{1, 2, 3}},
	}

	runVmTests(t, tests)
}

func TestHashMapExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`{ "a": 2, "b": 3 }`, map[string]int{"STRING: a": 2, "STRING: b": 3}},
	}

	runVmTests(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`let x = [0, 1, 2]; x[0]`, 0},
		{`let x = {"a": 20}; x["a"]`, 20},
		{`let var = {"a"+"b": "test"}; var["ab"]`, "test"},
		{`let x = {"a": 20}; x["b"]`, NULL},
	}

	runVmTests(t, tests)
}

func TestFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`fn() { return 10 }()`, 10},
		{`fn() {}()`, NULL},
		{`let x = fn() { 10 + 5 }(); x`, 15},
		{`let x = fn() { 10 + 5 }; let y = fn() { x() + 5 }(); y`, 20},
		{`let x = fn() { 10 + 5 }; let y = fn() { x }; y()()`, 15},
		{`fn(a) { a }(10)`, 10},
		{
			`
            let sum = fn(a, b) {
                let c = a + b;
                c;
            };
            sum(1,2) + sum(-1, -2);
            `,
			0,
		},
	}

	runVmTests(t, tests)
}

func TestLocalVariables(t *testing.T) {
	tests := []vmTestCase{
		{
			`let x = fn() { let num = 10; num }
		    let num = 5;
		    x();`,
			10,
		},
		{
			`let x = fn() { let num = 10; num }
		    let num = 5;
		    x();
		    num;`,
			5,
		},
		{
			`
            let num = 5;
            let numTwo = 1;
            let x = fn() { let num = 10; num + numTwo }
            x();`,
			11,
		},
	}

	runVmTests(t, tests)
}

func TestBuiltins(t *testing.T) {
	tests := []vmTestCase{
		{
			`push([], 1)`,
			[]int{1},
		},
		{
			`len([1, 5])`,
			2,
		},
		{
			`isEmpty([])`,
			true,
		},
		{
			`isEmpty([3])`,
			false,
		},
	}

	runVmTests(t, tests)
}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
            let newClosure = fn(a) {
                fn() { a; };
            };
            let closure = newClosure(99);
            closure();
            `,
			expected: 99,
		},
	}

	runVmTests(t, tests)
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let countDown = fn(x) {
				if (x == 0) {
					return 0;
				} else {
					countDown(x - 1)
				}
			};
			countDown(1)
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper = fn() {
				let countDown = fn(x) {
					if (x == 0) {
						return 0;
					} else {
						countDown(x - 1);
					}
				}
				countDown(1);
			};
			wrapper();
			`,
			expected: 0,
		},
	}

	runVmTests(t, tests)
}

func TestErrorHandling(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			let x = fn(a, b) { return a + b };
			x(1, 2, 3);
			`,
			expected: fmt.Errorf("wrong number of arguments: expected 2, got 3"),
		},
		{
			input: `
			let x = [];
			x[0]
			`,
			expected: fmt.Errorf("index 0 out of bounds for length 0"),
		},
		{
			input: `
			if (10) {};
			`,
			expected: fmt.Errorf("type missmatch: expected BOOLEAN, got INT"),
		},
		{
			input: `
			10 + "test"
			`,
			expected: fmt.Errorf("type missmatch: INT, STRING"),
		},
		{
			input: `
			len([], [])
			`,
			expected: fmt.Errorf("wrong number of arguments: expected 1. Got 2"),
		},
	}

	testVmError(t, tests)
}

func testVmError(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()

		if err == nil {
			t.Fatalf("expected error")
		}

		expectedError := tt.expected.(error)
		if err.Error() != expectedError.Error() {
			t.Fatalf("wrong error: expected '%s', got '%s'", expectedError, err)
		}
	}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPopped()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case *object.Null:
		err := testNullObject(actual)
		if err != nil {
			t.Errorf("testNullObject failed: %s", err)
		}
	case []int:
		err := testArrayObject(expected, actual)
		if err != nil {
			t.Errorf("testArrayObject failed: %s", err)
		}
	case map[string]int:
		err := testMapObject(expected, actual)
		if err != nil {
			t.Errorf("testMapObject failed: %s", err)
		}
	default:
		t.Errorf("tests for type %T not supported", expected)
	}
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
	}

	return nil
}

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not String. got=%T (%v)",
			actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%s, want=%s",
			result.Value, expected)
	}

	return nil
}

func testNullObject(actual object.Object) error {
	_, ok := actual.(*object.Null)
	if !ok {
		return fmt.Errorf("object is not Null. got=%T (%v)",
			actual, actual)
	}

	return nil
}

func testArrayObject(expected []int, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got=%T (%v)",
			actual, actual)
	}

	for i, e := range result.Elements {
		elem, ok := e.(*object.Integer)
		if !ok {
			return fmt.Errorf("element is not int at index %d. got=%T (%v)", i, elem, elem)
		}

		if int(elem.Value) != expected[i] {
			return fmt.Errorf("element has wrong value at index %d. got=%d, want=%d", i, elem.Value, expected[i])
		}
	}

	return nil
}

func testMapObject(expected map[string]int, actual object.Object) error {
	result, ok := actual.(*object.Map)
	if !ok {
		return fmt.Errorf("object is not Map. got=%T (%v)",
			actual, actual)
	}

	for key, value := range result.Entries {
		expectedValue, ok := expected[key]
		if !ok {
			return fmt.Errorf("key does not exist. got=%s", key)
		}

		actualValue, ok := value.(*object.Integer)
		if !ok {
			return fmt.Errorf("value is not int for key %s. got=%T (%v)", key, value, value)
		}

		if int(actualValue.Value) != expectedValue {
			return fmt.Errorf("entry has wrong value for key %s. got=%d, want=%d", key, actualValue.Value, expectedValue)
		}
	}

	return nil
}
