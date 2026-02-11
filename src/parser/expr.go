package parser

import (
	"fmt"
	"strconv"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	token := p.currentToken()
	tokenKind := token.Kind
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		p.panic(fmt.Sprintf("Unexpected token (nud) near: %s (%s)\n", tokenKind.ToString(), token.Value))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			p.panic(fmt.Sprintf("Unexpected token (led) near: %s (%s)\n", tokenKind.ToString(), token.Value))
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		val := p.advance().Value
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return ast.IntExpr{
				Value: i,
			}
		}

		number, _ := strconv.ParseFloat(val, 64)
		return ast.FloatExpr{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpr{Value: p.advance().Value}
	case lexer.TRUE:
		p.advance()
		return ast.BoolExpr{Value: true}
	case lexer.FALSE:
		p.advance()
		return ast.BoolExpr{Value: false}
	default:
		p.addErr(fmt.Sprintf("Cannot create primary expression from %s\n", p.currentTokenKind().ToString()))
		return ast.UnknowPrimary{}
	}
}

func parse_symbol_expr(p *parser) ast.Expr {
	isReference := p.currentTokenKind() == lexer.STAR
	if isReference {
		p.advance()
	}

	return ast.SymbolExpr{Value: p.advance().Value, IsReference: isReference}
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
		Operator: operator,
		Right:    rightExpr,
		Assignee: left,
	}
}

func parser_prefix_expr(p *parser) ast.Expr {
	operator := p.advance()
	rightExpr := parse_expr(p, default_bp)

	return ast.PrefixExpr{
		Operator: operator,
		Right:    rightExpr,
	}
}

func parse_postfix_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	operator := p.advance()

	return ast.AssignmentExpr{
		Assignee: left,
		Operator: operator,
		Right:    ast.IntExpr{Value: 1},
	}
}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance()
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)

	return expr
}

func parse_struct_instantiation_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	symbol, err := lib.ExpectType[ast.SymbolExpr](left)
	structIdentifier := symbol.Value

	if err != nil {
		p.addErr(err.Error())
	}

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

func parse_array_instantiation_expr(p *parser) ast.Expr {
	p.expect(lexer.OPEN_BRACKET)
	var elements = []ast.Expr{}
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_BRACKET {
		elements = append(elements, parse_expr(p, logical))
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_BRACKET)

	return ast.ArrayInstantiationExpr{
		Elements: elements,
	}
}

func parse_fn_call_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	var arguments = []ast.FnCallArg{}

	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var argumentIdentifier string

		if p.peekNextKind() == lexer.COLON {
			argumentIdentifier = p.expect(lexer.IDENTIFIER).Value
			p.advance()
		}

		expr := parse_expr(p, default_bp)

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
		Caller:    left,
		Arguments: arguments,
	}
}

func parse_fn_declare_anonymous_expr(p *parser) ast.Expr {
	p.expect(lexer.FN)
	return parse_fn_declare_expr(p)
}
func parse_fn_declare_expr(p *parser) ast.Expr {
	var arguments = map[string]ast.FnArg{}
	var typeArg ast.Type

	if p.currentTokenKind() == lexer.LESS {
		p.advance()
		typeArg = parse_type(p, default_bp)
		p.expect(lexer.GREATER)
	}

	p.expect(lexer.OPEN_PAREN)

	var arg_index = 0
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var explicitType ast.Type
		isMutable := p.currentTokenKind() == lexer.MUT
		if isMutable {
			p.advance()
		}

		isReference := p.currentTokenKind() == lexer.STAR
		if isReference {
			p.advance()
		}

		argumentIdentifier := p.expect(lexer.IDENTIFIER).Value
		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			explicitType = parse_type(p, default_bp)
		}

		_, exists := arguments[argumentIdentifier]

		if exists {
			p.addErr(fmt.Sprintf("Argument %s already exists in function", argumentIdentifier))
		}

		arguments[argumentIdentifier] = ast.FnArg{
			Identifier:  argumentIdentifier,
			Position:    arg_index,
			IsMutable:   isMutable,
			IsReference: isReference,
			Type:        explicitType,
		}

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA)
			arg_index++
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	var returnType ast.Type

	if p.currentTokenKind() == lexer.R_ARROW {
		p.advance()
		returnType = parse_type(p, default_bp)
	}

	p.expect(lexer.OPEN_CURLY)

	body := parse_block_stmt(p)

	p.expect(lexer.CLOSE_CURLY)

	return ast.FnDeclareExpr{
		Arguments:  arguments,
		Type:       typeArg,
		ReturnType: returnType,
		Body:       body,
	}
}

func parse_chain_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.expect(lexer.DOT)

	return ast.ChainExpr{
		Assignee: left,
		Member:   parse_expr(p, default_bp),
	}
}

func parse_is_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	p.expect(lexer.IS)
	right := parse_type(p, bp)

	return ast.IsTypeExpr{
		Left:  left,
		Right: right,
	}
}
