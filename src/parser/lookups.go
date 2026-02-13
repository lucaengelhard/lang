package parser

import "github.com/lucaengelhard/lang/src/lexer"

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

type stmt_handler func(p *Parser) Statement
type stmt_map map[lexer.TokenKind]stmt_handler

var stmt_lookup = stmt_map{}

func stmt(kind lexer.TokenKind, handler stmt_handler) {
	stmt_lookup[kind] = handler
}

type nud_handler func(p *Parser) Expression
type nud_map map[lexer.TokenKind]nud_handler

var nud_lookup = nud_map{}

func nud(kind lexer.TokenKind, handler nud_handler) {
	nud_lookup[kind] = handler
}

type bp_map map[lexer.TokenKind]binding_power

var bp_lookup = bp_map{}

type led_handler func(p *Parser, left Expression, bp binding_power) Expression
type led_map map[lexer.TokenKind]led_handler

var led_lookup = led_map{}

func led(kind lexer.TokenKind, handler led_handler, bp binding_power) {
	bp_lookup[kind] = bp
	led_lookup[kind] = handler
}

func createTokenLookups() {
	stmt(lexer.LET, parse_declaration_stmt)

	nud(lexer.TRUE, parse_boolean_expr)
	nud(lexer.FALSE, parse_boolean_expr)
}
