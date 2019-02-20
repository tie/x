package parser

import (
	"github.com/tie/x/config/token"
)

// TokenLine is a non-empty sequence of text tokens.
type TokenLine []token.Token

func (toks TokenLine) Directive() token.Token {
	if len(toks) <= 0 {
		panic("empty token line")
	}
	return toks[0]
}
