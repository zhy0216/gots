# goTS

A TypeScript-like language that compiles to Go. Write TypeScript-flavored code, get a native Go binary.

> Full documentation: [docs/](../docs/) — rendered site available via `make docs-serve` (VitePress).

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

Requires Go 1.25+.

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
gots build program.gts --emit-go  # Output Go source instead
gots emit-go program.gts          # Generate Go source code
gots repl                         # Interactive REPL
gots version                      # Show version
```

## Language Features

- **Types**: `number`, `int`, `float`, `string`, `boolean`, `void`, `null`, `any`, arrays, object types, function types, tuples, unions, intersections, literal types, `Map<K,V>`, `Set<T>`, `Date`, `RegExp`
- **Generics**: generic functions and classes with inference, explicit type arguments, and default type parameters
- **Functions**: declarations, expressions, arrows, closures, higher-order functions
- **Classes**: inheritance, `super`, method overriding, decorators
- **Interfaces**: structural typing
- **Async**: `async`/`await`, `Promise<T>`, `setTimeout`, `setInterval`, `queueMicrotask` (Go-scheduler-based event loop)
- **Control flow**: `if`/`else`, `while`, `for`, `for-of`, `switch`, `try`/`catch`/`finally`/`throw`
- **Modern syntax**: template literals, destructuring, spread, optional chaining (`?.`), nullish coalescing (`??`)
- **Modules**: named/default/namespace imports, re-exports
- **Go interop**: `import { ... } from "go:strings"` pulls Go stdlib functions; built-in SQL support via `connect()`
- **Standard library**: `Math`, `Number`, `JSON`, `Date`, `Object`, plus rich `String`, `Array`, `Map`, `Set`, `RegExp` methods
- **Type inference**: numeric literals default to `number`; `int`/`float` are opt-in for Go interop

### Examples

```typescript
// Functions and closures
function makeAdder(x: int): Function {
    return function(y: int): int {
        return x + y
    }
}
let add5 = makeAdder(5)
println(add5(10))  // 15

// Generics
function identity<T>(x: T): T {
    return x
}
let n: int = identity(42)

class Box<T> {
    value: T
    constructor(v: T) { this.value = v }
    get(): T { return this.value }
}

// Interfaces (structural typing)
interface Drawable {
    draw(): void
}
class Circle {
    radius: float
    constructor(radius: float) { this.radius = radius }
    draw(): void { println("circle") }
}
let shape: Drawable = new Circle(5.0)

// Async / event loop
async function compute(): Promise<int> {
    let a: int = await fetchValue()
    return a * 2
}
setTimeout(function(): void { println("later") }, 10)

// SQL
const db = connect("./test.db")
db.sql`CREATE TABLE IF NOT EXISTS users (name TEXT NOT NULL)`
const users = db.sql<User[]>`SELECT name FROM users`

// Go interop
import { ToUpper, Split } from "go:strings"
let upper: string = ToUpper("hello")

// Optional chaining + nullish coalescing
let name: string | null = null
println(name ?? "default")
```

### Numeric Type Rules

- Numeric literals (e.g., `42`, `3.14`) default to `number` type
- `int` and `float` are opt-in explicit annotations for Go interop
- Numeric literals can be assigned to any numeric type (`int`, `float`, `number`)
- `int` → `number`: allowed (widening); `number` → `int`: use `toint()`
- Division (`/`) always returns `float`; modulo (`%`) requires `int` operands
- `len()` returns `int`

## Built-in Functions

| Function | Description |
|----------|-------------|
| `println(x)` / `print(x)` | Print with/without newline |
| `len(arr)` | Array/string length |
| `push(arr, x)` / `pop(arr)` | Append / remove last element |
| `typeof(x)` | Get type as string |
| `tostring(x)` / `toint(x)` / `tofloat(x)` | Type conversion |
| `parseInt(x)` / `parseFloat(x)` | String to number |
| `sqrt(x)` / `floor(x)` / `ceil(x)` / `abs(x)` | Basic math |
| `isNaN(x)` / `isFinite(x)` | Number checks |
| `setTimeout(fn, ms)` / `setInterval(fn, ms)` / `clearTimeout(id)` / `clearInterval(id)` | Timers |
| `queueMicrotask(fn)` | Microtask queue |
| `connect(path)` | Open SQLite database |

See [docs/built-in-reference.md](../docs/built-in-reference.md) for the full reference (`Math`, `Number`, `JSON`, `Date`, `Object`, `String`, `Array`, `Map`, `Set`, `RegExp`).

## Type Mapping to Go

| GTS Type | Go Type |
|----------|---------|
| `number` | `float64` |
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
| `RegExp` | `*regexp.Regexp` |
| `Map<K,V>` / `Set<T>` | `map[K]V` / `map[T]struct{}` |
| `Date` | `time.Time` |
| `Promise<T>` | `GTS_Promise[T]` (generated) |

## Architecture

```
Source (.gts) → Lexer → Parser → TypedAST Builder → Go Code Generator → go build → Native Binary
```

| Package | Purpose |
|---------|---------|
| `pkg/token` | Token definitions |
| `pkg/lexer` | Tokenization with line/column tracking |
| `pkg/ast` | AST node definitions |
| `pkg/parser` | Pratt parser with error recovery |
| `pkg/types` | Type definitions and utilities |
| `pkg/typed` | Type-annotated AST + type checker |
| `pkg/codegen` | Go source code generator |
| `pkg/declaration` | Go stdlib `.d.gts` declaration loader |
| `pkg/module` | Local module loader |
| `cmd/gots` | CLI entry point |

## Documentation

The full documentation lives in [docs/](../docs/) and is rendered with VitePress:

```bash
make docs-install   # npm install (once)
make docs-serve     # start dev server at http://localhost:5173
make docs-build     # build static site to ../docs/.vitepress/dist
```

## Development

```bash
make test            # Run all tests
make test-unit       # Run unit tests only
make test-integration # Run integration tests
make test-test262    # Run test262 conformance tests
make check           # Run fmt, vet, and tests
make clean           # Clean build artifacts
```

## License

MIT
