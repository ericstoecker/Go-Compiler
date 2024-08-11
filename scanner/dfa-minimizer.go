package scanner

import (
	"compiler/token"
	"slices"
)

type DfaMinimizer struct {
}

// optimize using maps as set to avoid costly 'state contained in set x' computations
func (m *DfaMinimizer) Minimize(dfa *Dfa) *Dfa {
	characters, dfaStates := extractStatesAndCharacters(dfa)

	partition := partitionIntoAcceptingAndNonAccepting(dfa, dfaStates)
	worklist := slices.Clone(partition)

	for len(worklist) > 0 {
		currentItem := worklist[0]
		worklist = worklist[1:]

		for _, char := range characters {
			image := computeImage(dfa.Transitions, char, dfaStates, currentItem)
			if len(image) == 0 {
				continue
			}

			newPartition := make([][]int, 0)
			newPartition = append(newPartition, partition...)
			for _, elementInPartition := range partition {
				intersectionWithImage := intersection(elementInPartition, image)
				if len(intersectionWithImage) == 0 {
					continue
				}

				remainder := difference(elementInPartition, intersectionWithImage)
				if len(remainder) == 0 {
					continue
				}

				if index := findIndex(elementInPartition, newPartition); index != -1 {
					newPartition = slices.Delete(newPartition, index, index+1)
					newPartition = append(newPartition, intersectionWithImage, remainder)
				}

				if index := findIndex(currentItem, worklist); index != -1 {
					worklist = slices.Delete(worklist, index, index+1)
					worklist = append(worklist, intersectionWithImage, remainder)
				} else if len(intersectionWithImage) <= len(remainder) {
					worklist = append(worklist, intersectionWithImage)
				} else {
					worklist = append(worklist, remainder)
				}
			}

			partition = newPartition
		}
	}

	return constructDfaFromPartition(partition, dfa, characters)
}

func partitionIntoAcceptingAndNonAccepting(dfa *Dfa, dfaStates []int) [][]int {
	acceptingStateSets := make(map[token.TokenType][]int)
	nonacceptingStates := make([]int, 0)
	for _, state := range dfaStates {
		statesType, ok := dfa.TypeTable[state]
		if ok {
			acceptingStateSets[statesType] = append(acceptingStateSets[statesType], state)
		} else {
			nonacceptingStates = append(nonacceptingStates, state)
		}
	}

	result := make([][]int, len(acceptingStateSets)+1)
	i := 0
	for _, acceptingStates := range acceptingStateSets {
		result[i] = acceptingStates
		i++
	}
	result[len(acceptingStateSets)] = nonacceptingStates

	return result
}

func constructDfaFromPartition(partition [][]int, dfa *Dfa, characters []string) *Dfa {
	typeTable := make(map[int]token.TokenType)
	transitions := make(map[string]map[int]int)
	acceptingStates := make([]int, 0)
	initialState := 0
	for index, stateSet := range partition {
		state := stateSet[0]

		isAcceptingState := slices.Contains(dfa.AcceptingStates, state)
		if isAcceptingState {
			acceptingStates = append(acceptingStates, index)
			typeTable[index] = dfa.TypeTable[state]
		}

		isInitialState := slices.Contains(stateSet, dfa.InitialState)
		if isInitialState {
			initialState = index
		}

		for _, char := range characters {
			resultOfTransition, ok := dfa.Transitions[char][state]
			if ok {
				indexOfSetContaingResult := findIndexInPartition(partition, resultOfTransition)
				if transitions[char] == nil {
					transitions[char] = make(map[int]int)
				}
				transitions[char][index] = indexOfSetContaingResult
			}
		}
	}

	return &Dfa{Transitions: transitions, AcceptingStates: acceptingStates, InitialState: initialState, TypeTable: typeTable}
}

func findIndexInPartition(partitions [][]int, state int) int {
	for index, states := range partitions {
		slices.Sort(states)
		_, isContained := slices.BinarySearch(states, state)
		if isContained {
			return index
		}
	}

	return -1
}

func computeImage(transitions map[string]map[int]int, char string, dfaStates, states []int) []int {
	statesSet := make(map[int]struct{})
	for _, state := range states {
		statesSet[state] = struct{}{}
	}

	image := make([]int, 0)
	for _, state := range dfaStates {
		resultOfTransition, ok := transitions[char][state]

		if ok {
			_, ok = statesSet[resultOfTransition]
			if ok {
				image = append(image, state)
			}
		}
	}

	return image
}

func intersection(a, b []int) []int {
	slices.Sort(b)

	states := make([]int, 0)
	for _, state := range a {
		_, containedInB := slices.BinarySearch(b, state)
		if containedInB {
			states = append(states, state)
		}
	}

	return states
}

func difference(a, b []int) []int {
	slices.Sort(b)

	states := make([]int, 0)
	for _, state := range a {
		_, containedInB := slices.BinarySearch(b, state)
		if !containedInB {
			states = append(states, state)
		}
	}

	return states
}

func extractStatesAndCharacters(dfa *Dfa) ([]string, []int) {
	characters := make([]string, 0)
	states := make([]int, 0)
	for char, charMappings := range dfa.Transitions {
		characters = append(characters, char)
		for stateFrom, stateTo := range charMappings {
			states = append(states, stateFrom, stateTo)
		}
	}

	slices.Sort(characters)
	slices.Sort(states)
	return characters, slices.Compact(states)
}
