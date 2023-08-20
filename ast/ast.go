package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

// 语句
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token //Token.Ident    Token.Literal() == Value
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatment struct {
	Token token.Token //token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatment) statementNode() {}

func (ls *LetStatment) TokenLiteral() string {
	return ls.Token.Literal
}
