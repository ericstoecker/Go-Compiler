package parsergenerator

import (
	"bytes"
	"compiler/grammar"
	"compiler/token"
)

type LrItem struct {
	left      grammar.Category
	right     []grammar.Category
	lookahead token.TokenType

	position int
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

	out.WriteString(", ")
	out.WriteString(string(li.lookahead))
	out.WriteString("]")

	return out.String()
}
