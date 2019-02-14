package qlex

import (
	"fmt"
)

// Position represents a specific offset in the input stream.
type Position struct {
	Offset int
	Line   int
	Column int
}

func (p Position) String() string {
	return fmt.Sprintf("%v:%v (byte %v)", p.Line+1, p.Column+1, p.Offset+1)
}
