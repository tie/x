package parser

import (
	"io"

	"github.com/tie/x/config/lexer"
	"github.com/tie/x/config/token"
)

// Unit is a sequence of sections.  The first element is always a top-level section.
type Unit []Section

func (u Unit) TopLevel() Section {
	return u[0]
}

func (u Unit) Sections() []Section {
	return u[1:]
}

// Section is a sequence of statements.  Such sequence may be empty iff section has no explicit header.
type Section []Statement

// Statement is a non-empty sequence of text tokens.
type Statement []token.Token

func (s Statement) Directive() string {
	return s[0].Val
}

// CheckFunc checks syntax of a statement and reports errors.
type CheckFunc func(stmt Statement) error

// Syntax defines rules for checking syntax of sections.
// Rules are defined by CheckFunc check functions.
type Syntax struct {
	// TopLevel checks syntax of top-level statements, i.e. those without explicit section.
	TopLevel CheckFunc
	// Sections maps section header directive keyword to the section syntax checker.
	Sections map[string]CheckFunc
}

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(r io.RuneReader) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(r),
	}
}

func (p *Parser) NextStatement() (Statement, error) {
	stmt := Statement{}
	for {
		tok, err := p.lexer.NextToken()
		if err != nil {
			// suppress eof if statement is not empty
			if err == io.EOF && len(stmt) > 0 {
				err = nil
			}
			return stmt, err
		}
		switch tok.Typ {
		case token.SepToken:
			// skip empty statements
			if len(stmt) <= 0 {
				continue
			}
			return stmt, nil
		case token.TextToken:
			// line folding
			if tok.Val == "\\\n" {
				continue
			}
			stmt = append(stmt, tok)
		}
	}
}

func Parse(r io.RuneReader, syn Syntax) (unit Unit, err error) {
	var section Section
	p := NewParser(r)
	check := syn.TopLevel
	for {
		stmt, err := p.NextStatement()
		if err != nil {
			if err == io.EOF {
				err = nil
				// don't forget to emit on eof
				if len(section) > 0 || len(unit) <= 0 {
					unit = append(unit, section)
				}
			}
			return unit, err
		}
		dir := stmt.Directive()
		if nextSectionCheck, ok := syn.Sections[dir]; ok {
			// emit top-level or non-empty sections
			if len(section) > 0 || len(unit) <= 0 {
				unit = append(unit, section)
				section = Section{}
			}
			check = nextSectionCheck
		}
		if check != nil {
			if err := check(stmt); err != nil {
				return unit, err
			}
			section = append(section, stmt)
		}
	}
}
