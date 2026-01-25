# goTS

A TypeScript-like language that compiles to Go.

## Quick Start

```bash
# Build
make build

# Run a program
./gots run test/example.gts

# Compile to native binary
./gots build test/example.gts -o myapp

# Start REPL
./gots repl
```

## Installation

```bash
# Install to $GOPATH/bin
make install

# Or build locally
make build
```

## Usage

```bash
gots run program.gts              # Compile and run
gots build program.gts            # Compile to native binary
gots build program.gts -o myapp   # Specify output name
gots emit-go program.gts          # Output generated Go code
gots repl                         # Interactive REPL
```

## Language Features

### Variables and Types

```typescript
let x: int = 42
let pi: float = 3.14
const name: string = "goTS"
let flag: boolean = true
```

### Functions

```typescript
function add(a: int, b: int): int {
    return a + b
}

// Higher-order functions
function apply(f: Function, x: int): int {
    return f(x)
}
```

### Classes

```typescript
class Animal {
    name: string

    constructor(name: string) {
        this.name = name
    }

    speak(): void {
        println(this.name)
    }
}

class Dog extends Animal {
    speak(): void {
        println(this.name + " barks")
    }
}

let dog = new Dog("Rex")
dog.speak()
```

### Arrays

```typescript
let arr: int[] = [1, 2, 3]
push(arr, 4)
let last: int = pop(arr)
println(len(arr))
```

### Control Flow

```typescript
if (x > 0) {
    println("positive")
} else {
    println("non-positive")
}

for (let i: int = 0; i < 10; i = i + 1) {
    println(i)
}

for (let item of items) {
    println(item)
}
```

### Type Aliases

```typescript
type Point = { x: int, y: int }
let p: Point = { x: 10, y: 20 }
```

## Built-in Functions

| Function | Description |
|----------|-------------|
| `println(x)` | Print with newline |
| `print(x)` | Print without newline |
| `len(arr)` | Array/string length |
| `push(arr, x)` | Append to array |
| `pop(arr)` | Remove and return last element |
| `typeof(x)` | Get type as string |
| `tostring(x)` | Convert to string |
| `toint(x)` | Convert to int |
| `tofloat(x)` | Convert to float |
| `sqrt(x)` | Square root |
| `floor(x)` | Floor |
| `ceil(x)` | Ceiling |
| `abs(x)` | Absolute value |

## Development

```bash
make test            # Run all tests
make test-unit       # Run unit tests only
make test-integration # Run integration tests
make check           # Run fmt, vet, and tests
make clean           # Clean build artifacts
```

## License

MIT
