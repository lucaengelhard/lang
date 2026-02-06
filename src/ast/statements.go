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
	Identifier string
	Arguments  map[string]FnArg
	ReturnType Type
	Body       []Stmt
}

func (n FnStmt) stmt() {}
