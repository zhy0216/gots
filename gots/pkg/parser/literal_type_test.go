package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseLiteralTypeInTypeAlias(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "string literal type",
			input: `type Direction = "north" | "south" | "east" | "west"`,
		},
		{
			name:  "number literal type",
			input: `type One = 1`,
		},
		{
			name:  "mixed literal union",
			input: `type Status = "active" | "inactive" | 0 | 1`,
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

			// The type should be parsed (either LiteralType or UnionType containing LiteralTypes)
			if typeAlias.AliasType == nil {
				t.Fatal("aliasType should not be nil")
			}
		})
	}
}

func TestParseLiteralTypeAlone(t *testing.T) {
	input := `let x: "hello" = "hello"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	varDecl := program.Statements[0].(*ast.VarDecl)

	literalType, ok := varDecl.VarType.(*ast.LiteralType)
	if !ok {
		t.Fatalf("varType is not *ast.LiteralType, got %T", varDecl.VarType)
	}

	if literalType.Value != "hello" {
		t.Errorf("literal value = %q, want %q", literalType.Value, "hello")
	}
}

func TestParseNumberLiteralType(t *testing.T) {
	input := `let x: 42 = 42`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	varDecl := program.Statements[0].(*ast.VarDecl)

	literalType, ok := varDecl.VarType.(*ast.LiteralType)
	if !ok {
		t.Fatalf("varType is not *ast.LiteralType, got %T", varDecl.VarType)
	}

	if literalType.Value != "42" {
		t.Errorf("literal value = %q, want %q", literalType.Value, "42")
	}
}

func TestParseBooleanLiteralType(t *testing.T) {
	input := `let x: true = true`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	varDecl := program.Statements[0].(*ast.VarDecl)

	literalType, ok := varDecl.VarType.(*ast.LiteralType)
	if !ok {
		t.Fatalf("varType is not *ast.LiteralType, got %T", varDecl.VarType)
	}

	if literalType.Value != "true" {
		t.Errorf("literal value = %q, want %q", literalType.Value, "true")
	}
}
