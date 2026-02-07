package lexer

import (
	"fmt"
	"slices"
)

type TokenKind int

const (
	EOF TokenKind = iota

	NUMBER
	STRING
	IDENTIFIER

	OPEN_BRACKET
	CLOSE_BRACKET
	OPEN_CURLY
	CLOSE_CURLY
	OPEN_PAREN
	CLOSE_PAREN

	ASSIGNMENT
	EQUALS
	NOT
	NOT_EQUALS

	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	OR
	AND

	DOT
	SEMI_COLON
	COMMA
	COLON
	QUESTION
	SPREAD

	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS

	PLUS
	MINUS
	SLASH
	STAR
	PERCENT

	R_ARROW
	L_ARROW

	LET
	MUT
	IMPORT
	FN
	IF
	ELSE
	FOR
	WHILE
	EXPORT
	INTERFACE
	IN // TODO: For In Loop?
	TRUE
	FALSE
	STRUCT
	PUBLIC
	STATIC
	ENUM
	IS

	RETURN
	CONTINUE
	BREAK

	_TokenCount
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"let":       LET,
	"mut":       MUT,
	"import":    IMPORT,
	"fn":        FN,
	"if":        IF,
	"else":      ELSE,
	"for":       FOR,
	"while":     WHILE,
	"export":    EXPORT,
	"in":        IN,
	"true":      TRUE,
	"false":     FALSE,
	"struct":    STRUCT,
	"enum":      ENUM,
	"interface": INTERFACE,
	"public":    PUBLIC,
	"static":    STATIC,
	"is":        IS,
	"return":    RETURN,
	"continue":  CONTINUE,
	"break":     BREAK,
}

var token_string_lu map[TokenKind]string = map[TokenKind]string{
	EOF:            "eof",
	NUMBER:         "number",
	STRING:         "string",
	IDENTIFIER:     "identifier",
	OPEN_BRACKET:   "open_bracket",
	CLOSE_BRACKET:  "close_bracket",
	OPEN_CURLY:     "open_curly",
	CLOSE_CURLY:    "close_curly",
	OPEN_PAREN:     "open_paren",
	CLOSE_PAREN:    "close_paren",
	ASSIGNMENT:     "assignment",
	EQUALS:         "equals",
	NOT_EQUALS:     "not_equals",
	NOT:            "not",
	LESS:           "less",
	LESS_EQUALS:    "less_equals",
	GREATER:        "greater",
	GREATER_EQUALS: "greater_equals",
	OR:             "or",
	AND:            "and",
	DOT:            "dot",
	SEMI_COLON:     "semi_colon",
	COLON:          "colon",
	QUESTION:       "question",
	COMMA:          "comma",
	PLUS_PLUS:      "plus_plus",
	MINUS_MINUS:    "minus_minus",
	PLUS_EQUALS:    "plus_equals",
	MINUS_EQUALS:   "minus_equals",
	PLUS:           "plus",
	MINUS:          "dash",
	SLASH:          "slash",
	STAR:           "star",
	PERCENT:        "percent",
	R_ARROW:        "right_arrow",
	L_ARROW:        "left_arrow",
	SPREAD:         "spread",
}

func InitTokenLookup() {
	for value, kind := range reserved_lu {
		_, exists := token_string_lu[kind]
		if exists {
			panic(fmt.Sprintf("Definition conflict: %s already defined", kind.ToString()))
		}

		token_string_lu[kind] = value
	}

	for kind := range _TokenCount {
		_, exists := token_string_lu[kind]
		if !exists {
			panic(fmt.Sprintf("Token %s not in lookup", kind.ToString()))
		}
	}
}

func IsReserved(identifier string) bool {
	_, exists := reserved_lu[identifier]
	return exists
}

func (kind TokenKind) ToString() string {
	res, exist := token_string_lu[kind]

	if !exist {
		return fmt.Sprintf("unknown(%d)", kind)
	}

	return res
}

type TokenPosition struct {
	Line int
	Col  int
}

type Token struct {
	Kind     TokenKind
	Value    string
	Position TokenPosition
}

func (token Token) Is(kinds ...TokenKind) bool {
	return slices.Contains(kinds, token.Kind)
}

func (token Token) Debug() {
	if token.Is(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", token.Kind.ToString(), token.Value)
	} else {
		fmt.Printf("%s ()\n", token.Kind.ToString())
	}
}

func NewToken(kind TokenKind, value string, position TokenPosition) Token {
	return Token{Kind: kind, Value: value, Position: position}
}
