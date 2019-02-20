package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/internal/tokenh"
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
				expectTokens(
					tokenh.Text("a", "1:1(+0)", "1:2(+1)"),
				),
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
				expectTokens(
					tokenh.Text("a", "1:1(+0)", "1:2(+1)"),
					tokenh.Sep("\n", "1:2(+1)", "2:1(+2)"),
				),
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
				expectTokens(
					tokenh.Text("a", "1:1(+0)", "1:2(+1)"),
					tokenh.Space(" ", "1:2(+1)", "1:3(+2)"),
					tokenh.Text("b", "1:3(+2)", "1:4(+3)"),
				),
				expectEOF,
			},
		},
	})
}
