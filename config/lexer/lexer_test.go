package lexer

import (
	"io"
	"testing"

	"github.com/tie/x/config/internal/testingh"
	"github.com/tie/x/config/token"
)

type (
	LexerTest struct {
		Name string
		Input []testingh.ReadRune
		Passes []LexerTestPass
	}
	LexerTestPass func(*testing.T, *Lexer)
)

func RunLexerTests(t *testing.T, cases []LexerTest) {
	for _, c := range cases {
		r := testingh.NewRuneReader(c.Input)
		l, passes := NewLexer(r), c.Passes
		t.Run(c.Name, func(t *testing.T) {
			for _, pass := range passes {
				pass(t, l)
			}
		})
	}
}

func expectTokens(toks []token.Token) LexerTestPass {
	return func(t *testing.T, l *Lexer) {
		for _, tok := range toks {
			ntok, err := l.NextToken()
			if err != nil {
				t.Fatalf("expected %s token, got %s error", tok, err)
			}
			if ntok != tok {
				t.Fatalf("expected %s token, got %s token", tok, ntok)
			}
		}
	}
}

func expectEOF(t *testing.T, l *Lexer) {
	tok, err := l.NextToken()
	if err != io.EOF {
		if err != nil {
			t.Fatalf("expected EOF error, got %s error", err)
		} else {
			t.Fatalf("expected EOF error, got %s token", tok)
		}
	}
}

func expectPanicOnAccept(t *testing.T, l *Lexer) {
	if testingh.DoesNotPanic(l.accept) {
		t.Error("expected panic in Lexer.accept")
	}
}

func expectReadSuccess(t *testing.T, l *Lexer) {
	_, err := l.read()
	if err != nil {
		t.Fatalf("unexpected %s error from Lexer.read", err)
	}
}

func expectPeekError(t *testing.T, l *Lexer) {
	_, err := l.peek()
	if err == nil {
		t.Fatal("expected error from Lexer.peek")
	}
}
