package parser

import (
	"reflect"
	"testing"

	"github.com/ishwar00/JackAnalyzer/ast"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

func TestReturnSta(t *testing.T) {
	tests := []struct {
		input        string
		expReturnSta *ast.ReturnSta
	}{
		{
			input: "return;",
			expReturnSta: &ast.ReturnSta{
				Token: token.Token{
					Literal:  "return",
					Type:     token.RETURN,
					OnLine:   0,
					OnColumn: 0,
				},
			},
		},
		{
			input: "return 34;",
			expReturnSta: &ast.ReturnSta{
				Token: token.Token{
					Literal:  "return",
					Type:     token.RETURN,
					OnLine:   0,
					OnColumn: 0,
				},
				Expression: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "34",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 7,
					},
					Value: 34,
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)

		stmt := p.parseReturnSta()
		if stmt == nil {
			t.Fatalf("tests[%d]: stmt is not *ast.ReturnSta, got=%v", i, stmt)
		}

		if stmt.Token != tt.expReturnSta.Token {
			t.Fatalf("tests[%d]: stmt.Token is not %+v, got=%+v",
				i, tt.expReturnSta.Token, stmt.Token)
		}

		if !reflect.DeepEqual(stmt.Expression, tt.expReturnSta.Expression) {
			t.Fatalf("tests[%d]: stmt.Expression is not %+v, got=%+v",
				i, tt.expReturnSta.Expression, stmt.Expression)
		}
	}
}

func TestWhileSta(t *testing.T) {
	tests := []struct {
		input       string
		expWhileSta *ast.WhileSta
	}{
		{
			input: "while(343) {}",
			expWhileSta: &ast.WhileSta{
				Token: token.Token{
					Literal:  "while",
					Type:     token.WHILE,
					OnLine:   0,
					OnColumn: 0,
				},
				Condition: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "343",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 6,
					},
					Value: 343,
				},
				Statements: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 11,
					},
				},
			},
		},
		{
			input: "while(abc) { let deep = source; }",
			expWhileSta: &ast.WhileSta{
				Token: token.Token{
					Literal:  "while",
					Type:     token.WHILE,
					OnLine:   0,
					OnColumn: 0,
				},
				Condition: &ast.VarNameExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
					Name: "abc",
				},
				Statements: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 11,
					},
					Statements: []ast.Statement{
						&ast.LetSta{
							Token: token.Token{
								Literal:  "let",
								Type:     token.LET,
								OnLine:   0,
								OnColumn: 13,
							},
							VarName: ast.VarNameExp{
								Token: token.Token{
									Literal:  "deep",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 17,
								},
								Name: "deep",
							},
							Expression: &ast.VarNameExp{
								Token: token.Token{
									Literal:  "source",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 24,
								},
								Name: "source",
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		stmt := p.parseWhileSta()
		if stmt == nil {
			t.Fatalf("tests[%d]: stmt is not *ast.WhileSta, got=%v", i, stmt)
		}

		if stmt.Token != tt.expWhileSta.Token {
			t.Fatalf("tests[%d]: stmt.Token is not %+v, got=%+v",
				i, tt.expWhileSta.Token, stmt.Token)
		}

		if !reflect.DeepEqual(stmt.Condition, tt.expWhileSta.Condition) {
			t.Fatalf("tests[%d]: stmt.Condition is not '%+v', got=%+v",
				i, tt.expWhileSta.Condition, stmt.Condition)
		}

		if !reflect.DeepEqual(stmt.Statements, tt.expWhileSta.Statements) {
			t.Fatalf("tests[%d]: stmt.Statements is not %v, got='%v'",
				i, stmt.Statements, tt.expWhileSta.Statements)
		}
	}
}

func TestStatementBlock(t *testing.T) {
	tests := []struct {
		input             string
		expStatementBlock ast.StatementBlockSta
	}{
		{
			input: "{ let a = 345; let a_b = abc; }",
			expStatementBlock: ast.StatementBlockSta{
				Token: token.Token{
					Literal:  "{",
					Type:     token.LBRACE,
					OnLine:   0,
					OnColumn: 0,
				},
				Statements: []ast.Statement{
					&ast.LetSta{
						Token: token.Token{
							Literal:  "let",
							Type:     token.LET,
							OnLine:   0,
							OnColumn: 2,
						},
						VarName: ast.VarNameExp{
							Token: token.Token{
								Literal:  "a",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 6,
							},
							Name: "a",
						},
						Expression: &ast.IntConstantExp{
							Token: token.Token{
								Literal:  "345",
								Type:     token.INT,
								OnLine:   0,
								OnColumn: 10,
							},
							Value: 345,
						},
					},
					&ast.LetSta{
						Token: token.Token{
							Literal:  "let",
							Type:     token.LET,
							OnLine:   0,
							OnColumn: 15,
						},
						VarName: ast.VarNameExp{
							Token: token.Token{
								Literal:  "a_b",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 19,
							},
							Name: "a_b",
						},
						Expression: &ast.VarNameExp{
							Token: token.Token{
								Literal:  "abc",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 25,
							},
							Name: "abc",
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)

		sb := p.parseStatementBlock()
		if sb == nil {
			t.Fatalf("tests[%d]: sb is not *ast.StatementBlockSta, got=%v", i, sb)
		}
		if sb.Token != tt.expStatementBlock.Token {
			t.Fatalf("tests[%d]: sb.Token is not '%+v', got='%+v'",
				i, tt.expStatementBlock.Token, sb.Token)
		}

		for j, stmt := range tt.expStatementBlock.Statements {
			if !reflect.DeepEqual(stmt, sb.Statements[j]) { // comparing interfaces here
				t.Fatalf("tests[%d], %dth statement: stmt is not '%+v', got='%v'", i, j, stmt, sb.Statements[j])
			}
		}
	}
}

func TestIfElseSta(t *testing.T) {
	tests := []struct {
		input        string
		expIfElseSta *ast.IfElseSta
	}{
		{
			input: "if(343){}",
			expIfElseSta: &ast.IfElseSta{
				Token: token.Token{
					Literal:  "if",
					Type:     token.IF,
					OnLine:   0,
					OnColumn: 0,
				},
				Condition: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "343",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 3,
					},
					Value: 343,
				},
				Then: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 7,
					},
				},
			},
		},
		{
			input: "if(abc){ let hey = 34; }",
			expIfElseSta: &ast.IfElseSta{
				Token: token.Token{
					Literal:  "if",
					Type:     token.IF,
					OnLine:   0,
					OnColumn: 0,
				},
				Condition: &ast.VarNameExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 3,
					},
					Name: "abc",
				},
				Then: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 7,
					},
					Statements: []ast.Statement{
						&ast.LetSta{
							Token: token.Token{
								Literal:  "let",
								Type:     token.LET,
								OnLine:   0,
								OnColumn: 9,
							},
							VarName: ast.VarNameExp{
								Token: token.Token{
									Literal:  "hey",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 13,
								},
								Name: "hey",
							},
							Expression: &ast.IntConstantExp{
								Token: token.Token{
									Literal:  "34",
									Type:     token.INT,
									OnLine:   0,
									OnColumn: 19,
								},
								Value: 34,
							},
						},
					},
				},
			},
		},
		{
			input: "if(abc){ let hey = 34; } else { let yo = 3; }",
			expIfElseSta: &ast.IfElseSta{
				Token: token.Token{
					Literal:  "if",
					Type:     token.IF,
					OnLine:   0,
					OnColumn: 0,
				},
				Condition: &ast.VarNameExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 3,
					},
					Name: "abc",
				},
				Then: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 7,
					},
					Statements: []ast.Statement{
						&ast.LetSta{
							Token: token.Token{
								Literal:  "let",
								Type:     token.LET,
								OnLine:   0,
								OnColumn: 9,
							},
							VarName: ast.VarNameExp{
								Token: token.Token{
									Literal:  "hey",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 13,
								},
								Name: "hey",
							},
							Expression: &ast.IntConstantExp{
								Token: token.Token{
									Literal:  "34",
									Type:     token.INT,
									OnLine:   0,
									OnColumn: 19,
								},
								Value: 34,
							},
						},
					},
				},
				Else: &ast.StatementBlockSta{
					Token: token.Token{
						Literal:  "{",
						Type:     token.LBRACE,
						OnLine:   0,
						OnColumn: 30,
					},
					Statements: []ast.Statement{
						&ast.LetSta{
							Token: token.Token{
								Literal:  "let",
								Type:     token.LET,
								OnLine:   0,
								OnColumn: 32,
							},
							VarName: ast.VarNameExp{
								Token: token.Token{
									Literal:  "yo",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 36,
								},
								Name: "yo",
							},
							Expression: &ast.IntConstantExp{
								Token: token.Token{
									Literal:  "3",
									Type:     token.INT,
									OnLine:   0,
									OnColumn: 41,
								},
								Value: 3,
							},
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)

		stmt := p.parseIfElseSta()
		if stmt == nil {
			t.Fatalf("tests[%d]: stmt is not *ast.IfElseSta, got=%v", i, stmt)
		}

		if stmt.Token != tt.expIfElseSta.Token {
			t.Fatalf("tests[%d]: stmt.Token is not %+v, got=%+v",
				i, tt.expIfElseSta.Token, stmt.Token)
		}

		if !reflect.DeepEqual(stmt.Condition, tt.expIfElseSta.Condition) {
			t.Fatalf("tests[%d]:stmt.Condition is not %+v, got=%+v",
				i, tt.expIfElseSta.Condition, stmt.Condition)
		}

		if !reflect.DeepEqual(stmt.Then, tt.expIfElseSta.Then) {
			t.Fatalf("tests[%d]:stmt.Then is not %+v, got=%+v",
				i, tt.expIfElseSta.Then, stmt.Then)
		}

		if !reflect.DeepEqual(stmt.Else, tt.expIfElseSta.Else) {
			t.Fatalf("tests[%d]:stmt.Else is not %v, got=%v",
				i, *tt.expIfElseSta.Else, *stmt.Else)
			// NOTE: not checking for nil pointer
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input     string
		expLetSta ast.LetSta
	}{
		{ // test 0
			input: "let a = 34;",
			expLetSta: ast.LetSta{
				Token: token.Token{
					Literal:  "let",
					Type:     token.LET,
					OnLine:   0,
					OnColumn: 0,
				},
				VarName: ast.VarNameExp{
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 4,
					},
					Name: "a",
				},
				Expression: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "34",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 8,
					},
					Value: 34,
				},
			},
		},
		{ // test 1
			input: "let a = abc;",
			expLetSta: ast.LetSta{
				Token: token.Token{
					Literal:  "let",
					Type:     token.LET,
					OnLine:   0,
					OnColumn: 0,
				},
				VarName: ast.VarNameExp{
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 4,
					},
					Name: "a",
				},
				Expression: &ast.VarNameExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
					Name: "abc",
				},
			},
		},
		{ // test 2
			input: "let a[3] = 34;",
			expLetSta: ast.LetSta{
				Token: token.Token{
					Literal:  "let",
					Type:     token.LET,
					OnLine:   0,
					OnColumn: 0,
				},
				VarName: ast.VarNameExp{
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 4,
					},
					Name: "a",
				},
				Index: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "3",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 6,
					},
					Value: 3,
				},
				Expression: &ast.IntConstantExp{
					Token: token.Token{
						Literal:  "34",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: 34,
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		stmt := p.parseLetStatement()
		letSta, ok := stmt.(*ast.LetSta)
		if !ok {
			t.Fatalf("tests[%d]: stmt is not *ast.LetSta, got=%T", i, stmt)
		}

		if tt.expLetSta.Token != letSta.Token {
			t.Fatalf("letSta.Token is not %+v, got=%+v", tt.expLetSta, letSta.Token)
		}

		if tt.expLetSta.VarName.Token != letSta.VarName.Token {
			t.Fatalf("letSta.VarName.Token is not %+v, got=%+v",
				tt.expLetSta.VarName.Token, letSta.VarName.Token)
		}

		if tt.expLetSta.VarName.Name != letSta.VarName.Name {
			t.Fatalf("letSta.VarName.Name is not %s, got=%s",
				tt.expLetSta.VarName.Name, letSta.VarName.Name)
		}

		if !reflect.DeepEqual(letSta.Index, tt.expLetSta.Index) {
			t.Fatalf("letSta.Index is not %v, got=%v", tt.expLetSta.Index, letSta.Index)
		}

		if !reflect.DeepEqual(tt.expLetSta.Expression, letSta.Expression) {
			t.Fatalf("letSta.Expression is not %+T, got=%+T",
				tt.expLetSta.Expression, letSta.Expression)
		}
	}
}

func TestParseVarName(t *testing.T) {
	tests := []struct {
		input      string
		expVarName ast.VarNameExp
	}{
		{
			input: "ab_c",
			expVarName: ast.VarNameExp{
				Token: token.Token{
					Literal:  "ab_c",
					Type:     token.IDENT,
					OnLine:   0,
					OnColumn: 0,
				},
				Name: "ab_c",
			},
		},
		{
			input: "  _b9_c",
			expVarName: ast.VarNameExp{
				Token: token.Token{
					Literal:  "_b9_c",
					Type:     token.IDENT,
					OnLine:   0,
					OnColumn: 2,
				},
				Name: "_b9_c",
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseVarName()
		varName, ok := exp.(*ast.VarNameExp)
		if !ok {
			t.Fatalf("tests[%d]: exp is not *ast.VarNameExp, got=%T", i, exp)
		}

		if varName.Name != tt.expVarName.Name {
			t.Fatalf("tests[%d]: varName.Name is not %s, got=%s", i, tt.expVarName.Name, varName.Name)
		}

		if tt.expVarName.Token != varName.Token {
			t.Fatalf("varName.Token is not %+v, got=%+v", tt.expVarName.Token, varName.Token)
		}
	}
}

func TestParseInteger(t *testing.T) {
	tests := []struct {
		input         string
		expExpression ast.IntConstantExp
	}{
		{
			input: "345",
			expExpression: ast.IntConstantExp{
				Token: token.Token{
					Literal:  "345",
					Type:     token.INT,
					OnLine:   0,
					OnColumn: 0,
				},
				Value: 345,
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseExpression(LOWEST)
		intExp, ok := exp.(*ast.IntConstantExp)
		if !ok {
			t.Fatalf("tests[%d]: exp, is not *ast.IntConstantExp, got=%T", i, exp)
		}
		if intExp.Value != tt.expExpression.Value {
			t.Fatalf("tests[%d]: inExp.Value is not %d, got=%d", i, tt.expExpression.Value, intExp.Value)
		}

		if tt.expExpression.Token != intExp.Token {
			t.Fatalf("tests[%d]: intExp.Token is not %+v, got=%+v", i, tt.expExpression.Token, intExp.Token)
		}
	}
}

func TestVarDec(t *testing.T) {
	tests := []struct {
		input             string
		expToken          token.Token          // expected Token
		expDataType       token.Token          // expected dataType
		expIdentifierDecs []*ast.IdentifierDec // expected identifier declaration
	}{
		{ // test[0]
			input: "var int a, b,c;", // input
			expToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 0,
			},
			expDataType: token.Token{
				Literal:  "int",
				Type:     token.INT,
				OnLine:   0,
				OnColumn: 4,
			},
			expIdentifierDecs: []*ast.IdentifierDec{
				{ // identifier a
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
					Literal: "a",
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
					Literal: "b",
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
					Literal: "c",
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
			expToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 2,
			},
			expDataType: token.Token{
				Literal:  "Ball",
				Type:     token.IDENT,
				OnLine:   0,
				OnColumn: 6,
			},
			expIdentifierDecs: []*ast.IdentifierDec{
				{ // identifier a_b
					Token: token.Token{
						Literal:  "a_b",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Literal: "a_b",
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
					Literal: "__a",
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
					Literal: "ab9c",
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
			expToken: token.Token{ // Token
				Literal:  "var",
				Type:     token.VAR,
				OnLine:   0,
				OnColumn: 2,
			},
			expDataType: token.Token{
				Literal:  "Ball",
				Type:     token.IDENT,
				OnLine:   0,
				OnColumn: 6,
			},
			expIdentifierDecs: []*ast.IdentifierDec{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Literal: "foo_bar",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
				},
			},
		},
		{ // test[3]
			input: "  static char foo_bar;", // input
			expToken: token.Token{ // Token
				Literal:  "static",
				Type:     token.STATIC,
				OnLine:   0,
				OnColumn: 2,
			},
			expDataType: token.Token{
				Literal:  "char",
				Type:     token.CHAR,
				OnLine:   0,
				OnColumn: 9,
			},
			expIdentifierDecs: []*ast.IdentifierDec{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 14,
					},
					Literal: "foo_bar",
					DataType: token.Token{
						Literal:  "char",
						Type:     token.CHAR,
						OnLine:   0,
						OnColumn: 9,
					},
				},
			},
		},
		{ // test[4]
			input: "  field Ball foo_bar;", // input
			expToken: token.Token{ // Token
				Literal:  "field",
				Type:     token.FIELD,
				OnLine:   0,
				OnColumn: 2,
			},
			expDataType: token.Token{
				Literal:  "Ball",
				Type:     token.IDENT,
				OnLine:   0,
				OnColumn: 8,
			},
			expIdentifierDecs: []*ast.IdentifierDec{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 13,
					},
					Literal: "foo_bar",
					DataType: token.Token{
						Literal:  "Ball",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		varDec := p.parseVarDec()
		if p.errors.QueueSize() != 0 {
			p.ReportErrors()
			t.Fatal("terminating testing")
		}

		if varDec == nil {
			t.Fatalf("test[%d]: varDec is nil", i)
		}

		if tt.expToken != varDec.Token {
			t.Fatalf("test[%d]: varDec.Token is not %v, got=%v", i, tt.expToken, varDec.Token)
		}

		if tt.expDataType != varDec.DataType {
			t.Fatalf("test[%d]: varDec.DataType is not %v, got=%v",
				i, tt.expDataType, varDec.DataType)
		}

		for j, expectedIdentifer := range tt.expIdentifierDecs {
			if !testIdentifierDec(t, expectedIdentifer, varDec.IdentifierDecs[j]) {
				t.Fatalf("test[%d]: varDec.IdentifierDecs[%d] is not %v, got=%v",
					i, j, expectedIdentifer, varDec.IdentifierDecs[j])
			}
		}
	}
}

func testIdentifierDec(t *testing.T, ref *ast.IdentifierDec, tok *ast.IdentifierDec) bool {
	if ref.Token != tok.Token {
		t.Errorf("tok.Token is not %v, got=%v", ref.Token, tok.Token)
		return false
	}

	if ref.Literal != tok.Literal {
		t.Errorf("tok.Literal is not %s, got=%s", ref.Literal, tok.Literal)
		return false
	}

	if ref.DataType != tok.DataType {
		t.Errorf("tok.Type is not %v, got=%v", ref.DataType, tok.DataType)
		return false
	}
	return true
}
