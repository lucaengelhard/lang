package parser

import (
	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

type binding_power int

const (
	default_bp binding_power = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp binding_power) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]binding_power

var stmt_lu = stmt_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var bp_lu = bp_lookup{}

func led(kind lexer.TokenKind, bp binding_power, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	bp_lu[kind] = primary
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = default_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)

	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)

	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.MINUS, additive, parse_binary_expr)
	led(lexer.STAR, multiplicative, parse_binary_expr)
	led(lexer.SLASH, multiplicative, parse_binary_expr)
	led(lexer.PERCENT, multiplicative, parse_binary_expr)

	led(lexer.ASSIGNMENT, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)

	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)
	nud(lexer.MINUS, parser_prefix_expr)
	nud(lexer.OPEN_PAREN, parse_grouping_expr)

	led(lexer.OPEN_CURLY, call, parse_struct_instantiation_expr)
	led(lexer.OPEN_PAREN, call, parse_fn_call_expr)
	nud(lexer.OPEN_BRACKET, parse_array_instantiation_expr)

	stmt(lexer.LET, parse_declartion_stmt)
	stmt(lexer.STRUCT, parse_struct_stmt)
	stmt(lexer.FN, parse_fn_stmt)
	stmt(lexer.IF, parse_if_stmt)
	stmt(lexer.WHILE, parse_while_stmt)
	stmt(lexer.FOR, parse_for_stmt)
	stmt(lexer.RETURN, parse_return_stmt)
	stmt(lexer.CONTINUE, parse_continue_stmt)
	stmt(lexer.BREAK, parse_break_stmt)
}
