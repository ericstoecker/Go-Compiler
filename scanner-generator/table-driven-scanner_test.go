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

	nfaToDfaConverter := NewNfaToDfaConverter(nfa)
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
		scanner := New(tt.input, dfa)
		token := scanner.NextToken()
		if token != tt.expectedToken {
			t.Logf("current input: %s", tt.input)
			t.Errorf("expected: %v, got: %v", tt.expectedToken, token)
		}
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
	dfa := NewNfaToDfaConverter(combinedNfa).Convert()

	tests := []struct {
		inputs        string
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
		scanner := New(tt.inputs, dfa)
		token := scanner.NextToken()
		if token != tt.expectedToken {
			t.Logf("current input: %s", tt.inputs)
			t.Errorf("expected: %v, got: %v", tt.expectedToken, token)
		}
	}
}
