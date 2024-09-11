package parsergenerator

import (
	"compiler/grammar"
	"compiler/scanner"
	"fmt"
	"slices"
)

type stackItem struct {
	category grammar.Category
	state    int
}

type LrParser struct {
	actionTable map[int]map[grammar.Category]parseAction
	gotoTable   map[int]map[grammar.Category]int
}

func New(actionTable map[int]map[grammar.Category]parseAction, gotoTable map[int]map[grammar.Category]int) *LrParser {
	return &LrParser{
		actionTable: actionTable,
		gotoTable:   gotoTable,
	}
}

func (lr *LrParser) Parse(input string) error {
	stack := []*stackItem{
		{
			"",
			-1,
		},
		{
			"Goal",
			0,
		},
	}
	s := scanner.NewHandcodedScanner(input)
	token := s.NextToken()

	// todo add nil checks for table access
	shouldContinue := true
	for shouldContinue {
		currentState := stack[len(stack)-1].state
		if currentState == -1 {
			return fmt.Errorf("error when parsing input: state was -1")
		}

		actionForStateAndType := lr.actionTable[currentState][grammar.Category(token.Type)]
		if actionForStateAndType == nil {
			return fmt.Errorf("error when parsing input: no action for state %d and token %s",
				currentState, token.Type)
		}

		switch action := actionForStateAndType.(type) {
		case *reduce:
			stack = slices.Delete(stack, len(stack)-action.lenRightSide, len(stack))

			currentState := stack[len(stack)-1].state

			stack = append(stack, &stackItem{
				action.toCategory,
				lr.gotoTable[currentState][action.toCategory],
			})
		case *shift:
			stack = append(stack, &stackItem{
				grammar.Category(token.Type),
				action.toState,
			})

			token = s.NextToken()
		case *accept:
			shouldContinue = false
			// todo can this be removed?
		default:
			return fmt.Errorf("error when parsing input")
		}
	}

	return nil
}
