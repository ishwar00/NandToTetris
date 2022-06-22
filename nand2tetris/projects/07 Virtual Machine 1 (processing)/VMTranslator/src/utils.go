package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/// prepares error to report
func PrepError(instrInfo InstructionInfo, err error) error {
	return fmt.Errorf(`
Error on line:%v | %v
  	           ^^^^^^^ %v`, instrInfo.OnLine, instrInfo.Instruction, err.Error())
}

// creates file with the same name and directory as filePath with file extension asm,
// eg: if filePath is example/path/Prog.vm, then it creates a file example/path/Prog.asm
// it will overwrite if already such file exists.
// returns function which must be called to close the created file.
func CreateOutputFile(filePath string) (*os.File, func()) {
	directory, fileName := filepath.Split(filePath)
	hackFilePath := filepath.Join(directory, strings.Split(fileName, ".")[0]+".asm")
	hackFile, err := os.Create(hackFilePath)
	if err != nil {
		panic(err)
	}
	return hackFile, func() {
		hackFile.Close()
	}
}
