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
	expectNoErrors(t, `let x: int = 42;`)
	expectNoErrors(t, `let x: float = 3.14;`)
	expectError(t, `let x: string = 42;`, "cannot assign int to string")
}

func TestStringLiteral(t *testing.T) {
	expectNoErrors(t, `let x: string = "hello";`)
	expectError(t, `let x: int = "hello";`, "cannot assign string to int")
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
	expectNoErrors(t, `let x: int = 1 + 2;`)
	expectNoErrors(t, `let x: int = 3 - 1;`)
	expectNoErrors(t, `let x: int = 2 * 3;`)
	expectNoErrors(t, `let x: float = 6 / 2;`)
	expectNoErrors(t, `let x: int = 7 % 3;`)

	// String concatenation
	expectNoErrors(t, `let x: string = "hello" + " world";`)

	// Type errors in arithmetic
	expectError(t, `let x: int = "a" - 1;`, "requires numbers")
	expectError(t, `let x: int = true * 2;`, "requires numbers")

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
	expectNoErrors(t, `let x: int = -5;`)
	expectNoErrors(t, `let x: boolean = !true;`)
	expectError(t, `let x: int = -"hello";`, "requires number")
	expectError(t, `let x: boolean = !42;`, "requires boolean")
}

// ----------------------------------------------------------------------------
// Variable Tests
// ----------------------------------------------------------------------------

func TestVariableDeclaration(t *testing.T) {
	expectNoErrors(t, `let x: int = 1;`)
	expectNoErrors(t, `const PI: float = 3.14;`)
	expectError(t, `let x: int = "hello";`, "cannot assign string to int")
}

func TestVariableUsage(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 1;
		let y: int = x + 2;
	`)
	expectError(t, `let x: int = y;`, "undefined variable")
}

func TestVariableScoping(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 1;
		{
			let y: int = x + 1;
		}
	`)
}

// ----------------------------------------------------------------------------
// Control Flow Tests
// ----------------------------------------------------------------------------

func TestIfStatement(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 1;
		if (x > 0) {
			let y: int = 2;
		}
	`)
	expectError(t, `
		if (42) {
			let x: int = 1;
		}
	`, "must be boolean")
}

func TestWhileStatement(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 0;
		while (x < 10) {
			x = x + 1;
		}
	`)
	expectError(t, `
		while ("true") {
			let x: int = 1;
		}
	`, "must be boolean")
}

func TestForStatement(t *testing.T) {
	expectNoErrors(t, `
		let sum: int = 0;
		let i: int = 0;
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
		function add(a: int, b: int): int {
			return a + b;
		}
	`)
}

func TestFunctionCall(t *testing.T) {
	expectNoErrors(t, `
		function add(a: int, b: int): int {
			return a + b;
		}
		let x: int = add(1, 2);
	`)

	expectError(t, `
		function add(a: int, b: int): int {
			return a + b;
		}
		let x: int = add(1);
	`, "expected 2 arguments")

	expectError(t, `
		function add(a: int, b: int): int {
			return a + b;
		}
		let x: int = add("a", 2);
	`, "cannot pass string as int")
}

func TestReturnType(t *testing.T) {
	expectError(t, `
		function foo(): int {
			return "hello";
		}
	`, "cannot return string")

	expectError(t, `
		function foo(): int {
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
	expectNoErrors(t, `let arr: int[] = [1, 2, 3];`)
	expectNoErrors(t, `let arr: string[] = ["a", "b"];`)
}

func TestArrayIndexing(t *testing.T) {
	expectNoErrors(t, `
		let arr: int[] = [1, 2, 3];
		let x: int = arr[0];
	`)
	expectError(t, `
		let arr: int[] = [1, 2, 3];
		let x: int = arr["a"];
	`, "must be int")
}

func TestArrayAssignment(t *testing.T) {
	expectNoErrors(t, `
		let arr: int[] = [1, 2, 3];
		arr[0] = 42;
	`)
	expectError(t, `
		let arr: int[] = [1, 2, 3];
		arr[0] = "hello";
	`, "cannot assign string")
}

// ----------------------------------------------------------------------------
// Object Tests
// ----------------------------------------------------------------------------

func TestObjectLiteral(t *testing.T) {
	expectNoErrors(t, `
		let point: { x: int, y: int } = { x: 1, y: 2 };
	`)
}

func TestObjectPropertyAccess(t *testing.T) {
	expectNoErrors(t, `
		let point: { x: int, y: int } = { x: 1, y: 2 };
		let x: int = point.x;
	`)
	expectError(t, `
		let point: { x: int, y: int } = { x: 1, y: 2 };
		let z: int = point.z;
	`, "does not exist")
}

// ----------------------------------------------------------------------------
// Class Tests
// ----------------------------------------------------------------------------

func TestClassDeclaration(t *testing.T) {
	expectNoErrors(t, `
		class Point {
			x: int;
			y: int;

			constructor(x: int, y: int) {
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
			x: int;
			y: int;
			constructor(x: int, y: int) {
				this.x = x;
				this.y = y;
			}
		}
		let p: Point = new Point(1, 2);
	`)

	expectError(t, `
		class Point {
			x: int;
			constructor(x: int) {
				this.x = x;
			}
		}
		let p: Point = new Point("hello");
	`, "cannot pass string as int")
}

func TestThisExpression(t *testing.T) {
	expectNoErrors(t, `
		class Counter {
			value: int;
			constructor() {
				this.value = 0;
			}
			get(): int {
				return this.value;
			}
		}
	`)

	expectError(t, `let x: int = this.value;`, "this' outside of class")
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
		type Point = { x: int, y: int };
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
		let obj: { x: int } | null = null;
		let x: int = obj.x;
	`, "cannot access property on potentially null")
}

// ----------------------------------------------------------------------------
// Function Type Tests
// ----------------------------------------------------------------------------

func TestFunctionExpression(t *testing.T) {
	expectNoErrors(t, `
		let add: (int, int) => int = function(a: int, b: int): int {
			return a + b;
		};
	`)
}

func TestHigherOrderFunction(t *testing.T) {
	expectNoErrors(t, `
		function apply(f: (int) => int, x: int): int {
			return f(x);
		}

		let double: (int) => int = function(x: int): int {
			return x * 2;
		};

		let result: int = apply(double, 5);
	`)
}

// ============================================================
// Type Inference Tests
// ============================================================

func TestTypeInferenceNumber(t *testing.T) {
	expectNoErrors(t, `let x = 10;`)
}

func TestTypeInferenceString(t *testing.T) {
	expectNoErrors(t, `let name = "hello";`)
}

func TestTypeInferenceBoolean(t *testing.T) {
	expectNoErrors(t, `let flag = true;`)
}

func TestTypeInferenceArray(t *testing.T) {
	expectNoErrors(t, `let arr = [1, 2, 3];`)
}

func TestTypeInferenceObject(t *testing.T) {
	expectNoErrors(t, `let obj = { x: 1, y: 2 };`)
}

func TestTypeInferenceReassignment(t *testing.T) {
	// Inferred type should be enforced on reassignment
	expectError(t, `
		let x = 10;
		x = "hello";
	`, "cannot assign")
}

func TestTypeInferenceWithExplicitType(t *testing.T) {
	expectNoErrors(t, `let x: int = 10;`)
}

// ============================================================
// Compound Assignment Tests
// ============================================================

func TestCompoundAssignmentNumber(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 10;
		x += 5;
		x -= 3;
		x *= 2;
		x /= 4;
		x %= 3;
	`)
}

func TestCompoundAssignmentString(t *testing.T) {
	expectNoErrors(t, `
		let s: string = "hello";
		s += " world";
	`)
}

func TestCompoundAssignmentTypeError(t *testing.T) {
	expectError(t, `
		let x: int = 10;
		x += "hello";
	`, "requires numbers")
}

// ============================================================
// Arrow Function Tests
// ============================================================

func TestArrowFunctionExpression(t *testing.T) {
	expectNoErrors(t, `
		let add = (a: int, b: int): int => a + b;
	`)
}

func TestArrowFunctionBlock(t *testing.T) {
	expectNoErrors(t, `
		let double = (x: int): int => {
			return x * 2;
		};
	`)
}

func TestArrowFunctionNoParams(t *testing.T) {
	expectNoErrors(t, `
		let getZero = (): int => 0;
	`)
}

// ============================================================
// Increment/Decrement Tests
// ============================================================

func TestIncrementDecrement(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 10;
		x++;
		x--;
		++x;
		--x;
	`)
}

func TestIncrementDecrementTypeError(t *testing.T) {
	expectError(t, `
		let s: string = "hello";
		s++;
	`, "requires number")
}

// ============================================================
// For-of Loop Tests
// ============================================================

func TestForOfArray(t *testing.T) {
	expectNoErrors(t, `
		let arr: int[] = [1, 2, 3];
		for (let item of arr) {
			let x: int = item;
		}
	`)
}

func TestForOfString(t *testing.T) {
	expectNoErrors(t, `
		let s: string = "hello";
		for (let char of s) {
			let c: string = char;
		}
	`)
}

// ============================================================
// Switch Statement Tests
// ============================================================

func TestSwitchStatement(t *testing.T) {
	expectNoErrors(t, `
		let x: int = 1;
		switch (x) {
			case 1:
				let a: int = 1;
				break;
			case 2:
				let b: int = 2;
				break;
			default:
				let c: int = 0;
		}
	`)
}

// ============================================================
// Nullish Coalescing Tests
// ============================================================

func TestNullishCoalescing(t *testing.T) {
	expectNoErrors(t, `
		let x: int | null = null;
		let y: int = x ?? 0;
	`)
}

// ============================================================
// Const Validation Tests
// ============================================================

func TestConstReassignment(t *testing.T) {
	// Direct reassignment to const should fail
	expectError(t, `
		const x: int = 10;
		x = 20;
	`, "cannot assign to const")

	// Reassignment to let should work
	expectNoErrors(t, `
		let x: int = 10;
		x = 20;
	`)
}

func TestConstCompoundAssign(t *testing.T) {
	// Compound assignment to const should fail
	expectError(t, `
		const x: int = 10;
		x += 5;
	`, "cannot assign to const")

	expectError(t, `
		const x: int = 10;
		x -= 5;
	`, "cannot assign to const")

	expectError(t, `
		const x: int = 10;
		x *= 2;
	`, "cannot assign to const")
}

func TestConstIncrement(t *testing.T) {
	// Increment/decrement on const should fail
	expectError(t, `
		const x: int = 10;
		x++;
	`, "cannot assign to const")

	expectError(t, `
		const x: int = 10;
		++x;
	`, "cannot assign to const")

	expectError(t, `
		const x: int = 10;
		x--;
	`, "cannot assign to const")

	expectError(t, `
		const x: int = 10;
		--x;
	`, "cannot assign to const")
}
