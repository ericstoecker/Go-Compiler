package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type PrefixParseFn func() ast.Expression
type InfixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	precedences          map[token.TokenType]int
	prefixParseFunctions map[token.TokenType]PrefixParseFn
	infixParseFunctions  map[token.TokenType]InfixParseFn

	Errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.precedences = make(map[token.TokenType]int)
	p.precedences[token.EQUALS] = EQUALS
	p.precedences[token.NOT_EQUALS] = EQUALS
	p.precedences[token.GREATER_EQUAL] = EQUALS
	p.precedences[token.LESS_EQUAL] = EQUALS
	p.precedences[token.GT] = LESSGREATER
	p.precedences[token.LT] = LESSGREATER
	p.precedences[token.PLUS] = SUM
	p.precedences[token.MINUS] = SUM
	p.precedences[token.ASTERIK] = PRODUCT
	p.precedences[token.SLASH] = PRODUCT
	p.precedences[token.LPAREN] = CALL

	p.prefixParseFunctions = make(map[token.TokenType]PrefixParseFn)
	p.prefixParseFunctions[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFunctions[token.BANG] = p.parsePrefixExpression
	p.prefixParseFunctions[token.IDENT] = p.parseIdentifier
	p.prefixParseFunctions[token.INT] = p.parseInteger
	p.prefixParseFunctions[token.TRUE] = p.parseBoolean
	p.prefixParseFunctions[token.FALSE] = p.parseBoolean
	p.prefixParseFunctions[token.IF] = p.parseIfExpression
	p.prefixParseFunctions[token.FUNCTION] = p.parseFunctionExpression
	p.prefixParseFunctions[token.LPAREN] = p.parseParen

	p.infixParseFunctions = make(map[token.TokenType]InfixParseFn)
	p.infixParseFunctions[token.EQUALS] = p.parseInfixExpression
	p.infixParseFunctions[token.NOT_EQUALS] = p.parseInfixExpression
	p.infixParseFunctions[token.GREATER_EQUAL] = p.parseInfixExpression
	p.infixParseFunctions[token.LESS_EQUAL] = p.parseInfixExpression
	p.infixParseFunctions[token.GT] = p.parseInfixExpression
	p.infixParseFunctions[token.LT] = p.parseInfixExpression
	p.infixParseFunctions[token.PLUS] = p.parseInfixExpression
	p.infixParseFunctions[token.MINUS] = p.parseInfixExpression
	p.infixParseFunctions[token.ASTERIK] = p.parseInfixExpression
	p.infixParseFunctions[token.SLASH] = p.parseInfixExpression

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	p.Errors = make([]string, 0)
	var statements []ast.Statement

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		p.nextToken()
	}

	return &ast.Program{Statements: statements}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	expressionStatement := &ast.ExpressionStatement{Token: p.currentToken}

	expressionStatement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expressionStatement
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFunctions[p.currentToken.Type]
	if !ok {
		msg := fmt.Sprintf("No prefix parse function for token '%s' with literal '%s'", p.currentToken.Type, p.currentToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	leftExpr := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFunctions[p.peekToken.Type]
		if !ok {
			return leftExpr
		}
		p.nextToken()

		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExpression := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Type}

	p.nextToken()
	prefixExpression.Right = p.parseExpression(PREFIX)

	return prefixExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	if p.peekTokenIs(token.LPAREN) {
		return p.parseCallExpression()
	}
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Error when trying to parse %s to int", p.currentToken.Literal)
		p.Errors = append(p.Errors, msg)
	}
	return &ast.IntegerExpression{Token: p.currentToken, Value: value}
}

func (p *Parser) parseBoolean() ast.Expression {
	value, err := strconv.ParseBool(p.currentToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("Error when trying to parse %s to bool", p.currentToken.Literal)
		p.Errors = append(p.Errors, msg)
	}
	return &ast.BooleanExpression{Token: p.currentToken, Value: value}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExpr := &ast.InfixExpression{Left: left, Operator: p.currentToken.Type}

	precedence := p.currentPrecedence()
	p.nextToken()

	infixExpr.Right = p.parseExpression(precedence)

	return infixExpr
}

func (p *Parser) currentPrecedence() int {
	if currentPrecedence, ok := p.precedences[p.currentToken.Type]; ok {
		return currentPrecedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if currentPrecedence, ok := p.precedences[p.peekToken.Type]; ok {
		return currentPrecedence
	}
	return LOWEST
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		p.nextToken()
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	function := &ast.FunctionExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	function.Parameters = p.parseParameters()

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	function.Body = p.parseBlockStatement()

	return function
}

func (p *Parser) parseParen() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseCallExpression() ast.Expression {
	call := &ast.CallExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()

	arguments := make([]ast.Expression, 0)
	for p.currentToken.Type != token.RPAREN {
		argument := p.parseExpression(LOWEST)
		if argument != nil {
			arguments = append(arguments, argument)
		}
		p.nextToken()
	}
	call.Arguments = arguments

	return call
}

func (p *Parser) parseParameters() []*ast.Identifier {
	p.nextToken()

	params := make([]*ast.Identifier, 0)
	for p.peekToken.Type != token.RPAREN {
		param := p.parseIdentifier()
		if param != nil {
			params = append(params, param.(*ast.Identifier))
		}

		p.nextToken()
	}

	return params
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{Token: p.currentToken}
	p.nextToken()

	statements := make([]ast.Statement, 0)
	for p.currentToken.Type != token.RBRACE {
		statement := p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		p.nextToken()
	}
	blockStatement.Statements = statements

	return blockStatement
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	if !p.currentTokenIs(token.SEMICOLON) {
		returnStmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) expectPeek(expected token.TokenType) bool {
	if p.peekTokenIs(expected) {
		p.nextToken()
		return true
	}
	msg := fmt.Sprintf("Expected next token to be %s. Got '%s'", expected, p.peekToken)
	p.Errors = append(p.Errors, msg)
	return false
}

// func (p *Parser) Error() string {}

func (p *Parser) peekTokenIs(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}

func (p *Parser) currentTokenIs(expected token.TokenType) bool {
	return p.currentToken.Type == expected
}
