package parsergenerator

import (
	"compiler/grammar"
	"testing"
)

func TestGeneratedLrParserOnValidInputs(t *testing.T) {
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

	tests := []string{
		"()",
		"(())",
		"(()())",
		"((()))",
		"((((((()))))))", // Test deeply nested structure
		"()()()",
		"(())()",
	}

	tg := &TableGenerator{}
	action, goTo := tg.generateParseTables(productions)

	lrParser := New(action, goTo)

	for _, tt := range tests {
		err := lrParser.Parse(tt)

		if err != nil {
			t.Fatalf("error when parsing valid input '%s': %v", tt, err)
		}
	}
}

func TestGeneratedLrParserForBasicIfGrammar(t *testing.T) {
	productions := []grammar.Production{
		&grammar.NonTerminal{
			Name: GOAL,
			RightSide: &grammar.Identifier{
				Name: IF,
			},
		},
		&grammar.NonTerminal{
			Name: IF,
			RightSide: &grammar.Sequence{
				Items: []*grammar.Identifier{
					{Name: IF},
					{Name: LPAREN},
					{Name: RPAREN},
					{Name: LBRACE},
					{Name: RBRACE},
				},
			},
		},
		&grammar.Terminal{Name: IF, Regexp: "if"},
		&grammar.Terminal{Name: LPAREN, Regexp: "\\("},
		&grammar.Terminal{Name: RPAREN, Regexp: "\\)"},
		&grammar.Terminal{Name: LBRACE, Regexp: "\\{"},
		&grammar.Terminal{Name: RBRACE, Regexp: "\\}"},
	}

	tests := []string{
		"if () {}",
	}

	tg := &TableGenerator{}
	action, goTo := tg.generateParseTables(productions)

	lrParser := New(action, goTo)

	for _, tt := range tests {
		err := lrParser.Parse(tt)

		if err != nil {
			t.Fatalf("error when parsing valid input '%s': %v", tt, err)
		}
	}
}
