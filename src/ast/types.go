package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) _type() {}

type ArrayType struct {
	Type Type
}

func (t ArrayType) _type() {}
