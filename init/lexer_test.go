package main

import (
	"github.com/tie/x/qlex"
	"io"
	"strings"
	"testing"
)

func TestSpecialCases(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"Empty": {
			"",
			[]qlex.Token{},
		},
		"SepEof": {
			"\n",
			[]qlex.Token{
				{SepToken, "\n", qlex.Position{0, 0, 0}, qlex.Position{1, 1, 0}},
			},
		},
		"SpaceEof": {
			" ",
			[]qlex.Token{
				{SpaceToken, " ", qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
			},
		},
		"TextEof": {
			"a",
			[]qlex.Token{
				{TextToken, "a", qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
			},
		},
	})
}

func TestComment(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"#",
			[]qlex.Token{
				{CommentToken, "#", qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
			},
		},
		"Sep": {
			"#\n",
			[]qlex.Token{
				{CommentToken, "#", qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
				{SepToken, "\n", qlex.Position{1, 0, 1}, qlex.Position{2, 1, 0}},
			},
		},
		"Space": {
			"# ",
			[]qlex.Token{
				{CommentToken, "# ", qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
	})
}

func TestTextEscaping(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			"\\",
			[]qlex.Token{
				{TextToken, "\\", qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
			},
		},
		"Sep": {
			"\\\n",
			[]qlex.Token{
				{TextToken, "\\\n", qlex.Position{0, 0, 0}, qlex.Position{2, 1, 0}},
			},
		},
		"Space": {
			`\ `,
			[]qlex.Token{
				{TextToken, `\ `, qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
		"Comment": {
			`\#`,
			[]qlex.Token{
				{TextToken, `\#`, qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
		"Text": {
			`\text`,
			[]qlex.Token{
				{TextToken, `\text`, qlex.Position{0, 0, 0}, qlex.Position{5, 0, 5}},
			},
		},
		"QuoteAndEOF": {
			`\"`,
			[]qlex.Token{
				{TextToken, `\"`, qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
		// we also add sep because EOF terminates quoted string
		"QuoteNotEOF": {
			`\"` + "\n",
			[]qlex.Token{
				{TextToken, `\"`, qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
				{SepToken, "\n", qlex.Position{2, 0, 2}, qlex.Position{3, 1, 0}},
			},
		},
		"Escape": {
			"\\\\",
			[]qlex.Token{
				{TextToken, "\\\\", qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
	})
}

func TestQuotes(t *testing.T) {
	runTestCases(t, map[string]testCase{
		"EOF": {
			`"`,
			[]qlex.Token{
				{TextToken, `"`, qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
			},
		},
		"Empty": {
			`""`,
			[]qlex.Token{
				{TextToken, `""`, qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
		"Sep": {
			"\"\n\"",
			[]qlex.Token{
				{TextToken, `"`, qlex.Position{0, 0, 0}, qlex.Position{1, 0, 1}},
				{SepToken, "\n", qlex.Position{1, 0, 1}, qlex.Position{2, 1, 0}},
				{TextToken, `"`, qlex.Position{2, 1, 0}, qlex.Position{3, 1, 1}},
			},
		},
		"Space": {
			`" "`,
			[]qlex.Token{
				{TextToken, `" "`, qlex.Position{0, 0, 0}, qlex.Position{3, 0, 3}},
			},
		},
		"TextWithSpaces": {
			`" a "`,
			[]qlex.Token{
				{TextToken, `" a "`, qlex.Position{0, 0, 0}, qlex.Position{5, 0, 5}},
			},
		},
		"Escape": {
			`" \" "`,
			[]qlex.Token{
				{TextToken, `" \" "`, qlex.Position{0, 0, 0}, qlex.Position{6, 0, 6}},
			},
		},
		"EscapeEOF": {
			"\"\\",
			[]qlex.Token{
				{TextToken, "\"\\", qlex.Position{0, 0, 0}, qlex.Position{2, 0, 2}},
			},
		},
		"Multiple": {
			`" " " " "`,
			[]qlex.Token{
				{TextToken, `" "`, qlex.Position{0, 0, 0}, qlex.Position{3, 0, 3}},
				{SpaceToken, ` `, qlex.Position{3, 0, 3}, qlex.Position{4, 0, 4}},
				{TextToken, `" "`, qlex.Position{4, 0, 4}, qlex.Position{7, 0, 7}},
				{SpaceToken, ` `, qlex.Position{7, 0, 7}, qlex.Position{8, 0, 8}},
				{TextToken, `"`, qlex.Position{8, 0, 8}, qlex.Position{9, 0, 9}},
			},
		},
	})
}

type testCase struct {
	text string
	toks []qlex.Token
}

func runTestCases(t *testing.T, cases map[string]testCase) {
	run := func(t *testing.T, text string, toks []qlex.Token) {
		r := newTestReader(t, strings.NewReader(text))
		l := qlex.NewLexer(r, initState)
		for _, tok := range toks {
			ntok, eof := l.NextToken()
			if eof {
				t.Fatalf("expected %s, found EOF", tok)
			}
			if ntok != tok {
				t.Fatalf("expected %s, got %s", tok, ntok)
			}
		}
		tok, eof := l.NextToken()
		if !eof {
			t.Fatalf("expected EOF, found %s", tok)
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
	Reader io.Reader
	Test   *testing.T
	Err    error
}

func newTestReader(t *testing.T, r io.Reader) *testReader {
	return &testReader{r, t, nil}
}

func (r *testReader) Read(p []byte) (n int, err error) {
	r.Test.Helper()
	if r.Err != nil {
		r.Test.Errorf("Read called after %q error", r.Err)
	}
	n, err = r.Reader.Read(p)
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
