package config

import (
	"io"

	"github.com/tie/x/config/parser"
)

func Parse(r io.RuneReader) (parser.Unit, error) {
	return parser.Parse(r, Lexicon)
}
