package typechecker

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/sanity-io/litter"
)

var errors = make([]errorhandling.Error, 0)

func set_err(pos ast.Position, message string) {
	errors = append(errors, errorhandling.Error{
		Message:  "Type error -> " + message,
		Position: pos.Start,
	})
}

func Init(node ast.Stmt) []errorhandling.Error {
	createOpLookup()
	createHandlerLookup()
	createMatchLookup()

	root := createEnv(nil)
	check(node, root)
	return errors
}

func check(node any, env *env) ast.Type {
	handler, exists := node_handler_lu[reflect.TypeOf(node)]

	if !exists {
		litter.D(node)
		set_err(ast.Position{}, fmt.Sprintf("Node %s unknown to typechecker :(", reflect.TypeOf(node)))
		return ast.CreateUnsetType()
	}

	return handler(node, env)
}
