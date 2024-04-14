package token

type TokenType string

type Token struct {
    Type TokenType
    Literal string
}

const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"
    IDENT = "IDENT"
    INT = "INT"

    ASSIGN = "="
    PLUS = "+"

    COMMA = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"

    FUNCTION = "FUNCTION"
    LET = "LET"
)

func LookupIdentifier(identifier string) TokenType {
    if tok, ok := keywords[identifier]; ok {
        return tok
    }
    return IDENT
}

var keywords = map[string]TokenType {
    "fn": FUNCTION,
    "let": LET,
}
