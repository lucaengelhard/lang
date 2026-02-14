package ast

type BlockStmt struct {
	Body []Stmt
	Position
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
	Position
}

func (n ExpressionStmt) stmt() {}

type DeclarationStmt struct {
	Identifier    string
	IsMutable     bool
	AssignedValue Expr
	Type          Type
	Position
}

func (n DeclarationStmt) stmt() {}

type StructPropertyModifier struct {
	Name string
	Position
}

type StructProperty struct {
	Name      string
	Type      Type
	Modifiers map[string]StructPropertyModifier
	Position
}

type StructStmt struct {
	Identifier string
	Type       Type
	Properties map[string]StructProperty
	Position
}

func (n StructStmt) stmt() {}

type InterfaceStmt struct {
	Identifier string
	TypeArg    Type
	SingleType Type
	StructType map[string]StructProperty
	Position
}

func (n InterfaceStmt) stmt() {}

type EnumStmt struct {
	Identifier string
	Elements   map[string]int
	Position
}

func (n EnumStmt) stmt() {}

type IfStmt struct {
	Condition Expr
	True      BlockStmt
	False     BlockStmt
	Position
}

func (n IfStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      BlockStmt
	Position
}

func (n WhileStmt) stmt() {}

type ForStmt struct {
	Assignment Stmt
	Condition  Stmt
	Increment  Expr
	Body       BlockStmt
	Position
}

func (n ForStmt) stmt() {}

type ReturnStmt struct {
	Value Expr
	Position
}

func (n ReturnStmt) stmt() {}

type ContinueStmt struct {
	Position
}

func (n ContinueStmt) stmt() {}

type BreakStmt struct {
	Position
}

func (n BreakStmt) stmt() {}

type ImportStmt struct {
	Identifier string
	Items      []string
	Path       string
	Position
}

func (n ImportStmt) stmt() {}
