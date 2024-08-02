package scannergenerator

import (
	"compiler/token"
	"testing"
)

func TestNfaToDfaConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected Dfa
	}{
		{
			"ab",
			Dfa{
				Transitions: map[string]map[int]int{
					"a": {0: 1},
					"b": {1: 2},
				},
				AcceptingStates: []int{2},
			},
		},
		{
			"a|b",
			Dfa{
				Transitions: map[string]map[int]int{
					"a": {0: 1},
					"b": {0: 2},
				},
				AcceptingStates: []int{1, 2},
			},
		},
		{
			"a*",
			Dfa{
				Transitions: map[string]map[int]int{
					"a": {0: 1, 1: 1},
				},
				AcceptingStates: []int{0, 1},
			},
		},
		{
			"aa*",
			Dfa{
				Transitions: map[string]map[int]int{
					"a": {0: 1, 1: 2, 2: 2},
				},
				AcceptingStates: []int{1, 2},
			},
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := NewRegexpToNfaConverter(tt.input)
		nfa, _ := regexpToNfaConverter.Convert()

		nfaToDfaConverter := NewNfaToDfaConverter(nfa, map[token.TokenType]int{})
		dfa := nfaToDfaConverter.Convert()
		t.Logf("current input: %s", tt.input)
		t.Logf("expected: %v", tt.expected)
		t.Logf("result: %v", dfa)

		if dfa == nil {
			t.Fatalf("expected dfa not to be nil")
		}

		testAcceptingStates(t, tt.expected.AcceptingStates, dfa.AcceptingStates)

		dfaTransitions := dfa.Transitions
		testTransitions(t, tt.expected.Transitions, dfaTransitions)
	}
}

func TestMultistateNfaToDfaConversion(t *testing.T) {
	tests := []struct {
		input      []string
		tokenTypes []token.TokenType
		expected   Dfa
	}{
		{
			[]string{"a", "b"},
			[]token.TokenType{"first", "second"},
			Dfa{
				Transitions: map[string]map[int]int{
					"a": {0: 1},
					"b": {0: 2},
				},
				AcceptingStates: []int{1, 2},
				TypeTable:       map[int]token.TokenType{1: "first", 2: "second"},
			},
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := NewRegexpToNfaConverter(tt.input[0])
		nfa, _ := regexpToNfaConverter.Convert()
		nfa.TypeTable = map[int]token.TokenType{1: tt.tokenTypes[0]}

		secondConverter := NewRegexpToNfaConverter(tt.input[1])
		secondNfa, _ := secondConverter.Convert()
		secondNfa.TypeTable = map[int]token.TokenType{1: tt.tokenTypes[1]}

		nfaToDfaConverter := NewNfaToDfaConverter(nfa.UnionDistinct(secondNfa), map[token.TokenType]int{})
		dfa := nfaToDfaConverter.Convert()
		t.Logf("current input: %s", tt.input)
		t.Logf("expected: %v", tt.expected)
		t.Logf("result: %v", dfa)

		if dfa == nil {
			t.Fatalf("expected dfa not to be nil")
		}

		testAcceptingStates(t, tt.expected.AcceptingStates, dfa.AcceptingStates)

		dfaTransitions := dfa.Transitions
		testTransitions(t, tt.expected.Transitions, dfaTransitions)

		testTypeTables(t, tt.expected.TypeTable, dfa.TypeTable)
	}

}

func testTypeTables(t *testing.T, expected, actual map[int]token.TokenType) {
	for state, expectedType := range expected {
		actualType, ok := actual[state]
		if !ok {
			t.Fatalf("expected type for state %d but was not defined", state)
		}

		if expectedType != actualType {
			t.Fatalf("expected type for state %d to be %s. Got %s", state, expectedType, actualType)
		}
	}
}

func testTransitions(t *testing.T, expected, actual map[string]map[int]int) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expected dfa transitions not to be nil")
	}

	if len(actual) != len(expected) {
		t.Fatalf("sizes of transition tables differ. Expected %d. Got %d", len(actual), len(expected))
	}

	for symbol, transitionsForSymbol := range expected {
		resultMappings := actual[symbol]
		if resultMappings == nil {
			t.Fatalf("expected transitions for symbol '%s' but was nil", symbol)
		}

		if len(resultMappings) != len(transitionsForSymbol) {
			t.Fatalf("transitions for symbol '%s' differ in size. Expected %v. Got %v", symbol, transitionsForSymbol, resultMappings)
		}

		for state, expectedTransition := range transitionsForSymbol {
			actualTransition, ok := resultMappings[state]
			if !ok {
				t.Fatalf("expected transition from state %d under symbol '%s' but was not defined",
					state, symbol)
			}

			if expectedTransition != actualTransition {
				t.Fatalf("expected transition from state %d to state %d under symbol '%s'. Got %d",
					state, expectedTransition, symbol, actualTransition)
			}
		}
	}
}
func testAcceptingStates(t *testing.T, expected, actual []int) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expected accepting states not to be nil")
	}

	if len(expected) != len(actual) {
		t.Fatalf("sizes of accepting states differ. Expected %d. Got %d", len(expected), len(actual))
	}

	for i, expectedState := range expected {
		actualState := actual[i]
		if expectedState != actualState {
			t.Fatalf("expected accepting state %d. Got %d", expectedState, actualState)
		}
	}
}
