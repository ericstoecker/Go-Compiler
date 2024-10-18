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
	productions := []grammar.Production{
		// Program -> Statements
		grammar.NewNonTerminal("Program", &grammar.Identifier{Name: "Statements"}, func(nodes []ast.Node) ast.Node {
			return &ast.Program{Statements: nodes[0].([]ast.Statement)}
		}),

		// Statements -> Statement Statements
		grammar.NewNonTerminal("Statements", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "Statement"},
				{Name: "Statements"},
			},
		}, func(nodes []ast.Node) ast.Node {
			stmts := []ast.Statement{nodes[0].(ast.Statement)}
			stmts = append(stmts, nodes[1].([]ast.Statement)...)
			return stmts
		}),

		// Statements -> ε
		grammar.NewNonTerminal("Statements", &grammar.Empty{}, func(nodes []ast.Node) ast.Node {
			return []ast.Statement{}
		}),

		// Statement -> LetStatement
		grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "LetStatement"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Statement)
		}),

		// Statement -> ReturnStatement
		grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "ReturnStatement"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Statement)
		}),

		// Statement -> ExpressionStatement
		grammar.NewNonTerminal("Statement", &grammar.Identifier{Name: "ExpressionStatement"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Statement)
		}),

		// LetStatement -> 'let' IDENT '=' Expression ';'
		grammar.NewNonTerminal("LetStatement", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LET},
				{Name: token.IDENT},
				{Name: token.ASSIGN},
				{Name: "Expression"},
				{Name: token.SEMICOLON},
			},
		}, func(nodes []ast.Node) ast.Node {
			identToken := nodes[1].(token.Token)
			return &ast.LetStatement{
				Token: nodes[0].(token.Token),
				Name:  &ast.Identifier{Token: identToken, Value: identToken.Literal},
				Value: nodes[3].(ast.Expression),
			}
		}),

		// ReturnStatement -> 'return' Expression? ';'
		grammar.NewNonTerminal("ReturnStatement", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.RETURN},
				{Name: "ExpressionOpt"},
				{Name: token.SEMICOLON},
			},
		}, func(nodes []ast.Node) ast.Node {
			var returnValue ast.Expression
			if nodes[1] != nil {
				returnValue = nodes[1].(ast.Expression)
			}
			return &ast.ReturnStatement{
				Token:       nodes[0].(token.Token),
				ReturnValue: returnValue,
			}
		}),

		// ExpressionOpt -> Expression | ε
		grammar.NewNonTerminal("ExpressionOpt", &grammar.Identifier{Name: "Expression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("ExpressionOpt", &grammar.Empty{}, func(nodes []ast.Node) ast.Node {
			return nil
		}),

		// ExpressionStatement -> Expression ';'
		grammar.NewNonTerminal("ExpressionStatement", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "Expression"},
				{Name: token.SEMICOLON},
			},
		}, func(nodes []ast.Node) ast.Node {
			expr := nodes[0].(ast.Expression)
			return &ast.ExpressionStatement{
				Token:      expr.TokenLiteral(),
				Expression: expr,
			}
		}),

		// Expression rules with precedence levels (from lowest to highest)
		// Expression -> OrExpression
		grammar.NewNonTerminal("Expression", &grammar.Identifier{Name: "OrExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// OrExpression -> OrExpression '||' AndExpression
		grammar.NewNonTerminal("OrExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "OrExpression"},
				{Name: token.OR},
				{Name: "AndExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// OrExpression -> AndExpression
		grammar.NewNonTerminal("OrExpression", &grammar.Identifier{Name: "AndExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// AndExpression -> AndExpression '&&' EqualityExpression
		grammar.NewNonTerminal("AndExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "AndExpression"},
				{Name: token.AND},
				{Name: "EqualityExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// AndExpression -> EqualityExpression
		grammar.NewNonTerminal("AndExpression", &grammar.Identifier{Name: "EqualityExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// EqualityExpression -> EqualityExpression ( '==' | '!=' ) RelationalExpression
		grammar.NewNonTerminal("EqualityExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "EqualityExpression"},
				{Name: "EqOperator"},
				{Name: "RelationalExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// EqualityExpression -> RelationalExpression
		grammar.NewNonTerminal("EqualityExpression", &grammar.Identifier{Name: "RelationalExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// EqOperator -> '==' | '!='
		grammar.NewNonTerminal("EqOperator", &grammar.Identifier{Name: token.EQUALS}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("EqOperator", &grammar.Identifier{Name: token.NOT_EQUALS}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),

		// RelationalExpression -> RelationalExpression ( '<' | '>' | '<=' | '>=' ) AdditiveExpression
		grammar.NewNonTerminal("RelationalExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "RelationalExpression"},
				{Name: "RelOperator"},
				{Name: "AdditiveExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// RelationalExpression -> AdditiveExpression
		grammar.NewNonTerminal("RelationalExpression", &grammar.Identifier{Name: "AdditiveExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// RelOperator -> '<' | '>' | '<=' | '>='
		grammar.NewNonTerminal("RelOperator", &grammar.Identifier{Name: token.LT}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("RelOperator", &grammar.Identifier{Name: token.GT}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("RelOperator", &grammar.Identifier{Name: token.LESS_EQUAL}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("RelOperator", &grammar.Identifier{Name: token.GREATER_EQUAL}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),

		// AdditiveExpression -> AdditiveExpression ( '+' | '-' ) MultiplicativeExpression
		grammar.NewNonTerminal("AdditiveExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "AdditiveExpression"},
				{Name: "AddOperator"},
				{Name: "MultiplicativeExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// AdditiveExpression -> MultiplicativeExpression
		grammar.NewNonTerminal("AdditiveExpression", &grammar.Identifier{Name: "MultiplicativeExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// AddOperator -> '+' | '-'
		grammar.NewNonTerminal("AddOperator", &grammar.Identifier{Name: token.PLUS}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("AddOperator", &grammar.Identifier{Name: token.MINUS}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),

		// MultiplicativeExpression -> MultiplicativeExpression ( '*' | '/' ) UnaryExpression
		grammar.NewNonTerminal("MultiplicativeExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "MultiplicativeExpression"},
				{Name: "MulOperator"},
				{Name: "UnaryExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			left := nodes[0].(ast.Expression)
			operator := nodes[1].(token.Token)
			right := nodes[2].(ast.Expression)
			return &ast.InfixExpression{
				Token:    operator,
				Left:     left,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// MultiplicativeExpression -> UnaryExpression
		grammar.NewNonTerminal("MultiplicativeExpression", &grammar.Identifier{Name: "UnaryExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// MulOperator -> '*' | '/'
		grammar.NewNonTerminal("MulOperator", &grammar.Identifier{Name: token.ASTERISK}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("MulOperator", &grammar.Identifier{Name: token.SLASH}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),

		// UnaryExpression -> ( '!' | '-' ) UnaryExpression
		grammar.NewNonTerminal("UnaryExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "UnaryOperator"},
				{Name: "UnaryExpression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			operator := nodes[0].(token.Token)
			right := nodes[1].(ast.Expression)
			return &ast.PrefixExpression{
				Token:    operator,
				Operator: operator.Literal,
				Right:    right,
			}
		}),

		// UnaryExpression -> CallExpression
		grammar.NewNonTerminal("UnaryExpression", &grammar.Identifier{Name: "CallExpression"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// UnaryOperator -> '!' | '-'
		grammar.NewNonTerminal("UnaryOperator", &grammar.Identifier{Name: token.BANG}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("UnaryOperator", &grammar.Identifier{Name: token.MINUS}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),

		// CallExpression -> PrimaryExpression CallExpressionPrime
		grammar.NewNonTerminal("CallExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "PrimaryExpression"},
				{Name: "CallExpressionPrime"},
			},
		}, func(nodes []ast.Node) ast.Node {
			expr := nodes[0].(ast.Expression)
			suffixes := nodes[1].([]ast.Node)
			for _, suffix := range suffixes {
				switch s := suffix.(type) {
				case *ast.CallExpression:
					s.Left = expr
					expr = s
				case *ast.IndexExpression:
					s.Left = expr
					expr = s
				default:
					panic(fmt.Sprintf("Unexpected suffix type %T", s))
				}
			}
			return expr
		}),

		// CallExpressionPrime -> ( Arguments | IndexExpression )*
		grammar.NewNonTerminal("CallExpressionPrime", &grammar.ZeroOrMore{
			Item: &grammar.Choice{
				Items: []grammar.RightSide{
					&grammar.Identifier{Name: "Arguments"},
					&grammar.Identifier{Name: "IndexExpression"},
				},
			},
		}, func(nodes []ast.Node) ast.Node {
			return nodes
		}),

		// PrimaryExpression -> IDENT
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: token.IDENT}, func(nodes []ast.Node) ast.Node {
			tok := nodes[0].(token.Token)
			return &ast.Identifier{Token: tok, Value: tok.Literal}
		}),

		// PrimaryExpression -> INT
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: token.INT}, func(nodes []ast.Node) ast.Node {
			tok := nodes[0].(token.Token)
			value, err := strconv.ParseInt(tok.Literal, 0, 64)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse integer: %s", err))
			}
			return &ast.IntegerLiteral{Token: tok, Value: value}
		}),

		// PrimaryExpression -> STRING
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: token.STRING}, func(nodes []ast.Node) ast.Node {
			tok := nodes[0].(token.Token)
			return &ast.StringLiteral{Token: tok, Value: tok.Literal}
		}),

		// PrimaryExpression -> 'true' | 'false'
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: token.TRUE}, func(nodes []ast.Node) ast.Node {
			tok := nodes[0].(token.Token)
			return &ast.BooleanLiteral{Token: tok, Value: true}
		}),
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: token.FALSE}, func(nodes []ast.Node) ast.Node {
			tok := nodes[0].(token.Token)
			return &ast.BooleanLiteral{Token: tok, Value: false}
		}),

		// PrimaryExpression -> '(' Expression ')'
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LPAREN},
				{Name: "Expression"},
				{Name: token.RPAREN},
			},
		}, func(nodes []ast.Node) ast.Node {
			return nodes[1].(ast.Expression)
		}),

		// PrimaryExpression -> ArrayLiteral
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: "ArrayLiteral"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// PrimaryExpression -> MapLiteral
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: "MapLiteral"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// PrimaryExpression -> FunctionLiteral
		grammar.NewNonTerminal("PrimaryExpression", &grammar.Identifier{Name: "FunctionLiteral"}, func(nodes []ast.Node) ast.Node {
			return nodes[0].(ast.Expression)
		}),

		// ArrayLiteral -> '[' ExpressionList? ']'
		grammar.NewNonTerminal("ArrayLiteral", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LBRACKET},
				{Name: "ExpressionListOpt"},
				{Name: token.RBRACKET},
			},
		}, func(nodes []ast.Node) ast.Node {
			var elements []ast.Expression
			if nodes[1] != nil {
				elements = nodes[1].([]ast.Expression)
			}
			return &ast.ArrayLiteral{
				Token:    nodes[0].(token.Token),
				Elements: elements,
			}
		}),

		// MapLiteral -> '{' KeyValueList? '}'
		grammar.NewNonTerminal("MapLiteral", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LBRACE},
				{Name: "KeyValueListOpt"},
				{Name: token.RBRACE},
			},
		}, func(nodes []ast.Node) ast.Node {
			pairs := make(map[ast.Expression]ast.Expression)
			if nodes[1] != nil {
				pairs = nodes[1].(map[ast.Expression]ast.Expression)
			}
			return &ast.MapLiteral{
				Token:   nodes[0].(token.Token),
				Entries: pairs,
			}
		}),

		// FunctionLiteral -> 'fn' FunctionParameters BlockStatement
		grammar.NewNonTerminal("FunctionLiteral", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.FUNCTION},
				{Name: "FunctionParameters"},
				{Name: "BlockStatement"},
			},
		}, func(nodes []ast.Node) ast.Node {
			params := nodes[1].([]*ast.Identifier)
			body := nodes[2].(*ast.BlockStatement)
			return &ast.FunctionLiteral{
				Token:      nodes[0].(token.Token),
				Parameters: params,
				Body:       body,
			}
		}),

		// FunctionParameters -> '(' IdentifierList? ')'
		grammar.NewNonTerminal("FunctionParameters", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LPAREN},
				{Name: "IdentifierListOpt"},
				{Name: token.RPAREN},
			},
		}, func(nodes []ast.Node) ast.Node {
			var params []*ast.Identifier
			if nodes[1] != nil {
				params = nodes[1].([]*ast.Identifier)
			}
			return params
		}),

		// IdentifierListOpt -> IdentifierList | ε
		grammar.NewNonTerminal("IdentifierListOpt", &grammar.Identifier{Name: "IdentifierList"}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("IdentifierListOpt", &grammar.Empty{}, func(nodes []ast.Node) ast.Node {
			return nil
		}),

		// IdentifierList -> IDENT (',' IDENT)*
		grammar.NewNonTerminal("IdentifierList", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.IDENT},
				{Name: "IdentifierListPrime"},
			},
		}, func(nodes []ast.Node) ast.Node {
			params := []*ast.Identifier{
				{
					Token: nodes[0].(token.Token),
					Value: nodes[0].(token.Token).Literal,
				},
			}
			rest := nodes[1].([]*ast.Identifier)
			params = append(params, rest...)
			return params
		}),

		// IdentifierListPrime -> (',' IDENT)*
		grammar.NewNonTerminal("IdentifierListPrime", &grammar.ZeroOrMore{
			Item: &grammar.Sequence{
				Items: []*grammar.Identifier{
					{Name: token.COMMA},
					{Name: token.IDENT},
				},
			},
		}, func(nodes []ast.Node) ast.Node {
			idents := []*ast.Identifier{}
			for _, node := range nodes {
				pair := node.([]ast.Node)
				identToken := pair[1].(token.Token)
				idents = append(idents, &ast.Identifier{Token: identToken, Value: identToken.Literal})
			}
			return idents
		}),

		// ExpressionListOpt -> ExpressionList | ε
		grammar.NewNonTerminal("ExpressionListOpt", &grammar.Identifier{Name: "ExpressionList"}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("ExpressionListOpt", &grammar.Empty{}, func(nodes []ast.Node) ast.Node {
			return nil
		}),

		// ExpressionList -> Expression (',' Expression)*
		grammar.NewNonTerminal("ExpressionList", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "Expression"},
				{Name: "ExpressionListPrime"},
			},
		}, func(nodes []ast.Node) ast.Node {
			elements := []ast.Expression{nodes[0].(ast.Expression)}
			rest := nodes[1].([]ast.Expression)
			elements = append(elements, rest...)
			return elements
		}),

		// ExpressionListPrime -> (',' Expression)*
		grammar.NewNonTerminal("ExpressionListPrime", &grammar.ZeroOrMore{
			Item: &grammar.Sequence{
				Items: []*grammar.Identifier{
					{Name: token.COMMA},
					{Name: "Expression"},
				},
			},
		}, func(nodes []ast.Node) ast.Node {
			expressions := []ast.Expression{}
			for _, node := range nodes {
				pair := node.([]ast.Node)
				expr := pair[1].(ast.Expression)
				expressions = append(expressions, expr)
			}
			return expressions
		}),

		// KeyValueListOpt -> KeyValueList | ε
		grammar.NewNonTerminal("KeyValueListOpt", &grammar.Identifier{Name: "KeyValueList"}, func(nodes []ast.Node) ast.Node {
			return nodes[0]
		}),
		grammar.NewNonTerminal("KeyValueListOpt", &grammar.Empty{}, func(nodes []ast.Node) ast.Node {
			return nil
		}),

		// KeyValueList -> KeyValuePair (',' KeyValuePair)*
		grammar.NewNonTerminal("KeyValueList", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "KeyValuePair"},
				{Name: "KeyValueListPrime"},
			},
		}, func(nodes []ast.Node) ast.Node {
			pairs := make(map[ast.Expression]ast.Expression)
			firstPair := nodes[0].([]ast.Node)
			key := firstPair[0].(ast.Expression)
			value := firstPair[2].(ast.Expression)
			pairs[key] = value

			rest := nodes[1].([][]ast.Node)
			for _, pair := range rest {
				key = pair[1].(ast.Expression)
				value = pair[3].(ast.Expression)
				pairs[key] = value
			}
			return pairs
		}),

		// KeyValueListPrime -> (',' KeyValuePair)*
		grammar.NewNonTerminal("KeyValueListPrime", &grammar.ZeroOrMore{
			Item: &grammar.Sequence{
				Items: []*grammar.Identifier{
					{Name: token.COMMA},
					{Name: "KeyValuePair"},
				},
			},
		}, func(nodes []ast.Node) ast.Node {
			pairs := [][]ast.Node{}
			for _, node := range nodes {
				pair := node.([]ast.Node)
				pairs = append(pairs, pair)
			}
			return pairs
		}),

		// KeyValuePair -> Expression ':' Expression
		grammar.NewNonTerminal("KeyValuePair", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: "Expression"},
				{Name: token.COLON},
				{Name: "Expression"},
			},
		}, func(nodes []ast.Node) ast.Node {
			return nodes
		}),

		// BlockStatement -> '{' Statements '}'
		grammar.NewNonTerminal("BlockStatement", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LBRACE},
				{Name: "Statements"},
				{Name: token.RBRACE},
			},
		}, func(nodes []ast.Node) ast.Node {
			statements := nodes[1].([]ast.Statement)
			return &ast.BlockStatement{
				Token:      nodes[0].(token.Token),
				Statements: statements,
			}
		}),

		// IndexExpression -> '[' Expression ']'
		grammar.NewNonTerminal("IndexExpression", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LBRACKET},
				{Name: "Expression"},
				{Name: token.RBRACKET},
			},
		}, func(nodes []ast.Node) ast.Node {
			return &ast.IndexExpression{
				Token: nodes[0].(token.Token),
				Index: nodes[1].(ast.Expression),
			}
		}),

		// Arguments -> '(' ExpressionList? ')'
		grammar.NewNonTerminal("Arguments", &grammar.Sequence{
			Items: []*grammar.Identifier{
				{Name: token.LPAREN},
				{Name: "ExpressionListOpt"},
				{Name: token.RPAREN},
			},
		}, func(nodes []ast.Node) ast.Node {
			var args []ast.Expression
			if nodes[1] != nil {
				args = nodes[1].([]ast.Expression)
			} else {
				args = []ast.Expression{}
			}
			return &ast.CallExpression{
				Token:     nodes[0].(token.Token),
				Arguments: args,
			}
		}),

		// Add all necessary terminals with their regex patterns and handlers
		grammar.NewTerminal(token.LET, "let", nil),
		grammar.NewTerminal(token.RETURN, "return", nil),
		grammar.NewTerminal(token.FUNCTION, "fn", nil),
		grammar.NewTerminal(token.TRUE, "true", nil),
		grammar.NewTerminal(token.FALSE, "false", nil),
		grammar.NewTerminal(token.IDENT, "[a-zA-Z_][a-zA-Z0-9_]*", nil),
		grammar.NewTerminal(token.INT, "[0-9]+", nil),
		grammar.NewTerminal(token.STRING, "\"(\\\\.|[^\"])*\"", nil),
		grammar.NewTerminal(token.ASSIGN, "=", nil),
		grammar.NewTerminal(token.PLUS, "\\+", nil),
		grammar.NewTerminal(token.MINUS, "-", nil),
		grammar.NewTerminal(token.ASTERISK, "\\*", nil),
		grammar.NewTerminal(token.SLASH, "/", nil),
		grammar.NewTerminal(token.LT, "<", nil),
		grammar.NewTerminal(token.GT, ">", nil),
		grammar.NewTerminal(token.EQUALS, "==", nil),
		grammar.NewTerminal(token.NOT_EQUALS, "!=", nil),
		grammar.NewTerminal(token.GREATER_EQUAL, ">=", nil),
		grammar.NewTerminal(token.LESS_EQUAL, "<=", nil),
		grammar.NewTerminal(token.BANG, "!", nil),
		grammar.NewTerminal(token.AND, "&&", nil),
		grammar.NewTerminal(token.OR, "\\|\\|", nil),
		grammar.NewTerminal(token.COMMA, ",", nil),
		grammar.NewTerminal(token.COLON, ":", nil),
		grammar.NewTerminal(token.SEMICOLON, ";", nil),
		grammar.NewTerminal(token.LPAREN, "\\(", nil),
		grammar.NewTerminal(token.RPAREN, "\\)", nil),
		grammar.NewTerminal(token.LBRACE, "\\{", nil),
		grammar.NewTerminal(token.RBRACE, "\\}", nil),
		grammar.NewTerminal(token.LBRACKET, "\\[", nil),
		grammar.NewTerminal(token.RBRACKET, "\\]", nil),
	}

	return productions
}
