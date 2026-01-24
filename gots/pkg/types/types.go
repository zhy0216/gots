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
	KindInt PrimitiveKind = iota
	KindFloat
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
// Interface Type
// ----------------------------------------------------------------------------

// InterfaceMethod represents a method signature in an interface.
type InterfaceMethod struct {
	Name       string
	Params     []*Param
	ReturnType Type
}

// Interface represents an interface type.
type Interface struct {
	Name    string
	Methods map[string]*InterfaceMethod
}

func (i *Interface) typeNode() {}
func (i *Interface) String() string {
	return i.Name
}

func (i *Interface) Equals(other Type) bool {
	if o, ok := other.(*Interface); ok {
		return i.Name == o.Name
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

// IsNumeric checks if a type is int, float, or any (for dynamic typing support).
func IsNumeric(t Type) bool {
	t = Unwrap(t)
	if p, ok := t.(*Primitive); ok {
		return p.Kind == KindInt || p.Kind == KindFloat || p.Kind == KindAny
	}
	return false
}

// NumericResultType returns the result type for numeric operations.
// Returns any if either operand is any, int if both are int, else float.
func NumericResultType(left, right Type) Type {
	left = Unwrap(left)
	right = Unwrap(right)
	lp, lok := left.(*Primitive)
	rp, rok := right.(*Primitive)
	if lok && rok {
		if lp.Kind == KindAny || rp.Kind == KindAny {
			return AnyType
		}
		if lp.Kind == KindInt && rp.Kind == KindInt {
			return IntType
		}
	}
	return FloatType
}

// classImplementsInterface checks if a class implements an interface (structural typing).
func classImplementsInterface(class *Class, iface *Interface) bool {
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
	default:
		return t
	}
}

// typeNameForInstantiation returns a simple name for a type suitable for instantiation names.
func typeNameForInstantiation(t Type) string {
	switch typ := t.(type) {
	case *Primitive:
		return typ.String()
	case *Class:
		return typ.Name
	case *Array:
		return "arr_" + typeNameForInstantiation(typ.Element)
	case *Map:
		return "map_" + typeNameForInstantiation(typ.Key) + "_" + typeNameForInstantiation(typ.Value)
	case *Set:
		return "set_" + typeNameForInstantiation(typ.Element)
	default:
		return "unknown"
	}
}
