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

	ident := "\t"

	outputFile, closer, err := createOutputFile(file)
	if err != nil {
		return err
	}
	defer closer()

	outputFile.WriteString("<tokens>\n")
	for tok, buf := l.NextToken(), ""; tok.Type != token.EOF; tok = l.NextToken() {
		switch tok.Type {
		case token.CLASS, token.FUNCTION, token.METHOD, token.CONSTRUCTOR, token.FIELD,
			token.STATIC, token.VAR, token.CHAR, token.BOOLEAN,
			token.TRUE, token.FALSE, token.NULL, token.THIS, token.LET,
			token.IF, token.ELSE, token.RETURN, token.DO, token.VOID,
			token.WHILE:

			buf = fmt.Sprintf("%s<keyword>%v</keyword>\n", ident, tok.Literal)
		case token.LBRACE, token.RBRACE, token.LBRACK, token.RBRACK, token.RPAREN,
			token.LPAREN, token.PERIOD, token.COMMA, token.SEMICO, token.PLUS,
			token.MINUS, token.ASTERI, token.SLASH, token.EQ, token.TILDE,
			token.PIPE:

			buf = fmt.Sprintf("%s<symbol>%v</symbol>\n", ident, tok.Literal)
		case token.LT:
			buf = fmt.Sprintf("%s<symbol>&lt;</symbol>\n", ident)
		case token.GT:
			buf = fmt.Sprintf("%s<symbol>&gt;</symbol>\n", ident)
		case token.AMPERS:
			buf = fmt.Sprintf("%s<symbol>&amp;</symbol>\n", ident)
		case token.IDENT:
			buf = fmt.Sprintf("%s<identifier>%v</identifier>\n", ident, tok.Literal)
		case token.INT:
			buf = fmt.Sprintf("%s<integerConstant>%v</integerConstant>\n", ident, tok.Literal)
		case token.STR_CONST:
			buf = fmt.Sprintf("%s<stringConstant>%v</stringConstant>\n", ident, tok.Literal)
		default:
			buf = fmt.Sprintf("%s<illegal>%v</illegal>\n", ident, tok.Literal)
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
