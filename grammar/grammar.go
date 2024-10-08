package grammar

import "compiler/ast"

type Category string

type Production interface {
	production()
}

type NonTerminal struct {
	Name      Category
	RightSide RightSide
	Handler   func([]ast.Node) ast.Node
}

func NewNonTerminal(name Category, rightSide RightSide) *NonTerminal {
	return &NonTerminal{
		Name:      name,
		RightSide: rightSide,
	}
}

func (n *NonTerminal) production() {}

type Terminal struct {
	Name    Category
	Regexp  string
	Handler func(string) ast.Node
}

func (t *Terminal) production() {}

type RightSide interface {
	rightSide()
}

type Sequence struct {
	Items []*Identifier
}

func (s *Sequence) rightSide() {}

type Choice struct {
	Items []RightSide
}

func (c *Choice) rightSide() {}

type Identifier struct {
	Name Category
}

func (i *Identifier) rightSide() {}
