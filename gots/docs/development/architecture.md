# Architecture

goTS is a classic ahead-of-time compiler with four stages:

```
Source (.gts) → Lexer → Parser → TypedAST Builder → Go Code Generator → go build → Native Binary
```

Each stage is a separate Go package under `pkg/`, with `cmd/gots` as the CLI driver.

## Package Layout

| Package | Purpose |
|---------|---------|
| `pkg/token` | Token type definitions (keywords, operators, literals) |
| `pkg/lexer` | Tokenization with line/column tracking for error reporting |
| `pkg/ast` | Abstract syntax tree node definitions |
| `pkg/parser` | Pratt parser (precedence-climbing) with error recovery |
| `pkg/types` | Type definitions, assignability and unification logic |
| `pkg/typed` | Type-annotated AST builder — performs type checking during transformation |
| `pkg/codegen` | Go source generator (idiomatic Go, runtime helpers, SQL runtime) |
| `pkg/declaration` | Go stdlib `.d.gts` declaration loader (`go:strings`, `go:math`, ...) |
| `pkg/module` | Local module loader for `import ... from "./file"` |
| `cmd/gots` | CLI: `run`, `build`, `emit-go`, `repl` |

## The Pipeline in Detail

### 1. Lexer (`pkg/lexer`)

Produces a token stream with position info. Handles numeric literals (int/float), strings (double/single/backtick), template literals with `${...}` interpolation, regex literals (`/pattern/flags`), comments, and all keywords.

### 2. Parser (`pkg/parser`)

A Pratt parser (precedence-climbing) with two-token lookahead. Builds the untyped AST in `pkg/ast`. Collects multiple parse errors rather than failing fast, so users see every syntax problem at once.

### 3. TypedAST Builder (`pkg/typed`)

The heart of the type system. Transforms the AST into a type-annotated AST (`pkg/typed`), performing type checking on the way:

- **Multi-pass declaration collection**: type aliases, classes, interfaces, and imports are collected before function bodies are checked, allowing forward references.
- **Scope tracking**: variables, closures, and captured variables.
- **Type inference**: numeric literals default to `number`; return types and `let` bindings are inferred.
- **Structural typing**: interfaces are satisfied structurally.
- **Generics**: type parameter unification with inference and explicit type arguments.
- **Built-in registry**: signatures for global functions (`println`, `len`, ...), `Math`, `Number`, `JSON`, `Date`, `Object`, `Map`, `Set`, `RegExp`, `Promise`, string/array methods, timers, and `connect()`.
- **Go import resolution**: `import { X } from "go:pkg"` resolves through `pkg/declaration/stdlib/*.d.gts`.

Errors carry line/column positions and are aggregated before reporting.

### 4. Go Code Generator (`pkg/codegen`)

Emits idiomatic Go:

- **Naming**: exported names are capitalized (`count` → `Count`); Go reserved words get a `_` suffix; constructors become `NewClassName`; methods use a `this` receiver.
- **Classes** become structs with methods; inheritance flattens parent fields and methods.
- **Enums** become `type X int` + const blocks.
- **Tuples** become anonymous structs with `T0`, `T1` fields.
- **Unions/intersections** map to `interface{}` or merged structs; nullable types become pointers where possible.
- **Closures** map to Go closures with capture tracking.
- **Async** functions compile to `GTS_Promise[T]` values with goroutine-based scheduling.
- **SQL tagged templates** compile to parameterized queries via the generated SQL runtime (`pkg/codegen/sql.go`).
- **Runtime helpers** are emitted only when used (`gts_call`, `gts_toint`, `gts_setTimeout`, ...).

### 5. The Go Toolchain

The emitted code is self-contained Go (plus optional `modernc.org/sqlite` for SQL). The CLI writes it to a temporary module, runs `go run` / `go build`, and produces the final binary.

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

## Where Things Live

| Feature | Builder (type check) | Codegen (emit) |
|---------|----------------------|----------------|
| Generics | `pkg/typed/builder.go` — type parameter unification | `pkg/codegen/codegen.go` — monomorphized Go funcs |
| Async/await | `pkg/typed/builder.go` — `Promise<T>` typing | `pkg/codegen/codegen.go` — `GTS_Promise[T]` + event loop runtime |
| SQL | `pkg/typed/builder.go` — `SQLDatabase`/`SQLTransaction` types | `pkg/codegen/sql.go` — sqlite runtime + query generation |
| Go imports | `pkg/declaration/loader.go` — `.d.gts` resolution | `pkg/codegen/codegen.go` — direct Go call emission |
| Modules | `pkg/module/loader.go` | inlined/merged output |
