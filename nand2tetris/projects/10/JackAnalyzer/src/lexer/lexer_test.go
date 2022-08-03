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
		{token.LBRACE, "{", 0, 0, "test.jack"},
		{token.RBRACE, "}", 0, 1, "test.jack"},
		{token.LPAREN, "(", 0, 2, "test.jack"},
		{token.RPAREN, ")", 0, 3, "test.jack"},
		{token.LBRACK, "[", 0, 4, "test.jack"},
		{token.RBRACK, "]", 0, 5, "test.jack"},
		{token.PERIOD, ".", 0, 6, "test.jack"},
		{token.COMMA, ",", 0, 7, "test.jack"},
		{token.SEMICO, ";", 0, 8, "test.jack"},
		{token.PLUS, "+", 0, 9, "test.jack"},
		{token.MINUS, "-", 0, 10, "test.jack"},
		{token.ASTERI, "*", 0, 11, "test.jack"},
		{token.SLASH, "/", 0, 13, "test.jack"},
		{token.AMPERS, "&", 0, 14, "test.jack"},
		{token.PIPE, "|", 0, 15, "test.jack"},
		{token.LT, "<", 1, 0, "test.jack"},
		{token.GT, ">", 1, 1, "test.jack"},
		{token.PLUS, "+", 1, 2, "test.jack"},
		{token.TILDE, "~", 1, 3, "test.jack"},

		{token.CLASS, "class", 3, 0, "test.jack"},
		{token.IDENT, "main", 3, 6, "test.jack"},
		{token.LBRACE, "{", 3, 10, "test.jack"},
		{token.FIELD, "field", 4, 4, "test.jack"},
		{token.STATIC, "static", 4, 10, "test.jack"},
		{token.INT, "int", 4, 17, "test.jack"},
		{token.BOOLEAN, "boolean", 4, 21, "test.jack"},
		{token.CHAR, "char", 4, 29, "test.jack"},
		{token.IDENT, "c", 4, 34, "test.jack"},
		{token.COMMA, ",", 4, 35, "test.jack"},
		{token.IDENT, "d", 4, 37, "test.jack"},
		{token.SEMICO, ";", 4, 38, "test.jack"},
		{token.METHOD, "method", 5, 4, "test.jack"},
		{token.CONSTRUCTOR, "constructor", 5, 11, "test.jack"},
		{token.FUNCTION, "function", 5, 23, "test.jack"},
		{token.VOID, "void", 5, 32, "test.jack"},
		{token.IDENT, "func", 5, 37, "test.jack"},
		{token.LPAREN, "(", 5, 41, "test.jack"},
		{token.IDENT, "a", 5, 42, "test.jack"},
		{token.COMMA, ",", 5, 43, "test.jack"},
		{token.IDENT, "b", 5, 45, "test.jack"},
		{token.RPAREN, ")", 5, 46, "test.jack"},
		{token.LBRACE, "{", 5, 47, "test.jack"},
		{token.LET, "let", 6, 8, "test.jack"},
		{token.IDENT, "a_b", 6, 12, "test.jack"},
		{token.EQ, "=", 6, 15, "test.jack"},
		{token.IDENT, "greet", 6, 16, "test.jack"},
		{token.LPAREN, "(", 6, 21, "test.jack"},
		{token.RPAREN, ")", 6, 22, "test.jack"},
		{token.SEMICO, ";", 6, 23, "test.jack"},
		{token.DO, "do", 7, 8, "test.jack"},
		{token.IDENT, "greet", 7, 11, "test.jack"},
		{token.LPAREN, "(", 7, 16, "test.jack"},
		{token.RPAREN, ")", 7, 17, "test.jack"},
		{token.SEMICO, ";", 7, 18, "test.jack"},
		{token.VAR, "var", 8, 8, "test.jack"},
		{token.INT, "int", 8, 12, "test.jack"},
		{token.IDENT, "g", 8, 16, "test.jack"},
		{token.COMMA, ",", 8, 17, "test.jack"},
		{token.IDENT, "t", 8, 19, "test.jack"},
		{token.SEMICO, ";", 8, 20, "test.jack"},
		{token.TRUE, "true", 9, 8, "test.jack"},
		{token.COMMA, ",", 9, 12, "test.jack"},
		{token.FALSE, "false", 9, 13, "test.jack"},
		{token.COMMA, ",", 9, 18, "test.jack"},
		{token.NULL, "null", 9, 19, "test.jack"},
		{token.SEMICO, ";", 9, 23, "test.jack"},
		{token.RETURN, "return", 10, 8, "test.jack"},
		{token.THIS, "this", 10, 15, "test.jack"},
		{token.SEMICO, ";", 10, 19, "test.jack"},
		{token.IF, "if", 11, 8, "test.jack"},
		{token.LPAREN, "(", 11, 10, "test.jack"},
		{token.IDENT, "a", 11, 11, "test.jack"},
		{token.RPAREN, ")", 11, 12, "test.jack"},
		{token.LBRACE, "{", 11, 13, "test.jack"},
		{token.LET, "let", 11, 14, "test.jack"},
		{token.IDENT, "a", 11, 18, "test.jack"},
		{token.EQ, "=", 11, 19, "test.jack"},
		{token.FALSE, "false", 11, 20, "test.jack"},
		{token.SEMICO, ";", 11, 25, "test.jack"},
		{token.RBRACE, "}", 11, 26, "test.jack"},
		{token.ELSE, "else", 11, 27, "test.jack"},
		{token.LBRACE, "{", 11, 31, "test.jack"},
		{token.WHILE, "while", 11, 32, "test.jack"},
		{token.RBRACE, "}", 11, 37, "test.jack"},
		{token.LET, "let", 12, 8, "test.jack"},
		{token.IDENT, "a", 12, 12, "test.jack"},
		{token.EQ, "=", 12, 13, "test.jack"},
		{token.INT_CONST, "34322", 12, 14, "test.jack"},
		{token.SEMICO, ";", 12, 19, "test.jack"},
		{token.DO, "do", 13, 8, "test.jack"},
		{token.IDENT, "get", 13, 11, "test.jack"},
		{token.LPAREN, "(", 13, 14, "test.jack"},
		{token.INT_CONST, "3", 13, 15, "test.jack"},
		{token.COMMA, ",", 13, 16, "test.jack"},
		{token.INT_CONST, "5", 13, 17, "test.jack"},
		{token.RPAREN, ")", 13, 18, "test.jack"},
		{token.SEMICO, ";", 13, 19, "test.jack"},
		{token.LET, "let", 14, 8, "test.jack"},
		{token.IDENT, "str", 14, 12, "test.jack"},
		{token.EQ, "=", 14, 15, "test.jack"},
		{token.STR_CONST, "hey! there", 14, 17, "test.jack"},
		{token.SEMICO, ";", 14, 28, "test.jack"},
		{token.LET, "let", 15, 8, "test.jack"},
		{token.IDENT, "foo_4_bar", 15, 12, "test.jack"},
		{token.EQ, "=", 15, 22, "test.jack"},
		{token.STR_CONST, "deepsource", 15, 25, "test.jack"},
		{token.SEMICO, ";", 15, 36, "test.jack"},
		{token.RBRACE, "}", 16, 4, "test.jack"},
		{token.RBRACE, "}", 17, 0, "test.jack"},
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
