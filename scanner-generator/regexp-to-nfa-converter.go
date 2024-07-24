package scannergenerator

const EPSILON = "EPSILON"

type RegexpToNfaConverter struct {
	input    string
	position int
}

// maybe also return error
func (c *RegexpToNfaConverter) Convert() (result map[string]map[int]int) {
	return c.parseExpression()
}

func (c *RegexpToNfaConverter) parseExpression() map[string]map[int]int {
	leftExpr := c.prefixHandler()
	for c.position < len(c.input)-1 {
		c.position++
		leftExpr = c.parseConcatenation(leftExpr)
	}
	return leftExpr
}

func (c *RegexpToNfaConverter) prefixHandler() map[string]map[int]int {
	switch currentSymbol := string(c.input[c.position]); currentSymbol {
	case "(":
		return nil
	case "*":
		return nil
	case "|":
		return nil
	default:
		return c.convertSingleSymbol()
	}
}

func (c *RegexpToNfaConverter) convertSingleSymbol() map[string]map[int]int {
	currentSymbol := string(c.input[c.position])
	return map[string]map[int]int{
		currentSymbol: {0: 1},
	}
}

func (c *RegexpToNfaConverter) parseConcatenation(leftExpr map[string]map[int]int) map[string]map[int]int {
	rightExpr := c.parseExpression()

	highestStateInLeftExpr := findHighestState(leftExpr)
	for symbol, transitions := range rightExpr {
		if leftExpr[symbol] == nil {
			leftExpr[symbol] = make(map[int]int)
		}
		for stateFrom, stateTo := range transitions {
			leftExpr[symbol][stateFrom+highestStateInLeftExpr+1] = stateTo + highestStateInLeftExpr + 1
		}
	}

	if leftExpr[EPSILON] == nil {
		leftExpr[EPSILON] = make(map[int]int)
	}
	leftExpr[EPSILON][highestStateInLeftExpr] = highestStateInLeftExpr + 1

	return leftExpr
}

func findHighestState(expr map[string]map[int]int) int {
	highestState := 0
	for _, stateMappings := range expr {
		for stateFrom, stateTo := range stateMappings {
			if stateFrom > highestState {
				highestState = stateFrom
			}
			if stateTo > highestState {
				highestState = stateTo
			}
		}
	}
	return highestState
}
