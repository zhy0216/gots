// Package vm implements the virtual machine for GoTS.
package vm

import (
	"fmt"
	"math"

	"github.com/zhy0216/quickts/gots/pkg/bytecode"
)

// ValueType represents the type of a Value.
type ValueType byte

const (
	VAL_NULL ValueType = iota
	VAL_BOOL
	VAL_NUMBER
	VAL_OBJECT
)

// Value represents a runtime value in the VM.
type Value struct {
	Type ValueType
	data uint64
	obj  Object
}

// NullValue creates a null value.
func NullValue() Value {
	return Value{Type: VAL_NULL}
}

// BoolValue creates a boolean value.
func BoolValue(b bool) Value {
	var data uint64
	if b {
		data = 1
	}
	return Value{Type: VAL_BOOL, data: data}
}

// NumberValue creates a number value.
func NumberValue(n float64) Value {
	return Value{Type: VAL_NUMBER, data: math.Float64bits(n)}
}

// ObjectValue creates an object value.
func ObjectValue(obj Object) Value {
	return Value{Type: VAL_OBJECT, obj: obj}
}

// IsNull returns true if the value is null.
func (v Value) IsNull() bool {
	return v.Type == VAL_NULL
}

// IsBool returns true if the value is a boolean.
func (v Value) IsBool() bool {
	return v.Type == VAL_BOOL
}

// IsNumber returns true if the value is a number.
func (v Value) IsNumber() bool {
	return v.Type == VAL_NUMBER
}

// IsObject returns true if the value is an object.
func (v Value) IsObject() bool {
	return v.Type == VAL_OBJECT
}

// AsBool returns the boolean value.
func (v Value) AsBool() bool {
	return v.data == 1
}

// AsNumber returns the number value.
func (v Value) AsNumber() float64 {
	return math.Float64frombits(v.data)
}

// AsObject returns the object value.
func (v Value) AsObject() Object {
	return v.obj
}

// IsString returns true if the value is a string object.
func (v Value) IsString() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_STRING
}

// AsString returns the string value.
func (v Value) AsString() string {
	if s, ok := v.obj.(*ObjString); ok {
		return s.Value
	}
	return ""
}

// IsArray returns true if the value is an array object.
func (v Value) IsArray() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_ARRAY
}

// AsArray returns the array object.
func (v Value) AsArray() *ObjArray {
	if arr, ok := v.obj.(*ObjArray); ok {
		return arr
	}
	return nil
}

// IsClosure returns true if the value is a closure object.
func (v Value) IsClosure() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_CLOSURE
}

// AsClosure returns the closure object.
func (v Value) AsClosure() *ObjClosure {
	if c, ok := v.obj.(*ObjClosure); ok {
		return c
	}
	return nil
}

// IsClass returns true if the value is a class object.
func (v Value) IsClass() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_CLASS
}

// AsClass returns the class object.
func (v Value) AsClass() *ObjClass {
	if c, ok := v.obj.(*ObjClass); ok {
		return c
	}
	return nil
}

// IsInstance returns true if the value is an instance object.
func (v Value) IsInstance() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_INSTANCE
}

// AsInstance returns the instance object.
func (v Value) AsInstance() *ObjInstance {
	if i, ok := v.obj.(*ObjInstance); ok {
		return i
	}
	return nil
}

// IsBoundMethod returns true if the value is a bound method object.
func (v Value) IsBoundMethod() bool {
	if !v.IsObject() {
		return false
	}
	return v.obj.Type() == OBJ_BOUND_METHOD
}

// AsBoundMethod returns the bound method object.
func (v Value) AsBoundMethod() *ObjBoundMethod {
	if b, ok := v.obj.(*ObjBoundMethod); ok {
		return b
	}
	return nil
}

// String returns a string representation of the value.
func (v Value) String() string {
	switch v.Type {
	case VAL_NULL:
		return "null"
	case VAL_BOOL:
		if v.AsBool() {
			return "true"
		}
		return "false"
	case VAL_NUMBER:
		n := v.AsNumber()
		if n == math.Trunc(n) {
			return fmt.Sprintf("%.0f", n)
		}
		return fmt.Sprintf("%g", n)
	case VAL_OBJECT:
		return v.obj.String()
	}
	return "unknown"
}

// ValuesEqual returns true if two values are equal.
func ValuesEqual(a, b Value) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case VAL_NULL:
		return true
	case VAL_BOOL:
		return a.AsBool() == b.AsBool()
	case VAL_NUMBER:
		return a.AsNumber() == b.AsNumber()
	case VAL_OBJECT:
		return objectsEqual(a.obj, b.obj)
	}

	return false
}

func objectsEqual(a, b Object) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case OBJ_STRING:
		return a.(*ObjString).Value == b.(*ObjString).Value
	default:
		// Reference equality for other objects
		return a == b
	}
}

// IsTruthy returns true if the value is truthy.
func IsTruthy(v Value) bool {
	switch v.Type {
	case VAL_NULL:
		return false
	case VAL_BOOL:
		return v.AsBool()
	default:
		return true
	}
}

// ObjectType represents the type of an object.
type ObjectType byte

const (
	OBJ_STRING ObjectType = iota
	OBJ_ARRAY
	OBJ_OBJECT
	OBJ_FUNCTION
	OBJ_CLOSURE
	OBJ_UPVALUE
	OBJ_CLASS
	OBJ_INSTANCE
	OBJ_BOUND_METHOD
)

// Object is the interface for heap-allocated objects.
type Object interface {
	Type() ObjectType
	String() string
}

// ObjString represents a string object.
type ObjString struct {
	Value  string
	Hash   uint32
	marked bool
}

func (s *ObjString) Type() ObjectType   { return OBJ_STRING }
func (s *ObjString) String() string     { return s.Value }
func (s *ObjString) IsMarked() bool     { return s.marked }
func (s *ObjString) SetMarked(m bool)   { s.marked = m }

// NewObjString creates a new string object.
func NewObjString(value string) *ObjString {
	return &ObjString{
		Value: value,
		Hash:  hashString(value),
	}
}

func hashString(s string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= 16777619
	}
	return hash
}

// ObjArray represents an array object.
type ObjArray struct {
	Elements []Value
	marked   bool
}

func (a *ObjArray) Type() ObjectType { return OBJ_ARRAY }
func (a *ObjArray) String() string {
	return fmt.Sprintf("[Array(%d)]", len(a.Elements))
}
func (a *ObjArray) IsMarked() bool   { return a.marked }
func (a *ObjArray) SetMarked(m bool) { a.marked = m }

// NewObjArray creates a new array object.
func NewObjArray() *ObjArray {
	return &ObjArray{
		Elements: make([]Value, 0),
	}
}

// ObjObject represents a plain object (object literal).
type ObjObject struct {
	Fields map[string]Value
	marked bool
}

func (o *ObjObject) Type() ObjectType { return OBJ_OBJECT }
func (o *ObjObject) String() string   { return "[Object]" }
func (o *ObjObject) IsMarked() bool   { return o.marked }
func (o *ObjObject) SetMarked(m bool) { o.marked = m }

// NewObjObject creates a new plain object.
func NewObjObject() *ObjObject {
	return &ObjObject{
		Fields: make(map[string]Value),
	}
}

// ObjFunction represents a compiled function.
type ObjFunction struct {
	Name         string
	Arity        int
	UpvalueCount int
	Chunk        *bytecode.Chunk
	marked       bool
}

func (f *ObjFunction) Type() ObjectType { return OBJ_FUNCTION }
func (f *ObjFunction) String() string {
	if f.Name == "" {
		return "<script>"
	}
	return fmt.Sprintf("<fn %s>", f.Name)
}
func (f *ObjFunction) IsMarked() bool   { return f.marked }
func (f *ObjFunction) SetMarked(m bool) { f.marked = m }

// ObjClosure represents a closure (function + captured environment).
type ObjClosure struct {
	Function *ObjFunction
	Upvalues []*ObjUpvalue
	marked   bool
}

func (c *ObjClosure) Type() ObjectType { return OBJ_CLOSURE }
func (c *ObjClosure) String() string   { return c.Function.String() }
func (c *ObjClosure) IsMarked() bool   { return c.marked }
func (c *ObjClosure) SetMarked(m bool) { c.marked = m }

// NewObjClosure creates a new closure.
func NewObjClosure(fn *ObjFunction) *ObjClosure {
	upvalues := make([]*ObjUpvalue, fn.UpvalueCount)
	return &ObjClosure{
		Function: fn,
		Upvalues: upvalues,
	}
}

// ObjUpvalue represents a captured variable.
type ObjUpvalue struct {
	Location   *Value
	Closed     Value
	Next       *ObjUpvalue
	stackIndex int
	marked     bool
}

func (u *ObjUpvalue) Type() ObjectType { return OBJ_UPVALUE }
func (u *ObjUpvalue) String() string   { return "upvalue" }
func (u *ObjUpvalue) IsMarked() bool   { return u.marked }
func (u *ObjUpvalue) SetMarked(m bool) { u.marked = m }

// NewObjUpvalue creates a new upvalue pointing to a stack location.
func NewObjUpvalue(slot *Value) *ObjUpvalue {
	return &ObjUpvalue{
		Location: slot,
	}
}

// ObjClass represents a class definition.
type ObjClass struct {
	Name    string
	Super   *ObjClass
	Methods map[string]*ObjClosure
	marked  bool
}

func (c *ObjClass) Type() ObjectType { return OBJ_CLASS }
func (c *ObjClass) String() string   { return c.Name }
func (c *ObjClass) IsMarked() bool   { return c.marked }
func (c *ObjClass) SetMarked(m bool) { c.marked = m }

// NewObjClass creates a new class.
func NewObjClass(name string) *ObjClass {
	return &ObjClass{
		Name:    name,
		Methods: make(map[string]*ObjClosure),
	}
}

// ObjInstance represents an instance of a class.
type ObjInstance struct {
	Class  *ObjClass
	Fields map[string]Value
	marked bool
}

func (i *ObjInstance) Type() ObjectType { return OBJ_INSTANCE }
func (i *ObjInstance) String() string   { return fmt.Sprintf("%s instance", i.Class.Name) }
func (i *ObjInstance) IsMarked() bool   { return i.marked }
func (i *ObjInstance) SetMarked(m bool) { i.marked = m }

// NewObjInstance creates a new instance of a class.
func NewObjInstance(class *ObjClass) *ObjInstance {
	return &ObjInstance{
		Class:  class,
		Fields: make(map[string]Value),
	}
}

// ObjBoundMethod represents a method bound to an instance.
type ObjBoundMethod struct {
	Receiver Value
	Method   *ObjClosure
	marked   bool
}

func (b *ObjBoundMethod) Type() ObjectType { return OBJ_BOUND_METHOD }
func (b *ObjBoundMethod) String() string   { return b.Method.String() }
func (b *ObjBoundMethod) IsMarked() bool   { return b.marked }
func (b *ObjBoundMethod) SetMarked(m bool) { b.marked = m }
