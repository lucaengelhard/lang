package ast

import "github.com/lucaengelhard/lang/src/lexer"

type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type IntExpr struct {
	Value int64
}

func (n IntExpr) expr() {}

type BoolExpr struct {
	Value bool
}

func (n BoolExpr) expr() {}

type FloatExpr struct {
	Value float64
}

func (n FloatExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

type UnknowPrimary struct{}

func (n UnknowPrimary) expr() {}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator lexer.Token
	Right    Expr
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Right    Expr
}

func (n AssignmentExpr) expr() {}

type ChainExpr struct {
	Assignee Expr
	Member   Expr
}

func (n ChainExpr) expr() {}

type StructInstantiationExpr struct {
	StructIdentifier string
	Properties       map[string]Expr
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Elements []Expr
}

func (n ArrayInstantiationExpr) expr() {}

type FnCallArg struct {
	Identifier string
	Value      Expr
}

type FnCallExpr struct {
	Caller    Expr
	Arguments []FnCallArg
}

func (n FnCallExpr) expr() {}

type FnArg struct {
	Identifier string
	Position   int
	IsMutable  bool
	Type       Type
}

type FnDeclareExpr struct {
	Arguments  map[string]FnArg
	Type       Type
	ReturnType Type
	Body       BlockStmt
}

func (n FnDeclareExpr) expr() {}

type IsTypeExpr struct {
	Left  Expr
	Right Type
}

func (n IsTypeExpr) expr() {}
