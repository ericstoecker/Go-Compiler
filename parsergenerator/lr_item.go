package parsergenerator

import (
	"bytes"
	"compiler/grammar"
	"strconv"
)

type LrItem struct {
	left      grammar.Category
	right     []grammar.Category
	lookahead grammar.Category

	position   int
	precedence int
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
	out.WriteString("lookahead: ")
	out.WriteString(string(li.lookahead))
	out.WriteString(", ")
	out.WriteString("precedence: ")
	out.WriteString(strconv.Itoa(li.precedence))
	out.WriteString("]")

	return out.String()
}
