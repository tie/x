package testingh

import (
	"io"
	"testing"
)

// TestReader calls testing.Error() if io.Reader.ReadRune() that returned an error is followed by another io.Reader.ReadRune() call.
type TestReader struct {
	Reader io.RuneReader
	Test *testing.T
	Err error
}

func NewTestReader(t *testing.T, r io.RuneReader) *TestReader {
	return &TestReader{r, t, nil}
}

func (r *TestReader) ReadRune() (c rune, size int, err error) {
	r.Test.Helper()
	if r.Err != nil {
		r.Test.Errorf("Read called after %q error", r.Err)
	}
	c, size, err = r.Reader.ReadRune()
	if err != nil {
		if r.Err != nil {
			r.Test.Errorf(
				"Read returned different %q error after %q",
				err, r.Err,
			)
		}
		r.Err = err
	} else {
		if r.Err != nil {
			r.Test.Errorf("Read succeeded after %q error", r.Err)
			r.Err = nil
		}
	}
	return
}
