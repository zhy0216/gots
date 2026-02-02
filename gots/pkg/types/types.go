// Package types implements the type system for goTS.
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
	KindInt PrimitiveKind = iota
	KindFloat
	KindNumber
	KindString
	KindBoolean
	KindVoid
	KindNull
	KindAny  // Used for error recovery
	KindNever // Bottom type
)

// Primitive represents a primitive type (int, float, string, boolean, void, null).
type Primitive struct {
	Kind PrimitiveKind
}

func (p *Primitive) typeNode() {}
func (p *Primitive) String() string {
	switch p.Kind {
	case KindInt:
		return "int"
	case KindFloat:
		return "float"
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
	IntType     = &Primitive{Kind: KindInt}
	FloatType   = &Primitive{Kind: KindFloat}
	NumberType  = &Primitive{Kind: KindNumber}
	StringType  = &Primitive{Kind: KindString}
	BooleanType = &Primitive{Kind: KindBoolean}
	VoidType    = &Primitive{Kind: KindVoid}
	NullType    = &Primitive{Kind: KindNull}
	AnyType     = &Primitive{Kind: KindAny}
	NeverType   = &Primitive{Kind: KindNever}
)

// ----------------------------------------------------------------------------
// RegExp Type
// ----------------------------------------------------------------------------

// RegExp represents the RegExp built-in type.
type RegExp struct{}

func (r *RegExp) typeNode() {}
func (r *RegExp) String() string {
	return "RegExp"
}

func (r *RegExp) Equals(other Type) bool {
	_, ok := other.(*RegExp)
	return ok
}

// ----------------------------------------------------------------------------
// Console Type
// ----------------------------------------------------------------------------

// Console represents the global console object.
type Console struct{}

func (c *Console) typeNode() {}
func (c *Console) String() string {
	return "Console"
}

func (c *Console) Equals(other Type) bool {
	_, ok := other.(*Console)
	return ok
}

// ----------------------------------------------------------------------------
// Date Type
// ----------------------------------------------------------------------------

// Date represents the built-in Date type.
type Date struct{}

func (d *Date) typeNode() {}
func (d *Date) String() string {
	return "Date"
}

func (d *Date) Equals(other Type) bool {
	_, ok := other.(*Date)
	return ok
}

// DateType is the singleton Date type.
var DateType = &Date{}

// ----------------------------------------------------------------------------
// SQLDatabase Type
// ----------------------------------------------------------------------------

// SQLDatabase represents a database connection (maps to *sql.DB in Go).
type SQLDatabase struct{}

func (s *SQLDatabase) typeNode()          {}
func (s *SQLDatabase) String() string     { return "SQLDatabase" }
func (s *SQLDatabase) Equals(t Type) bool { _, ok := t.(*SQLDatabase); return ok }

// ----------------------------------------------------------------------------
// SQLTransaction Type
// ----------------------------------------------------------------------------

// SQLTransaction represents a database transaction (maps to *sql.Tx in Go).
type SQLTransaction struct{}

func (s *SQLTransaction) typeNode()          {}
func (s *SQLTransaction) String() string     { return "SQLTransaction" }
func (s *SQLTransaction) Equals(t Type) bool { _, ok := t.(*SQLTransaction); return ok }

// Singleton instances for SQL types.
var SQLDatabaseType = &SQLDatabase{}
var SQLTransactionType = &SQLTransaction{}

// ----------------------------------------------------------------------------
// BuiltinObject Type (for Math, JSON, etc.)
// ----------------------------------------------------------------------------

// BuiltinObject represents a built-in global object type (Math, JSON, etc.).
type BuiltinObject struct {
	Name string
}

func (b *BuiltinObject) typeNode() {}
func (b *BuiltinObject) String() string {
	return b.Name
}

func (b *BuiltinObject) Equals(other Type) bool {
	if o, ok := other.(*BuiltinObject); ok {
		return b.Name == o.Name
	}
	return false
}

// RegExpType is the singleton RegExp type.
var RegExpType = &RegExp{}

// ----------------------------------------------------------------------------
// Literal Type
// ----------------------------------------------------------------------------

// Literal represents a literal type (e.g., "hello", 42, true).
// This is a singleton type that only accepts a specific value.
type Literal struct {
	Kind  PrimitiveKind // The base kind (int, float, string, boolean)
	Value string        // The literal value as a string
}

func (l *Literal) typeNode() {}
func (l *Literal) String() string {
	return l.Value
}

func (l *Literal) Equals(other Type) bool {
	if ol, ok := other.(*Literal); ok {
		return l.Kind == ol.Kind && l.Value == ol.Value
	}
	return false
}

// BaseType returns the base primitive type for this literal.
func (l *Literal) BaseType() *Primitive {
	return &Primitive{Kind: l.Kind}
}

// ----------------------------------------------------------------------------
// Tuple Type
// ----------------------------------------------------------------------------

// Tuple represents a tuple type (e.g., [string, int] or [string, ...int[]]).
// Tuples are fixed-length arrays where each position has a specific type.
type Tuple struct {
	Elements []Type // Fixed-position element types
	Rest     Type   // Optional rest element type (the element type, not the array type)
}

func (t *Tuple) typeNode() {}
func (t *Tuple) String() string {
	elements := make([]string, len(t.Elements))
	for i, e := range t.Elements {
		elements[i] = e.String()
	}
	if t.Rest != nil {
		return fmt.Sprintf("[%s, ...%s[]]", strings.Join(elements, ", "), t.Rest.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (t *Tuple) Equals(other Type) bool {
	if ot, ok := other.(*Tuple); ok {
		if len(t.Elements) != len(ot.Elements) {
			return false
		}
		for i, e := range t.Elements {
			if !e.Equals(ot.Elements[i]) {
				return false
			}
		}
		// Check rest element
		if t.Rest == nil && ot.Rest == nil {
			return true
		}
		if t.Rest == nil || ot.Rest == nil {
			return false
		}
		return t.Rest.Equals(ot.Rest)
	}
	return false
}

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
// Map Type
// ----------------------------------------------------------------------------

// Map represents a map type (e.g., Map<string, int>).
type Map struct {
	Key   Type
	Value Type
}

func (m *Map) typeNode() {}
func (m *Map) String() string {
	return fmt.Sprintf("Map<%s, %s>", m.Key.String(), m.Value.String())
}

func (m *Map) Equals(other Type) bool {
	if o, ok := other.(*Map); ok {
		return m.Key.Equals(o.Key) && m.Value.Equals(o.Value)
	}
	return false
}

// ----------------------------------------------------------------------------
// Set Type
// ----------------------------------------------------------------------------

// Set represents a set type (e.g., Set<int>).
type Set struct {
	Element Type
}

func (s *Set) typeNode() {}
func (s *Set) String() string {
	return fmt.Sprintf("Set<%s>", s.Element.String())
}

func (s *Set) Equals(other Type) bool {
	if o, ok := other.(*Set); ok {
		return s.Element.Equals(o.Element)
	}
	return false
}

// ----------------------------------------------------------------------------
// Promise Type
// ----------------------------------------------------------------------------

// Promise represents a Promise<T> type.
type Promise struct {
	Value Type // The resolved value type T
}

func (p *Promise) typeNode() {}
func (p *Promise) String() string {
	return fmt.Sprintf("Promise<%s>", p.Value.String())
}

func (p *Promise) Equals(other Type) bool {
	if op, ok := other.(*Promise); ok {
		return p.Value.Equals(op.Value)
	}
	return false
}

// Unwrap returns the inner value type.
func (p *Promise) Unwrap() Type {
	return p.Value
}

// ----------------------------------------------------------------------------
// Enum Type
// ----------------------------------------------------------------------------

// EnumMember represents a member of an enum type.
type EnumMember struct {
	Name  string
	Value int
}

// Enum represents an enum type (e.g., enum Color { Red, Green, Blue }).
type Enum struct {
	Name    string
	Members []*EnumMember
}

func (e *Enum) typeNode() {}
func (e *Enum) String() string {
	return e.Name
}

func (e *Enum) Equals(other Type) bool {
	if o, ok := other.(*Enum); ok {
		return e.Name == o.Name
	}
	return false
}

// GetMember returns the member with the given name, or nil if not found.
func (e *Enum) GetMember(name string) *EnumMember {
	for _, m := range e.Members {
		if m.Name == name {
			return m
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Interface Type
// ----------------------------------------------------------------------------

// InterfaceField represents a field declaration in an interface.
type InterfaceField struct {
	Name string
	Type Type
}

// InterfaceMethod represents a method signature in an interface.
type InterfaceMethod struct {
	Name       string
	Params     []*Param
	ReturnType Type
}

// Interface represents an interface type.
type Interface struct {
	Name       string
	Fields     map[string]*InterfaceField
	FieldOrder []string // Preserves declaration order (needed for codegen)
	Methods    map[string]*InterfaceMethod
}

func (i *Interface) typeNode() {}
func (i *Interface) String() string {
	return i.Name
}

func (i *Interface) Equals(other Type) bool {
	if o, ok := other.(*Interface); ok {
		if i.Name != o.Name {
			return false
		}
		// Compare fields
		if len(i.Fields) != len(o.Fields) {
			return false
		}
		for name, field := range i.Fields {
			otherField, ok := o.Fields[name]
			if !ok {
				return false
			}
			if !field.Type.Equals(otherField.Type) {
				return false
			}
		}
		return true
	}
	return false
}

// HasMethod checks if the interface has a method with the given name.
func (i *Interface) HasMethod(name string) bool {
	_, ok := i.Methods[name]
	return ok
}

// GetMethod returns the method with the given name, or nil if not found.
func (i *Interface) GetMethod(name string) *InterfaceMethod {
	return i.Methods[name]
}

// HasField checks if the interface has a field with the given name.
func (i *Interface) HasField(name string) bool {
	if i.Fields == nil {
		return false
	}
	_, ok := i.Fields[name]
	return ok
}

// GetField returns the field with the given name, or nil if not found.
func (i *Interface) GetField(name string) *InterfaceField {
	if i.Fields == nil {
		return nil
	}
	return i.Fields[name]
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
// Union Type
// ----------------------------------------------------------------------------

// Union represents a union of multiple types (e.g., string | int | boolean).
type Union struct {
	Types []Type
}

func (u *Union) typeNode() {}
func (u *Union) String() string {
	types := make([]string, len(u.Types))
	for i, t := range u.Types {
		types[i] = t.String()
	}
	return strings.Join(types, " | ")
}

func (u *Union) Equals(other Type) bool {
	if ou, ok := other.(*Union); ok {
		if len(u.Types) != len(ou.Types) {
			return false
		}
		// Check all types in order
		for i, t := range u.Types {
			if !t.Equals(ou.Types[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// Contains checks if the union contains a specific type.
func (u *Union) Contains(t Type) bool {
	for _, ut := range u.Types {
		if ut.Equals(t) {
			return true
		}
	}
	return false
}

// ContainsNull checks if null is one of the union members.
func (u *Union) ContainsNull() bool {
	for _, t := range u.Types {
		if p, ok := t.(*Primitive); ok && p.Kind == KindNull {
			return true
		}
	}
	return false
}

// NonNullTypes returns all types in the union except null.
func (u *Union) NonNullTypes() []Type {
	var result []Type
	for _, t := range u.Types {
		if p, ok := t.(*Primitive); ok && p.Kind == KindNull {
			continue
		}
		result = append(result, t)
	}
	return result
}

// ----------------------------------------------------------------------------
// Intersection Type
// ----------------------------------------------------------------------------

// Intersection represents an intersection of multiple types (e.g., A & B).
// In Go, this is typically implemented as a struct with merged fields for object types.
type Intersection struct {
	Types []Type
}

func (i *Intersection) typeNode() {}
func (i *Intersection) String() string {
	types := make([]string, len(i.Types))
	for idx, t := range i.Types {
		types[idx] = t.String()
	}
	return strings.Join(types, " & ")
}

func (i *Intersection) Equals(other Type) bool {
	if oi, ok := other.(*Intersection); ok {
		if len(i.Types) != len(oi.Types) {
			return false
		}
		// Check all types in order
		for idx, t := range i.Types {
			if !t.Equals(oi.Types[idx]) {
				return false
			}
		}
		return true
	}
	return false
}

// MergeAsObject tries to merge the intersection types as an object type.
// Returns nil if the types cannot be merged (e.g., primitive & primitive).
func (i *Intersection) MergeAsObject() *Object {
	// Collect all properties from object types
	props := make(map[string]*Property)

	for _, t := range i.Types {
		t = Unwrap(t)
		if obj, ok := t.(*Object); ok {
			// Merge object properties
			for name, prop := range obj.Properties {
				if existing, exists := props[name]; exists {
					// If property exists, types must match
					if !existing.Type.Equals(prop.Type) {
						return nil // Conflicting property types
					}
				} else {
					props[name] = prop
				}
			}
		} else {
			// Non-object types in intersection cannot be merged
			return nil
		}
	}

	return &Object{Properties: props}
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
	Name            string
	Super           *Class
	Fields          map[string]*Field
	Methods         map[string]*Method
	Constructor     *Function
	GenericBaseName string // Original name if instantiated from a generic class (e.g., "Box")
	TypeArgs        []Type // Type arguments if instantiated (e.g., [int])
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
	if union, ok := t.(*Union); ok {
		return union.ContainsNull()
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

// MakeUnion creates a union type from multiple types, flattening nested unions.
func MakeUnion(types ...Type) Type {
	if len(types) == 0 {
		return NeverType
	}
	if len(types) == 1 {
		return types[0]
	}

	// Flatten nested unions
	var flattened []Type
	seen := make(map[string]bool)

	for _, t := range types {
		t = Unwrap(t)
		if union, ok := t.(*Union); ok {
			// Flatten nested union
			for _, ut := range union.Types {
				key := ut.String()
				if !seen[key] {
					flattened = append(flattened, ut)
					seen[key] = true
				}
			}
		} else {
			key := t.String()
			if !seen[key] {
				flattened = append(flattened, t)
				seen[key] = true
			}
		}
	}

	if len(flattened) == 1 {
		return flattened[0]
	}

	return &Union{Types: flattened}
}

// MakeIntersection creates an intersection type from multiple types.
// If all types are objects, it merges them into a single object.
func MakeIntersection(types ...Type) Type {
	if len(types) == 0 {
		return NeverType
	}
	if len(types) == 1 {
		return types[0]
	}

	// Try to merge as object
	intersection := &Intersection{Types: types}
	if merged := intersection.MergeAsObject(); merged != nil {
		return merged
	}

	return intersection
}

// IsAssignableTo checks if a type is assignable to another.
// This handles subtyping, including:
// - Null to nullable types
// - Subclass to superclass
// - Structural compatibility for objects
// - Function type compatibility with any
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

	// Literal types are assignable to their base primitive type
	if fromLit, ok := from.(*Literal); ok {
		if toPrim, ok := to.(*Primitive); ok {
			if fromLit.Kind == toPrim.Kind {
				return true
			}
			// Numeric literal types can be assigned to any numeric type (int, float, number)
			fromIsNumeric := fromLit.Kind == KindInt || fromLit.Kind == KindFloat || fromLit.Kind == KindNumber
			toIsNumeric := toPrim.Kind == KindInt || toPrim.Kind == KindFloat || toPrim.Kind == KindNumber
			if fromIsNumeric && toIsNumeric {
				return true
			}
		}
		// Literal to literal - same kind is compatible (for comparisons)
		if toLit, ok := to.(*Literal); ok {
			if fromLit.Kind == toLit.Kind {
				return true
			}
			// Numeric literals are compatible with each other
			fromIsNumeric := fromLit.Kind == KindInt || fromLit.Kind == KindFloat || fromLit.Kind == KindNumber
			toIsNumeric := toLit.Kind == KindInt || toLit.Kind == KindFloat || toLit.Kind == KindNumber
			if fromIsNumeric && toIsNumeric {
				return true
			}
		}
	}

	// Numeric type compatibility:
	// int -> number: allowed (widening)
	// float -> number: allowed (equivalent)
	// number -> float: allowed (equivalent)
	// number -> int: NOT allowed (use toint())
	if fromPrim, ok := from.(*Primitive); ok {
		if toPrim, ok := to.(*Primitive); ok {
			// int or float can be assigned to number
			if toPrim.Kind == KindNumber {
				if fromPrim.Kind == KindInt || fromPrim.Kind == KindFloat {
					return true
				}
			}
			// number can be assigned to float (they're equivalent at runtime)
			if fromPrim.Kind == KindNumber && toPrim.Kind == KindFloat {
				return true
			}
		}
	}

	// Null is assignable to nullable types
	if p, ok := from.(*Primitive); ok && p.Kind == KindNull {
		if _, ok := to.(*Nullable); ok {
			return true
		}
		// Null is also assignable to union types containing null
		if toUnion, ok := to.(*Union); ok {
			return toUnion.ContainsNull()
		}
	}

	// Non-null type is assignable to its nullable version
	if nullable, ok := to.(*Nullable); ok {
		if from.Equals(nullable.Inner) {
			return true
		}
	}

	// Union type assignability
	// A type is assignable to a union if it's assignable to any member of the union
	if toUnion, ok := to.(*Union); ok {
		for _, t := range toUnion.Types {
			if IsAssignableTo(from, t) {
				return true
			}
		}
		return false
	}

	// A union is assignable to another type if all members are assignable to that type
	if fromUnion, ok := from.(*Union); ok {
		for _, t := range fromUnion.Types {
			if !IsAssignableTo(t, to) {
				return false
			}
		}
		return true
	}

	// Intersection type assignability
	// An intersection A & B is assignable to A, to B, or to any supertype
	if fromIntersection, ok := from.(*Intersection); ok {
		// Intersection is assignable to any of its component types
		for _, t := range fromIntersection.Types {
			if IsAssignableTo(t, to) {
				return true
			}
		}
		// Also check if merged form is assignable
		if merged := fromIntersection.MergeAsObject(); merged != nil {
			return IsAssignableTo(merged, to)
		}
		return false
	}

	// A type is assignable to an intersection if it's assignable to all members
	if toIntersection, ok := to.(*Intersection); ok {
		for _, t := range toIntersection.Types {
			if !IsAssignableTo(from, t) {
				return false
			}
		}
		return true
	}

	// Class subtyping
	if fromClass, ok := from.(*Class); ok {
		if toClass, ok := to.(*Class); ok {
			return fromClass.IsSubclassOf(toClass)
		}
		// Class can be assigned to interface if it implements the interface
		if toInterface, ok := to.(*Interface); ok {
			return classImplementsInterface(fromClass, toInterface)
		}
	}

	// Array covariance
	if fromArr, ok := from.(*Array); ok {
		if toArr, ok := to.(*Array); ok {
			return IsAssignableTo(fromArr.Element, toArr.Element)
		}
	}

	// Tuple type assignability
	if fromTuple, ok := from.(*Tuple); ok {
		if toTuple, ok := to.(*Tuple); ok {
			// Tuples must have the same number of elements
			if len(fromTuple.Elements) != len(toTuple.Elements) {
				return false
			}
			// Each element must be assignable
			for i := range fromTuple.Elements {
				if !IsAssignableTo(fromTuple.Elements[i], toTuple.Elements[i]) {
					return false
				}
			}
			// Check rest elements
			if fromTuple.Rest == nil && toTuple.Rest == nil {
				return true
			}
			if fromTuple.Rest == nil || toTuple.Rest == nil {
				return false
			}
			return IsAssignableTo(fromTuple.Rest, toTuple.Rest)
		}
	}

	// Array to tuple assignment (for array literals)
	// An array is assignable to a tuple if the array element type is assignable to all tuple element types
	if fromArr, ok := from.(*Array); ok {
		if toTuple, ok := to.(*Tuple); ok {
			// Array element type must be assignable to each tuple element type
			for _, elemType := range toTuple.Elements {
				if !IsAssignableTo(fromArr.Element, elemType) {
					return false
				}
			}
			return true
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
		// Object can be assigned to interface if it has all required fields
		if toIface, ok := to.(*Interface); ok {
			// Check that the object has all interface fields with compatible types
			for fieldName, ifaceField := range toIface.Fields {
				objProp := fromObj.Properties[fieldName]
				if objProp == nil {
					return false
				}
				if !IsAssignableTo(objProp.Type, ifaceField.Type) {
					return false
				}
			}
			// Check that the object has all interface methods (would need function-typed properties)
			// For now, objects with only fields can satisfy field-only interfaces
			if len(toIface.Methods) > 0 {
				return false
			}
			return true
		}
	}

	// Function type compatibility
	// A function is assignable to another function if:
	// - They have the same number of parameters (or target has any params)
	// - Each parameter is compatible (contravariant, but we allow any)
	// - Return type is compatible (covariant, but we allow any)
	if fromFn, ok := from.(*Function); ok {
		if toFn, ok := to.(*Function); ok {
			// If target function has any-typed single param, accept any function
			if len(toFn.Params) == 1 && toFn.Params[0].Type.Equals(AnyType) && toFn.ReturnType.Equals(AnyType) {
				return true
			}
			// If param counts differ, not assignable (unless target is generic any function)
			if len(fromFn.Params) != len(toFn.Params) {
				return false
			}
			// Check each parameter (contravariant, but we accept any on either side)
			for i := range fromFn.Params {
				fromP := fromFn.Params[i].Type
				toP := toFn.Params[i].Type
				// If either is any, it's compatible
				if fromP.Equals(AnyType) || toP.Equals(AnyType) {
					continue
				}
				// Contravariance: to's param should be assignable to from's param
				if !IsAssignableTo(toP, fromP) {
					return false
				}
			}
			// Check return type (covariant)
			fromRet := fromFn.ReturnType
			toRet := toFn.ReturnType
			if fromRet.Equals(AnyType) || toRet.Equals(AnyType) {
				return true
			}
			return IsAssignableTo(fromRet, toRet)
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

	// One nullable, one not - check if non-nullable matches inner type
	// e.g., LUB(string, string | null) = string | null
	if aNullable, ok := a.(*Nullable); ok {
		if aNullable.Inner.Equals(b) {
			return a // Return the nullable type
		}
	}
	if bNullable, ok := b.(*Nullable); ok {
		if bNullable.Inner.Equals(a) {
			return b // Return the nullable type
		}
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

// IsNumeric checks if a type is int, float, number, or any (for dynamic typing support).
func IsNumeric(t Type) bool {
	t = Unwrap(t)
	if p, ok := t.(*Primitive); ok {
		return p.Kind == KindInt || p.Kind == KindFloat || p.Kind == KindNumber || p.Kind == KindAny
	}
	if l, ok := t.(*Literal); ok {
		return l.Kind == KindInt || l.Kind == KindFloat || l.Kind == KindNumber
	}
	return false
}

// NumericResultType returns the result type for numeric operations.
// Returns any if either operand is any, number if either is number,
// If both are literals, result is a number literal.
// If one is a specific type (int/float) and the other is a literal, use the specific type.
// int if both are int, else float.
func NumericResultType(left, right Type) Type {
	left = Unwrap(left)
	right = Unwrap(right)

	// Check if both operands are literals
	_, leftIsLit := left.(*Literal)
	_, rightIsLit := right.(*Literal)
	if leftIsLit && rightIsLit {
		// Both are literals - result is a number literal
		return &Literal{Kind: KindNumber, Value: ""}
	}

	// Helper to get the primitive kind from a type
	getKind := func(t Type) (PrimitiveKind, bool) {
		if p, ok := t.(*Primitive); ok {
			return p.Kind, true
		}
		if l, ok := t.(*Literal); ok {
			return l.Kind, true
		}
		return 0, false
	}

	leftKind, lok := getKind(left)
	rightKind, rok := getKind(right)
	if lok && rok {
		if leftKind == KindAny || rightKind == KindAny {
			return AnyType
		}

		// If one is a specific type (int/float) and the other is a literal,
		// the result inherits the specific type
		if leftIsLit && !rightIsLit {
			if rightKind == KindInt {
				return IntType
			}
			if rightKind == KindFloat {
				return FloatType
			}
		}
		if rightIsLit && !leftIsLit {
			if leftKind == KindInt {
				return IntType
			}
			if leftKind == KindFloat {
				return FloatType
			}
		}

		// If either is number, result is number
		if leftKind == KindNumber || rightKind == KindNumber {
			return NumberType
		}
		if leftKind == KindInt && rightKind == KindInt {
			return IntType
		}
	}
	return FloatType
}

// classImplementsInterface checks if a class implements an interface (structural typing).
func classImplementsInterface(class *Class, iface *Interface) bool {
	// Check fields
	for fieldName, ifaceField := range iface.Fields {
		classField := class.GetField(fieldName)
		if classField == nil {
			return false
		}
		if !IsAssignableTo(classField.Type, ifaceField.Type) {
			return false
		}
	}

	// Check methods
	for methodName, ifaceMethod := range iface.Methods {
		classMethod := class.GetMethod(methodName)
		if classMethod == nil {
			return false
		}

		// Check parameter count
		if len(classMethod.Params) != len(ifaceMethod.Params) {
			return false
		}

		// Check parameter types
		for i, ifaceParam := range ifaceMethod.Params {
			if !classMethod.Params[i].Type.Equals(ifaceParam.Type) {
				return false
			}
		}

		// Check return type
		if !classMethod.ReturnType.Equals(ifaceMethod.ReturnType) {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------------------
// Type Parameter
// ----------------------------------------------------------------------------

// TypeParameter represents a generic type parameter (e.g., T in function identity<T>).
type TypeParameter struct {
	Name       string
	Constraint Type // Optional constraint (e.g., T extends Comparable)
	Default    Type // Optional default type (e.g., T = string)
}

func (t *TypeParameter) typeNode() {}
func (t *TypeParameter) String() string {
	if t.Constraint != nil {
		return fmt.Sprintf("%s extends %s", t.Name, t.Constraint.String())
	}
	return t.Name
}

func (t *TypeParameter) Equals(other Type) bool {
	if ot, ok := other.(*TypeParameter); ok {
		return t.Name == ot.Name
	}
	return false
}

// SatisfiesConstraint checks if a type satisfies this type parameter's constraint.
func (t *TypeParameter) SatisfiesConstraint(typ Type) bool {
	if t.Constraint == nil {
		return true
	}
	return IsAssignableTo(typ, t.Constraint)
}

// ----------------------------------------------------------------------------
// Generic Function Type
// ----------------------------------------------------------------------------

// GenericFunction represents a generic function type.
type GenericFunction struct {
	TypeParams []*TypeParameter
	Params     []*Param
	ReturnType Type
}

func (g *GenericFunction) typeNode() {}
func (g *GenericFunction) String() string {
	typeParams := make([]string, len(g.TypeParams))
	for i, tp := range g.TypeParams {
		typeParams[i] = tp.String()
	}
	params := make([]string, len(g.Params))
	for i, p := range g.Params {
		params[i] = p.Type.String()
	}
	return fmt.Sprintf("<%s>(%s) => %s", strings.Join(typeParams, ", "), strings.Join(params, ", "), g.ReturnType.String())
}

func (g *GenericFunction) Equals(other Type) bool {
	if og, ok := other.(*GenericFunction); ok {
		if len(g.TypeParams) != len(og.TypeParams) {
			return false
		}
		if len(g.Params) != len(og.Params) {
			return false
		}
		for i, tp := range g.TypeParams {
			if !tp.Equals(og.TypeParams[i]) {
				return false
			}
		}
		for i, p := range g.Params {
			if !p.Type.Equals(og.Params[i].Type) {
				return false
			}
		}
		return g.ReturnType.Equals(og.ReturnType)
	}
	return false
}

// Instantiate creates a concrete function type by substituting type parameters.
func (g *GenericFunction) Instantiate(typeArgs []Type) (*Function, error) {
	// Fill in defaults for missing type args
	if len(typeArgs) < len(g.TypeParams) {
		filled := make([]Type, len(g.TypeParams))
		copy(filled, typeArgs)
		for i := len(typeArgs); i < len(g.TypeParams); i++ {
			if g.TypeParams[i].Default != nil {
				filled[i] = g.TypeParams[i].Default
			} else {
				return nil, fmt.Errorf("expected %d type arguments, got %d", len(g.TypeParams), len(typeArgs))
			}
		}
		typeArgs = filled
	}
	if len(typeArgs) != len(g.TypeParams) {
		return nil, fmt.Errorf("expected %d type arguments, got %d", len(g.TypeParams), len(typeArgs))
	}

	// Build substitution map
	subst := make(map[string]Type)
	for i, tp := range g.TypeParams {
		if !tp.SatisfiesConstraint(typeArgs[i]) {
			return nil, fmt.Errorf("type %s does not satisfy constraint %s", typeArgs[i].String(), tp.Constraint.String())
		}
		subst[tp.Name] = typeArgs[i]
	}

	// Substitute in params
	params := make([]*Param, len(g.Params))
	for i, p := range g.Params {
		params[i] = &Param{
			Name: p.Name,
			Type: substituteType(p.Type, subst),
		}
	}

	return &Function{
		Params:     params,
		ReturnType: substituteType(g.ReturnType, subst),
	}, nil
}

// ----------------------------------------------------------------------------
// Generic Class Type
// ----------------------------------------------------------------------------

// GenericClass represents a generic class type (e.g., Stack<T>).
type GenericClass struct {
	Name        string
	TypeParams  []*TypeParameter
	Super       *Class
	Fields      map[string]*Field
	Methods     map[string]*Method
	Constructor *Function
}

func (g *GenericClass) typeNode() {}
func (g *GenericClass) String() string {
	typeParams := make([]string, len(g.TypeParams))
	for i, tp := range g.TypeParams {
		typeParams[i] = tp.String()
	}
	return fmt.Sprintf("%s<%s>", g.Name, strings.Join(typeParams, ", "))
}

func (g *GenericClass) Equals(other Type) bool {
	if og, ok := other.(*GenericClass); ok {
		return g.Name == og.Name
	}
	return false
}

// Instantiate creates a concrete class type by substituting type parameters.
func (g *GenericClass) Instantiate(typeArgs []Type) (*Class, error) {
	if len(typeArgs) != len(g.TypeParams) {
		return nil, fmt.Errorf("expected %d type arguments, got %d", len(g.TypeParams), len(typeArgs))
	}

	// Build substitution map
	subst := make(map[string]Type)
	for i, tp := range g.TypeParams {
		if !tp.SatisfiesConstraint(typeArgs[i]) {
			return nil, fmt.Errorf("type %s does not satisfy constraint %s", typeArgs[i].String(), tp.Constraint.String())
		}
		subst[tp.Name] = typeArgs[i]
	}

	// Create instantiated class name (e.g., Stack_int)
	argNames := make([]string, len(typeArgs))
	for i, t := range typeArgs {
		argNames[i] = typeNameForInstantiation(t)
	}
	instName := g.Name + "_" + strings.Join(argNames, "_")

	// Substitute in fields
	fields := make(map[string]*Field)
	for name, f := range g.Fields {
		fields[name] = &Field{
			Name: f.Name,
			Type: substituteType(f.Type, subst),
		}
	}

	// Substitute in methods
	methods := make(map[string]*Method)
	for name, m := range g.Methods {
		params := make([]*Param, len(m.Params))
		for i, p := range m.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: substituteType(p.Type, subst),
			}
		}
		methods[name] = &Method{
			Name:       m.Name,
			Params:     params,
			ReturnType: substituteType(m.ReturnType, subst),
		}
	}

	// Substitute in constructor
	var constructor *Function
	if g.Constructor != nil {
		params := make([]*Param, len(g.Constructor.Params))
		for i, p := range g.Constructor.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: substituteType(p.Type, subst),
			}
		}
		constructor = &Function{
			Params:     params,
			ReturnType: g.Constructor.ReturnType,
		}
	}

	return &Class{
		Name:            instName,
		Super:           g.Super,
		Fields:          fields,
		Methods:         methods,
		Constructor:     constructor,
		GenericBaseName: g.Name,
		TypeArgs:        typeArgs,
	}, nil
}

// substituteType replaces type parameters with their concrete types.
func substituteType(t Type, subst map[string]Type) Type {
	switch typ := t.(type) {
	case *TypeParameter:
		if concrete, ok := subst[typ.Name]; ok {
			return concrete
		}
		return t
	case *Array:
		return &Array{Element: substituteType(typ.Element, subst)}
	case *Map:
		return &Map{
			Key:   substituteType(typ.Key, subst),
			Value: substituteType(typ.Value, subst),
		}
	case *Set:
		return &Set{Element: substituteType(typ.Element, subst)}
	case *Promise:
		return &Promise{Value: substituteType(typ.Value, subst)}
	case *Nullable:
		return &Nullable{Inner: substituteType(typ.Inner, subst)}
	case *Function:
		params := make([]*Param, len(typ.Params))
		for i, p := range typ.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: substituteType(p.Type, subst),
			}
		}
		return &Function{
			Params:     params,
			ReturnType: substituteType(typ.ReturnType, subst),
		}
	case *Union:
		types := make([]Type, len(typ.Types))
		for i, t := range typ.Types {
			types[i] = substituteType(t, subst)
		}
		return &Union{Types: types}
	case *Intersection:
		types := make([]Type, len(typ.Types))
		for i, t := range typ.Types {
			types[i] = substituteType(t, subst)
		}
		return &Intersection{Types: types}
	case *Tuple:
		elements := make([]Type, len(typ.Elements))
		for i, e := range typ.Elements {
			elements[i] = substituteType(e, subst)
		}
		var rest Type
		if typ.Rest != nil {
			rest = substituteType(typ.Rest, subst)
		}
		return &Tuple{Elements: elements, Rest: rest}
	case *Object:
		props := make(map[string]*Property)
		for name, prop := range typ.Properties {
			props[name] = &Property{
				Name:     prop.Name,
				Type:     substituteType(prop.Type, subst),
				Optional: prop.Optional,
			}
		}
		return &Object{Properties: props}
	case *Literal:
		return t
	case *Alias:
		return &Alias{
			Name:     typ.Name,
			Resolved: substituteType(typ.Resolved, subst),
		}
	default:
		return t
	}
}

// typeNameForInstantiation returns a simple name for a type suitable for instantiation names.
func typeNameForInstantiation(t Type) string {
	switch typ := t.(type) {
	case *Primitive:
		return typ.String()
	case *Literal:
		// Widen literal to its base type name
		return typ.BaseType().String()
	case *Class:
		return typ.Name
	case *Array:
		return "arr_" + typeNameForInstantiation(typ.Element)
	case *Map:
		return "map_" + typeNameForInstantiation(typ.Key) + "_" + typeNameForInstantiation(typ.Value)
	case *Set:
		return "set_" + typeNameForInstantiation(typ.Element)
	case *Promise:
		return "promise_" + typeNameForInstantiation(typ.Value)
	default:
		return "unknown"
	}
}

// WidenLiteral converts literal types to their base primitive types.
// Non-literal types are returned unchanged.
func WidenLiteral(t Type) Type {
	if lit, ok := t.(*Literal); ok {
		return lit.BaseType()
	}
	return t
}

// ----------------------------------------------------------------------------
// Generic Alias Type
// ----------------------------------------------------------------------------

// GenericAlias represents a generic type alias (e.g., type Pair<T, U> = { first: T, second: U }).
type GenericAlias struct {
	Name       string
	TypeParams []*TypeParameter
	Body       Type // The uninstantiated body type (contains TypeParameter references)
}

func (g *GenericAlias) typeNode() {}
func (g *GenericAlias) String() string {
	typeParams := make([]string, len(g.TypeParams))
	for i, tp := range g.TypeParams {
		typeParams[i] = tp.String()
	}
	return fmt.Sprintf("%s<%s>", g.Name, strings.Join(typeParams, ", "))
}

func (g *GenericAlias) Equals(other Type) bool {
	if og, ok := other.(*GenericAlias); ok {
		return g.Name == og.Name
	}
	return false
}

// Instantiate creates a concrete type by substituting type parameters.
func (g *GenericAlias) Instantiate(typeArgs []Type) (Type, error) {
	if len(typeArgs) != len(g.TypeParams) {
		return nil, fmt.Errorf("expected %d type arguments, got %d", len(g.TypeParams), len(typeArgs))
	}

	subst := make(map[string]Type)
	for i, tp := range g.TypeParams {
		if !tp.SatisfiesConstraint(typeArgs[i]) {
			return nil, fmt.Errorf("type %s does not satisfy constraint %s", typeArgs[i].String(), tp.Constraint.String())
		}
		subst[tp.Name] = typeArgs[i]
	}

	return substituteType(g.Body, subst), nil
}

// ----------------------------------------------------------------------------
// Generic Interface Type
// ----------------------------------------------------------------------------

// GenericInterface represents a generic interface type (e.g., interface Container<T> { get(): T }).
type GenericInterface struct {
	Name       string
	TypeParams []*TypeParameter
	Fields     map[string]*InterfaceField  // Uninstantiated fields (may contain TypeParameter references)
	FieldOrder []string                     // Preserves declaration order
	Methods    map[string]*InterfaceMethod // Uninstantiated methods (contain TypeParameter references)
}

func (g *GenericInterface) typeNode() {}
func (g *GenericInterface) String() string {
	typeParams := make([]string, len(g.TypeParams))
	for i, tp := range g.TypeParams {
		typeParams[i] = tp.String()
	}
	return fmt.Sprintf("%s<%s>", g.Name, strings.Join(typeParams, ", "))
}

func (g *GenericInterface) Equals(other Type) bool {
	if og, ok := other.(*GenericInterface); ok {
		return g.Name == og.Name
	}
	return false
}

// Instantiate creates a concrete interface type by substituting type parameters.
func (g *GenericInterface) Instantiate(typeArgs []Type) (*Interface, error) {
	if len(typeArgs) != len(g.TypeParams) {
		return nil, fmt.Errorf("expected %d type arguments, got %d", len(g.TypeParams), len(typeArgs))
	}

	subst := make(map[string]Type)
	for i, tp := range g.TypeParams {
		if !tp.SatisfiesConstraint(typeArgs[i]) {
			return nil, fmt.Errorf("type %s does not satisfy constraint %s", typeArgs[i].String(), tp.Constraint.String())
		}
		subst[tp.Name] = typeArgs[i]
	}

	// Build instantiated name
	argNames := make([]string, len(typeArgs))
	for i, t := range typeArgs {
		argNames[i] = typeNameForInstantiation(t)
	}
	instName := g.Name + "_" + strings.Join(argNames, "_")

	methods := make(map[string]*InterfaceMethod)
	for name, m := range g.Methods {
		params := make([]*Param, len(m.Params))
		for i, p := range m.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: substituteType(p.Type, subst),
			}
		}
		methods[name] = &InterfaceMethod{
			Name:       m.Name,
			Params:     params,
			ReturnType: substituteType(m.ReturnType, subst),
		}
	}

	// Substitute in fields
	fields := make(map[string]*InterfaceField)
	fieldOrder := make([]string, len(g.FieldOrder))
	copy(fieldOrder, g.FieldOrder)
	for name, f := range g.Fields {
		fields[name] = &InterfaceField{
			Name: f.Name,
			Type: substituteType(f.Type, subst),
		}
	}

	return &Interface{
		Name:       instName,
		Fields:     fields,
		FieldOrder: fieldOrder,
		Methods:    methods,
	}, nil
}
