package scannergenerator

const EPSILON = "EPSILON"

const (
	_ int = iota
	LOWEST
	ALTERNATION
	CONCATENATION
)

type RegexpToNfaConverter struct {
	input    string
	position int

	precedences map[string]int
}

func New(input string) *RegexpToNfaConverter {
	precedences := map[string]int{
		"|": ALTERNATION,
	}

	return &RegexpToNfaConverter{
		input:       input,
		position:    0,
		precedences: precedences,
	}
}

func (c *RegexpToNfaConverter) Convert() (result map[string]map[int][]int) {
	nfa := c.parseExpression(LOWEST)
	if nfa != nil {
		return nfa.Transitions
	}
	return nil
}

func (c *RegexpToNfaConverter) parseExpression(precedence int) *Nfa {
	left := c.prefixHandler()
	for c.position < len(c.input)-1 && precedence < c.peekPrecedence() {
		c.position++
		left = c.parseInfixExpression(left)
	}
	return left
}

func (c *RegexpToNfaConverter) peekPrecedence() int {
	if c.position+1 >= len(c.input) {
		return CONCATENATION
	}

	nextSymbol := string(c.input[c.position+1])
	if precedence, ok := c.precedences[nextSymbol]; ok {
		return precedence
	}
	return CONCATENATION
}

func (c *RegexpToNfaConverter) prefixHandler() *Nfa {
	switch currentSymbol := string(c.input[c.position]); currentSymbol {
	case "(":
		return nil
	case "*":
		return nil
	default:
		return c.convertSingleSymbol()
	}
}

func (c *RegexpToNfaConverter) convertSingleSymbol() *Nfa {
	currentSymbol := string(c.input[c.position])
	return &Nfa{
		Transitions: map[string]map[int][]int{
			currentSymbol: {0: []int{1}},
		},
		InitialState: 0,
		FinalState:   1,
	}
}

func (c *RegexpToNfaConverter) parseInfixExpression(left *Nfa) *Nfa {
	switch currentSymbol := string(c.input[c.position]); currentSymbol {
	case "|":
		return c.parseAlternation(left)
	default:
		return c.parseConcatenation(left)
	}
}

func (c *RegexpToNfaConverter) parseAlternation(left *Nfa) *Nfa {
	c.position++
	right := c.parseExpression(ALTERNATION)

	return left.union(right)
}

func (c *RegexpToNfaConverter) parseConcatenation(left *Nfa) *Nfa {
	right := c.parseExpression(CONCATENATION)

	return left.concatenation(right)
}
