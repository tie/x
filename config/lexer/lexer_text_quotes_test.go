package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestLexerTextQuotes(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "EOF",
			Input: []testingh.ReadRune{
				// "\""
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\"", "1:1(+0)", "1:2(+1)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Empty",
			Input: []testingh.ReadRune{
				// "\"\""
				{Rune: '"', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\"\"", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Sep",
			Input: []testingh.ReadRune{
				// "\"\n\""
				{Rune: '"', Size: 1},
				{Rune: '\n', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\"", "1:1(+0)", "1:2(+1)"),
					token.Tok(token.SepToken, "\n", "1:2(+1)", "2:1(+2)"),
					token.Tok(token.TextToken, "\"", "2:1(+2)", "2:2(+3)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Space",
			Input: []testingh.ReadRune{
				// "\" \""
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\" \"", "1:1(+0)", "1:4(+3)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "TextWithSpaces",
			Input: []testingh.ReadRune{
				// "\" a \""
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: 'a', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\" a \"", "1:1(+0)", "1:6(+5)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Escape",
			Input: []testingh.ReadRune{
				// "\" \\\" \""
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '\\', Size: 1},
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\" \\\" \"", "1:1(+0)", "1:7(+6)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "EscapeEOF",
			Input: []testingh.ReadRune{
				// "\"\\"
				{Rune: '"', Size: 1},
				{Rune: '\\', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\"\\", "1:1(+0)", "1:3(+2)"),
				}),
				expectEOF,
			},
		},
		{
			Name: "Multiple",
			Input: []testingh.ReadRune{
				// "\" \" \" \" \""
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Rune: ' ', Size: 1},
				{Rune: '"', Size: 1},
				{Error: io.EOF},
			},
			Passes: []LexerTestPass{
				expectTokens([]token.Token{
					token.Tok(token.TextToken, "\" \"", "1:1(+0)", "1:4(+3)"),
					token.Tok(token.SpaceToken, " ", "1:4(+3)", "1:5(+4)"),
					token.Tok(token.TextToken, "\" \"", "1:5(+4)", "1:8(+7)"),
					token.Tok(token.SpaceToken, " ", "1:8(+7)", "1:9(+8)"),
					token.Tok(token.TextToken, "\"", "1:9(+8)", "1:10(+9)"),
				}),
				expectEOF,
			},
		},
	})
}
