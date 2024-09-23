package parsergenerator

import (
	"compiler/grammar"
	"compiler/scanner"
	"compiler/token"
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
	dfa         *scanner.Dfa
}

func New(productions []grammar.Production) *LrParser {
	// Generate parse tables
	tg := NewTableGenerator()
	actionTable, gotoTable := tg.generateParseTables(productions)

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

func extractTokenClassifications(productions []grammar.Production) []token.TokenClassification {
	var tokenClassifications []token.TokenClassification
	precedence := 1

	for _, prod := range productions {
		switch t := prod.(type) {
		case *grammar.Terminal:
			tc := token.TokenClassification{
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
		},
		{
			"Goal",
			0,
		},
	}
	s := scanner.NewTableDrivenScanner(input, lr.dfa)
	token := s.NextToken()

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
		default:
			return fmt.Errorf("error when parsing input")
		}
	}

	return nil
}
