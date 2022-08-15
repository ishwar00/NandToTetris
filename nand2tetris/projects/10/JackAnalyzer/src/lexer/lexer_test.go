package lexer

import (
	"fmt"
	"testing"

	"github.com/ishwar00/JackAnalyzer/token"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		expectedType     token.TokenType
		expectedLiteral  string
		expectedOnLine   int
		expectedOnColumn int
		expectedInFile   string
	}{
		{token.LBRACE, "{", 5, 0, "test.jack"},
		{token.RBRACE, "}", 5, 1, "test.jack"},
		{token.LPAREN, "(", 5, 2, "test.jack"},
		{token.RPAREN, ")", 5, 3, "test.jack"},
		{token.LBRACK, "[", 5, 4, "test.jack"},
		{token.RBRACK, "]", 5, 5, "test.jack"},
		{token.PERIOD, ".", 5, 6, "test.jack"},
		{token.COMMA, ",", 5, 7, "test.jack"},
		{token.SEMICO, ";", 5, 8, "test.jack"},
		{token.PLUS, "+", 5, 9, "test.jack"},
		{token.MINUS, "-", 5, 10, "test.jack"},
		{token.ASTERI, "*", 5, 11, "test.jack"},
		{token.SLASH, "/", 5, 13, "test.jack"},
		{token.AMPERS, "&", 5, 14, "test.jack"},
		{token.PIPE, "|", 5, 15, "test.jack"},
		{token.LT, "<", 6, 0, "test.jack"},
		{token.GT, ">", 6, 1, "test.jack"},
		{token.PLUS, "+", 6, 2, "test.jack"},
		{token.TILDE, "~", 6, 3, "test.jack"},

		{token.CLASS, "class", 8, 0, "test.jack"},
		{token.IDENT, "main", 8, 6, "test.jack"},
		{token.LBRACE, "{", 8, 10, "test.jack"},
		{token.FIELD, "field", 9, 4, "test.jack"},
		{token.STATIC, "static", 9, 10, "test.jack"},
		{token.INT, "int", 9, 17, "test.jack"},
		{token.BOOLEAN, "boolean", 9, 21, "test.jack"},
		{token.CHAR, "char", 9, 29, "test.jack"},
		{token.IDENT, "c", 9, 34, "test.jack"},
		{token.COMMA, ",", 9, 35, "test.jack"},
		{token.IDENT, "d", 9, 37, "test.jack"},
		{token.SEMICO, ";", 9, 38, "test.jack"},
		{token.METHOD, "method", 10, 4, "test.jack"},
		{token.CONSTRUCTOR, "constructor", 10, 11, "test.jack"},
		{token.FUNCTION, "function", 10, 23, "test.jack"},
		{token.VOID, "void", 10, 32, "test.jack"},
		{token.IDENT, "func", 10, 37, "test.jack"},
		{token.LPAREN, "(", 10, 41, "test.jack"},
		{token.IDENT, "a", 10, 42, "test.jack"},
		{token.COMMA, ",", 10, 43, "test.jack"},
		{token.IDENT, "b", 10, 45, "test.jack"},
		{token.RPAREN, ")", 10, 46, "test.jack"},
		{token.LBRACE, "{", 10, 47, "test.jack"},
		{token.LET, "let", 11, 8, "test.jack"},
		{token.IDENT, "a_b", 11, 12, "test.jack"},
		{token.EQ, "=", 11, 15, "test.jack"},
		{token.IDENT, "greet", 11, 16, "test.jack"},
		{token.LPAREN, "(", 11, 21, "test.jack"},
		{token.RPAREN, ")", 11, 22, "test.jack"},
		{token.SEMICO, ";", 11, 23, "test.jack"},
		{token.DO, "do", 12, 8, "test.jack"},
		{token.IDENT, "greet", 12, 11, "test.jack"},
		{token.LPAREN, "(", 12, 16, "test.jack"},
		{token.RPAREN, ")", 12, 17, "test.jack"},
		{token.SEMICO, ";", 12, 18, "test.jack"},
		{token.VAR, "var", 13, 8, "test.jack"},
		{token.INT, "int", 13, 12, "test.jack"},
		{token.IDENT, "g", 13, 16, "test.jack"},
		{token.COMMA, ",", 13, 17, "test.jack"},
		{token.IDENT, "t", 13, 19, "test.jack"},
		{token.SEMICO, ";", 13, 20, "test.jack"},
		{token.TRUE, "true", 14, 8, "test.jack"},
		{token.COMMA, ",", 14, 12, "test.jack"},
		{token.FALSE, "false", 14, 13, "test.jack"},
		{token.COMMA, ",", 14, 18, "test.jack"},
		{token.NULL, "null", 14, 19, "test.jack"},
		{token.SEMICO, ";", 14, 23, "test.jack"},
		{token.RETURN, "return", 15, 8, "test.jack"},
		{token.THIS, "this", 15, 15, "test.jack"},
		{token.SEMICO, ";", 15, 19, "test.jack"},
		{token.IF, "if", 16, 8, "test.jack"},
		{token.LPAREN, "(", 16, 10, "test.jack"},
		{token.IDENT, "a", 16, 11, "test.jack"},
		{token.RPAREN, ")", 16, 12, "test.jack"},
		{token.LBRACE, "{", 16, 13, "test.jack"},
		{token.LET, "let", 16, 14, "test.jack"},
		{token.IDENT, "a", 16, 18, "test.jack"},
		{token.EQ, "=", 16, 19, "test.jack"},
		{token.FALSE, "false", 16, 20, "test.jack"},
		{token.SEMICO, ";", 16, 25, "test.jack"},
		{token.RBRACE, "}", 16, 26, "test.jack"},
		{token.ELSE, "else", 16, 27, "test.jack"},
		{token.LBRACE, "{", 16, 31, "test.jack"},
		{token.WHILE, "while", 16, 32, "test.jack"},
		{token.RBRACE, "}", 16, 37, "test.jack"},
		{token.LET, "let", 17, 8, "test.jack"},
		{token.IDENT, "a", 17, 12, "test.jack"},
		{token.EQ, "=", 17, 13, "test.jack"},
		{token.INT, "34322", 17, 14, "test.jack"},
		{token.SEMICO, ";", 17, 19, "test.jack"},
		{token.DO, "do", 18, 8, "test.jack"},
		{token.IDENT, "get", 18, 11, "test.jack"},
		{token.LPAREN, "(", 18, 14, "test.jack"},
		{token.INT, "3", 18, 15, "test.jack"},
		{token.COMMA, ",", 18, 16, "test.jack"},
		{token.INT, "5", 18, 17, "test.jack"},
		{token.RPAREN, ")", 18, 18, "test.jack"},
		{token.SEMICO, ";", 18, 19, "test.jack"},
		{token.LET, "let", 19, 8, "test.jack"},
		{token.IDENT, "str", 19, 12, "test.jack"},
		{token.EQ, "=", 19, 15, "test.jack"},
		{token.STR_CONST, "hey! there", 19, 17, "test.jack"},
		{token.SEMICO, ";", 19, 28, "test.jack"},
		{token.LET, "let", 20, 8, "test.jack"},
		{token.IDENT, "foo_4_bar", 20, 12, "test.jack"},
		{token.EQ, "=", 20, 22, "test.jack"},
		{token.STR_CONST, "deepsource", 20, 25, "test.jack"},
		{token.SEMICO, ";", 20, 36, "test.jack"},
		{token.RBRACE, "}", 21, 4, "test.jack"},
		{token.RBRACE, "}", 24, 5, "test.jack"},
		{token.RBRACE, "}", 25, 0, "test.jack"},
	}

	l, err := LexFile("test.jack")
	if err != nil {
		t.Fatal(err)
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			fmt.Println(tok)
			t.Fatalf("tests[%d]: expectedType=%v, but got=%v\n",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			fmt.Println(tok)
			t.Fatalf("tests[%d]: expectedLiteral=%v, but got=%v\n",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.OnLine != tt.expectedOnLine {
			fmt.Println(tok)
			t.Fatalf("tests[%d]: expectedOnLine=%v, but got=%v\n",
				i, tt.expectedOnLine, tok.OnLine)
		}

		if tok.OnColumn != tt.expectedOnColumn {
			fmt.Println(tok)
			t.Fatalf("tests[%d]: expectedOnColumn=%v, but got=%v\n",
				i, tt.expectedOnColumn, tok.OnColumn)
		}

		if tok.InFile != tt.expectedInFile {
			fmt.Println(tok)
			t.Fatalf("tests[%d]: expectedInFile=%v, but got=%v\n",
				i, tt.expectedInFile, tok.InFile)
		}
	}
}
