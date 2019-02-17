package main

import (
	"github.com/tie/x/testingh"
	"testing"
)

func TestPos(t *testing.T) {
	panicCases := []string{
		"",
		"abc",
		"0:0",
		"1:1",
		"0:1(0)",
		"1:0(0)",
		"1:1(-1)",
	}
	for _, s := range panicCases {
		var p Position
		f := func() {
			p = Pos(s)
		}
		if testingh.DoesNotPanic(f) {
			t.Errorf("expected panic from Pos(%q), but got %s position", s, p)
		}
	}
	validCases := map[string]Position{
		"1:1(-0)": {0, 0, 0},
		"1:1(+0)": {0, 0, 0},
		"1:2(+1)": {1, 0, 1},
		"2:1(+1)": {1, 1, 0},
	}
	for s, p := range validCases {
		np := Pos(s)
		if np != p {
			t.Errorf("expected %s position, got %s position", p, np)
		}
	}
}
