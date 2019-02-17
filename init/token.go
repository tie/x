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

// TokenLine is a non-empty sequence of text tokens.
type TokenLine []Token

func (toks TokenLine) Directive() Token {
	if len(toks) <= 0 {
		panic("empty token line")
	}
	return toks[0]
}

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

type Position struct {
	Offset, Line, Column int
}

func Pos(s string) Position {
	var off, line, col int
	_, err := fmt.Sscanf(s, "%d:%d(%d)", &line, &col, &off)
	if err != nil {
		panic(err)
	}
	if line <= 0 {
		panic("invalid line, must be >= 1")
	}
	if col <= 0 {
		panic("invalid column, must be >= 1")
	}
	if off < 0 {
		panic("invalid offset, must be >= 0")
	}
	line -= 1
	col -= 1
	return Position{off, line, col}
}

func (p Position) String() string {
	return fmt.Sprintf(
		"%d:%d(%+d)",
		p.Line+1, p.Column+1,
		p.Offset,
	)
}
