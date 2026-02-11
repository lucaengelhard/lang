package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/interpreter"
	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/parser"
)

func main() {
	bytes, _ := os.ReadFile(os.Args[1])
	tokens := lexer.Tokenize(string(bytes), true)

	/* for _, t := range tokens {
		fmt.Printf("%s -> %s\n", t.Kind.ToString(), t.Value)
	} */

	ast := parser.Parse(tokens)
	//litter.Dump(ast)

	interpreter.Init(ast)

	/* res, err := typecheck.Check(ast, nil)

	for _, e := range err {
		fmt.Println(e.Error())
	}

	litter.Dump(res) */
}
