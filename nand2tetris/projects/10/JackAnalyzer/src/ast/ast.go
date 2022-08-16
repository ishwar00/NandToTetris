package ast

import (
	"bytes"
	"fmt"
	"reflect"

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

type IdentifierDec struct {
	Token    token.Token
	Literal  string      // yeah, it is literal, Literal == Token.Literal
	DataType token.Token // dataType of data referred by identifier
}

func (i *IdentifierDec) GetToken() token.Token { return i.Token }
func (i *IdentifierDec) String() string        { return i.Token.Literal }
func (i *IdentifierDec) DefinitionNode()       {}

// static|field|var type varName1, varName2 ... ;
type VarDec struct {
	Token          token.Token // var, static, field
	DataType       token.Token // int, boolean, char, Array, Ball
	IdentifierDecs []*IdentifierDec
}

func (vd *VarDec) GetToken() token.Token { return vd.Token }
func (vd *VarDec) declarationNode()      {}

func (vd *VarDec) String() string {
	if vd == nil {
		return ""
	}

	var out bytes.Buffer

	buf := fmt.Sprintf("%s %s ", vd.Token.Literal, vd.DataType.Literal)
	out.WriteString(buf)

	for i, varName := range vd.IdentifierDecs {
		buf = fmt.Sprintf("%s: %s", varName.Literal, varName.DataType.Literal)
		out.WriteString(buf)
		if i+1 < len(vd.IdentifierDecs) {
			out.WriteString(", ")
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
func (cd *ClassDec) declarationNode()      {}

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
	out.WriteString("}")

	return out.String()
}

type IntConstantExp struct {
	Token token.Token
	Value int16
}

func (i *IntConstantExp) expressionNode() {}

func (i *IntConstantExp) String() string {
	if i == nil {
		return ""
	}
	return i.Token.Literal
}

func (i *IntConstantExp) GetToken() token.Token { return i.Token }

type VarNameExp struct {
	Token      token.Token
	Name       string
	Expression Expression
}

func (v *VarNameExp) expressionNode() {}

func (v *VarNameExp) String() string {
	if v == nil {
		return ""
	}
	return v.Name
}

func (v *VarNameExp) GetToken() token.Token { return v.Token }

// let varName ('[' Index ']')? = Expression
type LetSta struct {
	Token      token.Token // let
	VarName    VarNameExp
	Index      Expression
	Expression Expression
}

func (l *LetSta) statementNode() {}

func (l *LetSta) GetToken() token.Token { return l.Token }

func (l *LetSta) String() string {
	if l == nil {
		return ""
	}

	var out bytes.Buffer

	out.WriteString("let " + l.VarName.Name)
	if l.Index != nil && !reflect.ValueOf(l.Index).IsZero() {
		out.WriteString("[" + l.Index.String() + "]")
	}
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
		if stmt != nil && !reflect.ValueOf(stmt).IsZero() {
			out.WriteString(stmt.String() + "\n")
		}
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

	out.WriteString("if(" + ie.Condition.String() + ") {\n")
	out.WriteString(ie.Then.String() + "\n} ")
	if ie.Else != nil {
		out.WriteString("else {\n" + ie.Else.String())
		out.WriteString("\n}")
	}
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

	out.WriteString("while(" + ws.Condition.String() + ") {\n")
	out.WriteString(ws.Statements.String())
	out.WriteString("\n}")

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
	out.WriteString("return ")
	if r.Expression != nil {
		out.WriteString(r.Expression.String())
	}
	out.WriteString(";")
	return out.String()
}
