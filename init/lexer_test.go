package main

import (
	"github.com/tie/x/testingh"
	"io"
	"strings"
	"testing"
)

func TestSpecialCases(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"Empty": {
			"",
			[]Token{},
		},
		"SepEOF": {
			"\n",
			[]Token{
				{SepToken, "\n", Pos("1:1(+0)"), Pos("2:1(+1)")},
			},
		},
	})
}

func TestText(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"a",
			[]Token{
				{TextToken, "a", Pos("1:1(+0)"), Pos("1:2(+1)")},
			},
		},
		"Sep": {
			"a\n",
			[]Token{
				{TextToken, "a", Pos("1:1(+0)"), Pos("1:2(+1)")},
				{SepToken, "\n", Pos("1:2(+1)"), Pos("2:1(+2)")},
			},
		},
		"Space": {
			"a b",
			[]Token{
				{TextToken, "a", Pos("1:1(+0)"), Pos("1:2(+1)")},
				{SpaceToken, " ", Pos("1:2(+1)"), Pos("1:3(+2)")},
				{TextToken, "b", Pos("1:3(+2)"), Pos("1:4(+3)")},
			},
		},
	})
}

func TestSpace(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			" ",
			[]Token{
				{SpaceToken, " ", Pos("1:1(+0)"), Pos("1:2(+1)")},
			},
		},
		// regression: did not terminate space token at the end of line
		"SepSpace": {
			" \n",
			[]Token{
				{SpaceToken, " ", Pos("1:1(+0)"), Pos("1:2(+1)")},
				{SepToken, "\n", Pos("1:2(+1)"), Pos("2:1(+2)")},
			},
		},
	})
}

func TestComment(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"#",
			[]Token{
				{CommentToken, "#", Pos("1:1(+0)"), Pos("1:2(+1)")},
			},
		},
		"Sep": {
			"#\n",
			[]Token{
				{CommentToken, "#", Pos("1:1(+0)"), Pos("1:2(+1)")},
				{SepToken, "\n", Pos("1:2(+1)"), Pos("2:1(+2)")},
			},
		},
		"Space": {
			"# ",
			[]Token{
				{CommentToken, "# ", Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
	})
}

func TestTextEscaping(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"\\",
			[]Token{
				{TextToken, "\\", Pos("1:1(+0)"), Pos("1:2(+1)")},
			},
		},
		// line folding: terminate text token even when escaping end of line, though don't emit SepToken
		"Sep": {
			"\\\n\\\n",
			[]Token{
				{TextToken, "\\\n", Pos("1:1(+0)"), Pos("2:1(+2)")},
				{TextToken, "\\\n", Pos("2:1(+2)"), Pos("3:1(+4)")},
			},
		},
		"Space": {
			`\ `,
			[]Token{
				{TextToken, `\ `, Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
		"Comment": {
			`\#`,
			[]Token{
				{TextToken, `\#`, Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
		"Text": {
			`\text`,
			[]Token{
				{TextToken, `\text`, Pos("1:1(+0)"), Pos("1:6(+5)")},
			},
		},
		"QuoteAndEOF": {
			`\"`,
			[]Token{
				{TextToken, `\"`, Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
		// we also add sep because EOF terminates quoted string
		"QuoteNotEOF": {
			`\"` + "\n",
			[]Token{
				{TextToken, `\"`, Pos("1:1(+0)"), Pos("1:3(+2)")},
				{SepToken, "\n", Pos("1:3(+2)"), Pos("2:1(+3)")},
			},
		},
		"Escape": {
			"\\\\",
			[]Token{
				{TextToken, "\\\\", Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
	})
}

func TestTextQuotes(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			`"`,
			[]Token{
				{TextToken, `"`, Pos("1:1(+0)"), Pos("1:2(+1)")},
			},
		},
		"Empty": {
			`""`,
			[]Token{
				{TextToken, `""`, Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
		"Sep": {
			"\"\n\"",
			[]Token{
				{TextToken, `"`, Pos("1:1(+0)"), Pos("1:2(+1)")},
				{SepToken, "\n", Pos("1:2(+1)"), Pos("2:1(+2)")},
				{TextToken, `"`, Pos("2:1(+2)"), Pos("2:2(+3)")},
			},
		},
		"Space": {
			`" "`,
			[]Token{
				{TextToken, `" "`, Pos("1:1(+0)"), Pos("1:4(+3)")},
			},
		},
		"TextWithSpaces": {
			`" a "`,
			[]Token{
				{TextToken, `" a "`, Pos("1:1(+0)"), Pos("1:6(+5)")},
			},
		},
		"Escape": {
			`" \" "`,
			[]Token{
				{TextToken, `" \" "`, Pos("1:1(+0)"), Pos("1:7(+6)")},
			},
		},
		"EscapeEOF": {
			"\"\\",
			[]Token{
				{TextToken, "\"\\", Pos("1:1(+0)"), Pos("1:3(+2)")},
			},
		},
		"Multiple": {
			`" " " " "`,
			[]Token{
				{TextToken, `" "`, Pos("1:1(+0)"), Pos("1:4(+3)")},
				{SpaceToken, ` `, Pos("1:4(+3)"), Pos("1:5(+4)")},
				{TextToken, `" "`, Pos("1:5(+4)"), Pos("1:8(+7)")},
				{SpaceToken, ` `, Pos("1:8(+7)"), Pos("1:9(+8)")},
				{TextToken, `"`, Pos("1:9(+8)"), Pos("1:10(+9)")},
			},
		},
	})
}

type testCase struct {
	text string
	toks []Token
}

func runTestCases(t *testing.T, cases map[string]testCase) {
	run := func(t *testing.T, text string, toks []Token) {
		r := testingh.NewTestReader(t, strings.NewReader(text))
		l := NewLexer(r)
		for _, tok := range toks {
			ntok, err := l.NextToken()
			if err != nil {
				t.Fatalf("expected %s token, got %s error", tok, err)
			}
			if ntok != tok {
				t.Fatalf("expected %s token, got %s token", tok, ntok)
			}
		}
		ntok, err := l.NextToken()
		if err != io.EOF {
			if err != nil {
				t.Fatalf("expected EOF error, got %s error", err)
			} else {
				t.Fatalf("expected EOF error, got %s token", ntok)
			}
		}
	}
	for name, params := range cases {
		text, toks := params.text, params.toks
		t.Run(name, func(t *testing.T) {
			run(t, text, toks)
		})
	}
}
