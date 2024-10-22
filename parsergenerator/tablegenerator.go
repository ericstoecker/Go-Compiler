package parsergenerator

import (
	"compiler/ast"
	"compiler/grammar"
	"compiler/token"
	"fmt"
	"slices"
)

type parseAction interface {
	parseAction()
}

type shift struct {
	toState    int
	lrItem     *LrItem
	handler    func(string) ast.Node
	precedence int
}

func (s *shift) parseAction() {}

type reduce struct {
	toCategory   grammar.Category
	lenRightSide int
	lrItem       *LrItem
	handler      func([]ast.Node) ast.Node
	precedence   int
}

func (r *reduce) parseAction() {}

type accept struct {
	lenRightSide int
	handler      func([]ast.Node) ast.Node
}

func (a *accept) parseAction() {}

func generateParseTables(productions []grammar.Production) (actionTable map[int]map[grammar.Category]parseAction, gotoTable map[int]map[grammar.Category]int) {
	if len(productions) == 0 {
		panic("no productions defined")
	}
	goalProduction := productions[0].(*grammar.NonTerminal).Name

	canonicalCollections, gotoTable := generateCanonicalCollection(productions)
	actionTable = make(map[int]map[grammar.Category]parseAction)

	for i, collection := range canonicalCollections {
		if actionTable[i] == nil {
			actionTable[i] = make(map[grammar.Category]parseAction)
		}

		keys := make([]string, 0, len(collection))
		for key := range collection {
			keys = append(keys, key)
		}
		slices.Sort(keys)

		for _, key := range keys {
			lrItem := collection[key]

			isComplete := len(lrItem.right) == lrItem.position
			if lrItem.left != goalProduction && isComplete {
				// TODO Probably need to handle existing shift actions
				if existingAction := actionTable[i][lrItem.lookahead]; isReduce(existingAction) {
					conflictMessage := fmt.Sprintf("Reduce-Reduce Conflict at state %d on lookahead %s:\nExisting action: %s\nNew reduce action: %s",
						i, lrItem.lookahead, actionToString(existingAction), actionToString(&reduce{toCategory: lrItem.left, lenRightSide: len(lrItem.right), lrItem: lrItem}))
					panic(conflictMessage)
				}

				handler := findHandlerForNonTerminal(lrItem.left, productions)

				actionTable[i][lrItem.lookahead] = &reduce{
					toCategory:   lrItem.left,
					lenRightSide: len(lrItem.right),
					lrItem:       lrItem,
					handler:      handler,
					precedence:   lrItem.precedence,
				}
			} else if !isComplete && isFollowedByTerminal(lrItem, productions) {
				followingTerminal := lrItem.right[lrItem.position]

				goToState, ok := gotoTable[i][followingTerminal]
				if !ok {
					panic("no goto state defined")
				}

				handler := findHandlerForTerminal(followingTerminal, productions)
				shiftAction := &shift{
					toState:    goToState,
					lrItem:     lrItem,
					handler:    handler,
					precedence: lrItem.precedence,
				}

				existingAction := actionTable[i][followingTerminal]
				if isReduce(existingAction) {
					action, err := resolveShiftReduceConflict(existingAction.(*reduce), shiftAction)
					if err != nil {
						panic(err)
					}

					actionTable[i][followingTerminal] = action
				} else {
					actionTable[i][followingTerminal] = shiftAction
				}
			} else if lrItem.left == goalProduction && isComplete && lrItem.lookahead == token.EOF {
				if actionTable[i][token.EOF] != nil {
					panic("redefined action")
				}

				handler := findHandlerForNonTerminal(lrItem.left, productions)

				actionTable[i][token.EOF] = &accept{
					lenRightSide: len(lrItem.right),
					handler:      handler,
				}
			}
		}
	}

	return
}

func resolveShiftReduceConflict(reduceAction *reduce, shiftAction *shift) (parseAction, error) {
	if shiftAction.precedence == reduceAction.precedence {
		if isLrItemEqual(shiftAction.lrItem, reduceAction.lrItem) {
			return reduceAction, nil
		}

		return nil, fmt.Errorf(
			"cannot resolve shift-reduce conflict: Actions have the same precedence %d and are not resulting from the same grammar rule.\nShift action: %s,\nReduce action: %s",
			shiftAction.precedence,
			actionToString(shiftAction),
			actionToString(reduceAction),
		)
	} else if shiftAction.precedence > reduceAction.precedence {
		return shiftAction, nil
	} else {
		return reduceAction, nil
	}
}

func isLrItemEqual(a, b *LrItem) bool {
	if a.left != b.left {
		return false
	}

	if len(a.right) != len(b.right) {
		return false
	}

	for i := range a.right {
		if a.right[i] != b.right[i] {
			return false
		}
	}

	return true
}

func findHandlerForTerminal(terminal grammar.Category, productions []grammar.Production) func(string) ast.Node {
	for _, production := range productions {
		switch typedProduction := production.(type) {
		case *grammar.Terminal:
			if typedProduction.Name == terminal {
				return typedProduction.Handler
			}
		}
	}

	return nil
}

func findHandlerForNonTerminal(nonTerminal grammar.Category, productions []grammar.Production) func([]ast.Node) ast.Node {
	for _, production := range productions {
		switch typedProduction := production.(type) {
		case *grammar.NonTerminal:
			if typedProduction.Name == nonTerminal {
				return typedProduction.Handler
			}
		}
	}

	return nil
}

func actionToString(action parseAction) string {
	switch a := action.(type) {
	case *shift:
		return fmt.Sprintf("shift to state %d on %s", a.toState, a.lrItem)
	case *reduce:
		return fmt.Sprintf("reduce using rule %s", a.lrItem)
	case *accept:
		return "accept"
	default:
		return "unknown action"
	}
}

func isReduce(action parseAction) bool {
	if action == nil {
		return false
	}

	_, ok := action.(*reduce)
	return ok
}

// todo use map instead of slice
func isFollowedByTerminal(lrItem *LrItem, productions []grammar.Production) bool {
	if len(lrItem.right) == lrItem.position {
		return false
	}

	production := findProduction(productions, lrItem.right[lrItem.position])
	if production == nil {
		panic(fmt.Sprintf("encountered nil production when checking if %s isFollowedByTerminal", lrItem.String()))
	}

	if _, isTerminal := production.(*grammar.Terminal); isTerminal {
		return true
	}

	return false
}

func findProduction(productions []grammar.Production, identifier grammar.Category) grammar.Production {
	for _, production := range productions {
		switch typedProduction := production.(type) {
		case *grammar.Terminal:
			if typedProduction.Name == identifier {
				return production
			}
		case *grammar.NonTerminal:
			if typedProduction.Name == identifier {
				return production
			}
		}
	}

	return nil
}

func generateCanonicalCollection(input []grammar.Production) ([]map[string]*LrItem, map[int]map[grammar.Category]int) {
	flattenedProductions, rootSymbol := groupProductionsByCategory(input)

	ccZero := make(map[string]*LrItem)
	for _, production := range flattenedProductions[rootSymbol] {
		lrItem := productionToLrItem(production)
		lrItem.lookahead = token.EOF
		ccZero[lrItem.String()] = lrItem
	}

	firstSets := first(flattenedProductions)
	ccZero = closure(ccZero, flattenedProductions, firstSets)

	transitions := make(map[int]map[grammar.Category]int)
	canonicalCollections := []map[string]*LrItem{ccZero}

	changed := true
	for changed {
		changed = false

		for i, collection := range canonicalCollections {
			_, ok := transitions[i]
			if !ok {
				transitions[i] = make(map[grammar.Category]int)
			}
			for _, lrItem := range collection {
				if isAtEnd := len(lrItem.right) == lrItem.position; isAtEnd {
					continue
				}

				x := lrItem.right[lrItem.position]
				temp := goTo(collection, x, flattenedProductions, firstSets)

				if index := findCollectionInCanonicalCollections(canonicalCollections, temp); index == -1 {
					tempsIndex := len(canonicalCollections)
					canonicalCollections = append(canonicalCollections, temp)

					transitions[i][x] = tempsIndex
					changed = true
				} else {
					transitions[i][x] = index
				}
			}
		}
	}

	return canonicalCollections, transitions
}

func findCollectionInCanonicalCollections(cc []map[string]*LrItem, e map[string]*LrItem) int {
	for i, collection := range cc {
		if isCollectionEqual(collection, e) {
			return i
		}
	}

	return -1
}

func isCollectionEqual(a, b map[string]*LrItem) bool {
	if len(a) != len(b) {
		return false
	}

	for key := range a {
		_, ok := b[key]
		if !ok {
			return false
		}
	}

	return true
}

func productionToLrItem(prod grammar.Production) *LrItem {
	var left grammar.Category
	right := make([]grammar.Category, 0)

	switch typedProduction := prod.(type) {
	case *grammar.Terminal:
		left = typedProduction.Name
		right = append(right, typedProduction.Name)
		return &LrItem{left: left, right: right, position: 0}
	case *grammar.NonTerminal:
		left = typedProduction.Name

		switch rightSide := typedProduction.RightSide.(type) {
		case *grammar.Identifier:
			right = append(right, rightSide.Name)
		case *grammar.Sequence:
			left = typedProduction.Name

			for _, item := range rightSide.Items {
				right = append(right, item.Name)
			}
		case *grammar.Choice:
			panic(fmt.Sprintf("encountered a %T when converting %T to LrItem but %T is not allowed",
				rightSide, prod, rightSide))
		}

		return &LrItem{left: left, right: right, position: 0, precedence: typedProduction.Precedence}
	default:
		panic(fmt.Sprintf("unexpected production of type %T when converting to LrItem", prod))
	}
}

func closure(s map[string]*LrItem, productions map[grammar.Category][]grammar.Production, first map[grammar.Category][]grammar.Category) map[string]*LrItem {
	changed := true
	for changed {
		changed = false

		for _, lrItem := range s {
			var lookahead grammar.Category
			if len(lrItem.right)-1 > lrItem.position {
				lookahead = lrItem.right[lrItem.position+1]
			} else {
				lookahead = lrItem.lookahead
			}

			if len(lrItem.right) == lrItem.position {
				continue
			}

			currentItem := lrItem.right[lrItem.position]
			for _, production := range productions[currentItem] {
				if _, isTerminal := production.(*grammar.Terminal); isTerminal {
					continue
				}
				for _, b := range first[lookahead] {
					resultingLrItem := productionToLrItem(production)
					resultingLrItem.lookahead = b

					lrItemString := resultingLrItem.String()
					_, ok := s[lrItemString]
					if !ok {
						changed = true
						s[lrItemString] = resultingLrItem
					}
				}
			}
		}
	}

	return s
}

func goTo(s map[string]*LrItem, x grammar.Category, productions map[grammar.Category][]grammar.Production, firstSets map[grammar.Category][]grammar.Category) map[string]*LrItem {
	result := make(map[string]*LrItem)
	for _, lrItem := range s {
		if isAtEnd := len(lrItem.right) == lrItem.position; isAtEnd {
			continue
		}

		isBeforeX := lrItem.right[lrItem.position] == x
		if isBeforeX {
			newItem := &LrItem{
				left:       lrItem.left,
				right:      lrItem.right,
				position:   lrItem.position + 1,
				lookahead:  lrItem.lookahead,
				precedence: lrItem.precedence,
			}

			result[newItem.String()] = newItem
		}
	}

	return closure(result, productions, firstSets)
}

func first(productionsGroupedBySymbol map[grammar.Category][]grammar.Production) map[grammar.Category][]grammar.Category {
	firstSets := make(map[grammar.Category][]grammar.Category)

	for _, productionsForSymbol := range productionsGroupedBySymbol {
		if len(productionsForSymbol) == 0 {
			panic("encountered empty production when constructing initial firstSet")
		}

		production := productionsForSymbol[0]
		switch typedProduction := production.(type) {
		case *grammar.Terminal:
			firstSets[typedProduction.Name] = []grammar.Category{typedProduction.Name}
		case *grammar.NonTerminal:
			firstSets[typedProduction.Name] = make([]grammar.Category, 0)
		}
	}
	firstSets[token.EOF] = []grammar.Category{token.EOF}

	changed := true
	for changed {
		changed = false

		for symbol, productionsForSymbol := range productionsGroupedBySymbol {
			for _, production := range productionsForSymbol {
				switch typedProduction := production.(type) {
				case *grammar.Terminal:
					continue
				case *grammar.NonTerminal:
					firstSetOfSymbol := firstSets[symbol]
					switch rightSide := typedProduction.RightSide.(type) {
					case *grammar.Choice:
						panic("encountered grammar.Choice when computing first sets")
					case *grammar.Identifier:
						firstSetOfRight, ok := firstSets[rightSide.Name]
						if !ok {
							continue
						}

						newFirstSet, newElementsAdded := appendFirstSet(firstSetOfSymbol, firstSetOfRight)
						if newElementsAdded {
							firstSets[symbol] = newFirstSet
							changed = true
						}
					case *grammar.Sequence:
						if len(rightSide.Items) == 0 {
							panic("encountered right side with no items when computing first sets")
						}

						firstItemInSequence := rightSide.Items[0].Name
						firstSetOfRight, ok := firstSets[firstItemInSequence]
						if !ok {
							continue
						}

						newFirstSet, newElementsAdded := appendFirstSet(firstSetOfSymbol, firstSetOfRight)
						if newElementsAdded {
							firstSets[symbol] = newFirstSet
							changed = true
						}
					}
				}
			}
		}
	}

	return firstSets
}

func appendFirstSet(existing []grammar.Category, newElements []grammar.Category) ([]grammar.Category, bool) {
	changed := false
	for _, newElement := range newElements {

		index := slices.Index(existing, newElement)
		if index == -1 {
			existing = append(existing, newElement)
			changed = true
		}
	}
	return existing, changed
}

func groupProductionsByCategory(productions []grammar.Production) (map[grammar.Category][]grammar.Production, grammar.Category) {
	productionsGroupedByName := make(map[grammar.Category][]grammar.Production)
	var initialCategory grammar.Category

	for i, production := range productions {
		var name grammar.Category

		switch typedProduction := production.(type) {
		case *grammar.Terminal:
			name = typedProduction.Name
			productionsGroupedByName[typedProduction.Name] = []grammar.Production{typedProduction}
		case *grammar.NonTerminal:
			name = typedProduction.Name
			productionsGroupedByName[typedProduction.Name] = flattenNonTerminal(typedProduction)
		default:
			panic(fmt.Sprintf("unexpected production of type %T when flattening productions", production))
		}

		if i == 0 {
			initialCategory = name
		}
	}

	return productionsGroupedByName, initialCategory
}

func flattenNonTerminal(nt *grammar.NonTerminal) []grammar.Production {
	switch rightSide := nt.RightSide.(type) {
	case *grammar.Sequence:
		return []grammar.Production{nt}
	case *grammar.Identifier:
		return []grammar.Production{nt}
	case *grammar.Choice:
		flattenedChoice := make([]grammar.Production, len(rightSide.Items))
		for i, production := range rightSide.Items {
			flattenedChoice[i] = grammar.NewNonTerminal(nt.Name, production, production.Precedence())
		}
		return flattenedChoice
	}
	panic(fmt.Sprintf("unexpected Right Side of type %T when flattening NonTerminal '%s'", nt, nt.Name))
}
