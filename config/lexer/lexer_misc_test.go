package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestLexerMisc(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "Empty",
			Input: []testingh.ReadRune{
				// ""
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectEOF,
			},
		},
		{
			Name: "SepEOF",
			Input: []testingh.ReadRune{
				// "\n"
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.SepToken, "\n", "1:1(+0)", "2:1(+1)"),
				}),
				expectEOF,
			},
		},
	})
}
