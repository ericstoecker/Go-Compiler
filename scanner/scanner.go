package scanner

import "compiler/token"

type Scanner interface {
	NextToken() token.Token
}
