package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/interpreter"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
)

func main() {
	// Args

	// Reading file
	bytes, _ := os.ReadFile(os.Args[1])
	source := string(bytes)

	// Create errors
	errors := make([]errorhandling.Error, 0)

	// Tokenizing
	tokens, lexer_errors := lexer.Tokenize(source)
	errors = append(errors, lexer_errors...)

	// AST-Building
	ast, parser_errors := parser.Parse(tokens)
	errors = append(errors, parser_errors...)

	// Typechecking and updating of ast

	// Interpretation / Compilation
	interpreter.Init(ast)

	// Error handling
	errorhandling.PrintErrors(source, errors)
}
