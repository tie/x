package main

import (
	"io"
	"strings"
	"unicode"
)

type Lexer struct {
	reader io.RuneReader
	buffer strings.Builder
	startPos, endPos Position
	next struct {
		r rune
		size int
		err error
	}
}

func NewLexer(r io.RuneReader) *Lexer {
	return &Lexer{
		reader: r,
	}
}

func (l *Lexer) NextToken() (Token, error) {
	r, err := l.peek()
	if err != nil {
		return Token{}, err
	}
	tok, err := l.nextState(r)
	if err == io.EOF {
		err = nil
	}
	return tok, err
}

func (l *Lexer) nextState(r rune) (Token, error) {
	if r == '#' {
		return l.commentState()
	}
	if r == '\n' {
		return l.sepState()
	}
	if unicode.IsSpace(r) {
		return l.spacesState()
	}
	return l.textState()
}

func (l *Lexer) emit(typ TokenType) Token {
	value := l.buffer.String()
	tok := Token{typ, value, l.startPos, l.endPos}
	l.startPos = l.endPos
	l.buffer.Reset()
	return tok
}

func (l *Lexer) peek() (rune, error) {
	r, size, err := l.next.r, l.next.size, l.next.err
	if err != nil || size > 0 {
		return r, err
	}
	r, size, err = l.reader.ReadRune()
	l.next.r, l.next.size, l.next.err = r, size, err
	return r, err
}

func (l *Lexer) accept() {
	r, size, err := l.next.r, l.next.size, l.next.err
	if err != nil || size <= 0 {
		// it's a bug: accept without peek or after error
		panic("nothing to accept")
	}
	l.buffer.WriteRune(r)
	l.endPos.Offset += size
	switch r {
	case '\n':
		l.endPos.Column = 0
		l.endPos.Line++
	default:
		l.endPos.Column++
	}
	// clear next rune
	l.next.size = 0
}

func (l *Lexer) read() (rune, error) {
	r, err := l.peek()
	if err != nil {
		return r, err
	}
	l.accept()
	return r, err
}

func (l *Lexer) sepState() (Token, error) {
	// assume separator (i.e. end of line)
	l.accept()
	return l.emit(SepToken), nil
}

func (l *Lexer) spacesState() (Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(SpaceToken), err
		}
		if !unicode.IsSpace(r) || r == '\n' {
			return l.emit(SpaceToken), nil
		}
		l.accept()
	}
}

func (l *Lexer) textState() (Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(TextToken), err
		}
		if unicode.IsSpace(r) || r == '#' {
			return l.emit(TextToken), nil
		}
		switch r {
		case '\\':
			// escape character
			l.accept()
			r, err := l.read()
			if err != nil {
				return l.emit(TextToken), err
			}
			// and terminate at the end of line
			if r == '\n' {
				return l.emit(TextToken), nil
			}
			continue
		case '"':
			err := l.quoteText()
			if err != nil {
				return l.emit(TextToken), err
			}
			continue
		}
		l.accept()
	}
}

func (l *Lexer) quoteText() error {
	// assume it's a quote character
	r, err := l.read()
	if err != nil {
		// it's a bug: quoteText must be called after peek
		panic(err)
	}
	quote := r
	for {
		r, err := l.peek()
		if err != nil {
			// unterminated quoted thing
			return err
		}
		switch r {
		case '\\':
			// escape character
			l.accept()
			_, err := l.read()
			if err != nil {
				return err
			}
			continue
		case '\n':
			// terminate at end of line
			return nil
		case quote:
			// closing quote
			l.accept()
			return nil
		}
		l.accept()
	}
}

func (l *Lexer) commentState() (Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(CommentToken), err
		}
		if r == '\n' {
			return l.emit(CommentToken), nil
		}
		l.accept()
	}
}
