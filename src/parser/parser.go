package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()
	createTokenTypeLookups()
	return &parser{
		tokens: tokens,
	}
}

func Parse(tokens []lexer.Token) ast.BlockStmt {
	body := make([]ast.Stmt, 0)

	p := createParser(tokens)

	for p.hasTokens() {
		body = append(body, parse_stmt(p))
	}

	return ast.BlockStmt{
		Body: body,
	}
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) printCurrentToken() {
	fmt.Println(p.currentTokenKind().ToString())
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.pos++
	return tk
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) expectError(expected lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expected {
		if err == nil {
			err = fmt.Sprintf("Expected %s but recieved %s instead\n", expected.ToString(), kind.ToString())
		}

		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(expected lexer.TokenKind) lexer.Token {
	return p.expectError(expected, nil)
}

func (p *parser) nextIsKind(expected lexer.TokenKind) bool {
	p.advance()
	if p.currentTokenKind() != expected {
		return false
	}

	p.advance()
	return true
}
