package parsergenerator

import (
	"compiler/grammar"
	"compiler/token"
	"testing"
)

const (
	GOAL       = "Goal"
	LIST       = "List"
	PAIR       = "Pair"
	LEFTPAREN  = "LeftParen"
	RIGHTPAREN = "RightParen"
)

func TestTableConstruction(t *testing.T) {
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
							&grammar.Identifier{
								Name: LIST,
							},
							&grammar.Identifier{
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
							&grammar.Identifier{
								Name: LEFTPAREN,
							},
							&grammar.Identifier{
								Name: LIST,
							},
							&grammar.Identifier{
								Name: RIGHTPAREN,
							},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							&grammar.Identifier{
								Name: LEFTPAREN,
							},
							&grammar.Identifier{
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
				lookahead: token.LPAREN,
				position:  0,
			},
			&LrItem{
				left: LIST,
				right: []grammar.Category{
					PAIR,
				},
				lookahead: token.LPAREN,
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
				lookahead: token.LPAREN,
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
				lookahead: token.LPAREN,
				position:  0,
			},
		},
	}

	tg := NewTableGenerator()
	result := tg.generateCanonicalCollection(productions)

	if len(result) != len(canonicalCollections) {
		t.Fatalf("canonical collections do not have same size. Expected %d. Got %d", len(canonicalCollections), len(result))
	}

	for i, canonicalCollection := range canonicalCollections {
		resultingCollection := result[i]

		t.Logf("actual collection %d: %v", i, resultingCollection)
		t.Logf("expected collection %d: %v", i, canonicalCollection)

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
