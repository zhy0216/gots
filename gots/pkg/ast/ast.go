// Package ast defines the Abstract Syntax Tree nodes for goTS.
package ast

import (
	"bytes"
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

// TemplateLiteral represents a template literal (e.g., `Hello, ${name}!`).
type TemplateLiteral struct {
	Token       token.Token  // The TEMPLATE_LITERAL, TEMPLATE_HEAD, or TEMPLATE_MIDDLE token
	Parts       []string     // Static string parts
	Expressions []Expression // Interpolated expressions
}

func (t *TemplateLiteral) expressionNode()      {}
func (t *TemplateLiteral) TokenLiteral() string { return t.Token.Literal }
func (t *TemplateLiteral) String() string {
	var out strings.Builder
	out.WriteString("`")
	for i, part := range t.Parts {
		out.WriteString(part)
		if i < len(t.Expressions) {
			out.WriteString("${")
			out.WriteString(t.Expressions[i].String())
			out.WriteString("}")
		}
	}
	out.WriteString("`")
	return out.String()
}

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
	IsAsync    bool // true if declared with 'async' keyword
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

// NewExpr represents a new expression (e.g., new Point(1, 2) or new Map<string, int>()).
type NewExpr struct {
	Token     token.Token // The 'new' token
	ClassName string
	TypeArgs  []Type // Type arguments for generic types (e.g., <string, int> for Map)
	Arguments []Expression
}

func (n *NewExpr) expressionNode()      {}
func (n *NewExpr) TokenLiteral() string { return n.Token.Literal }
func (n *NewExpr) String() string {
	args := make([]string, len(n.Arguments))
	for i, a := range n.Arguments {
		args[i] = a.String()
	}
	typeArgsStr := ""
	if len(n.TypeArgs) > 0 {
		typeArgs := make([]string, len(n.TypeArgs))
		for i, t := range n.TypeArgs {
			typeArgs[i] = t.String()
		}
		typeArgsStr = "<" + strings.Join(typeArgs, ", ") + ">"
	}
	return fmt.Sprintf("new %s%s(%s)", n.ClassName, typeArgsStr, strings.Join(args, ", "))
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
	IsAsync    bool       // true if declared with 'async' keyword
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
	// Destructuring support
	Pattern Pattern // Alternative to Name for destructuring
}

func (v *VarDecl) statementNode()       {}
func (v *VarDecl) TokenLiteral() string { return v.Token.Literal }
func (v *VarDecl) String() string {
	keyword := "let"
	if v.IsConst {
		keyword = "const"
	}
	if v.Pattern != nil {
		return fmt.Sprintf("%s %s = %s", keyword, v.Pattern.String(), v.Value.String())
	}
	return fmt.Sprintf("%s %s: %s = %s", keyword, v.Name, v.VarType.String(), v.Value.String())
}

// ----------------------------------------------------------------------------
// Patterns (for destructuring)
// ----------------------------------------------------------------------------

// Pattern is the interface for destructuring patterns.
type Pattern interface {
	Node
	patternNode()
}

// ArrayPattern represents array destructuring: [a, b, c]
type ArrayPattern struct {
	Token    token.Token // The '[' token
	Elements []Pattern   // Can be IdentPattern or nested patterns
}

func (a *ArrayPattern) patternNode()         {}
func (a *ArrayPattern) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayPattern) String() string {
	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		if e != nil {
			elements[i] = e.String()
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// ObjectPattern represents object destructuring: {x, y} or {x: newX, y: newY}
type ObjectPattern struct {
	Token      token.Token // The '{' token
	Properties []*PropertyPattern
}

func (o *ObjectPattern) patternNode()         {}
func (o *ObjectPattern) TokenLiteral() string { return o.Token.Literal }
func (o *ObjectPattern) String() string {
	props := make([]string, len(o.Properties))
	for i, p := range o.Properties {
		props[i] = p.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(props, ", "))
}

// PropertyPattern represents a property in object destructuring.
type PropertyPattern struct {
	Key   string  // Original property name
	Value Pattern // Target pattern (can be IdentPattern or nested)
}

func (p *PropertyPattern) String() string {
	if ident, ok := p.Value.(*IdentPattern); ok && ident.Name == p.Key {
		return p.Key // Shorthand: {x} instead of {x: x}
	}
	return fmt.Sprintf("%s: %s", p.Key, p.Value.String())
}

// IdentPattern represents an identifier in a destructuring pattern.
type IdentPattern struct {
	Token token.Token
	Name  string
}

func (i *IdentPattern) patternNode()         {}
func (i *IdentPattern) TokenLiteral() string { return i.Token.Literal }
func (i *IdentPattern) String() string       { return i.Name }

// SpreadExpr represents a spread expression (...arr).
type SpreadExpr struct {
	Token    token.Token // The '...' token
	Argument Expression  // The expression being spread
}

func (s *SpreadExpr) expressionNode()      {}
func (s *SpreadExpr) TokenLiteral() string { return s.Token.Literal }
func (s *SpreadExpr) String() string       { return fmt.Sprintf("...%s", s.Argument.String()) }

// AwaitExpr represents an await expression.
type AwaitExpr struct {
	Token    token.Token // The 'await' token
	Argument Expression  // The promise being awaited
}

func (a *AwaitExpr) expressionNode()      {}
func (a *AwaitExpr) TokenLiteral() string { return a.Token.Literal }
func (a *AwaitExpr) String() string       { return fmt.Sprintf("await %s", a.Argument.String()) }

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

// TryStmt represents a try/catch statement.
type TryStmt struct {
	Token      token.Token // The 'try' token
	TryBlock   *Block
	CatchParam string // The catch parameter name (e.g., 'e' in catch(e))
	CatchBlock *Block
}

func (t *TryStmt) statementNode()       {}
func (t *TryStmt) TokenLiteral() string { return "try" }
func (t *TryStmt) String() string {
	return fmt.Sprintf("try %s catch(%s) %s", t.TryBlock.String(), t.CatchParam, t.CatchBlock.String())
}

// ThrowStmt represents a throw statement.
type ThrowStmt struct {
	Token token.Token // The 'throw' token
	Value Expression
}

func (t *ThrowStmt) statementNode()       {}
func (t *ThrowStmt) TokenLiteral() string { return "throw" }
func (t *ThrowStmt) String() string {
	return fmt.Sprintf("throw %s;", t.Value.String())
}

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
	TypeParams []*TypeParam // Generic type parameters (e.g., <T, U>)
	Params     []*Parameter
	ReturnType Type
	Body       *Block
	IsAsync    bool // true if declared with 'async' keyword
}

func (f *FuncDecl) statementNode()       {}
func (f *FuncDecl) TokenLiteral() string { return f.Token.Literal }
func (f *FuncDecl) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = fmt.Sprintf("%s: %s", p.Name, p.ParamType.String())
	}
	typeParams := ""
	if len(f.TypeParams) > 0 {
		tps := make([]string, len(f.TypeParams))
		for i, tp := range f.TypeParams {
			tps[i] = tp.String()
		}
		typeParams = fmt.Sprintf("<%s>", strings.Join(tps, ", "))
	}
	return fmt.Sprintf("function %s%s(%s): %s %s",
		f.Name, typeParams, strings.Join(params, ", "), f.ReturnType.String(), f.Body.String())
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
	TypeParams  []*TypeParam // Generic type parameters (e.g., <T>)
	SuperClass  string       // Empty if no superclass
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
	if len(c.TypeParams) > 0 {
		tps := make([]string, len(c.TypeParams))
		for i, tp := range c.TypeParams {
			tps[i] = tp.String()
		}
		out.WriteString(fmt.Sprintf("<%s>", strings.Join(tps, ", ")))
	}
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

// EnumDecl represents an enum declaration.
type EnumDecl struct {
	Token   token.Token // The 'enum' token
	Name    string
	Members []*EnumMember
}

// EnumMember represents a single member of an enum.
type EnumMember struct {
	Name  string
	Value Expression // nil if auto-assigned
}

func (e *EnumDecl) statementNode()       {}
func (e *EnumDecl) TokenLiteral() string { return e.Token.Literal }
func (e *EnumDecl) String() string {
	var out bytes.Buffer
	out.WriteString("enum ")
	out.WriteString(e.Name)
	out.WriteString(" { ")
	for i, m := range e.Members {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(m.Name)
		if m.Value != nil {
			out.WriteString(" = ")
			out.WriteString(m.Value.String())
		}
	}
	out.WriteString(" }")
	return out.String()
}

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// PrimitiveKind represents the kind of primitive type.
type PrimitiveKind int

const (
	TypeInt PrimitiveKind = iota
	TypeFloat
	TypeString
	TypeBoolean
	TypeVoid
	TypeNull
)

// PrimitiveType represents a primitive type (int, float, string, boolean, void, null).
type PrimitiveType struct {
	Kind PrimitiveKind
}

func (p *PrimitiveType) typeNode()          {}
func (p *PrimitiveType) TokenLiteral() string { return p.String() }
func (p *PrimitiveType) String() string {
	switch p.Kind {
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
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

// UnionType represents a union of multiple types (e.g., string | int | boolean).
type UnionType struct {
	Types []Type
}

func (u *UnionType) typeNode()            {}
func (u *UnionType) TokenLiteral() string { return u.String() }
func (u *UnionType) String() string {
	types := make([]string, len(u.Types))
	for i, t := range u.Types {
		types[i] = t.String()
	}
	return strings.Join(types, " | ")
}

// IntersectionType represents an intersection of multiple types (e.g., A & B).
type IntersectionType struct {
	Types []Type
}

func (i *IntersectionType) typeNode()            {}
func (i *IntersectionType) TokenLiteral() string { return i.String() }
func (i *IntersectionType) String() string {
	types := make([]string, len(i.Types))
	for idx, t := range i.Types {
		types[idx] = t.String()
	}
	return strings.Join(types, " & ")
}

// LiteralType represents a literal type (e.g., "hello", 42, true).
type LiteralType struct {
	Kind  PrimitiveKind // TypeString, TypeInt, TypeFloat, TypeBoolean
	Value string        // The literal value as a string
}

func (l *LiteralType) typeNode()            {}
func (l *LiteralType) TokenLiteral() string { return l.Value }
func (l *LiteralType) String() string       { return l.Value }

// TupleType represents a tuple type (e.g., [string, number] or [string, ...number[]]).
type TupleType struct {
	Token       token.Token // The '[' token
	Elements    []Type      // Fixed-position element types
	RestElement Type        // Optional rest element type (e.g., int[] in [string, ...int[]])
}

func (t *TupleType) typeNode()            {}
func (t *TupleType) TokenLiteral() string { return t.Token.Literal }
func (t *TupleType) String() string {
	elements := make([]string, len(t.Elements))
	for i, e := range t.Elements {
		elements[i] = e.String()
	}
	if t.RestElement != nil {
		return fmt.Sprintf("[%s, ...%s]", strings.Join(elements, ", "), t.RestElement.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// NamedType represents a reference to a named type (type alias or class).
type NamedType struct {
	Name     string
	TypeArgs []Type // Type arguments for generic types (e.g., Stack<int>)
}

func (n *NamedType) typeNode()            {}
func (n *NamedType) TokenLiteral() string { return n.Name }
func (n *NamedType) String() string {
	if len(n.TypeArgs) == 0 {
		return n.Name
	}
	args := make([]string, len(n.TypeArgs))
	for i, t := range n.TypeArgs {
		args[i] = t.String()
	}
	return fmt.Sprintf("%s<%s>", n.Name, strings.Join(args, ", "))
}

// TypeParam represents a type parameter (e.g., T in function identity<T>).
type TypeParam struct {
	Name       string
	Constraint Type // Optional constraint (e.g., T extends Comparable)
}

func (t *TypeParam) typeNode()            {}
func (t *TypeParam) TokenLiteral() string { return t.Name }
func (t *TypeParam) String() string {
	if t.Constraint != nil {
		return fmt.Sprintf("%s extends %s", t.Name, t.Constraint.String())
	}
	return t.Name
}

// MapType represents a map type (e.g., Map<string, int>).
type MapType struct {
	KeyType   Type
	ValueType Type
}

func (m *MapType) typeNode()            {}
func (m *MapType) TokenLiteral() string { return m.String() }
func (m *MapType) String() string {
	return fmt.Sprintf("Map<%s, %s>", m.KeyType.String(), m.ValueType.String())
}

// SetType represents a set type (e.g., Set<int>).
type SetType struct {
	ElementType Type
}

func (s *SetType) typeNode()            {}
func (s *SetType) TokenLiteral() string { return s.String() }
func (s *SetType) String() string {
	return fmt.Sprintf("Set<%s>", s.ElementType.String())
}

// PromiseType represents a Promise<T> type.
type PromiseType struct {
	ResultType Type // The resolved value type T
}

func (p *PromiseType) typeNode()            {}
func (p *PromiseType) TokenLiteral() string { return p.String() }
func (p *PromiseType) String() string {
	return fmt.Sprintf("Promise<%s>", p.ResultType.String())
}

// InterfaceType represents an interface type reference.
type InterfaceType struct {
	Name string
}

func (i *InterfaceType) typeNode()            {}
func (i *InterfaceType) TokenLiteral() string { return i.Name }
func (i *InterfaceType) String() string       { return i.Name }

// ----------------------------------------------------------------------------
// Interface Declaration
// ----------------------------------------------------------------------------

// InterfaceDecl represents an interface declaration.
type InterfaceDecl struct {
	Token   token.Token
	Name    string
	Methods []*InterfaceMethod
}

func (i *InterfaceDecl) statementNode()       {}
func (i *InterfaceDecl) TokenLiteral() string { return i.Token.Literal }
func (i *InterfaceDecl) String() string {
	var methods []string
	for _, m := range i.Methods {
		methods = append(methods, m.String())
	}
	return fmt.Sprintf("interface %s { %s }", i.Name, strings.Join(methods, "; "))
}

// InterfaceMethod represents a method signature in an interface.
type InterfaceMethod struct {
	Name       string
	Params     []*Parameter
	ReturnType Type
}

func (m *InterfaceMethod) String() string {
	var params []string
	for _, p := range m.Params {
		if p.ParamType != nil {
			params = append(params, fmt.Sprintf("%s: %s", p.Name, p.ParamType.String()))
		} else {
			params = append(params, p.Name)
		}
	}
	if m.ReturnType != nil {
		return fmt.Sprintf("%s(%s): %s", m.Name, strings.Join(params, ", "), m.ReturnType.String())
	}
	return fmt.Sprintf("%s(%s)", m.Name, strings.Join(params, ", "))
}

// GoImportDecl represents an import from a Go package.
// e.g., import { Sprintf, Println } from "go:fmt"
type GoImportDecl struct {
	Token   token.Token // The 'import' token
	Names   []string    // The names being imported
	Package string      // The Go package path (without "go:" prefix)
}

func (g *GoImportDecl) statementNode()       {}
func (g *GoImportDecl) TokenLiteral() string { return g.Token.Literal }
func (g *GoImportDecl) String() string {
	return fmt.Sprintf("import { %s } from \"go:%s\"", strings.Join(g.Names, ", "), g.Package)
}

// ModuleImportDecl represents an import from a local module.
// e.g., import { add, Vector } from "./math"
type ModuleImportDecl struct {
	Token  token.Token // The 'import' token
	Names  []string    // The names being imported
	Path   string      // The module path (e.g., "./math")
}

func (m *ModuleImportDecl) statementNode()       {}
func (m *ModuleImportDecl) TokenLiteral() string { return m.Token.Literal }
func (m *ModuleImportDecl) String() string {
	return fmt.Sprintf("import { %s } from \"%s\"", strings.Join(m.Names, ", "), m.Path)
}

// ExportModifier is a marker that a declaration is exported.
// It wraps another statement (FuncDecl, ClassDecl, VarDecl, TypeAliasDecl).
type ExportModifier struct {
	Token   token.Token
	Decl    Statement // The declaration being exported
}

func (e *ExportModifier) statementNode()       {}
func (e *ExportModifier) TokenLiteral() string { return e.Token.Literal }
func (e *ExportModifier) String() string {
	return "export " + e.Decl.String()
}

// ReExportDecl represents a re-export statement.
// e.g., export { foo, bar } from "./module" or export * from "./module"
type ReExportDecl struct {
	Token      token.Token // The 'export' token
	Names      []string    // The names being re-exported (empty for wildcard)
	Path       string      // The module path
	IsWildcard bool        // true for "export *"
}

func (r *ReExportDecl) statementNode()       {}
func (r *ReExportDecl) TokenLiteral() string { return r.Token.Literal }
func (r *ReExportDecl) String() string {
	if r.IsWildcard {
		return fmt.Sprintf("export * from \"%s\"", r.Path)
	}
	return fmt.Sprintf("export { %s } from \"%s\"", strings.Join(r.Names, ", "), r.Path)
}

// DefaultExport represents a default export statement.
// e.g., export default class Foo {} or export default function() {}
type DefaultExport struct {
	Token token.Token // The 'export' token
	Decl  Statement   // The declaration being exported (can be class, function, or expression)
}

func (d *DefaultExport) statementNode()       {}
func (d *DefaultExport) TokenLiteral() string { return d.Token.Literal }
func (d *DefaultExport) String() string {
	return "export default " + d.Decl.String()
}

// DefaultImport represents a default import statement.
// e.g., import Foo from "./module"
type DefaultImport struct {
	Token token.Token // The 'import' token
	Name  string      // The name to bind the default export
	Path  string      // The module path
}

func (d *DefaultImport) statementNode()       {}
func (d *DefaultImport) TokenLiteral() string { return d.Token.Literal }
func (d *DefaultImport) String() string {
	return fmt.Sprintf("import %s from \"%s\"", d.Name, d.Path)
}

// NamespaceImport represents a namespace import statement.
// e.g., import * as utils from "./utils"
type NamespaceImport struct {
	Token token.Token // The 'import' token
	Alias string      // The namespace alias
	Path  string      // The module path
}

func (n *NamespaceImport) statementNode()       {}
func (n *NamespaceImport) TokenLiteral() string { return n.Token.Literal }
func (n *NamespaceImport) String() string {
	return fmt.Sprintf("import * as %s from \"%s\"", n.Alias, n.Path)
}
