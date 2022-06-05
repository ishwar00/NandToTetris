package main

import (
	"bufio"
	"os"
	"strings"
)

const (
	A_INSTRUCTION = 0
	C_INSTRUCTION = 1
	L_INSTRUCTION = 2
)

type InstructionInfo struct {
	Instr  string // 16 bits
	AtLine int
	Type   int // A, L, C INSTRUCTION type
}

type Parser struct {
	instrList  []InstructionInfo
	nextInstr  int
	totalInstr int
}

func (p *Parser) Initialize(fileName string) {
	*p = Parser{
		instrList:  []InstructionInfo{},
		nextInstr:  -1,
		totalInstr: 0,
	}

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		// removing comment, if any
		if at := strings.Index(line, "//"); at != -1 {
			line = line[:at]
		}
		line = strings.Trim(line, " ")
		if len(line) > 0 {
			instrInfo := InstructionInfo{
				Instr:  line,
				AtLine: lineNumber,
				Type:   p.instrType(line),
			}
			p.instrList = append(p.instrList, instrInfo)
		}
	}
	p.totalInstr = len(p.instrList)
}

func (p Parser) HasMoreLines() bool {
	return p.nextInstr+1 < p.totalInstr
}

func (p *Parser) Advance() {
	p.nextInstr += 1
}

func (p Parser) InstructionType() int {
	return p.instrList[p.nextInstr].Type
}

func (p Parser) instrType(instruction string) int {
	switch {
	case instruction[0] == '(':
		return L_INSTRUCTION
	case instruction[0] == '@':
		return A_INSTRUCTION
	case strings.Contains(instruction, "=") || strings.Contains(instruction, ";"):
		return C_INSTRUCTION
	default:
		return -1
	}
}

func (p Parser) Symbol() string {
	switch p.InstructionType() {
	case L_INSTRUCTION:
		instruction := p.instrList[p.nextInstr].Instr
		return strings.Trim(instruction, "()") // (symbol) -> symbol
	case A_INSTRUCTION:
		instruction := p.instrList[p.nextInstr].Instr
		return instruction[1:] // @vv...v -> vv...v
	default:
		return "" // does not contain any symbol
	}
}

func (p Parser) Dest() string {
	instruction := p.instrList[p.nextInstr].Instr
	if strings.Contains(instruction, "=") {
		return strings.Split(instruction, "=")[0]
	}
	return "" // destinition is not specified
}

func (p Parser) Comp() string {
	instruction := p.instrList[p.nextInstr].Instr
	if strings.Contains(instruction, "=") {
		return strings.Split(instruction, "=")[1] // dest=comp
	}
	return strings.Split(instruction, ";")[0] // comp;jump
}

func (p Parser) Jump() string {
	instruction := p.instrList[p.nextInstr].Instr
	if strings.Contains(instruction, ";") {
		return strings.Split(instruction, ";")[1] // comp;jump
	}
	return "" // jump is not specified
}

func (p Parser) GetInstrInfo() InstructionInfo {
	return p.instrList[p.nextInstr]
}

func (p Parser) GetTotalInstr() int {
	return p.totalInstr
}
