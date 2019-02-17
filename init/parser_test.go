package main

import (
	"github.com/tie/x/testingh"
	"io"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	runParserTestCases(t, map[string]parserTestCase{
		"Empty": {
			"",
			[]TokenLine{},
		},
		"EmptyLine": {
			"\n",
			[]TokenLine{},
		},
		"LineSep": {
			"a\nb",
			[]TokenLine{
				{
					{TextToken, "a", Pos("1:1(+0)"), Pos("1:2(+1)")},
				},
				{
					{TextToken, "b", Pos("2:1(+2)"), Pos("2:2(+3)")},
				},
			},
		},
		// TODO: add more parser tests
	})
}

type parserTestCase struct {
	text string
	toks []TokenLine
}

func runParserTestCases(t *testing.T, cases map[string]parserTestCase) {
	run := func(t *testing.T, text string, tokLines []TokenLine) {
		r := testingh.NewTestReader(t, strings.NewReader(text))
		p := NewParser(r)
		for _, toks := range tokLines {
			ntoks, err := p.NextLine()
			if err != nil {
				t.Fatalf("expected %s token line, got %s error", toks, err)
			}
			if len(ntoks) != len(toks) {
				t.Fatalf("expected %s token line, got %s token line", toks, ntoks)
			}
			errs := 0
			for i, tok := range toks {
				ntok := ntoks[i]
				if ntok != tok {
					t.Errorf("expected %s token, got %s token", tok, ntok)
					errs++
				}
			}
			if errs > 0 {
				t.FailNow()
			}
		}
		ntoks, err := p.NextLine()
		if err != io.EOF {
			if err != nil {
				t.Fatalf("expected EOF error, got %s error", err)
			} else {
				t.Fatalf("expected EOF error, got %s token line", ntoks)
			}
		}
	}
	for name, params := range cases {
		text, tokLines := params.text, params.toks
		t.Run(name, func(t *testing.T) {
			run(t, text, tokLines)
		})
	}
}
