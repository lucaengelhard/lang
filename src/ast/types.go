package ast

type SymbolType struct {
	Value string
}

func (t SymbolType) _type() {}

type IntLiteralType struct {
	Value int64
}

func (t IntLiteralType) _type() {}

type StringLiteralType struct {
	Value string
}

func (t StringLiteralType) _type() {}

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

type UnkownType struct{}

func (t UnkownType) _type() {}

type BlockType struct {
	Body []Type
}

func (t BlockType) _type() {}
