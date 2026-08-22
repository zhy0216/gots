# FAQ

## What exactly is goTS?

A TypeScript-flavored language that transpiles to Go. It has its own type system (statically checked at compile time) and generates idiomatic Go, so programs become native binaries compiled by the Go toolchain.

## How is it different from TypeScript?

| | TypeScript | goTS |
|---|---|---|
| Runtime | JavaScript engine | Go binary |
| Types | Erased at runtime | Enforced at compile time, mapped to Go types |
| Numeric default | `number` (double) | `number` (float64); `int`/`float` are explicit Go interop types |
| Go interop | None (well-typed bindings) | Direct: `import { ... } from "go:strings"` |
| SQL | Via libraries | Built-in tagged-template SQL |
| async | Promise-based event loop | Promise-based event loop on the Go scheduler |

goTS is not a TypeScript compiler — it is a separate language with a compatible *flavor*.

## Why compile to Go instead of JavaScript?

Go binaries: fast startup, static typing, easy deployment, goroutines, and access to the whole Go ecosystem — while keeping TypeScript-style syntax, classes, closures, and modern JS ergonomics.

## What does `number` map to?

`float64`. Numeric literals (`42`, `3.14`) have type `number` by default. Annotate with `int` or `float` when you want Go's integer or float semantics:

```typescript
let x = 42          // number (float64)
let i: int = 42     // Go int
let f: float = 3.14 // Go float64
```

- `int` → `number`: allowed (widening)
- `number` → `int`: not implicit — use `toint(x)`
- Division `/` always returns `float`; modulo `%` requires `int` operands

## Does goTS support generics?

Yes — generic functions and classes with inference, explicit type arguments, and default type parameters. See [Language Spec §14](/language-spec#14-generics).

## Does goTS support async/await?

Yes. `async` functions return `Promise<T>`, `await` unwraps them, and an event loop (built on the Go scheduler) drives timers and microtasks. See [Language Spec §18](/language-spec#18-async-await-and-the-event-loop).

## How does SQL work?

`connect(path)` returns a database handle; SQL is written as tagged template literals and typed via generics:

```typescript
interface User { id: int, name: string }
const db = connect("./app.db")
const users = db.sql<User[]>`SELECT id, name FROM users`
```

Backed by SQLite (pure Go, no cgo). See [Built-in Reference: SQL](/built-in-reference#sql-database-connect).

## What Go packages can I import?

The Go stdlib packages declared under `pkg/declaration/stdlib/` — `strings`, `math`, `fmt`, `os`, `time`, `encoding/json`, `net/http`, `regexp`, `sort`, `strconv`, `bufio`, `bytes`, `io`, `path/filepath`. See [Language Spec §20](/language-spec#20-go-interop).

## What are the differences from JavaScript semantics?

See [Language Spec §26.1](/language-spec#26-1-differences-from-javascript). Highlights:

- `switch` has no fallthrough (Go semantics)
- `Math.round` rounds half away from zero (Go), not half up (JS)
- No hoisting; variables must be declared before use
- `typeof` returns `"number"` for `int` and `float`
- `getMonth()` is 0-11 (JS convention)

## What's not supported?

`do-while`, generators, symbols, Proxy/Reflect, WeakMap/WeakSet, BigInt, typed arrays, `eval`. See [Language Spec §26](/language-spec#26-unsupported-features).

## How does the compiler report errors?

In three stages with positions:

```text
Parser errors:
  line 3: unexpected token ...
Type errors:
  line 7: string not assignable to int
Codegen error: ...
```

## How do I run the docs site?

```bash
make docs-install  # npm install in docs/
make docs-serve    # dev server at http://localhost:5173
make docs-build    # static build → docs/.vitepress/dist
```

## Where are the benchmarks?

`benchmark/` — goTS-compiled binaries vs Bun, measured with hyperfine. The generated binary beat Bun by ~1.4x on the binary-tree benchmark.
