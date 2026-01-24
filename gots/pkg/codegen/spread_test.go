package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestSpreadInArrayLiteralCodegen(t *testing.T) {
	input := `
let other: int[] = [1, 2]
let arr: int[] = [...other, 3, 4]
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
		t.Fatalf("type errors: %v", builder.Errors())
	}

	output, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	code := string(output)

	// The spread should generate append
	if !strings.Contains(code, "append(") {
		t.Errorf("expected output to contain append for spread, got:\n%s", code)
	}
	// Variable name is lowercase in Go for local variables
	if !strings.Contains(code, "other") {
		t.Errorf("expected output to contain 'other' variable, got:\n%s", code)
	}
}

func TestSpreadOnlyArrayCodegen(t *testing.T) {
	input := `
let source: int[] = [1, 2, 3]
let copy: int[] = [...source]
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
		t.Fatalf("type errors: %v", builder.Errors())
	}

	output, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	code := string(output)

	// Spreading a single array should create a copy using append
	if !strings.Contains(code, "source") {
		t.Errorf("expected output to contain 'source' variable, got:\n%s", code)
	}
	if !strings.Contains(code, "append(") {
		t.Errorf("expected output to contain append for spread copy, got:\n%s", code)
	}
}

func TestMultipleSpreadsCodegen(t *testing.T) {
	input := `
let a: int[] = [1]
let b: int[] = [2]
let combined: int[] = [...a, ...b, 3]
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
		t.Fatalf("type errors: %v", builder.Errors())
	}

	output, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	code := string(output)

	// Multiple spreads should use multiple appends
	if !strings.Contains(code, "append(") {
		t.Errorf("expected output to contain append for spreads, got:\n%s", code)
	}
}

func TestSpreadInFunctionCallCodegen(t *testing.T) {
	input := `
let args: int[] = [1, 2, 3]
println(...args)
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
		t.Fatalf("type errors: %v", builder.Errors())
	}

	output, err := Generate(typedProg)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}

	code := string(output)

	// Function call spread should use ... (lowercase variable name)
	if !strings.Contains(code, "args...") {
		t.Errorf("expected output to contain args... for spread in function call, got:\n%s", code)
	}
}
