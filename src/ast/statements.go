package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (n ExpressionStmt) stmt() {}

type DeclarationStmt struct {
	Identifier    string
	IsMutable     bool
	AssignedValue Expr
	Type          Type
}

func (n DeclarationStmt) stmt() {}

type StructPropertyModifier struct {
	Name string
}

type StructProperty struct {
	Name      string
	Type      Type
	Modifiers map[string]StructPropertyModifier
}

type StructStmt struct {
	Identifier string
	Properties map[string]StructProperty
}

func (n StructStmt) stmt() {}

type FnArg struct {
	Identifier string
	IsMutable  bool
	Type       Type
}

type FnStmt struct {
	Identifier  string
	Arguments   map[string]FnArg
	GenericType Type
	ReturnType  Type
	Body        []Stmt
}

func (n FnStmt) stmt() {}

type IfStmt struct {
	Condition Expr
	True      []Stmt
	False     []Stmt
}

func (n IfStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      []Stmt
}

func (n WhileStmt) stmt() {}

type ForStmt struct {
	Assignment Stmt
	Condition  Stmt
	Increment  Stmt
	Body       []Stmt
}

func (n ForStmt) stmt() {}

type ReturnStmt struct {
	Value Expr
}

func (n ReturnStmt) stmt() {}

type ContinueStmt struct{}

func (n ContinueStmt) stmt() {}

type BreakStmt struct{}

func (n BreakStmt) stmt() {}

type ImportStmt struct {
	Identifier string
	Items      []string
	Path       string
}

func (n ImportStmt) stmt() {}
