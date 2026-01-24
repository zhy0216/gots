package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

// Test re-exports: export { foo } from "./module"
func TestParseReExportNamed(t *testing.T) {
	input := `export { foo, bar } from "./module"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	reExport, ok := program.Statements[0].(*ast.ReExportDecl)
	if !ok {
		t.Fatalf("expected *ast.ReExportDecl, got %T", program.Statements[0])
	}

	if len(reExport.Names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(reExport.Names))
	}

	if reExport.Names[0] != "foo" {
		t.Errorf("expected name 'foo', got %q", reExport.Names[0])
	}

	if reExport.Names[1] != "bar" {
		t.Errorf("expected name 'bar', got %q", reExport.Names[1])
	}

	if reExport.Path != "./module" {
		t.Errorf("expected path './module', got %q", reExport.Path)
	}

	if reExport.IsWildcard {
		t.Error("expected IsWildcard to be false")
	}
}

// Test re-export all: export * from "./module"
func TestParseReExportAll(t *testing.T) {
	input := `export * from "./module"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	reExport, ok := program.Statements[0].(*ast.ReExportDecl)
	if !ok {
		t.Fatalf("expected *ast.ReExportDecl, got %T", program.Statements[0])
	}

	if !reExport.IsWildcard {
		t.Error("expected IsWildcard to be true")
	}

	if reExport.Path != "./module" {
		t.Errorf("expected path './module', got %q", reExport.Path)
	}
}

// Test default export: export default class Foo {}
func TestParseDefaultExportClass(t *testing.T) {
	input := `export default class Foo {
		x: int
		constructor(x: int) {
			this.x = x
		}
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	defaultExport, ok := program.Statements[0].(*ast.DefaultExport)
	if !ok {
		t.Fatalf("expected *ast.DefaultExport, got %T", program.Statements[0])
	}

	classDecl, ok := defaultExport.Decl.(*ast.ClassDecl)
	if !ok {
		t.Fatalf("expected ClassDecl, got %T", defaultExport.Decl)
	}

	if classDecl.Name != "Foo" {
		t.Errorf("expected class name 'Foo', got %q", classDecl.Name)
	}
}

// Test default export function
func TestParseDefaultExportFunction(t *testing.T) {
	input := `export default function add(a: int, b: int): int {
		return a + b
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	defaultExport, ok := program.Statements[0].(*ast.DefaultExport)
	if !ok {
		t.Fatalf("expected *ast.DefaultExport, got %T", program.Statements[0])
	}

	funcDecl, ok := defaultExport.Decl.(*ast.FuncDecl)
	if !ok {
		t.Fatalf("expected FuncDecl, got %T", defaultExport.Decl)
	}

	if funcDecl.Name != "add" {
		t.Errorf("expected function name 'add', got %q", funcDecl.Name)
	}
}

// Test default import: import Foo from "./module"
func TestParseDefaultImport(t *testing.T) {
	input := `import Foo from "./module"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	defaultImport, ok := program.Statements[0].(*ast.DefaultImport)
	if !ok {
		t.Fatalf("expected *ast.DefaultImport, got %T", program.Statements[0])
	}

	if defaultImport.Name != "Foo" {
		t.Errorf("expected name 'Foo', got %q", defaultImport.Name)
	}

	if defaultImport.Path != "./module" {
		t.Errorf("expected path './module', got %q", defaultImport.Path)
	}
}

// Test namespace import: import * as utils from "./utils"
func TestParseNamespaceImport(t *testing.T) {
	input := `import * as utils from "./utils"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	nsImport, ok := program.Statements[0].(*ast.NamespaceImport)
	if !ok {
		t.Fatalf("expected *ast.NamespaceImport, got %T", program.Statements[0])
	}

	if nsImport.Alias != "utils" {
		t.Errorf("expected alias 'utils', got %q", nsImport.Alias)
	}

	if nsImport.Path != "./utils" {
		t.Errorf("expected path './utils', got %q", nsImport.Path)
	}
}
