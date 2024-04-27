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

	expressionStatement.Expression = p.parseExpression()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return expressionStatement
}

func (p *Parser) parseExpression() ast.Expression {
	switch p.currentToken.Type {
	case token.IDENT:
		return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	case token.INT:
		value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to parse %s to int", p.currentToken.Literal)
			p.Errors = append(p.Errors, msg)
		}
		return &ast.IntegerExpression{Token: p.currentToken, Value: value}
	default:
		return nil
	}
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

func (p *Parser) peekTokenIs(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}

func (p *Parser) currentTokenIs(expected token.TokenType) bool {
	return p.currentToken.Type == expected
}
