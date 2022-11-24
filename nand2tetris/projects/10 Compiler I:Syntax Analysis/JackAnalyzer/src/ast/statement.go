package ast

import (
	"bytes"

	"github.com/ishwar00/JackAnalyzer/token"
)

// statementNode is a dummy node,
// just to make sure that we don't use Defintion(say) instead of Statement

// Statement interface is implemented by if, while, let and do statements
type Statement interface {
	Node
	statementNode()
}

// let varName ('[' Index ']')? = Expression
type LetSta struct {
	Token      token.Token // let
	Name       Expression  // IdentifierExp, InfixExpression
	Expression Expression
}

func (_ *LetSta) statementNode() {}

func (l *LetSta) GetToken() token.Token { return l.Token }

func (l *LetSta) String() string {
	if l == nil {
		return ""
	}

	var out bytes.Buffer

	out.WriteString("let " + l.Name.String())
	out.WriteString(" = " + l.Expression.String() + ";")

	return out.String()
}

type StatementBlockSta struct {
	Token      token.Token // {
	Statements []Statement
}

func (_ *StatementBlockSta) statementNode() {}

func (sb *StatementBlockSta) GetToken() token.Token { return sb.Token }

func (sb *StatementBlockSta) String() string {
	if sb == nil {
		return ""
	}

	var out bytes.Buffer

	tab := "    "
	for _, stmt := range sb.Statements {
		out.WriteString(tab + stmt.String() + "\n")
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

func (_ *IfElseSta) statementNode() {}

func (ie *IfElseSta) GetToken() token.Token { return ie.Token }

func (ie *IfElseSta) String() string {
	if ie == nil {
		return ""
	}

	var out bytes.Buffer

	out.WriteString("if(" + ie.Condition.String() + ") {\n")
	out.WriteString(ie.Then.String())
	if ie.Else != nil {
		out.WriteString("\n} else {\n" + ie.Else.String())
	}
	out.WriteString("\n}\n")
	return out.String()
}

type WhileSta struct {
	Token      token.Token
	Condition  Expression
	Statements *StatementBlockSta
}

func (_ *WhileSta) statementNode() {}

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

type DoSta struct {
	Token   token.Token // do
	SubCall Expression
}

func (_ *DoSta) statementNode() {}

func (d *DoSta) GetToken() token.Token { return d.Token }

func (d *DoSta) String() string {
	var out bytes.Buffer

	out.WriteString("do ")
	out.WriteString(d.SubCall.String() + ";")
	return out.String()
}

type ReturnSta struct {
	Token      token.Token
	Expression Expression
}

func (_ *ReturnSta) statementNode() {}

func (r *ReturnSta) GetToken() token.Token { return r.Token }

func (r *ReturnSta) String() string {
	if r == nil {
		return ""
	}

	var out bytes.Buffer
	out.WriteString("return ")
	if r.Expression != nil {
		out.WriteString(r.Expression.String())
	}
	out.WriteString(";")
	return out.String()
}
