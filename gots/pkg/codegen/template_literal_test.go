package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

func TestTemplateLiteralCodegen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple template literal",
			input:    "let x: string = `hello world`",
			expected: `var x string = "hello world"`,
		},
		{
			name:     "template with one expression",
			input:    "let name: string = \"Alice\"\nlet x: string = `Hello, ${name}!`",
			expected: `fmt.Sprintf("Hello, %v!", name)`,
		},
		{
			name:     "template with multiple expressions",
			input:    "let a: int = 1\nlet b: int = 2\nlet x: string = `${a} + ${b} = ${a + b}`",
			expected: `fmt.Sprintf("%v + %v = %v", a, b, (a + b))`,
		},
		{
			name:     "template with complex expression",
			input:    "let x: int = 5\nlet s: string = `Result: ${x * 2}`",
			expected: `fmt.Sprintf("Result: %v", (x * 2))`,
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
			typedProg := builder.Build(program)

			if len(builder.Errors()) > 0 {
				t.Fatalf("type errors: %v", builder.Errors())
			}

			output, err := Generate(typedProg)
			if err != nil {
				t.Fatalf("codegen error: %v", err)
			}

			if !strings.Contains(string(output), tt.expected) {
				t.Errorf("expected output to contain:\n%s\n\ngot:\n%s", tt.expected, string(output))
			}
		})
	}
}

func TestTemplateLiteralEndToEnd(t *testing.T) {
	input := `
let name: string = "World"
let greeting: string = ` + "`Hello, ${name}!`" + `
println(greeting)
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

	// Verify the output contains expected patterns
	code := string(output)
	if !strings.Contains(code, "fmt.Sprintf") {
		t.Errorf("expected output to contain fmt.Sprintf, got:\n%s", code)
	}
	if !strings.Contains(code, `"Hello, %v!"`) {
		t.Errorf("expected output to contain format string, got:\n%s", code)
	}
}
