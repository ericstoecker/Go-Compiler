package ast

import "compiler/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (lst *LetStatement) TokenLiteral() string {
	return lst.Token.Literal
}
func (lst *LetStatement) statementNode() {}

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}
func (ident *Identifier) expressionNode() {}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStmt *ReturnStatement) TokenLiteral() string {
	return returnStmt.Token.Literal
}
func (returnStmt *ReturnStatement) statementNode() {}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStmt *ExpressionStatement) TokenLiteral() string {
	return expressionStmt.Token.Literal
}
func (expressionStmt *ExpressionStatement) statementNode() {}

type IntegerExpression struct {
	Token token.Token
	Value int64
}

func (intExpression *IntegerExpression) TokenLiteral() string {
	return intExpression.Token.Literal
}
func (intExpression *IntegerExpression) expressionNode() {}
