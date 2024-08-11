package scanner

import (
	"compiler/token"
	"testing"
)

func TestAccepting(t *testing.T) {
	RegexpToNfaConverter := NewRegexpToNfaConverter("a(b|c)*")
	nfa, _ := RegexpToNfaConverter.Convert()
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

func TestSpecialCharacters(t *testing.T) {
	tests := []struct {
		regexp        string
		input         string
		expectedToken token.Token
	}{
		{
			"[1-3]([1-5])*",
			"1",
			token.Token{Type: "ACCEPT", Literal: "1"},
		},
		{
			"[1-3]([1-5])*",
			"2",
			token.Token{Type: "ACCEPT", Literal: "2"},
		},
		{
			"[1-3]([1-5])*",
			"3",
			token.Token{Type: "ACCEPT", Literal: "3"},
		},
		{
			"[1-3]([1-5])*",
			"4",
			token.Token{Type: token.ILLEGAL, Literal: "4"},
		},
		{
			"[1-3]([1-5])*",
			"15",
			token.Token{Type: "ACCEPT", Literal: "15"},
		},
		{
			"[1-3]([1-5])*",
			"151",
			token.Token{Type: "ACCEPT", Literal: "151"},
		},
		{
			"\\(",
			"(",
			token.Token{Type: "ACCEPT", Literal: "("},
		},
		{
			"(\\()",
			"(",
			token.Token{Type: "ACCEPT", Literal: "("},
		},
		{
			"\\)a",
			")a",
			token.Token{Type: "ACCEPT", Literal: ")a"},
		},
		{
			"[a-z]",
			"a",
			token.Token{Type: "ACCEPT", Literal: "a"},
		},
		{
			"[a-z]",
			"t",
			token.Token{Type: "ACCEPT", Literal: "t"},
		},
		{
			"[a-c]",
			"d",
			token.Token{Type: token.ILLEGAL, Literal: "d"},
		},
		{
			"[A-Z]",
			"C",
			token.Token{Type: "ACCEPT", Literal: "C"},
		},
		{
			"[A-Z]",
			"c",
			token.Token{Type: token.ILLEGAL, Literal: "c"},
		},
	}

	for _, tt := range tests {
		RegexpToNfaConverter := NewRegexpToNfaConverter(tt.regexp)
		nfa, _ := RegexpToNfaConverter.Convert()
		nfa.TypeTable = make(map[int]token.TokenType)
		for _, acceptingState := range nfa.AcceptingStates {
			nfa.TypeTable[acceptingState] = "ACCEPT"
		}

		nfaToDfaConverter := NewNfaToDfaConverter(nfa, map[token.TokenType]int{})
		dfa := nfaToDfaConverter.Convert()
		testNextChar(t, tt.input, dfa, tt.expectedToken)
	}
}

func TestMultipleCategories(t *testing.T) {
	firstConverter := NewRegexpToNfaConverter("ab*")
	firstNfa, _ := firstConverter.Convert()
	firstNfa.TypeTable = map[int]token.TokenType{
		5: "FIRST",
	}

	secondConverter := NewRegexpToNfaConverter("([c-d])*")
	secondNfa, _ := secondConverter.Convert()
	for _, acceptingState := range secondNfa.AcceptingStates {
		secondNfa.TypeTable[acceptingState] = "SECOND"
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
	firstNfa, _ := firstConverter.Convert()
	firstNfa.TypeTable = map[int]token.TokenType{
		1: "FIRST",
	}

	secondConverter := NewRegexpToNfaConverter("aa*")
	secondNfa, _ := secondConverter.Convert()
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

	scanner := NewTableDrivenScanner(input, dfa)
	token := scanner.NextToken()
	if token != expectedToken {
		t.Logf("current input: %s", input)
		t.Errorf("expected: %v, got: %v", expectedToken, token)
	}
}
