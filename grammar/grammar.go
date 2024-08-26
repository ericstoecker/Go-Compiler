package grammar

type Category string

type Production interface {
	production()
}

type NonTerminal struct {
	Name      Category
	RightSide RightSide
}

func NewNonTerminal(name Category, rightSide RightSide) *NonTerminal {
	return &NonTerminal{
		Name:      name,
		RightSide: rightSide,
	}
}

func (n *NonTerminal) production() {}

type Terminal struct {
	Name   Category
	Regexp string
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
