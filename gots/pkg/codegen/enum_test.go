package codegen

import (
	"strings"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/typed"
)

// normalizeSpaces collapses multiple spaces into one for comparison
func normalizeSpaces(s string) string {
	result := strings.Join(strings.Fields(s), " ")
	return result
}

func TestEnumCodegen(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "basic enum",
			input: `enum Color { Red, Green, Blue }`,
			expected: []string{
				"type Color int",
				"ColorRed Color = 0",
				"ColorGreen Color = 1",
				"ColorBlue Color = 2",
			},
		},
		{
			name:  "enum with explicit values",
			input: `enum Status { Pending = 1, Active = 2, Completed = 3 }`,
			expected: []string{
				"type Status int",
				"StatusPending Status = 1",
				"StatusActive Status = 2",
				"StatusCompleted Status = 3",
			},
		},
		{
			name:  "enum with mixed values",
			input: `enum Priority { Low, Medium = 5, High }`,
			expected: []string{
				"type Priority int",
				"PriorityLow Priority = 0",
				"PriorityMedium Priority = 5",
				"PriorityHigh Priority = 6",
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
			typedProg := builder.Build(program)
			if len(builder.Errors()) > 0 {
				t.Fatalf("type errors: %v", builder.Errors())
			}

			output, err := Generate(typedProg)
			if err != nil {
				t.Fatalf("codegen error: %v", err)
			}

			outputStr := string(output)
			normalizedOutput := normalizeSpaces(outputStr)
			for _, exp := range tt.expected {
				if !strings.Contains(normalizedOutput, normalizeSpaces(exp)) {
					t.Errorf("expected output to contain %q\nGot:\n%s", exp, outputStr)
				}
			}
		})
	}
}

func TestEnumMemberAccess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "enum member access",
			input: `
enum Color { Red, Green, Blue }
let c = Color.Red
`,
			expected: []string{
				"var c Color = ColorRed",
			},
		},
		{
			name: "enum member in expression",
			input: `
enum Status { Pending = 1, Active = 2 }
println(Status.Active)
`,
			expected: []string{
				"fmt.Println(StatusActive)",
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
			typedProg := builder.Build(program)
			if len(builder.Errors()) > 0 {
				t.Fatalf("type errors: %v", builder.Errors())
			}

			output, err := Generate(typedProg)
			if err != nil {
				t.Fatalf("codegen error: %v", err)
			}

			outputStr := string(output)
			normalizedOutput := normalizeSpaces(outputStr)
			for _, exp := range tt.expected {
				if !strings.Contains(normalizedOutput, normalizeSpaces(exp)) {
					t.Errorf("expected output to contain %q\nGot:\n%s", exp, outputStr)
				}
			}
		})
	}
}
