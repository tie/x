package main

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"
)

func main() {
	r := bufio.NewReaderSize(os.Stdin, utf8.UTFMax)
	unit, err := Parse(r, Lexicon)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(unit)
}
