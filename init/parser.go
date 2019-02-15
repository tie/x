package main

import (
	"io"
)

type Parser struct {
	lexer *Lexer
}

func NewParser(r io.RuneReader) *Parser {
	return &Parser{
		lexer: NewLexer(r),
	}
}

func (p *Parser) parseLine() (Line, error) {
	line := Line{}
	for {
		tok, err := p.lexer.NextToken()
		if err != nil {
			// suppress eof if line is not empty
			if err == io.EOF && len(line) > 0 {
				err = nil
			}
			return line, err
		}
		switch tok.Typ {
		case SepToken:
			if len(line) > 0 {
				return line, nil
			}
		case TextToken:
			// line folding
			if tok.Val == "\\\n" {
				continue
			}
			line = append(line, tok)
		case CommentToken, SpaceToken:
			continue
		}
	}
}
