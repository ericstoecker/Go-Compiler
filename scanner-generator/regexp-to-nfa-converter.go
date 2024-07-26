package scannergenerator

const EPSILON = "EPSILON"

const (
	_ int = iota
	LOWEST
	ALTERNATION
	CONCATENATION
	KLEENE
	PARENTHESIS
)

type RegexpToNfaConverter struct {
	regexp   string
	position int

	precedences map[string]int
}

func NewRegexpToNfaConverter(input string) *RegexpToNfaConverter {
	precedences := map[string]int{
		"|": ALTERNATION,
		"*": KLEENE,
		"(": PARENTHESIS,
	}

	return &RegexpToNfaConverter{
		regexp:      input,
		position:    0,
		precedences: precedences,
	}
}

func (c *RegexpToNfaConverter) Convert() *Nfa {
	return c.parseExpression(LOWEST)
}

func (c *RegexpToNfaConverter) parseExpression(precedence int) *Nfa {
	left := c.prefixHandler()

	for c.position < len(c.regexp)-1 && precedence < c.peekPrecedence() {
		c.position++
		if c.regexp[c.position] == ')' {
			return left
		}
		left = c.parseInfixExpression(left)
	}
	return left
}

func (c *RegexpToNfaConverter) peekPrecedence() int {
	if c.position+1 >= len(c.regexp) {
		return CONCATENATION
	}

	nextSymbol := string(c.regexp[c.position+1])
	if precedence, ok := c.precedences[nextSymbol]; ok {
		return precedence
	}
	return CONCATENATION
}

func (c *RegexpToNfaConverter) prefixHandler() *Nfa {
	switch currentSymbol := string(c.regexp[c.position]); currentSymbol {
	case "(":
		return c.parseParenthesis()
	default:
		return c.parseSingleSymbol()
	}
}

func (c *RegexpToNfaConverter) parseParenthesis() *Nfa {
	c.position++
	nfa := c.parseExpression(LOWEST)
	c.position += 2
	return nfa
}

func (c *RegexpToNfaConverter) parseSingleSymbol() *Nfa {
	currentSymbol := string(c.regexp[c.position])
	return &Nfa{
		Transitions: map[string]map[int][]int{
			currentSymbol: {0: []int{1}},
		},
		InitialState:    0,
		AcceptingStates: []int{1},
	}
}

func (c *RegexpToNfaConverter) parseInfixExpression(left *Nfa) *Nfa {
	switch currentSymbol := string(c.regexp[c.position]); currentSymbol {
	case "|":
		return c.parseAlternation(left)
	case "*":
		return c.parseKleeneStar(left)
	default:
		return c.parseConcatenation(left)
	}
}

func (c *RegexpToNfaConverter) parseKleeneStar(left *Nfa) *Nfa {
	c.position++
	return left.Kleene()
}

func (c *RegexpToNfaConverter) parseAlternation(left *Nfa) *Nfa {
	c.position++
	right := c.parseExpression(ALTERNATION)

	return left.Union(right)
}

func (c *RegexpToNfaConverter) parseConcatenation(left *Nfa) *Nfa {
	right := c.parseExpression(CONCATENATION)

	return left.Concatenation(right)
}
