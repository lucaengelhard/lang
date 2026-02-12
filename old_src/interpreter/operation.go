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

func create_binop_with_cast[From any, To any, Return any](token lexer.TokenKind, op func(l To, r To) Return, cast func(f From) To) {
	no_cast := func(l To, r To) Return {
		return op(l, r)
	}

	left := func(l From, r To) Return {
		return op(cast(l), r)
	}

	right := func(l To, r From) Return {
		return op(l, cast(r))
	}
	create_binop(token, no_cast)
	create_binop(token, left)
	create_binop(token, right)
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
	create_binop(lexer.PLUS, add[string])
	create_binop_with_cast(lexer.PLUS, add[string], int_to_str)
	create_binop(lexer.PLUS, add[int64])
	create_binop_with_cast(lexer.PLUS, add[float64], int_to_float)
	create_binop(lexer.MINUS, sub[int64])
	create_binop_with_cast(lexer.MINUS, sub, int_to_float) // Also create for not casted?
	create_binop(lexer.STAR, mult[int64])
	create_binop_with_cast(lexer.STAR, mult[float64], int_to_float)
	create_binop(lexer.SLASH, div[int64])
	create_binop_with_cast(lexer.SLASH, div[float64], int_to_float)
	create_binop(lexer.PERCENT, mod[int64])

	create_binop(lexer.EQUALS, eq[int64])
	create_binop(lexer.EQUALS, eq[float64])
	create_binop(lexer.GREATER, greater[int64])
	create_binop_with_cast(lexer.GREATER, greater[float64], int_to_float)
	create_binop(lexer.LESS, lesser[int64])
	create_binop_with_cast(lexer.LESS, lesser[float64], int_to_float)
	create_binop(lexer.GREATER_EQUALS, greater_eq[int64])
	create_binop_with_cast(lexer.GREATER_EQUALS, greater_eq[float64], int_to_float)
	create_binop(lexer.LESS_EQUALS, lesser_eq[int64])
	create_binop_with_cast(lexer.LESS_EQUALS, lesser_eq[float64], int_to_float)

	create_binop(lexer.OR, or)
	create_binop(lexer.AND, and)

	assignment_operation_lu[lexer.PLUS_EQUALS] = lexer.PLUS
	assignment_operation_lu[lexer.MINUS_EQUALS] = lexer.MINUS
	assignment_operation_lu[lexer.PLUS_PLUS] = lexer.PLUS
	assignment_operation_lu[lexer.MINUS_MINUS] = lexer.MINUS
}

func int_to_float(input int64) float64 {
	return float64(input)
}

func int_to_str(input int64) string {
	return fmt.Sprint(input)
}

func add[T lib.Arithmetic | ~string](l T, r T) T {
	return l + r
}

func sub[T lib.Arithmetic](l T, r T) T {
	return l - r
}

func mult[T lib.Arithmetic](l T, r T) T {
	return l * r
}

func div[T lib.Arithmetic](l T, r T) T {
	return l / r
}

func mod[T lib.Int](l T, r T) T {
	return l & r
}

func eq[T lib.Compareable](l T, r T) bool {
	return l == r
}

func greater[T lib.Orderable](l T, r T) bool {
	return l > r
}

func lesser[T lib.Orderable](l T, r T) bool {
	return l < r
}

func greater_eq[T lib.Orderable](l T, r T) bool {
	return l >= r
}

func lesser_eq[T lib.Orderable](l T, r T) bool {
	return l <= r
}

func and(l bool, r bool) bool {
	return l && r
}

func or(l bool, r bool) bool {
	return l || r
}
