package typecheck

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lib"
	"github.com/sanity-io/litter"
)

type env struct {
	Types        map[string]ast.Type
	Declarations map[string]ast.Type
	Parent       *env
}

func (env *env) getType(identifier string) (ast.Type, error) {
	res, exist := env.Types[identifier]

	if exist {
		return res, nil
	}

	if !exist {
		return ast.UnkownType{}, fmt.Errorf("Type %s doesn't exist", identifier)
	}

	return env.Parent.getType(identifier)
}

func (env *env) getDeclaration(identifier string) (ast.Type, error) {
	res, exist := env.Types[identifier]

	if exist {
		return res, nil
	}

	if !exist {
		return ast.UnkownType{}, fmt.Errorf("Type %s doesn't exist", identifier)
	}

	return env.Parent.getDeclaration(identifier)
}

func createEnv(parent *env) *env {
	return &env{
		Parent:       parent,
		Types:        map[string]ast.Type{},
		Declarations: map[string]ast.Type{},
	}
}

func check_handler[T any](node any, env *env, callback func(node T, env *env) (ast.Type, []error)) (ast.Type, []error) {
	errors := make([]error, 0)
	stmt, type_err := lib.ExpectType[T](node)
	res, cb_err := callback(stmt, env)
	if type_err != nil {
		errors = append(errors, type_err)
	}
	if cb_err != nil {
		errors = append(errors, cb_err...)
	}

	return res, errors
}

func compare_types(left ast.Type, right ast.Type) bool {
	litter.Dump(left)
	litter.Dump(right)

	return true
}

func Check(input any, env *env) (ast.Type, []error) {
	errors := make([]error, 0)
	var returnType ast.Type
	switch input.(type) {
	case ast.DeclarationStmt:
		res, err := check_handler(input, env, CheckDeclarationStmt)
		errors = append(errors, err...)
		returnType = res
	case ast.BlockStmt:
		res, err := check_handler(input, env, CheckBlockStmt)
		errors = append(errors, err...)
		returnType = res
	case ast.IntExpr:
		res, err := check_handler(input, env, CheckIntExpr)
		errors = append(errors, err...)
		returnType = res
	default:
		errors = append(errors, fmt.Errorf("Unhandeled type %s: ", reflect.TypeOf(input)))
	}

	return returnType, errors
}

func CheckBlockStmt(block ast.BlockStmt, env *env) (ast.Type, []error) {
	errors := make([]error, 0)
	body := make([]ast.Type, 0)

	for _, stmt := range block.Body {
		res, err := Check(stmt, createEnv(env))
		body = append(body, res)
		errors = append(errors, err...)
	}

	return ast.BlockType{
		Body: body,
	}, errors
}

func CheckDeclarationStmt(decl ast.DeclarationStmt, env *env) (ast.Type, []error) {
	errors := make([]error, 0)

	calculatedType, assignedErr := Check(decl.AssignedValue, env)

	if assignedErr != nil {
		errors = append(errors, assignedErr...)
	}

	if decl.Type != nil && !compare_types(decl.Type, calculatedType) {
		errors = append(errors, fmt.Errorf("Mismatched types")) // TODO: Make Handling better
		return ast.UnkownType{}, errors
	}

	env.Declarations[decl.Identifier] = calculatedType

	return calculatedType, errors
}

func CheckIntExpr(i ast.IntExpr, env *env) (ast.Type, []error) {
	return ast.IntLiteralType{
		Value: i.Value,
	}, make([]error, 0)
}
