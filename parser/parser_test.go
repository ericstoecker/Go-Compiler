package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
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
	expectNotNil(t, program)

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

func TestReturnStatement(t *testing.T) {
	input := `
    return 10;
    return foobar;
    return;
    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	expectNotNil(t, program)

	if len(program.Statements) != 3 {
		t.Fatalf("Program has wrong numer of statements. Expected 3. Got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Not a return statement")
		}
	}
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

func TestIdentifier(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	statements := program.Statements

	if len(statements) != 1 {
		t.Fatalf("Program has wrong numer of statements. Expected 1. Got %d", len(program.Statements))
	}

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T", statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Identifier. Got %T", stmt.Expression)
	}

	if ident.Token.Type != token.IDENT {
		t.Fatalf("Expected token.Type to be IDENT. Got '%s'", ident.Token.Type)
	}

	if ident.Value != "foobar" {
		t.Fatalf("Expected ident.Value to be 'foobar'. Got '%s'", ident.Value)
	}
}

func expectNotNil(t *testing.T, program *ast.Program) {
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
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
