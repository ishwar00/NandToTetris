package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ishwar00/JackAnalyzer/parserXML"
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
                parserxml.ParseIntoXML(filePath)
			}
		}
	} else if ext == ".jack" {
		parserxml.ParseIntoXML(os.Args[1])
	} else {
		errMsg := fmt.Errorf("invalid argument %s", path)
		panic(errMsg)
	}
}
