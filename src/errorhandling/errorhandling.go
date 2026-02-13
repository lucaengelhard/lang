package errorhandling

import (
	"fmt"

	"github.com/lucaengelhard/lang/src/lib"
)

type Error struct {
	Message      string
	Position     int
	TokenLiteral string
}

func PrintErrors(source string, errors []Error) {
	for _, err := range errors {
		row, col := lib.Int_to_file_pos(source, err.Position)
		print_error(err.Message, row, col)

	}
}

func print_error(message string, row int, col int) {
	fmt.Printf("[%v:%v]: %s\n", row, col, message)
}
