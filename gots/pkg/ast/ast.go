// Package ast defines the Abstract Syntax Tree nodes for GoTS.
package ast

import (
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

// Node is the interface that all AST nodes implement.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is the interface for statement nodes.
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface for expression nodes.
type Expression interface {
	Node
	expressionNode()
}

// Type is the interface for type nodes.
type Type interface {
	Node
	typeNode()
}

// ----------------------------------------------------------------------------
// Program
// ----------------------------------------------------------------------------

// Program is the root node of the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ----------------------------------------------------------------------------
// Expressions
// ----------------------------------------------------------------------------

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Token token.Token
	Value float64
}

func (n *NumberLiteral) expressionNode()      {}
func (n *NumberLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NumberLiteral) String() string       { return n.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return fmt.Sprintf("%q", s.Value) }

// BoolLiteral represents a boolean literal (true/false).
type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (b *BoolLiteral) expressionNode()      {}
func (b *BoolLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BoolLiteral) String() string       { return b.Token.Literal }

// NullLiteral represents the null literal.
type NullLiteral struct {
	Token token.Token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return "null" }
func (n *NullLiteral) String() string       { return "null" }

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token
	Name  string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Name }

// BinaryExpr represents a binary expression (e.g., a + b).
type BinaryExpr struct {
	Token token.Token // The operator token
	Left  Expression
	Op    token.Type
	Right Expression
}

func (b *BinaryExpr) expressionNode()      {}
func (b *BinaryExpr) TokenLiteral() string { return b.Token.Literal }
func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Token.Literal, b.Right.String())
}

// UnaryExpr represents a unary expression (e.g., -x, !flag).
type UnaryExpr struct {
	Token   token.Token // The operator token
	Op      token.Type
	Operand Expression
}

func (u *UnaryExpr) expressionNode()      {}
func (u *UnaryExpr) TokenLiteral() string { return u.Token.Literal }
func (u *UnaryExpr) String() string {
	return fmt.Sprintf("(%s%s)", u.Token.Literal, u.Operand.String())
}

// CallExpr represents a function call expression (e.g., fn() or fn?.()).
type CallExpr struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or PropertyExpr
	Arguments []Expression
	Optional  bool // true for optional chaining ?.()
}

func (c *CallExpr) expressionNode()      {}
func (c *CallExpr) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpr) String() string {
	args := make([]string, len(c.Arguments))
	for i, a := range c.Arguments {
		args[i] = a.String()
	}
	if c.Optional {
		return fmt.Sprintf("%s?.(%s)", c.Function.String(), strings.Join(args, ", "))
	}
	return fmt.Sprintf("%s(%s)", c.Function.String(), strings.Join(args, ", "))
}

// IndexExpr represents an index expression (e.g., arr[0] or arr?.[0]).
type IndexExpr struct {
	Token    token.Token // The '[' token
	Object   Expression
	Index    Expression
	Optional bool // true for optional chaining (?.[])
}

func (i *IndexExpr) expressionNode()      {}
func (i *IndexExpr) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpr) String() string {
	if i.Optional {
		return fmt.Sprintf("%s?.[%s]", i.Object.String(), i.Index.String())
	}
	return fmt.Sprintf("%s[%s]", i.Object.String(), i.Index.String())
}

// PropertyExpr represents a property access expression (e.g., obj.x or obj?.x).
type PropertyExpr struct {
	Token    token.Token // The '.' or '?.' token
	Object   Expression
	Property string
	Optional bool // true for optional chaining (?.)
}

func (p *PropertyExpr) expressionNode()      {}
func (p *PropertyExpr) TokenLiteral() string { return p.Token.Literal }
func (p *PropertyExpr) String() string {
	if p.Optional {
		return fmt.Sprintf("%s?.%s", p.Object.String(), p.Property)
	}
	return fmt.Sprintf("%s.%s", p.Object.String(), p.Property)
}

// ArrayLiteral represents an array literal (e.g., [1, 2, 3]).
type ArrayLiteral struct {
	Token    token.Token // The '[' token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// PropertyDef represents a property definition in an object literal.
type PropertyDef struct {
	Key   string
	Value Expression
}

// ObjectLiteral represents an object literal (e.g., {x: 1, y: 2}).
type ObjectLiteral struct {
	Token      token.Token // The '{' token
	Properties []*PropertyDef
}

func (o *ObjectLiteral) expressionNode()      {}
func (o *ObjectLiteral) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectLiteral) String() string {
	props := make([]string, len(o.Properties))
	for i, p := range o.Properties {
		props[i] = fmt.Sprintf("%s: %s", p.Key, p.Value.String())
	}
	return fmt.Sprintf("{%s}", strings.Join(props, ", "))
}

// FunctionExpr represents a function expression.
type FunctionExpr struct {
	Token      token.Token // The 'function' token
	Params     []*Parameter
	ReturnType Type
	Body       *Block
}

func (f *FunctionExpr) expressionNode()      {}
func (f *FunctionExpr) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionExpr) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = fmt.Sprintf("%s: %s", p.Name, p.ParamType.String())
	}
	return fmt.Sprintf("function(%s): %s { ... }", strings.Join(params, ", "), f.ReturnType.String())
}

// NewExpr represents a new expression (e.g., new Point(1, 2)).
type NewExpr struct {
	Token     token.Token // The 'new' token
	ClassName string
	Arguments []Expression
}

func (n *NewExpr) expressionNode()      {}
func (n *NewExpr) TokenLiteral() string { return n.Token.Literal }
func (n *NewExpr) String() string {
	args := make([]string, len(n.Arguments))
	for i, a := range n.Arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("new %s(%s)", n.ClassName, strings.Join(args, ", "))
}

// ThisExpr represents the 'this' keyword.
type ThisExpr struct {
	Token token.Token
}

func (t *ThisExpr) expressionNode()      {}
func (t *ThisExpr) TokenLiteral() string { return "this" }
func (t *ThisExpr) String() string       { return "this" }

// SuperExpr represents a super call.
type SuperExpr struct {
	Token     token.Token // The 'super' token
	Arguments []Expression
}

func (s *SuperExpr) expressionNode()      {}
func (s *SuperExpr) TokenLiteral() string { return "super" }
func (s *SuperExpr) String() string {
	args := make([]string, len(s.Arguments))
	for i, a := range s.Arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("super(%s)", strings.Join(args, ", "))
}

// AssignExpr represents an assignment expression.
type AssignExpr struct {
	Token  token.Token // The '=' token
	Target Expression  // Identifier, IndexExpr, or PropertyExpr
	Value  Expression
}

func (a *AssignExpr) expressionNode()      {}
func (a *AssignExpr) TokenLiteral() string { return a.Token.Literal }
func (a *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", a.Target.String(), a.Value.String())
}

// CompoundAssignExpr represents a compound assignment expression (+=, -=, etc.).
type CompoundAssignExpr struct {
	Token  token.Token // The operator token (+=, -=, etc.)
	Target Expression  // Identifier, IndexExpr, or PropertyExpr
	Op     token.Type  // The compound operator
	Value  Expression
}

func (c *CompoundAssignExpr) expressionNode()      {}
func (c *CompoundAssignExpr) TokenLiteral() string { return c.Token.Literal }
func (c *CompoundAssignExpr) String() string {
	return fmt.Sprintf("%s %s %s", c.Target.String(), c.Token.Literal, c.Value.String())
}

// ArrowFunctionExpr represents an arrow function expression.
type ArrowFunctionExpr struct {
	Token      token.Token // The '=>' token
	Params     []*Parameter
	ReturnType Type
	Body       *Block     // For block body: () => { ... }
	Expression Expression // For expression body: () => expr
}

func (a *ArrowFunctionExpr) expressionNode()      {}
func (a *ArrowFunctionExpr) TokenLiteral() string { return a.Token.Literal }
func (a *ArrowFunctionExpr) String() string {
	params := make([]string, len(a.Params))
	for i, p := range a.Params {
		params[i] = fmt.Sprintf("%s: %s", p.Name, p.ParamType.String())
	}
	if a.Body != nil {
		return fmt.Sprintf("(%s): %s => { ... }", strings.Join(params, ", "), a.ReturnType.String())
	}
	return fmt.Sprintf("(%s): %s => %s", strings.Join(params, ", "), a.ReturnType.String(), a.Expression.String())
}

// UpdateExpr represents an increment/decrement expression (++x, x++, --x, x--).
type UpdateExpr struct {
	Token   token.Token // The ++ or -- token
	Op      token.Type  // INCREMENT or DECREMENT
	Operand Expression  // Must be assignable (identifier, property, index)
	Prefix  bool        // true for ++x, false for x++
}

func (u *UpdateExpr) expressionNode()      {}
func (u *UpdateExpr) TokenLiteral() string { return u.Token.Literal }
func (u *UpdateExpr) String() string {
	if u.Prefix {
		return fmt.Sprintf("(%s%s)", u.Token.Literal, u.Operand.String())
	}
	return fmt.Sprintf("(%s%s)", u.Operand.String(), u.Token.Literal)
}

// ----------------------------------------------------------------------------
// Statements
// ----------------------------------------------------------------------------

// ExprStmt wraps an expression as a statement.
type ExprStmt struct {
	Token token.Token
	Expr  Expression
}

func (e *ExprStmt) statementNode()       {}
func (e *ExprStmt) TokenLiteral() string { return e.Token.Literal }
func (e *ExprStmt) String() string {
	if e.Expr != nil {
		return e.Expr.String() + ";"
	}
	return ""
}

// VarDecl represents a variable declaration (let or const).
type VarDecl struct {
	Token   token.Token // The 'let' or 'const' token
	Name    string
	VarType Type
	Value   Expression
	IsConst bool
}

func (v *VarDecl) statementNode()       {}
func (v *VarDecl) TokenLiteral() string { return v.Token.Literal }
func (v *VarDecl) String() string {
	keyword := "let"
	if v.IsConst {
		keyword = "const"
	}
	return fmt.Sprintf("%s %s: %s = %s", keyword, v.Name, v.VarType.String(), v.Value.String())
}

// Block represents a block of statements.
type Block struct {
	Token      token.Token // The '{' token
	Statements []Statement
}

func (b *Block) statementNode()       {}
func (b *Block) TokenLiteral() string { return b.Token.Literal }
func (b *Block) String() string {
	var out strings.Builder
	out.WriteString("{ ")
	for _, s := range b.Statements {
		out.WriteString(s.String())
		out.WriteString(" ")
	}
	out.WriteString("}")
	return out.String()
}

// IfStmt represents an if statement.
type IfStmt struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *Block
	Alternative Statement // *Block or *IfStmt (else if)
}

func (i *IfStmt) statementNode()       {}
func (i *IfStmt) TokenLiteral() string { return i.Token.Literal }
func (i *IfStmt) String() string {
	var out strings.Builder
	out.WriteString("if (")
	out.WriteString(i.Condition.String())
	out.WriteString(") ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

// WhileStmt represents a while statement.
type WhileStmt struct {
	Token     token.Token // The 'while' token
	Condition Expression
	Body      *Block
}

func (w *WhileStmt) statementNode()       {}
func (w *WhileStmt) TokenLiteral() string { return w.Token.Literal }
func (w *WhileStmt) String() string {
	return fmt.Sprintf("while (%s) %s", w.Condition.String(), w.Body.String())
}

// ForStmt represents a for statement.
type ForStmt struct {
	Token     token.Token // The 'for' token
	Init      *VarDecl
	Condition Expression
	Update    Expression
	Body      *Block
}

func (f *ForStmt) statementNode()       {}
func (f *ForStmt) TokenLiteral() string { return f.Token.Literal }
func (f *ForStmt) String() string {
	return fmt.Sprintf("for (%s; %s; %s) %s",
		f.Init.String(), f.Condition.String(), f.Update.String(), f.Body.String())
}

// ForOfStmt represents a for-of statement.
type ForOfStmt struct {
	Token    token.Token // The 'for' token
	Variable *VarDecl    // The loop variable declaration
	Iterable Expression  // The iterable expression
	Body     *Block
}

func (f *ForOfStmt) statementNode()       {}
func (f *ForOfStmt) TokenLiteral() string { return f.Token.Literal }
func (f *ForOfStmt) String() string {
	return fmt.Sprintf("for (let %s of %s) %s",
		f.Variable.Name, f.Iterable.String(), f.Body.String())
}

// SwitchStmt represents a switch statement.
type SwitchStmt struct {
	Token        token.Token // The 'switch' token
	Discriminant Expression
	Cases        []*CaseClause
}

func (s *SwitchStmt) statementNode()       {}
func (s *SwitchStmt) TokenLiteral() string { return s.Token.Literal }
func (s *SwitchStmt) String() string {
	return fmt.Sprintf("switch (%s) { ... }", s.Discriminant.String())
}

// CaseClause represents a case or default clause in a switch statement.
type CaseClause struct {
	Token      token.Token // The 'case' or 'default' token
	Test       Expression  // nil for default case
	Consequent []Statement
}

func (c *CaseClause) String() string {
	if c.Test == nil {
		return "default: ..."
	}
	return fmt.Sprintf("case %s: ...", c.Test.String())
}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Token token.Token // The 'return' token
	Value Expression  // nil for void return
}

func (r *ReturnStmt) statementNode()       {}
func (r *ReturnStmt) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStmt) String() string {
	if r.Value != nil {
		return fmt.Sprintf("return %s;", r.Value.String())
	}
	return "return;"
}

// BreakStmt represents a break statement.
type BreakStmt struct {
	Token token.Token
}

func (b *BreakStmt) statementNode()       {}
func (b *BreakStmt) TokenLiteral() string { return "break" }
func (b *BreakStmt) String() string       { return "break;" }

// ContinueStmt represents a continue statement.
type ContinueStmt struct {
	Token token.Token
}

func (c *ContinueStmt) statementNode()       {}
func (c *ContinueStmt) TokenLiteral() string { return "continue" }
func (c *ContinueStmt) String() string       { return "continue;" }

// ----------------------------------------------------------------------------
// Declarations
// ----------------------------------------------------------------------------

// Parameter represents a function parameter.
type Parameter struct {
	Name      string
	ParamType Type
}

// FuncDecl represents a function declaration.
type FuncDecl struct {
	Token      token.Token // The 'function' token
	Name       string
	Params     []*Parameter
	ReturnType Type
	Body       *Block
}

func (f *FuncDecl) statementNode()       {}
func (f *FuncDecl) TokenLiteral() string { return f.Token.Literal }
func (f *FuncDecl) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = fmt.Sprintf("%s: %s", p.Name, p.ParamType.String())
	}
	return fmt.Sprintf("function %s(%s): %s %s",
		f.Name, strings.Join(params, ", "), f.ReturnType.String(), f.Body.String())
}

// Field represents a class field.
type Field struct {
	Name      string
	FieldType Type
}

// Method represents a class method.
type Method struct {
	Name       string
	Params     []*Parameter
	ReturnType Type
	Body       *Block
}

// Constructor represents a class constructor.
type Constructor struct {
	Params []*Parameter
	Body   *Block
}

// ClassDecl represents a class declaration.
type ClassDecl struct {
	Token       token.Token // The 'class' token
	Name        string
	SuperClass  string // Empty if no superclass
	Fields      []*Field
	Constructor *Constructor
	Methods     []*Method
}

func (c *ClassDecl) statementNode()       {}
func (c *ClassDecl) TokenLiteral() string { return c.Token.Literal }
func (c *ClassDecl) String() string {
	var out strings.Builder
	out.WriteString("class ")
	out.WriteString(c.Name)
	if c.SuperClass != "" {
		out.WriteString(" extends ")
		out.WriteString(c.SuperClass)
	}
	out.WriteString(" { ... }")
	return out.String()
}

// TypeAliasDecl represents a type alias declaration.
type TypeAliasDecl struct {
	Token     token.Token // The 'type' token
	Name      string
	AliasType Type
}

func (t *TypeAliasDecl) statementNode()       {}
func (t *TypeAliasDecl) TokenLiteral() string { return t.Token.Literal }
func (t *TypeAliasDecl) String() string {
	return fmt.Sprintf("type %s = %s;", t.Name, t.AliasType.String())
}

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// PrimitiveKind represents the kind of primitive type.
type PrimitiveKind int

const (
	TypeNumber PrimitiveKind = iota
	TypeString
	TypeBoolean
	TypeVoid
	TypeNull
)

// PrimitiveType represents a primitive type (number, string, boolean, void, null).
type PrimitiveType struct {
	Kind PrimitiveKind
}

func (p *PrimitiveType) typeNode()          {}
func (p *PrimitiveType) TokenLiteral() string { return p.String() }
func (p *PrimitiveType) String() string {
	switch p.Kind {
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeBoolean:
		return "boolean"
	case TypeVoid:
		return "void"
	case TypeNull:
		return "null"
	}
	return "unknown"
}

// ArrayType represents an array type (e.g., number[]).
type ArrayType struct {
	ElementType Type
}

func (a *ArrayType) typeNode()            {}
func (a *ArrayType) TokenLiteral() string { return a.String() }
func (a *ArrayType) String() string {
	return fmt.Sprintf("%s[]", a.ElementType.String())
}

// ObjectTypeProperty represents a property in an object type.
type ObjectTypeProperty struct {
	Name      string
	PropType  Type
}

// ObjectType represents an object type (e.g., { x: number, y: number }).
type ObjectType struct {
	Properties []*ObjectTypeProperty
}

func (o *ObjectType) typeNode()            {}
func (o *ObjectType) TokenLiteral() string { return o.String() }
func (o *ObjectType) String() string {
	props := make([]string, len(o.Properties))
	for i, p := range o.Properties {
		props[i] = fmt.Sprintf("%s: %s", p.Name, p.PropType.String())
	}
	return fmt.Sprintf("{ %s }", strings.Join(props, ", "))
}

// FunctionType represents a function type (e.g., (a: number) => string).
type FunctionType struct {
	ParamTypes []Type
	ReturnType Type
}

func (f *FunctionType) typeNode()            {}
func (f *FunctionType) TokenLiteral() string { return f.String() }
func (f *FunctionType) String() string {
	params := make([]string, len(f.ParamTypes))
	for i, p := range f.ParamTypes {
		params[i] = p.String()
	}
	return fmt.Sprintf("(%s) => %s", strings.Join(params, ", "), f.ReturnType.String())
}

// NullableType represents a nullable type (e.g., string | null).
type NullableType struct {
	Inner Type
}

func (n *NullableType) typeNode()            {}
func (n *NullableType) TokenLiteral() string { return n.String() }
func (n *NullableType) String() string {
	return fmt.Sprintf("%s | null", n.Inner.String())
}

// NamedType represents a reference to a named type (type alias or class).
type NamedType struct {
	Name string
}

func (n *NamedType) typeNode()            {}
func (n *NamedType) TokenLiteral() string { return n.Name }
func (n *NamedType) String() string       { return n.Name }
