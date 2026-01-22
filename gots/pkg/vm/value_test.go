package vm

import "testing"

func TestValueNull(t *testing.T) {
	v := NullValue()

	if !v.IsNull() {
		t.Error("NullValue() should be null")
	}
	if v.IsBool() || v.IsNumber() || v.IsObject() {
		t.Error("NullValue() should not be bool, number, or object")
	}
	if v.String() != "null" {
		t.Errorf("NullValue().String() = %q, want %q", v.String(), "null")
	}
}

func TestValueBool(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		v := BoolValue(tt.input)

		if !v.IsBool() {
			t.Errorf("BoolValue(%v) should be bool", tt.input)
		}
		if v.IsNull() || v.IsNumber() || v.IsObject() {
			t.Error("BoolValue() should not be null, number, or object")
		}
		if v.AsBool() != tt.input {
			t.Errorf("BoolValue(%v).AsBool() = %v", tt.input, v.AsBool())
		}
		if v.String() != tt.expected {
			t.Errorf("BoolValue(%v).String() = %q, want %q", tt.input, v.String(), tt.expected)
		}
	}
}

func TestValueNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{42, "42"},
		{3.14, "3.14"},
		{0, "0"},
		{-7, "-7"},
	}

	for _, tt := range tests {
		v := NumberValue(tt.input)

		if !v.IsNumber() {
			t.Errorf("NumberValue(%v) should be number", tt.input)
		}
		if v.IsNull() || v.IsBool() || v.IsObject() {
			t.Error("NumberValue() should not be null, bool, or object")
		}
		if v.AsNumber() != tt.input {
			t.Errorf("NumberValue(%v).AsNumber() = %v", tt.input, v.AsNumber())
		}
	}
}

func TestValueEquality(t *testing.T) {
	tests := []struct {
		a, b     Value
		expected bool
	}{
		{NullValue(), NullValue(), true},
		{BoolValue(true), BoolValue(true), true},
		{BoolValue(false), BoolValue(false), true},
		{BoolValue(true), BoolValue(false), false},
		{NumberValue(42), NumberValue(42), true},
		{NumberValue(42), NumberValue(43), false},
		{NullValue(), BoolValue(false), false},
		{NullValue(), NumberValue(0), false},
		{BoolValue(false), NumberValue(0), false},
	}

	for i, tt := range tests {
		if ValuesEqual(tt.a, tt.b) != tt.expected {
			t.Errorf("test %d: ValuesEqual(%v, %v) = %v, want %v",
				i, tt.a, tt.b, !tt.expected, tt.expected)
		}
	}
}

func TestObjString(t *testing.T) {
	s := NewObjString("hello")

	if s.Value != "hello" {
		t.Errorf("ObjString.Value = %q, want %q", s.Value, "hello")
	}
	if s.Type() != OBJ_STRING {
		t.Errorf("ObjString.Type() = %v, want %v", s.Type(), OBJ_STRING)
	}

	v := ObjectValue(s)
	if !v.IsObject() {
		t.Error("ObjectValue should be object")
	}
	if !v.IsString() {
		t.Error("ObjectValue(ObjString) should be string")
	}
	if v.AsString() != "hello" {
		t.Errorf("v.AsString() = %q, want %q", v.AsString(), "hello")
	}
}

func TestObjArray(t *testing.T) {
	arr := NewObjArray()
	arr.Elements = append(arr.Elements, NumberValue(1))
	arr.Elements = append(arr.Elements, NumberValue(2))
	arr.Elements = append(arr.Elements, NumberValue(3))

	if arr.Type() != OBJ_ARRAY {
		t.Errorf("ObjArray.Type() = %v, want %v", arr.Type(), OBJ_ARRAY)
	}
	if len(arr.Elements) != 3 {
		t.Errorf("len(arr.Elements) = %d, want 3", len(arr.Elements))
	}

	v := ObjectValue(arr)
	if !v.IsArray() {
		t.Error("ObjectValue(ObjArray) should be array")
	}
}

func TestObjObject(t *testing.T) {
	obj := NewObjObject()
	obj.Fields["x"] = NumberValue(10)
	obj.Fields["y"] = NumberValue(20)

	if obj.Type() != OBJ_OBJECT {
		t.Errorf("ObjObject.Type() = %v, want %v", obj.Type(), OBJ_OBJECT)
	}

	x, ok := obj.Fields["x"]
	if !ok {
		t.Error("obj.Fields[\"x\"] should exist")
	}
	if x.AsNumber() != 10 {
		t.Errorf("obj.Fields[\"x\"] = %v, want 10", x.AsNumber())
	}
}

func TestObjFunction(t *testing.T) {
	fn := &ObjFunction{
		Name:  "add",
		Arity: 2,
	}

	if fn.Type() != OBJ_FUNCTION {
		t.Errorf("ObjFunction.Type() = %v, want %v", fn.Type(), OBJ_FUNCTION)
	}
	if fn.Name != "add" {
		t.Errorf("fn.Name = %q, want %q", fn.Name, "add")
	}
	if fn.Arity != 2 {
		t.Errorf("fn.Arity = %d, want 2", fn.Arity)
	}
}

func TestObjClosure(t *testing.T) {
	fn := &ObjFunction{Name: "test", Arity: 0}
	closure := NewObjClosure(fn)

	if closure.Type() != OBJ_CLOSURE {
		t.Errorf("ObjClosure.Type() = %v, want %v", closure.Type(), OBJ_CLOSURE)
	}
	if closure.Function != fn {
		t.Error("closure.Function should be fn")
	}

	v := ObjectValue(closure)
	if !v.IsClosure() {
		t.Error("ObjectValue(ObjClosure) should be closure")
	}
}

func TestObjClass(t *testing.T) {
	class := NewObjClass("Point")
	class.Methods["distance"] = NewObjClosure(&ObjFunction{Name: "distance", Arity: 1})

	if class.Type() != OBJ_CLASS {
		t.Errorf("ObjClass.Type() = %v, want %v", class.Type(), OBJ_CLASS)
	}
	if class.Name != "Point" {
		t.Errorf("class.Name = %q, want %q", class.Name, "Point")
	}
	if _, ok := class.Methods["distance"]; !ok {
		t.Error("class should have 'distance' method")
	}
}

func TestObjInstance(t *testing.T) {
	class := NewObjClass("Point")
	instance := NewObjInstance(class)
	instance.Fields["x"] = NumberValue(3)
	instance.Fields["y"] = NumberValue(4)

	if instance.Type() != OBJ_INSTANCE {
		t.Errorf("ObjInstance.Type() = %v, want %v", instance.Type(), OBJ_INSTANCE)
	}
	if instance.Class != class {
		t.Error("instance.Class should be class")
	}
	if instance.Fields["x"].AsNumber() != 3 {
		t.Errorf("instance.Fields[\"x\"] = %v, want 3", instance.Fields["x"].AsNumber())
	}
}

func TestValueIsTruthy(t *testing.T) {
	tests := []struct {
		value    Value
		expected bool
	}{
		{NullValue(), false},
		{BoolValue(true), true},
		{BoolValue(false), false},
		{NumberValue(0), true},    // Numbers are truthy (except we might change this)
		{NumberValue(1), true},
		{NumberValue(-1), true},
	}

	for _, tt := range tests {
		if IsTruthy(tt.value) != tt.expected {
			t.Errorf("IsTruthy(%v) = %v, want %v", tt.value, !tt.expected, tt.expected)
		}
	}
}
