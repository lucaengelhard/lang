package lexer

import (
	"fmt"
	"regexp"
	"strconv"
)

type regexHandler func(lex *lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
}

func (lex *lexer) advanceN(n int) {
	lex.pos += n
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func Tokenize(source string, init bool) []Token {
	lex := createLexer(source, init)

	for !lex.at_eof() {
		matched := false

		for _, pattern := range lex.patterns {
			location := pattern.regex.FindStringIndex(lex.remainder())

			if location != nil && location[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			panic(fmt.Sprintf("Lexer Error -> unrecognized token near %s\n", lex.remainder()))
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.get_file_pos()))
	return lex.Tokens
}

func createLexer(source string, init bool) *lexer {
	if init {
		InitTokenLookup()
	}
	return &lexer{pos: 0, source: source, Tokens: make([]Token, 0), patterns: []regexPattern{
		{regexp.MustCompile(`\s+`), skipHandler},
		{regexp.MustCompile(`\/\/.*`), skipHandler},
		{regexp.MustCompile(`\/\*[\s\S]*?\*\/`), skipHandler},
		{regexp.MustCompile(`"(?:[^"\\]|\\.)*"`), stringHandler},
		{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
		{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
		{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
		{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
		{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
		{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
		{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
		{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
		{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
		{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
		{regexp.MustCompile(`=`), defaultHandler(ASSIGNMENT, "=")},
		{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
		{regexp.MustCompile(`<-`), defaultHandler(L_ARROW, "<-")},
		{regexp.MustCompile(`->`), defaultHandler(R_ARROW, "->")},
		{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
		{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
		{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
		{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
		{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
		{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
		{regexp.MustCompile(`\.\.\.`), defaultHandler(SPREAD, "...")},
		{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
		{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
		{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
		{regexp.MustCompile(`\?`), defaultHandler(QUESTION, "?")},
		{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
		{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
		{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
		{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
		{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
		{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
		{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
		{regexp.MustCompile(`/`), defaultHandler(SLASH, "/")},
		{regexp.MustCompile(`\*`), defaultHandler(STAR, "*")},
		{regexp.MustCompile(`%`), defaultHandler(PERCENT, "%")},
	}}
}

func (lex *lexer) get_file_pos() TokenPosition {
	var line = 1
	var col = 1
	for index, r := range lex.source {
		if index >= lex.pos {
			break
		}
		col++
		if r == '\n' || r == '\r' {
			line++
			col = 0
		}
	}

	return TokenPosition{
		Line: line,
		Col:  col,
	}
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, regex *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value, lex.get_file_pos()))
	}
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	lex.push(NewToken(NUMBER, match, lex.get_file_pos()))
	lex.advanceN(len(match))
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	literal := lex.remainder()[match[0]:match[1]]
	unqouted, _ := strconv.Unquote(literal)

	format := regexp.MustCompile(`{(.*?)}`).FindAllStringIndex(unqouted, -1)

	strs := make([]string, 0)
	sub_exprs := make([][]Token, 0)

	var prev = 0
	for _, format_match := range format {
		left := unqouted[prev:format_match[0]]
		strs = append(strs, left)
		prev = format_match[1]

		substr := unqouted[format_match[0]+1 : format_match[1]-1]
		sub_exprs = append(sub_exprs, Tokenize(substr, false))
	}

	strs = append(strs, unqouted[prev:])

	lex.push(NewToken(STRING, strs[0], lex.get_file_pos()))

	for i, expr := range sub_exprs {
		lex.push(NewToken(PLUS, "+", lex.get_file_pos()))
		lex.push(NewToken(OPEN_PAREN, "(", lex.get_file_pos()))
		for _, tok := range expr {
			if tok.Kind != EOF {
				lex.push(tok)
			}
		}

		lex.push(NewToken(CLOSE_PAREN, ")", lex.get_file_pos()))

		lex.push(NewToken(PLUS, "+", lex.get_file_pos()))

		lex.push(NewToken(STRING, strs[i+1], lex.get_file_pos()))
	}

	lex.advanceN(len(literal))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if kind, exists := reserved_lu[match]; exists {
		lex.push(NewToken(kind, match, lex.get_file_pos()))
	} else {
		lex.push(NewToken(IDENTIFIER, match, lex.get_file_pos()))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}
