package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/ishwar00/JackAnalyzer/repl"
)


func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hi %s!, welcome to Jack REPL\n", user.Username)
    fmt.Println("Go ahead, type in some Jack")
    repl.Start(os.Stdin, os.Stdout)
}

