package ast

import (
	"bytes"
	"fmt"

	"github.com/ishwar00/JackAnalyzer/token"
)

type Node interface {
	GetToken() token.Token
	String() string
}

// declarationNode, statementNode, expressionNode are dummy nodes,
// just to make sure that we don't use Defintion(say) instead of Statement

// Declare interface will be implemented by
// * class declaration
// * static, var and field variables declaration
// * function, method and constructor declaration
type Declare interface {
	Node
	declarationNode()
}

// Statement will be implemented by if, while, let and do statements
type Statement interface {
	Node
	statementNode()
}

// Expression will be implemented by constants and function call etc
type Expression interface {
	Node
	expressionNode()
}

type IdentifierExp struct {
	Token token.Token
	Value string
}

func (i *IdentifierExp) expressionNode() {}

func (i *IdentifierExp) GetToken() token.Token { return i.Token }

func (i *IdentifierExp) String() string {
	if i == nil {
		return ""
	}
	return i.Value
}

// static|field|var type varName1, varName2 ... ;
type VarDec struct {
	Token          token.Token // var, static, field
	DataType       token.Token // int, boolean, char, Array, Ball etc
	IdentifierExps []*IdentifierExp
}

func (vd *VarDec) declarationNode() {}

func (vd *VarDec) GetToken() token.Token { return vd.Token }

func (vd *VarDec) String() string {
	if vd == nil {
		return ""
	}

	var out bytes.Buffer
	space := " "

	out.WriteString(vd.Token.Literal + space)
	out.WriteString(vd.DataType.Literal + space)

	for i, identifier := range vd.IdentifierExps {
		out.WriteString(identifier.String())
		if i+1 < len(vd.IdentifierExps) {
			out.WriteString(",")
		}
	}
	out.WriteString(";")
	return out.String()
}

type ClassDec struct {
	Token        token.Token // class token
	Name         string      // name of the class
	ClassVarDecs []*VarDec
}

func (cd *ClassDec) GetToken() token.Token { return cd.Token }

func (cd *ClassDec) declarationNode() {}

func (cs *ClassDec) String() string {
	if cs == nil {
		return ""
	}

	var out bytes.Buffer
	tab := "    "

	buf := fmt.Sprintf("class %s {\n", cs.Name)
	out.WriteString(buf)
	for _, varDec := range cs.ClassVarDecs {
		out.WriteString(tab + varDec.String() + "\n")
	}
	out.WriteString("\n}")

	return out.String()
}

type IntConstExp struct {
	Token token.Token
	Value int64
}

func (i *IntConstExp) expressionNode() {}

func (i *IntConstExp) String() string {
	if i == nil {
		return ""
	}
	return i.Token.Literal
}

func (i *IntConstExp) GetToken() token.Token { return i.Token }

// let varName ('[' Index ']')? = Expression
type LetSta struct {
	Token      token.Token // let
	Name       Expression  // IdentifierExp, InfixExpression
	Expression Expression
}

func (l *LetSta) statementNode() {}

func (l *LetSta) GetToken() token.Token { return l.Token }

func (l *LetSta) String() string {
	if l == nil {
		return ""
	}

	var out bytes.Buffer

	out.WriteString("let " + l.Name.String())
	out.WriteString("=" + l.Expression.String() + ";")

	return out.String()
}

type StatementBlockSta struct {
	Token      token.Token // {
	Statements []Statement
}

func (sb *StatementBlockSta) statementNode() {}

func (sb *StatementBlockSta) GetToken() token.Token { return sb.Token }

func (sb *StatementBlockSta) String() string {
	if sb == nil {
		return ""
	}

	var out bytes.Buffer

	for _, stmt := range sb.Statements {
		out.WriteString(stmt.String() + "\n")
	}
	return out.String()
}

// if(Condition) {
//     Then
// }(else {
//     Else
// })?
type IfElseSta struct {
	Token     token.Token // if
	Condition Expression
	Then      *StatementBlockSta
	Else      *StatementBlockSta
}

func (ie *IfElseSta) statementNode() {}

func (ie *IfElseSta) GetToken() token.Token { return ie.Token }

func (ie *IfElseSta) String() string {
	if ie == nil {
		return ""
	}

	var out bytes.Buffer
	tab := "    "

	out.WriteString("if(" + ie.Condition.String() + ") {\n")
	out.WriteString(tab + ie.Then.String())
	if ie.Else != nil {
		out.WriteString("\n} else {\n" + tab + ie.Else.String())
	}
	out.WriteString("\n}\n")
	return out.String()
}

type WhileSta struct {
	Token      token.Token
	Condition  Expression
	Statements *StatementBlockSta
}

func (ws *WhileSta) statementNode() {}

func (ws *WhileSta) GetToken() token.Token { return ws.Token }

func (ws *WhileSta) String() string {
	if ws == nil {
		return ""
	}

	var out bytes.Buffer
	tab := "    "

	out.WriteString("while(" + ws.Condition.String() + ") {\n")
	out.WriteString(tab + ws.Statements.String())
	out.WriteString("\n}\n")

	return out.String()
}

type ReturnSta struct {
	Token      token.Token
	Expression Expression
}

func (r *ReturnSta) statementNode() {}

func (r *ReturnSta) GetToken() token.Token { return r.Token }

func (r *ReturnSta) String() string {
	if r == nil {
		return ""
	}

	var out bytes.Buffer
	out.WriteString("return " + r.Expression.String() + ";")
	return out.String()
}

type PrefixExp struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExp) expressionNode() {}

func (p *PrefixExp) GetToken() token.Token { return p.Token }

func (p *PrefixExp) String() string {
	if p == nil {
		return ""
	}
	var out bytes.Buffer

	out.WriteString("(" + p.Operator)
	out.WriteString(p.Right.String() + ")")
	return out.String()
}

type InfixExp struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExp) expressionNode() {}

func (i *InfixExp) GetToken() token.Token { return i.Token }

func (i *InfixExp) String() string {
	if i == nil {
		return ""
	}
	var out bytes.Buffer
	mp := map[string]string{
		"(": ")",
		"[": "]",
	}

	out.WriteString("(" + i.Left.String())
	out.WriteString(i.Operator)
	out.WriteString(i.Right.String())
	if i, ok := mp[i.Operator]; ok {
		out.WriteString(i)
	}
	out.WriteString(")")

	return out.String()
}

type StrConstExp struct {
	Token token.Token
	Value string
}

func (s *StrConstExp) expressionNode() {}

func (s *StrConstExp) GetToken() token.Token { return s.Token }

func (s *StrConstExp) String() string {
	if s == nil {
		return ""
	}
	return s.Value
}

type KeywordConstExp struct {
	Token token.Token
	Value string
}

func (k *KeywordConstExp) expressionNode() {}

func (k *KeywordConstExp) GetToken() token.Token { return k.Token }

func (k *KeywordConstExp) String() string {
	if k == nil {
		return ""
	}
	return k.Value
}

type ExpressionListExp struct {
	Token       token.Token // first token of first expression
	Expressions []Expression
}

func (el *ExpressionListExp) expressionNode() {}

func (el *ExpressionListExp) GetToken() token.Token { return el.Token }

func (el *ExpressionListExp) String() string {
	if el == nil {
		return ""
	}
	var out bytes.Buffer

	for i, exp := range el.Expressions {
		out.WriteString(exp.String())
		if i+1 < len(el.Expressions) {
			out.WriteString(",")
		}
	}
	return out.String()
}
