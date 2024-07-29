package scannergenerator

import "compiler/token"

type ScannerGenerator struct {
}

func NewScannerGenerator() *ScannerGenerator {
	return &ScannerGenerator{}
}

func (s *ScannerGenerator) GenerateScanner(tokenClassifications []token.TokenClassification) *Dfa {
	var nfa *Nfa
	precedences := make(map[token.TokenType]int)
	for _, tokenClassification := range tokenClassifications {
		nfaForClassification := NewRegexpToNfaConverter(tokenClassification.Regexp).Convert()

		for _, state := range nfaForClassification.AcceptingStates {
			nfaForClassification.TypeTable[state] = tokenClassification.TokenType
		}

		if nfa == nil {
			nfa = nfaForClassification
		} else {
			nfa = nfa.UnionDistinct(nfaForClassification)
		}

		precedences[tokenClassification.TokenType] = tokenClassification.Precedence
	}

	dfa := NewNfaToDfaConverter(nfa, precedences).Convert()

	return dfa
}
