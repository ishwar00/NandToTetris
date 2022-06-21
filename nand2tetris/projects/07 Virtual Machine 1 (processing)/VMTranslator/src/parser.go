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
	instrInfo []InstructionInfo // instruction info
	nextInstr int               // next instruction
}

/// takes input file, reads through the program line by line and removes
/// comment and trims white spaces around the instruction, and puts in a slice.
func (p *Parser) Initializer(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() || filepath.Ext(filePath) != ".vm" {
		panic(fmt.Errorf(color.RedString("provided argument must be a file with .vm extension")))
	}

	*p = Parser{
		instrInfo: []InstructionInfo{},
		nextInstr: -1,
	}

	scanner := bufio.NewScanner(file)
	for lineNumber := 0; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		if at := strings.Index(line, "//"); at != -1 {
			line = line[:at]
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			instrInfo := InstructionInfo{
				Instruction: line,
				Type:        commandType(line),
				OnLine:      lineNumber,
			}
			p.instrInfo = append(p.instrInfo, instrInfo)
		}
	}
}

func (p *Parser) HasMoreLines() bool {
	return p.nextInstr+1 < len(p.instrInfo)
}

/// makes the next command as the current command
/// this method should be called only if HasMoreLines is true
/// Initially there is no current command.
func (p *Parser) Advance() {
	if p.HasMoreLines() {
		p.nextInstr += 1
	}
}

/// Returns a constant representing the type of the current command
/// if the current command if arithmentic-logic command, returns C_ARITHMETIC
func (p Parser) CommandType() int {
	instruction := p.instrInfo[p.nextInstr].Instruction
	return commandType(instruction)
}

func commandType(instruction string) int {
	fields := strings.Fields(instruction)
	switch fields[0] {
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return C_ARITHMETIC
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
