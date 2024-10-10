package parsergenerator

import (
	"compiler/ast"
	"compiler/grammar"
	"compiler/token"
	"fmt"
	"strconv"
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

func TestAstConstruction(t *testing.T) {
	productions := []grammar.Production{
		&grammar.NonTerminal{
			Name: GOAL,
			RightSide: &grammar.Identifier{
				Name: "Sum",
			},
			Handler: func(nodes []ast.Node) ast.Node {
				if len(nodes) != 1 {
					panic(fmt.Sprintf("Expected 1 node. Got %d", len(nodes)))
				}
				return nodes[0]
			},
		},
		&grammar.NonTerminal{
			Name: "Sum",
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
							{Name: "Sum"},
							{Name: "plus"},
							{Name: "number"},
						},
					},
				},
			},
			Handler: func(nodes []ast.Node) ast.Node {
				if len(nodes) != 3 {
					panic(fmt.Sprintf("Expected 3 nodes. Got %d", len(nodes)))
				}

				leftExpression, ok := nodes[0].(ast.Expression)
				if !ok {
					panic(fmt.Sprintf("Expected node to be an Expression. Got %T", nodes[0]))
				}

				rightExpression, ok := nodes[2].(ast.Expression)
				if !ok {
					panic(fmt.Sprintf("Expected node to be an Expression. Got %T", nodes[2]))
				}

				return &ast.InfixExpression{
					Left:     leftExpression,
					Operator: token.PLUS,
					Right:    rightExpression,
				}
			},
		},
		&grammar.Terminal{
			Name:   "number",
			Regexp: "[0-9]",
			Handler: func(s string) ast.Node {
				value, err := strconv.ParseInt(s, 0, 64)
				if err != nil {
					msg := fmt.Sprintf("Error when trying to parse %s to int", s)
					panic(msg)
				}
				return &ast.IntegerLiteral{
					Token: token.Token{
						Type:    "number",
						Literal: s,
					},
					Value: value,
				}
			},
		},
		&grammar.Terminal{
			Name:   "plus",
			Regexp: "+",
		},
	}

	tests := []string{
		"1 + 2 + 3",
	}

	lrParser := New(productions)

	for _, tt := range tests {
		node, err := lrParser.Parse(tt)

		if err != nil {
			t.Fatalf("error when parsing valid input '%s': %v", tt, err)
		}

		if node == nil {
			t.Fatalf("Expected node to be non-nil")
		}

		infixExpression, ok := node.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected node to be an InfixExpression. Got %T", node)
		}

		leftExpression, ok := infixExpression.Left.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected left to be an InfixExpression. Got %T", infixExpression.Left)
		}

		leftExpressionSummand, ok := leftExpression.Left.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("Expected left to be an IntegerLiteral. Got %T", leftExpression.Left)
		}

		if leftExpressionSummand.Value != 1 {
			t.Fatalf("Expected left summand to be 1. Got %d", leftExpressionSummand.Value)
		}

		rightExpressionSummand, ok := leftExpression.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("Expected right to be an IntegerLiteral. Got %T", leftExpression.Right)
		}

		if rightExpressionSummand.Value != 2 {
			t.Fatalf("Expected right summand to be 2. Got %d", rightExpressionSummand.Value)
		}

		_, ok = infixExpression.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("Expected right to be an IntegerLiteral. Got %T", infixExpression.Right)
		}

		if infixExpression.Operator != token.PLUS {
			t.Fatalf("Expected operator to be '+'. Got %s", infixExpression.Operator)
		}
	}
}
