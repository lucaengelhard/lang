package errorhandling

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/lexer"
	"github.com/lucaengelhard/lang/src/lib"
)

func PrintErrors(source string, lexerErrors []lexer.LexerError) {
	if len(lexerErrors) > 0 {
		fmt.Printf("%v Errors occured during lexing:\n", len(lexerErrors))
		for _, err := range lexerErrors {
			row, col := lib.Int_to_file_pos(source, err.Position)
			print_error(err.Message, row, col)

		}
	}
}

func print_error(message string, row int, col int) {
	fmt.Printf("-> [%v:%v]: %s\n", row, col, message)
}
