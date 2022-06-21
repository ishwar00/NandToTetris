package main

import (
	"fmt"

	"github.com/fatih/color"
)

/// prepares error to report
func PrepError(instrInfo InstructionInfo, err error) error {
	return fmt.Errorf(`
Error on line:%v | %v
  	           ^^^^^^^ %v`, instrInfo.OnLine, color.RedString(instrInfo.Instruction), color.RedString(err.Error()))
}
