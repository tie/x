package token

import (
	"fmt"
)

type Position struct {
	Offset, Line, Column int
}

func (p Position) String() string {
	return fmt.Sprintf(
		"%d:%d(%+d)",
		p.Line+1, p.Column+1,
		p.Offset,
	)
}
