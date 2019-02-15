package main

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
	return fmt.Sprintf("[%v %q %v:%v %v:%v]", t.Typ, t.Val, t.Pos.Line, t.Pos.Column, t.End.Line, t.End.Column)
}

type Position struct {
	Offset, Line, Column int
}

func (p Position) String() string {
	return fmt.Sprintf("%v:%v (byte %v)", p.Line+1, p.Column+1, p.Offset+1)
}
