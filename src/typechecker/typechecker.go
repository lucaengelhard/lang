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
	createMatchLookup()

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
	add_handler(string_handler)
	add_handler(binary_expr_handler)
	add_handler(declaration_handler)
	add_handler(assignment_handler)
	add_handler(array_instantiation_handler)
	add_handler(interface_handler)
	add_handler(struct_stmt_handler)
	add_handler(struct_instantiation_handler)

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

func string_handler(node ast.StringExpr, env *env) ast.Type {
	return ast.CreateBaseType(ast.STRING)
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

	// TODO: make more sophisticated equality check, so that order of array doesn't matter for example
	// Also partial matching doesn't work
	if !node.Type.IsUnset() && !match(env.get_type(node.Type.Name), computed) {
		set_err(node.Position, fmt.Sprintf("Type %s doesn't match %s (%s)", computed.ToString(), node.Type.ToString(), env.get_type(node.Type.Name).ToString()))
		return ast.CreateUnsetType()
	}

	if node.Type.IsUnset() {
		env.set(node.Identifier, computed, true, node.IsMutable)
	} else {
		env.set(node.Identifier, env.get_type(node.Type.Name), true, node.IsMutable)
	}

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

func array_instantiation_handler(node ast.ArrayInstantiationExpr, env *env) ast.Type {
	elements := make([]ast.Type, 0)

	for _, el := range node.Elements {
		computed := check(el, env)
		var exists = false

		for _, already_existing := range elements {
			if reflect.DeepEqual(already_existing, computed) {
				exists = true
			}
		}

		if !exists {
			elements = append(elements, computed)
		}
	}

	return ast.Type{
		Name:      ast.ARRAY,
		Arguments: elements,
	}
}

func wrap_property_type(identifer string, prop_type ast.Type) ast.Type {
	return ast.Type{Name: identifer, Arguments: []ast.Type{prop_type}}
}

func interface_handler(node ast.InterfaceStmt, env *env) ast.Type {
	if !node.SingleType.IsUnset() {
		env.set_type(node.Identifier, node.SingleType)
	} else {
		properties := make([]ast.Type, 0)

		for _, prop := range node.StructType {
			properties = append(properties, wrap_property_type(prop.Name, prop.Type))
		}

		env.set_type(node.Identifier, ast.Type{
			Name:      ast.DICT,
			Arguments: properties,
		})
	}

	return ast.CreateUnsetType()
}

func struct_stmt_handler(node ast.StructStmt, env *env) ast.Type {
	properties := make([]ast.Type, 0)

	for _, prop := range node.Properties {
		properties = append(properties, wrap_property_type(prop.Name, prop.Type))
	}

	env.set_type(node.Identifier, ast.Type{
		Name:      ast.STRUCT,
		Arguments: properties,
	})

	return ast.CreateUnsetType()
}

func struct_instantiation_handler(node ast.StructInstantiationExpr, env *env) ast.Type {
	struct_type := env.get_type(node.StructIdentifier)

	if struct_type.IsUnset() {
		return struct_type
	}

	for _, prop_type := range struct_type.Arguments {
		prop_val, exists := node.Properties[prop_type.Name]

		if !exists {
			set_err(node.Position, fmt.Sprintf("Property %s missing on struct", prop_type.Name))
			return struct_type
		}

		computed := wrap_property_type(prop_type.Name, check(prop_val, env))

		if !match(prop_type, computed) {
			set_err(node.Position, fmt.Sprintf("Property %s expected %s but got %s", prop_type.Name, prop_type.ToString(), computed.ToString()))
			return struct_type
		}

	}

	return struct_type
}
