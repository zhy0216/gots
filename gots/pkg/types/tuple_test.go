package types

import (
	"testing"
)

func TestTupleType(t *testing.T) {
	// Test basic tuple creation
	tuple := &Tuple{
		Elements: []Type{StringType, IntType},
	}

	if len(tuple.Elements) != 2 {
		t.Errorf("tuple has %d elements, want 2", len(tuple.Elements))
	}

	// Test String() output
	expected := "[string, int]"
	if tuple.String() != expected {
		t.Errorf("tuple.String() = %q, want %q", tuple.String(), expected)
	}
}

func TestTupleTypeWithRest(t *testing.T) {
	// Test tuple with rest element: [string, ...int[]]
	tuple := &Tuple{
		Elements: []Type{StringType},
		Rest:     IntType, // The element type of the rest array
	}

	expected := "[string, ...int[]]"
	if tuple.String() != expected {
		t.Errorf("tuple.String() = %q, want %q", tuple.String(), expected)
	}
}

func TestTupleTypeEquals(t *testing.T) {
	tuple1 := &Tuple{Elements: []Type{StringType, IntType}}
	tuple2 := &Tuple{Elements: []Type{StringType, IntType}}
	tuple3 := &Tuple{Elements: []Type{IntType, StringType}}
	tuple4 := &Tuple{Elements: []Type{StringType}}

	if !tuple1.Equals(tuple2) {
		t.Error("identical tuples should be equal")
	}

	if tuple1.Equals(tuple3) {
		t.Error("tuples with different element order should not be equal")
	}

	if tuple1.Equals(tuple4) {
		t.Error("tuples with different lengths should not be equal")
	}
}

func TestTupleTypeEqualsWithRest(t *testing.T) {
	tuple1 := &Tuple{Elements: []Type{StringType}, Rest: IntType}
	tuple2 := &Tuple{Elements: []Type{StringType}, Rest: IntType}
	tuple3 := &Tuple{Elements: []Type{StringType}, Rest: StringType}
	tuple4 := &Tuple{Elements: []Type{StringType}} // No rest

	if !tuple1.Equals(tuple2) {
		t.Error("identical tuples with rest should be equal")
	}

	if tuple1.Equals(tuple3) {
		t.Error("tuples with different rest types should not be equal")
	}

	if tuple1.Equals(tuple4) {
		t.Error("tuple with rest should not equal tuple without rest")
	}
}

func TestTupleAssignableToTuple(t *testing.T) {
	tuple1 := &Tuple{Elements: []Type{StringType, IntType}}
	tuple2 := &Tuple{Elements: []Type{StringType, IntType}}

	if !IsAssignableTo(tuple1, tuple2) {
		t.Error("identical tuples should be assignable")
	}
}

func TestTupleIsNotArray(t *testing.T) {
	tuple := &Tuple{Elements: []Type{StringType, IntType}}
	array := &Array{Element: StringType}

	if tuple.Equals(array) {
		t.Error("tuple should not equal array")
	}
}
