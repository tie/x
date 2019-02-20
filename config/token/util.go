package token

import (
	"fmt"
)

func Tok(typ TokenType, val string, pos, end string) Token {
	return Token{typ, val, Pos(pos), Pos(end)}
}

func Pos(s string) Position {
	var off, line, col int
	_, err := fmt.Sscanf(s, "%d:%d(%d)", &line, &col, &off)
	if err != nil {
		panic(err)
	}
	line -= 1
	col -= 1
	return Position{off, line, col}
}
