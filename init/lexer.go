package main

import (
	"github.com/tie/x/qlex"
	"unicode"
)

const (
	_ = iota
	SepToken
	SpaceToken
	TextToken
	CommentToken
)

func initState(l *qlex.Lexer) qlex.StateFunc {
	m := l.PeekRune()
	if m.IsEOF() {
		return eofState
	}
	r := m.Rune()
	switch r {
	case '#':
		return commentState
	case '\n':
		return sepState
	default:
		if isSpace(r) {
			return spacesState
		}
		return textState
	}
}

func eofState(l *qlex.Lexer) qlex.StateFunc {
	return nil
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isText(r rune) bool {
	return !isSpace(r) && r != '\n' && r != '#'
}

func sepState(l *qlex.Lexer) qlex.StateFunc {
	l.NextRune()
	l.Emit(SepToken)
	return initState
}

func spacesState(l *qlex.Lexer) (next qlex.StateFunc) {
	for {
		m := l.PeekRune()
		if m.IsEOF() {
			l.Emit(SpaceToken)
			return eofState
		}
		r := m.Rune()
		if !isSpace(r) {
			l.Emit(SpaceToken)
			return initState
		}
		l.NextRune()
	}
}

func textState(l *qlex.Lexer) qlex.StateFunc {
	for {
		m := l.PeekRune()
		if m.IsEOF() {
			l.Emit(TextToken)
			return eofState
		}
		r := m.Rune()
		switch r {
		case '\\':
			if !escapeText(l) {
				l.Emit(TextToken)
				return eofState
			}
		case '"':
			if !quoteText(l) {
				l.Emit(TextToken)
				return eofState
			}
		default:
			if !isText(r) {
				l.Emit(TextToken)
				return initState
			}
			l.NextRune()
		}
	}
}

func escapeText(l *qlex.Lexer) (ok bool) {
	l.NextRune()
	m := l.PeekRune()
	if m.IsEOF() {
		// escaping EOF, huh…
		return false
	}
	l.NextRune()
	return true
}

func quoteText(l *qlex.Lexer) (ok bool) {
	quote := l.NextRune().Rune()
	for {
		m := l.PeekRune()
		if m.IsEOF() {
			// unterminated quoted thing
			return false
		}
		r := m.Rune()
		switch r {
		case '\\':
			if !escapeText(l) {
				return false
			}
		case '\n':
			// terminate quoted token at end of line
			return true
		case quote:
			// consume closing quote
			l.NextRune()
			return true
		default:
			l.NextRune()
		}
	}
}

func commentState(l *qlex.Lexer) qlex.StateFunc {
	for {
		m := l.PeekRune()
		if m.IsEOF() {
			l.Emit(CommentToken)
			return eofState
		}
		r := m.Rune()
		if r == '\n' {
			l.Emit(CommentToken)
			return initState
		}
		l.NextRune()
	}
}
