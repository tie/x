package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestLexerSpace(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "EOF",
			Input: []testingh.ReadRune{
				// " "
				{Rune: ' ', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.SpaceToken, " ", "1:1(+0)", "1:2(+1)"),
				}),
				expectEOF,
			},
		},
		// regression: did not terminate space token at the end of line
		{
			Name: "SepSpace",
			Input: []testingh.ReadRune{
				// " \n"
				{Rune: ' ', Size: 1},
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.SpaceToken, " ", "1:1(+0)", "1:2(+1)"),
					token.Tok(token.SepToken, "\n", "1:2(+1)", "2:1(+2)"),
				}),
				expectEOF,
			},
		},
	})
}
