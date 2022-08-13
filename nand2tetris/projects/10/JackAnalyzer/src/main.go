package main

import (
	// "fmt"
	// "os"
	// "os/user"

	// "github.com/ishwar00/JackAnalyzer/repl"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	lexerxml "github.com/ishwar00/JackAnalyzer/lexerXML"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	//     panic(err)
	// }
	// fmt.Printf("Hi %s!, welcome to Jack REPL\n", user.Username)
	// fmt.Println("Go ahead, type in some Jack")
	// repl.Start(os.Stdin, os.Stdout)

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
				lexerxml.Run(filePath)
			}
		}

	} else if ext == ".jack" {
		lexerxml.Run(path)
	} else {
		errMsg := fmt.Errorf("invalid argument %s", path)
		panic(errMsg)
	}
}
