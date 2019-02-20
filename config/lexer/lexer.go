package lexer

import (
	"io"
	"strings"
	"unicode"

	"github.com/tie/x/config/token"
)

type Lexer struct {
	reader io.RuneReader
	buffer strings.Builder
	startPos, endPos token.Position
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

func (l *Lexer) NextToken() (token.Token, error) {
	r, err := l.peek()
	if err != nil {
		return token.Token{}, err
	}
	tok, err := l.nextState(r)
	if err == io.EOF {
		err = nil
	}
	return tok, err
}

func (l *Lexer) nextState(r rune) (token.Token, error) {
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

func (l *Lexer) emit(typ token.TokenType) token.Token {
	value := l.buffer.String()
	tok := token.Token{
		Typ: typ,
		Val: value,
		Pos: l.startPos,
		End: l.endPos,
	}
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

func (l *Lexer) sepState() (token.Token, error) {
	// assume separator (i.e. end of line)
	l.accept()
	return l.emit(token.SepToken), nil
}

func (l *Lexer) spacesState() (token.Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(token.SpaceToken), err
		}
		if !unicode.IsSpace(r) || r == '\n' {
			return l.emit(token.SpaceToken), nil
		}
		l.accept()
	}
}

func (l *Lexer) textState() (token.Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(token.TextToken), err
		}
		if unicode.IsSpace(r) || r == '#' {
			return l.emit(token.TextToken), nil
		}
		switch r {
		case '\\':
			// escape character
			l.accept()
			r, err := l.read()
			if err != nil {
				return l.emit(token.TextToken), err
			}
			// and terminate at the end of line
			if r == '\n' {
				return l.emit(token.TextToken), nil
			}
			continue
		case '"':
			l.accept()
			err := quoteText(l, r)
			if err != nil {
				return l.emit(token.TextToken), err
			}
			continue
		}
		l.accept()
	}
}

func quoteText(l *Lexer, quote rune) error {
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

func (l *Lexer) commentState() (token.Token, error) {
	for {
		r, err := l.peek()
		if err != nil {
			return l.emit(token.CommentToken), err
		}
		if r == '\n' {
			return l.emit(token.CommentToken), nil
		}
		l.accept()
	}
}
