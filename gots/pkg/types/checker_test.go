package types

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
)

func checkProgram(t *testing.T, input string) (*Checker, bool) {
	t.Helper()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	c := NewChecker()
	c.Check(program)
	return c, !c.HasErrors()
}

func expectNoErrors(t *testing.T, input string) {
	t.Helper()
	c, ok := checkProgram(t, input)
	if !ok {
		for _, err := range c.Errors() {
			t.Errorf("unexpected error: %s", err.String())
		}
	}
}

func expectError(t *testing.T, input string, expectedMsg string) {
	t.Helper()
	c, ok := checkProgram(t, input)
	if ok {
		t.Fatal("expected type error but got none")
	}

	for _, err := range c.Errors() {
		if contains(err.Message, expectedMsg) {
			return
		}
	}
	t.Errorf("expected error containing %q, got: %v", expectedMsg, c.Errors())
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// Expression Type Checking Tests
// ----------------------------------------------------------------------------

func TestNumberLiteral(t *testing.T) {
	expectNoErrors(t, `let x: number = 42;`)
	expectNoErrors(t, `let x: number = 3.14;`)
	expectError(t, `let x: string = 42;`, "cannot assign number to string")
}

func TestStringLiteral(t *testing.T) {
	expectNoErrors(t, `let x: string = "hello";`)
	expectError(t, `let x: number = "hello";`, "cannot assign string to number")
}

func TestBooleanLiteral(t *testing.T) {
	expectNoErrors(t, `let x: boolean = true;`)
	expectNoErrors(t, `let x: boolean = false;`)
	expectError(t, `let x: string = true;`, "cannot assign boolean to string")
}

func TestNullLiteral(t *testing.T) {
	expectNoErrors(t, `let x: string | null = null;`)
	expectError(t, `let x: string = null;`, "cannot assign null to string")
}

func TestBinaryExpressions(t *testing.T) {
	// Arithmetic
	expectNoErrors(t, `let x: number = 1 + 2;`)
	expectNoErrors(t, `let x: number = 3 - 1;`)
	expectNoErrors(t, `let x: number = 2 * 3;`)
	expectNoErrors(t, `let x: number = 6 / 2;`)
	expectNoErrors(t, `let x: number = 7 % 3;`)

	// String concatenation
	expectNoErrors(t, `let x: string = "hello" + " world";`)

	// Type errors in arithmetic
	expectError(t, `let x: number = "a" - 1;`, "requires numbers")
	expectError(t, `let x: number = true * 2;`, "requires numbers")

	// Comparison
	expectNoErrors(t, `let x: boolean = 1 < 2;`)
	expectNoErrors(t, `let x: boolean = 1 <= 2;`)
	expectNoErrors(t, `let x: boolean = 1 > 2;`)
	expectNoErrors(t, `let x: boolean = 1 >= 2;`)

	// Equality
	expectNoErrors(t, `let x: boolean = 1 == 2;`)
	expectNoErrors(t, `let x: boolean = "a" != "b";`)

	// Logical
	expectNoErrors(t, `let x: boolean = true && false;`)
	expectNoErrors(t, `let x: boolean = true || false;`)
	expectError(t, `let x: boolean = 1 && true;`, "requires booleans")
}

func TestUnaryExpressions(t *testing.T) {
	expectNoErrors(t, `let x: number = -5;`)
	expectNoErrors(t, `let x: boolean = !true;`)
	expectError(t, `let x: number = -"hello";`, "requires number")
	expectError(t, `let x: boolean = !42;`, "requires boolean")
}

// ----------------------------------------------------------------------------
// Variable Tests
// ----------------------------------------------------------------------------

func TestVariableDeclaration(t *testing.T) {
	expectNoErrors(t, `let x: number = 1;`)
	expectNoErrors(t, `const PI: number = 3.14;`)
	expectError(t, `let x: number = "hello";`, "cannot assign string to number")
}

func TestVariableUsage(t *testing.T) {
	expectNoErrors(t, `
		let x: number = 1;
		let y: number = x + 2;
	`)
	expectError(t, `let x: number = y;`, "undefined variable")
}

func TestVariableScoping(t *testing.T) {
	expectNoErrors(t, `
		let x: number = 1;
		{
			let y: number = x + 1;
		}
	`)
}

// ----------------------------------------------------------------------------
// Control Flow Tests
// ----------------------------------------------------------------------------

func TestIfStatement(t *testing.T) {
	expectNoErrors(t, `
		let x: number = 1;
		if (x > 0) {
			let y: number = 2;
		}
	`)
	expectError(t, `
		if (42) {
			let x: number = 1;
		}
	`, "must be boolean")
}

func TestWhileStatement(t *testing.T) {
	expectNoErrors(t, `
		let x: number = 0;
		while (x < 10) {
			x = x + 1;
		}
	`)
	expectError(t, `
		while ("true") {
			let x: number = 1;
		}
	`, "must be boolean")
}

func TestForStatement(t *testing.T) {
	expectNoErrors(t, `
		let sum: number = 0;
		let i: number = 0;
		while (i < 10) {
			sum = sum + i;
			i = i + 1;
		}
	`)
}

func TestBreakContinue(t *testing.T) {
	expectNoErrors(t, `
		while (true) {
			break;
		}
	`)
	expectNoErrors(t, `
		while (true) {
			continue;
		}
	`)
	expectError(t, `break;`, "break outside loop")
	expectError(t, `continue;`, "continue outside loop")
}

// ----------------------------------------------------------------------------
// Function Tests
// ----------------------------------------------------------------------------

func TestFunctionDeclaration(t *testing.T) {
	expectNoErrors(t, `
		function add(a: number, b: number): number {
			return a + b;
		}
	`)
}

func TestFunctionCall(t *testing.T) {
	expectNoErrors(t, `
		function add(a: number, b: number): number {
			return a + b;
		}
		let x: number = add(1, 2);
	`)

	expectError(t, `
		function add(a: number, b: number): number {
			return a + b;
		}
		let x: number = add(1);
	`, "expected 2 arguments")

	expectError(t, `
		function add(a: number, b: number): number {
			return a + b;
		}
		let x: number = add("a", 2);
	`, "cannot pass string as number")
}

func TestReturnType(t *testing.T) {
	expectError(t, `
		function foo(): number {
			return "hello";
		}
	`, "cannot return string")

	expectError(t, `
		function foo(): number {
			return;
		}
	`, "missing return value")
}

func TestReturnOutsideFunction(t *testing.T) {
	expectError(t, `return 1;`, "return outside function")
}

// ----------------------------------------------------------------------------
// Array Tests
// ----------------------------------------------------------------------------

func TestArrayLiteral(t *testing.T) {
	expectNoErrors(t, `let arr: number[] = [1, 2, 3];`)
	expectNoErrors(t, `let arr: string[] = ["a", "b"];`)
}

func TestArrayIndexing(t *testing.T) {
	expectNoErrors(t, `
		let arr: number[] = [1, 2, 3];
		let x: number = arr[0];
	`)
	expectError(t, `
		let arr: number[] = [1, 2, 3];
		let x: number = arr["a"];
	`, "must be number")
}

func TestArrayAssignment(t *testing.T) {
	expectNoErrors(t, `
		let arr: number[] = [1, 2, 3];
		arr[0] = 42;
	`)
	expectError(t, `
		let arr: number[] = [1, 2, 3];
		arr[0] = "hello";
	`, "cannot assign string")
}

// ----------------------------------------------------------------------------
// Object Tests
// ----------------------------------------------------------------------------

func TestObjectLiteral(t *testing.T) {
	expectNoErrors(t, `
		let point: { x: number, y: number } = { x: 1, y: 2 };
	`)
}

func TestObjectPropertyAccess(t *testing.T) {
	expectNoErrors(t, `
		let point: { x: number, y: number } = { x: 1, y: 2 };
		let x: number = point.x;
	`)
	expectError(t, `
		let point: { x: number, y: number } = { x: 1, y: 2 };
		let z: number = point.z;
	`, "does not exist")
}

// ----------------------------------------------------------------------------
// Class Tests
// ----------------------------------------------------------------------------

func TestClassDeclaration(t *testing.T) {
	expectNoErrors(t, `
		class Point {
			x: number;
			y: number;

			constructor(x: number, y: number) {
				this.x = x;
				this.y = y;
			}

			add(other: Point): Point {
				return new Point(this.x + other.x, this.y + other.y);
			}
		}
	`)
}

func TestNewExpression(t *testing.T) {
	expectNoErrors(t, `
		class Point {
			x: number;
			y: number;
			constructor(x: number, y: number) {
				this.x = x;
				this.y = y;
			}
		}
		let p: Point = new Point(1, 2);
	`)

	expectError(t, `
		class Point {
			x: number;
			constructor(x: number) {
				this.x = x;
			}
		}
		let p: Point = new Point("hello");
	`, "cannot pass string as number")
}

func TestThisExpression(t *testing.T) {
	expectNoErrors(t, `
		class Counter {
			value: number;
			constructor() {
				this.value = 0;
			}
			get(): number {
				return this.value;
			}
		}
	`)

	expectError(t, `let x: number = this.value;`, "this' outside of class")
}

func TestInheritance(t *testing.T) {
	expectNoErrors(t, `
		class Animal {
			name: string;
			constructor(name: string) {
				this.name = name;
			}
		}

		class Dog extends Animal {
			breed: string;
			constructor(name: string, breed: string) {
				super(name);
				this.breed = breed;
			}
		}

		let dog: Animal = new Dog("Buddy", "Lab");
	`)
}

// ----------------------------------------------------------------------------
// Type Alias Tests
// ----------------------------------------------------------------------------

func TestTypeAlias(t *testing.T) {
	expectNoErrors(t, `
		type Point = { x: number, y: number };
		let p: Point = { x: 1, y: 2 };
	`)
}

// ----------------------------------------------------------------------------
// Null Safety Tests
// ----------------------------------------------------------------------------

func TestNullableType(t *testing.T) {
	expectNoErrors(t, `let x: string | null = null;`)
	expectNoErrors(t, `let x: string | null = "hello";`)
}

func TestNullTypeNarrowing(t *testing.T) {
	expectNoErrors(t, `
		let x: string | null = "hello";
		if (x != null) {
			let y: string = x;
		}
	`)
}

func TestNullPropertyAccess(t *testing.T) {
	expectError(t, `
		let obj: { x: number } | null = null;
		let x: number = obj.x;
	`, "cannot access property on potentially null")
}

// ----------------------------------------------------------------------------
// Function Type Tests
// ----------------------------------------------------------------------------

func TestFunctionExpression(t *testing.T) {
	expectNoErrors(t, `
		let add: (number, number) => number = function(a: number, b: number): number {
			return a + b;
		};
	`)
}

func TestHigherOrderFunction(t *testing.T) {
	expectNoErrors(t, `
		function apply(f: (number) => number, x: number): number {
			return f(x);
		}

		let double: (number) => number = function(x: number): number {
			return x * 2;
		};

		let result: number = apply(double, 5);
	`)
}
