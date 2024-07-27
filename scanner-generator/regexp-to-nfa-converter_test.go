package scannergenerator

import (
	"slices"
	"testing"
)

func TestRegexpToNfaConversion(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]map[int][]int
	}{
		// {
		// 	"a",
		// 	map[string]map[int][]int{
		// 		"a": {0: []int{1}},
		// 	},
		// },
		// {
		// 	"b",
		// 	map[string]map[int][]int{
		// 		"b": {0: []int{1}},
		// 	},
		// },
		// {
		// 	"ab",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		EPSILON: {1: []int{2}},
		// 	},
		// },
		// {
		// 	"abc",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		"c":     {4: []int{5}},
		// 		EPSILON: {1: []int{2}, 3: []int{4}},
		// 	},
		// },
		// {
		// 	"a|b",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		EPSILON: {4: []int{0, 2}, 1: []int{5}, 3: []int{5}},
		// 	},
		// },
		// {
		// 	"ab|c",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		"c":     {4: []int{5}},
		// 		EPSILON: {1: []int{2}, 6: []int{0, 4}, 3: []int{7}, 5: []int{7}},
		// 	},
		// },
		// {
		// 	"a|bc",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		"c":     {4: []int{5}},
		// 		EPSILON: {6: []int{0, 2}, 1: []int{7}, 3: []int{4}, 5: []int{7}},
		// 	},
		// },
		// {
		// 	"a*",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		EPSILON: {2: []int{0, 3}, 1: []int{0, 3}},
		// 	},
		// },
		// {
		// 	"a(b|c)",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		"c":     {4: []int{5}},
		// 		EPSILON: {6: []int{2, 4}, 3: []int{7}, 5: []int{7}, 1: []int{6}},
		// 	},
		// },
		// {
		// 	"a(b|c)*",
		// 	map[string]map[int][]int{
		// 		"a":     {0: []int{1}},
		// 		"b":     {2: []int{3}},
		// 		"c":     {4: []int{5}},
		// 		EPSILON: {6: []int{2, 4}, 3: []int{7}, 5: []int{7}, 1: []int{8}, 8: []int{6, 9}, 7: []int{6, 9}},
		// 	},
		// },
		// {
		// 	"1|2|3",
		// 	map[string]map[int][]int{
		// 		"1":     {0: []int{1}},
		// 		"2":     {2: []int{3}},
		// 		"3":     {4: []int{5}},
		// 		EPSILON: {1: []int{9}, 3: []int{7}, 5: []int{7}, 6: []int{4, 2}, 7: []int{9}, 8: []int{2, 0}},
		// 	},
		// },
		// {
		// 	"[1-3]",
		// 	map[string]map[int][]int{
		// 		"1":     {0: []int{1}},
		// 		"2":     {2: []int{3}},
		// 		"3":     {4: []int{5}},
		// 		EPSILON: {1: []int{9}, 3: []int{7}, 5: []int{7}, 6: []int{4, 2}, 7: []int{9}, 8: []int{2, 0}},
		// 	},
		// },
	}

	for _, tt := range tests {
		regexpToNfaConverter := &RegexpToNfaConverter{regexp: tt.input}
		result := regexpToNfaConverter.Convert().Transitions

		t.Logf("current input: %s", tt.input)
		t.Logf("expected: %v", tt.expected)
		t.Logf("result: %v", result)
		if result == nil {
			t.Fatalf("expected result not to be nil")
		}

		for symbol, transitionsForSymbol := range tt.expected {
			resultMappings := result[symbol]
			if resultMappings == nil {
				t.Fatalf("expected transitions for symbol '%s' but was nil", symbol)
			}

			if len(resultMappings) != len(transitionsForSymbol) {
				t.Fatalf("transitions for symbol '%s' differ in size. Expected %v. Got %v", symbol, transitionsForSymbol, resultMappings)
			}

			for state, transitions := range transitionsForSymbol {
				mappingInResultForSymbol, ok := resultMappings[state]
				if !ok {
					t.Fatalf("expected transitions from state %d under symbol '%s' but was not defined",
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
		t.Fatalf("expected transitions %v but got %v", expected, actual)
	}

	slices.Sort(expected)
	slices.Sort(actual)
	for i, expectedState := range expected {
		if expectedState != actual[i] {
			t.Fatalf("expected transitions %v but got %v", expected, actual)
		}
	}
}
