package ast

import "strings"

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

// TODO: Maybe optional depth argument? So the level of recursion can be set?
func (t Type) ToString() string {
	var arg_string strings.Builder

	if len(t.Arguments) > 0 {
		arg_string.WriteString("<")

		for i, arg := range t.Arguments {
			arg_string.WriteString(arg.ToString())

			if i < len(t.Arguments)-1 {
				arg_string.WriteString(",")
			}
		}

		arg_string.WriteString(">")
	}

	return t.Name + arg_string.String()
}

type Position struct {
	Start int
	End   int
}

func CreatePosition(start int, end int) Position {
	return Position{Start: start, End: end}
}
