package ast

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
}

type Position struct {
	Start int
	End   int
}

func CreatePosition(start int, end int) Position {
	return Position{Start: start, End: end}
}
