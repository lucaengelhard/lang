package main

import (
	"fmt"
	"os"

	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
)

func main() {
	// Reading file
	bytes, _ := os.ReadFile(os.Args[1])
	source := string(bytes)

	// Tokenizing
	lexer := lexer.Tokenize(source)
	for _, t := range lexer.Tokens {
		fmt.Printf("%s -> %s\n", t.Kind.ToString(), t.Literal)
	}

	// AST-Building

	// Typechecking and updating of ast

	// Interpretation / Compilation

	// Error handling
	errorhandling.PrintErrors(source, lexer.Errors)
}
