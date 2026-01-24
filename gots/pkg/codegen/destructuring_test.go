package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestArrayDestructuringCodegen(t *testing.T) {
	input := `
let arr: int[] = [1, 2, 3]
let [a, b, c]: int[] = arr
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

	// Verify the output contains expected patterns
	if !strings.Contains(code, "_destructure_temp") {
		t.Errorf("expected output to contain temp variable, got:\n%s", code)
	}
	if !strings.Contains(code, "_destructure_temp[0]") {
		t.Errorf("expected output to contain indexed access, got:\n%s", code)
	}
}

func TestObjectDestructuringCodegen(t *testing.T) {
	input := `
let point: {x: int, y: int} = {x: 10, y: 20}
let {x, y}: {x: int, y: int} = point
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

	// Verify the output contains expected patterns
	if !strings.Contains(code, "_destructure_temp") {
		t.Errorf("expected output to contain temp variable, got:\n%s", code)
	}
	if !strings.Contains(code, "_destructure_temp.X") {
		t.Errorf("expected output to contain property access (exported), got:\n%s", code)
	}
}
