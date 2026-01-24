package types

import (
	"testing"
)

func TestUnionType(t *testing.T) {
	// Test Union creation
	union := MakeUnion(StringType, IntType)
	if u, ok := union.(*Union); !ok {
		t.Fatalf("MakeUnion did not return *Union, got %T", union)
	} else {
		if len(u.Types) != 2 {
			t.Fatalf("union has %d types, want 2", len(u.Types))
		}
	}
}

func TestUnionTypeEquals(t *testing.T) {
	union1 := &Union{Types: []Type{StringType, IntType}}
	union2 := &Union{Types: []Type{StringType, IntType}}
	union3 := &Union{Types: []Type{IntType, StringType}} // Different order

	if !union1.Equals(union2) {
		t.Error("identical unions should be equal")
	}

	// Note: With our current implementation, order matters for equality
	// This is a design choice - could be changed to set-based equality
	if union1.Equals(union3) {
		t.Error("unions with different order should not be equal (current implementation)")
	}
}

func TestUnionContains(t *testing.T) {
	union := &Union{Types: []Type{StringType, IntType, NullType}}

	if !union.Contains(StringType) {
		t.Error("union should contain StringType")
	}
	if !union.Contains(IntType) {
		t.Error("union should contain IntType")
	}
	if !union.Contains(NullType) {
		t.Error("union should contain NullType")
	}
	if union.Contains(BooleanType) {
		t.Error("union should not contain BooleanType")
	}
}

func TestUnionContainsNull(t *testing.T) {
	unionWithNull := &Union{Types: []Type{StringType, NullType}}
	unionWithoutNull := &Union{Types: []Type{StringType, IntType}}

	if !unionWithNull.ContainsNull() {
		t.Error("union with null should return true for ContainsNull")
	}
	if unionWithoutNull.ContainsNull() {
		t.Error("union without null should return false for ContainsNull")
	}
}

func TestUnionNonNullTypes(t *testing.T) {
	union := &Union{Types: []Type{StringType, IntType, NullType}}
	nonNull := union.NonNullTypes()

	if len(nonNull) != 2 {
		t.Fatalf("expected 2 non-null types, got %d", len(nonNull))
	}
	if !nonNull[0].Equals(StringType) {
		t.Error("first non-null type should be string")
	}
	if !nonNull[1].Equals(IntType) {
		t.Error("second non-null type should be int")
	}
}

func TestMakeUnionFlattensNested(t *testing.T) {
	inner := &Union{Types: []Type{IntType, FloatType}}
	outer := MakeUnion(StringType, inner)

	union, ok := outer.(*Union)
	if !ok {
		t.Fatalf("MakeUnion should return *Union, got %T", outer)
	}

	if len(union.Types) != 3 {
		t.Fatalf("flattened union should have 3 types, got %d", len(union.Types))
	}
}

func TestMakeUnionDeduplicates(t *testing.T) {
	result := MakeUnion(StringType, StringType, IntType)

	union, ok := result.(*Union)
	if !ok {
		t.Fatalf("MakeUnion should return *Union, got %T", result)
	}

	if len(union.Types) != 2 {
		t.Fatalf("union should have 2 unique types, got %d", len(union.Types))
	}
}

func TestMakeUnionSingleType(t *testing.T) {
	result := MakeUnion(StringType)

	if result != StringType {
		t.Error("MakeUnion with single type should return that type")
	}
}

func TestIsAssignableToUnion(t *testing.T) {
	union := &Union{Types: []Type{StringType, IntType}}

	// String is assignable to string | int
	if !IsAssignableTo(StringType, union) {
		t.Error("StringType should be assignable to string | int")
	}

	// Int is assignable to string | int
	if !IsAssignableTo(IntType, union) {
		t.Error("IntType should be assignable to string | int")
	}

	// Boolean is not assignable to string | int
	if IsAssignableTo(BooleanType, union) {
		t.Error("BooleanType should not be assignable to string | int")
	}
}

func TestUnionIsAssignableToType(t *testing.T) {
	// string | int is not assignable to just string
	union := &Union{Types: []Type{StringType, IntType}}

	if IsAssignableTo(union, StringType) {
		t.Error("string | int should not be assignable to string")
	}

	// string | int is assignable to any
	if !IsAssignableTo(union, AnyType) {
		t.Error("string | int should be assignable to any")
	}
}

func TestNullAssignableToUnionWithNull(t *testing.T) {
	union := &Union{Types: []Type{StringType, NullType}}

	if !IsAssignableTo(NullType, union) {
		t.Error("null should be assignable to string | null union")
	}
}

func TestIsNullableWithUnion(t *testing.T) {
	unionWithNull := &Union{Types: []Type{StringType, NullType}}
	unionWithoutNull := &Union{Types: []Type{StringType, IntType}}

	if !IsNullable(unionWithNull) {
		t.Error("union containing null should be nullable")
	}
	if IsNullable(unionWithoutNull) {
		t.Error("union not containing null should not be nullable")
	}
}

func TestUnionString(t *testing.T) {
	union := &Union{Types: []Type{StringType, IntType, BooleanType}}
	expected := "string | int | boolean"
	if union.String() != expected {
		t.Errorf("union.String() = %q, want %q", union.String(), expected)
	}
}
