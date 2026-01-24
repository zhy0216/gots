# goTS Starter Guide

A quick introduction to goTS, a statically-typed TypeScript subset that compiles to Go.

---

## What is goTS?

goTS is a minimal, statically-typed programming language with TypeScript-like syntax. It transpiles source code to Go, which is then compiled to native binaries.

**Key features:**
- TypeScript-compatible syntax (goTS is a valid TypeScript subset)
- Static type checking at compile time
- First-class functions and closures
- Classes with inheritance
- Compiles to native binaries via Go
- Interactive REPL

---

## Installation

### Building from Source

```bash
cd gots
go build -o gots ./cmd/gots
```

This creates the `gots` executable in the current directory.

---

## CLI Usage

### Run a Program

```bash
gots run program.gts
```

Or simply:

```bash
gots program.gts
```

### Build a Native Binary

```bash
gots build program.gts              # Creates program binary
gots build program.gts -o myapp     # Custom output name
```

### Generate Go Source Code

```bash
gots emit-go program.gts            # Creates program.go
gots emit-go program.gts output.go  # Custom output name
gots build program.gts --emit-go    # Alternative via build command
```

### Interactive REPL

```bash
gots repl
```

```
goTS REPL v0.2.0 (Go transpiler)
Type 'exit' or press Ctrl+D to quit

>>> let x: int = 42
>>> println(x * 2)
84
>>> exit
```

### Other Commands

```bash
gots version    # Show version
gots help       # Show help
```

---

## Language Basics

### Hello World

```typescript
println("Hello, World!");
```

### Variables

All variables require explicit type annotations.

```typescript
let x: int = 10;
let pi: float = 3.14159;
let name: string = "Alice";
let active: boolean = true;
const MAX: int = 100;
```

### Numeric Types

goTS has two numeric types: `int` and `float`.

```typescript
let count: int = 42;        // Integer
let price: float = 19.99;   // Floating point

// Type rules:
// int + int = int
// int + float = float
// float + float = float
// Division (/) always returns float
// Modulo (%) requires int operands
```

### Functions

```typescript
function add(a: int, b: int): int {
    return a + b;
}

function greet(name: string): void {
    println("Hello, " + name);
}

let result: int = add(5, 3);
greet("World");
```

### Control Flow

```typescript
// If statement
if (x > 10) {
    println("big");
} else {
    println("small");
}

// While loop
let i: int = 0;
while (i < 5) {
    println(i);
    i = i + 1;
}

// For loop
for (let j: int = 0; j < 5; j = j + 1) {
    println(j);
}
```

### Arrays

```typescript
let numbers: int[] = [1, 2, 3, 4, 5];
println(numbers[0]);     // 1
println(len(numbers));   // 5

push(numbers, 6);        // Append to array
let last: int = pop(numbers);  // Remove last element
```

### Objects

```typescript
type Point = { x: int, y: int };

let origin: Point = { x: 0, y: 0 };
println(origin.x);
origin.y = 10;
```

### Classes

```typescript
class Counter {
    count: int;

    constructor() {
        this.count = 0;
    }

    increment(): void {
        this.count = this.count + 1;
    }

    getCount(): int {
        return this.count;
    }
}

let c: Counter = new Counter();
c.increment();
c.increment();
println(c.getCount());  // 2
```

### Inheritance

```typescript
class Animal {
    name: string;

    constructor(name: string) {
        this.name = name;
    }

    speak(): void {
        println(this.name + " makes a sound");
    }
}

class Dog extends Animal {
    constructor(name: string) {
        super(name);
    }

    speak(): void {
        println(this.name + " barks");
    }
}

let dog: Dog = new Dog("Rex");
dog.speak();  // "Rex barks"
```

### Closures

```typescript
function makeCounter(): Function {
    let count: int = 0;
    return function(): int {
        count = count + 1;
        return count;
    };
}

let counter: Function = makeCounter();
println(counter());  // 1
println(counter());  // 2
println(counter());  // 3
```

### Nullable Types

```typescript
let name: string | null = null;

if (name != null) {
    println(name);  // Type narrowing: name is string here
}
```

---

## Built-in Functions

| Function      | Description                        |
|---------------|------------------------------------|
| `print(v)`    | Print value without newline        |
| `println(v)`  | Print value with newline           |
| `len(s)`      | Get string or array length         |
| `tostring(v)` | Convert value to string            |
| `toint(v)`    | Convert value to int               |
| `tofloat(v)`  | Convert value to float             |
| `push(arr, v)`| Append value to array              |
| `pop(arr)`    | Remove and return last element     |
| `typeof(v)`   | Get type name as string            |
| `sqrt(n)`     | Square root                        |
| `floor(n)`    | Round down to integer              |
| `ceil(n)`     | Round up to integer                |
| `abs(n)`      | Absolute value                     |

---

## Example Program

```typescript
// FizzBuzz in goTS

function fizzbuzz(n: int): void {
    for (let i: int = 1; i <= n; i = i + 1) {
        if (i % 15 == 0) {
            println("FizzBuzz");
        } else if (i % 3 == 0) {
            println("Fizz");
        } else if (i % 5 == 0) {
            println("Buzz");
        } else {
            println(i);
        }
    }
}

fizzbuzz(20);
```

---

## File Extension

| Extension | Description                     |
|-----------|---------------------------------|
| `.gts`    | goTS source file                |

---

## Further Reading

- [Language Specification](language-spec-v1.md) - Complete language reference

---

## Quick Reference

```typescript
// Types
int, float, string, boolean, null, void
int[], string[][]                  // Arrays
{ x: int, y: float }               // Object types
(a: int) => string                 // Function types
Function                           // Dynamic function type
string | null                      // Nullable types

// Operators
+ - * / %                          // Arithmetic
== != < > <= >=                    // Comparison
&& || !                            // Logical

// Keywords
let, const, function, return
if, else, while, for, break, continue
class, constructor, this, new, extends, super
type, true, false, null
```
