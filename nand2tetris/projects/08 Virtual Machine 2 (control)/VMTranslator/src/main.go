package main

import (
	"log"
	"os"

	"github.com/fatih/color"
)

const doc = `expected way to run VMTranslator is: ./VMTranslator <path to .vm file>`

func main() {
	if len(os.Args) != 2 {
		log.Fatal(color.RedString(doc))
	}

	if err := virtualMachine(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
