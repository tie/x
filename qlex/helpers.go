package qlex

import (
	"strings"
	"unicode/utf8"
)

func (l *Lexer) emitToken(typ TokenType, value string) {
	l.tokens <- Token{typ, value, l.startPos, l.endPos}
	l.startPos = l.endPos
	l.buffer.Reset()
}

// Emit sends the currently read token on the Lexer.Tokens channel.
func (l *Lexer) Emit(typ TokenType) {
	l.emitToken(typ, l.buffer.String())
}

// TokenStart returns where the first accepted but not yet emitted byte or rune
// is located in the input source.
func (l Lexer) TokenStart() Position {
	return l.startPos
}

// TokenStop returns where the last accepted but not yet emitted byte or rune
// is located in the input source.
func (l Lexer) TokenStop() Position {
	return l.endPos
}

func (l *Lexer) discardRune(r rune, nbBytes int) {
	l.endPos.Offset += nbBytes
	switch r {
	case '\n':
		l.endPos.Column = 0
		l.endPos.Line++
	default:
		l.endPos.Column++
	}
}

func (l *Lexer) recordRune(r rune, width int) {
	l.buffer.WriteRune(r)
	l.discardRune(r, width)
}

func (l *Lexer) discardByte(b byte) {
	l.endPos.Offset++
	switch b {
	case '\n':
		l.endPos.Column = 0
		l.endPos.Line++
	default:
		l.endPos.Column++
	}
}

func (l *Lexer) recordByte(b byte) {
	l.buffer.WriteByte(b)
	l.discardByte(b)
}

// NextRune returns the next rune in the input source and advances.
// If you don't want to advance, use PeekRune.
// The resulting rune is wrapped in a MaybeRune, which allows to check whether
// you reached the end of the input source instead of reading a rune.
func (l *Lexer) NextRune() MaybeRune {
	r, width, err := l.reader.ReadRune()
	if err != nil {
		return eofRune
	}
	l.recordRune(r, width)
	return MaybeRune(r)
}

// NextByte returns the next byte in the input source and advances.
// If you don't want to advance, use PeekByte.
// The resulting byte is wrapped in a MaybeByte, which allows to check whether
// you reached the end of the input source instead of reading a byte.
func (l *Lexer) NextByte() MaybeByte {
	b, err := l.reader.ReadByte()
	if err != nil {
		return MaybeByte{eof: true}
	}
	l.recordByte(b)
	return MaybeByte{value: b}
}

// PeekRune returns the next rune in the input source without advancing.
// If you want to advance too, use NextRune.
// The resulting rune is wrapped in a MaybeRune, which allows to check whether
// you reached the end of the input source instead of reading a rune.
func (l *Lexer) PeekRune() MaybeRune {
	bytes, err := l.reader.Peek(1)
	if err != nil {
		return eofRune
	}

	if len(bytes) > 0 {
		if !utf8.RuneStart(bytes[0]) {
			return eofRune
		}
	}

	if !utf8.FullRune(bytes) {
		for w := 2; w <= utf8.UTFMax; w++ {
			bytes, err = l.reader.Peek(w)
			if err != nil {
				return eofRune
			}
			if !utf8.FullRune(bytes) {
				continue
			}
		}
	}

	r, sz := utf8.DecodeRune(bytes)
	if r == utf8.RuneError && sz < 2 {
		return eofRune
	}
	return MaybeRune(r)
}

// PeekByte returns the next byte in the input source without advancing.
// If you want to advance too, use NextByte.
// The resulting byte is wrapped in a MaybeByte, which allows to check whether
// you reached the end of the input source instead of reading a byte.
func (l *Lexer) PeekByte() MaybeByte {
	bytes, err := l.reader.Peek(1)
	if err != nil {
		return MaybeByte{eof: true}
	}
	return MaybeByte{value: bytes[0]}
}

// AcceptRuneAmong consumes the next rune if it's from the valid set and the
// end of the input source is not reached.
// It returns true if a rune was consumed, false otherwise.
func (l *Lexer) AcceptRuneAmong(valid string) bool {
	m := l.PeekRune()
	if m.IsEOF() {
		return false
	}
	r := m.Rune()
	if strings.ContainsRune(valid, r) {
		l.NextRune()
		return true
	}
	return false
}

// AcceptByteAmong consumes the next byte if it's from the valid set and the end
// of the input source is not reached.
// It returns true if a byte was consumed, false otherwise.
func (l *Lexer) AcceptByteAmong(valid string) bool {
	m := l.PeekByte()
	if m.IsEOF() {
		return false
	}
	b := m.Byte()
	if strings.IndexByte(valid, b) >= 0 {
		l.NextByte()
		return true
	}
	return false
}

// AcceptRuneIf is a generalization of AcceptRuneAmong: instead of checking
// whether the rune is contained in a provided string, it checks whether the
// provided predicate function returns true for that rune.
func (l *Lexer) AcceptRuneIf(f func(rune) bool) bool {
	m := l.PeekRune()
	if m.IsEOF() {
		return false
	}
	r := m.Rune()
	if f(r) {
		l.NextRune()
		return true
	}
	return false
}

// AcceptByteIf is a generalization of AcceptByteAmong: instead of checking
// whether the byte is contained in the provided string, it checks whether the
// provided predicate function returns true for that byte.
func (l *Lexer) AcceptByteIf(f func(byte) bool) bool {
	m := l.PeekByte()
	if m.IsEOF() {
		return false
	}
	b := m.Byte()
	if f(b) {
		l.NextByte()
		return true
	}
	return false
}

// AcceptRunesAmong consumes a sequence of consecutive runes as long as they are
// from the valid set and the end of the input source is not reached.
// It returns true if at least one rune was accepted, false otherwise.
func (l *Lexer) AcceptRunesAmong(valid string) bool {
	atLeastOne := false
	for l.AcceptRuneAmong(valid) {
		atLeastOne = true
	}
	return atLeastOne
}

// AcceptBytesAmong consumes a sequence of consecutive bytes as long as they are
// from the valid set and the end of the input source is not reached.
// It returns true if at least one byte was accepted, false otherwise.
func (l *Lexer) AcceptBytesAmong(valid string) bool {
	atLeastOne := false
	for l.AcceptByteAmong(valid) {
		atLeastOne = true
	}
	return atLeastOne
}

// AcceptRunesIf is a generalization of AcceptRunesAmong that uses AcceptRuneIf
// instead of AcceptRuneAmong: successive runes are consumed as long as they
// satisfy the provided predicate function and the end of the input source is
// not reached.
// It returns true if at least one rune was accepted, false otherwise.
func (l *Lexer) AcceptRunesIf(f func(rune) bool) bool {
	atLeastOne := false
	for l.AcceptRuneIf(f) {
		atLeastOne = true
	}
	return atLeastOne
}

// AcceptBytesIf is a generalization of AcceptBytesAmong that uses AcceptByteIf
// instead of AcceptByteAmong: successive bytes are consumed as long as they
// satisfy the provided predicate function and the end of the input source is
// not reached.
// It returns true if at least one byte was accepted, false otherwise.
func (l *Lexer) AcceptBytesIf(f func(byte) bool) bool {
	atLeastOne := false
	for l.AcceptByteIf(f) {
		atLeastOne = true
	}
	return atLeastOne
}
