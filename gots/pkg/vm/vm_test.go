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

func TestVMGlobalVariable(t *testing.T) {
	vm := runVM(t, `let x: number = 42; x;`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 42 {
		t.Errorf("expected 42, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMGlobalVariableAssignment(t *testing.T) {
	vm := runVM(t, `let x: number = 1; x = 42; x;`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 42 {
		t.Errorf("expected 42, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMGlobalVariableInExpression(t *testing.T) {
	vm := runVM(t, `let x: number = 10; let y: number = 5; x + y;`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 15 {
		t.Errorf("expected 15, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMGlobalVariablePrintln(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `let x: number = 42; println(x);`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "42\n" {
		t.Errorf("expected '42\\n', got %q", output)
	}
}

func TestVMGlobalStringVariable(t *testing.T) {
	vm := runVM(t, `let s: string = "hello"; s;`)
	if !vm.lastPopped.IsString() {
		t.Fatalf("expected string, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsString() != "hello" {
		t.Errorf("expected 'hello', got %q", vm.lastPopped.AsString())
	}
}

func TestVMIfStatementTrue(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `if (true) { println(1); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "1\n" {
		t.Errorf("expected '1\\n', got %q", output)
	}
}

func TestVMIfStatementFalse(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `if (false) { println(1); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("expected no output, got %q", output)
	}
}

func TestVMIfElseStatementTrue(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `if (true) { println(1); } else { println(2); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "1\n" {
		t.Errorf("expected '1\\n', got %q", output)
	}
}

func TestVMIfElseStatementFalse(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `if (false) { println(1); } else { println(2); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "2\n" {
		t.Errorf("expected '2\\n', got %q", output)
	}
}

func TestVMIfWithCondition(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `let x: number = 5; if (x > 3) { println(1); } else { println(2); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "1\n" {
		t.Errorf("expected '1\\n', got %q", output)
	}
}

func TestVMWhileLoop(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `let x: number = 0; while (x < 3) { println(x); x = x + 1; }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "0\n1\n2\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestVMWhileLoopNeverExecutes(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `while (false) { println(1); }`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "" {
		t.Errorf("expected no output, got %q", output)
	}
}

func TestVMWhileLoopSum(t *testing.T) {
	vm := runVM(t, `let sum: number = 0; let i: number = 1; while (i <= 5) { sum = sum + i; i = i + 1; } sum;`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 15 {
		t.Errorf("expected 15 (1+2+3+4+5), got %v", vm.lastPopped.AsNumber())
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

// ============================================================
// Functions Tests
// ============================================================

func TestVMSimpleFunction(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `
		function greet(): void {
			println(42);
		}
		greet();
	`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	if output != "42\n" {
		t.Errorf("expected '42\\n', got %q", output)
	}
}

func TestVMFunctionWithReturn(t *testing.T) {
	vm := runVM(t, `
		function add(a: number, b: number): number {
			return a + b;
		}
		add(3, 4);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 7 {
		t.Errorf("expected 7, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMFunctionWithMultipleParameters(t *testing.T) {
	vm := runVM(t, `
		function sum3(a: number, b: number, c: number): number {
			return a + b + c;
		}
		sum3(10, 20, 30);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 60 {
		t.Errorf("expected 60, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMFunctionLocalVariables(t *testing.T) {
	vm := runVM(t, `
		function compute(x: number): number {
			let y: number = x * 2;
			let z: number = y + 10;
			return z;
		}
		compute(5);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 20 {
		t.Errorf("expected 20 (5*2+10), got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMRecursiveFunction(t *testing.T) {
	vm := runVM(t, `
		function factorial(n: number): number {
			if (n <= 1) {
				return 1;
			}
			return n * factorial(n - 1);
		}
		factorial(5);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 120 {
		t.Errorf("expected 120 (5!), got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMFibonacci(t *testing.T) {
	vm := runVM(t, `
		function fib(n: number): number {
			if (n < 2) {
				return n;
			}
			return fib(n - 1) + fib(n - 2);
		}
		fib(10);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 55 {
		t.Errorf("expected 55 (fib(10)), got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMNestedFunctionCalls(t *testing.T) {
	vm := runVM(t, `
		function double(x: number): number {
			return x * 2;
		}
		function quadruple(x: number): number {
			return double(double(x));
		}
		quadruple(5);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 20 {
		t.Errorf("expected 20, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMFunctionExpression(t *testing.T) {
	vm := runVM(t, `
		let square: Function = function(x: number): number {
			return x * x;
		};
		square(7);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 49 {
		t.Errorf("expected 49, got %v", vm.lastPopped.AsNumber())
	}
}

// ============================================================
// Closures Tests
// ============================================================

func TestVMClosureSimple(t *testing.T) {
	vm := runVM(t, `
		function makeAdder(x: number): number {
			return function(y: number): number {
				return x + y;
			};
		}
		let add5: Function = makeAdder(5);
		add5(10);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 15 {
		t.Errorf("expected 15, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMClosureCounter(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `
		function makeCounter(): number {
			let count: number = 0;
			return function(): number {
				count = count + 1;
				return count;
			};
		}
		let counter: Function = makeCounter();
		println(counter());
		println(counter());
		println(counter());
	`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "1\n2\n3\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestVMClosureMultipleInstances(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `
		function makeCounter(): number {
			let count: number = 0;
			return function(): number {
				count = count + 1;
				return count;
			};
		}
		let c1: Function = makeCounter();
		let c2: Function = makeCounter();
		println(c1());
		println(c1());
		println(c2());
		println(c1());
		println(c2());
	`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "1\n2\n1\n3\n2\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

func TestVMClosureNestedCapture(t *testing.T) {
	vm := runVM(t, `
		function outer(a: number): number {
			let b: number = 10;
			function middle(): number {
				let c: number = 100;
				return function(): number {
					return a + b + c;
				};
			}
			return middle();
		}
		let f: Function = outer(1);
		f();
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 111 {
		t.Errorf("expected 111 (1+10+100), got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMClosureCapturesParameter(t *testing.T) {
	vm := runVM(t, `
		function multiplier(factor: number): number {
			return function(x: number): number {
				return x * factor;
			};
		}
		let triple: Function = multiplier(3);
		triple(7);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 21 {
		t.Errorf("expected 21, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMClosureModifiesEnclosedVariable(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `
		function createAccumulator(): number {
			let total: number = 0;
			return function(x: number): number {
				total = total + x;
				return total;
			};
		}
		let acc: Function = createAccumulator();
		println(acc(5));
		println(acc(10));
		println(acc(3));
	`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "5\n15\n18\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

// ============================================================
// Local Variables in Blocks Tests
// ============================================================

func TestVMLocalVariablesInBlock(t *testing.T) {
	vm := runVM(t, `
		let x: number = 10;
		{
			let x: number = 20;
			x;
		}
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 20 {
		t.Errorf("expected 20 (inner x), got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMLocalVariableShadowing(t *testing.T) {
	var buf bytes.Buffer
	vm := newVMWithOutput(t, `
		let x: number = 10;
		println(x);
		{
			let x: number = 20;
			println(x);
		}
		println(x);
	`, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("VM error: %v", err)
	}

	output := buf.String()
	expected := "10\n20\n10\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}

// ============================================================
// Phase 5: Arrays Tests
// ============================================================

func TestVMArrayLiteral(t *testing.T) {
	vm := runVM(t, `
		let arr: number[] = [1, 2, 3];
		arr;
	`)
	if !vm.lastPopped.IsArray() {
		t.Fatalf("expected array, got %v", vm.lastPopped.Type)
	}
	arr := vm.lastPopped.AsArray()
	if len(arr.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr.Elements))
	}
}

func TestVMArrayIndexRead(t *testing.T) {
	vm := runVM(t, `
		let arr: number[] = [10, 20, 30];
		arr[1];
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 20 {
		t.Errorf("expected 20, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMArrayIndexWrite(t *testing.T) {
	vm := runVM(t, `
		let arr: number[] = [1, 2, 3];
		arr[1] = 42;
		arr[1];
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 42 {
		t.Errorf("expected 42, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMArrayLen(t *testing.T) {
	vm := runVM(t, `
		let arr: number[] = [1, 2, 3, 4, 5];
		len(arr);
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 5 {
		t.Errorf("expected 5, got %v", vm.lastPopped.AsNumber())
	}
}

// ============================================================
// Phase 5: Object Literals Tests
// ============================================================

func TestVMObjectLiteral(t *testing.T) {
	vm := runVM(t, `
		let obj: {x: number, y: number} = {x: 10, y: 20};
		obj.x;
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 10 {
		t.Errorf("expected 10, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMObjectPropertyWrite(t *testing.T) {
	vm := runVM(t, `
		let obj: {x: number} = {x: 10};
		obj.x = 42;
		obj.x;
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 42 {
		t.Errorf("expected 42, got %v", vm.lastPopped.AsNumber())
	}
}

// ============================================================
// Phase 5: Classes Tests
// ============================================================

func TestVMClassBasic(t *testing.T) {
	vm := runVM(t, `
		class Point {
			x: number;
			y: number;
			constructor(x: number, y: number) {
				this.x = x;
				this.y = y;
			}
		}
		let p: Point = new Point(3, 4);
		p.x;
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 3 {
		t.Errorf("expected 3, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMClassMethod(t *testing.T) {
	vm := runVM(t, `
		class Counter {
			value: number;
			constructor() {
				this.value = 0;
			}
			increment(): void {
				this.value = this.value + 1;
			}
			get(): number {
				return this.value;
			}
		}
		let c: Counter = new Counter();
		c.increment();
		c.increment();
		c.increment();
		c.get();
	`)
	if !vm.lastPopped.IsNumber() {
		t.Fatalf("expected number, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsNumber() != 3 {
		t.Errorf("expected 3, got %v", vm.lastPopped.AsNumber())
	}
}

func TestVMClassInheritance(t *testing.T) {
	vm := runVM(t, `
		class Animal {
			name: string;
			constructor(name: string) {
				this.name = name;
			}
			speak(): string {
				return this.name;
			}
		}
		class Dog extends Animal {
			constructor(name: string) {
				super(name);
			}
		}
		let d: Dog = new Dog("Buddy");
		d.speak();
	`)
	if !vm.lastPopped.IsString() {
		t.Fatalf("expected string, got %v", vm.lastPopped.Type)
	}
	if vm.lastPopped.AsString() != "Buddy" {
		t.Errorf("expected 'Buddy', got %v", vm.lastPopped.AsString())
	}
}
