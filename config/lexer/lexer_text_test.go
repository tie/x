package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestLexerText(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "EOF",
			Input: []testingh.ReadRune{
				// "a"
				{Rune: 'a', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "a", "1:1(+0)", "1:2(+1)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Sep",
			Input: []testingh.ReadRune{
				// "a\n"
				{Rune: 'a', Size: 1},
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "a", "1:1(+0)", "1:2(+1)"),
					token.Tok(token.SepToken, "\n", "1:2(+1)", "2:1(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Space",
			Input: []testingh.ReadRune{
				// "a b"
				{Rune: 'a', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: 'b', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "a", "1:1(+0)", "1:2(+1)"),
					token.Tok(token.SpaceToken, " ", "1:2(+1)", "1:3(+2)"),
					token.Tok(token.TextToken, "b", "1:3(+2)", "1:4(+3)"),
				}),
				expectEOF,
			},
		},
	})
}
