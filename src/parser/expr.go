package parser

import (
	"fmt"
	"strconv"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("Nud handler expected for token %s\n", tokenKind.ToString()))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("Led handler expected for token %s\n", tokenKind.ToString()))
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{Value: number}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{Value: p.advance().Value}
	default:
		panic(fmt.Sprintf("Cannot create primaty expression from %s\n", p.currentTokenKind().ToString()))
	}
}

func parse_binary_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operator := p.advance()
	right := parse_expr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func parse_assignment_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operator := p.advance()
	rightExpr := parse_expr(p, bp)

	return ast.AssignmentExpr{
		Operator:  operator,
		RightExpr: rightExpr,
		Assigne:   left,
	}
}

func parser_prefix_expr(p *parser) ast.Expr {
	operator := p.advance()
	rightExpr := parse_expr(p, default_bp)

	return ast.PrefixExpr{
		Operator:  operator,
		RightExpr: rightExpr,
	}
}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance()
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)

	return expr
}

func parse_struct_instantiation_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	structIdentifier := lib.ExpectType[ast.SymbolExpr](left).Value

	var properties = map[string]ast.Expr{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		propertyName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parse_expr(p, logical)

		properties[propertyName] = expr

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructInstantiationExpr{
		StructIdentifier: structIdentifier,
		Properties:       properties,
	}
}

// TODO: Change syntax to [number][10] or something similar
func parse_array_instantiation_expr(p *parser) ast.Expr {
	var elements = []ast.Expr{}

	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)

	arrayType := parse_type(p, default_bp)

	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		elements = append(elements, parse_expr(p, logical))
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}
	p.expect(lexer.CLOSE_CURLY)
	return ast.ArrayInstantiationExpr{
		Type:     arrayType,
		Elements: elements,
	}
}

func parse_fn_call_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	identifier := lib.ExpectType[ast.SymbolExpr](left).Value
	var arguments = []ast.FnCallArg{}

	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var argumentIdentifier string

		if p.peekNextKind() == lexer.COLON {
			argumentIdentifier = p.expect(lexer.IDENTIFIER).Value
			p.advance()
		}

		expr := parse_expr(p, logical)

		arguments = append(arguments, ast.FnCallArg{
			Identifier: argumentIdentifier,
			Value:      expr,
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	return ast.FnCallExpr{
		Identifier: identifier,
		Arguments:  arguments,
	}
}
