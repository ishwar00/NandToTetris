package parser

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/fatih/color"
	"github.com/ishwar00/JackAnalyzer/ast"
	errhandler "github.com/ishwar00/JackAnalyzer/errHandler"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

var green = color.GreenString
var yellow = color.YellowString
var red = color.RedString

type prefixFn func() ast.Expression
type infixFn func(ast.Expression) ast.Expression

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	prefixFns map[token.TokenType]prefixFn
	infixFns  map[token.TokenType]infixFn
	errors    errhandler.ErrHandler
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:         l,
		prefixFns: map[token.TokenType]prefixFn{},
		infixFns:  map[token.TokenType]infixFn{},
	}

	p.registerPrefixFn(token.INT, p.parseInteger)
	p.registerPrefixFn(token.IDENT, p.parseVarName)

	p.nextToken()
	p.nextToken()
	return p
}

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.ASTERI: PRODUCT,
	token.LPAREN: CALL,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixFn) {
	p.prefixFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixFn) {
	p.infixFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(tokenType token.TokenType) bool {
	return p.curToken.Type == tokenType
}

func (p *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) skipToSemicolon() {
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.SEMICO) {
		p.nextToken()
	}
}

func isNativeType(tok token.Token) bool {
	return tok.Type == token.INT || tok.Type == token.BOOLEAN || tok.Type == token.CHAR
}

func (p *Parser) addError(errMsg string, token token.Token) {
	e := errhandler.Error{
		ErrMsg:   errMsg,
		OnLine:   token.OnLine,
		OnColumn: token.OnColumn,
		Length:   len(token.Literal),
		File:     token.InFile,
	}
	p.errors.Add(e)
}

func (p *Parser) ReportErrors() {
	if p.errors.QueueSize() > 0 {
		os.Stdout.WriteString(red("syntax error(s)\n"))
		p.errors.ReportAll()
	}
}

func canBeType(tok token.Token) bool {
	if isNativeType(tok) {
		return true
	}
	return tok.Type == token.IDENT
}

// peeks for static or field variable declaration, if it finds
// var variable declaration, it adds error for var and skips
// whole var declaration and calls itself again.
func (p *Parser) peekClassVarDec() bool {
	if p.peekTokenIs(token.VAR) {
		errMsg := fmt.Sprintf("cannot use %s(used for local variable declaration), expected %s or %s",
			red("var"), green("static"), green("field"))

		p.addError(errMsg, p.peekToken)
		p.nextToken()       // consume semicolon ;
		p.skipToSemicolon() // skipping to next semicolon
		return p.peekClassVarDec()
	}
	return p.peekTokenIs(token.STATIC) || p.peekTokenIs(token.FIELD)
}

func (p *Parser) ParseProgram() ast.Declare {
	return p.parseClassDec()
}

// func (p *Parser) parseStatement() ast.Declare {
// 	switch p.curToken.Type {
// 	case token.IF:
// 	case token.DO:
// 	case token.LET:
//     case token.WHILE:
// 	default:
// 		return nil
// 	}
// }

// helper to add error message
func (p *Parser) expectedIdentifierErr(tok token.Token) {
	errMsg := "expected an identifier name"
	if token.IsKeyword(tok.Literal) {
		errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as identifier",
			yellow(tok.Literal))
	}
	p.addError(errMsg, tok)
}

// parses: static|field|var varName1, varName2, ... ;
func (p *Parser) parseVarDec() *ast.VarDec {
	varDec := &ast.VarDec{
		Token: p.curToken, // var, static, field keyword
	}
	// expected name of a 'data type' next, eg: int, boolean, char, Ball ... etc
	if !canBeType(p.peekToken) {
		errMsg := "expected an identifier representing name of a " + yellow("DataType") + " here, like int, boolean etc"
		if token.IsKeyword(p.peekToken.Literal) {
			errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as '%s' name",
				yellow(p.peekToken.Literal), yellow("DataType"))
		}
		p.addError(errMsg, p.peekToken)
		p.skipToSemicolon()
		return nil
	}

	p.nextToken()
	varDec.DataType = p.curToken

	// expecting identifier
	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		p.skipToSemicolon()
		return nil
	}
	p.nextToken()

	identifier := &ast.IdentifierDec{
		Token:    p.curToken,
		Literal:  p.curToken.Literal,
		DataType: varDec.DataType, // int, boolean, Ball ...
	}
	varDec.IdentifierDecs = append(varDec.IdentifierDecs, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // on ,(comma)
		if !p.peekTokenIs(token.IDENT) {
			p.expectedIdentifierErr(p.peekToken)
			p.skipToSemicolon()
			return nil
		}
		p.nextToken()
		identifier := &ast.IdentifierDec{
			Token:    p.curToken,
			Literal:  p.curToken.Literal,
			DataType: varDec.DataType,
		}
		varDec.IdentifierDecs = append(varDec.IdentifierDecs, identifier)
	}

	if !p.peekTokenIs(token.SEMICO) {
		errMsg := fmt.Sprintf("expected ';', but got '%s'", p.peekToken.Literal)
		p.addError(errMsg, p.peekToken)
		p.skipToSemicolon()
		return nil
	}
	p.nextToken()
	return varDec
}

func (p *Parser) parseClassDec() *ast.ClassDec {
	classDec := &ast.ClassDec{
		Token: p.curToken, // keyword class
	}

	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		return nil
	}

	p.nextToken()
	classDec.Name = p.curToken.Literal

	if !p.peekTokenIs(token.LBRACE) {
		errMsg := fmt.Sprintf("expected %s, but got %v", yellow("{"), red(p.peekToken.Literal))
		p.addError(errMsg, p.peekToken)
		return nil
	}
	p.nextToken()

	for p.peekClassVarDec() {
		p.nextToken()
		if varDec := p.parseVarDec(); varDec != nil {
			classDec.ClassVarDecs = append(classDec.ClassVarDecs, varDec)
		}
	}

	if !p.peekTokenIs(token.RBRACE) {
		errMsg := fmt.Sprintf("expected a %s, but got %+v", yellow("}"), red(p.peekToken.Literal))
		p.addError(errMsg, p.peekToken)
		return nil
	}
	return classDec
}

func (p *Parser) parseInteger() ast.Expression {
	v, err := strconv.ParseInt(p.curToken.Literal, 0, 16)
	if err != nil {
		p.addError(err.Error(), p.curToken)
		return nil
	}
	return &ast.IntConstantExp{
		Token: p.curToken,
		Value: int16(v),
	}
}

func (p *Parser) parseVarName() ast.Expression {
    return &ast.VarNameExp{ Token: p.curToken, Name: p.curToken.Literal }
}

func (p *Parser) parseLetStatement() ast.Statement {
    letSta := &ast.LetSta{Token: p.curToken}

    if !p.peekTokenIs(token.IDENT) {
        p.expectedIdentifierErr(p.peekToken)
        p.skipToSemicolon()
        return nil
    }

    p.nextToken()
    letSta.VarName = ast.VarNameExp{
        Token: p.curToken,
        Name: p.curToken.Literal,
    }

    if !p.peekTokenIs(token.EQ) {
        errMsg := "expected " + green("=") + " but got " + red(p.peekToken.Literal)
        p.addError(errMsg, p.peekToken)
        return nil
    }

    p.nextToken()
    p.nextToken()
    exp := p.parseExpression(LOWEST)
    if exp != nil && !reflect.ValueOf(exp).IsZero() {
        letSta.Expression = exp
        return letSta
    }
    return nil 
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixFns[p.curToken.Type]
	if !ok {
		errMsg := "expected an expression, but got " + red(p.curToken.Literal)
		p.addError(errMsg, p.curToken)
		return nil
	}
	leftExp := prefix()
	if !p.peekTokenIs(token.SEMICO) && precedence < p.peekPrecedence() {
		infix, ok := p.infixFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

    if p.peekTokenIs(token.SEMICO) {
        p.nextToken()
    }
	return leftExp
}




