package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseEnumDeclaration(t *testing.T) {
	input := `enum Color { Red, Green, Blue }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	enumDecl, ok := program.Statements[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("expected *ast.EnumDecl, got %T", program.Statements[0])
	}

	if enumDecl.Name != "Color" {
		t.Errorf("expected name 'Color', got %q", enumDecl.Name)
	}

	if len(enumDecl.Members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(enumDecl.Members))
	}

	expectedMembers := []string{"Red", "Green", "Blue"}
	for i, name := range expectedMembers {
		if enumDecl.Members[i].Name != name {
			t.Errorf("expected member[%d] name %q, got %q", i, name, enumDecl.Members[i].Name)
		}
		if enumDecl.Members[i].Value != nil {
			t.Errorf("expected member[%d] value nil, got %v", i, enumDecl.Members[i].Value)
		}
	}
}

func TestParseEnumWithValues(t *testing.T) {
	input := `enum Status { Pending = 1, Active = 2, Completed = 3 }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	enumDecl, ok := program.Statements[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("expected *ast.EnumDecl, got %T", program.Statements[0])
	}

	if enumDecl.Name != "Status" {
		t.Errorf("expected name 'Status', got %q", enumDecl.Name)
	}

	if len(enumDecl.Members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(enumDecl.Members))
	}

	// Check each member has a value
	for i, m := range enumDecl.Members {
		if m.Value == nil {
			t.Errorf("expected member[%d] to have a value", i)
		}
		numLit, ok := m.Value.(*ast.NumberLiteral)
		if !ok {
			t.Errorf("expected member[%d] value to be NumberLiteral, got %T", i, m.Value)
		}
		if numLit.Value != float64(i+1) {
			t.Errorf("expected member[%d] value %d, got %v", i, i+1, numLit.Value)
		}
	}
}

func TestParseEnumMixed(t *testing.T) {
	input := `enum Priority { Low, Medium = 5, High }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	enumDecl, ok := program.Statements[0].(*ast.EnumDecl)
	if !ok {
		t.Fatalf("expected *ast.EnumDecl, got %T", program.Statements[0])
	}

	if len(enumDecl.Members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(enumDecl.Members))
	}

	// Low has no explicit value
	if enumDecl.Members[0].Value != nil {
		t.Errorf("expected Low to have no value")
	}

	// Medium = 5
	if enumDecl.Members[1].Value == nil {
		t.Errorf("expected Medium to have a value")
	}

	// High has no explicit value
	if enumDecl.Members[2].Value != nil {
		t.Errorf("expected High to have no value")
	}
}
