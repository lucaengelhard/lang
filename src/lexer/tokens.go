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
	FROM
	FN
	IF
	ELSE
	FOR
	WHILE
	EXPORT
	INTERFACE
	IN
	TRUE
	FALSE
	STRUCT
	PUBLIC
	STATIC
	ENUM
	RETURN
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"let":       LET,
	"mut":       MUT,
	"import":    IMPORT,
	"from":      FROM,
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
	"return":    RETURN,
}

func IsReserved(identifier string) bool {
	_, exists := reserved_lu[identifier]
	return exists
}

func (kind TokenKind) ToString() string {
	switch kind {
	case EOF:
		return "eof"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case IDENTIFIER:
		return "identifier"
	case OPEN_BRACKET:
		return "open_bracket"
	case CLOSE_BRACKET:
		return "close_bracket"
	case OPEN_CURLY:
		return "open_curly"
	case CLOSE_CURLY:
		return "close_curly"
	case OPEN_PAREN:
		return "open_paren"
	case CLOSE_PAREN:
		return "close_paren"
	case ASSIGNMENT:
		return "assignment"
	case EQUALS:
		return "equals"
	case NOT_EQUALS:
		return "not_equals"
	case NOT:
		return "not"
	case LESS:
		return "less"
	case LESS_EQUALS:
		return "less_equals"
	case GREATER:
		return "greater"
	case GREATER_EQUALS:
		return "greater_equals"
	case OR:
		return "or"
	case AND:
		return "and"
	case DOT:
		return "dot"
	case SEMI_COLON:
		return "semi_colon"
	case COLON:
		return "colon"
	case QUESTION:
		return "question"
	case COMMA:
		return "comma"
	case PLUS_PLUS:
		return "plus_plus"
	case MINUS_MINUS:
		return "minus_minus"
	case PLUS_EQUALS:
		return "plus_equals"
	case MINUS_EQUALS:
		return "minus_equals"
	case PLUS:
		return "plus"
	case MINUS:
		return "dash"
	case SLASH:
		return "slash"
	case STAR:
		return "star"
	case PERCENT:
		return "percent"
	case LET:
		return "let"
	case MUT:
		return "mut"
	case IMPORT:
		return "import"
	case FROM:
		return "from"
	case FN:
		return "fn"
	case IF:
		return "if"
	case ELSE:
		return "else"
	case FOR:
		return "for"
	case WHILE:
		return "while"
	case EXPORT:
		return "export"
	case IN:
		return "in"
	case STRUCT:
		return "struct"
	case R_ARROW:
		return "right_arrow"
	case L_ARROW:
		return "left_arrow"
	case ENUM:
		return "enum"
	case INTERFACE:
		return "interface"
	case PUBLIC:
		return "public"
	case STATIC:
		return "static"
	case RETURN:
		return "return"
	default:
		return fmt.Sprintf("unknown(%d)", kind)
	}
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

func (token Token) is(kinds ...TokenKind) bool {
	return slices.Contains(kinds, token.Kind)
}

func (token Token) Debug() {
	if token.is(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s (%s)\n", token.Kind.ToString(), token.Value)
	} else {
		fmt.Printf("%s ()\n", token.Kind.ToString())
	}
}

func NewToken(kind TokenKind, value string, position TokenPosition) Token {
	return Token{Kind: kind, Value: value, Position: position}
}
