package ast

import (
	"bytes"

	"github.com/ishwar00/JackAnalyzer/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
}

type Program struct {
	Statements []Statement
}

// these two methods may not make sense, but Program meets Statement interface
func (p *Program) TokenLiteral() string { return p.Statements[0].TokenLiteral() }

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	return out.String()
}

type Identifier struct {
	Token    token.Token
	Value    string      // class name, variable name, subroutine name
	DataType token.Token // int, boolean, char, Array, Ball
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Token.Literal }

type VarDecStatement struct {
	Token       token.Token // var
	DataType    token.Token // int, boolean, char, Array, Ball
	Identifiers []*Identifier
}

func (vd *VarDecStatement) TokenLiteral() string { return vd.Token.Literal }
func (vd *VarDecStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("var ")
	buffer.WriteString(vd.DataType.Literal)
	buffer.WriteString(" ")
	for i, varName := range vd.Identifiers {
		buffer.WriteString(varName.Value)
		if i+1 < len(vd.Identifiers) {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString(";")
	return buffer.String()
}
