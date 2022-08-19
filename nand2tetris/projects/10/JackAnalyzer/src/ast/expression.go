package ast

import (
	"bytes"

	"github.com/ishwar00/JackAnalyzer/token"
)

// expressionNode is a dummy node,
// just to make sure that we don't use Defintion(say) instead of Statement

// Expression interface is implemented by constants and function call etc
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

    // okay may be DRY, but it is more readable
    if c, ok := mp[i.Operator]; ok {
        out.WriteString(i.Left.String())
        out.WriteString(i.Operator)
        if i.Right != nil {
            out.WriteString(i.Right.String())
        }
        out.WriteString(c)
    } else {
        out.WriteString("(" + i.Left.String())
        out.WriteString(i.Operator)
        out.WriteString(i.Right.String() + ")")
    }

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
    return "\"" + s.Value + "\""
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
