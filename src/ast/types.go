package ast

type Type struct {
	Name      string
	Arguments []Type
}

const (
	UNSET_TYPE = "__unset__"
	INTEGER    = "int"
	FLOAT      = "float"
	BOOL       = "bool"
	FUNCTION   = "func"
)

func CreateUnsetType() Type {
	return CreateBaseType(UNSET_TYPE)
}

func CreateBaseType(name string) Type {
	return Type{
		Name:      name,
		Arguments: make([]Type, 0),
	}
}

func (t Type) IsUnset() bool {
	return t.Name == UNSET_TYPE
}

type Position struct {
	Start int
	End   int
}

func CreatePosition(start int, end int) Position {
	return Position{Start: start, End: end}
}
