package ast

import "github.com/lucaengelhard/lang/src/lexer"

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assigne   Expr
	Operator  lexer.Token
	RightExpr Expr
}

func (n AssignmentExpr) expr() {}

type StructInstantiationExpr struct {
	StructIdentifier string
	Properties       map[string]Expr
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Type     Type
	Elements []Expr
}

func (n ArrayInstantiationExpr) expr() {}

type FnCallExpr struct {
	Identifier string
	Arguments  []Expr
}

func (n FnCallExpr) expr() {}
