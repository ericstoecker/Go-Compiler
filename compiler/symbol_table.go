package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GlobalScope"
	LocalScope    SymbolScope = "LocalScope"
	BuiltinScope  SymbolScope = "SymbolScope"
	FreeScope     SymbolScope = "FreeScope"
	FunctionScope SymbolScope = "FunctionScope"
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

	FreeSymbols []*Symbol
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]*Symbol)
	free := []*Symbol{}
	return &SymbolTable{symbols: s, FreeSymbols: free}
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

func (s *SymbolTable) DefineBuiltin(index int, name string) *Symbol {
	symbol := &Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.symbols[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) *Symbol {
	symbol := &Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.symbols[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(original *Symbol) *Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := &Symbol{Name: original.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	s.symbols[original.Name] = symbol
	return symbol
}

func (s *SymbolTable) RetrieveSymbol(name string) (*Symbol, bool) {
	symbol, ok := s.symbols[name]
	if !ok && s.outer != nil {
		symbol, ok = s.outer.RetrieveSymbol(name)
		if !ok {
			return symbol, ok
		}

		if symbol.Scope == GlobalScope || symbol.Scope == BuiltinScope {
			return symbol, ok
		}

		free := s.defineFree(symbol)
		return free, true
	}
	return symbol, ok
}
