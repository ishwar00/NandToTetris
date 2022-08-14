package errhandler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
)

type Error struct {
	ErrMsg   string
	OnLine   int
	OnColumn int
	// this is the length of highlighted text, indicating error
	// starting from onColumn to onColumn + length
	Length int
	File   string // file path, must be absolute
}

type ErrHandler struct {
	queue []Error
}

func (e *ErrHandler) Add(errMsg Error) error {
	absPath, err := filepath.Abs(errMsg.File)
	if err != nil {
		return err
	}
	errMsg.File = absPath
	e.queue = append(e.queue, errMsg)
	return nil
}

func (e ErrHandler) ReportAll() {
	fileErrs := map[string][]Error{}

	// classify errors by files
	for _, err := range e.queue {
		if _, ok := fileErrs[err.File]; !ok {
			fileErrs[err.File] = make([]Error, 0)
		}
		fileErrs[err.File] = append(fileErrs[err.File], err)
	}

	for file, errs := range fileErrs { // Eeee!! this is little more complicated than i anticipated
		errMsg := fmt.Sprintf("found %d error(s) in %s\n\n", len(errs), file)
		os.Stdout.WriteString(errMsg)

		sort.Slice(errs, func(i, j int) bool { return errs[i].OnLine < errs[j].OnLine })
		program := readFile(file)
		errBuf := map[string][]Error{} // buffer for grouping errors with same error message on same line

		curLine := errs[0].OnLine // buffering error messages of current line
		for _, err := range errs {
			if curLine == err.OnLine {
				if _, ok := errBuf[err.ErrMsg]; !ok {
					errBuf[err.ErrMsg] = make([]Error, 0) // creating empty buffer for Error objects
				}
				errBuf[err.ErrMsg] = append(errBuf[err.ErrMsg], err)
			} else {
				for errMsg := range errBuf {
					report(program[curLine], errBuf[errMsg])
				}
				errBuf = make(map[string][]Error)
				errBuf[err.ErrMsg] = []Error{err}
				curLine = err.OnLine
			}
		}
		for errMsg := range errBuf {
			report(program[curLine], errBuf[errMsg])
		}
	}
}

/* why errs is a slice?: there can be mulitple errors
with same message on a single line */
func report(line string, errs []Error) {
	sort.Slice(errs, func(i, j int) bool { return errs[i].OnColumn < errs[j].OnColumn })
	var errPointer string
	var errLine string
	last := 0
	for _, err := range errs {
		errLine += line[last:err.OnColumn]
		errPointer += createString(err.OnColumn-last, ' ')

		last = err.OnColumn + err.Length
		errLine += color.RedString(line[err.OnColumn:last])
		errPointer += createString(err.Length, '^')
	}
	errLine += line[last:]

	lineNumber := fmt.Sprintf(" %d| ", errs[0].OnLine+1)
	padding := createString(len(lineNumber)-2, ' ') + "| " // ----|
	errMsg := padding + "\n"
	errMsg += lineNumber + errLine + "\n"
	errMsg += padding + errPointer + "\n"
	errMsg += padding + errs[0].ErrMsg + "\n\n"
	os.Stdout.WriteString(errMsg)
}

func createString(length int, char rune) string {
	var str string
	for i := 0; i < length; i++ {
		str += string(char)
	}
	return str
}

func readFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err) // let's just panic as of now
	}

	scanner := bufio.NewScanner(file)
	var program []string
	for scanner.Scan() {
		program = append(program, scanner.Text())
	}
	return program
}

func (e *ErrHandler) QueueSize() int {
	return len(e.queue)
}

func (e *ErrHandler) ClearQueue() {
	e.queue = make([]Error, 0)
}
