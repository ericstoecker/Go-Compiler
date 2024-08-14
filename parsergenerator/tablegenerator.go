package parsergenerator

import (
	"compiler/grammar"
	"fmt"
	"log"
)

type TableGenerator struct{}

func NewTableGenerator() *TableGenerator {
	return &TableGenerator{}
}

func (tg *TableGenerator) generateCanonicalCollection(productions []grammar.Production) []map[string]*LrItem {
	flattenedProductions, rootSymbol := flattenProductions(productions)

	for key, productions := range flattenedProductions {
		log.Printf("productions for key %s", key)
		for _, production := range productions {
			lrItem := productionToLrItem(production)
			log.Print(lrItem.String())
		}
	}

	ccZero := make(map[string]*LrItem)
	for _, production := range flattenedProductions[rootSymbol] {
		lrItem := productionToLrItem(production)
		ccZero[lrItem.String()] = lrItem
	}

	result := closure(ccZero, flattenedProductions)

	return nil
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

func closure(s map[string]*LrItem, productions map[grammar.Category]grammar.Production) map[string]*LrItem {

	return nil
}

func goTo(map[string]grammar.Category) map[string]*LrItem {
	return nil
}

func first(symbol string) []string {
	return nil
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
