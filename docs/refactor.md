# Refactoring Plan: VM to Go Transpiler (AOT)

This document outlines the plan to refactor GoTS from a bytecode VM interpreter to an ahead-of-time (AOT) compiler that transpiles to Go source code.

## Overview

### Current Architecture

```
Source (.gts) → Lexer → Parser → TypeChecker → Compiler → Bytecode (.gtsb) → VM
```

### New Architecture

```
Source (.gts) → Lexer → Parser → TypeChecker → TypedAST → GoCodeGen → Go Source → go build → Native Binary
```

### Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Backend | Go transpilation | Full Go runtime integration, no cgo overhead |
| GC | Go's GC (mark-sweep based) | Battle-tested, concurrent, no custom implementation needed |
| Closures | Display closures | Compile-time environment capture via Go closures |
| Go interop | Native | Direct access to all Go libraries |

### Benefits

- **Performance**: Go compiler generates optimized native code
- **Go ecosystem**: Full access to Go stdlib and third-party packages
- **Debugging**: Standard Go tooling (delve, pprof, etc.)
- **Deployment**: Single static binary
- **Maintenance**: No VM/bytecode infrastructure to maintain

---

## Phase 1: Preserve Types in AST

**Goal**: Modify the type checker to annotate the AST with resolved types instead of discarding them.

### 1.1 Create TypedAST Package

Create `pkg/typed/` to hold type-annotated AST nodes:

```
pkg/typed/
├── ast.go        # TypedNode interface and wrappers
├── expr.go       # Typed expression nodes
├── stmt.go       # Typed statement nodes
└── types.go      # Type representation (reuse from pkg/types)
```

Each typed node wraps the original AST node with its resolved type:

```go
// pkg/typed/ast.go
package typed

import (
    "github.com/zhy0216/quickts/gots/pkg/ast"
    "github.com/zhy0216/quickts/gots/pkg/types"
)

type Node interface {
    AST() ast.Node
    Type() types.Type
}

type Expr struct {
    Ast      ast.Expression
    ExprType types.Type
}

type BinaryExpr struct {
    Ast         *ast.BinaryExpr
    Left        Expr
    Right       Expr
    Op          string
    ResultType  types.Type
}

type CallExpr struct {
    Ast        *ast.CallExpression
    Callee     Expr
    Args       []Expr
    ReturnType types.Type
}

type FunctionDecl struct {
    Ast        *ast.FunctionDeclaration
    Name       string
    Params     []TypedParam
    Body       *BlockStmt
    ReturnType types.Type
    Captures   []Capture  // For closures: what variables are captured
}

type Capture struct {
    Name     string
    Type     types.Type
    Depth    int        // Lexical depth where variable is defined
    Index    int        // Slot index in that scope
}
```

### 1.2 Modify Type Checker

Update `pkg/types/checker.go` to return typed AST:

```go
// Current signature
func (c *Checker) Check(program *ast.Program) []error

// New signature
func (c *Checker) Check(program *ast.Program) (*typed.Program, []error)
```

Changes needed in `pkg/types/checker.go`:
- Add methods to construct typed nodes during checking
- Track capture sets for closures during scope analysis
- Return `*typed.Program` containing all type-annotated declarations

### 1.3 Files to Modify

| File | Changes |
|------|---------|
| `pkg/types/checker.go` | Return TypedAST, track captures |
| `pkg/types/scope.go` | Add capture tracking for closures |
| `pkg/ast/ast.go` | No changes (keep original AST) |

---

## Phase 2: Go Code Generator

**Goal**: Create a new package that generates Go source code from TypedAST.

### 2.1 Create CodeGen Package

```
pkg/codegen/
├── codegen.go     # Main generator, program structure
├── expr.go        # Expression generation
├── stmt.go        # Statement generation
├── types.go       # GTS types → Go types mapping
├── runtime.go     # Runtime support code generation
├── names.go       # Name mangling and collision avoidance
└── builtins.go    # Built-in function mapping
```

### 2.2 Type Mapping

| GTS Type | Go Type |
|----------|---------|
| `number` | `float64` |
| `string` | `string` |
| `boolean` | `bool` |
| `null` | `*struct{}` (nil) or use option type |
| `T[]` | `[]T` |
| `{a: T, b: U}` | `struct{ A T; B U }` or `map[string]any` |
| `T \| null` | `*T` (pointer, nil = null) |
| `(a: T) => U` | `func(T) U` |
| `class C` | `type C struct { ... }` |

### 2.3 Core Generator Structure

```go
// pkg/codegen/codegen.go
package codegen

import (
    "bytes"
    "go/format"
    "github.com/zhy0216/quickts/gots/pkg/typed"
)

type Generator struct {
    buf        *bytes.Buffer
    indent     int
    imports    map[string]bool
    typeDecls  []string        // Forward declarations
    funcDecls  []string        // Generated functions
    mainBody   []string        // Main function body
}

func Generate(prog *typed.Program) ([]byte, error) {
    g := &Generator{
        buf:     new(bytes.Buffer),
        imports: make(map[string]bool),
    }

    g.genProgram(prog)

    return format.Source(g.buf.Bytes())
}

func (g *Generator) genProgram(prog *typed.Program) {
    g.writeln("package main")
    g.writeln("")

    // Imports
    g.genImports()

    // Runtime support
    g.genRuntime()

    // Type declarations (classes, type aliases)
    for _, decl := range prog.TypeDecls {
        g.genTypeDecl(decl)
    }

    // Function declarations
    for _, fn := range prog.Functions {
        g.genFunction(fn)
    }

    // Main function (top-level statements)
    g.writeln("func main() {")
    g.indent++
    for _, stmt := range prog.TopLevel {
        g.genStmt(stmt)
    }
    g.indent--
    g.writeln("}")
}
```

### 2.4 Expression Generation

```go
// pkg/codegen/expr.go

func (g *Generator) genExpr(e typed.Expr) string {
    switch expr := e.(type) {
    case *typed.NumberLit:
        return fmt.Sprintf("%v", expr.Value)

    case *typed.StringLit:
        return fmt.Sprintf("%q", expr.Value)

    case *typed.BinaryExpr:
        left := g.genExpr(expr.Left)
        right := g.genExpr(expr.Right)

        // Handle string concatenation
        if expr.Op == "+" && expr.Left.Type().IsString() {
            return fmt.Sprintf("(%s + %s)", left, right)
        }

        return fmt.Sprintf("(%s %s %s)", left, g.mapOp(expr.Op), right)

    case *typed.CallExpr:
        callee := g.genExpr(expr.Callee)
        args := g.genArgs(expr.Args)
        return fmt.Sprintf("%s(%s)", callee, args)

    case *typed.MemberExpr:
        obj := g.genExpr(expr.Object)
        return fmt.Sprintf("%s.%s", obj, g.exportName(expr.Property))

    case *typed.IndexExpr:
        obj := g.genExpr(expr.Object)
        index := g.genExpr(expr.Index)
        return fmt.Sprintf("%s[int(%s)]", obj, index)

    case *typed.FunctionExpr:
        return g.genClosure(expr)

    // ... more cases
    }
}
```

### 2.5 Statement Generation

```go
// pkg/codegen/stmt.go

func (g *Generator) genStmt(s typed.Stmt) {
    switch stmt := s.(type) {
    case *typed.VarDecl:
        if stmt.Init != nil {
            init := g.genExpr(stmt.Init)
            g.writeln("%s := %s", stmt.Name, init)
        } else {
            g.writeln("var %s %s", stmt.Name, g.goType(stmt.Type))
        }

    case *typed.IfStmt:
        cond := g.genExpr(stmt.Condition)
        g.writeln("if %s {", cond)
        g.indent++
        g.genStmt(stmt.Then)
        g.indent--
        if stmt.Else != nil {
            g.writeln("} else {")
            g.indent++
            g.genStmt(stmt.Else)
            g.indent--
        }
        g.writeln("}")

    case *typed.WhileStmt:
        cond := g.genExpr(stmt.Condition)
        g.writeln("for %s {", cond)
        g.indent++
        g.genStmt(stmt.Body)
        g.indent--
        g.writeln("}")

    case *typed.ReturnStmt:
        if stmt.Value != nil {
            val := g.genExpr(stmt.Value)
            g.writeln("return %s", val)
        } else {
            g.writeln("return")
        }

    // ... more cases
    }
}
```

### 2.6 Closure Generation (Display Closures)

```go
// pkg/codegen/codegen.go

func (g *Generator) genClosure(fn *typed.FunctionExpr) string {
    if len(fn.Captures) == 0 {
        // No captures, simple function
        return g.genSimpleFunc(fn)
    }

    // Generate closure that captures variables
    // Go handles this natively - captured vars become part of closure
    params := g.genParams(fn.Params)
    retType := g.goType(fn.ReturnType)

    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf("func(%s) %s {\n", params, retType))

    // Body - captured variables are accessed directly
    // Go's closure semantics match what we need
    for _, stmt := range fn.Body.Statements {
        buf.WriteString(g.genStmtToString(stmt))
    }

    buf.WriteString("}")
    return buf.String()
}
```

### 2.7 Built-in Function Mapping

```go
// pkg/codegen/builtins.go

var builtinMap = map[string]string{
    "println": "fmt.Println",
    "print":   "fmt.Print",
    "len":     "gts_len",      // Custom wrapper for unified len
    "push":    "gts_push",     // Append wrapper
    "pop":     "gts_pop",      // Pop wrapper
    "typeof":  "gts_typeof",   // Runtime type check
    "toString": "gts_toString",
    "toNumber": "gts_toNumber",
    "sqrt":    "math.Sqrt",
    "floor":   "math.Floor",
    "ceil":    "math.Ceil",
    "abs":     "math.Abs",
}

// Runtime helpers generated in output
const runtimeHelpers = `
func gts_len(v any) float64 {
    switch x := v.(type) {
    case string:
        return float64(len(x))
    case []any:
        return float64(len(x))
    default:
        panic("len: invalid type")
    }
}

func gts_push[T any](arr *[]T, val T) {
    *arr = append(*arr, val)
}

func gts_pop[T any](arr *[]T) T {
    n := len(*arr)
    val := (*arr)[n-1]
    *arr = (*arr)[:n-1]
    return val
}

func gts_typeof(v any) string {
    if v == nil {
        return "null"
    }
    switch v.(type) {
    case float64:
        return "number"
    case string:
        return "string"
    case bool:
        return "boolean"
    default:
        return "object"
    }
}
`
```

---

## Phase 3: Class and Object Support

**Goal**: Generate Go structs and methods for GTS classes.

### 3.1 Class to Struct Mapping

```typescript
// GTS input
class Point {
    x: number
    y: number

    constructor(x: number, y: number) {
        this.x = x
        this.y = y
    }

    distance(): number {
        return sqrt(this.x * this.x + this.y * this.y)
    }
}
```

```go
// Go output
type Point struct {
    X float64
    Y float64
}

func NewPoint(x float64, y float64) *Point {
    p := &Point{}
    p.X = x
    p.Y = y
    return p
}

func (this *Point) Distance() float64 {
    return math.Sqrt(this.X*this.X + this.Y*this.Y)
}
```

### 3.2 Inheritance

```typescript
// GTS input
class ColorPoint extends Point {
    color: string

    constructor(x: number, y: number, color: string) {
        super(x, y)
        this.color = color
    }
}
```

```go
// Go output - embedding for inheritance
type ColorPoint struct {
    Point  // Embedded
    Color string
}

func NewColorPoint(x float64, y float64, color string) *ColorPoint {
    cp := &ColorPoint{}
    cp.Point = *NewPoint(x, y)  // super() call
    cp.Color = color
    return cp
}
```

---

## Phase 4: Nullable Types and Union Types

**Goal**: Handle nullable types using Go pointers and interfaces.

### 4.1 Nullable Primitives

```typescript
// GTS
let x: number | null = null
x = 42
```

```go
// Go - use pointer for nullable
var x *float64 = nil
tmp := float64(42)
x = &tmp
```

### 4.2 Optional Chaining

```typescript
// GTS
let len = arr?.length
```

```go
// Go
var len float64
if arr != nil {
    len = float64(len(*arr))
}
```

### 4.3 Nullish Coalescing

```typescript
// GTS
let val = x ?? defaultValue
```

```go
// Go
var val float64
if x != nil {
    val = *x
} else {
    val = defaultValue
}
```

---

## Phase 5: Remove VM and Bytecode

**Goal**: Delete all VM-related code after the new codegen is working.

### 5.1 Packages to Remove

| Package | Reason |
|---------|--------|
| `pkg/vm/` | No longer interpreting bytecode |
| `pkg/bytecode/` | No longer generating bytecode |
| `pkg/compiler/` | Replaced by `pkg/codegen/` |

### 5.2 Files to Remove

```
pkg/vm/
├── vm.go           # DELETE
├── value.go        # DELETE
├── object.go       # DELETE
├── gc.go           # DELETE (using Go's GC)
├── builtins.go     # DELETE (reimplemented in codegen)
└── debug.go        # DELETE

pkg/bytecode/
├── chunk.go        # DELETE
├── opcode.go       # DELETE
├── binary.go       # DELETE
└── disasm.go       # DELETE

pkg/compiler/
├── compiler.go     # DELETE (replaced by codegen)
└── locals.go       # DELETE
```

### 5.3 Commands to Update

Update `cmd/gots/main.go`:

```go
// Old flow
source := readFile(path)
tokens := lexer.Lex(source)
ast := parser.Parse(tokens)
errors := checker.Check(ast)
chunk := compiler.Compile(ast)
vm.Run(chunk)

// New flow
source := readFile(path)
tokens := lexer.Lex(source)
ast := parser.Parse(tokens)
typedAST, errors := checker.Check(ast)
goCode := codegen.Generate(typedAST)
writeFile(outPath, goCode)
exec.Command("go", "build", outPath).Run()
```

---

## Phase 6: CLI and Tooling

**Goal**: Update the CLI to support the new compilation model.

### 6.1 New Commands

```bash
# Compile to Go source (for inspection)
gots build --emit-go program.gts -o program.go

# Compile to binary
gots build program.gts -o program

# Run directly (compile + execute)
gots run program.gts

# Format GTS source
gots fmt program.gts
```

### 6.2 Output Modes

| Flag | Output |
|------|--------|
| `--emit-go` | Generate Go source file |
| `--emit-ast` | Dump typed AST (for debugging) |
| (default) | Compile to native binary |

### 6.3 Build Process

```go
func buildBinary(source, output string) error {
    // 1. Generate Go code
    goCode, err := compile(source)
    if err != nil {
        return err
    }

    // 2. Write to temp file
    tmpDir, _ := os.MkdirTemp("", "gots-build-*")
    defer os.RemoveAll(tmpDir)

    goFile := filepath.Join(tmpDir, "main.go")
    os.WriteFile(goFile, goCode, 0644)

    // 3. Initialize go module
    exec.Command("go", "mod", "init", "gts_program").Run()

    // 4. Build
    cmd := exec.Command("go", "build", "-o", output, goFile)
    return cmd.Run()
}
```

---

## Phase 7: Go Library Integration

**Goal**: Allow importing and using Go packages from GTS.

### 7.1 Import Syntax

```typescript
// GTS
import { Println, Sprintf } from "fmt"
import { ReadFile } from "os"
import { Now } from "time"

let content = ReadFile("data.txt")
Println(Sprintf("Read at %v", Now()))
```

### 7.2 Generated Go

```go
import (
    "fmt"
    "os"
    "time"
)

func main() {
    content, _ := os.ReadFile("data.txt")
    fmt.Println(fmt.Sprintf("Read at %v", time.Now()))
}
```

### 7.3 Type Declarations for Go Packages

Create `.d.gts` declaration files for Go packages:

```typescript
// go/fmt.d.gts
declare module "fmt" {
    function Println(...args: any[]): void
    function Printf(format: string, ...args: any[]): void
    function Sprintf(format: string, ...args: any[]): string
}

// go/os.d.gts
declare module "os" {
    function ReadFile(name: string): [string, Error | null]
    function WriteFile(name: string, data: string): Error | null
}
```

---

## Implementation Order

### Milestone 1: Basic Transpilation (Week 1-2)
- [ ] Create `pkg/typed/` with typed AST nodes
- [ ] Modify type checker to return typed AST
- [ ] Create basic `pkg/codegen/` generator
- [ ] Support: variables, functions, if/while, basic expressions
- [ ] Generate and run simple programs

### Milestone 2: Complete Language Support (Week 3-4)
- [ ] Classes and inheritance
- [ ] Closures with capture
- [ ] Arrays and objects
- [ ] All operators and built-ins
- [ ] Break/continue

### Milestone 3: Nullable and Unions (Week 5)
- [ ] Nullable types via pointers
- [ ] Optional chaining
- [ ] Nullish coalescing
- [ ] Type narrowing in conditionals

### Milestone 4: Remove VM (Week 6)
- [ ] Delete `pkg/vm/`, `pkg/bytecode/`, `pkg/compiler/`
- [ ] Update all tests to use new codegen
- [ ] Update CLI commands
- [ ] Update documentation

### Milestone 5: Go Integration (Week 7-8)
- [ ] Go package imports
- [ ] Declaration files for common Go packages
- [ ] Error handling patterns
- [ ] Goroutine support (future)

---

## Testing Strategy

### Unit Tests
- Each codegen component tested in isolation
- Compare generated Go code against expected output

### Integration Tests
- Existing `.gts` test files should produce same output
- Run generated Go code and compare results to expected

### Regression Tests
```bash
# For each test case:
gots run test.gts > actual.txt
diff expected.txt actual.txt
```

---

## File Structure After Refactor

```
pkg/
├── lexer/           # KEEP - unchanged
├── parser/          # KEEP - unchanged
├── ast/             # KEEP - unchanged
├── types/           # MODIFY - return typed AST
│   ├── checker.go
│   ├── scope.go
│   └── types.go
├── typed/           # NEW - typed AST
│   ├── ast.go
│   ├── expr.go
│   ├── stmt.go
│   └── program.go
├── codegen/         # NEW - Go code generator
│   ├── codegen.go
│   ├── expr.go
│   ├── stmt.go
│   ├── types.go
│   ├── runtime.go
│   ├── names.go
│   └── builtins.go
├── vm/              # DELETE
├── bytecode/        # DELETE
└── compiler/        # DELETE

cmd/
└── gots/
    └── main.go      # MODIFY - new build pipeline
```

---

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Go codegen complexity | Start with subset, expand incrementally |
| Type mapping edge cases | Document limitations, fail clearly |
| Performance regression | Benchmark against VM, optimize hotspots |
| Breaking changes | Keep VM in parallel until codegen stable |

---

## Open Questions

1. **Error handling**: Should we use Go's `error` return pattern or panic/recover?
2. **Generics**: How to handle GTS generic-like patterns when targeting Go?
3. **Concurrency**: Future support for goroutines/channels?
4. **Modules**: Multi-file GTS programs and import resolution?

---

## References

- Go compiler internals: https://go.dev/src/cmd/compile/
- Display closures: https://en.wikipedia.org/wiki/Call_stack#Structure
- Go code generation patterns: `go/ast`, `go/format`, `go/printer`
