package scannergenerator

type RegexpToNfaConverter struct {
}

// maybe also return error
func (c *RegexpToNfaConverter) Convert(regexp string) (result map[string]map[int]int) {
	result = make(map[string]map[int]int)
	for _, char := range regexp {
		result[string(char)] = map[int]int{0: 1}
	}
	return
}
