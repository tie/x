package main

type (
	// ExpandFunc handles parsing of token line.
	ExpandFunc func(toks TokenLine) (Line, error)
	// UnknownFunc handles lines with unknown directives.
	UnknownFunc func(toks TokenLine) error
	// SectionLexicon defines rules for parsing Section.
	SectionLexicon struct {
		// ExpandFunc handles parsing of token line.
		ExpandFunc
		// UnknownFunc handles lines with unknown directives.
		//
		// Lines with unknown directives are skipped if:
		// - Unknown is nil
		// - Unknown() returns nil
		// Parsing terminates if:
		// - Unknown() returns non-nil error
		UnknownFunc
		Directives map[string]ExpandFunc
	}
	// UnitLexicon defines rules for parsing Unit.
	UnitLexicon map[string]SectionLexicon
)

var Lexicon = UnitLexicon{
	"import": {
		ExpandFunc: dummyExpand,
	},
	"on": {
		ExpandFunc: dummyExpand,
	},
	"service": {
		ExpandFunc: dummyExpand,
		Directives: map[string]ExpandFunc{
			"class": dummyExpand,
		},
	},
}

func dummyExpand(toks TokenLine) (Line, error) {
	var line Line
	for _, tok := range toks {
		line = append(line, tok.Val)
	}
	return line, nil
}
