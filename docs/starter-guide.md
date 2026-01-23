# GoTS Starter Guide

A quick introduction to GoTS, a statically-typed TypeScript subset that compiles to bytecode and runs on a Go-based virtual machine.

---

## What is GoTS?

GoTS is a minimal, statically-typed programming language with TypeScript-like syntax. It compiles source code to bytecode which runs on a custom VM written in Go.

**Key features:**
- TypeScript-compatible syntax (GoTS is a valid TypeScript subset)
- Static type checking at compile time
- First-class functions and closures
- Classes with inheritance
- Garbage collection
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

### Compile to Bytecode

```bash
gots compile program.gts              # Creates program.gtsb
gots compile program.gts output.gtsb  # Custom output name
```

### Execute Bytecode

```bash
gots exec program.gtsb
```

### Disassemble Bytecode

```bash
gots disasm program.gtsb
gots disasm program.gts    # Compiles first, then disassembles
```

### Interactive REPL

```bash
gots repl
```

```
GoTS REPL v0.1.0
Type 'exit' or press Ctrl+D to quit

>>> let x: number = 42
>>> println(x * 2)
84
>>> exit
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
let x: number = 10;
let name: string = "Alice";
let active: boolean = true;
const PI: number = 3.14159;
```

### Functions

```typescript
function add(a: number, b: number): number {
    return a + b;
}

function greet(name: string): void {
    println("Hello, " + name);
}

let result: number = add(5, 3);
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
let i: number = 0;
while (i < 5) {
    println(i);
    i = i + 1;
}

// For loop
for (let j: number = 0; j < 5; j = j + 1) {
    println(j);
}
```

### Arrays

```typescript
let numbers: number[] = [1, 2, 3, 4, 5];
println(numbers[0]);     // 1
println(len(numbers));   // 5

push(numbers, 6);        // Append to array
let last: number | null = pop(numbers);  // Remove last element
```

### Objects

```typescript
type Point = { x: number, y: number };

let origin: Point = { x: 0, y: 0 };
println(origin.x);
origin.y = 10;
```

### Classes

```typescript
class Counter {
    count: number;

    constructor() {
        this.count = 0;
    }

    increment(): void {
        this.count = this.count + 1;
    }

    getCount(): number {
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
function makeCounter(): () => number {
    let count: number = 0;
    return function(): number {
        count = count + 1;
        return count;
    };
}

let counter: () => number = makeCounter();
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

| Function    | Description                        |
|-------------|------------------------------------|
| `print(v)`  | Print value without newline        |
| `println(v)`| Print value with newline           |
| `len(s)`    | Get string or array length         |
| `toString(n)` | Convert number to string         |
| `toNumber(s)` | Parse string as number           |
| `push(arr, v)` | Append value to array           |
| `pop(arr)`  | Remove and return last element     |
| `sqrt(n)`   | Square root                        |
| `floor(n)`  | Round down to integer              |
| `ceil(n)`   | Round up to integer                |
| `abs(n)`    | Absolute value                     |

---

## Example Program

```typescript
// FizzBuzz in GoTS

function fizzbuzz(n: number): void {
    for (let i: number = 1; i <= n; i = i + 1) {
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

## File Extensions

| Extension | Description                     |
|-----------|---------------------------------|
| `.gts`    | GoTS source file                |
| `.gtsb`   | Compiled bytecode               |
| `.gtsb.gz`| Compressed bytecode             |

---

## Further Reading

- [Language Specification](language-spec-v1.md) - Complete language reference
- [Bytecode Specification](bytecode-spec-v1.md) - VM internals and opcodes
- [Implementation Plan](implementation-plan.md) - Project architecture

---

## Quick Reference

```typescript
// Types
number, string, boolean, null, void
number[], string[][]           // Arrays
{ x: number, y: number }       // Object types
(a: number) => string          // Function types
string | null                  // Nullable types

// Operators
+ - * / %                      // Arithmetic
== != < > <= >=                // Comparison
&& || !                        // Logical

// Keywords
let, const, function, return
if, else, while, for, break, continue
class, constructor, this, new, extends, super
type, true, false, null
```
