package scanner

import (
	"compiler/token"
	"slices"
)

type NfaToDfaConverter struct {
	nfa *Nfa

	TypePrecedences map[token.TokenType]int
}

func NewNfaToDfaConverter(nfa *Nfa, typePrecedences map[token.TokenType]int) *NfaToDfaConverter {
	return &NfaToDfaConverter{nfa: nfa, TypePrecedences: typePrecedences}
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
	typeTable := make(map[int]token.TokenType)
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

		highestRankingTokenType := -1
		for _, acceptingNfaState := range c.nfa.AcceptingStates {
			if slices.Contains(currentItem, acceptingNfaState) {
				acceptingStates = append(acceptingStates, currentItemsIndex)
				statesType, ok := c.nfa.TypeTable[acceptingNfaState]

				if ok && c.TypePrecedences[statesType] > highestRankingTokenType {
					highestRankingTokenType = c.TypePrecedences[statesType]
					typeTable[currentItemsIndex] = statesType
				}
			}
		}
	}

	// Build and return the DFA
	return &Dfa{
		Transitions:     transitions,
		AcceptingStates: acceptingStates,
		InitialState:    0,
		TypeTable:       typeTable,
	}
}

func (c *NfaToDfaConverter) followEpsilon(states []int) []int {
	visited := make(map[int]bool)
	result := make([]int, 0)
	for _, state := range states {
		closure := c.followEpsilonFromState(state, visited)
		result = append(result, closure...)
	}

	return filterDuplicates(result)
}

func (c *NfaToDfaConverter) followEpsilonFromState(state int, visited map[int]bool) []int {
	if visited[state] {
		return []int{}
	}
	visited[state] = true

	result := []int{state}

	epsilonTransitionsForState := c.nfa.Transitions[EPSILON][state]
	for _, neighbor := range epsilonTransitionsForState {
		result = append(result, c.followEpsilonFromState(neighbor, visited)...)
	}
	return result
}

func filterDuplicates(a []int) []int {
	slices.Sort(a)
	return slices.Compact(a)
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
