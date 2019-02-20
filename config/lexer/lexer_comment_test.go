package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
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
				expectTokens([]token.Token{
					token.Tok(token.CommentToken, "#", "1:1(+0)", "1:2(+1)"),
				}),
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
				expectTokens([]token.Token{
					token.Tok(token.CommentToken, "#", "1:1(+0)", "1:2(+1)"),
					token.Tok(token.SepToken, "\n", "1:2(+1)", "2:1(+2)"),
				}),
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
				expectTokens([]token.Token{
					token.Tok(token.CommentToken, "# ", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
	})
}
