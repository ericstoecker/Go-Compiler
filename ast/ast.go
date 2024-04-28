package ast

import (
	"bytes"
	"compiler/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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
func (lst *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("let ")
	out.WriteString(lst.TokenLiteral())
	out.WriteString(" = ")
	out.WriteString(lst.Value.String())

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (ident *Identifier) TokenLiteral() string {
	return ident.Token.Literal
}
func (ident *Identifier) expressionNode() {}
func (ident *Identifier) String() string {
	return ident.TokenLiteral()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (returnStmt *ReturnStatement) TokenLiteral() string {
	return returnStmt.Token.Literal
}
func (returnStmt *ReturnStatement) statementNode() {}
func (returnStmt *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(returnStmt.TokenLiteral() + " ")
	out.WriteString(returnStmt.ReturnValue.String())

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (expressionStmt *ExpressionStatement) TokenLiteral() string {
	return expressionStmt.Token.Literal
}
func (expressionStmt *ExpressionStatement) statementNode() {}
func (expressionStmt *ExpressionStatement) String() string {
	if expressionStmt.Expression != nil {
		return expressionStmt.Expression.String()
	}
	return ""
}

type IntegerExpression struct {
	Token token.Token
	Value int64
}

func (intExpression *IntegerExpression) TokenLiteral() string {
	return intExpression.Token.Literal
}
func (intExpression *IntegerExpression) expressionNode() {}
func (intExpression *IntegerExpression) String() string {
	return intExpression.TokenLiteral()
}

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (boolExpression *BooleanExpression) TokenLiteral() string {
	return boolExpression.Token.Literal
}
func (boolExpression *BooleanExpression) expressionNode() {}
func (boolExpression *BooleanExpression) String() string {
	return boolExpression.TokenLiteral()
}

type PrefixExpression struct {
	Token    token.Token
	Operator token.TokenType
	Right    Expression
}

func (prefixExpr *PrefixExpression) TokenLiteral() string {
	return prefixExpr.Token.Literal
}
func (prefixExpr *PrefixExpression) expressionNode() {}
func (prefixExpr *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(prefixExpr.TokenLiteral() + " ")
	out.WriteString(prefixExpr.Right.String())

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator token.TokenType
	Right    Expression
}

func (infixExpr *InfixExpression) TokenLiteral() string {
	return infixExpr.Token.Literal
}
func (infixExpr *InfixExpression) expressionNode() {}
func (infixExpr *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")

	out.WriteString(infixExpr.Left.String())
	out.WriteString(" " + string(infixExpr.Operator) + " ")
	out.WriteString(infixExpr.Right.String())

	out.WriteString(")")

	return out.String()
}
