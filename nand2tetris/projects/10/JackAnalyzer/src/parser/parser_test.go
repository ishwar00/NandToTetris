package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ishwar00/JackAnalyzer/ast"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

func TestArithmeticExp(t *testing.T) {
	tests := []struct {
		input         string
		expExpression string
	}{
		{
			"-a * b",
			"((-a)*b)",
		},
		{
			"~-a",
			"(~(-a))",
		},
		{
			"a + b + c",
			"((a+b)+c)",
		},
		{
			"a + b - c",
			"((a+b)-c)",
		},
		{
			"a * b * c",
			"((a*b)*c)",
		},
		{
			"a * b / c",
			"((a*b)/c)",
		},
		{
			"a + b / c",
			"(a+(b/c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a+(b*c))+(d/e))-f)",
		},
		{
			"5>4 = 3<4",
			"((5>4)=(3<4))",
		},
		{
			"3 + 4 * 5 = 3 * 1 + 4 * 5",
			"((3+(4*5))=((3*1)+(4*5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3>5 = false",
			"((3>5)=false)",
		},
		{
			"3 < 5 = true",
			"((3<5)=true)",
		},
		{
			"1 + (2 +3) + 4",
			"((1+(2+3))+4)",
		},
		{
			"(5 + 5) * 2",
			"((5+5)*2)",
		},
		{
			"2 / (5 + 5)",
			"(2/(5+5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5+5)*2)*(5+5))",
		},
		{
			"-(5 + 5)",
			"(-(5+5))",
		},
		{
			"~(true = true)",
			"(~(true=true))",
		},
		{
			"a + add(b * c) + d",
			"((a+(add((b*c))))+d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"(add(a,b,1,(2*3),(4+5),(add(6,(7*8)))))",
		},
		{
			"add(a + b + c * d / f + g)",
			"(add((((a+b)+((c*d)/f))+g)))",
		},
		{
			"add(arr[0], \"str\", varName, call(A[index]))",
			"(add((arr[0]),str,varName,(call((A[index])))))",
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseExpression(LOWEST)
		if exp.String() != tt.expExpression {
			t.Fatalf("tests[%d]: exp is not %+v, got=%+v", i, tt.expExpression, exp.String())
		}
	}
}

func TestGroupExp(t *testing.T) {
	tests := []struct {
		input         string
		expExpression ast.Expression
	}{
		{
			input: "(hey)",
			expExpression: &ast.IdentifierExp{
				Token: token.Token{
					Literal:  "hey",
					Type:     token.IDENT,
					OnColumn: 1,
				},
				Value: "hey",
			},
		},
		{
			input: "((hey))",
			expExpression: &ast.IdentifierExp{
				Token: token.Token{
					Literal:  "hey",
					Type:     token.IDENT,
					OnColumn: 2,
				},
				Value: "hey",
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)

		exp := p.parseExpression(LOWEST)

		if exp == nil || reflect.ValueOf(exp).IsZero() {
			t.Fatalf("tests[%d]: exp is nil", i)
		}

		if !reflect.DeepEqual(exp, tt.expExpression) {
			t.Fatalf("tests[%d]: exp is not %+v, got=%+v", i, exp, tt.expExpression)
		}
	}
}

func TestMethodCall(t *testing.T) {
	tests := []struct {
		input    string
		expInfix ast.Expression
	}{
		{
			input: "varName.mCall()",
			expInfix: &ast.InfixExp{
				Token: token.Token{
					Literal:  "(",
					Type:     token.LPAREN,
					OnColumn: 13,
				},
				Left: &ast.InfixExp{
					Token: token.Token{
						Literal:  ".",
						Type:     token.PERIOD,
						OnColumn: 7,
					},
					Left: &ast.IdentifierExp{
						Token: token.Token{
							Literal: "varName",
							Type:    token.IDENT,
						},
						Value: "varName",
					},
					Operator: ".",
					Right: &ast.IdentifierExp{
						Token: token.Token{
							Literal:  "mCall",
							Type:     token.IDENT,
							OnColumn: 8,
						},
						Value: "mCall",
					},
				},
				Operator: "(",
			},
		},
		{
			input: "foo.bar(4, hey(hi, 9))",
			expInfix: &ast.InfixExp{
				Token: token.Token{
					Literal:  "(",
					Type:     token.LPAREN,
					OnColumn: 7,
				},
				Left: &ast.InfixExp{
					Token: token.Token{
						Literal:  ".",
						Type:     token.PERIOD,
						OnColumn: 3,
					},
					Left: &ast.IdentifierExp{
						Token: token.Token{
							Literal: "foo",
							Type:    token.IDENT,
						},
						Value: "foo",
					},
					Operator: ".",
					Right: &ast.IdentifierExp{
						Token: token.Token{
							Literal:  "bar",
							Type:     token.IDENT,
							OnColumn: 4,
						},
						Value: "bar",
					},
				},
				Operator: "(",
				Right: &ast.ExpressionListExp{
					Token: token.Token{
						Literal:  "4",
						Type:     token.INT,
						OnColumn: 8,
					},
					Expressions: []ast.Expression{
						&ast.IntConstExp{
							Token: token.Token{
								Literal:  "4",
								Type:     token.INT,
								OnColumn: 8,
							},
							Value: 4,
						},
						&ast.InfixExp{
							Token: token.Token{
								Literal:  "(",
								Type:     token.LPAREN,
								OnColumn: 14,
							},
							Left: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "hey",
									Type:     token.IDENT,
									OnColumn: 11,
								},
								Value: "hey",
							},
							Operator: "(",
							Right: &ast.ExpressionListExp{
								Token: token.Token{
									Literal:  "hi",
									Type:     token.IDENT,
									OnColumn: 15,
								},
								Expressions: []ast.Expression{
									&ast.IdentifierExp{
										Token: token.Token{
											Literal:  "hi",
											Type:     token.IDENT,
											OnColumn: 15,
										},
										Value: "hi",
									},
									&ast.IntConstExp{
										Token: token.Token{
											Literal:  "9",
											Type:     token.INT,
											OnColumn: 19,
										},
										Value: 9,
									},
								},
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
		exp := p.parseExpression(LOWEST)

		if exp == nil || reflect.ValueOf(exp).IsZero() {
			t.Fatalf("tests[%d]: exp is nil", i)
		}

		if !reflect.DeepEqual(exp, tt.expInfix) {
			fmt.Printf("last %+v\n", p.curToken)
			t.Fatalf("tests[%d]: exp is not %+v, got=%+v", i, tt.expInfix, exp)
		}
	}
}

func TestSubroutineCall(t *testing.T) {
	tests := []struct {
		input            string
		expSubroutineExp *ast.InfixExp
	}{
		{
			input: "subName()",
			expSubroutineExp: &ast.InfixExp{
				Token: token.Token{
					Literal:  "(",
					Type:     token.LPAREN,
					OnColumn: 7,
				},
				Operator: "(",
				Left: &ast.IdentifierExp{
					Token: token.Token{
						Literal: "subName",
						Type:    token.IDENT,
					},
					Value: "subName",
				},
			},
		},
		{
			input: "Subroutine(3, 45, asd)",
			expSubroutineExp: &ast.InfixExp{
				Token: token.Token{
					Literal:  "(",
					Type:     token.LPAREN,
					OnColumn: 10,
				},
				Operator: "(",
				Left: &ast.IdentifierExp{
					Token: token.Token{
						Literal: "Subroutine",
						Type:    token.IDENT,
					},
					Value: "Subroutine",
				},
				Right: &ast.ExpressionListExp{
					Token: token.Token{
						Literal:  "3",
						Type:     token.INT,
						OnColumn: 11,
					},
					Expressions: []ast.Expression{
						&ast.IntConstExp{
							Token: token.Token{
								Literal:  "3",
								Type:     token.INT,
								OnColumn: 11,
							},
							Value: 3,
						},
						&ast.IntConstExp{
							Token: token.Token{
								Literal:  "45",
								Type:     token.INT,
								OnColumn: 14,
							},
							Value: 45,
						},
						&ast.IdentifierExp{
							Token: token.Token{
								Literal:  "asd",
								Type:     token.IDENT,
								OnColumn: 18,
							},
							Value: "asd",
						},
					},
				},
			},
		},
		{
			input: "subName(\"sdf\")",
			expSubroutineExp: &ast.InfixExp{
				Token: token.Token{
					Literal:  "(",
					Type:     token.LPAREN,
					OnColumn: 7,
				},
				Left: &ast.IdentifierExp{
					Token: token.Token{
						Literal: "subName",
						Type:    token.IDENT,
					},
					Value: "subName",
				},
				Operator: "(",
				Right: &ast.ExpressionListExp{
					Token: token.Token{
						Literal:  "sdf",
						Type:     token.STR_CONST,
						OnColumn: 9,
					},
					Expressions: []ast.Expression{
						&ast.StrConstExp{
							Token: token.Token{
								Literal:  "sdf",
								Type:     token.STR_CONST,
								OnColumn: 9,
							},
							Value: "sdf",
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseExpression(LOWEST)

		infixExp, ok := exp.(*ast.InfixExp)
		if !ok {
			t.Fatalf("exp is not *ast.InfixExp, got=%T", exp)
		}

		if infixExp.Token != tt.expSubroutineExp.Token {
			t.Fatalf("infixExp.Token is not %+v, got=%+v",
				tt.expSubroutineExp.Token, infixExp.Token)
		}

		if infixExp.Operator != tt.expSubroutineExp.Operator {
			t.Fatalf("infixExp.Operator is not %+v, got=%+v",
				tt.expSubroutineExp.Operator, infixExp.Operator)
		}

		if !reflect.DeepEqual(infixExp.Left, tt.expSubroutineExp.Left) {
			t.Fatalf("tests[%d]: infixExp.Left is not %+v, got=%+v",
				i, tt.expSubroutineExp.Left, infixExp.Left)
		}

		if !reflect.DeepEqual(tt.expSubroutineExp.Right, infixExp.Right) {
			t.Fatalf("tests[%d]: infixExp is not %+v, got=%+v",
				i, tt.expSubroutineExp, infixExp.Right)
		}

		// expList, ok := tt.expSubroutineExp.Right.(*ast.ExpressionListExp)
		// if !ok {
		//     t.Fatalf("tt.expSubroutineExp.Right is not *ast.ExpressionListExp, got=%T",
		//     tt.expSubroutineExp.Right)
		// }

		// list, ok := infixExp.Right.(*ast.ExpressionListExp)
		// if !ok {
		//     t.Fatalf("tests[%d]: infix.Right is not *ast.ExpressionListExp, got=%T", i, infixExp.Right)
		// }

		// for j, jthExp := range expList.Expressions {
		//     if !reflect.DeepEqual(jthExp, list.Expressions[j]) {
		//         t.Fatalf("tests[%d][%d]: exp is not %+v, got=%+v", i, j, tt.expSubroutineExp, exp)
		//     }
		// }
	}
}

func TestExpressionList(t *testing.T) {
	tests := []struct {
		input      string
		expExpList ast.ExpressionListExp
	}{
		{
			input: "2",
			expExpList: ast.ExpressionListExp{
				Token: token.Token{
					Literal: "2",
					Type:    token.INT,
				},
				Expressions: []ast.Expression{
					&ast.IntConstExp{
						Token: token.Token{
							Literal: "2",
							Type:    token.INT,
						},
						Value: 2,
					},
				},
			},
		},
		{
			input: "2, 3, 4",
			expExpList: ast.ExpressionListExp{
				Token: token.Token{
					Literal: "2",
					Type:    token.INT,
				},
				Expressions: []ast.Expression{
					&ast.IntConstExp{
						Token: token.Token{
							Literal: "2",
							Type:    token.INT,
						},
						Value: 2,
					},
					&ast.IntConstExp{
						Token: token.Token{
							Literal:  "3",
							Type:     token.INT,
							OnColumn: 3,
						},
						Value: 3,
					},
					&ast.IntConstExp{
						Token: token.Token{
							Literal:  "4",
							Type:     token.INT,
							OnColumn: 6,
						},
						Value: 4,
					},
				},
			},
		},
		{
			input: "a, 3, \"d__f\"",
			expExpList: ast.ExpressionListExp{
				Token: token.Token{
					Literal: "a",
					Type:    token.IDENT,
				},
				Expressions: []ast.Expression{
					&ast.IdentifierExp{
						Token: token.Token{
							Literal: "a",
							Type:    token.IDENT,
						},
						Value: "a",
					},
					&ast.IntConstExp{
						Token: token.Token{
							Literal:  "3",
							Type:     token.INT,
							OnColumn: 3,
						},
						Value: 3,
					},
					&ast.StrConstExp{
						Token: token.Token{
							Literal:  "d__f",
							Type:     token.STR_CONST,
							OnColumn: 7,
						},
						Value: "d__f",
					},
				},
			},
		},
	}
	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseExpressionList()
		expList, ok := exp.(*ast.ExpressionListExp)
		if !ok {
			t.Fatalf("tests[%d]: exp is not *ast.ExpressionListExp, got=%T", i, exp)
		}

		if expList.Token != tt.expExpList.Token {
			t.Fatalf("expList.Token is not %+v, got=%+v",
				tt.expExpList.Token, expList.Token)
		}

		for j, exp := range expList.Expressions {
			if !reflect.DeepEqual(exp, tt.expExpList.Expressions[j]) {
				t.Fatalf("tests[%d]: exp[%d], is not %+v, got=%+v",
					i, j, tt.expExpList.Expressions[j], exp)
			}
		}
	}
}

func TestArrayIndex(t *testing.T) {
	tests := []struct {
		input         string
		expArrayIndex ast.InfixExp
	}{
		{
			input: "arr[343]",
			expArrayIndex: ast.InfixExp{
				Token: token.Token{
					Literal:  "[",
					Type:     token.LBRACK,
					OnLine:   0,
					OnColumn: 3,
				},
				Operator: "[",
				Left: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "arr",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 0,
					},
					Value: "arr",
				},
				Right: &ast.IntConstExp{
					Token: token.Token{
						Literal:  "343",
						Type:     token.INT,
						OnLine:   0,
						OnColumn: 4,
					},
					Value: 343,
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)

		exp := p.parseExpression(LOWEST)
		arrExp, ok := exp.(*ast.InfixExp)
		if !ok {
			t.Fatalf("tests[%d]: arrExp is not *ast.InfixExp, got=%T", i, exp)
		}

		if arrExp.Token != tt.expArrayIndex.Token {
			t.Fatalf("tests[%d]: arrExp.Token is not %+v, got=%+v",
				i, tt.expArrayIndex.Token, arrExp.Token)
		}

		if arrExp.Operator != tt.expArrayIndex.Operator {
			t.Fatalf("tests[%d]: arrExp.Operator is not %+v, got=%+v",
				i, tt.expArrayIndex.Operator, arrExp.Operator)
		}

		if !reflect.DeepEqual(arrExp.Left, tt.expArrayIndex.Left) {
			t.Fatalf("tests[%d]: arrExp.Left is not %+v, got=%+v",
				i, tt.expArrayIndex.Left, arrExp.Left)
		}

		if !reflect.DeepEqual(arrExp.Right, tt.expArrayIndex.Right) {
			t.Fatalf("tests[%d]: arrExp.Index is not %+v, got=%+v",
				i, tt.expArrayIndex.Right, arrExp.Right)
		}
	}
}

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
				Expression: &ast.IntConstExp{
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
				Condition: &ast.IntConstExp{
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
				Condition: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 6,
					},
					Value: "abc",
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
							Name: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "deep",
									Type:     token.IDENT,
									OnColumn: 17,
								},
								Value: "deep",
							},
							Expression: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "source",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 24,
								},
								Value: "source",
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
						Name: &ast.IdentifierExp{
							Token: token.Token{
								Literal:  "a",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 6,
							},
							Value: "a",
						},
						Expression: &ast.IntConstExp{
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
						Name: &ast.IdentifierExp{
							Token: token.Token{
								Literal:  "a_b",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 19,
							},
							Value: "a_b",
						},
						Expression: &ast.IdentifierExp{
							Token: token.Token{
								Literal:  "abc",
								Type:     token.IDENT,
								OnLine:   0,
								OnColumn: 25,
							},
							Value: "abc",
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
				Condition: &ast.IntConstExp{
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
				Condition: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 3,
					},
					Value: "abc",
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
							Name: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "hey",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 13,
								},
								Value: "hey",
							},
							Expression: &ast.IntConstExp{
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
				Condition: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 3,
					},
					Value: "abc",
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
							Name: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "hey",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 13,
								},
								Value: "hey",
							},
							Expression: &ast.IntConstExp{
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
							Name: &ast.IdentifierExp{
								Token: token.Token{
									Literal:  "yo",
									Type:     token.IDENT,
									OnLine:   0,
									OnColumn: 36,
								},
								Value: "yo",
							},
							Expression: &ast.IntConstExp{
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
				Name: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 4,
					},
					Value: "a",
				},
				Expression: &ast.IntConstExp{
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
				Name: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 4,
					},
					Value: "a",
				},
				Expression: &ast.IdentifierExp{
					Token: token.Token{
						Literal:  "abc",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
					Value: "abc",
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
				Name: &ast.InfixExp{
					Token: token.Token{
						Literal:  "[",
						Type:     token.LBRACK,
						OnColumn: 5,
					},
					Operator: "[",
					Left: &ast.IdentifierExp{
						Token: token.Token{
							Literal:  "a",
							Type:     token.IDENT,
							OnColumn: 4,
						},
						Value: "a",
					},
					Right: &ast.IntConstExp{
						Token: token.Token{
							Literal:  "3",
							Type:     token.INT,
							OnColumn: 6,
						},
						Value: 3,
					},
				},
				Expression: &ast.IntConstExp{
					Token: token.Token{
						Literal:  "34",
						Type:     token.INT,
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
		letSta := p.parseLetStatement()

		if letSta == nil {
			t.Fatalf("tests[%d]: letSta is nil", i)
		}

		if tt.expLetSta.Token != letSta.Token {
			t.Fatalf("tests[%d]: letSta.Token is not %+v, got=%+v",
				i, tt.expLetSta, letSta.Token)
		}

		// v := reflect.ValueOf
		// v := reflect.TypeOf

		if !reflect.DeepEqual(tt.expLetSta.Name, letSta.Name) {
			t.Fatalf("tests[%d]:letSta.Name is not %+v, got=%+v",
				i, tt.expLetSta.Name.GetToken(), letSta.Name.GetToken())
		}

		if !reflect.DeepEqual(tt.expLetSta.Expression, letSta.Expression) {
			t.Fatalf("letSta.Expression is not %s, got=%s",
				tt.expLetSta.Expression, letSta.Expression)
		}
	}
}

func TestParseVarName(t *testing.T) {
	tests := []struct {
		input      string
		expVarName ast.IdentifierExp
	}{
		{
			input: "ab_c",
			expVarName: ast.IdentifierExp{
				Token: token.Token{
					Literal:  "ab_c",
					Type:     token.IDENT,
					OnLine:   0,
					OnColumn: 0,
				},
				Value: "ab_c",
			},
		},
		{
			input: "  _b9_c",
			expVarName: ast.IdentifierExp{
				Token: token.Token{
					Literal:  "_b9_c",
					Type:     token.IDENT,
					OnLine:   0,
					OnColumn: 2,
				},
				Value: "_b9_c",
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseIdentifierExp()
		varName, ok := exp.(*ast.IdentifierExp)
		if !ok {
			t.Fatalf("tests[%d]: exp is not *ast.IdentifierExp, got=%T", i, exp)
		}

		if varName.Value != tt.expVarName.Value {
			t.Fatalf("tests[%d]: varName.Value is not %s, got=%s",
				i, tt.expVarName.Value, varName.Value)
		}

		if tt.expVarName.Token != varName.Token {
			t.Fatalf("varName.Token is not %+v, got=%+v", tt.expVarName.Token, varName.Token)
		}
	}
}

func TestParseConstant(t *testing.T) {
	tests := []struct {
		input         string
		expExpression ast.Expression
	}{
		{
			input: "345",
			expExpression: &ast.IntConstExp{
				Token: token.Token{
					Literal:  "345",
					Type:     token.INT,
					OnLine:   0,
					OnColumn: 0,
				},
				Value: 345,
			},
		},
		{
			input: "\"strconstant\"",
			expExpression: &ast.StrConstExp{
				Token: token.Token{
					Literal:  "strconstant",
					Type:     token.STR_CONST,
					OnLine:   0,
					OnColumn: 1,
				},
				Value: "strconstant",
			},
		},
		{
			input: "this",
			expExpression: &ast.KeywordConstExp{
				Token: token.Token{
					Literal:  "this",
					Type:     token.THIS,
					OnLine:   0,
					OnColumn: 0,
				},
				Value: "this",
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		exp := p.parseExpression(LOWEST)
		if !reflect.DeepEqual(exp, tt.expExpression) {
			t.Fatalf("tests[%d]: exp is not %+v, got=%+v", i, exp, tt.expExpression)
		}
	}
}

func TestVarDec(t *testing.T) {
	tests := []struct {
		input             string
		expToken          token.Token          // expected Token
		expDataType       token.Token          // expected dataType
		expIdentifierExps []*ast.IdentifierExp // expected identifier declaration
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
			expIdentifierExps: []*ast.IdentifierExp{
				{ // identifier a
					Token: token.Token{
						Literal:  "a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 8,
					},
					Value: "a",
				},
				{ // identifier b
					Token: token.Token{
						Literal:  "b",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "b",
				},
				{ // identifier c
					Token: token.Token{
						Literal:  "c",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 13,
					},
					Value: "c",
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
			expIdentifierExps: []*ast.IdentifierExp{
				{ // identifier a_b
					Token: token.Token{
						Literal:  "a_b",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "a_b",
				},
				{ // identifier __a
					Token: token.Token{
						Literal:  "__a",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 16,
					},
					Value: "__a",
				},
				{ // identifier ab9c
					Token: token.Token{
						Literal:  "ab9c",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 20,
					},
					Value: "ab9c",
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
			expIdentifierExps: []*ast.IdentifierExp{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 11,
					},
					Value: "foo_bar",
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
			expIdentifierExps: []*ast.IdentifierExp{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 14,
					},
					Value: "foo_bar",
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
			expIdentifierExps: []*ast.IdentifierExp{
				{ // identifier foo_bar
					Token: token.Token{
						Literal:  "foo_bar",
						Type:     token.IDENT,
						OnLine:   0,
						OnColumn: 13,
					},
					Value: "foo_bar",
				},
			},
		},
	}

	for i, tt := range tests {
		l := lexer.LexString(tt.input)
		p := New(l)
		varDec := p.parseVarDec()

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

		for j, expectedIdentifer := range tt.expIdentifierExps {
			if !reflect.DeepEqual(expectedIdentifer, varDec.IdentifierExps[j]) {
				t.Fatalf("test[%d]: varDec.IdentifierExps[%d] is not %v, got=%v",
					i, j, expectedIdentifer, varDec.IdentifierExps[j])
			}
		}
	}
}
