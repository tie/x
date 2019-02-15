package main

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"
)

func lex(l *Lexer) <-chan Token {
	tokens := make(chan Token)
	go func() {
		defer close(tokens)
		for {
			tok, err := l.NextToken()
			if err != nil {
				break
			}
			tokens <- tok
		}
	}()
	return tokens
}

func main() {
	r := bufio.NewReaderSize(os.Stdin, utf8.UTFMax)
	l := NewLexer(r)

	line := []Token{}
	for tok := range lex(l) {
		switch tok.Typ {
		case SepToken:
			log.Println(line)
			line = []Token{}
		case SpaceToken:
			continue
		case TextToken:
			// line folding
			if tok.Val == "\\\n" {
				continue
			}
			line = append(line, tok)
		case CommentToken:
			continue
		}
	}
}
