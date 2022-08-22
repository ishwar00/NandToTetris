package lexerxml

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

func Run(file string) error {
	l, err := lexer.LexFile(file)
	if err != nil {
		return err
	}

	indent := "\t"

	outputFile, closer, err := createOutputFile(file)
	if err != nil {
		return err
	}
	defer closer()

	outputFile.WriteString("<tokens>\n")
	for tok, buf := l.NextToken(), ""; tok.Type != token.EOF; tok = l.NextToken() {
		value := int(tok.Type)
		// keyword range
		if 0 <= value && value <= 20 {
			value = 100
		}

		// symbol range
		if 21 <= value && value <= 36 {
			value = 200
		}

		switch value {
		case 100:
			buf = fmt.Sprintf("%s<keyword>%v</keyword>\n", indent, tok.Literal)
		case 200:
			buf = fmt.Sprintf("%s<symbol>%v</symbol>\n", indent, tok.Literal)
		case token.LT:
			buf = fmt.Sprintf("%s<symbol>&lt;</symbol>\n", indent)
		case token.GT:
			buf = fmt.Sprintf("%s<symbol>&gt;</symbol>\n", indent)
		case token.AMPERS:
			buf = fmt.Sprintf("%s<symbol>&amp;</symbol>\n", indent)
		case token.IDENT:
			buf = fmt.Sprintf("%s<identifier>%v</identifier>\n", indent, tok.Literal)
		case token.INT_CONST:
			buf = fmt.Sprintf("%s<integerConstant>%v</integerConstant>\n", indent, tok.Literal)
		case token.STR_CONST:
			buf = fmt.Sprintf("%s<stringConstant>%v</stringConstant>\n", indent, tok.Literal)
		default:
			buf = fmt.Sprintf("%s<illegal>%v</illegal>\n", indent, tok.Literal)
		}
		outputFile.WriteString(buf)
	}
	outputFile.WriteString("</tokens>")
	if l.FoundErrors() {
		l.ReportErrors()
	}
	return nil
}

func createOutputFile(path string) (*os.File, func(), error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return nil, func() {}, fmt.Errorf("%s, is not a file or it does not end with file extension .jack", path)
	}

	if ext != ".jack" {
		return nil, func() {}, fmt.Errorf("%s, file extension must end with .jack", path)
	}
	// TT just to avoid name collison with with testing files ending with T.xml
	// this lexer supposed to produce similar output to file ending with T.xml file
	outputPath := path[0:len(path)-len(ext)] + "TT.xml"
	file, err := os.Create(outputPath)
	return file, func() { file.Close() }, err
}
