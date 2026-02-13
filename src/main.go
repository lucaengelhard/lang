package main

import (
	"flag"
	"os"

	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/interpreter"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	// Args
	source_path := flag.String("source", os.Args[1], "Source File")
	interpret := flag.Bool("interpret", true, "Set interpretation flag")
	debug := flag.Bool("debug", true, "Enable debugging")
	flag.Parse()

	// Reading file
	bytes, _ := os.ReadFile(*source_path)
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
	if *interpret {
		interpreter.Init(ast)
	}

	if *debug {
		litter.D(ast)
	}

	// Error handling
	errorhandling.PrintErrors(source, errors)
}
