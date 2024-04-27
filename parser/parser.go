package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/token"
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
		statements = append(statements, statement)
	}

	return &ast.Program{Statements: statements}
}

func (p *Parser) parseExpression() ast.Expression {
	// switch p.currentToken.Type {
	// case token.INT:
	// 	return &ast.Identifier{}
	// }
	p.nextToken()
	return &ast.Identifier{}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}
	p.nextToken()
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		p.Errors = append(p.Errors, "Expected a "+string(token.ASSIGN))
	}
	p.nextToken()

	p.nextToken()

	if p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	if p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) expectPeek(expected token.TokenType) bool {
	if p.peekTokenIs(expected) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}

func (p *Parser) currentTokenIs(expected token.TokenType) bool {
	return p.currentToken.Type == expected
}
