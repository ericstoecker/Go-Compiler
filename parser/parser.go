package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
	"fmt"
	"strconv"
)

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	Errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

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

	expressionStatement.Expression = p.parseExpression(0)

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

	// check peekToken and see if an infix is available
	for !p.peekTokenIs(token.SEMICOLON) && precedence == 0 && p.peekToken.Type == token.PLUS {
		// if so => pass left side to parseInfixExpression
		p.nextToken()

		// parse infix expression
		infixExpr := &ast.InfixExpression{Left: leftExpr, Operator: p.currentToken.Type}
		p.nextToken()
		infixExpr.Right = p.parseExpression(1)

		leftExpr = infixExpr
	}
	// if not => return left
	return leftExpr
}

func (p *Parser) parseInfixExpression(leftExpr ast.Expression) *ast.InfixExpression {
	return nil
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
