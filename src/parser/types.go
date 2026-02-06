package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp binding_power) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]binding_power

var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}
var type_bp_lu = type_bp_lookup{}

func type_led(kind lexer.TokenKind, bp binding_power, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func type_nud(kind lexer.TokenKind, nud_fn type_nud_handler) {
	type_bp_lu[kind] = primary
	type_nud_lu[kind] = nud_fn
}

func createTokenTypeLookups() {
	type_nud(lexer.IDENTIFIER, parse_symbol_type)
	type_nud(lexer.OPEN_BRACKET, parse_array_type)
}

func parse_type(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		p.panic(fmt.Sprintf("Type Nud handler expected for token %s\n", tokenKind.ToString()))
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			p.panic(fmt.Sprintf("Type Led handler expected for token %s\n", tokenKind.ToString()))
		}

		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_symbol_type(p *parser) ast.Type {
	return ast.SymbolType{Value: p.expect(lexer.IDENTIFIER).Value}
}

// remove explicit array type but modifiy symbol type so its generic
func parse_array_type(p *parser) ast.Type {
	p.advance()
	p.expect(lexer.CLOSE_BRACKET)
	var arrayType = parse_type(p, default_bp)
	return ast.ArrayType{
		Type: arrayType,
	}
}
