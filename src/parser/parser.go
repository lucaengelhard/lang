package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
)

func Parse(tokens []lexer.Token) (Statement, []errorhandling.Error) {
	p := createParser(tokens)
	return parse_block_stmt(p), p.Errors
}

type Parser struct {
	tokens     []lexer.Token
	tokenIndex int
	Errors     []errorhandling.Error
	forceExit  bool
}

func (p *Parser) err(message string) {
	p.Errors = append(p.Errors, errorhandling.Error{
		Message:  "Parser error -> " + message,
		Position: p.currentToken().Position,
	})
}

func (p *Parser) panic(message string) {
	p.err(message)
	p.forceExit = true
}

func createParser(tokens []lexer.Token) *Parser {
	createTokenLookups()
	return &Parser{
		tokens:     tokens,
		tokenIndex: 0,
		Errors:     make([]errorhandling.Error, 0),
		forceExit:  false,
	}
}

func (p *Parser) currentToken() lexer.Token {
	return p.tokens[p.tokenIndex]
}

func (p *Parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *Parser) printCurrentToken() {
	fmt.Println(p.currentTokenKind().ToString())
}

func (p *Parser) hasTokens() bool {
	return p.tokenIndex < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *Parser) advance() lexer.Token {
	tk := p.currentToken()
	if !p.hasTokens() {
		p.tokenIndex++
	}

	return tk
}

func (p *Parser) peekNext() lexer.Token {
	return p.tokens[p.tokenIndex+1]
}

func (p *Parser) advanceIfKind(kind lexer.TokenKind) (lexer.Token, bool) {
	if p.peekNext().Kind == kind {
		return p.expect(kind), true
	}

	return p.currentToken(), false
}

func (p *Parser) expectError(expected lexer.TokenKind, msg any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expected {
		if msg == nil {
			msg = fmt.Sprintf("Expected token %s but recieved %s instead", expected.ToString(), kind.ToString())
		}
		p.err(fmt.Sprintf("%v", msg))
	}

	return p.advance()
}

func (p *Parser) expect(expected lexer.TokenKind) lexer.Token {
	return p.expectError(expected, nil)
}
