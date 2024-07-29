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

	SLASH   = "/"
	ASTERIK = "*"
	LT      = "<"
	GT      = ">"

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

type TokenClassification struct {
	Regexp     string
	TokenType  TokenType
	Precedence int
}

var TokenClassifications = []TokenClassification{
	{"=", ASSIGN, 1},
	{"+", PLUS, 1},
	{"-", MINUS, 1},
	{",", COMMA, 1},
	{";", SEMICOLON, 1},
	{":", COLON, 1},
	{"\\(", LPAREN, 1},
	{"\\)", RPAREN, 1},
	{"{", LBRACE, 1},
	{"}", RBRACE, 1},
	{"\\[", LBRACKET, 1},
	{"\\]", RBRACKET, 1},
	{">", GT, 1},
	{">=", GREATER_EQUAL, 1},
	{"<", LT, 1},
	{"<=", LESS_EQUAL, 1},
	{"==", EQUALS, 1},
	{"!", BANG, 1},
	{"!=", NOT_EQUALS, 1},
	{"&&", AND, 1},
	{"\\|\\|", OR, 1},
	{"/", SLASH, 1},
	{"let", LET, 2},
	{"return", RETURN, 2},
	{"fn", FUNCTION, 2},
	{"if", IF, 2},
	{"else", ELSE, 2},
	{"true", TRUE, 2},
	{"false", FALSE, 2},
	{"[a-z]([a-z]|[A-Z])*", IDENT, 1},
	{"[0-9]([0-9])*", INT, 1},
	{`"([a-z]|[A-Z]|[0-9]| )*"`, STRING, 1},
}
