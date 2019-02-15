package main

import (
	"bufio"
	"log"
	"os"
	"io"
	"unicode/utf8"
)

func main() {
	r := bufio.NewReaderSize(os.Stdin, utf8.UTFMax)
	p := NewParser(r)
	for {
		line, err := p.parseLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		log.Println(line)
	}
}
