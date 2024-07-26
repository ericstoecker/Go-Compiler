package scannergenerator

import (
	"compiler/token"
	"testing"
)

func TestAccepting(t *testing.T) {
	RegexpToNfaConverter := NewRegexpToNfaConverter("a(b|c)*")
	nfa := RegexpToNfaConverter.Convert()

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
