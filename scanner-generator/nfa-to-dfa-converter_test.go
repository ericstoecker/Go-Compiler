package scannergenerator

import "testing"

func TestNfaToDfaConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]map[int]int
	}{
		{
			"ab",
			map[string]map[int]int{
				"a": {0: 1},
				"b": {1: 2},
			},
		},
		{
			"a|b",
			map[string]map[int]int{
				"a": {0: 1},
				"b": {0: 2},
				// add accepting states testing
			},
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := NewRegexpToNfaConverter(tt.input)
		nfa := regexpToNfaConverter.Convert()

		nfaToDfaConverter := NewNfaToDfaConverter(nfa)
		result := nfaToDfaConverter.Convert()

		t.Logf("current input: %s", tt.input)
		t.Logf("expected: %v", tt.expected)
		t.Logf("result: %v", result)
		if result == nil {
			t.Fatalf("expected result not to be nil")
		}

		if len(result) != len(tt.expected) {
			t.Fatalf("sizes of transition tables differ. Expected %d. Got %d", len(result), len(tt.expected))
		}

		for symbol, transitionsForSymbol := range tt.expected {
			resultMappings := result[symbol]
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

}
