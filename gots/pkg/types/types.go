// Package types implements the type system for GoTS.
package types

import (
	"fmt"
	"strings"
)

// Type is the interface that all types implement.
type Type interface {
	typeNode()
	String() string
	Equals(Type) bool
}

// ----------------------------------------------------------------------------
// Primitive Types
// ----------------------------------------------------------------------------

// PrimitiveKind represents the kind of primitive type.
type PrimitiveKind int

const (
	KindNumber PrimitiveKind = iota
	KindString
	KindBoolean
	KindVoid
	KindNull
	KindAny  // Used for error recovery
	KindNever // Bottom type
)

// Primitive represents a primitive type (number, string, boolean, void, null).
type Primitive struct {
	Kind PrimitiveKind
}

func (p *Primitive) typeNode() {}
func (p *Primitive) String() string {
	switch p.Kind {
	case KindNumber:
		return "number"
	case KindString:
		return "string"
	case KindBoolean:
		return "boolean"
	case KindVoid:
		return "void"
	case KindNull:
		return "null"
	case KindAny:
		return "any"
	case KindNever:
		return "never"
	}
	return "unknown"
}

func (p *Primitive) Equals(other Type) bool {
	if o, ok := other.(*Primitive); ok {
		return p.Kind == o.Kind
	}
	return false
}

// Convenience constructors for primitive types
var (
	NumberType  = &Primitive{Kind: KindNumber}
	StringType  = &Primitive{Kind: KindString}
	BooleanType = &Primitive{Kind: KindBoolean}
	VoidType    = &Primitive{Kind: KindVoid}
	NullType    = &Primitive{Kind: KindNull}
	AnyType     = &Primitive{Kind: KindAny}
	NeverType   = &Primitive{Kind: KindNever}
)

// ----------------------------------------------------------------------------
// Array Type
// ----------------------------------------------------------------------------

// Array represents an array type (e.g., number[]).
type Array struct {
	Element Type
}

func (a *Array) typeNode() {}
func (a *Array) String() string {
	return fmt.Sprintf("%s[]", a.Element.String())
}

func (a *Array) Equals(other Type) bool {
	if o, ok := other.(*Array); ok {
		return a.Element.Equals(o.Element)
	}
	return false
}

// ----------------------------------------------------------------------------
// Object Type
// ----------------------------------------------------------------------------

// Property represents a property in an object type.
type Property struct {
	Name     string
	Type     Type
	Optional bool
}

// Object represents an object type (e.g., { x: number, y: number }).
type Object struct {
	Properties map[string]*Property
}

func (o *Object) typeNode() {}
func (o *Object) String() string {
	if len(o.Properties) == 0 {
		return "{}"
	}
	props := make([]string, 0, len(o.Properties))
	for name, prop := range o.Properties {
		optMark := ""
		if prop.Optional {
			optMark = "?"
		}
		props = append(props, fmt.Sprintf("%s%s: %s", name, optMark, prop.Type.String()))
	}
	return fmt.Sprintf("{ %s }", strings.Join(props, ", "))
}

func (o *Object) Equals(other Type) bool {
	if ot, ok := other.(*Object); ok {
		if len(o.Properties) != len(ot.Properties) {
			return false
		}
		for name, prop := range o.Properties {
			otherProp, exists := ot.Properties[name]
			if !exists {
				return false
			}
			if prop.Optional != otherProp.Optional {
				return false
			}
			if !prop.Type.Equals(otherProp.Type) {
				return false
			}
		}
		return true
	}
	return false
}

// GetProperty returns the type of a property, or nil if not found.
func (o *Object) GetProperty(name string) *Property {
	return o.Properties[name]
}

// ----------------------------------------------------------------------------
// Function Type
// ----------------------------------------------------------------------------

// Param represents a function parameter type.
type Param struct {
	Name string
	Type Type
}

// Function represents a function type.
type Function struct {
	Params     []*Param
	ReturnType Type
}

func (f *Function) typeNode() {}
func (f *Function) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.Type.String()
	}
	return fmt.Sprintf("(%s) => %s", strings.Join(params, ", "), f.ReturnType.String())
}

func (f *Function) Equals(other Type) bool {
	if of, ok := other.(*Function); ok {
		if len(f.Params) != len(of.Params) {
			return false
		}
		for i, p := range f.Params {
			if !p.Type.Equals(of.Params[i].Type) {
				return false
			}
		}
		return f.ReturnType.Equals(of.ReturnType)
	}
	return false
}

// ----------------------------------------------------------------------------
// Nullable Type
// ----------------------------------------------------------------------------

// Nullable represents a nullable type (e.g., string | null).
type Nullable struct {
	Inner Type
}

func (n *Nullable) typeNode() {}
func (n *Nullable) String() string {
	return fmt.Sprintf("%s | null", n.Inner.String())
}

func (n *Nullable) Equals(other Type) bool {
	if on, ok := other.(*Nullable); ok {
		return n.Inner.Equals(on.Inner)
	}
	return false
}

// Unwrap returns the inner non-nullable type.
func (n *Nullable) Unwrap() Type {
	return n.Inner
}

// ----------------------------------------------------------------------------
// Class Type
// ----------------------------------------------------------------------------

// Field represents a class field.
type Field struct {
	Name string
	Type Type
}

// Method represents a class method.
type Method struct {
	Name       string
	Params     []*Param
	ReturnType Type
}

// Class represents a class type.
type Class struct {
	Name        string
	Super       *Class
	Fields      map[string]*Field
	Methods     map[string]*Method
	Constructor *Function
}

func (c *Class) typeNode() {}
func (c *Class) String() string {
	return c.Name
}

func (c *Class) Equals(other Type) bool {
	if oc, ok := other.(*Class); ok {
		// Class types are equal if they have the same name (nominal typing)
		return c.Name == oc.Name
	}
	return false
}

// GetField returns the field type, searching up the inheritance chain.
func (c *Class) GetField(name string) *Field {
	if field, ok := c.Fields[name]; ok {
		return field
	}
	if c.Super != nil {
		return c.Super.GetField(name)
	}
	return nil
}

// GetMethod returns the method, searching up the inheritance chain.
func (c *Class) GetMethod(name string) *Method {
	if method, ok := c.Methods[name]; ok {
		return method
	}
	if c.Super != nil {
		return c.Super.GetMethod(name)
	}
	return nil
}

// IsSubclassOf checks if this class is a subclass of another.
func (c *Class) IsSubclassOf(other *Class) bool {
	if c == other {
		return true
	}
	if c.Super != nil {
		return c.Super.IsSubclassOf(other)
	}
	return false
}

// ----------------------------------------------------------------------------
// Type Alias
// ----------------------------------------------------------------------------

// Alias represents a type alias that wraps another type.
type Alias struct {
	Name     string
	Resolved Type
}

func (a *Alias) typeNode() {}
func (a *Alias) String() string {
	return a.Name
}

func (a *Alias) Equals(other Type) bool {
	// Type aliases are structurally typed - compare the resolved types
	return a.Resolved.Equals(other)
}

// Unwrap returns the resolved type.
func (a *Alias) Unwrap() Type {
	// Recursively unwrap nested aliases
	if nested, ok := a.Resolved.(*Alias); ok {
		return nested.Unwrap()
	}
	return a.Resolved
}

// ----------------------------------------------------------------------------
// Helper Functions
// ----------------------------------------------------------------------------

// Unwrap resolves type aliases to their underlying type.
func Unwrap(t Type) Type {
	if alias, ok := t.(*Alias); ok {
		return alias.Unwrap()
	}
	return t
}

// IsNullable checks if a type can be null.
func IsNullable(t Type) bool {
	t = Unwrap(t)
	if _, ok := t.(*Nullable); ok {
		return true
	}
	if p, ok := t.(*Primitive); ok {
		return p.Kind == KindNull
	}
	return false
}

// MakeNullable wraps a type in Nullable if not already nullable.
func MakeNullable(t Type) Type {
	if IsNullable(t) {
		return t
	}
	return &Nullable{Inner: t}
}

// IsAssignableTo checks if a type is assignable to another.
// This handles subtyping, including:
// - Null to nullable types
// - Subclass to superclass
// - Structural compatibility for objects
func IsAssignableTo(from, to Type) bool {
	from = Unwrap(from)
	to = Unwrap(to)

	// Any type is assignable to and from any (for error recovery)
	if _, ok := from.(*Primitive); ok && from.(*Primitive).Kind == KindAny {
		return true
	}
	if _, ok := to.(*Primitive); ok && to.(*Primitive).Kind == KindAny {
		return true
	}

	// Exact equality
	if from.Equals(to) {
		return true
	}

	// Null is assignable to nullable types
	if p, ok := from.(*Primitive); ok && p.Kind == KindNull {
		if _, ok := to.(*Nullable); ok {
			return true
		}
	}

	// Non-null type is assignable to its nullable version
	if nullable, ok := to.(*Nullable); ok {
		if from.Equals(nullable.Inner) {
			return true
		}
	}

	// Class subtyping
	if fromClass, ok := from.(*Class); ok {
		if toClass, ok := to.(*Class); ok {
			return fromClass.IsSubclassOf(toClass)
		}
	}

	// Array covariance
	if fromArr, ok := from.(*Array); ok {
		if toArr, ok := to.(*Array); ok {
			return IsAssignableTo(fromArr.Element, toArr.Element)
		}
	}

	// Object structural subtyping
	if fromObj, ok := from.(*Object); ok {
		if toObj, ok := to.(*Object); ok {
			// from must have all properties that to requires
			for name, toProp := range toObj.Properties {
				fromProp := fromObj.Properties[name]
				if fromProp == nil {
					if !toProp.Optional {
						return false
					}
					continue
				}
				if !IsAssignableTo(fromProp.Type, toProp.Type) {
					return false
				}
			}
			return true
		}
	}

	return false
}

// LeastUpperBound finds the least upper bound of two types.
// Used for inferring the type of conditional expressions.
func LeastUpperBound(a, b Type) Type {
	a = Unwrap(a)
	b = Unwrap(b)

	if a.Equals(b) {
		return a
	}

	// If one is null, result is nullable version of the other
	if p, ok := a.(*Primitive); ok && p.Kind == KindNull {
		return MakeNullable(b)
	}
	if p, ok := b.(*Primitive); ok && p.Kind == KindNull {
		return MakeNullable(a)
	}

	// Both nullable - LUB of inner types, then wrap
	if aNullable, ok := a.(*Nullable); ok {
		if bNullable, ok := b.(*Nullable); ok {
			inner := LeastUpperBound(aNullable.Inner, bNullable.Inner)
			return MakeNullable(inner)
		}
	}

	// Class inheritance - find common ancestor
	if aClass, ok := a.(*Class); ok {
		if bClass, ok := b.(*Class); ok {
			if aClass.IsSubclassOf(bClass) {
				return bClass
			}
			if bClass.IsSubclassOf(aClass) {
				return aClass
			}
			// Find common ancestor by walking up
			for super := aClass.Super; super != nil; super = super.Super {
				if bClass.IsSubclassOf(super) {
					return super
				}
			}
		}
	}

	// No common type found - return any (error recovery)
	return AnyType
}
