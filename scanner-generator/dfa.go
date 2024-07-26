package scannergenerator

type Dfa struct {
	Transitions     map[string]map[int]int
	InitialState    int
	AcceptingStates []int
}
