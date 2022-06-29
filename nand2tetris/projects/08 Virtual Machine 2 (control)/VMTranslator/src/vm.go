package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func virtualMachine(vmSource string) error {
	vmSourceFiles, err := vmFiles(vmSource)
	if err != nil {
		return err
	}

	var parser Parser
	var codeWriter CodeWriter

	codeWriter.Initialize(vmSource)
	codeWriter.InjectBootstrapCode()
	defer codeWriter.closer()

	for _, vmFilePath := range vmSourceFiles {
		parser.Initialize(vmFilePath)
		codeWriter.SetFilePath(vmFilePath)
		comment := "// "
		for parser.HasMoreLines() {
			parser.Advance()
			codeWriter.Write([]string{comment + parser.GetInstrInfo().Instruction})
			fields := strings.Fields(parser.GetInstrInfo().Instruction)
			switch commandType := parser.CommandType(); commandType {
			case C_PUSH, C_POP:
				segment := parser.Arg1()
				index, err := parser.Arg2()
				if err != nil { // failed to parse index
					return err
				}
				err = codeWriter.WritePushPop(commandType, segment, index)
				if err != nil { // failed while either writing into file or parsing
					return PrepError(parser.GetInstrInfo(), err)
				}
			case C_ARITHMETIC:
				command := parser.Arg1()
				err := codeWriter.WriteArithmetic(command)
				if err != nil {
					return PrepError(parser.GetInstrInfo(), err)
				}

			case C_LABEL:
				err := codeWriter.WriteLabel(fields[1], parser.CurrentFunction)
				if err != nil {
					return err
				}
			case C_GOTO:
				err := codeWriter.WriteGoto(fields[1], parser.CurrentFunction)
				if err != nil {
					return err
				}
			case C_IF:
				err := codeWriter.WriteIf(fields[1], parser.CurrentFunction)
				if err != nil {
					return err
				}
			case C_FUNCTION:
				nVars, err := strconv.ParseInt(fields[2], 10, 32)
				if err != nil {
					return PrepError(parser.GetInstrInfo(), err)
				}
				err = codeWriter.WriteFunction(parser.CurrentFunction, int(nVars))
				if err != nil {
					return err
				}

			case C_CALL:
				nArgs, err := strconv.ParseInt(fields[2], 10, 32)
				if err != nil {
					return PrepError(parser.GetInstrInfo(), err)
				}
				err = codeWriter.WriteCall(fields[1], parser.CurrentFunction, int(nArgs), parser.CallCount)
				if err != nil {
					return err
				}
			case C_RETURN:
				if err := codeWriter.WriteReturn(); err != nil {
					return err
				}
			default:
				errMsg := fmt.Errorf(color.RedString("unrecognized command"))
				return PrepError(parser.GetInstrInfo(), errMsg)
			}
		}
	}
	return nil
}

func vmFiles(vmSourcePath string) ([]string, error) {
	vmSource, err := os.Open(vmSourcePath)
	if err != nil {
		return nil, err
	}
	defer vmSource.Close()

	vmSourceInfo, err := vmSource.Stat()
	if err != nil {
		return nil, err
	}

	var vmFiles []string
	if vmSourceInfo.IsDir() {
		vmFileStats, err := ioutil.ReadDir(vmSourcePath)
		if err != nil {
			return nil, err
		}

		for _, fileStat := range vmFileStats {
			source := filepath.Join(vmSourcePath, fileStat.Name())
			if filepath.Ext(source) == ".vm" {
				vmFiles = append(vmFiles, source)
			}
		}
	} else {
		vmFiles = append(vmFiles, vmSourcePath)
	}
	return vmFiles, nil
}
