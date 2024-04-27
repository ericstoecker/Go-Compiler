package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
    let x = 5;
    let y = 10;
    let foobar = 6934;
    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program has wrong number of statements. Expected 3. Got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. Got '%q'", s.TokenLiteral())
		return false
	}
	lst, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Not a let statement.")
	}

	if lst.Name.Value != name {
		t.Errorf("Expected lst.Name.Value to be %s. Got '%s'", name, lst.Name.Value)
	}

	if lst.Name.TokenLiteral() != name {
		t.Errorf("Expected lst.Name.TokenLiteral() to be %s. Got '%s'", name, lst.Name.TokenLiteral())
	}
	return true
}

func TestErrorHandling(t *testing.T) {
	input := `
    let x 5
    `
	l := lexer.New(input)
	p := New(l)

	p.ParseProgram()

	errors := p.Errors
	if errors == nil {
		t.Fatalf("Expected error. Got nil")
	}

	if len(errors) != 1 {
		t.Fatalf("Expected 1 error. Got %d", len(errors))
	}
}
