package parser

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/lexer"
)

type parserError struct {
	msg  string
	line int
	col  int
}

func (err *parserError) ToString() string {
	return fmt.Sprintf("[%v:%v] -> %s", err.line, err.col, err.msg)
}

func CreateParserError(msg string, line int, col int) *parserError {
	return &parserError{
		msg,
		line,
		col,
	}
}

type parser struct {
	tokens []lexer.Token
	pos    int
	errors []*parserError
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

	if len(p.errors) > 0 {
		fmt.Println("Errors occured during parsing:")
		for _, err := range p.errors {
			fmt.Println(err.ToString())
		}
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

func (p *parser) peekNext() lexer.Token {
	return p.tokens[p.pos+1]
}

func (p *parser) peekNextKind() lexer.TokenKind {
	return p.peekNext().Kind
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

func (p *parser) addErr(msg string) {
	token := p.currentToken()
	position := token.Position

	p.errors = append(p.errors, CreateParserError(msg, position.Line, position.Col))
}

func (p *parser) panic(msg string) {
	token := p.currentToken()
	position := token.Position

	panic(CreateParserError(msg, position.Line, position.Col).ToString())
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
