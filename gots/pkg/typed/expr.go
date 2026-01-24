package typed

import (
	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Expr is the interface for all typed expressions.
type Expr interface {
	exprNode()
	Type() types.Type
}

// NumberLit represents a numeric literal.
type NumberLit struct {
	Value    float64
	ExprType types.Type
}

func (n *NumberLit) exprNode()        {}
func (n *NumberLit) Type() types.Type { return n.ExprType }

// StringLit represents a string literal.
type StringLit struct {
	Value    string
	ExprType types.Type
}

func (s *StringLit) exprNode()        {}
func (s *StringLit) Type() types.Type { return s.ExprType }

// TemplateLit represents a template literal with interpolations.
type TemplateLit struct {
	Parts       []string     // Static string parts
	Expressions []Expr       // Interpolated expressions
	ExprType    types.Type   // Always string
}

func (t *TemplateLit) exprNode()        {}
func (t *TemplateLit) Type() types.Type { return t.ExprType }

// BoolLit represents a boolean literal.
type BoolLit struct {
	Value    bool
	ExprType types.Type
}

func (b *BoolLit) exprNode()        {}
func (b *BoolLit) Type() types.Type { return b.ExprType }

// NullLit represents a null literal.
type NullLit struct {
	ExprType types.Type
}

func (n *NullLit) exprNode()        {}
func (n *NullLit) Type() types.Type { return n.ExprType }

// Ident represents an identifier reference.
type Ident struct {
	Name     string
	ExprType types.Type
}

func (i *Ident) exprNode()        {}
func (i *Ident) Type() types.Type { return i.ExprType }

// BinaryExpr represents a binary expression.
type BinaryExpr struct {
	Left     Expr
	Op       string // "+", "-", "*", "/", "%", "<", ">", "<=", ">=", "==", "!=", "&&", "||", "??"
	Right    Expr
	ExprType types.Type
}

func (b *BinaryExpr) exprNode()        {}
func (b *BinaryExpr) Type() types.Type { return b.ExprType }

// UnaryExpr represents a unary expression.
type UnaryExpr struct {
	Op       string // "-", "!"
	Operand  Expr
	ExprType types.Type
}

func (u *UnaryExpr) exprNode()        {}
func (u *UnaryExpr) Type() types.Type { return u.ExprType }

// SpreadExpr represents a spread expression (...arg).
type SpreadExpr struct {
	Argument Expr
	ExprType types.Type
}

func (s *SpreadExpr) exprNode()        {}
func (s *SpreadExpr) Type() types.Type { return s.ExprType }

// CallExpr represents a function call.
type CallExpr struct {
	Callee   Expr
	Args     []Expr
	Optional bool // true for ?.() optional chaining
	ExprType types.Type
}

func (c *CallExpr) exprNode()        {}
func (c *CallExpr) Type() types.Type { return c.ExprType }

// IndexExpr represents array/string indexing.
type IndexExpr struct {
	Object   Expr
	Index    Expr
	Optional bool // true for ?.[] optional chaining
	ExprType types.Type
}

func (i *IndexExpr) exprNode()        {}
func (i *IndexExpr) Type() types.Type { return i.ExprType }

// PropertyExpr represents property access.
type PropertyExpr struct {
	Object   Expr
	Property string
	Optional bool // true for ?. optional chaining
	ExprType types.Type
}

func (p *PropertyExpr) exprNode()        {}
func (p *PropertyExpr) Type() types.Type { return p.ExprType }

// ArrayLit represents an array literal.
type ArrayLit struct {
	Elements []Expr
	ExprType types.Type // *types.Array
}

func (a *ArrayLit) exprNode()        {}
func (a *ArrayLit) Type() types.Type { return a.ExprType }

// ObjectLit represents an object literal.
type ObjectLit struct {
	Properties []*PropertyInit
	ExprType   types.Type // *types.Object
}

type PropertyInit struct {
	Key   string
	Value Expr
}

func (o *ObjectLit) exprNode()        {}
func (o *ObjectLit) Type() types.Type { return o.ExprType }

// FuncExpr represents a function expression (including arrow functions).
type FuncExpr struct {
	Params   []*Param
	Body     *BlockStmt  // nil for expression-body arrow functions
	BodyExpr Expr        // non-nil for expression-body arrow functions
	Captures []*Capture  // Variables captured from enclosing scopes
	ExprType types.Type  // *types.Function
}

type Param struct {
	Name string
	Type types.Type
}

func (f *FuncExpr) exprNode()        {}
func (f *FuncExpr) Type() types.Type { return f.ExprType }

// NewExpr represents object instantiation.
type NewExpr struct {
	ClassName string
	TypeArgs  []types.Type // Type arguments for generic class instantiation
	Args      []Expr
	ExprType  types.Type // *types.Class
}

func (n *NewExpr) exprNode()        {}
func (n *NewExpr) Type() types.Type { return n.ExprType }

// ThisExpr represents the 'this' keyword.
type ThisExpr struct {
	ExprType types.Type // *types.Class
}

func (t *ThisExpr) exprNode()        {}
func (t *ThisExpr) Type() types.Type { return t.ExprType }

// SuperExpr represents a super() call.
type SuperExpr struct {
	Args     []Expr
	ExprType types.Type
}

func (s *SuperExpr) exprNode()        {}
func (s *SuperExpr) Type() types.Type { return s.ExprType }

// AssignExpr represents an assignment expression.
type AssignExpr struct {
	Target   Expr // Ident, IndexExpr, or PropertyExpr
	Value    Expr
	ExprType types.Type
}

func (a *AssignExpr) exprNode()        {}
func (a *AssignExpr) Type() types.Type { return a.ExprType }

// CompoundAssignExpr represents compound assignment (+=, -=, etc.).
type CompoundAssignExpr struct {
	Target   Expr   // Ident, IndexExpr, or PropertyExpr
	Op       string // "+=", "-=", "*=", "/=", "%="
	Value    Expr
	ExprType types.Type
}

func (c *CompoundAssignExpr) exprNode()        {}
func (c *CompoundAssignExpr) Type() types.Type { return c.ExprType }

// UpdateExpr represents increment/decrement (++, --).
type UpdateExpr struct {
	Op       string // "++" or "--"
	Operand  Expr   // Must be assignable
	Prefix   bool   // true for ++x, false for x++
	ExprType types.Type
}

func (u *UpdateExpr) exprNode()        {}
func (u *UpdateExpr) Type() types.Type { return u.ExprType }

// BuiltinCall represents a call to a built-in function.
type BuiltinCall struct {
	Name     string
	Args     []Expr
	ExprType types.Type
}

func (b *BuiltinCall) exprNode()        {}
func (b *BuiltinCall) Type() types.Type { return b.ExprType }

// MapLit represents a map literal.
type MapLit struct {
	Entries  []*MapEntry
	ExprType types.Type // *types.Map
}

type MapEntry struct {
	Key   Expr
	Value Expr
}

func (m *MapLit) exprNode()        {}
func (m *MapLit) Type() types.Type { return m.ExprType }

// SetLit represents a Set literal expression.
type SetLit struct {
	ExprType types.Type // *types.Set
}

func (s *SetLit) exprNode()        {}
func (s *SetLit) Type() types.Type { return s.ExprType }

// MethodCallExpr represents a method call on an object (e.g., map.get("key")).
type MethodCallExpr struct {
	Object   Expr
	Method   string
	Args     []Expr
	ExprType types.Type
}

func (m *MethodCallExpr) exprNode()        {}
func (m *MethodCallExpr) Type() types.Type { return m.ExprType }

// EnumMemberExpr represents access to an enum member (e.g., Color.Red).
type EnumMemberExpr struct {
	EnumName   string     // Name of the enum (e.g., "Color")
	MemberName string     // Name of the member (e.g., "Red")
	ExprType   types.Type // The enum type
}

func (e *EnumMemberExpr) exprNode()        {}
func (e *EnumMemberExpr) Type() types.Type { return e.ExprType }

// tokenOpToString converts token.Type to operator string.
func TokenOpToString(op ast.Node) string {
	return op.TokenLiteral()
}
