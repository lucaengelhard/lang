package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/ast"
	"github.com/lucaengelhard/lang/src/errorhandling"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
	"github.com/lucaengelhard/lang/src/typechecker"
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
	var abstract_syntax_tree ast.Stmt
	var parser_errors []errorhandling.Error
	if len(errors) == 0 {
		abstract_syntax_tree, parser_errors = parser.Parse(tokens)
		errors = append(errors, parser_errors...)
	}

	// Typechecking and updating of ast
	if len(errors) == 0 {
		type_errors := typechecker.Init(abstract_syntax_tree)
		errors = append(errors, type_errors...)
	}

	// Interpretation / Compilation
	/* if len(errors) == 0 {
		interpreter.Init(abstract_syntax_tree)
	} */

	// Error handling
	errorhandling.PrintErrors(source, errors)
}
