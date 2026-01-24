package declaration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/codegen"
	"github.com/zhy0216/quickts/gots/pkg/declaration"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

// End-to-end tests that verify the full pipeline:
// Declaration file → Type checking → Code generation

func TestE2EDeclarationToCodegen(t *testing.T) {
	// Create a temporary directory with a custom declaration
	tmpDir, err := os.MkdirTemp("", "gots-e2e-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a test declaration file for a custom package
	declContent := `
declare module "go:mylib" {
    function Greet(name: string): string
    function Add(a: int, b: int): int
    const Version: string
}
`
	declPath := filepath.Join(tmpDir, "go_mylib.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	// Add search path to default loader
	declaration.DefaultLoader.AddSearchPath(tmpDir)

	// goTS source that uses the declared package
	source := `
import { Greet, Add } from "go:mylib"

let greeting: string = Greet("World")
let sum: int = Add(1, 2)
println(greeting)
println(sum)
`

	// Parse
	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	// Type check
	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		t.Fatalf("Type errors: %v", builder.Errors())
	}

	// Generate Go code
	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		t.Fatalf("Code generation error: %v", err)
	}

	goCodeStr := string(goCode)

	// Verify the generated code contains expected imports
	if !strings.Contains(goCodeStr, `"mylib"`) {
		t.Error("Expected generated code to import mylib")
	}

	// Verify function calls are generated
	if !strings.Contains(goCodeStr, "mylib.Greet") {
		t.Error("Expected generated code to call mylib.Greet")
	}
	if !strings.Contains(goCodeStr, "mylib.Add") {
		t.Error("Expected generated code to call mylib.Add")
	}
}

func TestE2EDeclarationWithTypes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-e2e-types-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Declaration with various type patterns
	declContent := `
declare module "go:typedlib" {
    function GetStrings(): string[]
    function ProcessArray(arr: int[]): int
    function MaybeValue(): string | null
}
`
	declPath := filepath.Join(tmpDir, "go_typedlib.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	declaration.DefaultLoader.AddSearchPath(tmpDir)

	source := `
import { GetStrings, ProcessArray } from "go:typedlib"

let strings: string[] = GetStrings()
let arr: int[] = [1, 2, 3]
let result: int = ProcessArray(arr)
`

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		t.Fatalf("Type errors: %v", builder.Errors())
	}

	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		t.Fatalf("Code generation error: %v", err)
	}

	goCodeStr := string(goCode)

	// Verify array types are handled correctly
	if !strings.Contains(goCodeStr, "typedlib.GetStrings") {
		t.Error("Expected generated code to call typedlib.GetStrings")
	}
	if !strings.Contains(goCodeStr, "typedlib.ProcessArray") {
		t.Error("Expected generated code to call typedlib.ProcessArray")
	}
}

func TestE2EDeclarationTypeErrorDetection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-e2e-error-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Declaration with specific types
	declContent := `
declare module "go:strictlib" {
    function RequiresInt(n: int): void
    function RequiresString(s: string): void
}
`
	declPath := filepath.Join(tmpDir, "go_strictlib.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	declaration.DefaultLoader.AddSearchPath(tmpDir)

	// goTS source with type mismatch
	source := `
import { RequiresInt } from "go:strictlib"

RequiresInt("not an int")
`

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	_ = builder.Build(program)

	// Should have type errors
	if !builder.HasErrors() {
		t.Error("Expected type error for passing string to int parameter")
	}
}

func TestE2EMultipleDeclarationImports(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-e2e-multi-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// First declaration
	decl1 := `
declare module "go:lib1" {
    function Func1(): int
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go_lib1.d.gts"), []byte(decl1), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	// Second declaration
	decl2 := `
declare module "go:lib2" {
    function Func2(): string
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go_lib2.d.gts"), []byte(decl2), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	declaration.DefaultLoader.AddSearchPath(tmpDir)

	// Source using both
	source := `
import { Func1 } from "go:lib1"
import { Func2 } from "go:lib2"

let n: int = Func1()
let s: string = Func2()
println(n)
println(s)
`

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		t.Fatalf("Type errors: %v", builder.Errors())
	}

	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		t.Fatalf("Code generation error: %v", err)
	}

	goCodeStr := string(goCode)

	// Both imports should be in generated code
	if !strings.Contains(goCodeStr, "lib1") {
		t.Error("Expected generated code to import lib1")
	}
	if !strings.Contains(goCodeStr, "lib2") {
		t.Error("Expected generated code to import lib2")
	}
}

func TestE2EDeclarationWithConstants(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-e2e-const-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:constlib" {
    const MaxSize: int
    const DefaultName: string
    const Pi: float
}
`
	declPath := filepath.Join(tmpDir, "go_constlib.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	declaration.DefaultLoader.AddSearchPath(tmpDir)

	source := `
import { MaxSize, Pi } from "go:constlib"

let size: int = MaxSize
let pi: float = Pi
println(size)
println(pi)
`

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		t.Fatalf("Type errors: %v", builder.Errors())
	}

	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		t.Fatalf("Code generation error: %v", err)
	}

	goCodeStr := string(goCode)

	if !strings.Contains(goCodeStr, "constlib.MaxSize") {
		t.Error("Expected generated code to reference constlib.MaxSize")
	}
	if !strings.Contains(goCodeStr, "constlib.Pi") {
		t.Error("Expected generated code to reference constlib.Pi")
	}
}

func TestE2EStdlibFallback(t *testing.T) {
	// Test that standard library imports still work via the fallback registry
	source := `
import { ToUpper, ToLower } from "go:strings"
import { Sqrt } from "go:math"

let upper: string = ToUpper("hello")
let lower: string = ToLower("WORLD")
let root: float = Sqrt(16.0)
println(upper)
println(lower)
println(root)
`

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if builder.HasErrors() {
		t.Fatalf("Type errors: %v", builder.Errors())
	}

	goCode, err := codegen.Generate(typedProg)
	if err != nil {
		t.Fatalf("Code generation error: %v", err)
	}

	goCodeStr := string(goCode)

	// Verify stdlib imports work
	if !strings.Contains(goCodeStr, `"strings"`) {
		t.Error("Expected generated code to import strings")
	}
	if !strings.Contains(goCodeStr, `"math"`) {
		t.Error("Expected generated code to import math")
	}
	if !strings.Contains(goCodeStr, "strings.ToUpper") {
		t.Error("Expected generated code to call strings.ToUpper")
	}
}
