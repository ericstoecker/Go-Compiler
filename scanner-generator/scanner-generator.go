package scannergenerator

import (
	"compiler/token"
	"fmt"
)

type ScannerGenerator struct {
}

func NewScannerGenerator() *ScannerGenerator {
	return &ScannerGenerator{}
}

func (s *ScannerGenerator) GenerateScanner(tokenClassifications []token.TokenClassification) *Dfa {
	precedences := make(map[token.TokenType]int)
	nfas := make([]*Nfa, len(tokenClassifications))
	for i, tokenClassification := range tokenClassifications {
		nfaForClassification, err := NewRegexpToNfaConverter(tokenClassification.Regexp).Convert()
		if err != nil {
			panic(fmt.Sprintf("error when converting regexp '%s' to nfa: %v", tokenClassification.Regexp, err))
		}

		for _, state := range nfaForClassification.AcceptingStates {
			nfaForClassification.TypeTable[state] = tokenClassification.TokenType
		}

		nfas[i] = nfaForClassification

		precedences[tokenClassification.TokenType] = tokenClassification.Precedence
	}

	nfa := nfas[0].UnionDistinct(nfas[1:]...)
	dfa := NewNfaToDfaConverter(nfa, precedences).Convert()

	dfaMinimizer := &DfaMinimizer{}

	return dfaMinimizer.Minimize(dfa)
}
