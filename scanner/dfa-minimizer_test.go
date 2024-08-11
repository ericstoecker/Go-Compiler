package scanner

import (
	"compiler/token"
	"testing"
)

func TestMinimization(t *testing.T) {
	tests := []struct {
		regexp string
		dfa    Dfa
	}{
		{
			"a|b",
			Dfa{
				InitialState:    1,
				AcceptingStates: []int{0},
				Transitions: map[string]map[int]int{
					"a": {1: 0},
					"b": {1: 0},
				},
			},
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := NewRegexpToNfaConverter(tt.regexp)
		nfa, err := regexpToNfaConverter.Convert()
		if err != nil {
			t.Fatalf("error when converting regexp to nfa: %e", err)
		}

		nfaToDfaConverter := NewNfaToDfaConverter(nfa, make(map[token.TokenType]int))
		dfa := nfaToDfaConverter.Convert()
		dfa.TypeTable = map[int]token.TokenType{1: "test", 2: "test"}

		m := &DfaMinimizer{}
		minimizedDfa := m.Minimize(dfa)

		t.Logf("input: %s", tt.regexp)
		t.Logf("initial dfa: %v", dfa)
		t.Logf("result of minimization: %v", minimizedDfa)
		t.Logf("expected minimal dfa: %v", tt.dfa)

		if minimizedDfa == nil {
			t.Fatalf("expected minimized dfa not to be nil")
		}

		testAcceptingStates(t, tt.dfa.AcceptingStates, minimizedDfa.AcceptingStates)

		testTransitions(t, tt.dfa.Transitions, minimizedDfa.Transitions)
	}
}
