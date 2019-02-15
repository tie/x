package main

import (
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
				{SepToken, "\n", Position{0, 0, 0}, Position{1, 1, 0}},
			},
		},
	})
}

func TestText(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"a",
			[]Token{
				{TextToken, "a", Position{0, 0, 0}, Position{1, 0, 1}},
			},
		},
		"Sep": {
			"a\n",
			[]Token{
				{TextToken, "a", Position{0, 0, 0}, Position{1, 0, 1}},
				{SepToken, "\n", Position{1, 0, 1}, Position{2, 1, 0}},
			},
		},
		"Space": {
			"a b",
			[]Token{
				{TextToken, "a", Position{0, 0, 0}, Position{1, 0, 1}},
				{SpaceToken, " ", Position{1, 0, 1}, Position{2, 0, 2}},
				{TextToken, "b", Position{2, 0, 2}, Position{3, 0, 3}},
			},
		},
	})
}

func TestSpace(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			" ",
			[]Token{
				{SpaceToken, " ", Position{0, 0, 0}, Position{1, 0, 1}},
			},
		},
		// regression: did not terminate space token on end of line
		"SepSpace": {
			" \n",
			[]Token{
				{SpaceToken, " ", Position{0, 0, 0}, Position{1, 0, 1}},
				{SepToken, "\n", Position{1, 0, 1}, Position{2, 1, 0}},
			},
		},
	})
}

func TestComment(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"#",
			[]Token{
				{CommentToken, "#", Position{0, 0, 0}, Position{1, 0, 1}},
			},
		},
		"Sep": {
			"#\n",
			[]Token{
				{CommentToken, "#", Position{0, 0, 0}, Position{1, 0, 1}},
				{SepToken, "\n", Position{1, 0, 1}, Position{2, 1, 0}},
			},
		},
		"Space": {
			"# ",
			[]Token{
				{CommentToken, "# ", Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
	})
}

func TestTextEscaping(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"\\",
			[]Token{
				{TextToken, "\\", Position{0, 0, 0}, Position{1, 0, 1}},
			},
		},
		"Sep": {
			"\\\n",
			[]Token{
				{TextToken, "\\\n", Position{0, 0, 0}, Position{2, 1, 0}},
			},
		},
		"Space": {
			`\ `,
			[]Token{
				{TextToken, `\ `, Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
		"Comment": {
			`\#`,
			[]Token{
				{TextToken, `\#`, Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
		"Text": {
			`\text`,
			[]Token{
				{TextToken, `\text`, Position{0, 0, 0}, Position{5, 0, 5}},
			},
		},
		"QuoteAndEOF": {
			`\"`,
			[]Token{
				{TextToken, `\"`, Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
		// we also add sep because EOF terminates quoted string
		"QuoteNotEOF": {
			`\"` + "\n",
			[]Token{
				{TextToken, `\"`, Position{0, 0, 0}, Position{2, 0, 2}},
				{SepToken, "\n", Position{2, 0, 2}, Position{3, 1, 0}},
			},
		},
		"Escape": {
			"\\\\",
			[]Token{
				{TextToken, "\\\\", Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
	})
}

func TestTextQuotes(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			`"`,
			[]Token{
				{TextToken, `"`, Position{0, 0, 0}, Position{1, 0, 1}},
			},
		},
		"Empty": {
			`""`,
			[]Token{
				{TextToken, `""`, Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
		"Sep": {
			"\"\n\"",
			[]Token{
				{TextToken, `"`, Position{0, 0, 0}, Position{1, 0, 1}},
				{SepToken, "\n", Position{1, 0, 1}, Position{2, 1, 0}},
				{TextToken, `"`, Position{2, 1, 0}, Position{3, 1, 1}},
			},
		},
		"Space": {
			`" "`,
			[]Token{
				{TextToken, `" "`, Position{0, 0, 0}, Position{3, 0, 3}},
			},
		},
		"TextWithSpaces": {
			`" a "`,
			[]Token{
				{TextToken, `" a "`, Position{0, 0, 0}, Position{5, 0, 5}},
			},
		},
		"Escape": {
			`" \" "`,
			[]Token{
				{TextToken, `" \" "`, Position{0, 0, 0}, Position{6, 0, 6}},
			},
		},
		"EscapeEOF": {
			"\"\\",
			[]Token{
				{TextToken, "\"\\", Position{0, 0, 0}, Position{2, 0, 2}},
			},
		},
		"Multiple": {
			`" " " " "`,
			[]Token{
				{TextToken, `" "`, Position{0, 0, 0}, Position{3, 0, 3}},
				{SpaceToken, ` `, Position{3, 0, 3}, Position{4, 0, 4}},
				{TextToken, `" "`, Position{4, 0, 4}, Position{7, 0, 7}},
				{SpaceToken, ` `, Position{7, 0, 7}, Position{8, 0, 8}},
				{TextToken, `"`, Position{8, 0, 8}, Position{9, 0, 9}},
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
		r := newTestReader(t, strings.NewReader(text))
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
		tok, err := l.NextToken()
		if err != io.EOF {
			if err != nil {
				t.Fatalf("expected EOF error, got %s error", err)
			} else {
				t.Fatalf("expected EOF error, got %s token", tok)
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

// testReader calls testing.Error() if io.Reader.Read() that returned an error is followed by another io.Reader.Read() call.
type testReader struct {
	Reader io.RuneReader
	Test *testing.T
	Err error
}

func newTestReader(t *testing.T, r io.RuneReader) *testReader {
	return &testReader{r, t, nil}
}

func (r *testReader) ReadRune() (c rune, size int, err error) {
	r.Test.Helper()
	if r.Err != nil {
		r.Test.Errorf("Read called after %q error", r.Err)
	}
	c, size, err = r.Reader.ReadRune()
	if err != nil {
		if r.Err != nil {
			r.Test.Errorf(
				"Read returned different %q error after %q",
				err, r.Err,
			)
		}
		r.Err = err
	} else {
		if r.Err != nil {
			r.Test.Errorf("Read succeeded after %q error", r.Err)
			r.Err = nil
		}
	}
	return
}
