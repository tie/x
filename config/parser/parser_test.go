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

func expectStatements(stmtList []Statement) ParserTestPass {
	return func(t *testing.T, p *Parser) {
		for _, stmt := range stmtList {
			nstmt, err := p.NextStatement()
			if err != nil {
				t.Fatalf("expected %s statement, got %s error", stmt, err)
			}
			if len(nstmt) != len(stmt) {
				t.Fatalf("expected %s statement, got %s statement", stmt, nstmt)
			}
			errs := 0
			for i, tok := range stmt {
				ntok := nstmt[i]
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
	stmt, err := p.NextStatement()
	if err != io.EOF {
		if err != nil {
			t.Fatalf("expected EOF error, got %s error", err)
		} else {
			t.Fatalf("expected EOF error, got %s statement", stmt)
		}
	}
}
