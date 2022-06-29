package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type CodeWriter struct {
	outputFile    *os.File
	closer        func()
	segmentMap    map[string]string
	labelId       int
	currentVMFile string
}

func (cr *CodeWriter) Initialize(filePath string) error {
	asmFile, closer, err := CreateOutputFile(filePath)
	if err != nil {
		return err
	}

	*cr = CodeWriter{
		outputFile:    asmFile,
		closer:        closer,
		labelId:       0,
		currentVMFile: "null",
		segmentMap: map[string]string{
			"local":    "LCL",
			"argument": "ARG",
			"this":     "THIS",
			"that":     "THAT",
			"temp":     "TEMP",
		},
	}
	return nil
}

func (cr *CodeWriter) InjectBootstrapCode() {
	// 	 	 0	 |	261	 |	SP
	// 	 	 1	 |	261	 |	LCL
	// 	 	 2	 |	256	 |	ARG
	// 	 	 3	 |	-1	 |	THIS
	// 	 	 4	 |	-1	 |	THAT
	//	    ...	 |	...	 |
	// 		256	 |  1234 | fake return address
	// 		257	 |	261	 |	LCL
	// 		258	 |	256	 |	ARG
	// 		259	 |	-1	 |	THIS	// just before calling Sys.init
	// 		260	 |	-1	 |	THAT
	//		261	 |	...  |  <-- SP, LCL
	//			    ...
	//		jump to Sys.init

	code := []string{
		// SP = 261
		"@261",
		"D=A",
		"@SP",
		"M=D",

		// push fake address
		"@1234",
		"D=A",
		"@256",
		"M=D",

		// LCL = 261
		"@261",
		"D=A",
		"@LCL",
		"M=D",

		// ARG = 256
		"@256",
		"D=A",
		"@ARG",
		"M=D",

		// THIS = -1
		"@THIS",
		"M=-1",

		// THAT = -1
		"@THAT",
		"M=-1",

		// Jumping to Sys.init
		"@Sys.init",
		"0;JMP",
	}
	cr.Write(code)
}

// implements command: label SOME_LABEL
func (cr CodeWriter) WriteLabel(label, functionName string) error {
	code := []string{
		fmt.Sprintf("(%v$%v)", functionName, label),
	}
	return cr.Write(code)
}

// implements command: goto SOME_LABEL
func (cr CodeWriter) WriteGoto(label, functionName string) error {
	code := []string{
		fmt.Sprintf("@%v$%v", functionName, label),
		"0;JMP",
	}
	return cr.Write(code)
}

// implements command: if-goto SOME_LABEL
func (cr CodeWriter) WriteIf(label, functionName string) error {
	code := []string{
		"@SP",
		"M=M-1",
		"A=M",
		"D=M",
		fmt.Sprintf("@%v$%v", functionName, label),
		"D;JNE",
	}

	return cr.Write(code)
}

// implements command: function functionName Vargs
func (cr CodeWriter) WriteFunction(functionName string, nVars int) error {
	code := []string{
		fmt.Sprintf("(%v)", functionName),
	}

	pushCode := []string{ // pushed zero on stack
		"@SP",
		"A=M",
		"M=0",
		"@SP",
		"M=M+1",
	}

	for i := 0; i < nVars; i++ {
		code = append(code, pushCode...) // initializing local variables to zero
	}
	return cr.Write(code)
}

// implements command: call functionName nArgs
func (cr CodeWriter) WriteCall(calleeFunction, currentFunction string, nArgs int, callCount int) error {
	//					|	    ... 		 |
	//					|	    ... 		 |
	// 		ARG -->		| 		arg0		 |		caller stack just before jumping to callee's code block.
	//					|		...			 |		caller has pushed arguments for callee
	//					|		argn		 |
	// 					|	return address	 |  <----
	// 					|		LCL 	     | 		|
	// 					|		ARG 	     |		|	VM implements this
	// 					|		THIS	     | 		|
	// 					|		THAT	     |	<----
	// 	SP, LCL -->		| 		...			 |

	pushPointer := func(pointer string) []string {
		return []string{
			fmt.Sprintf("@%v", pointer),
			"D=M",
			"@SP",
			"A=M",
			"M=D", // *SP = pointer
			"@SP",
			"M=M+1",
		}
	}

	var code []string
	returnAddress := fmt.Sprintf("%v$ret.%v", currentFunction, callCount)

	code = append(code,
		fmt.Sprintf("@%v", returnAddress),
		"D=A",
		"@SP",
		"A=M",
		"M=D", // *SP = pointer
		"@SP",
		"M=M+1",
	)
	code = append(code, pushPointer("LCL")...)
	code = append(code, pushPointer("ARG")...)
	code = append(code, pushPointer("THIS")...)
	code = append(code, pushPointer("THAT")...)
	code = append(code,
		// setting argument section for callee
		// ARG = SP - 5 - nArgs
		"@SP",
		"D=M",
		"@5",
		"D=D-A", // SP = SP - 5
		fmt.Sprintf("@%v", nArgs),
		"D=D-A", // SP = SP - nArgs
		"@ARG",
		"M=D",

		// set LCL = SP
		"@SP",
		"D=M",
		"@LCL",
		"M=D",

		// goto functionName
		fmt.Sprintf("@%v", calleeFunction),
		"0;JMP",

		// injecting label (returnAddress)
		fmt.Sprintf("(%v)", returnAddress),
	)
	return cr.Write(code)
}

// implements command: return
func (cr CodeWriter) WriteReturn() error {

	//						 	    ARG  --> |		arg0	   |
	//										 |		...		   |
	//										 |		arg1	   |
	//							  		---> | return address  |
	//							  		|	 | 		LCL 	   |
	//							  frame | 	 |		ARG 	   |  this is callee's stack,
	//							  		|	 | 		THIS	   |
	//							  		---> | 		THAT	   |
	//								LCL	-->	 |		...		   | <- frame(pointer)
	//										 |		...		   |
	//									     |	return value   |
	//						  		 SP -->  |		...        |
	//						 					global stack
	//
	//								*	save return address
	//								*	overwrite arg0 with return value
	//								*   restore pointer's of caller LCL, ARG, THIS and THAT
	//								*   jump to return address

	code := []string{
		// R13 = frame = LCL
		"@LCL",
		"D=M",
		"@R13", // frame = R13 are pointers
		"M=D",

		// R14 = returnAddress = *(R13 - 5)
		"@R13",
		"D=M",
		"@5",
		"A=D-A", // D = R13 - 5
		"D=M",
		"@R14",
		"M=D", // R14 = returnAddress

		// *ARG = pop()
		"@SP",
		"M=M-1", // SP = SP - 1
		"A=M",
		"D=M", // D = pop()
		"@ARG",
		"A=M",
		"M=D",

		// SP = ARG + 1
		"@ARG",
		"M=M+1",
		"D=M",
		"@SP",
		"M=D",

		// THAT = *(frame - 1)
		"@R13",
		"M=M-1", // frame = frame - 1
		"A=M",
		"D=M", // D = *(frame)
		"@THAT",
		"M=D",

		// THIS = *(frame - 2)
		"@R13",
		"M=M-1", // frame = frame - 1
		"A=M",
		"D=M", // D = *(frame)
		"@THIS",
		"M=D",

		// ARG = *(frame - 3)
		"@R13",
		"M=M-1", // frame = frame - 1
		"A=M",
		"D=M", // D = *(frame)
		"@ARG",
		"M=D",

		// LCL = *(frame - 4)
		"@R13",
		"M=M-1", // frame = frame - 1
		"A=M",
		"D=M", // D = *(frame)
		"@LCL",
		"M=D",

		// goto return return
		"@R14",
		"A=M",
		"0;JMP",
	}
	return cr.Write(code)
}

// should be used to set to current compiling .vm file
func (cr *CodeWriter) SetFilePath(filePath string) {
	cr.currentVMFile = filePath
}

func (cr *CodeWriter) WritePushPop(command int, segment string, index int) error {
	var code []string
	pointer := []string{"THIS", "THAT"}
	fileName := filepath.Base(cr.currentVMFile)
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

	// | b | a |  |   <- global stack
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
		code = append(code, fmt.Sprintf("M=D%vM", mp[command])) // M = b op a
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
