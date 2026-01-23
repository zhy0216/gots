# GoTS (GoTypeScript)

A TypeScript-like language compiler and virtual machine written in Go.

## Project Overview

GoTS implements a complete compilation pipeline: lexer → parser → type checker → compiler → bytecode VM. It supports static typing, functions, closures, classes with inheritance, and garbage collection.

## Important Notes

**`../scaffold/` is for implementation reference only.** It contains:
- `quickjs/` - QuickJS JavaScript engine (C implementation reference)
- `typescript-go-main/` - TypeScript compiler in Go (reference for TypeScript semantics)

Do not modify files in scaffold/. Use them only as reference for understanding how language features should be implemented.

## Architecture

```
Source (.gts) → Lexer → Parser → Type Checker → Compiler → Bytecode (.gtsb) → VM
```

### Package Structure

| Package | Purpose |
|---------|---------|
| `pkg/token` | Token type definitions (85 token types) |
| `pkg/lexer` | Tokenization with line/column tracking |
| `pkg/ast` | AST node definitions (statements, expressions, types) |
| `pkg/parser` | Pratt parser with operator precedence |
| `pkg/types` | Type checker with scope and type narrowing |
| `pkg/bytecode` | Opcode definitions and binary serialization |
| `pkg/compiler` | AST to bytecode generation |
| `pkg/vm` | Stack-based bytecode interpreter with GC |
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
gots run program.gts           # Compile and execute
gots compile program.gts       # Compile to .gtsb
gots exec program.gtsb         # Execute bytecode
gots disasm program.gts        # Disassemble
gots repl                      # Interactive REPL
```

## Language Features

### Types
- Primitives: `number`, `string`, `boolean`, `void`, `null`
- Arrays: `number[]`, `string[]`
- Objects: `{x: number, y: string}`
- Functions: `(a: number) => string`
- Nullable: `string | null`
- Type aliases: `type Point = {x: number, y: number}`

### Syntax Examples

```typescript
// Variables
let x: number = 42;
const name: string = "GoTS";

// Functions
function factorial(n: number): number {
    if (n <= 1) { return 1; }
    return n * factorial(n - 1);
}

// Classes
class Animal {
    name: string;
    constructor(name: string) {
        this.name = name;
    }
    speak(): void {
        println(this.name);
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

// Arrays
let arr: number[] = [1, 2, 3];
push(arr, 4);
let last: number = pop(arr);
```

### Built-in Functions
`println`, `print`, `len`, `push`, `pop`, `typeof`, `tostring`, `tonumber`, `sqrt`, `floor`, `ceil`, `abs`

## Key Implementation Details

### Parser
- Pratt parsing with 10 precedence levels
- Two-token lookahead
- Error recovery (collects multiple errors)

### Compiler
- Max 256 local variables per scope
- Slot 0 reserved for function/this
- Upvalue chains for closure capture (Lua-style)

### VM
- Stack-based (256 elements max)
- 64 call frames max
- Mark-and-sweep GC (threshold starts at 1MB, grows 2x)

### Bytecode Format
- Magic: `GTSB` (0x47545342)
- Version: 1
- Supports gzip compression (.gtsb.gz)
- RLE-compressed line info

## Code Conventions

- `Obj*` prefix for heap-allocated objects (ObjString, ObjArray, etc.)
- `VAL_*` for value type constants
- `OP_*` for bytecode opcodes
- `TYPE_*` for compilation context types

## Testing

Each package has corresponding `*_test.go` files. Example programs in `test/` directory.

```bash
# Run specific package tests
go test -v ./pkg/lexer
go test -v ./pkg/parser
go test -v ./pkg/compiler
go test -v ./pkg/vm
```

## Development Workflow

1. Make changes to source
2. Run `go test ./...` to verify
3. Test with example: `go run ./cmd/gots run test/example.gts`
4. For bytecode changes, check disassembly: `go run ./cmd/gots disasm test/example.gts`
