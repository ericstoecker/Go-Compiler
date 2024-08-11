package grammar

type GrammarRule interface {
	grammarRule()
}

type NonTerminal struct {
	Name      string
	RightSide RightSide
}

func (n *NonTerminal) grammarRule() {}

type Terminal struct {
	Name   string
	Regexp string
}

func (t *Terminal) grammarRule() {}

type RightSide interface {
	rightSide()
}

type Sequence struct {
}

func (s *Sequence) rightSide() {}

type Choice struct {
}

func (c *Choice) rightSide() {}

type Identifier struct {
	Name string
}

func (i *Identifier) rightSide() {}
