package main

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"
)

func main() {
	l := NewLexer(bufio.NewReaderSize(os.Stdin, utf8.UTFMax))
	go l.Run()

	line := []Token{}
	for tok := range l.tokens {
		switch tok.Typ {
		case SepToken:
			log.Println(line)
			line = []Token{}
		case SpaceToken:
			continue
		case TextToken:
			// ignore escaped sep
			if tok.Val == "\\\n" {
				continue
			}
			line = append(line, tok)
		case CommentToken:
			continue
		}
	}
}
