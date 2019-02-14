package qlex

import (
	"fmt"
)

// MaybeRune wraps either a rune or an error.
type MaybeRune rune

const eofRune = -(iota + 1)

// MaybeByte wraps either a byte or an error.
type MaybeByte struct {
	value byte
	eof   bool
}

// Rune returns the rune represented by the MaybeRune, or panics.
func (m MaybeRune) Rune() rune {
	if m < 0 {
		panic("no rune read: reached end of file")
	}
	return rune(m)
}

// Byte returns the byte represented by the MaybeByte, or panics.
func (m MaybeByte) Byte() byte {
	if m.eof {
		panic("no byte read: reached end of file")
	}
	return m.value
}

// IsEOF returns whether a MaybeRune represents the fact of having reached the
// end of an input source.
func (m MaybeRune) IsEOF() bool {
	return m < 0
}

// IsEOF returns whether a MaybeByte represents the fact of having reached the
// end of an input source.
func (m MaybeByte) IsEOF() bool {
	return m.eof
}

func (m MaybeRune) String() string {
	if m.IsEOF() {
		return "EOF"
	}
	return fmt.Sprintf("%q", m.Rune())
}

func (m MaybeByte) String() string {
	if m.IsEOF() {
		return "EOF"
	}
	return fmt.Sprintf("%q", m.Byte())
}
