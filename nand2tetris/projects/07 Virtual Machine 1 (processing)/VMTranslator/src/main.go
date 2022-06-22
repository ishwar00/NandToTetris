package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf(color.RedString(`
expected way to run VMTranslator is: ./VMTranslator <path to .vm file>`,
		)))
	}

	if err := virtualMachine(os.Args[1]); err != nil {
		panic(err)
	}
}
