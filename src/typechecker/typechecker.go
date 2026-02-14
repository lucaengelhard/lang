package typechecker

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
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

	root := createEnv(nil)
	check(node, root)
	return errors
}

func createHandlerLookup() {
	add_handler(expr_stmt_handler)
	add_handler(block_handler)
	add_handler(symbol_handler)
	add_handler(int_handler)
	add_handler(float_handler)
	add_handler(bool_handler)
	add_handler(binary_expr_handler)
	add_handler(declaration_handler)
	add_handler(assignment_handler)

}

type handler func(node any, env *env) ast.Type
type handler_lookup map[reflect.Type]handler

var node_handler_lu = handler_lookup{}

func add_handler[Node any](handler func(node Node, env *env) ast.Type) {
	_, exists := node_handler_lu[reflect.TypeFor[Node]()]

	if exists {
		set_err(ast.Position{}, fmt.Sprintf("Node handler already exists for node of type %s", reflect.TypeFor[Node]()))
	}

	node_handler_lu[reflect.TypeFor[Node]()] = func(node any, env *env) ast.Type {
		return handler(node.(Node), env)
	}
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

func expr_stmt_handler(node ast.ExpressionStmt, env *env) ast.Type {
	return check(node.Expression, env)
}

func block_handler(node ast.BlockStmt, env *env) ast.Type {
	scope := createEnv(env)
	for _, stmt := range node.Body {
		check(stmt, scope)
	}
	return ast.CreateUnsetType()
}

func symbol_handler(node ast.SymbolExpr, env *env) ast.Type {
	val, err := env.get(node.Value)

	if err != nil {
		set_err(node.Position, err.Error())
		return ast.CreateUnsetType()
	}

	return val.Value
}

func int_handler(node ast.IntExpr, env *env) ast.Type {
	return ast.CreateBaseType(ast.INTEGER)
}

func float_handler(node ast.FloatExpr, env *env) ast.Type {
	return ast.CreateBaseType(ast.FLOAT)
}

func bool_handler(node ast.BoolExpr, env *env) ast.Type {
	return ast.CreateBaseType(ast.BOOL)
}

func binary_expr_handler(node ast.BinaryExpr, env *env) ast.Type {
	value, err := exec_type_op(node.Operator.Kind, check(node.Left, env), check(node.Right, env))

	if err != nil {
		set_err(node.Position, err.Error())
		return ast.CreateUnsetType()
	}

	return value
}

func declaration_handler(node ast.DeclarationStmt, env *env) ast.Type {
	computed := check(node.AssignedValue, env)

	if !node.Type.IsUnset() && !reflect.DeepEqual(node.Type, computed) {
		set_err(node.Position, fmt.Sprintf("Type %s doesn't match %s", node.Type.ToString(), computed.ToString()))
		return ast.CreateUnsetType()
	}

	env.set(node.Identifier, computed, true, node.IsMutable)

	return ast.CreateUnsetType()
}

func assignment_handler(node ast.AssignmentExpr, env *env) ast.Type {
	assignee := node.Assignee.(ast.SymbolExpr)
	right := check(node.Right, env)

	current, err := env.get(assignee.Value)

	if err != nil {
		set_err(node.Position, err.Error())
		return ast.CreateUnsetType()
	}

	op_token, op_token_exists := lexer.Assignment_operation_lu[node.Operator.Kind]

	if op_token_exists {
		computed, err := exec_type_op(op_token, current.Value, right)

		if err != nil {
			set_err(node.Position, err.Error())
		}

		err = env.set(assignee.Value, computed, false, false)

		if err != nil {
			set_err(node.Position, err.Error())
		}
	} else if node.Operator.Kind == lexer.ASSIGNMENT {
		err = env.set(assignee.Value, right, false, false)

		if err != nil {
			set_err(node.Position, err.Error())
		}
	} else {
		set_err(node.Position, fmt.Sprintf("Unknown assignment operator: %s", node.Operator.Kind.ToString()))
	}

	return ast.CreateUnsetType()
}
