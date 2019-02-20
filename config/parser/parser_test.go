package parser

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
)

type (
	ParserTest struct {
		Name string
		Input []testingh.ReadRune
		Passes []ParserTestPass
	}
	ParserTestPass func(*testing.T, *Parser)
)

func RunParserTests(t *testing.T, cases []ParserTest) {
	for _, c := range cases {
		r := testingh.NewRuneReader(c.Input)
		p, passes := NewParser(r), c.Passes
		t.Run(c.Name, func(t *testing.T) {
			for _, pass := range passes {
				pass(t, p)
			}
		})
	}
}

func expectLines(tokLines []TokenLine) ParserTestPass {
	return func(t *testing.T, p *Parser) {
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
	}
}

func expectEOF(t *testing.T, p *Parser) {
	tokLine, err := p.NextLine()
	if err != io.EOF {
		if err != nil {
			t.Fatalf("expected EOF error, got %s error", err)
		} else {
			t.Fatalf("expected EOF error, got %s token line", tokLine)
		}
	}
}
