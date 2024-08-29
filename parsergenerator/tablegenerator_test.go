package parsergenerator

import (
	"compiler/grammar"
	"compiler/token"
	"maps"
	"slices"
	"testing"
)

const (
	GOAL       = "Goal"
	LIST       = "List"
	PAIR       = "Pair"
	LEFTPAREN  = "("
	RIGHTPAREN = ")"
)

func festTableConstruction(t *testing.T) {
	productions := []grammar.Production{
		&grammar.NonTerminal{
			Name: GOAL,
			RightSide: &grammar.Identifier{
				Name: LIST,
			},
		},
		&grammar.NonTerminal{
			Name: LIST,
			RightSide: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{
								Name: LIST,
							},
							{
								Name: PAIR,
							},
						},
					},
					&grammar.Identifier{
						Name: PAIR,
					},
				},
			},
		},
		&grammar.NonTerminal{
			Name: PAIR,
			RightSide: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{
								Name: LEFTPAREN,
							},
							{
								Name: LIST,
							},
							{
								Name: RIGHTPAREN,
							},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{
								Name: LEFTPAREN,
							},
							{
								Name: RIGHTPAREN,
							},
						},
					},
				},
			},
		},
		&grammar.Terminal{
			Name:   LEFTPAREN,
			Regexp: "\\(",
		},
		&grammar.Terminal{
			Name:   RIGHTPAREN,
			Regexp: "\\)",
		},
	}

	canonicalCollections := [][]*LrItem{
		{
			&LrItem{
				left: GOAL,
				right: []grammar.Category{
					LIST,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					LIST,
					PAIR,
				},
				lookahead: LEFTPAREN,
				position:  0,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					PAIR,
				},
				lookahead: LEFTPAREN,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					LIST,
					PAIR,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					PAIR,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					LIST,
					RIGHTPAREN,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					LIST,
					RIGHTPAREN,
				},
				lookahead: LEFTPAREN,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					RIGHTPAREN,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					RIGHTPAREN,
				},
				lookahead: LEFTPAREN,
				position:  0,
			},
		},
		{
			&LrItem{
				left: GOAL,
				right: []grammar.Category{
					LIST,
				},
				lookahead: token.EOF,
				position:  1,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					LIST,
					PAIR,
				},
				lookahead: token.EOF,
				position:  1,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					LIST,
					PAIR,
				},
				lookahead: LEFTPAREN,
				position:  1,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					LIST,
					RIGHTPAREN,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					LIST,
					RIGHTPAREN,
				},
				lookahead: LEFTPAREN,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					RIGHTPAREN,
				},
				lookahead: token.EOF,
				position:  0,
			},
			&LrItem{
				left: PAIR,
				right: []grammar.Category{
					LEFTPAREN,
					RIGHTPAREN,
				},
				lookahead: LEFTPAREN,
				position:  0,
			},
		},
	}

	tg := NewTableGenerator()

	action, gotoTable := tg.generateParseTables(productions)
	t.Logf("%v", action)
	t.Logf("%v", gotoTable)

	result, _ := tg.generateCanonicalCollection(productions)

	if len(result) != len(canonicalCollections) {
		t.Fatalf("canonical collections do not have same size. Expected %d. Got %d", len(canonicalCollections), len(result))
	}

	for i, canonicalCollection := range canonicalCollections {
		resultingCollection := result[i]

		t.Logf("actual collection:")
		stringifiedResult := getKeys(resultingCollection)
		slices.Sort(stringifiedResult)
		for _, lrItem := range stringifiedResult {
			t.Logf(lrItem)
		}

		t.Logf("expected collection:")
		stringifiedCanonicalCollection := mapToStringRepresentation(canonicalCollection)
		slices.Sort(stringifiedCanonicalCollection)
		for _, lrItem := range stringifiedCanonicalCollection {
			t.Logf(lrItem)
		}

		if resultingCollection == nil {
			t.Fatalf("expected to have canonical collection %d. Got nil", i)
		}

		if len(resultingCollection) != len(canonicalCollection) {
			t.Fatalf("canonical collection %d has wrong size. Expected %d. Got %d", i, len(canonicalCollection), len(resultingCollection))
		}

		for _, lrItem := range canonicalCollection {

			stringifiedLrItem := lrItem.String()
			_, ok := resultingCollection[stringifiedLrItem]
			if !ok {
				t.Fatalf("expected the following item in collection %d: %s. But was not present", i, stringifiedLrItem)
			}
		}
	}
}

func getKeys(sourceMap map[string]*LrItem) []string {
	keys := make([]string, 0)
	for key := range maps.Keys(sourceMap) {
		keys = append(keys, key)
	}

	return keys
}

func mapToStringRepresentation(lrItems []*LrItem) []string {
	stringified := make([]string, len(lrItems))
	for i, lrItem := range lrItems {
		stringified[i] = lrItem.String()
	}
	return stringified
}
