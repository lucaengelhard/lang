package parser

import (
	"fmt"
	"reflect"
	"strconv"

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
	type_nud(lexer.NUMBER, parse_number_type)
	type_nud(lexer.STRING, parse_string_type)
	type_nud(lexer.OPEN_PAREN, parse_fn_type)
	type_led(lexer.LESS, call, parse_generic_type)
	type_led(lexer.IS, logical, parse_is_type)
}

func parse_type(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		p.err(fmt.Sprintf("Type Nud handler expected for token %s\n", tokenKind.ToString()))
		return nil
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			p.err(fmt.Sprintf("Type Led handler expected for token %s\n", tokenKind.ToString()))
			return nil
		}

		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_symbol_type(p *parser) ast.Type {
	return ast.SymbolType{Value: p.expect(lexer.IDENTIFIER).Literal}
}

func parse_string_type(p *parser) ast.Type {
	return ast.StringLiteralType{Value: p.expect(lexer.STRING).Literal}
}

func parse_number_type(p *parser) ast.Type {
	val := p.expect(lexer.NUMBER).Literal
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return ast.IntLiteralType{
			Value: i,
		}
	}

	p.err(fmt.Sprintf("Only integers allowed in Number Types %s", val))
	return ast.UnkownType{}
}

func parse_generic_type(p *parser, left ast.Type, bp binding_power) ast.Type {
	symbol, ok := left.(ast.SymbolType)

	if !ok {
		p.err(fmt.Sprintf("Type error: Expected %s got %s", reflect.TypeFor[ast.SymbolExpr](), reflect.TypeOf(symbol)))
		return nil
	}

	identifier := symbol.Value

	var arguments = []ast.Type{}
	p.expect(lexer.LESS)

	for p.hasTokens() && p.currentTokenKind() != lexer.GREATER {
		typeArg := parse_type(p, logical)
		arguments = append(arguments, typeArg)

		if p.currentTokenKind() != lexer.GREATER {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.GREATER)

	return ast.GenericType{
		Identifier: identifier,
		Arguments:  arguments,
	}
}

func parse_fn_type(p *parser) ast.Type {
	var arguments = map[string]ast.FnArg{}

	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var explicitType ast.Type
		isMutable := p.currentTokenKind() == lexer.MUT
		if isMutable {
			p.advance()
		}

		argumentIdentifier := p.expect(lexer.IDENTIFIER).Literal
		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			explicitType = parse_type(p, default_bp)
		}

		_, exists := arguments[argumentIdentifier]

		if exists {
			p.err(fmt.Sprintf("Argument %s already exists in function type", argumentIdentifier))
		}

		arguments[argumentIdentifier] = ast.FnArg{
			Identifier: argumentIdentifier,
			IsMutable:  isMutable,
			Type:       explicitType,
		}

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	var returnType ast.Type

	if p.currentTokenKind() == lexer.R_ARROW {
		p.advance()
		returnType = parse_type(p, default_bp)
	}

	return ast.FnType{
		Arguments:  arguments,
		ReturnType: returnType,
	}
}

func parse_is_type(p *parser, left ast.Type, bp binding_power) ast.Type {
	p.advance()
	return ast.IsType{
		Left:  left,
		Right: parse_type(p, bp),
	}
}
