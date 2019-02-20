package token

import (
	"fmt"
)

type TokenType string

const (
	SepToken = "separator"
	SpaceToken = "space"
	TextToken = "text"
	CommentToken = "comment"
)

type Token struct {
	Typ TokenType
	Val string
	Pos Position
	End Position
}

func (t Token) String() string {
	return fmt.Sprintf(
		"[%v %q %v %v (%+d)]",
		t.Typ, t.Val,
		t.Pos, t.End,
		t.End.Offset-t.Pos.Offset,
	)
}
