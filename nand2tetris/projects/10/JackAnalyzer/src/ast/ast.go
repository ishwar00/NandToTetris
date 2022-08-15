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

func (i *IntConstantExp) String() string { return i.Token.Literal }

func (i *IntConstantExp) GetToken() token.Token { return i.Token }

type VarNameExp struct {
	Token      token.Token
	Name       string
	Expression Expression
}

func (v *VarNameExp) expressionNode() {}

func (v *VarNameExp) String() string {
	return v.Name
}

func (v *VarNameExp) GetToken() token.Token { return v.Token }

// let varName = Expression
type LetSta struct {
	Token      token.Token // let
	VarName    VarNameExp
	Expression Expression
}

func (l *LetSta) statementNode() {}

func (l *LetSta) GetToken() token.Token { return l.Token }

func (l *LetSta) String() string {
	var out bytes.Buffer

	buf := fmt.Sprintf("let %s = ", l.VarName.Name)
	out.WriteString(buf)
	out.WriteString(l.Expression.String() + ";")

	return out.String()
}
