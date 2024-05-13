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

type StringExpression struct {
	Token token.Token
	Value string
}

func (str *StringExpression) TokenLiteral() string {
	return str.Token.Literal
}
func (str *StringExpression) expressionNode() {}
func (str *StringExpression) String() string {
	return str.TokenLiteral()
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

	out.WriteString("(")

	out.WriteString(prefixExpr.TokenLiteral())
	out.WriteString(prefixExpr.Right.String())

	out.WriteString(")")

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

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ifExpression *IfExpression) TokenLiteral() string {
	return ifExpression.Token.Literal
}
func (ifExpression *IfExpression) expressionNode() {}
func (ifExpression *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ifExpression.Condition.String())
	out.WriteString(") ")

	out.WriteString(ifExpression.Consequence.String())

	if ifExpression.Alternative != nil {
		out.WriteString(" ")
		out.WriteString(ifExpression.Alternative.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bstmt *BlockStatement) TokenLiteral() string {
	return bstmt.Token.Literal
}
func (bstmt *BlockStatement) statementNode() {}
func (bstmt *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, s := range bstmt.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")

	return out.String()

}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fn *FunctionLiteral) TokenLiteral() string {
	return fn.Token.Literal
}
func (fn *FunctionLiteral) expressionNode() {}
func (fn *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for _, s := range fn.Parameters {
		out.WriteString(s.String())
	}
	out.WriteString(")")

	out.WriteString(fn.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Arguments []Expression
}

func (call *CallExpression) TokenLiteral() string {
	return call.Token.Literal
}
func (call *CallExpression) expressionNode() {}
func (call *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(call.TokenLiteral())
	out.WriteString("(")
	for i, arg := range call.Arguments {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")

	return out.String()
}

type ArrayExpression struct {
	Token    token.Token
	Elements []Expression
}

func (arr *ArrayExpression) TokenLiteral() string {
	return arr.Token.Literal
}
func (arr *ArrayExpression) expressionNode() {}
func (arr *ArrayExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	for i, e := range arr.Elements {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(e.String())
	}

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ind *IndexExpression) TokenLiteral() string {
	return ind.Token.Literal
}
func (ind *IndexExpression) expressionNode() {}
func (ind *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ind.Left.String())
	out.WriteString("[")
	out.WriteString(ind.Index.String())
	out.WriteString("]")

	return out.String()
}

type MapExpression struct {
	Token   token.Token
	Entries map[Expression]Expression
}

func (mapExpr *MapExpression) TokenLiteral() string {
	return mapExpr.Token.Literal
}
func (mapExpr *MapExpression) expressionNode() {}
func (mapExpr *MapExpression) String() string {
	var out bytes.Buffer

	out.WriteString(mapExpr.TokenLiteral())
	out.WriteString("{ ")
	for key, value := range mapExpr.Entries {
		out.WriteString(key.String())
		out.WriteString(": ")
		out.WriteString(value.String())
		out.WriteString(" ")
	}
	out.WriteString("}")

	return out.String()
}
