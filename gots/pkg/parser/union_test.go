package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseUnionTypeInTypeAlias(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedTypes []string
	}{
		{
			name:          "two types",
			input:         `type StringOrNumber = string | int`,
			expectedTypes: []string{"string", "int"},
		},
		{
			name:          "three types",
			input:         `type Value = string | int | boolean`,
			expectedTypes: []string{"string", "int", "boolean"},
		},
		{
			name:          "with null",
			input:         `type OptionalString = string | null`,
			expectedTypes: []string{"string", "null"},
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

			unionType, ok := typeAlias.AliasType.(*ast.UnionType)
			if !ok {
				// Check if it's a NullableType for backward compatibility
				if nullable, ok := typeAlias.AliasType.(*ast.NullableType); ok && len(tt.expectedTypes) == 2 && tt.expectedTypes[1] == "null" {
					if nullable.Inner.String() != tt.expectedTypes[0] {
						t.Errorf("inner type = %s, want %s", nullable.Inner.String(), tt.expectedTypes[0])
					}
					return
				}
				t.Fatalf("aliasType is not *ast.UnionType, got %T", typeAlias.AliasType)
			}

			if len(unionType.Types) != len(tt.expectedTypes) {
				t.Fatalf("union has %d types, want %d", len(unionType.Types), len(tt.expectedTypes))
			}

			for i, expectedType := range tt.expectedTypes {
				if unionType.Types[i].String() != expectedType {
					t.Errorf("type[%d] = %s, want %s", i, unionType.Types[i].String(), expectedType)
				}
			}
		})
	}
}

func TestParseUnionTypeInVarDecl(t *testing.T) {
	input := `let x: string | int = "hello"`

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

	unionType, ok := varDecl.VarType.(*ast.UnionType)
	if !ok {
		t.Fatalf("varType is not *ast.UnionType, got %T", varDecl.VarType)
	}

	if len(unionType.Types) != 2 {
		t.Fatalf("union has %d types, want 2", len(unionType.Types))
	}

	if unionType.Types[0].String() != "string" {
		t.Errorf("type[0] = %s, want string", unionType.Types[0].String())
	}
	if unionType.Types[1].String() != "int" {
		t.Errorf("type[1] = %s, want int", unionType.Types[1].String())
	}
}

func TestParseUnionTypeInFunctionParam(t *testing.T) {
	input := `function foo(x: string | int): void { }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl, ok := program.Statements[0].(*ast.FuncDecl)
	if !ok {
		t.Fatalf("statement is not *ast.FuncDecl, got %T", program.Statements[0])
	}

	if len(funcDecl.Params) != 1 {
		t.Fatalf("function has %d params, want 1", len(funcDecl.Params))
	}

	unionType, ok := funcDecl.Params[0].ParamType.(*ast.UnionType)
	if !ok {
		t.Fatalf("paramType is not *ast.UnionType, got %T", funcDecl.Params[0].ParamType)
	}

	if len(unionType.Types) != 2 {
		t.Fatalf("union has %d types, want 2", len(unionType.Types))
	}
}

func TestParseUnionTypeArrays(t *testing.T) {
	input := `type ArrayUnion = int[] | string[]`

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

	// Check first type is int[]
	arr1, ok := unionType.Types[0].(*ast.ArrayType)
	if !ok {
		t.Fatalf("type[0] is not *ast.ArrayType, got %T", unionType.Types[0])
	}
	if arr1.ElementType.String() != "int" {
		t.Errorf("type[0] element = %s, want int", arr1.ElementType.String())
	}

	// Check second type is string[]
	arr2, ok := unionType.Types[1].(*ast.ArrayType)
	if !ok {
		t.Fatalf("type[1] is not *ast.ArrayType, got %T", unionType.Types[1])
	}
	if arr2.ElementType.String() != "string" {
		t.Errorf("type[1] element = %s, want string", arr2.ElementType.String())
	}
}
