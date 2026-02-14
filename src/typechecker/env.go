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

func createEnv(parent *env) *env {
	return &env{
		Parent:       parent,
		Declarations: map[string]*env_decl{},
	}
}
