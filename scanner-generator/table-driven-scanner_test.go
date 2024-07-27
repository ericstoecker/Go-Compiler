package scannergenerator

import (
	"compiler/token"
	"testing"
)

func TestAccepting(t *testing.T) {
	RegexpToNfaConverter := NewRegexpToNfaConverter("a(b|c)*")
	nfa := RegexpToNfaConverter.Convert()
	nfa.TypeTable = map[int]token.TokenType{
		9: "ACCEPT",
	}

	nfaToDfaConverter := NewNfaToDfaConverter(nfa, map[token.TokenType]int{})
	dfa := nfaToDfaConverter.Convert()

	tests := []struct {
		input         string
		expectedToken token.Token
	}{
		{
			"a",
			token.Token{Type: "ACCEPT", Literal: "a"},
		},
		{
			"ab",
			token.Token{Type: "ACCEPT", Literal: "ab"},
		},
		{
			"ac",
			token.Token{Type: "ACCEPT", Literal: "ac"},
		},
		{
			"abb",
			token.Token{Type: "ACCEPT", Literal: "abb"},
		},
		{
			"d",
			token.Token{Type: token.ILLEGAL, Literal: "d"},
		},
		{
			" a",
			token.Token{Type: "ACCEPT", Literal: "a"},
		},
	}

	for _, tt := range tests {
		testNextChar(t, tt.input, dfa, tt.expectedToken)
	}
}

func TestMultipleCategories(t *testing.T) {
	firstConverter := NewRegexpToNfaConverter("ab*")
	firstNfa := firstConverter.Convert()
	firstNfa.TypeTable = map[int]token.TokenType{
		5: "FIRST",
	}

	secondConverter := NewRegexpToNfaConverter("cd*")
	secondNfa := secondConverter.Convert()
	secondNfa.TypeTable = map[int]token.TokenType{
		5: "SECOND",
	}

	combinedNfa := firstNfa.UnionDistinct(secondNfa)
	dfa := NewNfaToDfaConverter(combinedNfa, map[token.TokenType]int{}).Convert()

	tests := []struct {
		input         string
		expectedToken token.Token
	}{
		{
			"ab",
			token.Token{Type: "FIRST", Literal: "ab"},
		},
		{
			"cd",
			token.Token{Type: "SECOND", Literal: "cd"},
		},
		{
			"ll",
			token.Token{Type: token.ILLEGAL, Literal: "l"},
		},
		{
			"ab cd",
			token.Token{Type: "FIRST", Literal: "ab"},
		},
		{
			"  abbbb",
			token.Token{Type: "FIRST", Literal: "abbbb"},
		},
	}

	for _, tt := range tests {
		testNextChar(t, tt.input, dfa, tt.expectedToken)
	}
}

func TestConflictingCategories(t *testing.T) {
	firstConverter := NewRegexpToNfaConverter("a")
	firstNfa := firstConverter.Convert()
	firstNfa.TypeTable = map[int]token.TokenType{
		1: "FIRST",
	}

	secondConverter := NewRegexpToNfaConverter("aa*")
	secondNfa := secondConverter.Convert()
	secondNfa.TypeTable = map[int]token.TokenType{
		5: "SECOND",
	}

	combinedNfa := firstNfa.UnionDistinct(secondNfa)
	typePrecedences := map[token.TokenType]int{
		"FIRST":  2,
		"SECOND": 1,
	}
	dfa := NewNfaToDfaConverter(combinedNfa, typePrecedences).Convert()

	tests := []struct {
		input         string
		expectedToken token.Token
	}{
		{
			"a",
			token.Token{Type: "FIRST", Literal: "a"},
		},
		{
			"aa",
			token.Token{Type: "SECOND", Literal: "aa"},
		},
	}

	for _, tt := range tests {
		testNextChar(t, tt.input, dfa, tt.expectedToken)
	}
}

func testNextChar(t *testing.T, input string, dfa *Dfa, expectedToken token.Token) {
	t.Helper()

	scanner := New(input, dfa)
	token := scanner.NextToken()
	if token != expectedToken {
		t.Logf("current input: %s", input)
		t.Errorf("expected: %v, got: %v", expectedToken, token)
	}
}
