package parser

import (
	"testing"

	"github.com/ishwar00/JackAnalyzer/ast"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

func TestVarDec(t *testing.T) {
	tests := []struct {
		input              string
		expectedToken      token.Token // var
		expectedDataType   token.Token
		expectedIdentifers []*ast.Identifier
	}{
		{ // test[0]
			input: "var int a, b,c;", // input
			expectedToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 0,
			},
			expectedDataType: token.Token{
				Literal:  "int",
				Type:     token.INT,
				OnLine:   0,
				OnColumn: 4,
			},
			expectedIdentifers: []*ast.Identifier{
				{ // identifier a
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
					Value: "a",
					DataType: token.Token{
						Literal:  "int",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 4,
					},
				},
				{ // identifier b
					Token: token.Token{
						Literal:  "b",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "b",
					DataType: token.Token{
						Literal:  "int",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 4,
					},
				},
				{ // identifier c
					Token: token.Token{
						Literal:  "c",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 13,
					},
					Value: "c",
					DataType: token.Token{
						Literal:  "int",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 4,
					},
				},
			},
		},
		{ // test[1]
			input: "  var Ball a_b, __a,ab9c;", // input
			expectedToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 2,
			},
			expectedDataType: token.Token{
				Literal:  "Ball",
				Type:     token.IDENT,
				OnLine:   0,
				OnColumn: 6,
			},
			expectedIdentifers: []*ast.Identifier{
				{ // identifier a_b
					Token: token.Token{
						Literal:  "a_b",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "a_b",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
				},
				{ // identifier __a
					Token: token.Token{
						Literal:  "__a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 16,
					},
					Value: "__a",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
				},
				{ // identifier ab9c
					Token: token.Token{
						Literal:  "ab9c",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 20,
					},
					Value: "ab9c",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
				},
			},
		},
		{ // test[2]
			input: "  var Ball foo_bar;", // input
			expectedToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 2,
			},
			expectedDataType: token.Token{
				Literal:  "Ball",
				Type:     token.IDENT,
				OnLine:   0,
				OnColumn: 6,
			},
			expectedIdentifers: []*ast.Identifier{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "foo_bar",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		if p.errors.QueueSize() != 0 {
			p.ReportErrors()
			t.Fatal("terminating testing")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d]: size of program.Statements is not %d, got=%d", i, 1, len(program.Statements))
		}

		stmt := program.Statements[0]
		varDec, ok := stmt.(*ast.VarDecStatement)
		if !ok {
			t.Fatalf("test[%d]: stmt is not *ast.VarDecStatement, got=%T", i, stmt)
		}

		if !testToken(t, tt.expectedToken, varDec.Token) {
			t.Fatalf("test[%d]: varDec.Token is not %v, got=%v", i, tt.expectedToken, varDec.Token)
		}

		if !testToken(t, tt.expectedDataType, varDec.DataType) {
			t.Fatalf("test[%d]: varDec.DataType is not %v, got=%v",
				i, tt.expectedDataType, varDec.DataType)
		}

		for j, expectedIdentifer := range tt.expectedIdentifers {
			if !testIdentifier(t, expectedIdentifer, varDec.Identifiers[j]) {
				t.Fatalf("test[%d]: varDec.Identifiers[%d] is not %v, got=%v",
					i, j, expectedIdentifer, varDec.Identifiers[j])
			}
		}
	}
}

func testIdentifier(t *testing.T, ref *ast.Identifier, tok *ast.Identifier) bool {
	if !testToken(t, ref.Token, tok.Token) {
		t.Errorf("tok.Token is not %v, got=%v", ref.Token, tok.Token)
		return false
	}

	if ref.Value != tok.Value {
		t.Errorf("tok.Value is not %s, got=%s", ref.Value, tok.Value)
		return false
	}

	if !testToken(t, ref.DataType, tok.DataType) {
		t.Errorf("tok.Type is not %v, got=%v", ref.DataType, tok.DataType)
		return false
	}
	return true
}

func testToken(t *testing.T, ref token.Token, tok token.Token) bool {
	if ref.Literal != tok.Literal {
		t.Errorf("tok.Literal is not %s, got=%s", ref.Literal, tok.Literal)
		return false
	}

	if ref.Type != tok.Type {
		t.Errorf("tok.Type is not %d, got=%d", ref.Type, tok.Type)
		return false
	}

	if ref.OnColumn != tok.OnColumn {
		t.Errorf("ref.OnColumn is not %d, but got=%d", ref.OnColumn, tok.OnColumn)
		return false
	}

	if ref.OnLine != tok.OnLine {
		t.Errorf("ref.OnLine is not %d, got=%d", ref.OnLine, tok.OnLine)
		return false
	}

	if ref.InFile != tok.InFile {
		t.Errorf("ref.InFile is not %s, got=%s", ref.InFile, tok.InFile)
		return false
	}
	return true
}
