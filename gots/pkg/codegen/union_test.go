package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestCodegenUnionType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "union type alias",
			input: `
type StringOrInt = string | int

function getValue(): StringOrInt {
    return 42
}
`,
			expected: []string{
				"type StringOrInt interface{}",
				"func getValue() interface{}",
				"return 42",
			},
		},
		{
			name: "union parameter",
			input: `
function process(value: string | int): void {
    println(value)
}
`,
			expected: []string{
				"func process(value interface{})",
			},
		},
		{
			name: "union variable",
			input: `
let x: string | int = "hello"
x = 42
`,
			expected: []string{
				"var x interface{} = \"hello\"",
				"x = 42",
			},
		},
		{
			name: "three way union",
			input: `
type Value = string | int | boolean
let v: Value = true
`,
			expected: []string{
				"type Value interface{}",
				"var v interface{} = true",
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

func TestCodegenUnionWithNull(t *testing.T) {
	input := `
type OptionalString = string | null
let x: OptionalString = null
x = "hello"
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

	// Since union with null can use nullable (pointer) or interface{}
	// Just check that it compiles without errors
	if !strings.Contains(code, "var x") {
		t.Errorf("generated code does not contain variable declaration\nGenerated:\n%s", code)
	}
}

func TestCodegenUnionArrayTypes(t *testing.T) {
	input := `
type ArrayUnion = int[] | string[]
let arr: ArrayUnion = [1, 2, 3]
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

	if !strings.Contains(code, "type ArrayUnion interface{}") {
		t.Errorf("generated code does not contain union type alias\nGenerated:\n%s", code)
	}
	if !strings.Contains(code, "var arr interface{}") {
		t.Errorf("generated code does not contain variable declaration\nGenerated:\n%s", code)
	}
}
