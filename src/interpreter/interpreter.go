package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lib"
)

type env struct {
	Declarations map[string]any
	Parent       *env
}

func (env *env) getDeclaration(identifier string) (any, error) {
	res, exist := env.Declarations[identifier]

	if exist {
		return res, nil
	}

	if !exist && env.Parent == nil {
		return ast.UnknowPrimary{}, fmt.Errorf("Type %s doesn't exist", identifier)
	}

	return env.Parent.getDeclaration(identifier)
}

func createEnv(parent *env) *env {
	return &env{
		Parent:       parent,
		Declarations: map[string]any{},
	}
}

func Init(node any) {
	createOpLookup()
	interpret(node, nil)
}

func interpret(node any, env *env) any {
	var result any
	switch node.(type) {
	case ast.BlockStmt:
		interpret_block(node, env)
	case ast.DeclarationStmt:
		interpret_declaration(node, env)
	case ast.SymbolExpr:
		result = interpret_symbol_expr(node, env)
	case ast.ExpressionStmt:
		e, _ := lib.ExpectType[ast.ExpressionStmt](node)
		result = interpret(e.Expression, env)
	case ast.IntExpr:
		i, _ := lib.ExpectType[ast.IntExpr](node)
		result = i.Value
	case ast.FloatExpr:
		i, _ := lib.ExpectType[ast.FloatExpr](node)
		result = i.Value
	case ast.BinaryExpr:
		result = interpret_binary_exp(node, env)
	default:
		fmt.Printf("Unhandled: %s\n", reflect.TypeOf(node))
	}

	fmt.Println(result)

	return result
}

func interpret_block(block any, env *env) {
	b, _ := lib.ExpectType[ast.BlockStmt](block)
	scope := createEnv(env)
	for _, stmt := range b.Body {
		interpret(stmt, scope)
	}
}

func interpret_declaration(decl any, env *env) {
	d, _ := lib.ExpectType[ast.DeclarationStmt](decl)

	env.Declarations[d.Identifier] = interpret(d.AssignedValue, env)
}

func interpret_symbol_expr(exp any, env *env) any {
	s, _ := lib.ExpectType[ast.SymbolExpr](exp)
	v, _ := env.Declarations[s.Value]
	return v
}

func interpret_binary_exp(expression any, env *env) any {
	e, _ := lib.ExpectType[ast.BinaryExpr](expression)
	left_result := interpret(e.Left, env)
	right_esult := interpret(e.Right, env)
	return execute_op(e.Operator.Kind, left_result, right_esult)
}
