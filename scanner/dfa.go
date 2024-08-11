package scanner

import "compiler/token"

type Dfa struct {
	Transitions     map[string]map[int]int
	InitialState    int
	AcceptingStates []int

	TypeTable map[int]token.TokenType
}
