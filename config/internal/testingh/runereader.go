package testingh

import (
	"errors"
)

var (
	ErrSequenceEnd = errors.New("cannot advance past the end of the sequence")
)

// ReadRune is a return value captured from io.RuneReader.ReadRune().
type ReadRune struct {
	Rune rune
	Size int
	Error error
}

// RuneReader is io.RuneReader that reproduces the sequence of ReadRune return values.
type RuneReader struct {
	// Seq is a sequence of values to be returned with successive calls to ReadRune.
	Seq []ReadRune
	// Pos is the position in sequnce.
	Pos int
}

func NewRuneReader(vals []ReadRune) *RuneReader {
	return &RuneReader{
		Seq: vals,
	}
}

// ReadRune returns current value in the sequence and advances to the next one.
// Error will be set to ErrSequenceEnd on attempt to advance past the end of the the sequence.
func (r *RuneReader) ReadRune() (rune, int, error) {
	var v ReadRune
	if r.Pos < len(r.Seq) {
		v = r.Seq[r.Pos]
		r.Pos++
	} else {
		v.Error = ErrSequenceEnd
	}
	return v.Rune, v.Size, v.Error
}
