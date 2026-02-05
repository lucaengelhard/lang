package parser

import (
	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parse_declartion_stmt(p *parser) ast.Stmt {
	isMutable := p.nextIsKind(lexer.MUT)
	identifier := p.expectError(lexer.IDENTIFIER, "Expected identifier in variable declaration").Value
	p.expect(lexer.ASSIGNMENT)
	assignedValue := parse_expr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.DeclarationStmt{
		IsMutable:     isMutable,
		Identifier:    identifier,
		AssignedValue: assignedValue,
	}
}
