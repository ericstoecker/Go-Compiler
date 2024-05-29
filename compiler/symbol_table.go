package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GlobalScope"
	LocalScope  SymbolScope = "LocalScope"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	outer *SymbolTable

	symbols        map[string]*Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{symbols: make(map[string]*Symbol), numDefinitions: 0}
}

func FromSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{outer: outer, symbols: make(map[string]*Symbol), numDefinitions: 0}
}

func (s *SymbolTable) UnwrapSymbolTable() *SymbolTable {
	return s.outer
}

func (s *SymbolTable) Define(name string) *Symbol {
	symbol := &Symbol{Name: name, Index: s.numDefinitions}

	if s.outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.symbols[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) RetrieveSymbol(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[name]
	if !ok {
		return s.outer.RetrieveSymbol(name)
	}
	return symbol, ok
}
