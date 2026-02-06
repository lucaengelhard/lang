package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) _type() {}

type GenericType struct {
	Identifier string
	Arguments  []Type
}

func (t GenericType) _type() {}
