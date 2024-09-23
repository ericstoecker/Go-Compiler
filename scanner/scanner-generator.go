package scanner

import (
	"compiler/token"
	"fmt"
)

type TokenClassification struct {
	Regexp     string
	TokenType  token.TokenType
	Precedence int
}

type ScannerGenerator struct {
}

func NewScannerGenerator() *ScannerGenerator {
	return &ScannerGenerator{}
}

func (s *ScannerGenerator) GenerateScanner(tokenClassifications []TokenClassification) *Dfa {
	precedences := make(map[token.TokenType]int)
	nfas := make([]*Nfa, len(tokenClassifications))
	for i, tokenClassification := range tokenClassifications {
		nfaForClassification, err := NewRegexpToNfaConverter(tokenClassification.Regexp).Convert()
		if err != nil {
			panic(fmt.Sprintf("error when converting regexp '%s' to nfa: %v", tokenClassification.Regexp, err))
		}

		for _, state := range nfaForClassification.AcceptingStates {
			nfaForClassification.TypeTable[state] = tokenClassification.TokenType
		}

		nfas[i] = nfaForClassification

		precedences[tokenClassification.TokenType] = tokenClassification.Precedence
	}

	nfa := nfas[0].UnionDistinct(nfas[1:]...)
	dfa := NewNfaToDfaConverter(nfa, precedences).Convert()

	dfaMinimizer := &DfaMinimizer{}

	return dfaMinimizer.Minimize(dfa)
}

var TokenClassifications = []TokenClassification{
	{"=", token.ASSIGN, 1},
	{"+", token.PLUS, 1},
	{"-", token.MINUS, 1},
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
	{"&&", token.AND, 1},
	{"\\|\\|", token.OR, 1},
	{"/", token.SLASH, 1},
	{"let", token.LET, 2},
	{"return", token.RETURN, 2},
	{"fn", token.FUNCTION, 2},
	{"if", token.IF, 2},
	{"else", token.ELSE, 2},
	{"true", token.TRUE, 2},
	{"false", token.FALSE, 2},
	{"[a-z]([a-z]|[A-Z])*", token.IDENT, 1},
	{"[0-9]([0-9])*", token.INT, 1},
	{`"([a-z]|[A-Z]|[0-9]| )*"`, token.STRING, 1},
}
