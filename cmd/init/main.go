package main

import (
	"bufio"
	"log"
	"os"
	"unicode/utf8"

	"github.com/tie/x/config"
)

func main() {
	r := bufio.NewReaderSize(os.Stdin, utf8.UTFMax)
	unit, err := config.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(unit)
}
