package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
)

type binop func(l, r any) any
type binop_lookup map[lexer.TokenKind]map[reflect.Type]map[reflect.Type]binop

var binop_lu = binop_lookup{}
var assignment_operation_lu = map[lexer.TokenKind]lexer.TokenKind{}

func create_binop[L any, R any, Ret any](token lexer.TokenKind, op func(l L, r R) Ret) {
	_, tk_map_exists := binop_lu[token]

	if !tk_map_exists {
		binop_lu[token] = map[reflect.Type]map[reflect.Type]binop{}
	}

	_, l_map_exists := binop_lu[token][reflect.TypeFor[L]()]

	if !l_map_exists {
		binop_lu[token][reflect.TypeFor[L]()] = map[reflect.Type]binop{}
	}

	binop_lu[token][reflect.TypeFor[L]()][reflect.TypeFor[R]()] = func(l, r any) any {
		valid_l, _ := lib.ExpectType[L](l)
		valid_r, _ := lib.ExpectType[R](r)

		return op(valid_l, valid_r)
	}
}

func get_op(token lexer.TokenKind, left any, right any) binop {
	err_str := fmt.Sprintf("No operation for %s and %s\n", reflect.TypeOf(left), reflect.TypeOf(right))
	tk, exists_tk := binop_lu[token]

	if !exists_tk {
		panic(err_str)
	}

	l, exist_l := tk[reflect.TypeOf(left)]

	if !exist_l {
		panic(err_str)
	}

	op, exists_r := l[reflect.TypeOf(right)]

	if !exists_r {
		panic(err_str)
	}

	return op
}

func execute_binop(token lexer.TokenKind, left any, right any) any {
	return get_op(token, left, right)(left, right)
}

func createOpLookup() {
	create_binop(lexer.PLUS, int_add)
	create_binop_with_cast(lexer.PLUS, float_add, int_to_float)
	create_binop(lexer.MINUS, int_sub)
	create_binop_with_cast(lexer.MINUS, float_minus, int_to_float)
	create_binop(lexer.STAR, int_mult)
	create_binop_with_cast(lexer.STAR, float_mult, int_to_float)
	create_binop(lexer.SLASH, int_div)
	create_binop_with_cast(lexer.SLASH, float_div, int_to_float)
	create_binop(lexer.PERCENT, int_mod)

	assignment_operation_lu[lexer.PLUS_EQUALS] = lexer.PLUS
	assignment_operation_lu[lexer.MINUS_EQUALS] = lexer.MINUS
	assignment_operation_lu[lexer.PLUS_PLUS] = lexer.PLUS
	assignment_operation_lu[lexer.MINUS_MINUS] = lexer.MINUS
}

func int_add(l int64, r int64) int64 {
	return l + r
}

func int_sub(l int64, r int64) int64 {
	return l - r
}

func int_mult(l int64, r int64) int64 {
	return l * r
}

func int_div(l int64, r int64) int64 {
	return l / r
}

func int_mod(l int64, r int64) int64 {
	return l % r
}

func float_add(l float64, r float64) float64 {
	return l + r
}

func float_minus(l float64, r float64) float64 {
	return l - r
}

func float_mult(l float64, r float64) float64 {
	return l * r
}

func float_div(l float64, r float64) float64 {
	return l / r
}

func int_to_float(input int64) float64 {
	return float64(input)
}

func create_binop_with_cast[From any, To any](token lexer.TokenKind, op func(l To, r To) To, cast func(f From) To) {
	no_cast := func(l To, r To) To {
		return op(l, r)
	}

	left := func(l From, r To) To {
		return op(cast(l), r)
	}

	right := func(l To, r From) To {
		return op(l, cast(r))
	}
	create_binop(token, no_cast)
	create_binop(token, left)
	create_binop(token, right)
}
