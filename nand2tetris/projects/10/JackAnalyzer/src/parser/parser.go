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

var (
	green  = color.GreenString
	yellow = color.YellowString
	red    = color.RedString
)

var statementTokens = []token.TokenType{
	token.LET,
	token.WHILE,
	token.IF,
	token.RETURN,
	token.DO,
}

type (
	prefixFn func() ast.Expression
	infixFn  func(ast.Expression) ast.Expression
)

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

	p.registerPrefixFn(token.INT_CONST, p.parseInteger)
	p.registerPrefixFn(token.IDENT, p.parseIdentifierExp)
	p.registerPrefixFn(token.STR_CONST, p.parseStrconstantExp)
	p.registerPrefixFn(token.TRUE, p.parseKeywordConstantExp)
	p.registerPrefixFn(token.FALSE, p.parseKeywordConstantExp)
	p.registerPrefixFn(token.THIS, p.parseKeywordConstantExp)
	p.registerPrefixFn(token.NULL, p.parseKeywordConstantExp)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExp)
	p.registerPrefixFn(token.TILDE, p.parsePrefixExp)
	p.registerPrefixFn(token.LPAREN, p.parseGroupExp)

	p.registerInfixFn(token.LBRACK, p.parseArrayIndex)
	p.registerInfixFn(token.LPAREN, p.parseSubroutineCall)
	p.registerInfixFn(token.PLUS, p.parseInfixExp)
	p.registerInfixFn(token.MINUS, p.parseInfixExp)
	p.registerInfixFn(token.SLASH, p.parseInfixExp)
	p.registerInfixFn(token.ASTERI, p.parseInfixExp)
	p.registerInfixFn(token.LT, p.parseInfixExp)
	p.registerInfixFn(token.GT, p.parseInfixExp)
	p.registerInfixFn(token.EQ, p.parseInfixExp)
	p.registerInfixFn(token.PERIOD, p.parseMethodCall)
	p.registerInfixFn(token.AMPERS, p.parseInfixExp)
	p.registerInfixFn(token.PIPE, p.parseInfixExp)

	p.nextToken()
	p.nextToken()
	return p
}

const (
	_ int = iota
	LOWEST
	OR                // bitwise or
	AND               // bitwise and
	EQUALS            // ==
	LESSGREATER       // > or <
	SUM               // +
	PRODUCT           // *
	PREFIX            // -X or ~X
	CALL_INDEX_PERIOD // myFunction(X), arr[4], foo.bar()
)

var precedences = map[token.TokenType]int{
	token.EQ:     EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.ASTERI: PRODUCT,
	token.LPAREN: CALL_INDEX_PERIOD,
	token.LBRACK: CALL_INDEX_PERIOD,
	token.PERIOD: CALL_INDEX_PERIOD,
	token.AMPERS: AND,
	token.PIPE:   OR,
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

func (p *Parser) expectedTypeErr(tok token.Token) {
	errMsg := "expected an identifier representing name of a " + yellow("DataType") + " here, like int, boolean etc"
	if token.IsKeyword(tok.Literal) {
		errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as '%s' name",
			yellow(tok.Literal), yellow("DataType"))
	}
	p.addError(errMsg, tok)
}

// peeks for static or field variable declaration, if it finds
// var variable declaration, it adds error for var and skips
// whole var declaration and calls itself again.
func (p *Parser) peekClassVarDec() bool {
	if p.peekTokenIs(token.VAR) {
		errMsg := fmt.Sprintf("cannot use %s(used for local variable declaration), expected %s or %s",
			red("var"), green("static"), green("field"))

		p.addError(errMsg, p.peekToken)
		p.nextToken()                 // consume semicolon ;
		p.skipToNextSta(token.SEMICO) // skipping to next semicolon
		return p.peekClassVarDec()
	}
	return p.peekTokenIs(token.STATIC) || p.peekTokenIs(token.FIELD)
}

func (p *Parser) peekStatement() bool {
	return p.peekTokenIs(token.LET) || p.peekTokenIs(token.IF) ||
		p.peekTokenIs(token.WHILE) || p.peekTokenIs(token.DO) || p.peekTokenIs(token.RETURN)
}

// skip to next toks, which ever comes first
func (p *Parser) skipToNextSta(toks ...token.TokenType) {
	toks = append(toks, statementTokens...)
	match := func() bool {
		for _, tok := range toks {
			if p.peekTokenIs(tok) {
				return true
			}
		}
		return false
	}

	for !p.peekTokenIs(token.EOF) && !match() {
		p.nextToken()
	}
}

func (p *Parser) ParseProgram() {
	// buf := p.parseIfElseSta()
	// if buf != nil && !reflect.ValueOf(buf).IsZero() {
	//     fmt.Println(buf.String())
	// }
	var buf interface {
		String() string
	}
	for !p.curTokenIs(token.EOF) {
		switch p.curToken.Type {
		case token.CLASS:
			// buf = p.parseClassDec()
		case token.IF:
			buf = p.parseIfElseSta()
		case token.WHILE:
			buf = p.parseWhileSta()
		case token.RETURN:
			buf = p.parseReturnSta()
		case token.LET:
			buf = p.parseLetStatement()
		}
		if buf != nil && !reflect.ValueOf(buf).IsZero() {
			fmt.Println(buf.String())
		}
		p.nextToken()
	}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IF:
		return p.parseIfElseSta()
	case token.DO:
		return nil
	case token.LET:
		return p.parseLetStatement()
	case token.WHILE:
		return p.parseWhileSta()
	case token.RETURN:
		return p.parseReturnSta()
	default:
		return nil
	}
}

// helper to add error message
func (p *Parser) expectedIdentifierErr(tok token.Token) {
	errMsg := "expected an identifier name"
	if token.IsKeyword(tok.Literal) {
		errMsg = fmt.Sprintf("cannot use reserved keyword '%s' as identifier",
			yellow(tok.Literal))
	}
	p.addError(errMsg, tok)
}

func (p *Parser) expectedErr(expected string, got token.Token) {
	errMsg := "expected " + green(expected) + " but got " + red(got.Literal)
	p.addError(errMsg, got)
}

// parses: static|field|var varName1, varName2, ... ;
func (p *Parser) parseVarDec() *ast.VarDec {
	varDec := &ast.VarDec{
		Token: p.curToken, // var, static, field keyword
	}
	// expected name of a 'data type' next, eg: int, boolean, char, Ball ... etc
	if !canBeType(p.peekToken) {
		p.expectedTypeErr(p.peekToken)
		p.skipToNextSta(token.SEMICO)
		return nil
	}

	p.nextToken()
	varDec.DataType = p.curToken

	// expecting identifier
	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		p.skipToNextSta(token.SEMICO)
		return nil
	}
	p.nextToken()

	identifier := &ast.IdentifierExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	varDec.IdentifierExps = append(varDec.IdentifierExps, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // on ,(comma)
		if !p.peekTokenIs(token.IDENT) {
			p.expectedIdentifierErr(p.peekToken)
			p.skipToNextSta(token.SEMICO)
			return nil
		}
		p.nextToken()
		identifier := &ast.IdentifierExp{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
		varDec.IdentifierExps = append(varDec.IdentifierExps, identifier)
	}

	if !p.peekTokenIs(token.SEMICO) {
		p.expectedErr(";", p.peekToken)
		p.skipToNextSta(token.SEMICO)
		return nil
	}
	p.nextToken()
	return varDec
}

func (p *Parser) peekSubroutine() bool {
	return p.peekToken.Type == token.METHOD ||
		p.peekToken.Type == token.CONSTRUCTOR ||
		p.peekToken.Type == token.FUNCTION
}

// class className {
//     static|field varName1, varName2, ...;
//     function|method|constructor type subName(type1 p1, type2 p2, ...) {
//         var v1, v2 ...;
//         if|while|do|return statements
//    }
// }
func (p *Parser) ParseClassDec() *ast.ClassDec {
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
		p.expectedErr("{", p.peekToken)
		return nil
	}
	p.nextToken()

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		if p.peekClassVarDec() {
			p.nextToken()
			if varDec := p.parseVarDec(); varDec != nil {
				classDec.ClassVarDecs = append(classDec.ClassVarDecs, varDec)
			}
		} else if p.peekSubroutine() {
			p.nextToken()
			subroutine := p.parseSubroutineBodyDec()
			if subroutine != nil {
				classDec.Subroutines = append(classDec.Subroutines, subroutine)
			}
		} else {
			p.expectedErr("expected method or function or constructor or static or field",
				p.peekToken)
		}
		fmt.Println(p.curToken)
	}

	if !p.peekTokenIs(token.RBRACE) {
		p.expectedErr("}", p.peekToken)
		return nil
	}
	return classDec
}

func (p *Parser) parseInteger() ast.Expression {
	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(err.Error(), p.curToken)
		return nil
	}
	return &ast.IntConstExp{
		Token: p.curToken,
		Value: v,
	}
}

func (p *Parser) parseIdentifierExp() ast.Expression {
	return &ast.IdentifierExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseMethodCall(varName ast.Expression) ast.Expression {
	infix := &ast.InfixExp{
		Token:    p.curToken,         // .
		Operator: p.curToken.Literal, // .
		Left:     varName,
	}

	if !p.peekTokenIs(token.IDENT) {
		p.expectedErr("Identifier", p.peekToken)
		p.skipToNextSta(token.SEMICO)
		return nil
	}
	precedence := precedences[p.curToken.Type]
	p.nextToken()
	infix.Right = p.parseExpression(precedence)
	return infix
}

// let Name ('[' Index ']')? = Expression
func (p *Parser) parseLetStatement() *ast.LetSta {
	letSta := &ast.LetSta{
		Token: p.curToken, // let
	}

	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		p.skipToNextSta(token.SEMICO)
		return nil
	}

	p.nextToken()
	letSta.Name = &ast.IdentifierExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(token.LBRACK) {
		p.nextToken()
		letSta.Name = p.parseArrayIndex(letSta.Name)
	}

	if !p.peekTokenIs(token.EQ) {
		p.expectedErr("=", p.peekToken)
		return nil
	}

	p.nextToken()
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if exp == nil || reflect.ValueOf(exp).IsZero() {
		return nil
	}

	letSta.Expression = exp
	if !p.peekTokenIs(token.SEMICO) {
		p.expectedErr(";", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	return letSta
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixFns[p.curToken.Type]
	if !ok {
		p.expectedErr("Expression", p.curToken)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICO) && precedence < p.peekPrecedence() {
		infix, ok := p.infixFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

// if(Condition) {
//    statements
// } (else {
//    statements
// })?
func (p *Parser) parseIfElseSta() *ast.IfElseSta {
	ifElse := &ast.IfElseSta{
		Token: p.curToken,
	}

	expErr := func(exp string, got token.Token) {
		p.expectedErr(exp, got)
		p.skipToNextSta()
	}

	if !p.peekTokenIs(token.LPAREN) {
		expErr("(", p.peekToken)
		p.skipToNextSta()
		return nil
	}

	p.nextToken()
	p.nextToken()
	condition := p.parseExpression(LOWEST)
	if condition == nil || reflect.ValueOf(condition).IsZero() {
		p.skipToNextSta()
		return nil
	}

	ifElse.Condition = condition

	if !p.peekTokenIs(token.RPAREN) {
		expErr(")", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()

	if !p.peekTokenIs(token.LBRACE) {
		expErr("{", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	ifElse.Then = p.parseStatementBlock()
	if ifElse.Then == nil {
		p.skipToNextSta()
		return nil
	}

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.peekTokenIs(token.LBRACE) {
			expErr("{", p.peekToken)
			p.skipToNextSta()
			return nil
		}
		p.nextToken()
		ifElse.Else = p.parseStatementBlock()
		if ifElse.Else == nil {
			p.skipToNextSta()
			return nil
		}
	}
	return ifElse
}

// { statement1; statement2; ... statementN; }
func (p *Parser) parseStatementBlock() *ast.StatementBlockSta {
	stmtBlock := &ast.StatementBlockSta{
		Token: p.curToken, // {
	}

	for p.peekStatement() {
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil && !reflect.ValueOf(stmt).IsZero() {
			stmtBlock.Statements = append(stmtBlock.Statements, stmt)
		}
	}

	if !p.peekTokenIs(token.RBRACE) {
		p.expectedErr("}", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	return stmtBlock
}

func (p *Parser) parseWhileSta() *ast.WhileSta {
	whileSta := &ast.WhileSta{
		Token: p.curToken,
	}

	if !p.peekTokenIs(token.LPAREN) {
		p.expectedErr("(", p.peekToken)
		p.skipToNextSta()
		return nil
	}

	p.nextToken()
	p.nextToken()
	condition := p.parseExpression(LOWEST)
	if condition == nil || reflect.ValueOf(condition).IsZero() {
		p.skipToNextSta()
		return nil
	}
	whileSta.Condition = condition
	if !p.peekTokenIs(token.RPAREN) {
		p.expectedErr(")", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	p.nextToken()
	whileSta.Statements = p.parseStatementBlock()
	if whileSta.Statements == nil {
		p.skipToNextSta()
		return nil
	}
	return whileSta
}

// return; return Expression;
func (p *Parser) parseReturnSta() *ast.ReturnSta {
	returnSta := &ast.ReturnSta{
		Token: p.curToken,
	}

	if !p.peekTokenIs(token.SEMICO) {
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		if exp == nil || reflect.ValueOf(exp).IsZero() {
			p.skipToNextSta()
			return nil
		}
		returnSta.Expression = exp
	}

	if !p.peekTokenIs(token.SEMICO) {
		p.expectedErr(";", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()

	return returnSta
}

func (p *Parser) parseArrayIndex(left ast.Expression) ast.Expression {
	arrayIndex := &ast.InfixExp{
		Token:    p.curToken,         // [
		Operator: p.curToken.Literal, // [
		Left:     left,
	}

	p.nextToken()
	arrayIndex.Right = p.parseExpression(CALL_INDEX_PERIOD)

	if !p.peekTokenIs(token.RBRACK) {
		p.expectedErr("]", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	return arrayIndex
}

func (p *Parser) parseStrconstantExp() ast.Expression {
	return &ast.StrConstExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseKeywordConstantExp() ast.Expression {
	return &ast.KeywordConstExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

// parses: expression0, expression1, ... expresssionN
// if it encounters an error while parsing ith expression
// it will skip to one of the 'skipTo' tokens, continues for more errors
func (p *Parser) parseExpressionList(skipTo ...token.TokenType) ast.Expression {
	expList := &ast.ExpressionListExp{
		Token: p.curToken,
	}
	err := false

	exp := p.parseExpression(LOWEST)
	if exp != nil && !reflect.ValueOf(exp).IsZero() {
		expList.Expressions = append(expList.Expressions, exp)
	} else {
		p.skipToNextSta(skipTo...)
		err = true
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exp := p.parseExpression(LOWEST)
		if exp != nil && !reflect.ValueOf(exp).IsZero() {
			expList.Expressions = append(expList.Expressions, exp)
		} else {
			err = true
			p.skipToNextSta(skipTo...)
		}
	}
	if err {
		return nil
	}
	return expList
}

func (p *Parser) parsePrefixExp() ast.Expression {
	exp := &ast.PrefixExp{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseInfixExp(leftExp ast.Expression) ast.Expression {
	exp := &ast.InfixExp{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     leftExp,
	}

	precedence := precedences[p.curToken.Type]
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) parseSubroutineCall(leftExp ast.Expression) ast.Expression {
	subCall := &ast.InfixExp{
		Token:    p.curToken,         // (
		Operator: p.curToken.Literal, // (
		Left:     leftExp,
	}

	if p.peekTokenIs(token.RPAREN) { // empty expression list
		p.nextToken()
		return subCall
	}

	p.nextToken()
	subCall.Right = p.parseExpressionList()
	if !p.peekTokenIs(token.RPAREN) {
		p.expectedErr(")", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	return subCall
}

func (p *Parser) parseGroupExp() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.RPAREN) {
		p.expectedErr(")", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) parseDoSta() *ast.DoSta {
	doSta := &ast.DoSta{
		Token: p.curToken,
	}

	p.nextToken()
	doSta.SubCall = p.parseExpression(LOWEST)
	if !p.peekTokenIs(token.SEMICO) {
		p.expectedErr(";", p.peekToken)
		p.skipToNextSta()
		return nil
	}
	return doSta
}

func (p *Parser) parseParameterDec() *ast.ParameterDec {
	parameter := &ast.ParameterDec{
		Token:    p.curToken,
		DataType: p.curToken.Literal,
	}

	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		p.skipToNextSta(token.COMMA, token.LPAREN)
		return nil
	}
	p.nextToken()
	parameter.Identifier = &ast.IdentifierExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	return parameter
}

func (p *Parser) parseParameterListDec() []*ast.ParameterDec {
	parameters := []*ast.ParameterDec{}
	err := false

	parameter := p.parseParameterDec()
	if parameter != nil {
		parameters = append(parameters, parameter)
	} else {
		err = true
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if !canBeType(p.peekToken) {
			p.expectedTypeErr(p.peekToken)
			p.skipToNextSta(token.COMMA, token.LPAREN)
			err = true
			continue
		}
		p.nextToken()
		parameter = p.parseParameterDec()
		if parameter == nil {
			err = true
			continue
		}
		parameters = append(parameters, parameter)
	}
	if err {
		return nil
	}
	return parameters
}

func (p *Parser) parseSubroutineDec() *ast.SubroutineDec {
	subroutine := &ast.SubroutineDec{
		Token: p.curToken,
		Type:  p.curToken.Literal,
	}

	if !canBeType(p.peekToken) {
		p.expectedTypeErr(p.peekToken)
		p.skipToNextSta(token.LBRACE)
		return nil
	}
	p.nextToken()
	subroutine.ReturnType = p.curToken

	if !p.peekTokenIs(token.IDENT) {
		p.expectedIdentifierErr(p.peekToken)
		p.skipToNextSta(token.LBRACE)
		return nil
	}
	p.nextToken()
	subroutine.SubName = &ast.IdentifierExp{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.peekTokenIs(token.LPAREN) {
		p.expectedErr("(", p.peekToken)
		p.skipToNextSta(token.LBRACE)
		return nil
	}
	p.nextToken()
	if p.peekTokenIs(token.RPAREN) { // zero parameters
		p.nextToken()
		return subroutine
	}

	if !canBeType(p.peekToken) {
		p.expectedTypeErr(p.peekToken)
		p.skipToNextSta(token.LBRACE)
		return nil
	}

	p.nextToken()
	subroutine.Parameters = p.parseParameterListDec()
	if !p.peekTokenIs(token.RPAREN) {
		p.expectedErr(")", p.peekToken)
		p.skipToNextSta(token.LBRACE)
		return nil
	}
	p.nextToken()
	subroutine.Body = *p.parseSubroutineBodyDec()

	return subroutine
}

func (p *Parser) parseSubroutineBodyDec() *ast.SubroutineBodyDec {
	body := &ast.SubroutineBodyDec{
		Token: p.curToken, // {
	}

	for !p.peekTokenIs(token.LBRACE) && !p.peekTokenIs(token.EOF) {
		if p.peekStatement() || p.peekTokenIs(token.VAR) {
			p.nextToken()
			if p.curToken.Type == token.VAR {
				varDec := p.parseVarDec()
				if varDec != nil {
					body.VarDecs = append(body.VarDecs, varDec)
				}
			} else {
				stmt := p.parseStatement()
				if stmt != nil && !reflect.ValueOf(stmt).IsZero() {
					body.Statements = append(body.Statements, stmt)
				}
			}
		} else {
			p.expectedErr("var or do or if or while or return", p.peekToken)
			p.skipToNextSta(token.LBRACE)
		}
	}
	return body
}
