package typed

import (
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Stmt is the interface for all typed statements.
type Stmt interface {
	stmtNode()
}

// ExprStmt wraps an expression as a statement.
type ExprStmt struct {
	Expr Expr
}

func (e *ExprStmt) stmtNode() {}

// VarDecl represents a variable declaration.
type VarDecl struct {
	Name    string
	VarType types.Type
	Init    Expr // nil if no initializer
	IsConst bool
}

func (v *VarDecl) stmtNode() {}

// BlockStmt represents a block of statements.
type BlockStmt struct {
	Stmts []Stmt
}

func (b *BlockStmt) stmtNode() {}

// IfStmt represents an if statement.
type IfStmt struct {
	Condition Expr
	Then      *BlockStmt
	Else      Stmt // *BlockStmt or *IfStmt (else if), or nil
}

func (i *IfStmt) stmtNode() {}

// WhileStmt represents a while loop.
type WhileStmt struct {
	Condition Expr
	Body      *BlockStmt
}

func (w *WhileStmt) stmtNode() {}

// ForStmt represents a for loop.
type ForStmt struct {
	Init      *VarDecl // nil if no init
	Condition Expr     // nil if no condition (infinite loop)
	Update    Expr     // nil if no update
	Body      *BlockStmt
}

func (f *ForStmt) stmtNode() {}

// ForOfStmt represents a for-of loop.
type ForOfStmt struct {
	Variable    *VarDecl // Loop variable
	Iterable    Expr     // Array or string to iterate
	ElementType types.Type
	Body        *BlockStmt
}

func (f *ForOfStmt) stmtNode() {}

// SwitchStmt represents a switch statement.
type SwitchStmt struct {
	Discriminant Expr
	Cases        []*CaseClause
}

func (s *SwitchStmt) stmtNode() {}

// CaseClause represents a case or default clause.
type CaseClause struct {
	Test  Expr   // nil for default case
	Stmts []Stmt // Consequent statements
}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Value Expr // nil for void return
}

func (r *ReturnStmt) stmtNode() {}

// BreakStmt represents a break statement.
type BreakStmt struct{}

func (b *BreakStmt) stmtNode() {}

// ContinueStmt represents a continue statement.
type ContinueStmt struct{}

func (c *ContinueStmt) stmtNode() {}

// TryStmt represents a try/catch statement.
type TryStmt struct {
	TryBlock   *BlockStmt
	CatchParam *VarDecl // The catch parameter (e.g., 'e' in catch(e))
	CatchBlock *BlockStmt
}

func (t *TryStmt) stmtNode() {}

// ThrowStmt represents a throw statement.
type ThrowStmt struct {
	Value Expr
}

func (t *ThrowStmt) stmtNode() {}

// FuncDecl represents a top-level function declaration.
type FuncDecl struct {
	Name       string
	Params     []*Param
	ReturnType types.Type
	Body       *BlockStmt
	Captures   []*Capture // For closures
}

func (f *FuncDecl) stmtNode() {}

// ClassDecl represents a class declaration.
type ClassDecl struct {
	Name        string
	Super       string // Empty if no superclass
	SuperClass  *types.Class
	Fields      []*FieldDecl
	Constructor *ConstructorDecl
	Methods     []*MethodDecl
}

func (c *ClassDecl) stmtNode() {}

// FieldDecl represents a class field.
type FieldDecl struct {
	Name string
	Type types.Type
}

// ConstructorDecl represents a class constructor.
type ConstructorDecl struct {
	Params []*Param
	Body   *BlockStmt
}

// MethodDecl represents a class method.
type MethodDecl struct {
	Name       string
	Params     []*Param
	ReturnType types.Type
	Body       *BlockStmt
}
