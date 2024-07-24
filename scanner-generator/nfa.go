package scannergenerator

type Nfa struct {
	Transitions  map[string]map[int][]int
	InitialState int
	FinalState   int
}

func (n *Nfa) concatenation(other *Nfa) *Nfa {
	for symbol, transitions := range other.Transitions {
		if n.Transitions[symbol] == nil {
			n.Transitions[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+n.FinalState+1)
			}
			n.Transitions[symbol][stateFrom+n.FinalState+1] = statesTo
		}
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}
	n.Transitions[EPSILON][n.FinalState] = []int{other.InitialState + n.FinalState + 1}

	return &Nfa{Transitions: n.Transitions, InitialState: n.InitialState, FinalState: n.FinalState + other.FinalState + 1}
}

func (n *Nfa) union(other *Nfa) *Nfa {
	for symbol, transitions := range other.Transitions {
		if n.Transitions[symbol] == nil {
			n.Transitions[symbol] = make(map[int][]int)
		}
		for stateFrom, stateTo := range transitions {
			statesTo := []int{}
			for _, state := range stateTo {
				statesTo = append(statesTo, state+n.FinalState+1)
			}
			n.Transitions[symbol][stateFrom+n.FinalState+1] = statesTo
		}
	}

	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	numberOfStatesInUnion := n.FinalState + other.FinalState + 2
	n.Transitions[EPSILON][numberOfStatesInUnion] = []int{n.FinalState + 1, n.InitialState}
	n.Transitions[EPSILON][n.FinalState] = []int{numberOfStatesInUnion + 1}
	n.Transitions[EPSILON][n.FinalState+other.FinalState+1] = []int{numberOfStatesInUnion + 1}

	return &Nfa{Transitions: n.Transitions, InitialState: numberOfStatesInUnion, FinalState: numberOfStatesInUnion + 1}
}

func (n *Nfa) kleene() *Nfa {
	if n.Transitions[EPSILON] == nil {
		n.Transitions[EPSILON] = make(map[int][]int)
	}

	initialState := n.FinalState + 1
	finalState := n.FinalState + 2

	n.Transitions[EPSILON][initialState] = []int{n.InitialState, finalState}
	n.Transitions[EPSILON][n.FinalState] = []int{n.InitialState, finalState}

	return &Nfa{Transitions: n.Transitions, InitialState: initialState, FinalState: finalState}
}
