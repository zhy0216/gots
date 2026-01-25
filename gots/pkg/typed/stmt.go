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
	Pattern Pattern // Alternative to Name for destructuring
}

func (v *VarDecl) stmtNode() {}

// ----------------------------------------------------------------------------
// Patterns (for destructuring)
// ----------------------------------------------------------------------------

// Pattern is the interface for destructuring patterns.
type Pattern interface {
	patternNode()
	Type() types.Type
}

// ArrayPattern represents array destructuring: [a, b, c]
type ArrayPattern struct {
	Elements    []Pattern
	PatternType types.Type // The type of the array being destructured
}

func (a *ArrayPattern) patternNode()        {}
func (a *ArrayPattern) Type() types.Type   { return a.PatternType }

// ObjectPattern represents object destructuring: {x, y}
type ObjectPattern struct {
	Properties  []*PropertyPattern
	PatternType types.Type // The type of the object being destructured
}

func (o *ObjectPattern) patternNode()       {}
func (o *ObjectPattern) Type() types.Type  { return o.PatternType }

// PropertyPattern represents a property in object destructuring.
type PropertyPattern struct {
	Key   string  // Original property name
	Value Pattern // Target pattern (can be IdentPattern or nested)
}

// IdentPattern represents an identifier in a destructuring pattern.
type IdentPattern struct {
	Name        string
	PatternType types.Type // The type of the bound variable
}

func (i *IdentPattern) patternNode()       {}
func (i *IdentPattern) Type() types.Type  { return i.PatternType }

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
	TypeParams []*types.TypeParameter // Generic type parameters
	Params     []*Param
	ReturnType types.Type
	Body       *BlockStmt
	Captures   []*Capture // For closures
	IsAsync    bool       // true if declared with 'async' keyword
}

func (f *FuncDecl) stmtNode() {}

// ClassDecl represents a class declaration.
type ClassDecl struct {
	Name        string
	TypeParams  []*types.TypeParameter // Generic type parameters
	Super       string                 // Empty if no superclass
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

// InterfaceDecl represents a typed interface declaration.
type InterfaceDecl struct {
	Name    string
	Methods []*InterfaceMethodDecl
}

func (i *InterfaceDecl) stmtNode() {}

// InterfaceMethodDecl represents a method signature in a typed interface.
type InterfaceMethodDecl struct {
	Name       string
	Params     []*Param
	ReturnType types.Type
}

// GoImportDecl represents an import from a Go package.
type GoImportDecl struct {
	Names   []string // Imported names
	Package string   // Go package path
}

func (g *GoImportDecl) stmtNode() {}

// ModuleImportDecl represents an import from a local module.
type ModuleImportDecl struct {
	Names []string // Imported names
	Path  string   // Module path (e.g., "./math")
}

func (m *ModuleImportDecl) stmtNode() {}

// DefaultImport represents a default import.
type DefaultImport struct {
	Name string // Name to bind the default export
	Path string // Module path
}

func (d *DefaultImport) stmtNode() {}

// NamespaceImport represents a namespace import.
type NamespaceImport struct {
	Alias string // Namespace alias
	Path  string // Module path
}

func (n *NamespaceImport) stmtNode() {}

// ReExportDecl represents a re-export statement.
type ReExportDecl struct {
	Names      []string // Names being re-exported (empty for wildcard)
	Path       string   // Module path
	IsWildcard bool     // true for "export *"
}

func (r *ReExportDecl) stmtNode() {}

// DefaultExport represents a default export.
type DefaultExport struct {
	Decl Stmt // The declaration being exported
}

func (d *DefaultExport) stmtNode() {}

// ExportModifier wraps a declaration that is exported.
type ExportModifier struct {
	Decl Stmt // The exported declaration
}

func (e *ExportModifier) stmtNode() {}
