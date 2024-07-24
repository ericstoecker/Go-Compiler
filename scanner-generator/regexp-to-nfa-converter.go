package scannergenerator

type RegexpToNfaConverter struct {
}

// maybe also return error
func (c *RegexpToNfaConverter) Convert(regexp string) map[string]map[int]int {
	return map[string]map[int]int{"a": {0: 1}}
}
