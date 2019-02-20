package parser

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

func TestParserMisc(t *testing.T) {
	RunParserTests(t, []ParserTest{
		{
			Name: "Empty",
			Input: []testingh.ReadRune{
				// ""
				{Error: io.EOF},
			},
			Passes: []ParserTestPass{
				expectEOF,
			},
		},
		{
			Name: "EmptyLine",
			Input: []testingh.ReadRune{
				// "\n"
				{Rune: '\n', Size: 1},
				{Error: io.EOF},
			},
			Passes: []ParserTestPass{
				expectEOF,
			},
		},
		{
			Name: "LineSep",
			Input: []testingh.ReadRune{
				// "a\nb"
				{Rune: 'a', Size: 1},
				{Rune: '\n', Size: 1},
				{Rune: 'b', Size: 1},
				{Error: io.EOF},
			},
			Passes: []ParserTestPass{
				expectLines([]TokenLine{
					{
						token.Tok(token.TextToken, "a", "1:1(+0)", "1:2(+1)"),
					},
					{
						token.Tok(token.TextToken, "b", "2:1(+2)", "2:2(+3)"),
					},
				}),
				expectEOF,
			},
		},
		// TODO: add more parser tests
	})
}
