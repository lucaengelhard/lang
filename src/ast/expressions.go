package ast

import "github.com/lucaengelhard/lang/src/lexer"

type NumberExpr struct {
	Value float64
	Position
}

func (n NumberExpr) expr() {}

type IntExpr struct {
	Value int64
	Position
}

func (n IntExpr) expr() {}

type BoolExpr struct {
	Value bool
	Position
}

func (n BoolExpr) expr() {}

type FloatExpr struct {
	Value float64
	Position
}

func (n FloatExpr) expr() {}

type StringExpr struct {
	Value string
	Position
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
	Position
}

func (n SymbolExpr) expr() {}

type UnknowPrimary struct{ Position }

func (n UnknowPrimary) expr() {}

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
	Position
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator lexer.Token
	Right    Expr
	Position
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Right    Expr
	Position
}

func (n AssignmentExpr) expr() {}

type ChainExpr struct {
	Assignee Expr
	Member   Expr
	Position
}

func (n ChainExpr) expr() {}

type StructInstantiationExpr struct {
	StructIdentifier string
	Properties       map[string]Expr
	Position
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Elements []Expr
	Position
}

func (n ArrayInstantiationExpr) expr() {}

type FnCallArg struct {
	Identifier string
	Value      Expr
	Position
}

type FnCallExpr struct {
	Caller    Expr
	Arguments []FnCallArg
	Position
}

func (n FnCallExpr) expr() {}

type FnArg struct {
	Identifier string
	ArgIndex   int
	IsMutable  bool
	Type       Type
	Position
}

type FnDeclareExpr struct {
	Arguments  map[string]FnArg
	Type       Type
	ReturnType Type
	Body       BlockStmt
	Position
}

func (n FnDeclareExpr) expr() {}

type IsTypeExpr struct {
	Left  Expr
	Right Type
	Position
}

func (n IsTypeExpr) expr() {}
