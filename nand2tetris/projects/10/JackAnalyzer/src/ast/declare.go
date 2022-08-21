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

// declarationNode is dummy node,
// just to make sure that we don't use Defintion(say) instead of Statement

// Declare interface is implemented by
// * class declaration
// * static, var and field variables declaration
// * function, method and constructor declaration
type Declare interface {
	Node
	declarationNode()
}

// static|field|var type varName1, varName2 ... ;
type VarDec struct {
	Token          token.Token // var, static, field
	DataType       token.Token // int, boolean, char, Array, Ball etc
	IdentifierExps []*IdentifierExp
}

func (_ *VarDec) declarationNode() {}

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
	Subroutines  []*SubroutineDec
}

func (cd *ClassDec) GetToken() token.Token { return cd.Token }

func (_ *ClassDec) declarationNode() {}

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

	for _, subroutine := range cs.Subroutines {
		out.WriteString(tab + subroutine.String() + "\n")
	}
	out.WriteString("\n}")

	return out.String()
}

type ParameterDec struct {
	Token      token.Token
	DataType   string
	Identifier *IdentifierExp
}

func (_ *ParameterDec) declarationNode() {}

func (p *ParameterDec) GetToken() token.Token { return p.Token }

func (p *ParameterDec) String() string {
	if p == nil {
		return ""
	}
	var out bytes.Buffer

	out.WriteString(p.DataType + " " + p.Identifier.String())
	return out.String()
}

// TODO: little clunky, do better
type SubroutineDec struct {
	Token      token.Token // function|method|constructor
	Type       string
	ReturnType token.Token
	SubName    *IdentifierExp // subroutine name
	Parameters []*ParameterDec
	Body       *SubroutineBodyDec
}

func (_ *SubroutineDec) declarationNode() {}

func (s *SubroutineDec) GetToken() token.Token { return s.Token }

func (s *SubroutineDec) String() string {
	if s == nil {
		return ""
	}
	var out bytes.Buffer

	out.WriteString(s.Type + " " + s.ReturnType.Literal + " ")
	out.WriteString(s.SubName.String() + "(")
	for i, param := range s.Parameters {
		out.WriteString(param.String())
		if i+1 < len(s.Parameters) {
			out.WriteString(",")
		}
	}
	out.WriteString(")")
	out.WriteString(s.Body.String())
	return out.String()
}

type SubroutineBodyDec struct {
	Token      token.Token
	VarDecs    []*VarDec
	Statements []Statement
}

func (_ *SubroutineBodyDec) declarationNode() {}

func (s *SubroutineBodyDec) GetToken() token.Token { return s.Token }

func (s *SubroutineBodyDec) String() string {
	if s == nil {
		return ""
	}
	var out bytes.Buffer

	tab := "    "
	out.WriteString("{\n")
	for _, stmt := range s.Statements {
		out.WriteString(tab + stmt.String() + "\n")
	}
	out.WriteString("\n}")
	return out.String()
}
