package scannergenerator

import "slices"

const EPSILON = "EPSILON"

type RegexpToNfaConverter struct {
	input    string
	position int
}

// maybe also return error
func (c *RegexpToNfaConverter) Convert() (result map[string]map[int][]int) {
	return c.parseExpression()
}

func (c *RegexpToNfaConverter) parseExpression() map[string]map[int][]int {
	leftExpr := c.prefixHandler()
	for c.position < len(c.input)-1 {
		c.position++
		leftExpr = c.parseInfixExpression(leftExpr)
	}
	return leftExpr
}

func (c *RegexpToNfaConverter) prefixHandler() map[string]map[int][]int {
	switch currentSymbol := string(c.input[c.position]); currentSymbol {
	case "(":
		return nil
	case "*":
		return nil
	default:
		return c.convertSingleSymbol()
	}
}

func (c *RegexpToNfaConverter) convertSingleSymbol() map[string]map[int][]int {
	currentSymbol := string(c.input[c.position])
	return map[string]map[int][]int{
		currentSymbol: {0: []int{1}},
	}
}

func (c *RegexpToNfaConverter) parseInfixExpression(leftExpr map[string]map[int][]int) map[string]map[int][]int {
	switch currentSymbol := string(c.input[c.position]); currentSymbol {
	case "|":
		return c.parseAlternation(leftExpr)
	default:
		return c.parseConcatenation(leftExpr)
	}
}

func (c *RegexpToNfaConverter) parseAlternation(leftExpr map[string]map[int][]int) map[string]map[int][]int {
	c.position++
	rightExpr := c.parseExpression()

	highestStateInLeftExpr := findHighestState(leftExpr)
	highestStateInRightExpr := findHighestState(rightExpr)
	lowestStateInRightExpr := findLowestState(rightExpr)

	for symbol, transitions := range rightExpr {
		if leftExpr[symbol] == nil {
			leftExpr[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+highestStateInLeftExpr+1)
			}
			leftExpr[symbol][stateFrom+highestStateInLeftExpr+1] = statesTo
		}
	}

	if leftExpr[EPSILON] == nil {
		leftExpr[EPSILON] = make(map[int][]int)
	}

	numberOfStatesInUnion := highestStateInLeftExpr + highestStateInRightExpr + 2
	leftExpr[EPSILON][numberOfStatesInUnion] = []int{highestStateInRightExpr + 1, lowestStateInRightExpr}
	leftExpr[EPSILON][highestStateInLeftExpr] = []int{numberOfStatesInUnion + 1}
	leftExpr[EPSILON][highestStateInLeftExpr+highestStateInRightExpr+1] = []int{numberOfStatesInUnion + 1}

	return leftExpr
}

func (c *RegexpToNfaConverter) parseConcatenation(leftExpr map[string]map[int][]int) map[string]map[int][]int {
	rightExpr := c.parseExpression()

	highestStateInLeftExpr := findHighestState(leftExpr)
	for symbol, transitions := range rightExpr {
		if leftExpr[symbol] == nil {
			leftExpr[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+highestStateInLeftExpr+1)
			}
			leftExpr[symbol][stateFrom+highestStateInLeftExpr+1] = statesTo
		}
	}

	if leftExpr[EPSILON] == nil {
		leftExpr[EPSILON] = make(map[int][]int)
	}
	leftExpr[EPSILON][highestStateInLeftExpr] = []int{highestStateInLeftExpr + 1}

	return leftExpr
}

func findHighestState(expr map[string]map[int][]int) int {
	highestState := 0
	for _, stateMappings := range expr {
		for stateFrom, stateTo := range stateMappings {
			if stateFrom > highestState {
				highestState = stateFrom
			}
			if slices.Max(stateTo) > highestState {
				highestState = slices.Max(stateTo)
			}
		}
	}
	return highestState
}

func findLowestState(expr map[string]map[int][]int) int {
	lowestState := 0
	for _, stateMappings := range expr {
		for stateFrom, stateTo := range stateMappings {
			if stateFrom < lowestState {
				lowestState = stateFrom
			}
			if slices.Min(stateTo) < lowestState {
				lowestState = slices.Min(stateTo)
			}
		}
	}
	return lowestState
}
