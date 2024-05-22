package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GlobalScope"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	symbols        map[string]*Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]*Symbol), numDefinitions: 0}
}

func (s *SymbolTable) Define(name string) *Symbol {
	symbol := &Symbol{Name: name, Scope: GlobalScope, Index: s.numDefinitions}

	s.symbols[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) RetrieveSymbol(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[name]
	return symbol, ok
}
