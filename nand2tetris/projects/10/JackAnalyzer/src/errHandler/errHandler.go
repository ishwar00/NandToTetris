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
	Length   int // this is the length of highlighted text, indicating error
	// starting from onColumn to onColumn + length
	File string // file path, must be absolute
}

var errQueue []Error

func init() {
	errQueue = make([]Error, 0)
}

func Add(errMsg Error) error {
	absPath, err := filepath.Abs(errMsg.File)
	if err != nil {
		return err
	}
	errMsg.File = absPath
	errQueue = append(errQueue, errMsg)
	return nil
}

func ReportAll() {
	fileErrs := map[string][]Error{}

	// classify errors on files
	for _, err := range errQueue {
		if _, ok := fileErrs[err.File]; !ok {
			fileErrs[err.File] = make([]Error, 0)
		}
		fileErrs[err.File] = append(fileErrs[err.File], err)
	}

	for file, errs := range fileErrs { // Eeee!! this is little more complicated than i anticipated
		errMsg := fmt.Sprintf("found %d errors in %s\n\n", len(errs), file)
		os.Stdout.WriteString(errMsg)

		sort.Slice(errs, func(i, j int) bool { return errs[i].OnLine < errs[j].OnLine })
		program := readFile(file)
		errBuf := map[string][]Error{} // buffer for grouping errors with same error message on same line

		curLine := errs[0].OnLine // buffering error messages of curren line
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
		errPointer += color.RedString(createString(err.Length, '^'))
	}
	errLine += line[last:]

	buf := fmt.Sprintf("on line %v:", errs[0].OnLine)
	errPointer = createString(len(buf), ' ') + errPointer
	errMsg := fmt.Sprintf("%s%s\n", buf, errLine)
	errMsg += errPointer
	errMsg += fmt.Sprintf("\n\t %s\n", color.GreenString(errs[0].ErrMsg))

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

func ClearBuffer() {
	errQueue = make([]Error, 0)
}