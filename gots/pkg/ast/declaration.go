// Package ast - Declaration file AST nodes for .d.gts files
package ast

import (
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

// ----------------------------------------------------------------------------
// Declaration File AST Nodes
// These nodes are used for .d.gts declaration files that describe Go packages
// ----------------------------------------------------------------------------

// DeclarationFile represents a complete .d.gts declaration file.
type DeclarationFile struct {
	Modules []*DeclareModule
}

func (d *DeclarationFile) TokenLiteral() string {
	if len(d.Modules) > 0 {
		return d.Modules[0].TokenLiteral()
	}
	return ""
}

func (d *DeclarationFile) String() string {
	var out strings.Builder
	for _, m := range d.Modules {
		out.WriteString(m.String())
		out.WriteString("\n")
	}
	return out.String()
}

// DeclareModule represents a module declaration block.
// e.g., declare module "go:strings" { ... }
type DeclareModule struct {
	Token   token.Token // The 'declare' token
	Name    string      // The module name (e.g., "go:strings")
	Members []Declaration
}

func (d *DeclareModule) statementNode()       {}
func (d *DeclareModule) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareModule) String() string {
	var out strings.Builder
	out.WriteString("declare module \"")
	out.WriteString(d.Name)
	out.WriteString("\" {\n")
	for _, m := range d.Members {
		out.WriteString("    ")
		out.WriteString(m.String())
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}

// Declaration is the interface for all declaration members.
type Declaration interface {
	Node
	declarationNode()
}

// DeclareFunction represents a function declaration.
// e.g., function ToUpper(s: string): string
type DeclareFunction struct {
	Token      token.Token // The 'function' token
	Name       string
	Params     []*Parameter
	ReturnType Type
	Variadic   bool // true if last param is ...rest
}

func (d *DeclareFunction) declarationNode()     {}
func (d *DeclareFunction) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareFunction) String() string {
	params := make([]string, len(d.Params))
	for i, p := range d.Params {
		prefix := ""
		if d.Variadic && i == len(d.Params)-1 {
			prefix = "..."
		}
		params[i] = fmt.Sprintf("%s%s: %s", prefix, p.Name, p.ParamType.String())
	}
	return fmt.Sprintf("function %s(%s): %s", d.Name, strings.Join(params, ", "), d.ReturnType.String())
}

// DeclareType represents a type alias declaration.
// e.g., type Error = { Error(): string }
type DeclareType struct {
	Token     token.Token // The 'type' token
	Name      string
	AliasType Type
}

func (d *DeclareType) declarationNode()     {}
func (d *DeclareType) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareType) String() string {
	return fmt.Sprintf("type %s = %s", d.Name, d.AliasType.String())
}

// DeclareInterface represents an interface declaration.
// e.g., interface Reader { Read(p: byte[]): (int, Error | null) }
type DeclareInterface struct {
	Token   token.Token // The 'interface' token
	Name    string
	Methods []*InterfaceMethod
}

func (d *DeclareInterface) declarationNode()     {}
func (d *DeclareInterface) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareInterface) String() string {
	var methods []string
	for _, m := range d.Methods {
		methods = append(methods, m.String())
	}
	return fmt.Sprintf("interface %s { %s }", d.Name, strings.Join(methods, "; "))
}

// DeclareConst represents a constant declaration.
// e.g., const Pi: float
type DeclareConst struct {
	Token     token.Token // The 'const' token
	Name      string
	ConstType Type
}

func (d *DeclareConst) declarationNode()     {}
func (d *DeclareConst) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareConst) String() string {
	return fmt.Sprintf("const %s: %s", d.Name, d.ConstType.String())
}

// DeclareClass represents a class declaration.
// e.g., class Regexp { Test(s: string): boolean }
type DeclareClass struct {
	Token      token.Token // The 'class' token
	Name       string
	SuperClass string // Empty if no superclass
	Fields     []*Field
	Methods    []*Method
}

func (d *DeclareClass) declarationNode()     {}
func (d *DeclareClass) TokenLiteral() string { return d.Token.Literal }
func (d *DeclareClass) String() string {
	var out strings.Builder
	out.WriteString("class ")
	out.WriteString(d.Name)
	if d.SuperClass != "" {
		out.WriteString(" extends ")
		out.WriteString(d.SuperClass)
	}
	out.WriteString(" { ... }")
	return out.String()
}

// TupleType represents a tuple type for Go's multiple return values.
// e.g., (int, Error | null)
type TupleType struct {
	Elements []Type
}

func (t *TupleType) typeNode()            {}
func (t *TupleType) TokenLiteral() string { return t.String() }
func (t *TupleType) String() string {
	elems := make([]string, len(t.Elements))
	for i, e := range t.Elements {
		elems[i] = e.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(elems, ", "))
}

// AnyType represents the 'any' type.
type AnyType struct{}

func (a *AnyType) typeNode()            {}
func (a *AnyType) TokenLiteral() string { return "any" }
func (a *AnyType) String() string       { return "any" }

// ByteType represents the 'byte' type (alias for int in Go).
type ByteType struct{}

func (b *ByteType) typeNode()            {}
func (b *ByteType) TokenLiteral() string { return "byte" }
func (b *ByteType) String() string       { return "byte" }
