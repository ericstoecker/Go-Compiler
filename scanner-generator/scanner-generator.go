package scannergenerator

import "compiler/token"

type TokenClassification struct {
	regexp     string
	tokenType  token.TokenType
	precedence int
}

type ScannerGenerator struct {
}

func NewScannerGenerator() *ScannerGenerator {
	return &ScannerGenerator{}
}

func (s *ScannerGenerator) GenerateScanner(tokenClassifications []TokenClassification) *Dfa {
	var nfa *Nfa
	precedences := make(map[token.TokenType]int)
	for _, tokenClassification := range tokenClassifications {
		nfaForClassification := NewRegexpToNfaConverter(tokenClassification.regexp).Convert()

		for _, state := range nfaForClassification.AcceptingStates {
			nfaForClassification.TypeTable[state] = tokenClassification.tokenType
		}

		if nfa == nil {
			nfa = nfaForClassification
		} else {
			nfa = nfa.UnionDistinct(nfaForClassification)
		}

		precedences[tokenClassification.tokenType] = tokenClassification.precedence
	}

	dfa := NewNfaToDfaConverter(nfa, precedences).Convert()

	return dfa
}
