package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/internal/tokenh"
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
				expectTokens(
					tokenh.Text("\"", "1:1(+0)", "1:2(+1)"),
				),
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
				expectTokens(
					tokenh.Text("\"\"", "1:1(+0)", "1:3(+2)"),
				),
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
				expectTokens(
					tokenh.Text("\"", "1:1(+0)", "1:2(+1)"),
					tokenh.Sep("\n", "1:2(+1)", "2:1(+2)"),
					tokenh.Text("\"", "2:1(+2)", "2:2(+3)"),
				),
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
				expectTokens(
					tokenh.Text("\" \"", "1:1(+0)", "1:4(+3)"),
				),
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
				expectTokens(
					tokenh.Text("\" a \"", "1:1(+0)", "1:6(+5)"),
				),
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
				expectTokens(
					tokenh.Text("\" \\\" \"", "1:1(+0)", "1:7(+6)"),
				),
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
				expectTokens(
					tokenh.Text("\"\\", "1:1(+0)", "1:3(+2)"),
				),
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
				expectTokens(
					tokenh.Text("\" \"", "1:1(+0)", "1:4(+3)"),
					tokenh.Space(" ", "1:4(+3)", "1:5(+4)"),
					tokenh.Text("\" \"", "1:5(+4)", "1:8(+7)"),
					tokenh.Space(" ", "1:8(+7)", "1:9(+8)"),
					tokenh.Text("\"", "1:9(+8)", "1:10(+9)"),
				),
				expectEOF,
			},
		},
	})
}
