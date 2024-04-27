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
	var statements []ast.Statement

	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		statements = append(statements, statement)
	}

	return &ast.Program{Statements: statements}
}

func (p *Parser) parseExpression() ast.Expression {
	return &ast.Identifier{}
}

func (p *Parser) parseStatement() ast.Statement {

	lstmt := &ast.LetStatement{Token: p.currentToken}
	p.nextToken()
	lstmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	p.nextToken()
	lstmt.Value = p.parseExpression()
	p.nextToken()
	p.nextToken()
	p.nextToken()

	return lstmt
}