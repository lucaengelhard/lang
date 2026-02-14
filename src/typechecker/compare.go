package typechecker

import (
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
)

func match(expected, input ast.Type) bool {
	return exec_match_op(expected, input) // || exec_match_op(b, a)
}

type match_op func(a, b ast.Type) bool

var match_lookup = map[string]map[string]match_op{}

func exec_match_op(a, b ast.Type) bool {
	op, exists := match_lookup[a.Name][b.Name]

	if a.Name == ast.UNION && b.Name != ast.UNION {
		for _, union_type := range a.Arguments {
			if match(union_type, b) {
				return true
			}
		}
	}

	if !exists {
		return reflect.DeepEqual(a, b)
	}

	return op(a, b)
}

func create_match_op(a, b string, op match_op) {
	_, a_map_exists := match_lookup[a]

	if !a_map_exists {
		match_lookup[a] = map[string]match_op{}
	}

	match_lookup[a][b] = op
}

func createMatchLookup() {
	create_match_op(ast.DICT, ast.STRUCT, match_dict_struct)
}

func match_dict_struct(input_dict, input_struct ast.Type) bool {

	for _, dict_prop := range input_dict.Arguments {
		var exists = false

		for _, struct_prop := range input_struct.Arguments {
			if dict_prop.Name == struct_prop.Name {
				exists = match(dict_prop, struct_prop)
			}
		}

		if !exists {
			return false
		}
	}

	return true
}
