# Getting Started

goTS (GoTypeScript) is a TypeScript-like language that compiles to Go. Write `.gts` files, get native binaries.

## Prerequisites

- **Go 1.25+** — the generated code is compiled by the Go toolchain
- Nothing else! The compiler is a single Go binary.

## Build from Source

```bash
git clone https://github.com/zhy0216/quickts
cd quickts/gots

# Build the gots binary into the current directory
make build

# Or install to $GOPATH/bin
make install
```

## Your First Program

Create `hello.gts`:

```typescript
function greet(name: string): void {
    println("Hello, " + name + "!")
}

greet("goTS")

let numbers: int[] = [1, 2, 3, 4, 5]
let doubled: int[] = numbers.map((x: int): int => x * 2)
println(doubled)  // [2, 4, 6, 8, 10]
```

Run it:

```bash
gots run hello.gts
# Hello, goTS!
# [2, 4, 6, 8, 10]
```

## Compile to a Binary

```bash
# Build a native executable
gots build hello.gts -o hello

./hello
# Hello, goTS!
# [2, 4, 6, 8, 10]
```

## See the Generated Go

```bash
gots emit-go hello.gts
# writes hello.go — inspect the Go code goTS generated
```

## Try the REPL

```bash
gots repl
```

```text
goTS REPL v0.2.0 (Go transpiler)
Type 'exit' or press Ctrl+D to quit

>>> let x: int = 42
>>> println(x * 2)
84
>>> exit
```

## Next Steps

- **[Language Spec](/language-spec)** — the complete language reference
- **[CLI Reference](/guide/cli)** — every command and flag
- **[Built-in Reference](/built-in-reference)** — the standard library
- **Examples** — the `test/` directory contains runnable programs for every feature:
  - `test/generics.gts` — generics
  - `test/interfaces.gts` — structural typing
  - `test/eventloop.gts` — async/await, timers, microtasks
  - `test/sql_test.gts` — typed SQL
  - `test/go_imports.gts` — Go stdlib imports
  - `test/y_combinator.gts` — higher-order functions
  - `test/optional_chaining.gts`, `test/nullish_coalescing.gts` — modern syntax
