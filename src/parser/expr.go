package parser

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

func parse_expr(p *parser, bp binding_power) ast.Expr {
	token := p.currentToken()
	tokenKind := token.Kind
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		p.err(fmt.Sprintf("Unexpected token (nud) near: %s (%s)\n", tokenKind.ToString(), token.Literal))
		return nil
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]

		if !exists {
			p.err(fmt.Sprintf("Unexpected token (led) near: %s (%s)\n", tokenKind.ToString(), token.Literal))
			return nil
		}

		left = led_fn(p, left, bp_lu[p.currentTokenKind()])
	}

	return left
}

func parse_boolean_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
	switch p.currentTokenKind() {
	case lexer.TRUE:
		p.advance()
		return ast.BoolExpr{Value: true, Position: pos}
	case lexer.FALSE:
		p.advance()
		return ast.BoolExpr{Value: false, Position: pos}
	default:
		p.err(fmt.Sprintf("Cannot create boolean expression from %s\n", p.currentTokenKind().ToString()))
		return ast.UnknowPrimary{}
	}
}

func parse_number_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
	val := p.advance().Literal
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return ast.IntExpr{
			Value:    i,
			Position: pos,
		}
	}

	number, _ := strconv.ParseFloat(val, 64)
	return ast.FloatExpr{
		Value:    number,
		Position: pos,
	}
}

func parse_string_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
	literal := p.advance().Literal
	return ast.StringExpr{Value: literal, Position: pos}
}

func parse_symbol_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()

	isReference := p.currentTokenKind() == lexer.STAR
	if isReference {
		p.advance()
	}

	return ast.SymbolExpr{Value: p.advance().Literal, IsReference: isReference, Position: pos}
}

func parse_binary_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()

	operator := p.advance()
	right := parse_expr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
		Position: pos,
	}
}

func parse_assignment_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()
	operator := p.advance()
	rightExpr := parse_expr(p, bp)

	return ast.AssignmentExpr{
		Operator: operator,
		Right:    rightExpr,
		Assignee: left,
		Position: pos,
	}
}

func parser_prefix_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
	operator := p.advance()
	rightExpr := parse_expr(p, default_bp)

	return ast.PrefixExpr{
		Operator: operator,
		Right:    rightExpr,
		Position: pos,
	}
}

func parse_postfix_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()
	operator := p.advance()

	return ast.AssignmentExpr{
		Assignee: left,
		Operator: operator,
		Right:    ast.IntExpr{Value: 1},
		Position: pos,
	}
}

func parse_grouping_expr(p *parser) ast.Expr {
	p.advance()
	expr := parse_expr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)

	return expr
}

func parse_struct_instantiation_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()

	symbol, ok := left.(ast.SymbolExpr)

	if !ok {
		p.err(fmt.Sprintf("Type error: Expected %s got %s", reflect.TypeFor[ast.SymbolExpr](), reflect.TypeOf(symbol)))
	}

	structIdentifier := symbol.Value

	var properties = map[string]ast.Expr{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		propertyName := p.expect(lexer.IDENTIFIER).Literal
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
		Position:         pos,
	}
}

func parse_array_instantiation_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
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
		Position: pos,
	}
}

func parse_fn_call_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()
	var arguments = []ast.FnCallArg{}

	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var argumentIdentifier string

		if p.peekNextKind() == lexer.COLON {
			argumentIdentifier = p.expect(lexer.IDENTIFIER).Literal
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
		Position:  pos,
	}
}

func parse_fn_declare_anonymous_expr(p *parser) ast.Expr {
	p.expect(lexer.FN)
	return parse_fn_declare_expr(p)
}
func parse_fn_declare_expr(p *parser) ast.Expr {
	pos := p.curentTokenPosition()
	var arguments = map[string]ast.FnArg{}
	var typeArg = ast.CreateUnsetType()

	if p.currentTokenKind() == lexer.LESS {
		p.advance()
		typeArg = parse_type(p, default_bp)
		p.expect(lexer.GREATER)
	}

	p.expect(lexer.OPEN_PAREN)

	var arg_index = 0
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var explicitType = ast.CreateUnsetType()
		isMutable := p.currentTokenKind() == lexer.MUT
		if isMutable {
			p.advance()
		}

		isReference := p.currentTokenKind() == lexer.STAR
		if isReference {
			p.advance()
		}

		argumentIdentifier := p.expect(lexer.IDENTIFIER).Literal
		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			explicitType = parse_type(p, default_bp)
		}

		_, exists := arguments[argumentIdentifier]

		if exists {
			p.err(fmt.Sprintf("Argument %s already exists in function", argumentIdentifier))
		}

		arguments[argumentIdentifier] = ast.FnArg{
			Identifier:  argumentIdentifier,
			ArgIndex:    arg_index,
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

	var returnType = ast.CreateUnsetType()

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
		Position:   pos,
	}
}

func parse_chain_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()

	p.expect(lexer.DOT)

	return ast.ChainExpr{
		Assignee: left,
		Member:   parse_expr(p, default_bp),
		Position: pos,
	}
}

func parse_is_expr(p *parser, left ast.Expr, bp binding_power) ast.Expr {
	pos := p.curentTokenPosition()
	p.expect(lexer.IS)
	right := parse_type(p, bp)

	return ast.IsTypeExpr{
		Left:     left,
		Right:    right,
		Position: pos,
	}
}
