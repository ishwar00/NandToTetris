package main

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

func Assembler(filePath string) error {
	color.Green("assembling...")

	var code Code
	var parser Parser
	var st SymbolTable
	code.Initialize()
	st.Initialize(filePath)
	parser.Initialize(filePath)
	memoryAllocator := 16
	binaryInstructions := []string{}

	err := st.DoPass1()
	if err != nil {
		return err
	}

	for parser.HasMoreLines() {
		parser.Advance()
		switch parser.InstructionType() {
		case L_INSTRUCTION:
			// Lable declaration (xxx) produces no code
			continue

		case A_INSTRUCTION:
			symbol := parser.Symbol()

			if IsSymbol(symbol) {
				if !st.Contains(symbol) {
					st.AddEntry(symbol, memoryAllocator)
					memoryAllocator++
				}

				address := st.GetAddress(symbol)
				binaryInstructions = append(binaryInstructions, BinaryRep(int64(address)))
			} else {
				v, err := strconv.ParseInt(symbol, 10, 32)
				if err != nil {
					errMsg := fmt.Errorf("could not parse %v into integer\n %w", color.RedString(symbol), err)
					return PrepError(parser.GetInstrInfo(), errMsg)
				}
				binaryInstructions = append(binaryInstructions, BinaryRep(int64(v)))
			}

		case C_INSTRUCTION:
			dest := parser.Dest()
			comp := parser.Comp()
			jump := parser.Jump()
			binaryCode := "111" + code.Comp(comp) + code.Dest(dest) + code.Jump(jump)
			if len(binaryCode) != 16 {
				errMsg := fmt.Errorf("unknown instruction: %v", parser.GetInstrInfo().Instr)
				return PrepError(parser.GetInstrInfo(), errMsg)
			}
			binaryInstructions = append(binaryInstructions, binaryCode)
		default:
			return PrepError(parser.GetInstrInfo(), fmt.Errorf("alien instruction: failed to classify the instruction"))
		}
	}

	WriteToFile(binaryInstructions, filePath)
	TI := strconv.Itoa(parser.GetTotalInstr())
	fmt.Println(color.GreenString("processed"), color.YellowString(TI), color.GreenString("instructions"))
	return nil
}
