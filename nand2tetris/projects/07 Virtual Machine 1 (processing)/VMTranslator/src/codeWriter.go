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

func (cr *CodeWriter) Write(code []string) error {
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
	case "add", "sub", "and", "or":
		// pops two top stack values, {add, sub, and, or} them, pushes result back to stack
		code = cr.codeForAddSubAndOr(command)
	case "neg", "not":
		// pops top stack value, does arithmetic or logical negation on it, pushes result back onto the stack
		mp := map[string]string{
			"neg": "-",
			"not": "!",
		}
		code = []string{
			"@SP",
			"A=M",
			"A=A-1",
			fmt.Sprintf("M=%vM", mp[command]), // *SP = -(*SP) or !(*SP)
		}
	case "eq", "gt", "lt":
		code = cr.codeForGtLtEq(command)
	}
	return cr.Write(code)
}

func (cr CodeWriter) codeForAddSubAndOr(command string) []string {
	mp := map[string]string{
		"add": "+",
		"sub": "-",
		"or":  "|",
		"and": "&",
	}

	// | b | a |  |   <- stack
	// 			SP    <- stack pointer
	//  b op a
	code := []string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M", // D = a
		"A=A-1",
		// M = b
	}

	if command == "sub" {
		code = append(code, "M=M-D") // M = b - a
	} else {
		code = append(code, cr.swapDM()...)                     // injecting swaping code
		code = append(code, fmt.Sprintf("M=D%vM", mp[command])) // M = b op a
	}
	return code
}

// this will swap values of registers D and M
// say D has a, M has b
func (cr CodeWriter) swapDM() []string {
	code := []string{
		"M=D+M", // M = a + b
		"D=M-D", // D = a + b - a, => D = b
		"M=M-D", // M = a + b - b, => M = a
	}
	return code
}

func (cr *CodeWriter) codeForGtLtEq(command string) []string {
	mp := map[string]string{
		"eq": "JEQ",
		"gt": "JGT",
		"lt": "JLT",
	}
	// | b | a |  |   <- stack
	// 			SP    <- stack pointer
	//  b op a
	code := []string{
		"@SP",
		"M=M-1", // SP = SP - 1
		"A=M",
		"D=M", // D = *SP, D = a
		"A=A-1",
		"D=M-D", //  b - a
		fmt.Sprintf("@TRUE_%v", cr.labelId),
		fmt.Sprintf("D;%v", mp[command]), // jump if (b - a op 0), op belongs to { '<', '>', '==' }
		"D=0",                            // false
		fmt.Sprintf("@DONE_%v", cr.labelId),
		"0;JMP",
		fmt.Sprintf("(TRUE_%v)", cr.labelId),
		"D=-1", // true
		fmt.Sprintf("(DONE_%v)", cr.labelId),
		"@SP",
		"A=M",
		"A=A-1",
		"M=D",
	}
	cr.labelId++
	return code
}

func (cr CodeWriter) Close() {
	cr.closer()
}
