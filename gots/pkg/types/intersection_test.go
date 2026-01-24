package types

import (
	"testing"
)

func TestIntersectionType(t *testing.T) {
	// Create two object types
	objA := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
	}}
	objB := &Object{Properties: map[string]*Property{
		"y": {Name: "y", Type: IntType},
	}}

	intersection := MakeIntersection(objA, objB)

	// When both are objects, they should be merged
	merged, ok := intersection.(*Object)
	if !ok {
		t.Fatalf("MakeIntersection of two objects should return *Object, got %T", intersection)
	}

	if len(merged.Properties) != 2 {
		t.Fatalf("merged object should have 2 properties, got %d", len(merged.Properties))
	}

	if merged.Properties["x"] == nil || merged.Properties["y"] == nil {
		t.Error("merged object should have both x and y properties")
	}
}

func TestIntersectionTypeSingle(t *testing.T) {
	result := MakeIntersection(StringType)

	if result != StringType {
		t.Error("MakeIntersection with single type should return that type")
	}
}

func TestIntersectionTypeNonObjects(t *testing.T) {
	// Intersection of non-objects should keep as Intersection
	intersection := &Intersection{Types: []Type{StringType, IntType}}

	// This can't be merged
	merged := intersection.MergeAsObject()
	if merged != nil {
		t.Error("intersection of primitives should not be mergeable")
	}
}

func TestIntersectionTypeEquals(t *testing.T) {
	inter1 := &Intersection{Types: []Type{StringType, IntType}}
	inter2 := &Intersection{Types: []Type{StringType, IntType}}
	inter3 := &Intersection{Types: []Type{IntType, StringType}}

	if !inter1.Equals(inter2) {
		t.Error("identical intersections should be equal")
	}

	// Different order means not equal (similar to union)
	if inter1.Equals(inter3) {
		t.Error("intersections with different order should not be equal")
	}
}

func TestIntersectionString(t *testing.T) {
	intersection := &Intersection{Types: []Type{StringType, IntType, BooleanType}}
	expected := "string & int & boolean"
	if intersection.String() != expected {
		t.Errorf("intersection.String() = %q, want %q", intersection.String(), expected)
	}
}

func TestIntersectionAssignableToComponent(t *testing.T) {
	objA := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
	}}
	objB := &Object{Properties: map[string]*Property{
		"y": {Name: "y", Type: IntType},
	}}

	intersection := &Intersection{Types: []Type{objA, objB}}

	// An intersection should be assignable to any of its components
	if !IsAssignableTo(intersection, objA) {
		t.Error("intersection should be assignable to objA")
	}
	if !IsAssignableTo(intersection, objB) {
		t.Error("intersection should be assignable to objB")
	}
}

func TestAssignableToIntersection(t *testing.T) {
	objA := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
	}}
	objB := &Object{Properties: map[string]*Property{
		"y": {Name: "y", Type: IntType},
	}}
	objAB := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
		"y": {Name: "y", Type: IntType},
	}}

	intersection := &Intersection{Types: []Type{objA, objB}}

	// An object with all required properties should be assignable to the intersection
	if !IsAssignableTo(objAB, intersection) {
		t.Error("objAB should be assignable to intersection of objA and objB")
	}

	// An object with only some properties should not be assignable
	if IsAssignableTo(objA, intersection) {
		t.Error("objA alone should not be assignable to intersection")
	}
}

func TestMergeObjectIntersection(t *testing.T) {
	objA := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
	}}
	objB := &Object{Properties: map[string]*Property{
		"y": {Name: "y", Type: StringType},
	}}

	intersection := &Intersection{Types: []Type{objA, objB}}
	merged := intersection.MergeAsObject()

	if merged == nil {
		t.Fatal("MergeAsObject should succeed for compatible object types")
	}

	if len(merged.Properties) != 2 {
		t.Fatalf("merged should have 2 properties, got %d", len(merged.Properties))
	}

	if !merged.Properties["x"].Type.Equals(IntType) {
		t.Error("merged.x should be int")
	}
	if !merged.Properties["y"].Type.Equals(StringType) {
		t.Error("merged.y should be string")
	}
}

func TestMergeConflictingObjectIntersection(t *testing.T) {
	objA := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: IntType},
	}}
	objB := &Object{Properties: map[string]*Property{
		"x": {Name: "x", Type: StringType}, // Conflict: same property, different type
	}}

	intersection := &Intersection{Types: []Type{objA, objB}}
	merged := intersection.MergeAsObject()

	if merged != nil {
		t.Error("MergeAsObject should fail for conflicting property types")
	}
}
