package parser

import (
	"github.com/lucaengelhard/lang/src/lexer"
)

type Statement interface {
	stmt()
}

type ErrorStmt struct{}

func (n ErrorStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expression
}

func (n ExpressionStmt) stmt() {}

func parse_stmt(p *Parser) Statement {
	stmt_fn, exists := stmt_lookup[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)

	p.advanceIfKind(lexer.SEMI_COLON)

	return ExpressionStmt{
		Expression: expression,
	}
}

type BlockStmt struct {
	Body []Statement
}

func (n BlockStmt) stmt() {}

func parse_block_stmt(p *Parser) Statement {
	body := make([]Statement, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parse_stmt(p))
	}

	return BlockStmt{
		Body: body,
	}
}

type DeclarationStmt struct {
	Identifier     string
	IsMutable      bool
	AssignedValue  Expression
	Type           Type
	SourcePosition int
}

func (n DeclarationStmt) stmt() {}

func parse_declaration_stmt(p *Parser) Statement {
	keyword := p.expect(lexer.LET)
	_, isMutable := p.advanceIfKind(lexer.MUT)
	identifier := p.expect(lexer.IDENTIFIER)

	var explicitType Type
	if _, exexplicitTypeSet := p.advanceIfKind(lexer.COLON); exexplicitTypeSet {
		explicitType = parse_type(p, default_bp)
	}

	p.expect(lexer.ASSIGNMENT)
	assignedValue := parse_expr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return DeclarationStmt{
		Identifier:     identifier.Literal,
		IsMutable:      isMutable,
		SourcePosition: keyword.Position,
		Type:           explicitType,
		AssignedValue:  assignedValue,
	}
}
