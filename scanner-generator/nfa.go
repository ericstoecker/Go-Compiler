package scannergenerator

import (
	"compiler/token"
)

type Nfa struct {
	Transitions     map[string]map[int][]int
	InitialState    int
	AcceptingStates []int
	NumberOfStates  int

	TypeTable map[int]token.TokenType
}

func NfaFromSingleSymbol(symbol string) *Nfa {
	return &Nfa{
		Transitions: map[string]map[int][]int{
			symbol: {0: []int{1}},
		},
		InitialState:    0,
		AcceptingStates: []int{1},
		NumberOfStates:  2,
	}
}

func (n *Nfa) Concatenation(other *Nfa) *Nfa {
	if n.NumberOfStates == 0 {
		panic("n.NumberOfStates is 0")
	}
	if other.NumberOfStates == 0 {
		panic("other.NumberOfStates is 0")
	}

	for symbol, transitions := range other.Transitions {
		if n.Transitions[symbol] == nil {
			n.Transitions[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+n.NumberOfStates)
			}
			n.Transitions[symbol][stateFrom+n.NumberOfStates] = statesTo
		}
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{other.InitialState + n.NumberOfStates}
	}

	acceptingStates := make([]int, 0)
	for _, state := range other.AcceptingStates {
		acceptingStates = append(acceptingStates, state+n.NumberOfStates)
	}
	return &Nfa{Transitions: n.Transitions, InitialState: n.InitialState, AcceptingStates: acceptingStates, NumberOfStates: n.NumberOfStates + other.NumberOfStates}
}

func (n *Nfa) Union(others ...*Nfa) *Nfa {
	if n.NumberOfStates == 0 {
		panic("n.NumberOfStates is 0")
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	totalNumberOfStates := n.NumberOfStates
	for _, other := range others {
		totalNumberOfStates += other.NumberOfStates
	}

	currentNumberOfStates := n.NumberOfStates
	for _, other := range others {
		if other.NumberOfStates == 0 {
			panic("other.NumberOfStates is 0")
		}

		for symbol, transitions := range other.Transitions {
			if n.Transitions[symbol] == nil {
				n.Transitions[symbol] = make(map[int][]int)
			}
			for stateFrom, stateTo := range transitions {
				statesTo := []int{}
				for _, state := range stateTo {
					statesTo = append(statesTo, state+currentNumberOfStates)
				}
				n.Transitions[symbol][stateFrom+currentNumberOfStates] = statesTo
			}
		}

		n.Transitions[EPSILON][totalNumberOfStates] = append(n.Transitions[EPSILON][totalNumberOfStates], other.InitialState+currentNumberOfStates)
		for _, state := range other.AcceptingStates {
			n.Transitions[EPSILON][state+currentNumberOfStates] = []int{totalNumberOfStates + 1}
		}

		currentNumberOfStates += other.NumberOfStates
	}

	n.Transitions[EPSILON][totalNumberOfStates] = append(n.Transitions[EPSILON][totalNumberOfStates], n.InitialState)
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{totalNumberOfStates + 1}
	}

	return &Nfa{Transitions: n.Transitions, InitialState: currentNumberOfStates, AcceptingStates: []int{currentNumberOfStates + 1}, NumberOfStates: currentNumberOfStates + 2}
}

func (n *Nfa) UnionDistinct(others ...*Nfa) *Nfa {
	if n.NumberOfStates == 0 {
		panic("n.NumberOfStates is 0")
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	totalNumberOfStates := n.NumberOfStates
	for _, other := range others {
		totalNumberOfStates += other.NumberOfStates
	}

	currentNumberOfStates := n.NumberOfStates
	acceptingStates := make([]int, 0)
	typeTable := make(map[int]token.TokenType)
	for _, other := range others {
		if other.NumberOfStates == 0 {
			panic("other.NumberOfStates is 0")
		}

		for symbol, transitions := range other.Transitions {
			if n.Transitions[symbol] == nil {
				n.Transitions[symbol] = make(map[int][]int)
			}
			for stateFrom, stateTo := range transitions {
				statesTo := []int{}
				for _, state := range stateTo {
					statesTo = append(statesTo, state+currentNumberOfStates)
				}
				n.Transitions[symbol][stateFrom+currentNumberOfStates] = statesTo
			}
		}

		n.Transitions[EPSILON][totalNumberOfStates] = append(n.Transitions[EPSILON][totalNumberOfStates], other.InitialState+currentNumberOfStates)

		acceptingStates = append(acceptingStates, n.AcceptingStates...)
		for _, state := range other.AcceptingStates {
			acceptingStates = append(acceptingStates, state+currentNumberOfStates)
		}

		for state, tokenType := range other.TypeTable {
			typeTable[state+currentNumberOfStates] = tokenType
		}

		currentNumberOfStates += other.NumberOfStates
	}

	n.Transitions[EPSILON][totalNumberOfStates] = append(n.Transitions[EPSILON][totalNumberOfStates], n.InitialState)

	for state, tokenType := range n.TypeTable {
		typeTable[state] = tokenType
	}

	return &Nfa{Transitions: n.Transitions, InitialState: totalNumberOfStates, AcceptingStates: acceptingStates, TypeTable: typeTable, NumberOfStates: totalNumberOfStates + 1}
}

func (n *Nfa) Kleene() *Nfa {
	if n.NumberOfStates == 0 {
		panic("n.NumberOfStates is 0")
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	initialState := n.NumberOfStates
	finalState := n.NumberOfStates + 1

	n.Transitions[EPSILON][initialState] = []int{n.InitialState, finalState}
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{n.InitialState, finalState}
	}

	return &Nfa{Transitions: n.Transitions, InitialState: initialState, AcceptingStates: []int{finalState}, NumberOfStates: n.NumberOfStates + 2}
}
