package parser

import (
	"io"

	"github.com/tie/x/config/lexer"
	"github.com/tie/x/config/token"
)

type (
	Unit []Section
	Section []Line
	Line []string
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(r io.RuneReader) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(r),
	}
}

func (p *Parser) NextLine() (TokenLine, error) {
	toks := TokenLine{}
	for {
		tok, err := p.lexer.NextToken()
		if err != nil {
			// suppress eof if line is not empty
			if err == io.EOF && len(toks) > 0 {
				err = nil
			}
			return toks, err
		}
		switch tok.Typ {
		case token.SepToken:
			// skip empty lines
			if len(toks) <= 0 {
				continue
			}
			return toks, nil
		case token.TextToken:
			// line folding
			if tok.Val == "\\\n" {
				continue
			}
			toks = append(toks, tok)
		}
	}
}

func Parse(r io.RuneReader, lex UnitLexicon) (Unit, error) {
	p := NewParser(r)
	var (
		unit Unit
		section Section
		sectionLex SectionLexicon
	)
	for {
		toks, err := p.NextLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				// don't forget to emit non-empty section on eof
				if len(section) > 0 {
					unit = append(unit, section)
				}
			}
			return unit, err
		}
		dirTok := toks.Directive()
		var expand ExpandFunc
		if newSectionLex, ok := lex[dirTok.Val]; ok {
			// it's a new section
			if len(section) > 0 {
				unit = append(unit, section)
				section = Section{}
			}
			sectionLex = newSectionLex
			expand = sectionLex.ExpandFunc
		} else {
			// it's a section directive
			expand, ok = sectionLex.Directives[dirTok.Val]
			if !ok {
				handleUnknown := sectionLex.UnknownFunc
				// terminate if unknown directive is an error in this context
				if handleUnknown != nil {
					err := handleUnknown(toks)
					if err != nil {
						return unit, err
					}
				}
				// skip if handler is not defined or did not return error
				continue
			}
		}
		line, err := expand(toks)
		if err != nil {
			return unit, err
		}
		section = append(section, line)
	}
}
