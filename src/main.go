package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./test/main.lang")
	tokens := lexer.Tokenize(string(bytes))

	ast := parser.Parse(tokens)

	litter.Dump(ast)
}
