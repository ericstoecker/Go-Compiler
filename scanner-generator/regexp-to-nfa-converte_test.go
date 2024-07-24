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
	}

	regexpToNfaConverter := &RegexpToNfaConverter{}
	for _, tt := range tests {
		result := regexpToNfaConverter.Convert(tt.input)

		t.Logf("current input: %s", tt.input)
		if result == nil {
			t.Errorf("expected result not to be nil")
		}

		for symbol, stateMappings := range tt.expected {
			resultMappings := result[symbol]
			if resultMappings == nil {
				t.Errorf("expected mappings for symbol '%s' but was nil", symbol)
			}

			if len(resultMappings) != len(stateMappings) {
				t.Errorf("state mappings differ in size. Expected %v. Got %v", stateMappings, resultMappings)
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
