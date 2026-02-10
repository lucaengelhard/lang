package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
	"github.com/sanity-io/litter"
)

type env_decl struct {
	Identifier string
	Mutable    bool
	Value      any
}

type env struct {
	Declarations map[string]env_decl
	Parent       *env
}

func (env *env) get(identifier string) (*env_decl, error) {
	res, exist := env.Declarations[identifier]

	if exist {
		return &res, nil
	}

	if !exist && env.Parent == nil {
		return &env_decl{}, fmt.Errorf("Variable %s doesn't exist\n", identifier)
	}

	return env.Parent.get(identifier)
}

func (env *env) set(identifer string, value any, isNew bool, isMutable bool) {
	if isNew {
		if _, exists := env.Declarations[identifer]; exists {
			panic(fmt.Sprintf("%s already exist in scope\n", identifer))
		}
		env.Declarations[identifer] = env_decl{
			Identifier: identifer,
			Mutable:    isMutable,
			Value:      value,
		}
		return
	}

	decl, err := env.get(identifer)

	if err != nil {
		panic(err)
	}

	if !decl.Mutable {
		panic(fmt.Sprintf("%s is not mutable\n", identifer))
	}
	decl.Value = value
}

func createEnv(parent *env) *env {
	return &env{
		Parent:       parent,
		Declarations: map[string]env_decl{},
	}
}

func Init(node any) {
	createOpLookup()

	interpret(node, createStdEnv())
}

func interpret(node any, env *env) (any, any) {
	var result any
	var return_value any
	switch node.(type) {
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
		e, _ := lib.ExpectType[ast.ExpressionStmt](node)
		result, _ = interpret(e.Expression, env)
	case ast.IntExpr:
		i, _ := lib.ExpectType[ast.IntExpr](node)
		result = i.Value
	case ast.FloatExpr:
		i, _ := lib.ExpectType[ast.FloatExpr](node)
		result = i.Value
	case ast.BinaryExpr:
		result = interpret_binary_exp(node, env)
	case ast.AssignmentExpr:
		interpret_assignment(node, env)
	case ast.PrefixExpr:
		result = interpret_prefix_expr(node, env)
	case ast.IfStmt:
		return_value = interpret_if_stmt(node, env)
	case ast.ReturnStmt:
		stmt, _ := lib.ExpectType[ast.ReturnStmt](node)
		return_value, _ = interpret(stmt.Value, env)
	default:
		fmt.Printf("Unhandled: %s\n", reflect.TypeOf(node))
		litter.Dump(node)
	}

	return result, return_value
}

func interpret_block(input any, env *env) any {
	block, _ := lib.ExpectType[ast.BlockStmt](input)
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
	declaration, _ := lib.ExpectType[ast.DeclarationStmt](input)
	val, _ := interpret(declaration.AssignedValue, env)
	env.set(declaration.Identifier, val, true, declaration.IsMutable)
}

func interpret_fn_declaration(input any, env *env) func(args ...FnCallArg) any {
	declaration, _ := lib.ExpectType[ast.FnDeclareExpr](input)
	block, _ := lib.ExpectType[ast.BlockStmt](declaration.Body)
	position_arg_map := make([]ast.FnArg, len(declaration.Arguments))

	for _, arg := range declaration.Arguments {
		position_arg_map[arg.Position] = arg
	}

	return func(args ...FnCallArg) any {
		scope := createEnv(env)

		for index, arg := range args {
			if arg.Identifier != "" {
				named_arg, exists := declaration.Arguments[arg.Identifier]

				if !exists {
					panic(fmt.Sprintf("Argument %s doesn't exist on function", arg.Identifier))
				}

				scope.set(named_arg.Identifier, arg.Value, true, named_arg.IsMutable)

				continue
			}

			positional_arg := position_arg_map[index]
			scope.set(positional_arg.Identifier, arg.Value, true, positional_arg.IsMutable)
		}

		for _, stmt := range block.Body {
			_, ret := interpret(stmt, env)
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
}

func interpret_fn_call(input any, env *env) any {
	call, _ := lib.ExpectType[ast.FnCallExpr](input)
	caller_symbol, _ := lib.ExpectType[ast.SymbolExpr](call.Caller)

	declaration, _ := env.get(caller_symbol.Value)
	fn, _ := lib.ExpectType[func(args ...FnCallArg) any](declaration.Value)

	args := make([]FnCallArg, 0)

	for _, arg := range call.Arguments {
		val, _ := interpret(arg.Value, env)
		args = append(args, FnCallArg{
			Identifier: arg.Identifier,
			Value:      val,
		})
	}

	return fn(args...)
}

func interpret_assignment(input any, env *env) {
	assignment, _ := lib.ExpectType[ast.AssignmentExpr](input)
	assignee, _ := lib.ExpectType[ast.SymbolExpr](assignment.Assignee)
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
	symbol, _ := lib.ExpectType[ast.SymbolExpr](input)
	value, err := env.get(symbol.Value)

	if err != nil {
		panic(err)
	}

	return value.Value
}

func interpret_binary_exp(input any, env *env) any {
	expression, _ := lib.ExpectType[ast.BinaryExpr](input)
	left_result, _ := interpret(expression.Left, env)
	right_result, _ := interpret(expression.Right, env)

	return execute_binop(expression.Operator.Kind, left_result, right_result)
}

func interpret_prefix_expr(input any, env *env) any {
	expression, _ := lib.ExpectType[ast.PrefixExpr](input)
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
	stmt, _ := lib.ExpectType[ast.IfStmt](input)

	cond, _ := interpret(stmt.Condition, env)
	decision, _ := lib.ExpectType[bool](cond)

	if decision {
		_, return_value = interpret(stmt.True, env)
	} else {
		_, return_value = interpret(stmt.False, env)
	}

	return return_value
}
