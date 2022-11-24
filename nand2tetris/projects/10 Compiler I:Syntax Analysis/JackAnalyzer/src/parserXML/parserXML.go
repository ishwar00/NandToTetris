package parserxml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ishwar00/JackAnalyzer/ast"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/parser"
	"github.com/ishwar00/JackAnalyzer/token"
)

// implementation is not pretty and neat.
func ParseIntoXML(jackFilePath string) {
	l, err := lexer.LexFile(jackFilePath)
	if err != nil {
		panic(err)
	}
	if l.FoundErrors() {
		l.ReportErrors()
	}
	p := parser.New(l)
	AST := p.ParseClassDec()
	if p.HasErrors() {
		p.ReportErrors()
		return
	}

	if AST == nil {
		panic("ParseIntoXML: AST is nil")
	}
	file, closer, err := createOutputFile(jackFilePath)
	if err != nil {
		panic(err)
	}
	defer closer()
	// walkClass(AST, os.Stdout)
	walkClass(AST, file)
}

func writeTok(tok token.Token, out io.Writer) {
	tag := tagOf(tok.Type)
	switch tok.Type {
	case token.LT:
		tok.Literal = fmt.Sprintf("&lt;")
	case token.GT:
		tok.Literal = fmt.Sprintf("&gt;")
	case token.AMPERS:
		tok.Literal = fmt.Sprintf("&amp;")
	}
	buf := fmt.Sprintf("<%s>%s</%s>\n", tag, tok.Literal, tag)
	out.Write([]byte(buf))
}

func tagOf(tokenType token.TokenType) string {
	value := int(tokenType)
	// keyword range
	if 0 <= value && value <= 20 {
		value = token.KEYWORDS
	}

	// symbol range
	if 21 <= value && value <= 39 {
		value = token.SYMBOLS
	}

	switch value {
	case token.KEYWORDS:
		return "keyword"
	case token.SYMBOLS:
		return "symbol"
	case token.IDENT:
		return "identifier"
	case token.INT_CONST:
		return "integerConstant"
	case token.STR_CONST:
		return "stringConstant"
	default:
		return "illegal"
	}
}

func walkClass(class *ast.ClassDec, out io.Writer) {
	out.Write([]byte("<class>\n"))
	writeTok(class.Token, out)
	writeTok(token.Token{Literal: class.Name, Type: token.IDENT}, out)
	writeTok(token.Token{Literal: "{", Type: token.LBRACE}, out)

	for _, varDec := range class.ClassVarDecs {
		walkVarDec(varDec, out)
	}

	for _, subDec := range class.Subroutines {
		walkSubDec(subDec, out)
	}

	writeTok(token.Token{Literal: "}", Type: token.RBRACE}, out)
	out.Write([]byte("</class>\n"))
}

func walkVarDec(varDec *ast.VarDec, out io.Writer) {
	tag := "<classVarDec>\n"
	if varDec.Token.Literal == "var" {
		tag = "<varDec>\n"
	}
	out.Write([]byte(tag))

	writeTok(varDec.Token, out)    // var|static|field
	writeTok(varDec.DataType, out) // int, boolean, char, Ball etc

	for i, identifier := range varDec.IdentifierExps {
		writeTok(identifier.Token, out)
		if i+1 < len(varDec.IdentifierExps) {
			writeTok(token.Token{Literal: ",", Type: token.COMMA}, out)
		}
	}
	writeTok(token.Token{Literal: ";", Type: token.SEMICO}, out)

	tag = "</classVarDec>\n"
	if varDec.Token.Literal == "var" {
		tag = "</varDec>\n"
	}
	out.Write([]byte(tag))
}

func walkSubDec(subroutine *ast.SubroutineDec, out io.Writer) {
	out.Write([]byte("<subroutineDec>\n"))

	writeTok(subroutine.Token, out) // var|static|field
	writeTok(subroutine.ReturnType, out)
	writeTok(subroutine.SubName.Token, out)

	writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
	walkParameterList(subroutine.Parameters, out)
	writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)

	walkSubroutineBody(subroutine.Body, out)

	out.Write([]byte("</subroutineDec>\n"))
}

func walkParameterList(parameters []*ast.ParameterDec, out io.Writer) {
	out.Write([]byte("<parameterList>\n"))

	for i, param := range parameters {
		writeTok(param.Token, out)
		writeTok(param.Identifier.Token, out)
		if i+1 < len(parameters) {
			writeTok(token.Token{Literal: ",", Type: token.COMMA}, out)
		}
	}

	out.Write([]byte("</parameterList>\n"))
}

func walkSubroutineBody(body *ast.SubroutineBodyDec, out io.Writer) {
	out.Write([]byte("<subroutineBody>\n"))

	writeTok(token.Token{Literal: "{", Type: token.LBRACE}, out)

	for _, varDec := range body.VarDecs {
		walkVarDec(varDec, out)
	}

	walkStatements(body.Statements, out)

	writeTok(token.Token{Literal: "}", Type: token.RBRACE}, out)

	out.Write([]byte("</subroutineBody>\n"))
}

func walkStatements(stmts []ast.Statement, out io.Writer) {
	out.Write([]byte("<statements>\n"))

	for _, stmt := range stmts {
		switch v := stmt.(type) {
		case *ast.LetSta:
			walkLetSta(v, out)
		case *ast.IfElseSta:
			walkIfElseSta(v, out)
		case *ast.WhileSta:
			walkWhileSta(v, out)
		case *ast.ReturnSta:
			walkReturnSta(v, out)
		case *ast.DoSta:
			walkDoSta(v, out)
		default:
		}
	}

	out.Write([]byte("</statements>\n"))
}

func walkDoSta(do *ast.DoSta, out io.Writer) {
	out.Write([]byte("<doStatement>\n"))
	writeTok(do.Token, out)
	walkDoSubroutine(do.SubCall, out)
	writeTok(token.Token{Literal: ";", Type: token.SEMICO}, out)
	out.Write([]byte("</doStatement>\n"))
}

func walkDoSubroutine(exp ast.Expression, out io.Writer) {
	switch v := exp.(type) {
	case *ast.InfixExp:
		if v.Token.Literal == "(" {
			walkDoSubroutine(v.Left, out)
			writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
			if v.Right != nil {
				walkExpressionList(v.Right.(*ast.ExpressionListExp), out)
			} else {
				out.Write([]byte("<expressionList>\n"))
				out.Write([]byte("</expressionList>\n"))
			}
			writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)
		} else {
			walkDoSubroutine(v.Left, out)
			writeTok(token.Token{Literal: ".", Type: token.PERIOD}, out)
			writeTok(v.Right.GetToken(), out)
		}
	case *ast.IdentifierExp:
		writeTok(v.Token, out)
	default:
		errMsg := fmt.Sprintf("walkDoSubrutine got something unexpected exp %+v", exp)
		panic(errMsg)
	}
}

func walkWhileSta(whileSta *ast.WhileSta, out io.Writer) {
	out.Write([]byte("<whileStatement>\n"))

	writeTok(whileSta.Token, out)
	writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
	out.Write([]byte("<expression>\n"))
	walkExpression(whileSta.Condition, out)
	out.Write([]byte("</expression>\n"))
	writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)
	writeTok(token.Token{Literal: "{", Type: token.LBRACE}, out)
	if whileSta.Statements != nil {
		walkStatements(whileSta.Statements.Statements, out)
	}
	writeTok(token.Token{Literal: "}", Type: token.RBRACE}, out)

	out.Write([]byte("</whileStatement>\n"))
}

func walkReturnSta(returnSta *ast.ReturnSta, out io.Writer) {
	out.Write([]byte("<returnStatement>\n"))
	writeTok(returnSta.Token, out)
	if returnSta.Expression != nil {
		out.Write([]byte("<expression>\n"))
		walkExpression(returnSta.Expression, out)
		out.Write([]byte("</expression>\n"))
	}
	writeTok(token.Token{Literal: ";", Type: token.SEMICO}, out)
	out.Write([]byte("</returnStatement>\n"))
}

func walkIfElseSta(ifElse *ast.IfElseSta, out io.Writer) {
	out.Write([]byte("<ifStatement>\n"))

	writeTok(ifElse.Token, out) // if
	writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
	out.Write([]byte("<expression>\n"))
	walkExpression(ifElse.Condition, out) // condition
	out.Write([]byte("</expression>\n"))
	writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)
	writeTok(token.Token{Literal: "{", Type: token.LBRACE}, out)
	walkStatements(ifElse.Then.Statements, out)
	writeTok(token.Token{Literal: "}", Type: token.RBRACE}, out)

	if ifElse.Else != nil {
		writeTok(token.Token{Literal: "else", Type: token.ELSE}, out)
		writeTok(token.Token{Literal: "{", Type: token.LBRACE}, out)
		walkStatements(ifElse.Else.Statements, out)
		writeTok(token.Token{Literal: "}", Type: token.RBRACE}, out)
	}

	out.Write([]byte("</ifStatement>\n"))
}

func walkLetSta(stmt *ast.LetSta, out io.Writer) {
	out.Write([]byte("<letStatement>\n"))

	writeTok(stmt.Token, out) // let
	// walkExpression(stmt.Name, out)
	walkName(stmt.Name, out)
	writeTok(token.Token{Literal: "=", Type: token.EQ}, out)
	out.Write([]byte("<expression>\n"))
	walkExpression(stmt.Expression, out)
	out.Write([]byte("</expression>\n"))
	writeTok(token.Token{Literal: ";", Type: token.SEMICO}, out)

	out.Write([]byte("</letStatement>\n"))
}

func walkName(name ast.Expression, out io.Writer) {
	if name.GetToken().Literal == "[" {
		infix := name.(*ast.InfixExp)
		writeTok(infix.Left.GetToken(), out) // name
		writeTok(name.GetToken(), out)       // [
		out.Write([]byte("<expression>\n"))
		walkExpression(infix.Right, out)
		out.Write([]byte("</expression>\n"))
		writeTok(token.Token{Literal: "]", Type: token.RBRACK}, out)
	} else {
		writeTok(name.GetToken(), out) // this has to be name
	}
}

func walkExpression(exp ast.Expression, out io.Writer) {
	writeTerm := func(tok token.Token) {
		out.Write([]byte("<term>\n"))
		writeTok(tok, out)
		out.Write([]byte("</term>\n"))
	}

	switch v := exp.(type) {
	case *ast.PrefixExp:
		walkPrefixExp(v, out)
	case *ast.InfixExp:
		walkInfixExp(v, out)
	case *ast.IntConstExp:
		writeTerm(v.Token)
	case *ast.StrConstExp:
		writeTerm(v.Token)
	case *ast.KeywordConstExp:
		writeTerm(v.Token)
	case *ast.IdentifierExp:
		writeTerm(v.Token)
	case *ast.ExpressionListExp:
		walkExpressionList(v, out)
	default:
		errMsg := fmt.Sprintf(">>>> missing implementation for %T", exp)
		panic(errMsg)
	}
}

//        root
//    >= /    \ <
//   left      right
// it will return true if root operator's precedence is less or equal
// to left's operator precedence
func checkLeftPrecedence(left ast.Expression, root token.Token) bool {
	rootPrec := parser.Precedences[root.Type]
	switch v := left.(type) {
	case *ast.PrefixExp:
		return parser.PREFIX >= rootPrec
	case *ast.InfixExp:
		return parser.Precedences[v.Token.Type] >= rootPrec
	default:
		return true
	}
}

//        root
//    >= /    \ <
//   left      right
// it will return true if root operator's precedenc is less than right's
// operator's precedence
func checkRightPrecedence(root token.Token, right ast.Expression) bool {
	rootPrec := parser.Precedences[root.Type]
	switch v := right.(type) {
	case *ast.PrefixExp:
		return parser.PREFIX > rootPrec
	case *ast.InfixExp:
		return parser.Precedences[v.Token.Type] > rootPrec
	default:
		// int, constants etc
		return true
	}
}

func walkInfixExp(exp *ast.InfixExp, out io.Writer) {
	mp := map[string]token.Token{
		"(": {Literal: ")", Type: token.RPAREN},
		"[": {Literal: "]", Type: token.RBRACK},
	}

	leftClose := !checkLeftPrecedence(exp.Left, exp.Token)
	rightClose := !checkRightPrecedence(exp.Token, exp.Right)
	c, ok := mp[exp.Operator]
	if ok {
		out.Write([]byte("<term>\n"))
	}

	methodCall := exp.Operator == "."
	if methodCall && exp.Left.GetToken().Literal != "." {
		writeTok(exp.Left.GetToken(), out)
	} else if exp.Operator == "[" {
		writeTok(exp.Left.GetToken(), out)
	} else {
		writeExpression(exp.Left, leftClose, out)
	}

	writeTok(exp.Token, out) // operator

	if methodCall {
		writeTok(exp.Right.GetToken(), out)
	} else {
		if exp.Operator == "[" {
			out.Write([]byte("<expression>\n"))
		}
		writeExpression(exp.Right, rightClose, out)
		if exp.Operator == "[" {
			out.Write([]byte("</expression>\n"))
		}
	}

	if ok {
		writeTok(c, out) // closing ], )
		out.Write([]byte("</term>\n"))
	}
}

func writeExpression(exp ast.Expression, c bool, out io.Writer) {
	if c {
		out.Write([]byte("<term>\n"))
		writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
	}
	if exp != nil {
		walkExpression(exp, out)
	}
	if c {
		writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)
		out.Write([]byte("</term>\n"))
	}
}

func walkPrefixExp(exp *ast.PrefixExp, out io.Writer) {
	sur := false
	//         root
	//        /    \
	// operator      expression(is infix expression)
	// wrap expression(infix) with parentheses if root's precedence
	// is greater tha expression
	if infix, ok := exp.Right.(*ast.InfixExp); ok {
		rootPrec := parser.PREFIX
		rightPrec := parser.Precedences[infix.Token.Type]
		if rootPrec > rightPrec {
			sur = true
		}
	}

	// term: (unaryOp term)
	out.Write([]byte("<term>\n"))
	writeTok(exp.Token, out)
	if sur {
		out.Write([]byte("<term>\n"))
		writeTok(token.Token{Literal: "(", Type: token.LPAREN}, out)
	}
	walkExpression(exp.Right, out)
	if sur {
		writeTok(token.Token{Literal: ")", Type: token.RPAREN}, out)
		out.Write([]byte("</term>\n"))
	}
	out.Write([]byte("</term>\n"))
}

func walkExpressionList(expList *ast.ExpressionListExp, out io.Writer) {
	out.Write([]byte("<expressionList>\n"))

	expressions := expList.Expressions
	for i, exp := range expressions {
		out.Write([]byte("<expression>\n"))
		walkExpression(exp, out)
		out.Write([]byte("</expression>\n"))
		if i+1 < len(expressions) {
			writeTok(token.Token{Literal: ",", Type: token.COMMA}, out)
		}
	}

	out.Write([]byte("</expressionList>\n"))
}

func createOutputFile(path string) (*os.File, func(), error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return nil, func() {}, fmt.Errorf("%s, is not a file or it does not end with file extension .jack", path)
	}

	if ext != ".jack" {
		return nil, func() {}, fmt.Errorf("%s, file extension must end with .jack", path)
	}
	// TTT just to avoid name collison with with testing files ending with T.xml and TT.xml
	// this lexer supposed to produce similar output to file ending with .xml file
	outputPath := path[0:len(path)-len(ext)] + "TTT.xml"
	file, err := os.Create(outputPath)
	return file, func() { file.Close() }, err
}
