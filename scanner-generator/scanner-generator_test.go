package scannergenerator

import (
	"compiler/token"
	"testing"
)

var tokenClassifications = []TokenClassification{
	{"=", token.ASSIGN, 1},
	{"+", token.PLUS, 1},
	{",", token.COMMA, 1},
	{";", token.SEMICOLON, 1},
	{":", token.COLON, 1},
	{"\\(", token.LPAREN, 1},
	{"\\)", token.RPAREN, 1},
	{"{", token.LBRACE, 1},
	{"}", token.RBRACE, 1},
	{"\\[", token.LBRACKET, 1},
	{"\\]", token.RBRACKET, 1},
	{">", token.GT, 1},
	{">=", token.GREATER_EQUAL, 1},
	{"<", token.LT, 1},
	{"<=", token.LESS_EQUAL, 1},
	{"==", token.EQUALS, 1},
	{"!", token.BANG, 1},
	{"!=", token.NOT_EQUALS, 1},
	{"/", token.SLASH, 1},
	{"let", token.LET, 2},
	{"return", token.RETURN, 2},
	{"fn", token.FUNCTION, 2},
	{"\"([a-z])*\"", token.STRING, 1},
	{"[a-z]([a-z]|[A-Z])*", token.IDENT, 1},
	{"[0-9]([0-9])*", token.INT, 1},
}

func TestGeneratedScanner(t *testing.T) {
	input := `let fiveTen = 55;
	let ten = 10;

	let add = fn(x, y) {
	    x + y;
	};

	let result = add(five, ten);

	==
	!=
	<=
	>=
	!
	return
	[]
	:

	>
	<
	/
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "fiveTen"},
		{token.ASSIGN, "="},
		{token.INT, "55"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EQUALS, "=="},
		{token.NOT_EQUALS, "!="},
		{token.LESS_EQUAL, "<="},
		{token.GREATER_EQUAL, ">="},
		{token.BANG, "!"},
		{token.RETURN, "return"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		// {token.STRING, "\"test\""},
		{token.COLON, ":"},
		{token.GT, ">"},
		{token.LT, "<"},
		{token.SLASH, "/"},
		{token.EOF, ""},
	}

	scannerGenerator := NewScannerGenerator()
	dfa := scannerGenerator.GenerateScanner(tokenClassifications)
	s := New(input, dfa)

	for i, tt := range tests {
		tok := s.NextToken()

		t.Logf("result: %v", tok)
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - type wrong. expected=%q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
