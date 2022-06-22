package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type CodeWriter struct {
	outputFile *os.File
	closer     func()
	segmentMap map[string]string
	labelId    int
}

func (cr *CodeWriter) Initializer(filePath string) {
	asmFile, closer := CreateOutputFile(filePath)
	*cr = CodeWriter{
		outputFile: asmFile,
		closer:     closer,
		labelId:    0,
		segmentMap: map[string]string{
			"local":    "LCL",
			"argument": "ARG",
			"this":     "THIS",
			"that":     "THAT",
			"temp":     "TEMP",
		},
	}
}

func (cr *CodeWriter) WritePushPop(command int, segment string, index int) error {
	var code []string
	pointer := []string{"THIS", "THAT"}
	fileName := filepath.Base(cr.outputFile.Name())
	fileName = strings.Split(fileName, ".")[0]
	switch command {
	case C_PUSH:
		switch segment {
		case "pointer":
			code = []string{
				fmt.Sprintf("@%v", pointer[index]),
				"D=M",
				"@SP",
				"A=M",
				"M=D", // *SP = pointer[index]
				"@SP",
				"M=M+1", // SP++
			}

		case "argument", "local", "this", "that":
			code = []string{
				fmt.Sprintf("@%v", cr.segmentMap[segment]),
				"D=M",
				fmt.Sprintf("@%v", index),
				"A=D+A", // A = segment + index
				"D=M",   // D = segment[index]
				"@SP",
				"A=M",
				"M=D", // *SP = segment[index]
				"@SP",
				"M=M+1", // SP++
			}

		case "temp":
			code = []string{
				"@5",
				"D=A",
				fmt.Sprintf("@%v", index),
				"A=D+A", // 5 + index
				"D=M",   // D = temp[index]
				"@SP",
				"A=M",
				"M=D", // *SP = D = temp[index]
				"@SP",
				"M=M+1", // SP++
			}

		case "constant":
			code = []string{
				fmt.Sprintf("@%v", index),
				"D=A", // D = index
				"@SP",
				"A=M",
				"M=D", // *SP = index
				"@SP",
				"M=M+1", // SP++
			}

		case "static":
			code = []string{
				fmt.Sprintf("@%v.%v", fileName, index), // @Foo.3
				"D=M",
				"@SP",
				"A=M",
				"M=D", // *SP = static[i]
				"@SP",
				"M=M+1",
			}
		default:
			return fmt.Errorf("unknown segment %v encountered ", color.RedString(segment))
		}
	case C_POP:
		switch segment {
		case "pointer":
			code = []string{
				"@SP",
				"M=M-1",
				"A=M",
				"D=M",
				fmt.Sprintf("@%v", pointer[index]),
				"M=D", // pointer[index] = *SP
			}
		case "argument", "this", "that", "local":
			code = []string{
				fmt.Sprintf("@%v", cr.segmentMap[segment]),
				"D=M",
				fmt.Sprintf("@%v", index),
				"D=D+A", // D = segment + i
				"@R13",
				"M=D",
				"@SP",
				"M=M-1", // SP = SP - 1
				"A=M",
				"D=M", // D = *SP
				"@R13",
				"A=M",
				"M=D", // *R13 = D
			}

		case "temp":
			code = []string{
				"@5",
				"D=A",
				fmt.Sprintf("@%v", index),
				"D=D+A",
				"@R13",
				"M=D",
				"@SP",
				"M=M-1",
				"A=M",
				"D=M", // D = *SP
				"@R13",
				"A=M",
				"M=D", // temp[index] = *R13 = D = *SP
			}
		case "static":
			code = []string{
				"@SP",
				"M=M-1",
				"A=M",
				"D=M",
				fmt.Sprintf("@%v.%v", fileName, index),
				"M=D",
			}
		default:
			// this block is not supposed to be executed ever, if it did, it's an unknown segment error
			return fmt.Errorf("unknown segment %v encountered ", color.RedString(segment))
		}
	default:
		// this block is not supposed to be executed ever, if it did, it's an unknown segment error
		return fmt.Errorf("unknown segment %v encountered ", color.RedString(segment))
	}
	return cr.Write(code)
}

func (cr CodeWriter) Write(code []string) error {
	tab := "\t"
	newLine := "\n"
	for _, asm := range code {
		if asm = strings.TrimSpace(asm); (asm[0] != '(') && (!strings.Contains(asm, "//")) {
			asm = tab + asm
		}

		_, err := cr.outputFile.Write([]byte(asm + newLine))
		if err != nil {
			return fmt.Errorf(color.RedString(err.Error()))
		}
	}
	return nil
}

func (cr *CodeWriter) WriteArithmetic(command string) error {
	var code []string
	switch command {
	case "add":
		// pops two top stack values, adds them, pushes result back to stack
		code = []string{
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"D=M", // D = *SP
			"@R13",
			"M=D", // R13 = *SP = D, first stack top value
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"D=M", // D = *SP, second stack top value
			"@R13",
			"A=M",   // A as data register
			"D=D+A", // adding two top stack values
			"@SP",
			"A=M",
			"M=D", // *SP = result of addition
			"@SP",
			"M=M+1",
		}
	case "sub":
		// pops top two stack values, say y and x, does x - y and pushes result back onto stack
		code = []string{
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"D=M", // D = *SP
			"@R13",
			"M=D", // R13 = *SP = D, first stack top value
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"D=M", // D = *SP, second stack top value
			"@R13",
			"A=M",   // A as data register
			"D=D-A", // adding two top stack values
			"@SP",
			"A=M",
			"M=D", // *SP = result of addition
			"@SP",
			"M=M+1",
		}
	case "neg":
		// pops top stack value, does arithmetic negation on it, pushes result back onto the stack
		code = []string{
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"M=-M", // *SP = -(*SP)
			"@SP",
			"M=M+1", // SP++
		}
	case "eq", "gt", "lt":
		code = cr.codeForGtLtEq(command)
	case "and", "or":
		code = cr.codeForAndOr(command)
	case "not":
		code = []string{
			"@SP",
			"M=M-1", // SP = SP - 1
			"A=M",
			"M=!M", // *SP = !(*SP)
			"@SP",
			"M=M+1", // SP++
		}
	}
	return cr.Write(code)
}

func (cr CodeWriter) codeForAndOr(command string) []string {
	mp := map[string]string{
		"and": "&",
		"or":  "|",
	}

	code := []string{
		"@SP",
		"M=M-1", // SP = SP - 1
		"A=M",
		"D=M",
		"@R13",
		"M=D", // R13 = *SP
		"@SP",
		"M=M-1", // SP = Sp - 1
		"A=M",
		"D=M",
		"@R13",
		fmt.Sprintf("D=D%vM", mp[command]), // D = x command y
		"@SP",
		"A=M",
		"M=D",
		"@SP",
		"M=M+1",
	}
	return code
}

func (cr *CodeWriter) codeForGtLtEq(command string) []string {
	mp := map[string]string{
		"eq": "JEQ",
		"gt": "JGT",
		"lt": "JLT",
	}
	code := []string{
		"@SP",
		"M=M-1", // SP = SP - 1
		"A=M",
		"D=M", // D = *SP, first top stack value, call it y
		"@R13",
		"M=D", // R13 = *SP = D, R13 = y
		"@SP",
		"M=M-1", // SP = SP - 1
		"A=M",
		"D=M", // D = *SP, call it x
		"@R13",
		"D=D-M", //  x - y
		fmt.Sprintf("@TRUE%v", cr.labelId),
		fmt.Sprintf("D;%v", mp[command]), // jump if (x - y commnad 0), command could be '<', '>', '=='
		"@0",                             // false result
		"D=A",
		fmt.Sprintf("@PUSH%v", cr.labelId),
		"0;JMP",
		fmt.Sprintf("(TRUE%v)", cr.labelId),
		"D=-1",                              // true result
		fmt.Sprintf("(PUSH%v)", cr.labelId), // pushed result onto stack
		"@SP",
		"A=M",
		"M=D", // *SP = (x command y)
		"@SP",
		"M=M+1", // SP++
	}
	cr.labelId++
	return code
}

func (cr CodeWriter) Close() {
	cr.closer()
}
