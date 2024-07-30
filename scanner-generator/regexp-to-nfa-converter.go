package scannergenerator

import (
	"compiler/token"
	"fmt"
)

const EPSILON = "EPSILON"

const (
	_ int = iota
	CLOSING
	LOWEST
	ALTERNATION
	CONCATENATION
	KLEENE
	PARENTHESIS
)

type RegexpToNfaConverter struct {
	regexp   string
	position int

	ch byte

	precedences map[string]int
}

func NewRegexpToNfaConverter(input string) *RegexpToNfaConverter {
	precedences := map[string]int{
		"|": ALTERNATION,
		"*": KLEENE,
		"(": PARENTHESIS,
		")": CLOSING,
	}

	converter := &RegexpToNfaConverter{
		regexp:      input,
		position:    -1,
		precedences: precedences,
	}
	converter.readCharacter()

	return converter
}

func (c *RegexpToNfaConverter) Convert() (*Nfa, error) {
	nfa, err := c.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if c.position < len(c.regexp)-1 && c.regexp[c.position+1] == ')' {
		return nil, fmt.Errorf("expected opening ')'")
	}

	if nfa.TypeTable == nil {
		nfa.TypeTable = make(map[int]token.TokenType)
	}
	return nfa, nil
}

func (c *RegexpToNfaConverter) parseExpression(precedence int) (*Nfa, error) {
	left, err := c.prefixHandler()
	if err != nil {
		return nil, err
	}

	for c.position < len(c.regexp)-1 && precedence < c.peekPrecedence() {
		c.readCharacter()

		left, err = c.parseInfixExpression(left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
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

func (c *RegexpToNfaConverter) prefixHandler() (*Nfa, error) {
	switch currentSymbol := c.ch; currentSymbol {
	case '(':
		return c.parseParenthesis()
	case '[':
		return c.parseRange()
	case '\\':
		return c.parseEscapedSymbol(), nil
	case 0:
		return nil, nil
	default:
		return c.parseSingleSymbol(), nil
	}
}

func (c *RegexpToNfaConverter) parseEscapedSymbol() *Nfa {
	c.readCharacter()
	return c.parseSingleSymbol()
}

func (c *RegexpToNfaConverter) parseRange() (*Nfa, error) {
	c.readCharacter()
	lowerBound := c.ch

	c.readCharacter()
	if c.ch != '-' {
		return nil, fmt.Errorf("expected '-' in between range")
	}

	c.readCharacter()
	upperBound := c.ch
	if lowerBound >= upperBound {
		return nil, fmt.Errorf("lower bound greater or equal to upper bound '[%s-%s]'", string(lowerBound), string(upperBound))
	}

	c.readCharacter()
	if c.ch != ']' {
		return nil, fmt.Errorf("expected closing ']' for range")
	}

	nfasForSymbols := make([]*Nfa, upperBound-lowerBound+1)
	for i := lowerBound; i <= upperBound; i++ {
		symbolNfa := NfaFromSingleSymbol(string(i))
		nfasForSymbols[i-lowerBound] = symbolNfa
	}

	return nfasForSymbols[0].Union(nfasForSymbols[1:]...), nil
}

func (c *RegexpToNfaConverter) parseParenthesis() (*Nfa, error) {
	c.readCharacter()

	nfa, err := c.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	c.readCharacter()
	if string(c.ch) != ")" {
		return nil, fmt.Errorf("expected closing ')'")
	}

	return nfa, nil
}

func (c *RegexpToNfaConverter) parseSingleSymbol() *Nfa {
	currentSymbol := string(c.ch)
	return NfaFromSingleSymbol(currentSymbol)
}

func (c *RegexpToNfaConverter) parseInfixExpression(left *Nfa) (*Nfa, error) {
	switch currentSymbol := string(c.regexp[c.position]); currentSymbol {
	case "|":
		return c.parseAlternation(left)
	case "*":
		return c.parseKleeneStar(left), nil
	default:
		return c.parseConcatenation(left)
	}
}

func (c *RegexpToNfaConverter) parseKleeneStar(left *Nfa) *Nfa {
	return left.Kleene()
}

func (c *RegexpToNfaConverter) parseAlternation(left *Nfa) (*Nfa, error) {
	c.readCharacter()
	right, err := c.parseExpression(ALTERNATION)
	if err != nil {
		return nil, err
	}
	if right == nil {
		return nil, fmt.Errorf("expected right side of |")
	}

	return left.Union(right), nil
}

func (c *RegexpToNfaConverter) parseConcatenation(left *Nfa) (*Nfa, error) {
	right, err := c.parseExpression(CONCATENATION)
	if err != nil {
		return nil, err
	}

	return left.Concatenation(right), nil
}

func (c *RegexpToNfaConverter) readCharacter() {
	if c.position >= len(c.regexp)-1 {
		c.ch = 0
	} else {
		c.position++
		c.ch = c.regexp[c.position]
	}
}
