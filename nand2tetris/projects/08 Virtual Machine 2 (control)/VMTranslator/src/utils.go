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
func CreateOutputFile(outputPath string) (*os.File, func(), error) {
	outputSink, err := os.Open(outputPath)
	if err != nil {
		return nil, nil, err
	}

	defer outputSink.Close()
	sinkInfo, err := outputSink.Stat()
	if err != nil {
		return nil, nil, err
	}

	baseName := filepath.Base(outputPath)
	dirName := filepath.Dir(outputPath)

	var hackFilePath string
	if sinkInfo.IsDir() {
		hackFilePath = filepath.Join(outputPath, baseName+".asm")
	} else {
		baseName = strings.Split(baseName, ".")[0]
		hackFilePath = filepath.Join(dirName, baseName+".asm")
	}
	hackFile, err := os.Create(hackFilePath)
	if err != nil {
		return nil, nil, err
	}
	return hackFile, func() {
		hackFile.Close()
	}, nil
}
