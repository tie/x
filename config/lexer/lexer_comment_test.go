package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/internal/tokenh"
)

func TestLexerComment(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "EOF",
			Input: []testingh.ReadRune{
				// "#"
				{Rune: '#', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens(
					tokenh.Comment("#", "1:1(+0)", "1:2(+1)"),
				),
				expectEOF,
			},
		},
		{
			Name: "Sep",
			Input: []testingh.ReadRune{
				// "#\n"
				{Rune: '#', Size: 1},
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens(
					tokenh.Comment("#", "1:1(+0)", "1:2(+1)"),
					tokenh.Sep("\n", "1:2(+1)", "2:1(+2)"),
				),
				expectEOF,
			},
		},
		{
			Name: "Space",
			Input: []testingh.ReadRune{
				// "# "
				{Rune: '#', Size: 1},
				{Rune: ' ', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens(
					tokenh.Comment("# ", "1:1(+0)", "1:3(+2)"),
				),
				expectEOF,
			},
		},
	})
}
