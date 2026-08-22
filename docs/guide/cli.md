# CLI Reference

The `gots` binary exposes a small set of commands. Running `gots` without arguments (or `gots help`) prints usage.

## Commands

### run

Compile and run a `.gts` source file.

```bash
gots run <file.gts>
```

Compiles to Go in a temporary directory and executes it with `go run`. Errors are reported as **parser errors**, **type errors**, or **codegen errors** with line/column information.

```bash
gots run test/example.gts
```

::: tip
If the program imports Go packages that need module support (e.g. SQL via `connect()`), `gots` initializes a temporary Go module and runs `go mod tidy` automatically.
:::

### build

Compile to a native binary (or Go source with `--emit-go`).

```bash
gots build <file.gts> [-o <output>] [--emit-go]
```

| Flag | Description |
|------|-------------|
| `-o <output>` | Output file path. Defaults to the input name without the `.gts` extension |
| `--emit-go` | Write Go source instead of compiling a binary |

```bash
gots build app.gts                 # builds ./app
gots build app.gts -o bin/myapp    # builds ./bin/myapp
gots build app.gts --emit-go       # writes app.go
```

### emit-go

Generate Go source code without compiling.

```bash
gots emit-go <file.gts> [output.go]
```

```bash
gots emit-go app.gts              # writes app.go
gots emit-go app.gts gen/app.go   # writes gen/app.go
```

### repl

Start an interactive session.

```bash
gots repl
```

The REPL is stateful: `let`, `const`, `function`, and `class` declarations persist across lines. Multi-line input (blocks) is supported.

### version

```bash
gots version
# gots version 0.2.0 (Go transpiler)
```

### Shorthand

Passing a path directly runs it:

```bash
gots program.gts    # equivalent to: gots run program.gts
```

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Parse/type/codegen error, build failure, or usage error |
| `N` | The compiled program exited with status `N` (propagated by `gots run`) |

## Compilation Pipeline

Every command funnels through the same pipeline:

```
.gts source → Lexer → Parser → TypedAST Builder (type checking) → Go Code Generator → go toolchain
```

Compile errors are reported in three stages:

```text
Parser errors:   syntax problems (collected, not fail-fast)
Type errors:     type mismatches from the TypedAST builder
Codegen error:   Go generation failures
```

## Makefile Equivalents

The project Makefile wraps these commands:

```bash
make run FILE=test/example.gts    # gots run
make emit-go FILE=test/example.gts # gots emit-go
make build                          # go build ./cmd/gots
make install                        # go install ./cmd/gots
make repl                           # go run ./cmd/gots repl
```
