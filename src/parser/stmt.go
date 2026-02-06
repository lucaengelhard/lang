package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

func parse_stmt(p *parser, with_semicolon ...bool) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)

	if (len(with_semicolon) > 0 && with_semicolon[0]) || len(with_semicolon) == 0 {
		p.expect(lexer.SEMI_COLON)
	}

	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parse_declartion_stmt(p *parser) ast.Stmt {
	var explicitType ast.Type

	isMutable := p.nextIsKind(lexer.MUT)
	identifier := p.expect(lexer.IDENTIFIER).Value

	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		explicitType = parse_type(p, default_bp)
	}

	p.expect(lexer.ASSIGNMENT)
	assignedValue := parse_expr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.DeclarationStmt{
		IsMutable:     isMutable,
		Identifier:    identifier,
		AssignedValue: assignedValue,
		Type:          explicitType,
	}
}

func parse_struct_stmt(p *parser) ast.Stmt {
	p.expect(lexer.STRUCT)
	identifier := p.expect(lexer.IDENTIFIER).Value
	var properties = map[string]ast.StructProperty{}

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		if p.currentTokenKind() == lexer.IDENTIFIER || lexer.IsReserved(p.currentToken().Value) {
			var propertyModifiers = map[string]ast.StructPropertyModifier{}
			for lexer.IsReserved(p.currentToken().Value) {
				propertyModifiers[p.currentToken().Value] = ast.StructPropertyModifier{
					Name: p.currentToken().Value,
				}
				p.advance()
			}

			propertyName := p.expect(lexer.IDENTIFIER).Value
			p.expect(lexer.COLON)
			propertyType := parse_type(p, default_bp)
			p.expect(lexer.SEMI_COLON)

			_, exists := properties[propertyName]

			if exists {
				p.addErr(fmt.Sprintf("Property %s already exists on struct %s", propertyName, identifier))
			}

			properties[propertyName] = ast.StructProperty{
				Name:      propertyName,
				Type:      propertyType,
				Modifiers: propertyModifiers,
			}

			continue
		}

		p.addErr("This souldn't be reached :( so i wrote bad struct code")
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.StructStmt{
		Identifier: identifier,
		Properties: properties,
	}
}

func parse_fn_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FN)
	identifier := p.expect(lexer.IDENTIFIER).Value
	var arguments = map[string]ast.FnArg{}
	var genericType ast.Type

	if p.currentTokenKind() == lexer.LESS {
		p.advance()
		genericType = parse_type(p, default_bp)
		p.expect(lexer.GREATER)
	}

	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		var explicitType ast.Type
		isMutable := p.currentTokenKind() == lexer.MUT
		if isMutable {
			p.advance()
		}

		argumentIdentifier := p.expect(lexer.IDENTIFIER).Value
		if p.currentTokenKind() == lexer.COLON {
			p.advance()
			explicitType = parse_type(p, default_bp)
		}

		_, exists := arguments[argumentIdentifier]

		if exists {
			p.addErr(fmt.Sprintf("Argument %s already exists in function %s", argumentIdentifier, identifier))
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

	body := make([]ast.Stmt, 0)

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)

	return ast.FnStmt{
		Identifier:  identifier,
		Arguments:   arguments,
		GenericType: genericType,
		ReturnType:  returnType,
		Body:        body,
	}
}

func parse_if_stmt(p *parser) ast.Stmt {
	p.expect(lexer.IF)
	p.expect(lexer.OPEN_PAREN)
	cond := parse_expr(p, assignment)
	p.expect(lexer.CLOSE_PAREN)

	true_stmt := make([]ast.Stmt, 0)
	false_stmt := make([]ast.Stmt, 0)

	p.expect(lexer.OPEN_CURLY)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		true_stmt = append(true_stmt, parse_stmt(p))
	}
	p.expect(lexer.CLOSE_CURLY)

	if p.currentTokenKind() == lexer.ELSE {
		p.advance()
		p.expect(lexer.OPEN_CURLY)

		for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
			false_stmt = append(true_stmt, parse_stmt(p))
		}

		p.expect(lexer.CLOSE_CURLY)
	}

	return ast.IfStmt{
		Condition: cond,
		True:      true_stmt,
		False:     false_stmt,
	}
}

func parse_while_stmt(p *parser) ast.Stmt {
	p.expect(lexer.WHILE)
	p.expect(lexer.OPEN_PAREN)
	cond := parse_expr(p, assignment)
	p.expect(lexer.CLOSE_PAREN)

	p.expect(lexer.OPEN_CURLY)
	body := make([]ast.Stmt, 0)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}
	p.expect(lexer.CLOSE_CURLY)

	return ast.WhileStmt{
		Condition: cond,
		Body:      body,
	}
}

func parse_for_stmt(p *parser) ast.Stmt {
	p.expect(lexer.FOR)
	p.expect(lexer.OPEN_PAREN)
	assignemt := parse_stmt(p)
	cond := parse_stmt(p)
	incr := parse_stmt(p, false)
	p.expect(lexer.CLOSE_PAREN)

	p.expect(lexer.OPEN_CURLY)
	body := make([]ast.Stmt, 0)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}
	p.expect(lexer.CLOSE_CURLY)

	return ast.ForStmt{
		Assignment: assignemt,
		Condition:  cond,
		Increment:  incr,
		Body:       body,
	}
}

func parse_return_stmt(p *parser) ast.Stmt {
	p.expect(lexer.RETURN)
	expr := parse_expr(p, logical)
	p.expect(lexer.SEMI_COLON)

	return ast.ReturnStmt{
		Value: expr,
	}
}

func parse_continue_stmt(p *parser) ast.Stmt {
	p.expect(lexer.CONTINUE)
	p.expect(lexer.SEMI_COLON)
	return ast.ContinueStmt{}
}

func parse_break_stmt(p *parser) ast.Stmt {
	p.expect(lexer.BREAK)
	p.expect(lexer.SEMI_COLON)
	return ast.BreakStmt{}
}

func parse_import_stmt(p *parser) ast.Stmt {
	p.expect(lexer.IMPORT)
	path := p.expect(lexer.STRING).Value
	p.expect(lexer.R_ARROW)

	var identifier string
	items := make([]string, 0)

	if p.currentTokenKind() == lexer.OPEN_CURLY {
		p.advance()
		for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
			items = append(items, p.expect(lexer.IDENTIFIER).Value)

			if p.currentTokenKind() != lexer.CLOSE_CURLY {
				p.expect(lexer.COMMA)
			}
		}
		p.expect(lexer.CLOSE_CURLY)
	} else {
		identifier = p.expect(lexer.IDENTIFIER).Value
	}

	return ast.ImportStmt{
		Path:       path,
		Identifier: identifier,
		Items:      items,
	}
}
