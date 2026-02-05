package main

import (
	"os"

	"github.com/lucaengelhard/lang/src/lexer"
)

func main() {
	bytes, _ := os.ReadFile("./test/main.lang")
	tokens := lexer.Tokenize(string(bytes))

	for _, token := range tokens {
		token.Debug()
	}
}
