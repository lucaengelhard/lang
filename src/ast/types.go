package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) _type() {}

type FnType struct {
	Arguments  map[string]FnArg
	ReturnType Type
}

func (t FnType) _type() {}

type GenericType struct {
	Identifier string
	Arguments  []Type
}

func (t GenericType) _type() {}

type IsType struct {
	Left  Type
	Right Type
}

func (t IsType) _type() {}
