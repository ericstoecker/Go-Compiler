package scannergenerator

import "testing"

func TestConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]map[int]int
	}{
		{
			"a",
			map[string]map[int]int{
				"a": {0: 1},
			},
		},
		{
			"b",
			map[string]map[int]int{
				"b": {0: 1},
			},
		},
		{
			"ab",
			map[string]map[int]int{
				"a":     {0: 1},
				EPSILON: {1: 2},
				"b":     {2: 3},
			},
		},
		{
			"abc",
			map[string]map[int]int{
				"a":     {0: 1},
				EPSILON: {1: 2, 3: 4},
				"b":     {2: 3},
				"c":     {4: 5},
			},
		},
	}

	for _, tt := range tests {
		regexpToNfaConverter := &RegexpToNfaConverter{input: tt.input}
		result := regexpToNfaConverter.Convert()

		t.Logf("current input: %s", tt.input)
		t.Logf("expected: %v", tt.expected)
		t.Logf("result: %v", result)
		if result == nil {
			t.Errorf("expected result not to be nil")
		}

		for symbol, stateMappings := range tt.expected {
			resultMappings := result[symbol]
			if resultMappings == nil {
				t.Errorf("expected mappings for symbol '%s' but was nil", symbol)
			}

			if len(resultMappings) != len(stateMappings) {
				t.Errorf("transitions for symbol '%s' differ in size. Expected %v. Got %v", symbol, stateMappings, resultMappings)
			}

			for stateFrom, stateTo := range stateMappings {
				mappingInResultForSymbol, ok := resultMappings[stateFrom]
				if !ok {
					t.Errorf("expected mapping from state %d under symbol '%s' but was not defined",
						stateFrom, symbol)
				} else if stateTo != mappingInResultForSymbol {
					t.Errorf("state mapping from state %d under symbol '%s' differs. Expected %d. Got %d",
						stateFrom, symbol, stateTo, mappingInResultForSymbol)
				}
			}
		}
	}
}
