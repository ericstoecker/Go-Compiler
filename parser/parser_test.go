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
	checkParserErrors(t, p)

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
	checkParserErrors(t, p)

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
	checkParserErrors(t, p)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

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

func TestIntegerLiteralExpression(t *testing.T) {
	input := "10;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T", statements[0])
	}

	integer, ok := stmt.Expression.(*ast.IntegerExpression)
	if !ok {
		t.Fatalf("Expected IntegerExpression. Got %T", stmt.Expression)
	}

	if integer.Token.Type != token.INT {
		t.Fatalf("Expected token.Type to be INT. Got '%s'", integer.Token.Type)
	}

	if integer.Value != 10 {
		t.Fatalf("Expected integer.Value to be '10'. Got '%d'", integer.Value)
	}
}

// extend this
func TestPrefixExpressions(t *testing.T) {
	input := `
    -5;
    `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T", statements[0])
	}

	prefixExpr, ok := stmt.Expression.(*ast.PrefixExpression)
	if !ok {
		t.Fatalf("Expected PrefixExpression. Got %T", statements[0])
	}

	operator := prefixExpr.Operator
	if operator != token.MINUS {
		t.Fatalf("Exprected token.MINUS as operator. Got %s", operator)
	}

	rightExpr, ok := prefixExpr.Right.(*ast.IntegerExpression)
	if !ok {
		t.Fatalf("Expected IntegerExpression. Got %T", prefixExpr.Right)
	}

	if rightExpr.Value != 5 {
		t.Fatalf("Expected Value to be 5. Got %d", rightExpr.Value)
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + 5",
			"(5 + 5)",
		},
		{
			"5 + 5 + 5",
			"((5 + 5) + 5)",
		},
		{
			"5 - 5",
			"(5 - 5)",
		},
		{
			"2 * 10",
			"(2 * 10)",
		},
		{
			"3 / 2",
			"(3 / 2)",
		},
		{
			"3 + 5 * 7",
			"(3 + (5 * 7))",
		},
		{
			"2 / 3 - 2",
			"((2 / 3) - 2)",
		},
		{
			"10 * 3 / 4",
			"((10 * 3) / 4)",
		},
		{
			"true == false",
			"(true == false)",
		},
		{
			"10 != false",
			"(10 != false)",
		},
		{
			"4 < 2 == 3 > 1",
			"((4 < 2) == (3 > 1))",
		},
		{
			"4 >= 2",
			"(4 >= 2)",
		},
		{
			"0 <= 2",
			"(0 <= 2)",
		},
		{
			"!false == true",
			"((!false) == true)",
		},
	}

	for _, tt := range tests {
		testInfixExpression(t, tt.input, tt.expected)
	}
}

func TestIfExpression(t *testing.T) {
	input := `
    if (true) {5 + 5}
    `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T", statements[0])
	}

	ifExpr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected IfExpression. Got %T", stmt)
	}

	_, ok = ifExpr.Condition.(*ast.BooleanExpression)
	if !ok {
		t.Fatalf("Expected BooleanExpression. Got %T", ifExpr.Condition)
	}

	t.Log(program.String())

	blockExpr := ifExpr.Consequence.Statements[0]
	exprStmt, ok := blockExpr.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement in BlockStatement. Got %T", blockExpr)
	}

	_, ok = exprStmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Expected InfixExpression. Got %T", exprStmt)
	}

	alternative := ifExpr.Alternative
	if alternative != nil {
		t.Fatalf("Expected alternative to be nil. Got %T", alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `
    if (true) {5 + 5} else { 10 }
    `

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

	stmt, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement. Got %T", statements[0])
	}

	ifExpr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("Expected IfExpression. Got %T", stmt)
	}

	_, ok = ifExpr.Condition.(*ast.BooleanExpression)
	if !ok {
		t.Fatalf("Expected BooleanExpression. Got %T", ifExpr.Condition)
	}

	blockExpr := ifExpr.Consequence.Statements[0]
	exprStmt, ok := blockExpr.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement in BlockStatement. Got %T", blockExpr)
	}

	_, ok = exprStmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("Expected InfixExpression. Got %T", exprStmt)
	}

	alternativeBlock := ifExpr.Alternative.Statements[0]
	alternative, ok := alternativeBlock.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement in alternative. Got %T", alternativeBlock)
	}

	_, ok = alternative.Expression.(*ast.IntegerExpression)
	if !ok {
		t.Fatalf("Expected IntegerExpression in alternative. Got %T", alternative.Expression)
	}
}

func TestFunctionDefinition(t *testing.T) {
	input := `
    fn (arg, anotherArg) {
        return arg + anotherArg;
    }
    `
	program := constructProgram(input, t)

	statements := program.Statements
	expectProgramLength(t, statements, 1)

	exprStatement, ok := statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected FunctionExpression. Got %T", statements[0])
	}

	fnExpr, ok := exprStatement.Expression.(*ast.FunctionExpression)
	if !ok {
		t.Fatalf("Expected FunctionExpression. Got %T", exprStatement.Expression)
	}

	parameters := fnExpr.Parameters
	if len(parameters) != 2 {
		t.Fatalf("Expected 2 paramaters. Got %d", len(parameters))
	}

	firstParam := parameters[0]
	if firstParam == nil {
		t.Fatalf("Expected first parameter to be defined. Got nil")
	}
	secondParam := parameters[0]
	if secondParam == nil {
		t.Fatalf("Expected second parameter to be defined. Got nil")
	}

	body := fnExpr.Body
	if body == nil {
		t.Fatalf("Expected body to be defined. Got nil")
	}

	returnStmt, ok := body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("Expected ReturnStatement. Got %T", returnStmt)
	}
}

func constructProgram(input string, t *testing.T) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	return program
}

func testInfixExpression(t *testing.T, input string, expected string) {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	programString := program.String()
	if programString != expected {
		t.Fatalf("Expected '%s'. Got '%s'", expected, programString)
	}
}

func expectProgramLength(t *testing.T, statements []ast.Statement, expected int) {
	if len(statements) != expected {
		t.Fatalf("Program has wrong numer of statements. Expected %d. Got %d", expected, len(statements))
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

func checkParserErrors(t *testing.T, p *Parser) {
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
