package main

import (
	"github.com/tie/x/qlex"
	"log"
	"os"
)

func expand(s string) string {
	return os.Expand(s, func(key string) string {
		return ""
	})
}

func main() {
	line := []qlex.Token{}
	for tok := range qlex.Lex(os.Stdin, initState) {
		switch tok.Typ {
		case SepToken:
			log.Println(line)
			line = []qlex.Token{}
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
