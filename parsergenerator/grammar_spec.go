package parsergenerator

import (
	"compiler/ast"
	"compiler/grammar"
	"compiler/token"
	"fmt"
	"strconv"
)

// GrammarSpec returns the grammar productions with AST construction handlers.
func GrammarSpec() []grammar.Production {
	productions := []grammar.Production{}

	// Program -> Statements
	prod := grammar.NewNonTerminal("Program", &grammar.Identifier{Name: "Statements"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return &ast.Program{Statements: nodes[0].([]ast.Statement)}
	}
	productions = append(productions, prod)

	// Statements -> Statement Statements
	prod = grammar.NewNonTerminal("Statements", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "Statement"},
			{Name: "Statements"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		stmts := []ast.Statement{nodes[0].(ast.Statement)}
		stmts = append(stmts, nodes[1].([]ast.Statement)...)
		return stmts
	}
	productions = append(productions, prod)

	// Statements -> ε
	prod = grammar.NewNonTerminal("Statements", &grammar.Empty{}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return []ast.Statement{}
	}
	productions = append(productions, prod)

	// Statement -> LetStatement
	prod = grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "LetStatement"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0].(ast.Statement)
	}
	productions = append(productions, prod)

	// Statement -> ReturnStatement
	prod = grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "ReturnStatement"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0].(ast.Statement)
	}
	productions = append(productions, prod)

	// Statement -> ExpressionStatement
	prod = grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "ExpressionStatement"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0].(ast.Statement)
	}
	productions = append(productions, prod)

	// LetStatement -> 'let' IDENT '=' Expression ';'
	prod = grammar.NewNonTerminal("LetStatement", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "let"},
			{Name: "IDENT"},
			{Name: "="},
			{Name: "Expression"},
			{Name: ";"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		ident := nodes[1].(*ast.Identifier)
		return &ast.LetStatement{
			Token: nodes[0].(token.Token),
			Name:  ident,
			Value: nodes[3].(ast.Expression),
		}
	}
	productions = append(productions, prod)

	// ReturnStatement -> 'return' ExpressionOpt ';'
	prod = grammar.NewNonTerminal("ReturnStatement", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "return"},
			{Name: "ExpressionOpt"},
			{Name: ";"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		var returnValue ast.Expression
		if nodes[1] != nil {
			returnValue = nodes[1].(ast.Expression)
		}
		return &ast.ReturnStatement{
			Token:       nodes[0].(token.Token),
			ReturnValue: returnValue,
		}
	}
	productions = append(productions, prod)

	// ExpressionOpt -> Expression
	prod = grammar.NewNonTerminal("ExpressionOpt", &grammar.Identifier{Name: "Expression"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}
	productions = append(productions, prod)

	// ExpressionOpt -> ε
	prod = grammar.NewNonTerminal("ExpressionOpt", &grammar.Empty{}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nil
	}
	productions = append(productions, prod)

	// ExpressionStatement -> Expression ';'
	prod = grammar.NewNonTerminal("ExpressionStatement", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "Expression"},
			{Name: ";"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		expr := nodes[0].(ast.Expression)
		return &ast.ExpressionStatement{
			Token:      expr.TokenLiteral(),
			Expression: expr,
		}
	}
	productions = append(productions, prod)

	// Expression Precedence Levels
	// Defining expressions with explicit levels
	// Expression -> AssignmentExpression
	prod = grammar.NewNonTerminal("Expression", &grammar.Identifier{Name: "AssignmentExpression"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}
	productions = append(productions, prod)

	// AssignmentExpression -> LogicalOrExpression
	prod = grammar.NewNonTerminal("AssignmentExpression", &grammar.Identifier{Name: "LogicalOrExpression"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}
	productions = append(productions, prod)

	// LogicalOrExpression -> LogicalAndExpression ('||' LogicalAndExpression)*
	prod = grammar.NewNonTerminal("LogicalOrExpression", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "LogicalAndExpression"},
			{Name: "LogicalOrExpressionPrime"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		left := nodes[0].(ast.Expression)
		rightNodes := nodes[1].([]ast.Node)
		for i := 0; i < len(rightNodes); i += 2 {
			operator := rightNodes[i].(token.Token)
			right := rightNodes[i+1].(ast.Expression)
			left = &ast.InfixExpression{
				Token:    operator,
				Operator: operator.Literal,
				Left:     left,
				Right:    right,
			}
		}
		return left
	}
	productions = append(productions, prod)

	// LogicalOrExpressionPrime -> ('||' LogicalAndExpression)*
	prod = grammar.NewNonTerminal("LogicalOrExpressionPrime", &grammar.ZeroOrMore{
		Item: &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "||"},
				{Name: "LogicalAndExpression"},
			},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes
	}
	productions = append(productions, prod)

	// LogicalAndExpression -> EqualityExpression ('&&' EqualityExpression)*
	prod = grammar.NewNonTerminal("LogicalAndExpression", &grammar.Sequence{
		Items: []*grammar.Identifier{
			{Name: "EqualityExpression"},
			{Name: "LogicalAndExpressionPrime"},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		left := nodes[0].(ast.Expression)
		rightNodes := nodes[1].([]ast.Node)
		for i := 0; i < len(rightNodes); i += 2 {
			operator := rightNodes[i].(token.Token)
			right := rightNodes[i+1].(ast.Expression)
			left = &ast.InfixExpression{
				Token:    operator,
				Operator: operator.Literal,
				Left:     left,
				Right:    right,
			}
		}
		return left
	}
	productions = append(productions, prod)

	// LogicalAndExpressionPrime -> ('&&' EqualityExpression)*
	prod = grammar.NewNonTerminal("LogicalAndExpressionPrime", &grammar.ZeroOrMore{
		Item: &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "&&"},
				{Name: "EqualityExpression"},
			},
		},
	}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes
	}
	productions = append(productions, prod)

	// EqualityExpression -> RelationalExpression (('==' | '!=') RelationalExpression)*
	// Define similar productions for EqualityExpression and higher precedence expressions

	// PrimaryExpression -> Literal | Identifier | '(' Expression ')'
	// Etc.

	// Literal -> IntegerLiteral | BooleanLiteral | StringLiteral
	// IntegerLiteral -> INT
	prod = grammar.NewNonTerminal("IntegerLiteral", &grammar.Identifier{Name: "INT"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		tok := nodes[0].(token.Token)
		value, err := strconv.ParseInt(tok.Literal, 0, 64)
		if err != nil {
			panic(fmt.Sprintf("could not parse %q as integer", tok.Literal))
		}
		return &ast.IntegerLiteral{
			Token: tok,
			Value: value,
		}
	}
	productions = append(productions, prod)

	// BooleanLiteral -> 'true' | 'false'
	prod = grammar.NewNonTerminal("BooleanLiteral", &grammar.Identifier{Name: "BOOLEAN"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}
	productions = append(productions, prod)

	// Identifier -> IDENT
	prod = grammar.NewNonTerminal("Identifier", &grammar.Identifier{Name: "IDENT"}, 0)
	prod.Handler = func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}
	productions = append(productions, prod)

	// Add more productions to fully define the grammar as per your language specification.

	// Define all necessary terminals with their handlers
	productions = append(productions,
		&grammar.Terminal{
			Name:   "let",
			Regexp: "let",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.LET, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "return",
			Regexp: "return",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.RETURN, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "=",
			Regexp: "=",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.ASSIGN, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   ";",
			Regexp: ";",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.SEMICOLON, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "||",
			Regexp: "\\|\\|",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.OR, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "&&",
			Regexp: "&&",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.AND, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "==",
			Regexp: "==",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.EQUALS, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "!=",
			Regexp: "!=",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.NOT_EQUALS, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "IDENT",
			Regexp: "[a-zA-Z_][a-zA-Z0-9_]*",
			Handler: func(literal string) ast.Node {
				return &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: literal},
					Value: literal,
				}
			},
		},
		&grammar.Terminal{
			Name:   "INT",
			Regexp: "[0-9]+",
			Handler: func(literal string) ast.Node {
				value, err := strconv.ParseInt(literal, 0, 64)
				if err != nil {
					panic(fmt.Sprintf("could not parse %q as integer", literal))
				}
				return token.Token{Type: token.INT, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "true",
			Regexp: "true",
			Handler: func(literal string) ast.Node {
				return &ast.BooleanLiteral{
					Token: token.Token{Type: token.TRUE, Literal: literal},
					Value: true,
				}
			},
		},
		&grammar.Terminal{
			Name:   "false",
			Regexp: "false",
			Handler: func(literal string) ast.Node {
				return &ast.BooleanLiteral{
					Token: token.Token{Type: token.FALSE, Literal: literal},
					Value: false,
				}
			},
		},
		&grammar.Terminal{
			Name:   "+",
			Regexp: "\\+",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.PLUS, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "-",
			Regexp: "-",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.MINUS, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "*",
			Regexp: "\\*",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.ASTERISK, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "/",
			Regexp: "/",
			Handler: func(literal string) ast.Node {
				return token.Token{Type: token.SLASH, Literal: literal}
			},
		},
		&grammar.Terminal{
			Name:   "BOOLEAN",
			Regexp: "true|false",
			Handler: func(literal string) ast.Node {
				value := false
				if literal == "true" {
					value = true
				}
				return &ast.BooleanLiteral{
					Token: token.Token{Type: token.BOOLEAN, Literal: literal},
					Value: value,
				}
			},
		},
		// ... Add other terminals as needed ...
	)

	return productions
}
