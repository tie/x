package config

import (
	"github.com/tie/x/config/parser"
)

var Lexicon = parser.UnitLexicon{
	"import": {
		ExpandFunc: dummyExpand,
	},
	"on": {
		ExpandFunc: dummyExpand,
	},
	"service": {
		ExpandFunc: dummyExpand,
		Directives: map[string]parser.ExpandFunc{
			"class": dummyExpand,
		},
	},
}

func dummyExpand(tokLine parser.TokenLine) (parser.Line, error) {
	var line parser.Line
	for _, tok := range tokLine {
		line = append(line, tok.Val)
	}
	return line, nil
}
