package lexer

import (
	"fmt"
	"slices"
)

type Token struct {
	Kind     TokenKind
	Position int
	Literal  string
}

func (token Token) Is(kinds ...TokenKind) bool {
	return slices.Contains(kinds, token.Kind)
}

func NewToken(kind TokenKind, literal string, position int) Token {
	return Token{Kind: kind, Literal: literal, Position: position}
}

type TokenKind int

func (kind TokenKind) ToString() string {
	res, exist := token_string_lookup[kind]

	if !exist {
		return fmt.Sprintf("unknown(%d)", kind)
	}

	return res
}

const (
	EOF TokenKind = iota
	IDENTIFIER

	// Datatypes
	NUMBER
	STRING

	// Operators
	ASSIGNMENT
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS
	// Arithmetic Operators
	PLUS
	MINUS
	SLASH
	STAR
	PERCENT
	// Boolean Operators
	OR
	AND
	EQUALS
	NOT
	NOT_EQUALS
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	// Symbols
	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	R_ARROW
	L_ARROW

	DOT
	SEMI_COLON
	COMMA
	COLON
	QUESTION
	SPREAD

	// Keywords
	// Variables
	LET
	MUT
	// Modules
	IMPORT
	EXPORT
	// Statements
	FN
	IF
	ELSE
	FOR
	WHILE
	// Values
	TRUE
	FALSE
	// Types
	INTERFACE
	STRUCT
	ENUM
	IS
	// Control flow
	RETURN
	CONTINUE
	BREAK

	_TokenCount
)

var reserved_lookup = map[string]TokenKind{
	"let":       LET,
	"mut":       MUT,
	"import":    IMPORT,
	"export":    EXPORT,
	"fn":        FN,
	"if":        IF,
	"else":      ELSE,
	"for":       FOR,
	"while":     WHILE,
	"true":      TRUE,
	"false":     FALSE,
	"interface": INTERFACE,
	"struct":    STRUCT,
	"enum":      ENUM,
	"is":        IS,
	"return":    RETURN,
	"continue":  CONTINUE,
	"break":     BREAK,
}

var token_string_lookup = map[TokenKind]string{
	EOF:            "eof",
	IDENTIFIER:     "identifier",
	NUMBER:         "number",
	STRING:         "string",
	ASSIGNMENT:     "assignment",
	PLUS_PLUS:      "plus_plus",
	MINUS_MINUS:    "minus_minus",
	PLUS_EQUALS:    "plus_equals",
	MINUS_EQUALS:   "minus_equals",
	PLUS:           "plus",
	MINUS:          "dash",
	SLASH:          "slash",
	STAR:           "star",
	PERCENT:        "percent",
	OR:             "or",
	AND:            "and",
	EQUALS:         "equals",
	NOT_EQUALS:     "not_equals",
	NOT:            "not",
	LESS:           "less",
	LESS_EQUALS:    "less_equals",
	GREATER:        "greater",
	GREATER_EQUALS: "greater_equals",
	OPEN_BRACKET:   "open_bracket",
	CLOSE_BRACKET:  "close_bracket",
	OPEN_CURLY:     "open_curly",
	CLOSE_CURLY:    "close_curly",
	OPEN_PAREN:     "open_paren",
	CLOSE_PAREN:    "close_paren",
	R_ARROW:        "right_arrow",
	L_ARROW:        "left_arrow",
	DOT:            "dot",
	SEMI_COLON:     "semi_colon",
	COLON:          "colon",
	QUESTION:       "question",
	COMMA:          "comma",
	SPREAD:         "spread",
}

func InitTokenLookup() {
	for value, kind := range reserved_lookup {
		_, exists := token_string_lookup[kind]
		if exists {
			panic(fmt.Sprintf("Definition conflict: %s already defined", kind.ToString()))
		}

		token_string_lookup[kind] = value
	}

	for kind := range _TokenCount {
		_, exists := token_string_lookup[kind]
		if !exists {
			panic(fmt.Sprintf("Token %s not in lookup", kind.ToString()))
		}
	}
}
