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

type PrefixParseFn func() *ast.PrefixExpression
type InfixParseFn func(ast.Expression) *ast.InfixExpression

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	precedences         map[token.TokenType]int
	prefixParseFn       map[token.TokenType]PrefixParseFn
	infixParseFunctions map[token.TokenType]InfixParseFn

	Errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	p.precedences = make(map[token.TokenType]int)
	p.precedences[token.PLUS] = SUM
	p.precedences[token.MINUS] = SUM

	p.infixParseFunctions = make(map[token.TokenType]InfixParseFn)
	p.infixParseFunctions[token.PLUS] = p.parseInfixExpression
	p.infixParseFunctions[token.MINUS] = p.parseInfixExpression

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
	var leftExpr ast.Expression

	switch p.currentToken.Type {
	case token.MINUS:
		prefixExpression := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Type}

		p.nextToken()
		prefixExpression.Right = p.parseExpression(1)

		leftExpr = prefixExpression
	case token.IDENT:
		leftExpr = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.INT:
		value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to parse %s to int", p.currentToken.Literal)
			p.Errors = append(p.Errors, msg)
		}
		leftExpr = &ast.IntegerExpression{Token: p.currentToken, Value: value}
	default:
		// this will need to be error NO PREFIX PARSER
		leftExpr = nil
	}

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

func (p *Parser) parseInfixExpression(left ast.Expression) *ast.InfixExpression {
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
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	if !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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
