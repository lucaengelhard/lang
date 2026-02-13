package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	index  int
	errors []errorhandling.Error
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()
	createTokenTypeLookups()
	return &parser{
		tokens: tokens,
		index:  0,
		errors: make([]errorhandling.Error, 0),
	}
}

func Parse(tokens []lexer.Token) (ast.Stmt, []errorhandling.Error) {
	p := createParser(tokens)
	body := parse_block_stmt(p)
	return body, p.errors
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.index]
}

func (p *parser) printCurrentToken() {
	fmt.Println(p.currentTokenKind().ToString())
}

func (p *parser) peekNext() lexer.Token {
	return p.tokens[p.index+1]
}

func (p *parser) peekNextKind() lexer.TokenKind {
	return p.peekNext().Kind
}

func (p *parser) advance() lexer.Token {
	tk := p.currentToken()
	p.index++
	return tk
}

func (p *parser) hasTokens() bool {
	return p.index < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) addErr(msg string) {
	token := p.currentToken()
	p.errors = append(p.errors, errorhandling.Error{Message: "Parser error -> " + msg, Position: token.Position})
}

func (p *parser) expectError(expected lexer.TokenKind, msg any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expected {
		if msg == nil {
			msg = fmt.Sprintf("Expected token %s but recieved %s instead", expected.ToString(), kind.ToString())
		}
		p.addErr(fmt.Sprintf("%v", msg))
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
