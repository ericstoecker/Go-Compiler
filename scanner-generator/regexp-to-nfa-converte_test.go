package scannergenerator

import (
	"slices"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]map[int][]int
	}{
		{
			"a",
			map[string]map[int][]int{
				"a": {0: []int{1}},
			},
		},
		{
			"b",
			map[string]map[int][]int{
				"b": {0: []int{1}},
			},
		},
		{
			"ab",
			map[string]map[int][]int{
				"a":     {0: []int{1}},
				"b":     {2: []int{3}},
				EPSILON: {1: []int{2}},
			},
		},
		{
			"abc",
			map[string]map[int][]int{
				"a":     {0: []int{1}},
				"b":     {2: []int{3}},
				"c":     {4: []int{5}},
				EPSILON: {1: []int{2}, 3: []int{4}},
			},
		},
		{
			"a|b",
			map[string]map[int][]int{
				"a":     {0: []int{1}},
				"b":     {2: []int{3}},
				EPSILON: {4: []int{0, 2}, 1: []int{5}, 3: []int{5}},
			},
		},
		{
			"ab|c",
			map[string]map[int][]int{
				"a":     {0: []int{1}},
				"b":     {2: []int{3}},
				"c":     {4: []int{5}},
				EPSILON: {1: []int{2}, 6: []int{0, 4}, 3: []int{7}, 5: []int{7}},
			},
		},
		{
			"a|bc",
			map[string]map[int][]int{
				"a":     {0: []int{1}},
				"b":     {2: []int{3}},
				"c":     {4: []int{5}},
				EPSILON: {6: []int{0, 2}, 1: []int{7}, 3: []int{4}, 5: []int{7}},
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

		for symbol, transitionsForSymbol := range tt.expected {
			resultMappings := result[symbol]
			if resultMappings == nil {
				t.Errorf("expected transitions for symbol '%s' but was nil", symbol)
			}

			if len(resultMappings) != len(transitionsForSymbol) {
				t.Errorf("transitions for symbol '%s' differ in size. Expected %v. Got %v", symbol, transitionsForSymbol, resultMappings)
			}

			for state, transitions := range transitionsForSymbol {
				mappingInResultForSymbol, ok := resultMappings[state]
				if !ok {
					t.Errorf("expected transitions from state %d under symbol '%s' but was not defined",
						state, symbol)
				} else {
					compareTransitions(transitions, mappingInResultForSymbol, t)
				}
			}
		}
	}
}

func compareTransitions(expected, actual []int, t *testing.T) {
	if len(expected) != len(actual) {
		t.Errorf("expected transitions %v but got %v", expected, actual)
	}

	slices.Sort(expected)
	slices.Sort(actual)
	for i, expectedState := range expected {
		if expectedState != actual[i] {
			t.Errorf("expected transitions %v but got %v", expected, actual)
		}
	}
}
