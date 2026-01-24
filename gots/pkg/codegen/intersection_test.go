package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestCodegenIntersectionType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "intersection type alias with object types",
			input: `
type AB = { x: int } & { y: int }
let obj: AB = { x: 1, y: 2 }
`,
			expected: []string{
				"type AB struct",
				"X int",
				"Y int",
			},
		},
		{
			name: "intersection parameter",
			input: `
function process(value: { x: int } & { y: int }): void {
    println(value)
}
`,
			expected: []string{
				"func process(value struct",
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

func TestCodegenIntersectionObjectMerge(t *testing.T) {
	input := `
type Point = { x: int, y: int }
type Named = { name: string }
type NamedPoint = Point & Named
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

	// The intersection should generate a struct with merged properties
	if !strings.Contains(code, "type NamedPoint struct") {
		t.Errorf("generated code should contain NamedPoint struct\nGenerated:\n%s", code)
	}
	if !strings.Contains(code, "X int") || !strings.Contains(code, "Y int") || !strings.Contains(code, "Name string") {
		t.Errorf("NamedPoint should have all merged properties\nGenerated:\n%s", code)
	}
}
