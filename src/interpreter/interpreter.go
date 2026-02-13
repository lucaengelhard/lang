package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/sanity-io/litter"
)

type env_decl struct {
	Identifier  string
	IsMutable   bool
	IsReference bool
	Value       any
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

func (env *env) set(identifer string, value any, isNew bool, isMutable bool) {
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

func Init(node any) {
	createOpLookup()

	interpret(node, createStdEnv())
}

func interpret(node any, env *env) (any, any) {
	var result any
	var return_value any
	switch node := node.(type) {
	case ast.BlockStmt:
		return_value = interpret_block(node, env)
	case ast.DeclarationStmt:
		interpret_declaration(node, env)
	case ast.FnDeclareExpr:
		result = interpret_fn_declaration(node, env)
	case ast.FnCallExpr:
		result = interpret_fn_call(node, env)
	case ast.SymbolExpr:
		result = interpret_symbol_expr(node, env)
	case ast.ExpressionStmt:
		result, _ = interpret(node.Expression, env)
	case ast.IntExpr:
		result = node.Value
	case ast.FloatExpr:
		result = node.Value
	case ast.StringExpr:
		result = node.Value
	case ast.ArrayInstantiationExpr:
		result = interpret_arr_instantiation(node, env)
	case ast.BinaryExpr:
		result = interpret_binary_exp(node, env)
	case ast.AssignmentExpr:
		interpret_assignment(node, env)
	case ast.PrefixExpr:
		result = interpret_prefix_expr(node, env)
	case ast.IfStmt:
		return_value = interpret_if_stmt(node, env)
	case ast.ForStmt:
		return_value = interpret_for_stmt(node, env)
	case ast.WhileStmt:
		return_value = interpret_while_stmt(node, env)
	case ast.ReturnStmt:
		return_value, _ = interpret(node.Value, env)
	default:
		fmt.Printf("Unhandled: %s\n", reflect.TypeOf(node))
		litter.Dump(node)
	}

	return result, return_value
}

func interpret_block(input any, env *env) any {
	block, _ := input.(ast.BlockStmt)
	scope := createEnv(env)
	for _, stmt := range block.Body {
		_, return_value := interpret(stmt, scope)
		if return_value != nil {
			return return_value
		}
	}

	return nil
}

func interpret_declaration(input any, env *env) {
	declaration, _ := input.(ast.DeclarationStmt)
	val, _ := interpret(declaration.AssignedValue, env)
	env.set(declaration.Identifier, val, true, declaration.IsMutable)
}

func interpret_fn_declaration(input any, env *env) func(args ...FnCallArg) any {
	declaration, _ := input.(ast.FnDeclareExpr)
	block := declaration.Body
	position_arg_map := make([]ast.FnArg, len(declaration.Arguments))

	for _, arg := range declaration.Arguments {
		position_arg_map[arg.Position] = arg
	}

	return func(args ...FnCallArg) any {
		scope := createEnv(env)
		var NAMED_ARG_FLAG = false
		for index, passed_arg := range args {
			var definition_arg ast.FnArg

			if passed_arg.Identifier != "" {
				named_arg, exists := declaration.Arguments[passed_arg.Identifier]

				if !exists {
					panic(fmt.Sprintf("Argument %s doesn't exist on function", passed_arg.Identifier))
				}

				definition_arg = named_arg
				NAMED_ARG_FLAG = true
			} else {
				if NAMED_ARG_FLAG {
					panic("Positional arguments not allowed after named arguments have been set")
				}
				definition_arg = position_arg_map[index]
			}

			if definition_arg.IsReference {
				if passed_arg.Reference == nil {
					panic(fmt.Sprintf("Expected reference for argument %s (%v)\n", definition_arg.Identifier, definition_arg.Position))
				}

				if definition_arg.IsMutable && !passed_arg.Reference.IsMutable {
					panic(fmt.Sprintf("Expected mutable reference for argument %s (%v)\n", definition_arg.Identifier, definition_arg.Position))
				}

				scope.set_ref(definition_arg.Identifier, passed_arg.Reference)
			}

			if !definition_arg.IsReference {
				if passed_arg.Reference != nil {
					panic(fmt.Sprintf("Expected argument %s (%v) to be passed by value, got reference\n", definition_arg.Identifier, definition_arg.Position))
				}

				scope.set(definition_arg.Identifier, passed_arg.Value, true, definition_arg.IsMutable)
			}
		}

		for _, stmt := range block.Body {
			_, ret := interpret(stmt, scope)
			if ret != nil {
				return ret
			}
		}
		return nil
	}
}

type FnCallArg struct {
	Identifier string
	Value      any
	Reference  *env_decl
}

func interpret_fn_call(input any, env *env) any {
	call, _ := input.(ast.FnCallExpr)
	caller_symbol := call.Caller.(ast.SymbolExpr)

	declaration, _ := env.get(caller_symbol.Value)
	fn, _ := declaration.Value.(func(args ...FnCallArg) any)

	args := make([]FnCallArg, 0)

	for _, arg := range call.Arguments {
		val, _ := interpret(arg.Value, env)
		var reference *env_decl
		var identifier = arg.Identifier

		switch arg.Value.(type) {
		case ast.SymbolExpr:
			symbol, _ := arg.Value.(ast.SymbolExpr)

			if symbol.IsReference {
				reference, _ = env.get(symbol.Value)
			}
		}

		args = append(args, FnCallArg{
			Identifier: identifier,
			Value:      val,
			Reference:  reference,
		})
	}

	return fn(args...)
}

func interpret_assignment(input any, env *env) {
	assignment, _ := input.(ast.AssignmentExpr)
	assignee, _ := assignment.Assignee.(ast.SymbolExpr)
	right_result, _ := interpret(assignment.Right, env)

	current, error := env.get(assignee.Value)

	if error != nil {
		panic(error)
	}

	op_token, op_token_exists := assignment_operation_lu[assignment.Operator.Kind]

	if op_token_exists {
		env.set(assignee.Value, execute_binop(op_token, current.Value, right_result), false, false)
	} else {
		env.set(assignee.Value, right_result, false, false)
	}

}

func interpret_symbol_expr(input any, env *env) any {
	symbol, _ := input.(ast.SymbolExpr)
	value, err := env.get(symbol.Value)

	if err != nil {
		panic(err)
	}

	return value.Value
}

func interpret_binary_exp(input any, env *env) any {
	expression, _ := input.(ast.BinaryExpr)
	left_result, _ := interpret(expression.Left, env)
	right_result, _ := interpret(expression.Right, env)

	return execute_binop(expression.Operator.Kind, left_result, right_result)
}

func interpret_prefix_expr(input any, env *env) any {
	expression, _ := input.(ast.PrefixExpr)
	right_result, _ := interpret(expression.Right, env)

	switch expression.Operator.Kind {
	case lexer.MINUS:
		return execute_binop(lexer.STAR, int64(-1), right_result)

	default:
		fmt.Printf("Unhandled prefix: %s\n", expression.Operator.Kind.ToString())
		return nil
	}
}

func interpret_if_stmt(input any, env *env) any {
	var return_value any
	stmt, _ := input.(ast.IfStmt)

	cond, _ := interpret(stmt.Condition, env)
	decision, _ := cond.(bool)

	if decision {
		_, return_value = interpret(stmt.True, env)
	} else {
		_, return_value = interpret(stmt.False, env)
	}

	return return_value
}

func interpret_for_stmt(input any, env *env) any {
	stmt, _ := input.(ast.ForStmt)
	scope := createEnv(env)

	interpret(stmt.Assignment, scope)

	var ret any
	for true {
		result, _ := interpret(stmt.Condition, scope)
		condition, _ := result.(bool)

		if !condition {
			break
		}

		_, ret = interpret(stmt.Body, scope)

		if ret != nil {
			break
		}

		interpret(stmt.Increment, scope)

	}

	return ret
}

func interpret_while_stmt(input any, env *env) any {
	stmt, _ := input.(ast.WhileStmt)
	scope := createEnv(env)

	var ret any
	for true {
		result, _ := interpret(stmt.Condition, scope)
		condition, _ := result.(bool)

		if !condition {
			break
		}

		_, ret = interpret(stmt.Body, scope)

		if ret != nil {
			break
		}

	}

	return ret
}

func interpret_arr_instantiation(input any, env *env) any {
	arr, _ := input.(ast.ArrayInstantiationExpr)

	res := make([]any, 0)

	for _, el := range arr.Elements {
		el_res, _ := interpret(el, env)
		res = append(res, el_res)
	}

	return res
}
