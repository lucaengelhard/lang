package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
)

func main() {
	// Reading file
	bytes, _ := os.ReadFile(os.Args[1])
	source := string(bytes)

	errors := make([]errorhandling.Error, 0)

	// Tokenizing
	lexer := lexer.Tokenize(source)
	errors = append(errors, lexer.Errors...)

	// AST-Building

	// Typechecking and updating of ast

	// Interpretation / Compilation

	// Error handling
	errorhandling.PrintErrors(source, errors)
}
