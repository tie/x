package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestLexerTextEscaping(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "EOF",
			Input: []testingh.ReadRune{
				// "\\"
				{Rune: '\\', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\", "1:1(+0)", "1:2(+1)"),
				}),
				expectEOF,
			},
		},
		// line folding: terminate text token even when escaping end of line, though don't emit SepToken
		{
			Name: "Sep",
			Input: []testingh.ReadRune{
				// "\\\n\\\n"
				{Rune: '\\', Size: 1},
				{Rune: '\n', Size: 1},
				{Rune: '\\', Size: 1},
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\\n", "1:1(+0)", "2:1(+2)"),
					token.Tok(token.TextToken, "\\\n", "2:1(+2)", "3:1(+4)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Space",
			Input: []testingh.ReadRune{
				// "\\ "
				{Rune: '\\', Size: 1},
				{Rune: ' ', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\ ", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Comment",
			Input: []testingh.ReadRune{
				// "\\#"
				{Rune: '\\', Size: 1},
				{Rune: '#', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\#", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Text",
			Input: []testingh.ReadRune{
				// "\\n"
				{Rune: '\\', Size: 1},
				{Rune: 'n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\n", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "QuoteAndEOF",
			Input: []testingh.ReadRune{
				// "\\\""
				{Rune: '\\', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\\"", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		// we also add sep because EOF terminates quoted string
		{
			Name: "QuoteNotEOF",
			Input: []testingh.ReadRune{
				// "\\\"\n"
				{Rune: '\\', Size: 1},
				{Rune: '"', Size: 1},
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\\"", "1:1(+0)", "1:3(+2)"),
					token.Tok(token.SepToken, "\n", "1:3(+2)", "2:1(+3)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Escape",
			Input: []testingh.ReadRune{
				// "\\\\"
				{Rune: '\\', Size: 1},
				{Rune: '\\', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\\\\", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
	})
}
