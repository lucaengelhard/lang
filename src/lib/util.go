package lib

func IsType[T any](input any) bool {
	_, ok := input.(T)

	return ok
}

func Int_to_file_pos(source string, pos int) (row int, col int) {
	row = 1
	col = 1
	for index, r := range source {
		if index >= pos {
			break
		}
		col++
		if r == '\n' || r == '\r' {
			row++
			col = 0
		}
	}
	return row, col
}
