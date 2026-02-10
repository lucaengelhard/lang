package interpreter

import (
	"fmt"
	"reflect"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
	"github.com/sanity-io/litter"
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
	case ast.AssignmentExpr:
		interpret_assignment(node, env)
	case ast.PrefixExpr:
		result = interpret_prefix_expr(node, env)
	default:
		fmt.Printf("Unhandled: %s\n", reflect.TypeOf(node))
		litter.Dump(node)

	}

	fmt.Println(result)

	return result
}

func interpret_block(input any, env *env) {
	block, _ := lib.ExpectType[ast.BlockStmt](input)
	scope := createEnv(env)
	for _, stmt := range block.Body {
		interpret(stmt, scope)
	}
}

func interpret_declaration(input any, env *env) {
	declaration, _ := lib.ExpectType[ast.DeclarationStmt](input)

	env.Declarations[declaration.Identifier] = interpret(declaration.AssignedValue, env)
}

func interpret_assignment(input any, env *env) {
	assignment, _ := lib.ExpectType[ast.AssignmentExpr](input)
	assignee, _ := lib.ExpectType[ast.SymbolExpr](assignment.Assignee)
	right_result := interpret(assignment.Right, env)

	current, exists := env.Declarations[assignee.Value]

	if !exists {
		panic(fmt.Sprintf("Variable %s doesn't exist in the current scope\n", assignee.Value))
	}

	op_token, op_token_exists := assignment_operation_lu[assignment.Operator.Kind]

	if op_token_exists {
		env.Declarations[assignee.Value] = execute_binop(op_token, current, right_result)
	} else {
		env.Declarations[assignee.Value] = right_result
	}
}

func interpret_symbol_expr(input any, env *env) any {
	symbol, _ := lib.ExpectType[ast.SymbolExpr](input)
	value, _ := env.Declarations[symbol.Value]
	return value
}

func interpret_binary_exp(input any, env *env) any {
	expression, _ := lib.ExpectType[ast.BinaryExpr](input)
	left_result := interpret(expression.Left, env)
	right_esult := interpret(expression.Right, env)
	return execute_binop(expression.Operator.Kind, left_result, right_esult)
}

func interpret_prefix_expr(input any, env *env) any {
	expression, _ := lib.ExpectType[ast.PrefixExpr](input)
	right_result := interpret(expression.Right, env)

	switch expression.Operator.Kind {
	case lexer.MINUS:
		return execute_binop(lexer.STAR, int64(-1), right_result)

	default:
		fmt.Printf("Unhandled prefix: %s\n", expression.Operator.Kind.ToString())
		return nil
	}
}
