package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type LexerError struct {
	Message  string
	Position int
}

type Lexer struct {
	patterns  []regexPattern
	Tokens    []Token
	source    string
	pos       int
	Errors    []LexerError
	forceExit bool
}

func (lex *Lexer) advanceN(n int) {
	lex.pos += n
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) remainder() string {
	return lex.source[lex.pos:]
}

func (lex *Lexer) at_eof() bool {
	return lex.pos >= len(lex.source)
}

func (lex *Lexer) err(message string) {
	lex.Errors = append(lex.Errors, LexerError{
		Message:  message,
		Position: lex.pos,
	})
}

func (lex *Lexer) panic(message string) {
	lex.err(message)
	lex.forceExit = true
}

func (lex *Lexer) Print() {
	fmt.Printf("[Lexer]: %v Tokens\n", len(lex.Tokens))
	for _, t := range lex.Tokens {
		fmt.Printf("%s -> %s\n", t.Kind.ToString(), t.Literal)
	}
}

func Tokenize(source string) *Lexer {
	lex := createLexer(source)

	for !lex.at_eof() && !lex.forceExit {
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
			lex.panic(fmt.Sprintf("unrecognized token at %s", lex.remainder()))
		}
	}

	lex.push(NewToken(EOF, "EOF", lex.pos))
	return lex
}

func createLexer(source string) *Lexer {
	InitTokenLookup()
	return &Lexer{pos: 0, source: source, Tokens: make([]Token, 0), Errors: make([]LexerError, 0), patterns: []regexPattern{
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

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *Lexer, regex *regexp.Regexp) {
		lex.advanceN(len(value))
		lex.push(NewToken(kind, value, lex.pos))
	}
}

func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())

	lex.push(NewToken(NUMBER, match, lex.pos))
	lex.advanceN(len(match))
}

func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	literal := lex.remainder()[match[0]:match[1]]

	/*
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
		} */

	lex.push(NewToken(STRING, literal, lex.pos))
	lex.advanceN(len(literal))
}

func symbolHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	if kind, exists := reserved_lookup[match]; exists {
		lex.push(NewToken(kind, match, lex.pos))
	} else {
		lex.push(NewToken(IDENTIFIER, match, lex.pos))
	}

	lex.advanceN(len(match))
}

func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainder())
	lex.advanceN(match[1])
}
