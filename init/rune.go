package main

import (
	"unicode"
)

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isText(r rune) bool {
	return !isSpace(r) && r != '#'
}
