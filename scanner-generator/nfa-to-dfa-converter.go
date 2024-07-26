package scannergenerator

import (
	"slices"
)

type NfaToDfaConverter struct {
	nfa *Nfa
}

func NewNfaToDfaConverter(nfa *Nfa) *NfaToDfaConverter {
	return &NfaToDfaConverter{nfa: nfa}
}

func (c *NfaToDfaConverter) Convert() *Dfa {
	characters := make([]string, 0)
	transitions := make(map[string]map[int]int)
	for char := range c.nfa.Transitions {
		if char == EPSILON {
			continue
		}

		characters = append(characters, char)
		transitions[char] = make(map[int]int)
	}
	slices.Sort(characters) // Sort characters to ensure deterministic order

	currentItem := c.followEpsilon([]int{c.nfa.InitialState})
	dfaStates := [][]int{currentItem}
	workList := [][]int{currentItem}

	acceptingStates := make([]int, 0)
	for len(workList) != 0 {
		currentItem, workList = workList[0], workList[1:]
		currentItemsIndex := findIndex(currentItem, dfaStates)

		for _, char := range characters {
			temp := c.followEpsilon(c.delta(currentItem, char))
			tempsIndex := findIndex(temp, dfaStates)

			if len(temp) == 0 {
				continue
			}

			if tempsIndex == -1 {
				dfaStates = append(dfaStates, temp)
				workList = append(workList, temp)
				transitions[char][currentItemsIndex] = len(dfaStates) - 1
			} else {
				transitions[char][currentItemsIndex] = tempsIndex
			}

		}

		if slices.Contains(currentItem, c.nfa.FinalState) {
			acceptingStates = append(acceptingStates, currentItemsIndex)
		}
	}

	return &Dfa{
		Transitions:     transitions,
		InitialState:    0,
		AcceptingStates: acceptingStates,
	}
}

func (c *NfaToDfaConverter) followEpsilon(states []int) []int {
	epsilonTransitions := c.nfa.Transitions[EPSILON]
	if len(epsilonTransitions) == 0 {
		return states
	}

	result := make([]int, 0)
	result = append(result, states...)
	for _, state := range states {
		epsilonTransitionsForCurrentState := c.followEpsilonFromState(state)
		result = append(result, epsilonTransitionsForCurrentState...)
	}

	return filterDuplicates(result)
}

func filterDuplicates(a []int) []int {
	slices.Sort(a)
	return slices.Compact(a)
}

func (c *NfaToDfaConverter) followEpsilonFromState(state int) []int {
	epsilonTransitions := c.nfa.Transitions[EPSILON]
	if len(epsilonTransitions) == 0 {
		return []int{state}
	}

	epsilonTransitionsForState := epsilonTransitions[state]
	if len(epsilonTransitionsForState) == 0 {
		return []int{state}
	}

	directNeighborsThroughEpsilon := c.delta([]int{state}, EPSILON)
	result := make([]int, 0)
	for _, neighboringState := range directNeighborsThroughEpsilon {
		indirectNeighbors := c.followEpsilonFromState(neighboringState)
		result = append(result, indirectNeighbors...)
		result = append(result, neighboringState)
	}
	return result
}

func (c *NfaToDfaConverter) delta(states []int, char string) []int {
	transitionsForCharacter := c.nfa.Transitions[char]
	if len(transitionsForCharacter) == 0 {
		return make([]int, 0)
	}

	result := make([]int, 0)
	for _, state := range states {
		transitionsForCurrentState, ok := transitionsForCharacter[state]
		if ok {
			result = append(result, transitionsForCurrentState...)
		}
	}
	return result
}

func findIndex(searchedStates []int, dfaStates [][]int) int {
	slices.Sort(searchedStates)
	for index, states := range dfaStates {
		slices.Sort(states)
		if slices.Equal(states, searchedStates) {
			return index
		}
	}
	return -1
}