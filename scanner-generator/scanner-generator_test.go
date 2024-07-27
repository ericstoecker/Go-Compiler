package scannergenerator

import (
	"compiler/token"
	"testing"
)

var tokenClassifications = []TokenClassification{
	{"let", token.LET, 2},
	{"[a-z]([a-z])*", token.IDENT, 1},
}

func TestGeneratedScanner(t *testing.T) {
	input := `let five`
	// let ten = 10;

	// let add = fn(x, y) {
	//     x + y;
	// };

	// let result = add(five, ten);

	// ==
	// !=
	// <=
	// >=
	// !
	// return
	// []
	// "test"
	// :
	// `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		// {token.ASSIGN, "="},
		// {token.INT, "5"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "ten"},
		// {token.ASSIGN, "="},
		// {token.INT, "10"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "add"},
		// {token.ASSIGN, "="},
		// {token.FUNCTION, "fn"},
		// {token.LPAREN, "("},
		// {token.IDENT, "x"},
		// {token.COMMA, ","},
		// {token.IDENT, "y"},
		// {token.RPAREN, ")"},
		// {token.LBRACE, "{"},
		// {token.IDENT, "x"},
		// {token.PLUS, "+"},
		// {token.IDENT, "y"},
		// {token.SEMICOLON, ";"},
		// {token.RBRACE, "}"},
		// {token.SEMICOLON, ";"},
		// {token.LET, "let"},
		// {token.IDENT, "result"},
		// {token.ASSIGN, "="},
		// {token.IDENT, "add"},
		// {token.LPAREN, "("},
		// {token.IDENT, "five"},
		// {token.COMMA, ","},
		// {token.IDENT, "ten"},
		// {token.RPAREN, ")"},
		// {token.SEMICOLON, ";"},
		// {token.EQUALS, "=="},
		// {token.NOT_EQUALS, "!="},
		// {token.LESS_EQUAL, "<="},
		// {token.GREATER_EQUAL, ">="},
		// {token.BANG, "!"},
		// {token.RETURN, "return"},
		// {token.LBRACKET, "["},
		// {token.RBRACKET, "]"},
		// {token.STRING, "test"},
		// {token.COLON, ":"},
		// {token.EOF, ""},
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
