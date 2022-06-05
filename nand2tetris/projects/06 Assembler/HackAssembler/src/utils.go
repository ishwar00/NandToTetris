package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

func PrepError(instrInfo InstructionInfo, err error) error {
	return fmt.Errorf(`
Error on line:%v | %v
  	           ^^^^^^^ %v`, instrInfo.AtLine, color.RedString(instrInfo.Instr), color.RedString(err.Error()))
}

func BinaryRep(value int64) string {
	binary := []byte{'0'}
	for i := 14; i >= 0; i = i - 1 {
		if value&(1<<i) == 0 {
			binary = append(binary, '0')
		} else {
			binary = append(binary, '1')
		}
	}
	return string(binary)
}

func WriteToFile(binaryInstructions []string, filePath string) {
	hackFile, Close := CreateOutputFile(filePath)
	defer Close()
	for i, instruction := range binaryInstructions {
		if len(instruction) != 16 {
			err := fmt.Errorf(color.RedString("%vth instruction length is not 16", i))
			panic(err)
		}
		_, err := hackFile.Write([]byte(instruction + "\n"))
		if err != nil {
			panic(err)
		}
	}
}

func CreateOutputFile(filePath string) (*os.File, func()) {
	directory, fileName := filepath.Split(filePath)
	hackFilePath := filepath.Join(directory, strings.Split(fileName, ".")[0]+".hack")
	hackFile, err := os.Create(hackFilePath)
	if err != nil {
		panic(err)
	}
	return hackFile, func() {
		hackFile.Close()
	}
}

func IsSymbol(symbol string) bool {
	if len(symbol) == 0 || unicode.IsDigit(rune(symbol[0])) {
		return false
	}

	for _, char := range symbol {
		if unicode.IsDigit(char) ||
			unicode.IsLetter(char) ||
			strings.ContainsAny(string(char), ".$:_") {
			continue
		}
		return false
	}
	return true
}
