package parsergenerator

import (
	"compiler/grammar"
	"compiler/token"
	"fmt"
	"log"
	"slices"
)

type TableGenerator struct{}

func NewTableGenerator() *TableGenerator {
	return &TableGenerator{}
}

func (tg *TableGenerator) generateCanonicalCollection(productions []grammar.Production) []map[string]*LrItem {
	flattenedProductions, rootSymbol := flattenProductions(productions)

	for _, productions := range flattenedProductions {
		for _, production := range productions {
			lrItem := productionToLrItem(production)
			log.Print(lrItem.String())
		}
	}

	ccZero := make(map[string]*LrItem)
	for _, production := range flattenedProductions[rootSymbol] {
		lrItem := productionToLrItem(production)
		lrItem.lookahead = token.EOF
		ccZero[lrItem.String()] = lrItem
	}

	firstSets := first(flattenedProductions)
	ccZero = closure(ccZero, flattenedProductions, firstSets)

	return []map[string]*LrItem{ccZero}
}

func productionToLrItem(prod grammar.Production) *LrItem {

	var left grammar.Category
	right := make([]grammar.Category, 0)

	switch typedProduction := prod.(type) {
	case *grammar.Terminal:
		left = typedProduction.Name
		right = append(right, typedProduction.Name)
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
	}

	return &LrItem{left: left, right: right, position: 0}
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

func goTo(map[string]grammar.Category) map[string]*LrItem {
	return nil
}

func first(productionsGroupedBySymbol map[grammar.Category][]grammar.Production) map[grammar.Category][]grammar.Category {
	firstSets := make(map[grammar.Category][]grammar.Category)

	// todo no need to loop over inner elements
	// simply use first item
	// at least one has to exist due to earlier construction
	// otherwise panic
	for _, productionsForSymbol := range productionsGroupedBySymbol {
		for _, production := range productionsForSymbol {
			switch typedProduction := production.(type) {
			case *grammar.Terminal:
				firstSets[typedProduction.Name] = []grammar.Category{typedProduction.Name}
			case *grammar.NonTerminal:
				firstSets[typedProduction.Name] = make([]grammar.Category, 0)
			}
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
					switch rightSide := typedProduction.RightSide.(type) {
					case *grammar.Choice:
						panic("encountered grammar.Choice when computing first sets")
					case *grammar.Identifier:
						firstSetOfRight, ok := firstSets[rightSide.Name]
						if !ok {
							continue
						}

						newFirstSet, newElementsAdded := appendFirstSet(firstSets[symbol], firstSetOfRight)
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

						newFirstSet, newElementsAdded := appendFirstSet(firstSets[symbol], firstSetOfRight)
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

func flattenProductions(productions []grammar.Production) (map[grammar.Category][]grammar.Production, grammar.Category) {
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
			flattenedChoice[i] = grammar.NewNonTerminal(nt.Name, production)
		}
		return flattenedChoice
	}
	panic(fmt.Sprintf("unexpected Right Side of type %T when flattening NonTerminal", nt))
}
