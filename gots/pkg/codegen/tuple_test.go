package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestCodegenTupleType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "tuple type alias",
			input: `
type Pair = [string, int]
let p: Pair = ["hello", 42]
`,
			expected: []string{
				"type Pair struct",
				"T0 string",
				"T1 int",
			},
		},
		{
			name: "tuple parameter",
			input: `
function process(pair: [string, int]): void {
    println(pair)
}
`,
			expected: []string{
				"func process(pair struct",
			},
		},
		{
			name: "tuple variable",
			input: `
let point: [int, int] = [1, 2]
`,
			expected: []string{
				"point struct",
				"T0 int",
				"T1 int",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("parser errors: %v", p.Errors())
			}

			builder := typed.NewBuilder()
			typedProgram := builder.Build(program)
			if len(builder.Errors()) > 0 {
				t.Fatalf("type checker errors: %v", builder.Errors())
			}

			output, err := Generate(typedProgram)
			if err != nil {
				t.Fatalf("codegen error: %v", err)
			}
			code := string(output)

			for _, expected := range tt.expected {
				if !strings.Contains(code, expected) {
					t.Errorf("generated code does not contain %q\nGenerated:\n%s", expected, code)
				}
			}
		})
	}
}

func TestCodegenTupleWithThreeElements(t *testing.T) {
	input := `
type Triple = [string, int, boolean]
let t: Triple = ["hello", 42, true]
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProgram := builder.Build(program)
	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	output, err := Generate(typedProgram)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}
	code := string(output)

	// Check that the struct has three fields
	if !strings.Contains(code, "T0 string") {
		t.Errorf("missing T0 field\nGenerated:\n%s", code)
	}
	if !strings.Contains(code, "T1 int") {
		t.Errorf("missing T1 field\nGenerated:\n%s", code)
	}
	if !strings.Contains(code, "T2 bool") {
		t.Errorf("missing T2 field\nGenerated:\n%s", code)
	}
}

func TestCodegenNestedTuple(t *testing.T) {
	input := `
type Nested = [[int, int], string]
let n: Nested = [[1, 2], "test"]
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	builder := typed.NewBuilder()
	typedProgram := builder.Build(program)
	if len(builder.Errors()) > 0 {
		t.Fatalf("type checker errors: %v", builder.Errors())
	}

	output, err := Generate(typedProgram)
	if err != nil {
		t.Fatalf("codegen error: %v", err)
	}
	code := string(output)

	// The nested tuple should generate a nested struct
	if !strings.Contains(code, "type Nested struct") {
		t.Errorf("missing Nested struct\nGenerated:\n%s", code)
	}
}
