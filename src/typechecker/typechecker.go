package typechecker

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/sanity-io/litter"
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

func (env *env) set(identifer string, value ast.Type, isNew bool, isMutable bool) {
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
		return
	}

	decl, err := env.get(identifer)

	if err != nil {
		panic(err)
	}

	if !decl.IsMutable {
		panic(fmt.Sprintf("%s is not mutable\n", identifer))
	}
	decl.Value = value
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

func check(node any, env *env) ast.Type {
	switch node := node.(type) {
	case ast.BlockStmt:
		scope := createEnv(env)
		for _, stmt := range node.Body {
			check(stmt, scope)
		}
		return ast.CreateUnsetType()
	case ast.DeclarationStmt:
		computed := check(node.AssignedValue, env)

		if !node.Type.IsUnset() && !reflect.DeepEqual(node.Type, computed) {
			set_err(node.Position, fmt.Sprintf("Type %s doesn't match %s", node.Type.Name, computed.Name))
		}

		env.set(node.Identifier, computed, true, node.IsMutable)
		return ast.CreateUnsetType()

	case ast.SymbolExpr:
		val, err := env.get(node.Value)

		if err != nil {
			set_err(node.Position, err.Error())
			return ast.CreateUnsetType()
		}

		return val.Value

	case ast.IntExpr:
		return ast.CreateBaseType(ast.INTEGER)

	case ast.FloatExpr:
		return ast.CreateBaseType(ast.FLOAT)

	case ast.BoolExpr:
		return ast.CreateBaseType(ast.BOOL)

	case ast.BinaryExpr:
		value, err := exec_type_op(node.Operator.Kind, check(node.Left, env), check(node.Right, env))

		if err != nil {
			set_err(node.Position, err.Error())
		}

		return value
	}

	litter.D(node)
	set_err(ast.Position{}, "Node unknown to typechecker :(")
	return ast.CreateUnsetType()
}

var errors = make([]errorhandling.Error, 0)

func set_err(pos ast.Position, message string) {
	errors = append(errors, errorhandling.Error{
		Message:  message,
		Position: pos.Start,
	})
}

func Init(node ast.Stmt) []errorhandling.Error {
	createOpLookup()

	root := createEnv(nil)
	check(node, root)
	return errors
}
