# GoTS Implementation Plan

A phased approach to building the GoTS compiler and VM in Go.

---

## Overview

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                        GoTS Pipeline                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   Source     ┌───────┐    ┌────────┐    ┌───────────────┐   │
│   Code   ──▶ │ Lexer │ ─▶ │ Parser │ ─▶ │ Type Checker  │   │
│   (.gts)     └───────┘    └────────┘    └───────────────┘   │
│                Tokens        AST          Typed AST          │
│                                               │              │
│                                               ▼              │
│   Result   ┌──────┐    ┌───────────┐    ┌──────────┐        │
│      ◀──── │  VM  │ ◀─ │ Bytecode  │ ◀─ │ Compiler │        │
│            └──────┘    │  (.gtsb)  │    └──────────┘        │
│                        └───────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
gots/
├── cmd/
│   └── gots/
│       └── main.go           # CLI entry point
├── pkg/
│   ├── token/
│   │   └── token.go          # Token types and definitions
│   ├── lexer/
│   │   └── lexer.go          # Lexical analysis
│   ├── ast/
│   │   └── ast.go            # AST node definitions
│   ├── parser/
│   │   └── parser.go         # Syntax analysis
│   ├── types/
│   │   ├── types.go          # Type representations
│   │   └── checker.go        # Type checking
│   ├── compiler/
│   │   ├── compiler.go       # Bytecode generation
│   │   └── opcode.go         # Opcode definitions
│   ├── vm/
│   │   ├── vm.go             # Virtual machine
│   │   ├── value.go          # Value representation
│   │   └── object.go         # Heap objects
│   └── bytecode/
│       ├── chunk.go          # Bytecode chunk
│       └── binary.go         # Binary format read/write
├── test/
│   └── *.gts                 # Test programs
├── docs/
│   ├── language-spec-v1.md
│   └── bytecode-spec-v1.md
└── go.mod
```

---

## Phase 1: Foundation

**Goal**: Lex and parse basic expressions, compile to bytecode, execute in VM.

### Milestone 1.1: Project Setup & Tokens

**Files**: `token/token.go`

```go
// Define all token types
type TokenType int

const (
    // Literals
    TOKEN_NUMBER
    TOKEN_STRING
    TOKEN_TRUE
    TOKEN_FALSE
    TOKEN_NULL

    // Operators
    TOKEN_PLUS
    TOKEN_MINUS
    // ... etc

    // Keywords
    TOKEN_LET
    TOKEN_CONST
    TOKEN_FUNCTION
    // ... etc
)

type Token struct {
    Type    TokenType
    Lexeme  string
    Line    int
    Column  int
}
```

**Deliverable**: Token type definitions

---

### Milestone 1.2: Lexer

**Files**: `lexer/lexer.go`

**Features**:
- Single-character tokens: `+ - * / % ( ) { } [ ] ; , . :`
- Multi-character tokens: `== != <= >= && || =>`
- Keywords: `let`, `const`, `function`, `if`, `else`, `while`, `for`, etc.
- Literals: numbers, strings (single and double quote), identifiers
- Comments: `//` and `/* */`
- Error reporting with line/column

**Test**: Lex sample programs, verify token stream

---

### Milestone 1.3: AST Definitions

**Files**: `ast/ast.go`

**Node Types**:
```go
// Expressions
type NumberLiteral struct { Value float64 }
type StringLiteral struct { Value string }
type BoolLiteral struct { Value bool }
type NullLiteral struct {}
type Identifier struct { Name string }
type BinaryExpr struct { Left, Right Expr; Op TokenType }
type UnaryExpr struct { Operand Expr; Op TokenType }
type CallExpr struct { Callee Expr; Args []Expr }
type IndexExpr struct { Object, Index Expr }
type PropertyExpr struct { Object Expr; Name string }
type ArrayLiteral struct { Elements []Expr }
type ObjectLiteral struct { Properties []Property }
type FunctionExpr struct { Params []Param; ReturnType Type; Body *Block }
type NewExpr struct { Class string; Args []Expr }
type ThisExpr struct {}
type AssignExpr struct { Target Expr; Value Expr }

// Statements
type ExprStmt struct { Expr Expr }
type VarDecl struct { Name string; Type Type; Init Expr; IsConst bool }
type Block struct { Stmts []Stmt }
type IfStmt struct { Cond Expr; Then *Block; Else Stmt }
type WhileStmt struct { Cond Expr; Body *Block }
type ForStmt struct { Init *VarDecl; Cond Expr; Update Expr; Body *Block }
type ReturnStmt struct { Value Expr }
type BreakStmt struct {}
type ContinueStmt struct {}

// Declarations
type FuncDecl struct { Name string; Params []Param; ReturnType Type; Body *Block }
type ClassDecl struct { Name string; Super string; Fields []Field; Constructor *Constructor; Methods []Method }
type TypeAlias struct { Name string; Type Type }

// Program
type Program struct { Decls []Decl }
```

**Deliverable**: Complete AST type definitions

---

### Milestone 1.4: Expression Parser

**Files**: `parser/parser.go`

**Features**:
- Pratt parser (precedence climbing) for expressions
- Operator precedence handling
- Grouping with parentheses
- Literals: numbers, strings, booleans, null
- Identifiers
- Binary and unary operators
- Error recovery and reporting

**Test**: Parse expressions, print AST

---

### Milestone 1.5: Value & Object System

**Files**: `vm/value.go`, `vm/object.go`

```go
// NaN-boxing or tagged union for values
type Value struct {
    Type ValueType
    data uint64
}

// Object types
type Object interface {
    Type() ObjectType
}

type ObjString struct { Value string; Hash uint32 }
type ObjArray struct { Elements []Value }
// ... etc
```

**Deliverable**: Value representation with constructors and accessors

---

### Milestone 1.6: Bytecode & Chunk

**Files**: `bytecode/chunk.go`, `compiler/opcode.go`

```go
type Chunk struct {
    Code      []byte
    Constants []Value
    Lines     []int
}

func (c *Chunk) Write(byte, line int)
func (c *Chunk) AddConstant(Value) int
```

**Deliverable**: Chunk with constant pool, opcode definitions

---

### Milestone 1.7: Expression Compiler

**Files**: `compiler/compiler.go`

**Features**:
- Compile literals to `OP_CONSTANT`, `OP_TRUE`, etc.
- Compile binary ops to arithmetic/comparison opcodes
- Compile unary ops
- Constant pool management

**Test**: Compile `1 + 2 * 3`, inspect bytecode

---

### Milestone 1.8: Basic VM

**Files**: `vm/vm.go`

**Features**:
- Stack operations: push, pop, peek
- Execute arithmetic opcodes
- Execute comparison opcodes
- Execute `OP_PRINT`/`OP_PRINTLN`

**Test**: Run `println(1 + 2 * 3)` end-to-end

---

## Phase 2: Variables & Control Flow

**Goal**: Support variables, scoping, and control flow statements.

### Milestone 2.1: Statement Parser

**Features**:
- `let` and `const` declarations
- Expression statements
- Block statements
- `if`/`else` statements
- `while` statements
- `for` statements
- `break`/`continue`

**Test**: Parse complete programs with control flow

---

### Milestone 2.2: Local Variables

**Compiler Features**:
- Track locals in current scope
- `OP_GET_LOCAL`, `OP_SET_LOCAL`
- Scope stack for block scoping
- `OP_POP`, `OP_POPN` for scope exit

**VM Features**:
- Stack slots for locals
- Local variable access

**Test**: Variable declaration, assignment, shadowing

---

### Milestone 2.3: Global Variables

**Compiler Features**:
- Distinguish globals from locals
- `OP_GET_GLOBAL`, `OP_SET_GLOBAL`
- Global name table

**VM Features**:
- Globals map
- Global variable access

**Test**: Global variables, access from functions

---

### Milestone 2.4: Control Flow Compilation

**Features**:
- `OP_JUMP`, `OP_JUMP_BACK`
- `OP_JUMP_IF_FALSE`, `OP_JUMP_IF_TRUE`
- Patch jump offsets after compiling body
- Loop break/continue with jump stack

**Test**: If statements, while loops, for loops, nested loops with break/continue

---

## Phase 3: Functions & Closures

**Goal**: Support function declarations, calls, and closures.

### Milestone 3.1: Function Parser

**Features**:
- Function declarations
- Function expressions
- Parameter lists with types
- Return type annotations
- `return` statements

**Test**: Parse function declarations and expressions

---

### Milestone 3.2: Function Compilation

**Compiler Features**:
- Compile function body to separate chunk
- Create `ObjFunction` objects
- `OP_CLOSURE` for function creation
- Track function nesting

**Test**: Compile simple functions

---

### Milestone 3.3: Function Calls

**Compiler Features**:
- Compile call expressions
- `OP_CALL` with argument count

**VM Features**:
- Call frames
- Parameter passing
- `OP_RETURN`
- Stack management

**Test**: Call functions, recursion (factorial)

---

### Milestone 3.4: Closures

**Compiler Features**:
- Resolve upvalues (captured variables)
- Emit upvalue descriptors with `OP_CLOSURE`
- `OP_GET_UPVALUE`, `OP_SET_UPVALUE`

**VM Features**:
- `ObjUpvalue` for captured variables
- Open upvalue list
- `OP_CLOSE_UPVALUE` when variables go out of scope

**Test**: Counter closure, nested closures

---

## Phase 4: Type System

**Goal**: Implement static type checking.

### Milestone 4.1: Type Representations

**Files**: `types/types.go`

```go
type Type interface { typeNode() }

type PrimitiveType struct { Kind PrimitiveKind } // number, string, boolean, void, null
type ArrayType struct { Element Type }
type ObjectType struct { Properties map[string]Type }
type FunctionType struct { Params []Type; Return Type }
type ClassType struct { Name string; Super *ClassType; Fields, Methods map[string]Type }
type NullableType struct { Inner Type }
type NamedType struct { Name string; Resolved Type }
```

**Deliverable**: Type AST and equality checking

---

### Milestone 4.2: Type Parser

**Features**:
- Parse type annotations in declarations
- Parse function types: `(a: number) => string`
- Parse array types: `number[]`, `string[][]`
- Parse object types: `{ x: number, y: number }`
- Parse nullable types: `string | null`
- Parse type aliases

**Test**: Parse complex type annotations

---

### Milestone 4.3: Type Checker - Expressions

**Files**: `types/checker.go`

**Features**:
- Infer/check types of literals
- Check binary operator types
- Check function call argument types
- Check property access types
- Check array indexing types
- Build symbol table

**Test**: Type check expressions, catch type errors

---

### Milestone 4.4: Type Checker - Statements

**Features**:
- Check variable declarations
- Check assignments
- Check return types
- Check if/while condition is boolean
- Scope tracking for variables

**Test**: Type check complete programs

---

### Milestone 4.5: Type Checker - Functions & Classes

**Features**:
- Check function signatures
- Check class field types
- Check method signatures
- Check inheritance compatibility
- Check `this` type in methods

**Test**: Type check classes, inheritance, method overriding

---

### Milestone 4.6: Null Safety

**Features**:
- Track nullable vs non-nullable types
- Narrow types after null checks
- Prevent operations on nullable without check

**Test**: Null safety errors, type narrowing in if blocks

---

## Phase 5: Classes & Objects

**Goal**: Support classes, objects, and arrays.

### Milestone 5.1: Object Literals

**Compiler Features**:
- `OP_OBJECT` to create object
- Property initialization

**VM Features**:
- `ObjObject` type
- Property access/mutation

**Test**: Create and use object literals

---

### Milestone 5.2: Array Literals

**Compiler Features**:
- `OP_ARRAY` to create array
- `OP_GET_INDEX`, `OP_SET_INDEX`

**VM Features**:
- `ObjArray` type
- Index bounds checking

**Test**: Create arrays, index access, mutation

---

### Milestone 5.3: Class Parser

**Features**:
- Parse class declarations
- Parse fields with types
- Parse constructor
- Parse methods
- Parse `extends`

**Test**: Parse class hierarchies

---

### Milestone 5.4: Class Compilation

**Compiler Features**:
- `OP_CLASS` to create class
- `OP_METHOD` to define methods
- `OP_INHERIT` for inheritance
- Constructor as special method

**VM Features**:
- `ObjClass` type
- Method table

**Test**: Define classes with methods

---

### Milestone 5.5: Instance Creation

**Compiler Features**:
- `new ClassName(args)` compilation
- Call constructor after allocation

**VM Features**:
- `ObjInstance` type
- Field storage
- `this` binding

**Test**: Create instances, access fields

---

### Milestone 5.6: Method Calls

**Compiler Features**:
- `OP_GET_PROPERTY` for method lookup
- `OP_INVOKE` for optimized method calls

**VM Features**:
- Method binding
- `this` in call frame

**Test**: Method calls, chained calls

---

### Milestone 5.7: Inheritance

**Compiler Features**:
- `OP_GET_SUPER` for super method lookup
- `OP_SUPER_INVOKE` for super calls
- `super(args)` in constructors

**VM Features**:
- Method resolution order
- Super method dispatch

**Test**: Inheritance, method override, super calls

---

## Phase 6: Runtime & Polish

**Goal**: Built-in functions, GC, error handling, CLI.

### Milestone 6.1: Built-in Functions

**Features**:
- `print`, `println`
- `len` (string and array)
- `toString`, `toNumber`
- `push`, `pop`
- `sqrt`, `floor`, `ceil`, `abs`

**Implementation**: `OP_BUILTIN` with function ID

**Test**: All built-in functions

---

### Milestone 6.2: Garbage Collection

**Features**:
- Track all allocated objects
- Mark roots (stack, globals, frames, upvalues)
- Trace references
- Sweep unmarked objects
- Trigger on allocation threshold

**Test**: Programs that create many objects, verify no leaks

---

### Milestone 6.3: Error Handling

**Features**:
- Runtime error with line info
- Stack traces
- Graceful error messages
- Compile errors with location

**Test**: Verify helpful error messages

---

### Milestone 6.4: Binary Format

**Files**: `bytecode/binary.go`

**Features**:
- Write compiled module to `.gtsb` file
- Read `.gtsb` file into module
- Magic number, version check

**Test**: Compile, save, load, run

---

### Milestone 6.5: CLI

**Files**: `cmd/gots/main.go`

**Commands**:
```
gots run <file.gts>      # Compile and run
gots compile <file.gts>  # Compile to .gtsb
gots exec <file.gtsb>    # Execute bytecode
gots repl                # Interactive mode
gots disasm <file.gtsb>  # Disassemble bytecode
```

**Test**: All CLI commands

---

### Milestone 6.6: REPL

**Features**:
- Read-eval-print loop
- Multi-line input
- History
- Error recovery

**Test**: Interactive usage

---

## Phase 7: Testing & Documentation

### Milestone 7.1: Test Suite

**Structure**:
```
test/
├── lexer/           # Lexer unit tests
├── parser/          # Parser unit tests
├── checker/         # Type checker unit tests
├── compiler/        # Compiler unit tests
├── vm/              # VM unit tests
├── integration/     # End-to-end tests
│   ├── basic/       # Expressions, variables
│   ├── control/     # Control flow
│   ├── functions/   # Functions, closures
│   ├── classes/     # Classes, inheritance
│   ├── types/       # Type checking
│   └── errors/      # Error handling
└── benchmark/       # Performance tests
```

**Test Runner**: Go's built-in testing + custom harness for `.gts` files

---

### Milestone 7.2: Documentation

- Update language spec with any changes
- Update bytecode spec with any changes
- Write user guide
- Document CLI usage
- Add code comments

---

## Implementation Order Summary

| Phase | Milestones | Key Deliverable |
|-------|------------|-----------------|
| 1 | 1.1 - 1.8 | Execute `println(1 + 2 * 3)` |
| 2 | 2.1 - 2.4 | Variables and control flow |
| 3 | 3.1 - 3.4 | Functions and closures |
| 4 | 4.1 - 4.6 | Static type checking |
| 5 | 5.1 - 5.7 | Classes and objects |
| 6 | 6.1 - 6.6 | Built-ins, GC, CLI |
| 7 | 7.1 - 7.2 | Tests and documentation |

---

## Dependencies Between Phases

```
Phase 1 (Foundation)
    │
    ▼
Phase 2 (Variables & Control)
    │
    ├─────────────────┐
    ▼                 ▼
Phase 3 (Functions)   Phase 4 (Types)
    │                 │
    └────────┬────────┘
             ▼
      Phase 5 (Classes)
             │
             ▼
      Phase 6 (Runtime)
             │
             ▼
      Phase 7 (Polish)
```

**Note**: Phase 4 (Types) can be developed in parallel with Phase 3, then integrated before Phase 5.

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Closure complexity | Follow Lua's upvalue design closely |
| Type checker complexity | Start simple, add features incrementally |
| GC bugs | Test with stress tests, use write barriers if needed |
| Performance | Profile after correctness, optimize hot paths |

---

## Success Criteria per Phase

| Phase | Criterion |
|-------|-----------|
| 1 | Can evaluate arithmetic expressions |
| 2 | Can run FizzBuzz |
| 3 | Can run recursive factorial |
| 4 | Catches type errors at compile time |
| 5 | Can run LinkedList example from spec |
| 6 | Full language working with GC |
| 7 | Comprehensive test coverage |

---

## Estimated Complexity

| Component | Lines of Code (est.) | Complexity |
|-----------|---------------------|------------|
| Lexer | 300-400 | Low |
| Parser | 800-1000 | Medium |
| AST | 300-400 | Low |
| Type System | 600-800 | Medium-High |
| Type Checker | 800-1000 | High |
| Compiler | 1000-1200 | Medium-High |
| VM | 800-1000 | Medium |
| Objects/GC | 400-500 | Medium |
| Built-ins | 200-300 | Low |
| CLI | 200-300 | Low |
| **Total** | **5500-7000** | - |

---

## Next Steps

1. Create Go module and directory structure
2. Implement Milestone 1.1 (Tokens)
3. Implement Milestone 1.2 (Lexer)
4. Proceed through Phase 1 milestones
5. Iterate with tests at each milestone
