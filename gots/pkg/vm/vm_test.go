package vm

import (
	"bytes"
	"testing"

	"github.com/pocketlang/gots/pkg/compiler"
	"github.com/pocketlang/gots/pkg/lexer"
	"github.com/pocketlang/gots/pkg/parser"
)

func TestVMNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42;", 42},
		{"3.14;", 3.14},
	}

	for _, tt := range tests {
		vm := runVM(t, tt.input)
		if !vm.lastPopped.IsNumber() {
			t.Fatalf("expected number, got %v", vm.lastPopped.Type)
		}
		if vm.lastPopped.AsNumber() != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, vm.lastPopped.AsNumber())
		}
	}
}

func TestVMBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		vm := runVM(t, tt.input)
		if !vm.lastPopped.IsBool() {
			t.Fatalf("expected bool, got %v", vm.lastPopped.Type)
		}
		if vm.lastPopped.AsBool() != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, vm.lastPopped.AsBool())
		}
	}
}

func TestVMNullLiteral(t *testing.T) {
	vm := runVM(t, "null;")
	if !vm.lastPopped.IsNull() {
		t.Errorf("expected null, got %v", vm.lastPopped.Type)
	}
}

func TestVMStringLiteral(t *testing.T) {
	vm := runVM(t, `"hello";`)
	if !vm.lastPopped.IsString() {
		t.Fatalf("expected string, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsString() != "hello" {
		t.Errorf("expected 'hello', got %q", vm.lastPopped.AsString())
	}
}

func TestVMArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1 + 2;", 3},
		{"5 - 3;", 2},
		{"4 * 3;", 12},
		{"10 / 2;", 5},
		{"10 % 3;", 1},
		{"-5;", -5},
		{"--5;", 5},
		{"2 + 3 * 4;", 14},
		{"(2 + 3) * 4;", 20},
	}

	for _, tt := range tests {
		vm := runVM(t, tt.input)
		if !vm.lastPopped.IsNumber() {
			t.Fatalf("input %q: expected number, got %v", tt.input, vm.lastPopped.Type)
		}
		if vm.lastPopped.AsNumber() != tt.expected {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, vm.lastPopped.AsNumber())
		}
	}
}

func TestVMComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"1 == 1;", true},
		{"1 == 2;", false},
		{"1 != 2;", true},
		{"1 != 1;", false},
		{"1 < 2;", true},
		{"2 < 1;", false},
		{"1 <= 1;", true},
		{"1 <= 2;", true},
		{"2 <= 1;", false},
		{"2 > 1;", true},
		{"1 > 2;", false},
		{"1 >= 1;", true},
		{"2 >= 1;", true},
		{"1 >= 2;", false},
	}

	for _, tt := range tests {
		vm := runVM(t, tt.input)
		if !vm.lastPopped.IsBool() {
			t.Fatalf("input %q: expected bool, got %v", tt.input, vm.lastPopped.Type)
		}
		if vm.lastPopped.AsBool() != tt.expected {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, vm.lastPopped.AsBool())
		}
	}
}

func TestVMLogical(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
	}

	for _, tt := range tests {
		vm := runVM(t, tt.input)
		if !vm.lastPopped.IsBool() {
			t.Fatalf("input %q: expected bool, got %v", tt.input, vm.lastPopped.Type)
		}
		if vm.lastPopped.AsBool() != tt.expected {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, vm.lastPopped.AsBool())
		}
	}
}

func TestVMStringConcat(t *testing.T) {
	vm := runVM(t, `"hello" + "world";`)
	if !vm.lastPopped.IsString() {
		t.Fatalf("expected string, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsString() != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", vm.lastPopped.AsString())
	}
}

func TestVMPrintln(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, "println(42);", &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "42\n" {
		t.Errorf("expected '42\\n', got %q", output)
	}
}

func TestVMPrintlnString(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `println("hello");`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", output)
	}
}

func TestVMMultipleStatements(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `println(1 + 2); println(3 * 4);`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "3\n12\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

// Helper functions

func runVM(t *testing.T, source string) *VM {
	t.Helper()
	vm := compileAndCreateVM(t, source)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}
	return vm
}

func newVMWithOutput(t *testing.T, source string, output *bytes.Buffer) *VM {
	t.Helper()
	vm := compileAndCreateVM(t, source)
	vm.output = output
	return vm
}

func compileAndCreateVM(t *testing.T, source string) *VM {
	t.Helper()

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	c := compiler.New()
	chunk, err := c.Compile(program)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	vm := New(chunk)
	return vm
}
