package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/parser"
)

func main() {
	if len(os.Args) != 2 {
		panic("program requires arugments")
	}
	path := os.Args[1]
	ext := filepath.Ext(path)
	if ext == "" {
		// may be it is directory
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			ext := filepath.Ext(file.Name())
			if !file.IsDir() && ext == ".jack" {
				filePath := filepath.Join(path, file.Name())
				l, err := lexer.LexFile(filePath)
				if err != nil {
					panic(err)
				}
				p := parser.New(l)
				program := p.ParseProgram()
				fmt.Println(program.String())
			}
		}

	} else if ext == ".jack" {
		l, err := lexer.LexFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		p := parser.New(l)
		program := p.ParseProgram()
		l.ReportErrors()
		p.ReportErrors()
		fmt.Println(program.Statements)
	} else {
		errMsg := fmt.Errorf("invalid argument %s", path)
		panic(errMsg)
	}
}
