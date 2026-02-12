package parser

import "github.com/lucaengelhard/lang/old_src/lexer"

type Parser struct {
	tokens     []lexer.Token
	tokenIndex int
}
