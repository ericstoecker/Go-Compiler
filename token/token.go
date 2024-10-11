package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT   = "IDENT"
	INT     = "INT"
	STRING  = "STRING"

	ASSIGN = "="
	BANG   = "!"
	PLUS   = "+"
	MINUS  = "-"

	SLASH    = "/"
	ASTERISK = "*"
	LT       = "<"
	GT       = ">"

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"

	EQUALS        = "=="
	NOT_EQUALS    = "!="
	LESS_EQUAL    = "<="
	GREATER_EQUAL = ">="
	AND           = "&&"
	OR            = "||"

	IF   = "if"
	ELSE = "else"

	RETURN = "return"

	TRUE  = "true"
	FALSE = "false"
)

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENT
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}
