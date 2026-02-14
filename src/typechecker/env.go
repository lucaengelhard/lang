package typechecker

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
)

type env_decl struct {
	Identifier  string
	IsMutable   bool
	IsReference bool
	Value       ast.Type
}

type env struct {
	Declarations map[string]*env_decl
	Parent       *env
	Types        map[string]ast.Type
}

func (env *env) get(identifier string) (*env_decl, error) {
	res, exist := env.Declarations[identifier]

	if exist {
		return res, nil
	}

	if !exist && env.Parent == nil {
		return &env_decl{}, fmt.Errorf("Variable %s doesn't exist\n", identifier)
	}

	return env.Parent.get(identifier)
}

func (env *env) set(identifer string, value ast.Type, isNew bool, isMutable bool) error {
	if isNew {
		if _, exists := env.Declarations[identifer]; exists {
			panic(fmt.Sprintf("%s already exists in scope\n", identifer))
		}

		env.Declarations[identifer] = &env_decl{
			Identifier:  identifer,
			IsMutable:   isMutable,
			IsReference: false,
			Value:       value,
		}
		return nil
	}

	decl, err := env.get(identifer)

	if err != nil {
		return err
	}

	if !decl.IsMutable {
		return fmt.Errorf("%s is not mutable\n", identifer)
	}
	decl.Value = value
	return nil
}

func (env *env) set_ref(identifer string, ref *env_decl) {
	env.Declarations[identifer] = ref
}

func (env *env) get_root() *env {
	if env.Parent == nil {
		return env
	}

	return env.Parent.get_root()
}

func (env *env) set_type(identifer string, t ast.Type) {
	root := env.get_root()

	_, exists := root.Types[identifer]

	if exists {
		set_err(ast.Position{}, fmt.Sprintf("Type %s already exists", identifer))
		return
	}

	root.Types[identifer] = t
}

func (env *env) get_type(identifer string) ast.Type {
	root := env.get_root()

	t, exists := root.Types[identifer]

	if !exists {
		set_err(ast.Position{}, fmt.Sprintf("Type %s doesn't exist", identifer))
		return ast.CreateUnsetType()
	}

	return t
}

func createEnv(parent *env) *env {
	var types map[string]ast.Type

	if parent == nil {
		types = map[string]ast.Type{}
	}

	return &env{
		Parent:       parent,
		Declarations: map[string]*env_decl{},
		Types:        types,
	}
}
