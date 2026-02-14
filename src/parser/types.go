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
	type_nud(lexer.STAR, parse_ref_type)
	/* 	type_nud(lexer.NUMBER, parse_number_type)
	   	type_nud(lexer.STRING, parse_string_type) */
	//type_nud(lexer.OPEN_PAREN, parse_fn_type)
}

func parse_type(p *parser, bp binding_power) ast.Type {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		p.err(fmt.Sprintf("Type Nud handler expected for token %s\n", tokenKind.ToString()))
		return ast.CreateUnsetType()
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]

		if !exists {
			p.err(fmt.Sprintf("Type Led handler expected for token %s\n", tokenKind.ToString()))
			return ast.CreateUnsetType()
		}

		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_symbol_type(p *parser) ast.Type {
	ident := p.expect(lexer.IDENTIFIER)
	args := make([]ast.Type, 0)

	if p.currentTokenKind() == lexer.LESS {
		p.advance()
		for p.hasTokens() && p.currentTokenKind() != lexer.GREATER {
			args = append(args, parse_type(p, logical))

			if p.currentTokenKind() != lexer.GREATER {
				p.expect(lexer.COMMA)
			}
		}
		p.expect(lexer.GREATER)
	}

	return ast.Type{Name: ident.Literal, Arguments: args}
}

func parse_ref_type(p *parser) ast.Type {
	p.expect(lexer.STAR)

	return ast.Type{Name: ast.REFERENCE, Arguments: []ast.Type{parse_type(p, logical)}}
}

/* func parse_string_type(p *parser) typechecker.Type {
	return typechecker.Type{Name: "string"}
}

func parse_number_type(p *parser) typechecker.Type {
	val := p.expect(lexer.NUMBER).Literal
	if _, err := strconv.ParseInt(val, 10, 64); err == nil {
		return typechecker.Type{Name: "int"}
	}

	return typechecker.Type{Name: "float"}
} */

/* func parse_fn_type(p *parser) typechecker.Type {
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

	var returnType typechecker.Type

	if p.currentTokenKind() == lexer.R_ARROW {
		p.advance()
		returnType = parse_type(p, default_bp)
	}

	return ast.FnType{
		Arguments:  arguments,
		ReturnType: returnType,
	}
} */
