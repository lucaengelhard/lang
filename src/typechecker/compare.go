package typechecker

import (
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
)

func match(a, b ast.Type) bool {
	return reflect.DeepEqual(a, b)
}
