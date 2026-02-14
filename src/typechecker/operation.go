package typechecker

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

var type_binop_lookup = map[lexer.TokenKind]map[string]map[string]ast.Type{}

func exec_type_op(token lexer.TokenKind, left, right ast.Type) (ast.Type, error) {
	err_str := fmt.Sprintf("No type operation for %s and %s\n", left.Name, right.Name)

	tk, exists_tk := type_binop_lookup[token]

	if !exists_tk {
		return ast.CreateUnsetType(), fmt.Errorf("%s", err_str)
	}

	l, exist_l := tk[left.Name]

	if !exist_l {
		return ast.CreateUnsetType(), fmt.Errorf("%s", err_str)

	}

	return_type, exists_r := l[right.Name]

	if !exists_r {
		return ast.CreateUnsetType(), fmt.Errorf("%s", err_str)

	}

	return return_type, nil
}

func create_type_binop(token lexer.TokenKind, left, right, return_type ast.Type) {
	_, token_map_exists := type_binop_lookup[token]

	if !token_map_exists {
		type_binop_lookup[token] = map[string]map[string]ast.Type{}
	}

	_, left_map_exists := type_binop_lookup[token][left.Name]

	if !left_map_exists {
		type_binop_lookup[token][left.Name] = map[string]ast.Type{}
	}

	type_binop_lookup[token][left.Name][right.Name] = return_type
}

func create_commuative_type_binop(token lexer.TokenKind, a, b, return_type ast.Type) {
	create_type_binop(token, a, b, return_type)
	create_type_binop(token, b, a, return_type)
}

func createOpLookup() {
	create_type_binop(lexer.PLUS, ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.INTEGER))
	create_commuative_type_binop(lexer.PLUS, ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.FLOAT), ast.CreateBaseType(ast.FLOAT))
	create_type_binop(lexer.MINUS, ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.INTEGER))
	create_commuative_type_binop(lexer.MINUS, ast.CreateBaseType(ast.INTEGER), ast.CreateBaseType(ast.FLOAT), ast.CreateBaseType(ast.FLOAT))
}
