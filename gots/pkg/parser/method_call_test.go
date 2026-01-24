package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseStringMethodCall(t *testing.T) {
	input := `
let str: string = "hello"
let result: string[] = str.split(",")
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	// Second statement should be variable declaration with method call
	stmt, ok := program.Statements[1].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[1])
	}

	callExpr, ok := stmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", stmt.Value)
	}

	propExpr, ok := callExpr.Function.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expected PropertyExpr as function, got %T", callExpr.Function)
	}

	if propExpr.Property != "split" {
		t.Errorf("expected method name 'split', got %s", propExpr.Property)
	}

	objIdent, ok := propExpr.Object.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected Identifier as object, got %T", propExpr.Object)
	}

	if objIdent.Name != "str" {
		t.Errorf("expected object name 'str', got %s", objIdent.Name)
	}

	if len(callExpr.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(callExpr.Arguments))
	}
}

func TestParseArrayMethodCall(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3]
let doubled: int[] = arr.map((x: int): int => x * 2)
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[1].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[1])
	}

	callExpr, ok := stmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", stmt.Value)
	}

	propExpr, ok := callExpr.Function.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expected PropertyExpr as function, got %T", callExpr.Function)
	}

	if propExpr.Property != "map" {
		t.Errorf("expected method name 'map', got %s", propExpr.Property)
	}
}

func TestParseChainedMethodCalls(t *testing.T) {
	input := `
let result: string = "  hello  ".trim().split("e").join("-")
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[0])
	}

	// Should be a chain: join is the outermost call
	callExpr, ok := stmt.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr, got %T", stmt.Value)
	}

	propExpr, ok := callExpr.Function.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expected PropertyExpr, got %T", callExpr.Function)
	}

	if propExpr.Property != "join" {
		t.Errorf("expected outermost method 'join', got %s", propExpr.Property)
	}
}
