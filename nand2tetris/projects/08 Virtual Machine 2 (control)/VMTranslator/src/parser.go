package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	C_ARITHMETIC = 0
	C_PUSH       = 1
	C_POP        = 2
	C_LABEL      = 3
	C_GOTO       = 4
	C_IF         = 5
	C_FUNCTION   = 6
	C_RETURN     = 7
	C_CALL       = 8
)

type InstructionInfo struct {
	Instruction string
	Type        int
	OnLine      int
}

type Parser struct {
	instrInfo       []InstructionInfo // instruction info
	CurrentFunction string            // name of function being parsed
	CallCount       int
	nextInstr       int // next instruction
}

/// takes input file, reads through the program line by line and removes
/// comment and trims white spaces around the instruction, and puts in a slice.
func (p *Parser) Initialize(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.IsDir() || filepath.Ext(filePath) != ".vm" {
		panic(fmt.Errorf(color.RedString("provided argument must be a file with .vm extension")))
	}

	*p = Parser{
		instrInfo:       []InstructionInfo{},
		nextInstr:       -1,
		CurrentFunction: "null",
		CallCount:       0,
	}

	scanner := bufio.NewScanner(file)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		instruction := scanner.Text()
		if at := strings.Index(instruction, "//"); at != -1 {
			instruction = instruction[:at]
		}
		instruction = strings.TrimSpace(instruction)
		if len(instruction) > 0 {
			instrInfo := InstructionInfo{
				Instruction: instruction,
				Type:        p.commandType(instruction),
				OnLine:      lineNumber,
			}
			p.instrInfo = append(p.instrInfo, instrInfo)
		}
	}
	return nil
}

func (p Parser) HasMoreLines() bool {
	return p.nextInstr+1 < len(p.instrInfo)
}

/// makes the next command as the current command
/// this method should be called only if HasMoreLines is true
/// Initially there is no current command.
func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.nextInstr += 1
		instruction := p.GetInstrInfo().Instruction
		switch p.commandType(instruction) {
		case C_FUNCTION:
			p.CallCount = 0
			p.CurrentFunction = strings.Fields(instruction)[1] // function name
		case C_CALL:
			p.CallCount += 1
		}
	}
}

/// Returns a constant representing the type of the current command
/// if the current command if arithmentic-logic command, returns C_ARITHMETIC
func (p Parser) CommandType() int {
	instruction := p.instrInfo[p.nextInstr].Instruction
	return p.commandType(instruction)
}

func (p Parser) commandType(instruction string) int {
	fields := strings.Fields(instruction)
	switch fields[0] {
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	case "call":
		return C_CALL
	case "if-goto":
		return C_IF
	case "goto":
		return C_GOTO
	case "label":
		return C_LABEL
	default:
		return -1
	}
}

func (p Parser) GetInstrInfo() InstructionInfo {
	return p.instrInfo[p.nextInstr]
}

func (p Parser) Arg1() string {
	fields := strings.Fields(p.GetInstrInfo().Instruction)
	switch p.CommandType() {
	case C_POP, C_PUSH:
		return fields[1] // segment
	case C_ARITHMETIC:
		return fields[0] // command itself
	default:
		return ""
	}
}

func (p Parser) Arg2() (int, error) {
	fields := strings.Fields(p.GetInstrInfo().Instruction)
	switch p.CommandType() {
	case C_PUSH, C_POP:
		v, err := strconv.ParseInt(fields[2], 10, 32)
		if err != nil {
			return -1, PrepError(p.GetInstrInfo(), err)
		}
		return int(v), nil
	default:
		return -1, PrepError(p.GetInstrInfo(), fmt.Errorf("Arg2 is only called on instructions of type C_PUSH and C_POP"))
	}
}
