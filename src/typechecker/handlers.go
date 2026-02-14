package typechecker

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

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
	add_handler(fn_declare_handler)
	add_handler(return_handler)
	add_handler(if_handler)
	add_handler(fn_call_handler)
	add_handler(deref_handler)
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

func expr_stmt_handler(node ast.ExpressionStmt, env *env) ast.Type {
	return check(node.Expression, env)
}

func block_handler(node ast.BlockStmt, env *env) ast.Type {
	scope := createEnv(env)
	var return_type = ast.CreateUnsetType()
	for _, stmt := range node.Body {
		_, isReturn := stmt.(ast.ReturnStmt)
		_, isIf := stmt.(ast.IfStmt)
		computed := check(stmt, scope)
		if isReturn || isIf {
			return_type = computed
		}
	}

	return return_type
}

func symbol_handler(node ast.SymbolExpr, env *env) ast.Type {
	val, err := env.get(node.Value)

	if err != nil {
		set_err(node.Position, err.Error())
		return ast.CreateUnsetType()
	}

	if node.IsReference {
		return ast.Type{Name: ast.REFERENCE, Arguments: []ast.Type{val.Value}}
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
	var assigned_type = ast.CreateUnsetType()

	// TODO: make more sophisticated equality check, so that order of array doesn't matter for example
	// Also partial matching doesn't work
	if !node.Type.IsUnset() && !match(env.get_type(node.Type.Name), computed) {
		set_err(node.Position, fmt.Sprintf("Type %s doesn't match %s (%s)", computed.ToString(), node.Type.ToString(), env.get_type(node.Type.Name).ToString()))
		return ast.CreateUnsetType()
	}

	if node.Type.IsUnset() {
		assigned_type = computed
	} else {
		assigned_type = env.get_type(node.Type.Name)
	}

	if node.IsMutable {
		assigned_type = assigned_type.Mutable()
	}

	env.set(node.Identifier, assigned_type, true, node.IsMutable)

	return ast.CreateUnsetType()
}

func assignment_handler(node ast.AssignmentExpr, env *env) ast.Type {
	assignee := node.Assignee.(ast.SymbolExpr)
	current_declaration, err := env.get(assignee.Value)

	if err != nil {
		set_err(node.Position, err.Error())
		return ast.CreateUnsetType()
	}

	if !current_declaration.Value.Is(ast.MUTABLE) {
		set_err(node.Position, fmt.Sprintf("%s is not mutable", assignee.Value))
		return ast.CreateUnsetType()
	}

	stripped_current := current_declaration.Value.Arguments[0]

	right := check(node.Right, env)

	op_token, op_token_exists := lexer.Assignment_operation_lu[node.Operator.Kind]

	if op_token_exists {
		computed, err := exec_type_op(op_token, stripped_current, right)

		if err != nil {
			set_err(node.Position, fmt.Sprintf("Type %s is not assignable to variable of type %s (%s)", right.ToString(), stripped_current.ToString(), err.Error()))
			return ast.CreateUnsetType()
		}

		err = env.set(assignee.Value, computed.Mutable(), false, false)

		if err != nil {
			set_err(node.Position, err.Error())
		}
	} else if node.Operator.Kind == lexer.ASSIGNMENT {
		err = env.set(assignee.Value, right.Mutable(), false, false)

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

func fn_declare_handler(node ast.FnDeclareExpr, env *env) ast.Type {
	args := make([]ast.Type, 0)
	scope := createEnv(env)
	var return_type = node.ReturnType

	for _, arg := range node.Arguments {
		args = append(args, wrap_property_type(arg.Identifier, arg.Type))
		scope.set(arg.Identifier, arg.Type, true, arg.IsMutable)
	}

	computed_return_type := check(node.Body, scope)

	if !return_type.IsUnset() && !match(return_type, computed_return_type) {
		set_err(node.Position, fmt.Sprintf("Type %s doesn't match %s (%s)", computed_return_type.ToString(), return_type.ToString(), env.get_type(return_type.Name).ToString()))
	} else {
		return_type = computed_return_type
	}

	return ast.Type{
		Name: ast.FUNCTION,
		Arguments: []ast.Type{
			{Name: ast.FUNCTION_ARG, Arguments: args},
			wrap_property_type(ast.FUNCTION_RETURN, return_type),
		},
	}
}

func fn_call_handler(node ast.FnCallExpr, env *env) ast.Type {
	caller, _ := node.Caller.(ast.SymbolExpr)
	declaration, err := env.get(caller.Value)
	var return_type = ast.CreateUnsetType()

	if err != nil {
		set_err(node.Position, fmt.Sprintf("%s not found", caller.Value))
		return ast.CreateUnsetType()
	}

	if declaration.Value.Name != ast.FUNCTION {
		set_err(node.Position, fmt.Sprintf("%s not a function", caller.Value))
		return ast.CreateUnsetType()
	}

	for _, type_arg := range declaration.Value.Arguments {
		if type_arg.Name == ast.FUNCTION_ARG {
			if len(type_arg.Arguments) < len(node.Arguments) {
				set_err(node.Position, fmt.Sprintf("Too many arguments. Expected %d, got %d", len(type_arg.Arguments), len(node.Arguments)))

			}

			if len(type_arg.Arguments) > len(node.Arguments) {
				set_err(node.Position, fmt.Sprintf("Missing arguments. Expected %d, got %d", len(type_arg.Arguments), len(node.Arguments)))
			}

			// TODO: Handle named args
			for index, arg := range node.Arguments {
				expected := type_arg.Arguments[index].Arguments[0]
				computed := check(arg.Value, env)

				if !match(expected, computed) {
					set_err(node.Position, fmt.Sprintf("Mismatched argument (%d). Expected %s, got %s", index, expected.ToString(), computed.ToString()))
				}
			}
		}

		if type_arg.Name == ast.FUNCTION_RETURN {
			return_type = type_arg.Arguments[0]
		}
	}

	return return_type
}

func return_handler(node ast.ReturnStmt, env *env) ast.Type {
	return check(node.Value, env)
}

func deref_handler(node ast.DerefExpr, env *env) ast.Type {
	ref := check(node.Ref, env)

	if ref.Name != ast.REFERENCE {
		set_err(node.Position, fmt.Sprintf("Can't deference a variable that's not a reference (%s)", ref.ToString()))
		return ast.CreateUnsetType()
	}

	return ref.Arguments[0]
}

func if_handler(node ast.IfStmt, env *env) ast.Type {
	true_return := check(node.True, env)
	false_return := check(node.False, env)

	return ast.Type{Name: ast.UNION, Arguments: []ast.Type{true_return, false_return}}
}
