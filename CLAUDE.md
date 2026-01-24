# GoTS (GoTypeScript)

A TypeScript-like language that compiles to Go.

## Project Overview

GoTS implements a complete compilation pipeline: lexer → parser → type checker → Go code generator. It transpiles TypeScript-like source code to Go, providing static typing, functions, closures, classes with inheritance, and access to the Go ecosystem.

## Important Notes

**`../scaffold/` is for implementation reference only.** It contains:
- `quickjs/` - QuickJS JavaScript engine (C implementation reference)
- `typescript-go-main/` - TypeScript compiler in Go (reference for TypeScript semantics)

Do not modify files in scaffold/. Use them only as reference for understanding how language features should be implemented.

## Architecture

```
Source (.gts) → Lexer → Parser → TypedAST Builder → Go Code Generator → go build → Native Binary
```

### Package Structure

| Package | Purpose |
|---------|---------|
| `pkg/token` | Token type definitions (85 token types) |
| `pkg/lexer` | Tokenization with line/column tracking |
| `pkg/ast` | AST node definitions (statements, expressions, types) |
| `pkg/parser` | Pratt parser with operator precedence |
| `pkg/types` | Type definitions and utilities |
| `pkg/typed` | Type-annotated AST with builder |
| `pkg/codegen` | Go source code generator |
| `cmd/gots` | CLI entry point |

## Build & Run

```bash
# Build
go build -o gots ./cmd/gots

# Run tests
go test ./...

# Run with verbose
go test -v ./pkg/...
```

## CLI Commands

```bash
gots run program.gts              # Compile and run
gots build program.gts            # Compile to native binary
gots build program.gts -o myapp   # Specify output name
gots build program.gts --emit-go  # Output Go source instead
gots emit-go program.gts          # Generate Go source code
gots repl                         # Interactive REPL
```

## Language Features

### Types
- Primitives: `int`, `float`, `string`, `boolean`, `void`, `null`
- Arrays: `int[]`, `float[]`, `string[]`
- Objects: `{x: int, y: string}`
- Functions: `(a: int) => string`, `Function` (dynamic)
- Nullable: `string | null`
- Type aliases: `type Point = {x: int, y: int}`

### Type Mapping to Go

| GTS Type | Go Type |
|----------|---------|
| `int` | `int` |
| `float` | `float64` |
| `string` | `string` |
| `boolean` | `bool` |
| `void` | (no return) |
| `null` | `nil` / `interface{}` |
| `Function` | `interface{}` |
| `T[]` | `[]T` |
| `T \| null` | `*T` |
| `class C` | `*C` (struct pointer) |

### Numeric Type Rules
- Integer literals (e.g., `42`) have type `int`
- Decimal literals (e.g., `3.14`) have type `float`
- `int + int = int`, `int + float = float`, `float + float = float`
- Division (`/`) always returns `float`
- Modulo (`%`) requires `int` operands
- Array indexing requires `int`
- `len()` returns `int`

### Syntax Examples

```typescript
// Variables
let x: int = 42
let pi: float = 3.14159
const name: string = "GoTS"

// Functions
function factorial(n: int): int {
    if (n <= 1) { return 1 }
    return n * factorial(n - 1)
}

// Higher-order functions
function curry_add(a: int): Function {
    return function(b: int): int {
        return a + b
    }
}

// Classes
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
    constructor(name: string) {
        super(name)
    }
    speak(): void {
        println(this.name + " barks")
    }
}

// Arrays
let arr: int[] = [1, 2, 3]
push(arr, 4)
let last: int = pop(arr)
```

### Built-in Functions
`println`, `print`, `len`, `push`, `pop`, `typeof`, `tostring`, `toint`, `tofloat`, `sqrt`, `floor`, `ceil`, `abs`

## Key Implementation Details

### Parser
- Pratt parsing with 10 precedence levels
- Two-token lookahead
- Error recovery (collects multiple errors)

### TypedAST Builder
- Transforms AST to type-annotated AST
- Performs type checking during transformation
- Tracks scope and variable types
- Collects closure capture information
- Allows `any` type in arithmetic for dynamic typing support

### Code Generator
- Generates idiomatic Go code
- Maps GTS types to Go types
- Classes become Go structs with methods
- Closures map directly to Go closures
- Runtime helpers for dynamic operations (`gts_call`, `gts_toint`, etc.)
- Automatic type assertions for `any` operands in arithmetic

## Code Conventions

- Exported Go names are capitalized (e.g., `count` → `Count`)
- Constructor functions are named `NewClassName`
- Method receivers use `this` pointer
- Go reserved words get `_` suffix

## Testing

Each package has corresponding `*_test.go` files. Example programs in `test/` directory.

```bash
# Run specific package tests
go test -v ./pkg/lexer
go test -v ./pkg/parser
go test -v ./pkg/types
go test -v ./pkg/codegen

# Run example programs
gots run test/example.gts
gots run test/higher_order.gts
gots run test/y_combinator.gts
```

## Development Workflow

1. Make changes to source
2. Run `go test ./pkg/...` to verify
3. Test with example: `go run ./cmd/gots run test/example.gts`
4. Check generated Go: `go run ./cmd/gots emit-go test/example.gts`

## Generated Code Example

Input (`test.gts`):
```typescript
function add(a: int, b: int): int {
    return a + b
}
println(add(2, 3))
```

Output (generated Go):
```go
package main

import "fmt"

func add(a int, b int) int {
    return (a + b)
}

func main() {
    fmt.Println(add(2, 3))
}
```
