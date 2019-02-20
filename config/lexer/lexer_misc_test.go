package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/internal/tokenh"
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
				expectTokens(
					tokenh.Sep("\n", "1:1(+0)", "2:1(+1)"),
				),
				expectEOF,
			},
		},
	})
}
