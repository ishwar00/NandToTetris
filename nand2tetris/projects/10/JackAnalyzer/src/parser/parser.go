package parser

import (
	"fmt"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/ishwar00/JackAnalyzer/ast"
	errhandler "github.com/ishwar00/JackAnalyzer/errHandler"
	"github.com/ishwar00/JackAnalyzer/lexer"
	"github.com/ishwar00/JackAnalyzer/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    errhandler.ErrHandler
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l: l,
	}
	p.nextToken()
	p.nextToken()
	return p
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
		os.Stdout.WriteString("syntax error(s)\n")
		p.errors.ReportAll()
	}
}

func canBeType(tok token.Token) bool {
	if isNativeType(tok) {
		return true
	}
	return tok.Type == token.IDENT
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil && !reflect.ValueOf(stmt).IsZero() {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarDec()
	default:
		return nil
	}
}

func (p *Parser) parseVarDec() *ast.VarDecStatement {
	varDec := &ast.VarDecStatement{
		Token: p.curToken, // 'var' keyword
	}
	// expected name of a 'data type' next, eg: int, boolean, char, Ball ...
	if !canBeType(p.peekToken) {
		errMsg := "expected an identifier representing name of a " + color.YellowString("DataType") + " here, eg: int, boolean etc"
		if token.IsKeyword(p.peekToken.Literal) {
			errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as '%s' name",
				color.YellowString(p.peekToken.Literal), color.YellowString("DataType"))
		}
		p.addError(errMsg, p.peekToken)
		p.skipToSemicolon()
		return nil
	}

	p.nextToken()
	varDec.DataType = p.curToken

	// expecting identifier
	if !p.peekTokenIs(token.IDENT) {
		errMsg := "expected an identifier name"
		if token.IsKeyword(p.peekToken.Literal) {
			errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as identifier",
				color.YellowString(p.peekToken.Literal))
		}
		p.addError(errMsg, p.peekToken)
		p.skipToSemicolon()
		return nil
	}
	p.nextToken()

	identifier := &ast.Identifier{
		Token:    p.curToken,
		Value:    p.curToken.Literal,
		DataType: varDec.DataType, // int, boolean, Ball ...
	}
	varDec.Identifiers = append(varDec.Identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ,(comma)
		if !p.peekTokenIs(token.IDENT) {
			errMsg := "expected an identifier"
			if token.IsKeyword(p.peekToken.Literal) {
				errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as identifier",
					color.YellowString(p.peekToken.Literal))
			}
			p.addError(errMsg, p.peekToken)
			p.skipToSemicolon()
			return nil
		}
		p.nextToken()
		identifier := &ast.Identifier{
			Token:    p.curToken,
			Value:    p.curToken.Literal,
			DataType: varDec.DataType,
		}
		varDec.Identifiers = append(varDec.Identifiers, identifier)
	}

	if !p.peekTokenIs(token.SEMICO) {
		errMsg := fmt.Sprintf("expected ';', but got '%s'", p.peekToken.Literal)
		p.addError(errMsg, p.peekToken)
		p.skipToSemicolon()
		return nil
	}
	return varDec
}
