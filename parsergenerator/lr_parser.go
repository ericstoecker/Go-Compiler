package parsergenerator

import (
	"compiler/ast"
	"compiler/grammar"
	"compiler/scanner"
	"compiler/token"
	"fmt"
	"slices"
)

type stackItem struct {
	category grammar.Category
	state    int
	node     ast.Node
}

type LrParser struct {
	actionTable map[int]map[grammar.Category]parseAction
	gotoTable   map[int]map[grammar.Category]int
	dfa         *scanner.Dfa
}

func New(productions []grammar.Production) *LrParser {
	// Generate parse tables
	actionTable, gotoTable := generateParseTables(productions)

	// Extract token classifications from terminals
	tokenClassifications := extractTokenClassifications(productions)

	// Generate DFA for the scanner
	sg := scanner.NewScannerGenerator()
	dfa := sg.GenerateScanner(tokenClassifications)

	return &LrParser{
		actionTable: actionTable,
		gotoTable:   gotoTable,
		dfa:         dfa,
	}
}

func extractTokenClassifications(productions []grammar.Production) []scanner.TokenClassification {
	var tokenClassifications []scanner.TokenClassification
	precedence := 1

	for _, prod := range productions {
		switch t := prod.(type) {
		case *grammar.Terminal:
			tc := scanner.TokenClassification{
				TokenType:  token.TokenType(t.Name),
				Regexp:     t.Regexp,
				Precedence: precedence,
			}
			tokenClassifications = append(tokenClassifications, tc)
		}
	}

	return tokenClassifications
}

func (lr *LrParser) Parse(input string) error {
	stack := []*stackItem{
		{
			"",
			-1,
			nil,
		},
		{
			"Goal",
			0,
			nil,
		},
	}
	s := scanner.NewTableDrivenScanner(input, lr.dfa)
	nextToken := s.NextToken()

	shouldContinue := true
	for shouldContinue {
		currentState := stack[len(stack)-1].state
		if currentState == -1 {
			return fmt.Errorf("error when parsing input: state was -1")
		}

		actionForStateAndType := lr.actionTable[currentState][grammar.Category(nextToken.Type)]
		if actionForStateAndType == nil {
			return fmt.Errorf("error when parsing input: no action for state %d and nextToken %s",
				currentState, nextToken.Type)
		}

		switch action := actionForStateAndType.(type) {
		case *reduce:
			var node ast.Node
			if hasNonTerminalHandler := action.lrItem.nonTerminalHandler != nil; hasNonTerminalHandler {
				nodesFromStack := extractNodes(stack, action.lenRightSide)
				node = action.lrItem.nonTerminalHandler(nodesFromStack)
			}

			stack = slices.Delete(stack, len(stack)-action.lenRightSide, len(stack))
			currentState := stack[len(stack)-1].state

			stack = append(stack, &stackItem{
				action.toCategory,
				lr.gotoTable[currentState][action.toCategory],
				node,
			})
		case *shift:
			var node ast.Node
			if hasTerminalHandler := action.lrItem.terminalHandler != nil; hasTerminalHandler {
				node = action.lrItem.terminalHandler(nextToken.Literal)
			}

			stack = append(stack, &stackItem{
				grammar.Category(nextToken.Type),
				action.toState,
				node,
			})

			nextToken = s.NextToken()
		case *accept:
			shouldContinue = false
		default:
			return fmt.Errorf("error when parsing input")
		}
	}

	return nil
}

func extractNodes(stack []*stackItem, from int) []ast.Node {
	relevantStackItems := stack[len(stack)-from : len(stack)]

	nodes := make([]ast.Node, len(relevantStackItems))
	for i, stackItem := range relevantStackItems {
		nodes[i] = stackItem.node
	}

	return nodes
}
