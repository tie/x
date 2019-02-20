package tokenh

import (
	"fmt"

	"github.com/tie/x/config/token"
)

func Sep(val string, pos, end string) token.Token {
	return tok(token.SepToken, val, pos, end)
}

func Space(val string, pos, end string) token.Token {
	return tok(token.SpaceToken, val, pos, end)
}

func Text(val string, pos, end string) token.Token {
	return tok(token.TextToken, val, pos, end)
}

func Comment(val string, pos, end string) token.Token {
	return tok(token.CommentToken, val, pos, end)
}

func tok(typ token.TokenType, val string, a, b string) token.Token {
	return token.Token{
		Typ: typ,
		Val: val,
		Pos: pos(a),
		End: pos(b),
	}
}

func pos(s string) token.Position {
	var off, line, col int
	_, err := fmt.Sscanf(s, "%d:%d(%d)", &line, &col, &off)
	if err != nil {
		panic(err)
	}
	line -= 1
	col -= 1
	return token.Position{
		Line: line,
		Column: col,
		Offset: off,
	}
}
