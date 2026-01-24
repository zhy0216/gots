package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

// Test string.split() method
func TestStringMethod_Split(t *testing.T) {
	input := `
let str: string = "a,b,c"
let parts: string[] = str.split(",")
println(len(parts))
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.Split
	if !strings.Contains(output, "strings.Split") {
		t.Errorf("expected strings.Split call in output, got:\n%s", output)
	}
}

// Test array.join() method
func TestStringMethod_Join(t *testing.T) {
	input := `
let parts: string[] = ["a", "b", "c"]
let str: string = parts.join(",")
println(str)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.Join
	if !strings.Contains(output, "strings.Join") {
		t.Errorf("expected strings.Join call in output, got:\n%s", output)
	}
}

// Test string.replace() method
func TestStringMethod_Replace(t *testing.T) {
	input := `
let str: string = "hello world"
let newStr: string = str.replace("world", "goTS")
println(newStr)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.Replace
	if !strings.Contains(output, "strings.Replace") {
		t.Errorf("expected strings.Replace call in output, got:\n%s", output)
	}
}

// Test string.trim() method
func TestStringMethod_Trim(t *testing.T) {
	input := `
let str: string = "  hello  "
let trimmed: string = str.trim()
println(trimmed)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.TrimSpace
	if !strings.Contains(output, "strings.TrimSpace") {
		t.Errorf("expected strings.TrimSpace call in output, got:\n%s", output)
	}
}

// Test string.startsWith() method
func TestStringMethod_StartsWith(t *testing.T) {
	input := `
let str: string = "hello world"
let starts: boolean = str.startsWith("hello")
println(starts)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.HasPrefix
	if !strings.Contains(output, "strings.HasPrefix") {
		t.Errorf("expected strings.HasPrefix call in output, got:\n%s", output)
	}
}

// Test string.endsWith() method
func TestStringMethod_EndsWith(t *testing.T) {
	input := `
let str: string = "hello world"
let ends: boolean = str.endsWith("world")
println(ends)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.HasSuffix
	if !strings.Contains(output, "strings.HasSuffix") {
		t.Errorf("expected strings.HasSuffix call in output, got:\n%s", output)
	}
}

// Test string.includes() method
func TestStringMethod_Includes(t *testing.T) {
	input := `
let str: string = "hello world"
let has: boolean = str.includes("world")
println(has)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.Contains
	if !strings.Contains(output, "strings.Contains") {
		t.Errorf("expected strings.Contains call in output, got:\n%s", output)
	}
}

// Test string.toLowerCase() method
func TestStringMethod_ToLowerCase(t *testing.T) {
	input := `
let str: string = "HELLO"
let lower: string = str.toLowerCase()
println(lower)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.ToLower
	if !strings.Contains(output, "strings.ToLower") {
		t.Errorf("expected strings.ToLower call in output, got:\n%s", output)
	}
}

// Test string.toUpperCase() method
func TestStringMethod_ToUpperCase(t *testing.T) {
	input := `
let str: string = "hello"
let upper: string = str.toUpperCase()
println(upper)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.ToUpper
	if !strings.Contains(output, "strings.ToUpper") {
		t.Errorf("expected strings.ToUpper call in output, got:\n%s", output)
	}
}

// Test string.indexOf() method
func TestStringMethod_IndexOf(t *testing.T) {
	input := `
let str: string = "hello world"
let idx: int = str.indexOf("world")
println(idx)
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProg := builder.Build(program)

	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	code, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	output := string(code)

	// Should use strings.Index
	if !strings.Contains(output, "strings.Index") {
		t.Errorf("expected strings.Index call in output, got:\n%s", output)
	}
}
