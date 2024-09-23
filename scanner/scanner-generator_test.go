package scanner

import (
	"compiler/token"
	"testing"
)

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
	-

	&&
	||

	if
	else
	true
	false

	"test1T EST2"
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
		{token.COLON, ":"},
		{token.GT, ">"},
		{token.LT, "<"},
		{token.SLASH, "/"},
		{token.MINUS, "-"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.STRING, "\"test1T EST2\""},
		{token.EOF, ""},
	}

	scannerGenerator := NewScannerGenerator()
	dfa := scannerGenerator.GenerateScanner(TokenClassifications)
	s := NewTableDrivenScanner(input, dfa)

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
