package scannergenerator

import (
	"compiler/token"
	"slices"
)

type Nfa struct {
	Transitions     map[string]map[int][]int
	InitialState    int
	AcceptingStates []int

	TypeTable map[int]token.TokenType
}

func (n *Nfa) concatenation(other *Nfa) *Nfa {
	highestState := slices.Max(n.AcceptingStates)
	for symbol, transitions := range other.Transitions {
		if n.Transitions[symbol] == nil {
			n.Transitions[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+highestState+1)
			}
			n.Transitions[symbol][stateFrom+highestState+1] = statesTo
		}
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{other.InitialState + highestState + 1}
	}

	acceptingStates := make([]int, 0)
	for _, state := range other.AcceptingStates {
		acceptingStates = append(acceptingStates, state+highestState+1)
	}
	return &Nfa{Transitions: n.Transitions, InitialState: n.InitialState, AcceptingStates: acceptingStates}
}

func (n *Nfa) union(other *Nfa) *Nfa {
	highestState := slices.Max(n.AcceptingStates)
	for symbol, transitions := range other.Transitions {
		if n.Transitions[symbol] == nil {
			n.Transitions[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+highestState+1)
			}
			n.Transitions[symbol][stateFrom+highestState+1] = statesTo
		}
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	highestStateInOther := slices.Max(other.AcceptingStates)
	numberOfStatesInUnion := highestState + highestStateInOther + 2
	n.Transitions[EPSILON][numberOfStatesInUnion] = []int{highestState + 1, n.InitialState}
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{numberOfStatesInUnion + 1}
	}
	for _, state := range other.AcceptingStates {
		n.Transitions[EPSILON][state+highestState+1] = []int{numberOfStatesInUnion + 1}
	}

	return &Nfa{Transitions: n.Transitions, InitialState: numberOfStatesInUnion, AcceptingStates: []int{numberOfStatesInUnion + 1}}
}

func (n *Nfa) unionDistinct(other *Nfa) *Nfa {
	return nil
}

func (n *Nfa) kleene() *Nfa {
	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	highestState := slices.Max(n.AcceptingStates)
	initialState := highestState + 1
	finalState := highestState + 2

	n.Transitions[EPSILON][initialState] = []int{n.InitialState, finalState}
	for _, state := range n.AcceptingStates {
		n.Transitions[EPSILON][state] = []int{n.InitialState, finalState}
	}

	return &Nfa{Transitions: n.Transitions, InitialState: initialState, AcceptingStates: []int{finalState}}
}
