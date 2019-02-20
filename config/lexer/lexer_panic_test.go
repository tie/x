package lexer

import (
	"testing"

	"github.com/tie/x/config/internal/testingh"
)

func TestLexerPanic(t *testing.T) {
	RunLexerTests(t, []LexerTest{
		{
			Name: "AcceptWithoutPeek",
			Passes: []LexerTestPass{
				expectPanicOnAccept,
			},
		},
		{
			Name: "ReadThenAcceptWithoutPeek",
			Input: []testingh.ReadRune{
				{Rune: 'a', Size: 1},
			},
			Passes: []LexerTestPass{
				expectReadSuccess,
				expectPanicOnAccept,
			},
		},
		{
			Name: "PeekThenAcceptAfterError",
			Passes: []LexerTestPass{
				expectPeekError,
				expectPanicOnAccept,
			},
		},
	})
}
