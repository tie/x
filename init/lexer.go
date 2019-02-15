package main

import (
	"io"
	"strings"
)

type StateFunc func() StateFunc

type Lexer struct {
	state StateFunc
	tokens chan Token
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
	l := Lexer{
		reader: r,
		// TODO: lexer without goroutines
		tokens: make(chan Token, 1),
	}
	l.state = l.initState
	return &l
}

func (l *Lexer) Run() {
	for l.state != nil {
		l.state = l.state()
	}
	close(l.tokens)
}

func (l *Lexer) NextToken() (tok Token, eof bool) {
	for l.state != nil {
		select {
		case token := <-l.tokens:
			return token, false
		default:
			l.state = l.state()
		}
	}
	return Token{}, true
}

func (l *Lexer) emit(typ TokenType) {
	value := l.buffer.String()
	l.tokens <- Token{typ, value, l.startPos, l.endPos}
	l.startPos = l.endPos
	l.buffer.Reset()
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

func (l *Lexer) initState() StateFunc {
	r, err := l.peek()
	if err != nil {
		return l.eofState
	}
	switch r {
	case '#':
		return l.commentState
	case '\n':
		return l.sepState
	default:
		if isSpace(r) {
			return l.spacesState
		}
		return l.textState
	}
}

func (l *Lexer) eofState() StateFunc {
	return nil
}

func (l *Lexer) sepState() StateFunc {
	// assume separator (i.e. end of line)
	l.accept()
	l.emit(SepToken)
	return l.initState
}

func (l *Lexer) spacesState() StateFunc {
	for {
		r, err := l.peek()
		if err != nil {
			l.emit(SpaceToken)
			return l.eofState
		}
		if !isSpace(r) {
			l.emit(SpaceToken)
			return l.initState
		}
		l.accept()
	}
}

func (l *Lexer) textState() StateFunc {
	for {
		r, err := l.peek()
		if err != nil {
			l.emit(TextToken)
			return l.eofState
		}
		switch r {
		case '\\':
			if err := l.escapeText(); err != nil {
				l.emit(TextToken)
				return l.eofState
			}
		case '"':
			if err := l.quoteText(); err != nil {
				l.emit(TextToken)
				return l.eofState
			}
		default:
			if !isText(r) {
				l.emit(TextToken)
				return l.initState
			}
			l.accept()
		}
	}
}

func (l *Lexer) escapeText() error {
	// assume it's escape character (i.e. backslash)
	if _, err := l.read(); err != nil {
		return err
	}
	// accept next rune
	_, err := l.read()
	return err
}

func (l *Lexer) quoteText() error {
	// assume it's a quote character
	r, err := l.read()
	if err != nil {
		return err
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
			if err := l.escapeText(); err != nil {
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

func (l *Lexer) commentState() StateFunc {
	for {
		r, err := l.peek()
		if err != nil {
			l.emit(CommentToken)
			return l.eofState
		}
		if r == '\n' {
			l.emit(CommentToken)
			return l.initState
		}
		l.accept()
	}
}
