package scannergenerator

import "compiler/token"

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
	nfa := c.parseExpression(LOWEST)
	if nfa.TypeTable == nil {
		nfa.TypeTable = make(map[int]token.TokenType)
	}

	return nfa
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
	case "[":
		return c.parseRange()
	case "\\":
		return c.parseEscapedSymbol()
	default:
		return c.parseSingleSymbol()
	}
}

func (c *RegexpToNfaConverter) parseEscapedSymbol() *Nfa {
	c.position++
	return c.parseSingleSymbol()
}

func (c *RegexpToNfaConverter) parseRange() *Nfa {
	c.position++
	lowerBound := c.regexp[c.position]
	c.position += 2
	upperBound := c.regexp[c.position]
	c.position += 2

	if lowerBound > upperBound {
		panic("Lower bound is greater than upper bound")
	}

	var nfa *Nfa
	for i := lowerBound; i <= upperBound; i++ {
		symbolNfa := NfaFromSingleSymbol(string(i))
		if nfa == nil {
			nfa = symbolNfa
		} else {
			nfa = nfa.Union(symbolNfa)
		}
	}
	return nfa
}

func (c *RegexpToNfaConverter) parseParenthesis() *Nfa {
	c.position++
	nfa := c.parseExpression(LOWEST)
	c.position += 2
	return nfa
}

func (c *RegexpToNfaConverter) parseSingleSymbol() *Nfa {
	currentSymbol := string(c.regexp[c.position])
	return NfaFromSingleSymbol(currentSymbol)
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
