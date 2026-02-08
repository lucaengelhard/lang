package main

import (
	"fmt"
	"os"

	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
	"github.com/lucaengelhard/lang/src/typecheck"
	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile(os.Args[1])
	tokens := lexer.Tokenize(string(bytes))

	ast := parser.Parse(tokens)

	res, err := typecheck.Check(ast, nil)

	for _, e := range err {
		fmt.Println(e.Error())
	}

	litter.Dump(res)
}
