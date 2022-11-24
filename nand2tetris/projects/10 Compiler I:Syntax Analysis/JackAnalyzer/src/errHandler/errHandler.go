package errhandler

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// implementation is inspired from vlang
// https://github.com/vlang/v/blob/6b0743bb07f8d003ff703dc48ae2120cd8bc152a/vlib/v/util/errors.v#L20

// error_context_before - how many lines of source context to print before the pointer line
// error_context_after - ^^^ same, but after
const (
	error_context_before = 2
	error_context_after  = 2
)

type Error struct {
	ErrMsg   string
	OnLine   int
	OnColumn int

	// Length: length of text, pointed by pointer line
	// starting from onColumn to onColumn + length
	Length int
	File   string
}

// filepath:line:col: error_message
func (e *Error) format() string {
	return fmt.Sprintf("%s:%d:%d: %s: %s",
		e.File, e.OnLine+1, e.OnColumn+1, color.RedString("error"), e.ErrMsg)
}

type ErrHandler struct {
	error_count int
	fileErrs    map[string][]Error // key: fileName, value: slice of Error in file fileName
}

// TODO: Add does not need to return error
func (eh *ErrHandler) Add(errMsg Error) {
	if eh.fileErrs == nil {
		eh.fileErrs = make(map[string][]Error)
	}
	fileName := errMsg.File
	_, ok := eh.fileErrs[fileName]
	if !ok {
		eh.fileErrs[fileName] = make([]Error, 0)
	}
	eh.fileErrs[fileName] = append(eh.fileErrs[fileName], errMsg)
	eh.error_count++
}

func (e ErrHandler) ReportAll() {
	for file, errs := range e.fileErrs {
		cmpr := func(i, j int) bool {
			return errs[i].OnLine < errs[j].OnLine
		}
		sort.Slice(errs, cmpr)
		source := readFile(file)

		tab := "    "
		for _, err := range errs {
			bline := min(err.OnLine, max(0, err.OnLine-error_context_before))
			aline := min(len(source)-1, err.OnLine+error_context_after)

			formated_errMsg := fmt.Sprintf("\n%s\n", err.format())
			os.Stdout.WriteString(formated_errMsg)

			for onLine := bline; onLine <= aline; onLine++ {
				sline := strings.ReplaceAll(source[onLine], "\t", tab)
				cxt_line := fmt.Sprintf("%5d | %s\n", onLine+1, sline)
				os.Stdout.WriteString(cxt_line)
				if err.OnLine == onLine {
					offset := strings.Repeat(" ", err.OnColumn-1)
					pointer_str := strings.Repeat("^", err.Length)
					pointer_line := fmt.Sprintln("      |", offset, color.RedString(pointer_str))
					os.Stdout.WriteString(pointer_line)
				}
			}
		}
	}
}

func (eh *ErrHandler) Error_count() int {
	return eh.error_count
}

func max(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func readFile(filePath string) []string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	var source []string
	for scanner.Scan() {
		source = append(source, scanner.Text())
	}
	return source
}
