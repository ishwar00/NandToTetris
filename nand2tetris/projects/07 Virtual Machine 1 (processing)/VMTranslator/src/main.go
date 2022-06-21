package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf(color.RedString(`
expected way to run VMTranslator is: ./VMTranslator <path to vm file>`,
		)))
	}

	vmFilPath := os.Args[1]
	var parser Parser
	var codeWriter CodeWriter
	parser.Initializer(vmFilPath)
	codeWriter.Initializer(vmFilPath)
	for parser.HasMoreLines() {
		parser.Advance()
		codeWriter.Write([]string{"// " + parser.GetInstrInfo().Instruction})
		switch commandType := parser.CommandType(); commandType {
		case C_PUSH, C_POP:
			segment := parser.Arg1()
			index, err := parser.Arg2()
			if err != nil {
				panic(err)
			}
			err = codeWriter.WritePushPop(commandType, segment, index)
			if err != nil {
				panic(PrepError(parser.GetInstrInfo(), fmt.Errorf("failed while processing")))
			}
		case C_ARITHMETIC:
			command := parser.Arg1()
			err := codeWriter.WriteArithmetic(command)
			if err != nil {
				panic(PrepError(parser.GetInstrInfo(), fmt.Errorf("failed while processing")))
			}
		default:
			panic(color.RedString("missing implementation, VMTranslator is only implemented for C_PUSH, C_POP and C_ARITHMETIC type commands"))
		}
	}

	codeWriter.Write([]string{ // an infinite loop
		"(END)",
		"@END",
		"0;JMP",
	})
}
