package main

import (
	"fmt"

	"github.com/fatih/color"
)

type SymbolTable struct {
	ST       map[string]int
	fileName string
}

func (st *SymbolTable) Initialize(fileName string) {
	*st = SymbolTable{
		ST: map[string]int{
			"R0":     0,
			"R1":     1,
			"R2":     2,
			"R3":     3,
			"R4":     4,
			"R5":     5,
			"R6":     6,
			"R7":     7,
			"R8":     8,
			"R9":     9,
			"R10":    10,
			"R11":    11,
			"R12":    12,
			"R13":    13,
			"R14":    14,
			"R15":    15,
			"SP":     0,
			"LCL":    1,
			"ARG":    2,
			"THIS":   3,
			"THAT":   4,
			"SCREEN": 16384,
			"KBD":    24576,
		},
		fileName: fileName,
	}
}

func (st *SymbolTable) AddEntry(symbol string, address int) {
	st.ST[symbol] = address
}

func (st SymbolTable) Contains(symbol string) bool {
	_, ok := st.ST[symbol]
	return ok
}

func (st SymbolTable) GetAddress(symbol string) int {
	return st.ST[symbol]
}

func (st *SymbolTable) DoPass1() error {
	var parser Parser
	parser.Initialize(st.fileName)

	address := -1
	for parser.HasMoreLines() {
		parser.Advance()
		if parser.InstructionType() == L_INSTRUCTION {
			if st.Contains(parser.Symbol()) {
				errMsg := fmt.Errorf(color.RedString("duplicate label %v found, labels must be unique"), parser.Symbol())
				return PrepError(parser.GetInstrInfo(), errMsg)
			}
			st.AddEntry(parser.Symbol(), address+1)
			continue
		}
		address++
	}
	return nil
}
