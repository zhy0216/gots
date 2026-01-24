package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

// Test array.map() method
func TestArrayMethod_Map(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3]
let doubled: int[] = arr.map((x: int): int => x * 2)
println(doubled)
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

	// Should generate inline map with result array
	if !strings.Contains(output, "result := make([]int") {
		t.Errorf("expected inline map implementation, got:\n%s", output)
	}
}

// Test array.filter() method
func TestArrayMethod_Filter(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let evens: int[] = arr.filter((x: int): boolean => x % 2 == 0)
println(evens)
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

	// Should generate inline filter
	if !strings.Contains(output, "result := make([]int, 0)") {
		t.Errorf("expected inline filter implementation, got:\n%s", output)
	}
}

// Test array.reduce() method
func TestArrayMethod_Reduce(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let sum: int = arr.reduce((acc: int, x: int): int => acc + x, 0)
println(sum)
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

	// Should generate inline reduce
	if !strings.Contains(output, "acc :=") {
		t.Errorf("expected inline reduce implementation, got:\n%s", output)
	}
}

// Test array.find() method
func TestArrayMethod_Find(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let found: int | null = arr.find((x: int): boolean => x > 3)
println(found)
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

	// Should generate inline find
	if !strings.Contains(output, "return &v") || !strings.Contains(output, "return nil") {
		t.Errorf("expected inline find implementation, got:\n%s", output)
	}
}

// Test array.findIndex() method
func TestArrayMethod_FindIndex(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let idx: int = arr.findIndex((x: int): boolean => x > 3)
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

	// Should generate inline findIndex
	if !strings.Contains(output, "return i") || !strings.Contains(output, "return -1") {
		t.Errorf("expected inline findIndex implementation, got:\n%s", output)
	}
}

// Test array.some() method
func TestArrayMethod_Some(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let hasEven: boolean = arr.some((x: int): boolean => x % 2 == 0)
println(hasEven)
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

	// Should generate inline some
	if !strings.Contains(output, "return true") || !strings.Contains(output, "return false") {
		t.Errorf("expected inline some implementation, got:\n%s", output)
	}
}

// Test array.every() method
func TestArrayMethod_Every(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3, 4, 5]
let allPositive: boolean = arr.every((x: int): boolean => x > 0)
println(allPositive)
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

	// Should generate inline every
	if !strings.Contains(output, "return false") || !strings.Contains(output, "return true") {
		t.Errorf("expected inline every implementation, got:\n%s", output)
	}
}
