package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseIntersectionTypeInTypeAlias(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTypes []string
	}{
		{
			name:          "two types",
			input:         `type Combined = Point & Named`,
			expectedTypes: []string{"Point", "Named"},
		},
		{
			name:          "three types",
			input:         `type FullType = A & B & C`,
			expectedTypes: []string{"A", "B", "C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has %d statements, want 1", len(program.Statements))
			}

			typeAlias, ok := program.Statements[0].(*ast.TypeAliasDecl)
			if !ok {
				t.Fatalf("statement is not *ast.TypeAliasDecl, got %T", program.Statements[0])
			}

			intersectionType, ok := typeAlias.AliasType.(*ast.IntersectionType)
			if !ok {
				t.Fatalf("aliasType is not *ast.IntersectionType, got %T", typeAlias.AliasType)
			}

			if len(intersectionType.Types) != len(tt.expectedTypes) {
				t.Fatalf("intersection has %d types, want %d", len(intersectionType.Types), len(tt.expectedTypes))
			}

			for i, expectedType := range tt.expectedTypes {
				if intersectionType.Types[i].String() != expectedType {
					t.Errorf("type[%d] = %s, want %s", i, intersectionType.Types[i].String(), expectedType)
				}
			}
		})
	}
}

func TestParseIntersectionTypeInVarDecl(t *testing.T) {
	input := `let x: Foo & Bar = { a: 1, b: 2 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has %d statements, want 1", len(program.Statements))
	}

	varDecl, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("statement is not *ast.VarDecl, got %T", program.Statements[0])
	}

	intersectionType, ok := varDecl.VarType.(*ast.IntersectionType)
	if !ok {
		t.Fatalf("varType is not *ast.IntersectionType, got %T", varDecl.VarType)
	}

	if len(intersectionType.Types) != 2 {
		t.Fatalf("intersection has %d types, want 2", len(intersectionType.Types))
	}

	if intersectionType.Types[0].String() != "Foo" {
		t.Errorf("type[0] = %s, want Foo", intersectionType.Types[0].String())
	}
	if intersectionType.Types[1].String() != "Bar" {
		t.Errorf("type[1] = %s, want Bar", intersectionType.Types[1].String())
	}
}

func TestParseIntersectionTypeWithObjectTypes(t *testing.T) {
	input := `type Combined = { x: int } & { y: int }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	typeAlias := program.Statements[0].(*ast.TypeAliasDecl)
	intersectionType, ok := typeAlias.AliasType.(*ast.IntersectionType)
	if !ok {
		t.Fatalf("aliasType is not *ast.IntersectionType, got %T", typeAlias.AliasType)
	}

	if len(intersectionType.Types) != 2 {
		t.Fatalf("intersection has %d types, want 2", len(intersectionType.Types))
	}

	// Check first type is an object type
	obj1, ok := intersectionType.Types[0].(*ast.ObjectType)
	if !ok {
		t.Fatalf("type[0] is not *ast.ObjectType, got %T", intersectionType.Types[0])
	}
	if len(obj1.Properties) != 1 || obj1.Properties[0].Name != "x" {
		t.Errorf("type[0] should have property x")
	}

	// Check second type is an object type
	obj2, ok := intersectionType.Types[1].(*ast.ObjectType)
	if !ok {
		t.Fatalf("type[1] is not *ast.ObjectType, got %T", intersectionType.Types[1])
	}
	if len(obj2.Properties) != 1 || obj2.Properties[0].Name != "y" {
		t.Errorf("type[1] should have property y")
	}
}

func TestParseMixedUnionAndIntersection(t *testing.T) {
	// Union has lower precedence than intersection
	// So A & B | C should be (A & B) | C
	input := `type Mixed = A & B | C`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	typeAlias := program.Statements[0].(*ast.TypeAliasDecl)
	unionType, ok := typeAlias.AliasType.(*ast.UnionType)
	if !ok {
		t.Fatalf("aliasType is not *ast.UnionType, got %T", typeAlias.AliasType)
	}

	if len(unionType.Types) != 2 {
		t.Fatalf("union has %d types, want 2", len(unionType.Types))
	}

	// First type should be an intersection
	intersectionType, ok := unionType.Types[0].(*ast.IntersectionType)
	if !ok {
		t.Fatalf("union.Types[0] is not *ast.IntersectionType, got %T", unionType.Types[0])
	}

	if len(intersectionType.Types) != 2 {
		t.Fatalf("intersection has %d types, want 2", len(intersectionType.Types))
	}
}
