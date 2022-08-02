package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

const PROMPT = "> "

func Start(in io.Reader, _ io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Println("Exit > .exit")

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == ".exit" {
			return
		}
		l := lexer.LexString(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
