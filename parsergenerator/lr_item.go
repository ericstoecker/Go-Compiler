package parsergenerator

import (
	"bytes"
	"compiler/ast"
	"compiler/grammar"
)

type LrItem struct {
	left      grammar.Category
	right     []grammar.Category
	lookahead grammar.Category

	position int

	terminalHandler    func(string) ast.Node
	nonTerminalHandler func([]ast.Node) ast.Node
}

func (li *LrItem) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	out.WriteString(string(li.left))
	out.WriteString(" =>")

	for i, item := range li.right {
		out.WriteString(" ")

		if i == li.position {
			out.WriteString("● ")
		}
		out.WriteString(string(item))
	}

	if li.position == len(li.right) {
		out.WriteString(" ●")
	}

	out.WriteString(", ")
	out.WriteString(string(li.lookahead))
	out.WriteString("]")

	return out.String()
}
