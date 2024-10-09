package parsergenerator

import (
	"compiler/grammar"
	"testing"
)

const (
	GOAL       = "Goal"
	LIST       = "List"
	PAIR       = "Pair"
	LEFTPAREN  = "("
	RIGHTPAREN = ")"
	LPAREN     = "("
	RPAREN     = ")"
	IF         = "IF"
	LBRACE     = "{"
	RBRACE     = "}"
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

	lrParser := New(productions)

	for _, tt := range tests {
		_, err := lrParser.Parse(tt)

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
				Name: "IfElse",
			},
		},
		&grammar.NonTerminal{
			Name: "IfElse",
			RightSide: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Identifier{
						Name: IF,
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: IF},
							{Name: "ElseIf"},
						},
					},
				},
			},
		},
		&grammar.NonTerminal{
			Name: "ElseIf",
			RightSide: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "else"},
							{Name: LBRACE},
							{Name: RBRACE},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "else"},
							{Name: "if"},
							{Name: LPAREN},
							{Name: RPAREN},
							{Name: LBRACE},
							{Name: RBRACE},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "else"},
							{Name: "if"},
							{Name: LPAREN},
							{Name: RPAREN},
							{Name: LBRACE},
							{Name: RBRACE},
							{Name: "ElseIf"},
						},
					},
				},
			},
		},
		&grammar.NonTerminal{
			Name: IF,
			RightSide: &grammar.Sequence{
				Items: []*grammar.Identifier{
					{Name: "if"},
					{Name: LPAREN},
					{Name: RPAREN},
					{Name: LBRACE},
					{Name: RBRACE},
				},
			},
		},
		&grammar.Terminal{Name: "if", Regexp: "if"},
		&grammar.Terminal{Name: "else", Regexp: "else"},
		&grammar.Terminal{Name: LPAREN, Regexp: "\\("},
		&grammar.Terminal{Name: RPAREN, Regexp: "\\)"},
		&grammar.Terminal{Name: LBRACE, Regexp: "\\{"},
		&grammar.Terminal{Name: RBRACE, Regexp: "\\}"},
	}

	tests := []string{
		"if () {}",
		"if () {} else {}",
		"if () {} else if () {}",
		"if () {} else if () {} else {}",
		"if () {} else if () {} else if () {} else {}",
	}

	lrParser := New(productions)

	for _, tt := range tests {
		_, err := lrParser.Parse(tt)

		if err != nil {
			t.Fatalf("error when parsing valid input '%s': %v", tt, err)
		}
	}
}

func TestAmbiguousGrammar(t *testing.T) {
	productions := []grammar.Production{
		&grammar.NonTerminal{
			Name: GOAL,
			RightSide: &grammar.Identifier{
				Name: "sum",
			},
		},
		&grammar.NonTerminal{
			Name: "sum",
			RightSide: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "number"},
							{Name: "plus"},
							{Name: "number"},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "sum"},
							{Name: "plus"},
							{Name: "number"},
						},
					},
					&grammar.Sequence{
						Items: []*grammar.Identifier{
							{Name: "number"},
							{Name: "plus"},
							{Name: "sum"},
						},
					},
				},
			},
		},
		&grammar.Terminal{Name: "number", Regexp: "[0-9]"},
		&grammar.Terminal{Name: "plus", Regexp: "+"},
	}

	tests := []string{
		"1 + 2",
	}

	lrParser := New(productions)

	for _, tt := range tests {
		_, err := lrParser.Parse(tt)

		if err != nil {
			t.Fatalf("error when parsing valid input '%s': %v", tt, err)
		}
	}
}
