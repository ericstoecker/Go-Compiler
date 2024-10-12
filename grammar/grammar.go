package grammar

import "compiler/ast"

type Category string

type Production interface {
	production()
}

type NonTerminal struct {
	Name       Category
	RightSide  RightSide
	Handler    func([]ast.Node) ast.Node
	Precedence int
}

func NewNonTerminal(name Category, rightSide RightSide, precedence int) *NonTerminal {
	return &NonTerminal{
		Name:       name,
		RightSide:  rightSide,
		Precedence: precedence,
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
	Precedence() int
}

type Sequence struct {
	Items      []*Identifier
	precedence int
}

func NewSequence(items []*Identifier, precedence int) *Sequence {
	return &Sequence{
		Items:      items,
		precedence: precedence,
	}
}

func (s *Sequence) Precedence() int {
	return s.precedence
}

func (s *Sequence) rightSide() {}

type Choice struct {
	Items []RightSide
}

func (c *Choice) Precedence() int {
	return 0
}

func (c *Choice) rightSide() {}

type Identifier struct {
	Name Category
}

func (i *Identifier) rightSide() {}

func (i *Identifier) Precedence() int {
	return 0
}
