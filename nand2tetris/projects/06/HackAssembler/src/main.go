package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const doc = `
 HackAssembler example/path/prog1.asm example1/path1/prog2.asm ...
 
 When given to your assembler as a command line argument, one or more
 progi.asm file containing a Hack assembly language program, it will be
 translated into the correct Hack binary code and stored in a file named
 Progi.hack, located in the same folder as the source file \n
 (if a file by this name exists, it is overwritten).
 `

func main() {
	// filePaths
	// 	"../../add/Add.asm"
	// 	"../../max/Max.asm"
	// 	"../../max/MaxL.asm"
	// 	"../../pong/Pong.asm"
	// 	"../../pong/PongL.asm"
	// 	"../../rect/Rect.asm"
	// 	"../../rect/RectL.asm"

	if len(os.Args) == 1 {
		err := fmt.Errorf(color.RedString("program need arguments as .asm file paths, run with flag --help"))
		panic(err)
	}

	if len(os.Args) == 2 && strings.Trim(os.Args[1], " ") == "--help" {
		fmt.Println(doc)
		return
	}

	for _, filePath := range os.Args[1:] {
		err := Assembler(filePath)
		if err != nil {
			color.Red("terminating assembler...")
			panic(err)
		}
		color.Green("finished assembling...")
		fmt.Println("")
	}

}
