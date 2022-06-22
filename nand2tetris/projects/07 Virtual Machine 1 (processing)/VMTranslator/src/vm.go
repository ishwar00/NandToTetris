package main

import (
	"fmt"

	"github.com/fatih/color"
)

func virtualMachine(vmFilePath string) error {
	var parser Parser
	var codeWriter CodeWriter
	parser.Initializer(vmFilePath)
	codeWriter.Initializer(vmFilePath)
	defer codeWriter.closer()
	comment := "// "
	for parser.HasMoreLines() {
		parser.Advance()
		codeWriter.Write([]string{comment + parser.GetInstrInfo().Instruction})
		switch commandType := parser.CommandType(); commandType {
		case C_PUSH, C_POP:
			segment := parser.Arg1()
			index, err := parser.Arg2()
			if err != nil { // failed parse index
				return err
			}
			err = codeWriter.WritePushPop(commandType, segment, index)
			if err != nil { // failed while either writing into file or parsing
				return PrepError(parser.GetInstrInfo(), err)
			}
		case C_ARITHMETIC:
			command := parser.Arg1()
			err := codeWriter.WriteArithmetic(command)
			if err != nil {
				return PrepError(parser.GetInstrInfo(), err)
			}
		default:
			errMsg := fmt.Errorf(color.RedString("missing implementation, VMTranslator is only implemented for C_PUSH, C_POP and C_ARITHMETIC type commands"))
			return PrepError(parser.GetInstrInfo(), errMsg)
		}
	}

	codeWriter.Write([]string{ // an infinite loop
		"(END)",
		"@END",
		"0;JMP",
	})
	return nil
}
