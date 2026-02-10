package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
)

type operation func(l, r any) any
type op_lookup map[lexer.TokenKind]map[reflect.Type]map[reflect.Type]operation

var op_lu = op_lookup{}

func create_op[L any, R any, Ret any](token lexer.TokenKind, op func(l L, r R) Ret) {
	_, tk_map_exists := op_lu[token]

	if !tk_map_exists {
		op_lu[token] = map[reflect.Type]map[reflect.Type]operation{}
	}

	_, l_map_exists := op_lu[token][reflect.TypeFor[L]()]

	if !l_map_exists {
		op_lu[token][reflect.TypeFor[L]()] = map[reflect.Type]operation{}
	}

	op_lu[token][reflect.TypeFor[L]()][reflect.TypeFor[R]()] = func(l, r any) any {
		valid_l, _ := lib.ExpectType[L](l)
		valid_r, _ := lib.ExpectType[R](r)

		return op(valid_l, valid_r)
	}
}

func get_op(token lexer.TokenKind, left any, right any) operation {
	err_str := fmt.Sprintf("No operation for %s and %s\n", reflect.TypeOf(left), reflect.TypeOf(right))
	tk, exists_tk := op_lu[token]

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

func execute_op(token lexer.TokenKind, left any, right any) any {
	return get_op(token, left, right)(left, right)
}

func createOpLookup() {
	create_op(lexer.PLUS, int_add)
	with_cast(lexer.PLUS, float_add, int_to_float)
	create_op(lexer.MINUS, int_sub)
	with_cast(lexer.MINUS, float_minus, int_to_float)
	create_op(lexer.STAR, int_mult)
	with_cast(lexer.STAR, float_mult, int_to_float)
	create_op(lexer.SLASH, int_div)
	with_cast(lexer.SLASH, float_div, int_to_float)
	create_op(lexer.PERCENT, int_mod)
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

func with_cast[From any, To any](token lexer.TokenKind, op func(l To, r To) To, cast func(f From) To) {
	no_cast := func(l To, r To) To {
		return op(l, r)
	}

	left := func(l From, r To) To {
		return op(cast(l), r)
	}

	right := func(l To, r From) To {
		return op(l, cast(r))
	}
	create_op(token, no_cast)
	create_op(token, left)
	create_op(token, right)
}
