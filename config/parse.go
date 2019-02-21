package config

import (
	"io"

	"github.com/tie/x/config/parser"
)

func Parse(r io.RuneReader) (parser.Unit, error) {
	return parser.Parse(r, syntax)
}

var syntax = parser.Syntax{
	TopLevel: dummyCheck,
	Sections: map[string]parser.CheckFunc{
		"on": dummyCheck,
		"import": dummyCheck,
		"service": dummyCheck,
	},
}

func dummyCheck(stmt parser.Statement) error {
	return nil
}
