# GoTS Implementation Plan

This document outlines the implementation plan for features not yet implemented in GoTS.

---

## Priority Overview

| Priority | Feature | Complexity | Status |
|----------|---------|------------|--------|
| 🔴 High | **Type inference for variables** | Medium | TODO |
| 🔴 High | Arrow functions | Low | TODO |
| 🔴 High | Compound assignment (`+=`, `-=`, etc.) | Low | TODO |
| 🔴 High | Array methods (map, filter, reduce, etc.) | Medium | TODO |
| 🔴 High | String methods | Medium | TODO |
| 🟡 Medium | Increment/decrement (`++`, `--`) | Low | TODO |
| 🟡 Medium | Optional chaining (`?.`) | Medium | TODO |
| 🟡 Medium | Nullish coalescing (`??`) | Low | TODO |
| 🟡 Medium | Switch statement | Medium | TODO |
| 🟡 Medium | For-of loop | Low | TODO |
| 🟢 Low | Try/catch/finally | High | TODO |
| 🟢 Low | Template literals | Medium | TODO |
| 🟢 Low | Static class members | Medium | TODO |
| 🟢 Low | Interfaces | Medium | TODO |
| ⚪ Very Low | Generics | Very High | TODO |

---

## High Priority Features

### 1. Type Inference for Variables

**Design Principle:** Infer types for local variables while **requiring explicit signatures for functions**. This follows the Rust/Go philosophy - explicit at boundaries, inferred internally.

> **Note:** Type inference is **compile-time only**. It happens in the type checker phase. The compiler and VM are not affected - they don't know or care about types. The bytecode is identical whether types are explicit or inferred.

**Syntax:**
```typescript
// Variable inference - type annotation optional
let x = 10;                      // inferred as number
let name = "hello";              // inferred as string
let flag = true;                 // inferred as boolean
let arr = [1, 2, 3];             // inferred as number[]
let obj = { x: 1, y: 2 };        // inferred as { x: number, y: number }
let maybeNull = null;            // inferred as null (needs explicit type for nullable)

// Explicit type still allowed
let y: number = 20;

// const inference
const PI = 3.14159;              // inferred as number

// Functions MUST have explicit signatures (no inference)
function add(a: number, b: number): number {  // ✅ Required
    return a + b;
}

function bad(a, b) {             // ❌ Error: parameter types required
    return a + b;
}

// Arrow functions also require types
let double = (x: number): number => x * 2;    // ✅ Required
```

**What Gets Inferred:**
- `let` and `const` variable initializers
- Array literal element types
- Object literal property types
- Return type of expressions (for type checking, not declarations)

**What Requires Explicit Types:**
- Function parameters (always)
- Function return types (always)
- Class fields (always)
- Method parameters and return types (always)
- Variables without initializers: `let x: number;`
- Nullable variables: `let x: string | null = null;`

**Implementation Steps:**

1. **AST** (`pkg/ast/ast.go`)
   - Make `VarType` field nullable in `VarDecl`:
   ```go
   type VarDecl struct {
       Token   token.Token
       Name    string
       VarType Type       // nil when type should be inferred
       Value   Expression
       IsConst bool
   }
   ```

2. **Parser** (`pkg/parser/parser.go`)
   - Modify `parseVarDecl()` to make type annotation optional:
   ```go
   // Current: let x: number = 10;
   // New:     let x = 10;  OR  let x: number = 10;

   func (p *Parser) parseVarDecl() *ast.VarDecl {
       // ... parse name ...

       var varType ast.Type
       if p.peekTokenIs(token.COLON) {
           p.nextToken() // consume ':'
           varType = p.parseType()
       }
       // varType is nil if no annotation

       // ... parse = value ...
   }
   ```

3. **Type Checker** (`pkg/types/checker.go`)
   - Add `inferType()` function for expressions
   - Modify `checkVarDecl()` to infer when type is nil:
   ```go
   func (c *Checker) checkVarDecl(decl *ast.VarDecl) {
       if decl.Value == nil && decl.VarType == nil {
           c.error(decl.Token.Line, decl.Token.Column,
               "variable declaration requires type annotation or initializer")
           return
       }

       var declaredType Type
       if decl.VarType != nil {
           // Explicit type annotation
           declaredType = c.resolveType(decl.VarType)
           if decl.Value != nil {
               initType := c.checkExpr(decl.Value)
               if !IsAssignableTo(initType, declaredType) {
                   c.error(...)
               }
           }
       } else {
           // Infer type from initializer
           declaredType = c.inferType(decl.Value)
       }

       c.scope.Define(decl.Name, declaredType)
   }

   func (c *Checker) inferType(expr ast.Expression) Type {
       switch e := expr.(type) {
       case *ast.NumberLiteral:
           return NumberType
       case *ast.StringLiteral:
           return StringType
       case *ast.BoolLiteral:
           return BooleanType
       case *ast.NullLiteral:
           return NullType  // Note: bare null has limited usefulness
       case *ast.ArrayLiteral:
           return c.inferArrayType(e)
       case *ast.ObjectLiteral:
           return c.inferObjectType(e)
       case *ast.Identifier:
           if typ, found := c.scope.Lookup(e.Name); found {
               return typ
           }
           return AnyType
       case *ast.CallExpr:
           return c.checkCallExpr(e)  // Use return type
       case *ast.BinaryExpr:
           return c.checkBinaryExpr(e)
       // ... other cases ...
       default:
           return c.checkExpr(expr)
       }
   }
   ```

4. **Error Messages**
   - "variable declaration requires type annotation or initializer"
   - "cannot infer type of null literal, use explicit type annotation"
   - "function parameter 'x' requires type annotation"
   - "function return type required"

**Edge Cases to Handle:**

```typescript
// 1. Null literal alone - require explicit type
let x = null;                    // Error or infer as null?
let x: string | null = null;     // ✅ Explicit nullable

// 2. Empty array - cannot infer element type
let arr = [];                    // Error: cannot infer element type
let arr: number[] = [];          // ✅ Explicit type

// 3. Mixed array
let mixed = [1, "two"];          // Error or infer union?
                                 // Recommend: Error, require explicit

// 4. Reassignment must match inferred type
let x = 10;
x = "hello";                     // Error: cannot assign string to number

// 5. Function expressions still need types
let fn = function(x: number): number { return x; };  // ✅
let fn = function(x) { return x; };                  // ❌ Error
```

**Files to modify:**
- `pkg/ast/ast.go` - Make VarType nullable
- `pkg/parser/parser.go` - Optional type annotation
- `pkg/types/checker.go` - Add inference logic

> **Note:** Compiler (`pkg/compiler/compiler.go`) and VM are **not affected**. They already work without type information.

**Estimated effort:** 4-6 hours

---

### 2. Arrow Functions

**Syntax:**
```typescript
let add = (a: number, b: number): number => a + b;
let double = (x: number): number => { return x * 2; };
```

**Implementation Steps:**

1. **Parser** (`pkg/parser/parser.go`)
   - Detect `(params) =>` pattern in expression parsing
   - Support both expression body `=> expr` and block body `=> { stmts }`
   - Reuse `FunctionExpr` AST node or create `ArrowFunctionExpr`

2. **Type Checker** (`pkg/types/checker.go`)
   - Same handling as function expressions

3. **Compiler** (`pkg/compiler/compiler.go`)
   - Same as function expression compilation

**Files to modify:**
- `pkg/parser/parser.go`

**Estimated effort:** 2-3 hours

---

### 3. Compound Assignment Operators

**Syntax:**
```typescript
x += 5;    // x = x + 5
x -= 3;    // x = x - 3
x *= 2;    // x = x * 2
x /= 4;    // x = x / 4
x %= 3;    // x = x % 3
```

**Implementation Steps:**

1. **Lexer** (`pkg/lexer/lexer.go`)
   - Add recognition for `+=`, `-=`, `*=`, `/=`, `%=`

2. **Tokens** (`pkg/token/token.go`)
   - Add tokens: `PLUS_ASSIGN`, `MINUS_ASSIGN`, `STAR_ASSIGN`, `SLASH_ASSIGN`, `PERCENT_ASSIGN`

3. **Parser** (`pkg/parser/parser.go`)
   - Desugar `x += y` to `x = x + y` during parsing
   - Handle compound assignment on properties and indices

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/parser/parser.go`

**Estimated effort:** 2-3 hours

---

### 4. Array Methods

**Syntax:**
```typescript
arr.map((x) => x * 2);
arr.filter((x) => x > 0);
arr.reduce((acc, x) => acc + x, 0);
arr.forEach((x) => println(x));
arr.find((x) => x > 5);
arr.includes(value);
arr.indexOf(value);
arr.slice(start, end);
arr.concat(other);
arr.join(separator);
arr.reverse();
arr.sort();
```

**Implementation Steps:**

1. **VM** (`pkg/vm/vm.go`)
   - Add method dispatch for array objects
   - Implement each method as native function

2. **Type Checker** (`pkg/types/checker.go`)
   - Add array method types to type system
   - Handle generic callback types

3. **Bytecode** (`pkg/bytecode/opcode.go`)
   - Consider adding `OP_ARRAY_METHOD` or use existing `OP_INVOKE`

**Files to modify:**
- `pkg/vm/vm.go`
- `pkg/vm/value.go`
- `pkg/types/checker.go`

**Estimated effort:** 6-8 hours

---

### 5. String Methods

**Syntax:**
```typescript
str.toUpperCase();
str.toLowerCase();
str.split(delimiter);
str.trim();
str.trimStart();
str.trimEnd();
str.substring(start, end);
str.indexOf(substr);
str.lastIndexOf(substr);
str.replace(old, new);
str.startsWith(prefix);
str.endsWith(suffix);
str.includes(substr);
str.repeat(count);
str.charAt(index);
str.charCodeAt(index);
```

**Implementation Steps:**

1. **VM** (`pkg/vm/vm.go`)
   - Add method dispatch for string objects
   - Implement each method using Go's strings package

2. **Type Checker** (`pkg/types/checker.go`)
   - Add string method types

**Files to modify:**
- `pkg/vm/vm.go`
- `pkg/vm/value.go`
- `pkg/types/checker.go`

**Estimated effort:** 4-6 hours

---

## Medium Priority Features

### 6. Increment/Decrement Operators

**Syntax:**
```typescript
i++;    // post-increment
++i;    // pre-increment
i--;    // post-decrement
--i;    // pre-decrement
```

**Implementation Steps:**

1. **Tokens** (`pkg/token/token.go`)
   - Add `PLUS_PLUS`, `MINUS_MINUS`

2. **Lexer** (`pkg/lexer/lexer.go`)
   - Recognize `++` and `--`

3. **AST** (`pkg/ast/ast.go`)
   - Add `UpdateExpr` node with `Prefix` boolean and `Op` field

4. **Parser** (`pkg/parser/parser.go`)
   - Handle prefix in unary expression parsing
   - Handle postfix after identifiers/properties

5. **Type Checker** (`pkg/types/checker.go`)
   - Operand must be number
   - Target must be assignable (identifier, property, index)

6. **Compiler** (`pkg/compiler/compiler.go`)
   - Pre: load, increment, store, (value on stack)
   - Post: load, dup, increment, store, pop (old value on stack)

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/ast/ast.go`
- `pkg/parser/parser.go`
- `pkg/types/checker.go`
- `pkg/compiler/compiler.go`

**Estimated effort:** 3-4 hours

---

### 7. Optional Chaining

**Syntax:**
```typescript
let name = user?.profile?.name;
let result = obj?.method?.(arg);
arr?.[index];
```

**Implementation Steps:**

1. **Tokens** (`pkg/token/token.go`)
   - Add `QUESTION_DOT` for `?.`

2. **Lexer** (`pkg/lexer/lexer.go`)
   - Recognize `?.`

3. **AST** (`pkg/ast/ast.go`)
   - Add `Optional` boolean flag to `PropertyExpr`, `IndexExpr`, `CallExpr`

4. **Parser** (`pkg/parser/parser.go`)
   - Handle `?.` in property/index/call parsing

5. **Compiler** (`pkg/compiler/compiler.go`)
   - Generate: check null → jump to end (push null) or continue chain

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/ast/ast.go`
- `pkg/parser/parser.go`
- `pkg/types/checker.go`
- `pkg/compiler/compiler.go`

**Estimated effort:** 4-5 hours

---

### 8. Nullish Coalescing

**Syntax:**
```typescript
let value = maybeNull ?? defaultValue;
```

**Implementation Steps:**

1. **Tokens** (`pkg/token/token.go`)
   - Add `NULLISH_COALESCE` for `??`

2. **Lexer** (`pkg/lexer/lexer.go`)
   - Recognize `??`

3. **Parser** (`pkg/parser/parser.go`)
   - Add `??` as binary operator with appropriate precedence (lower than `||`)

4. **Compiler** (`pkg/compiler/compiler.go`)
   - Generate: eval left → dup → check null → jump if not null → pop → eval right

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/parser/parser.go`
- `pkg/compiler/compiler.go`

**Estimated effort:** 2-3 hours

---

### 9. Switch Statement

**Syntax:**
```typescript
switch (value) {
    case 1:
        println("one");
        break;
    case 2:
    case 3:
        println("two or three");
        break;
    default:
        println("other");
}
```

**Implementation Steps:**

1. **Tokens** (`pkg/token/token.go`)
   - Add `SWITCH`, `CASE`, `DEFAULT`

2. **Lexer** (`pkg/lexer/lexer.go`)
   - Add keywords to keyword map

3. **AST** (`pkg/ast/ast.go`)
   ```go
   type SwitchStmt struct {
       Token       token.Token
       Discriminant Expression
       Cases       []*CaseClause
   }

   type CaseClause struct {
       Token      token.Token
       Test       Expression  // nil for default
       Consequent []Statement
   }
   ```

4. **Parser** (`pkg/parser/parser.go`)
   - Add `parseSwitchStatement()`

5. **Type Checker** (`pkg/types/checker.go`)
   - Case values must be compatible with discriminant type

6. **Compiler** (`pkg/compiler/compiler.go`)
   - Generate as chain of if-else or jump table

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/ast/ast.go`
- `pkg/parser/parser.go`
- `pkg/types/checker.go`
- `pkg/compiler/compiler.go`

**Estimated effort:** 4-5 hours

---

### 10. For-of Loop

**Syntax:**
```typescript
for (let item of array) {
    println(item);
}

for (let char of string) {
    println(char);
}
```

**Implementation Steps:**

1. **Tokens** (`pkg/token/token.go`)
   - Add `OF`

2. **Lexer** (`pkg/lexer/lexer.go`)
   - Add `of` to keyword map

3. **AST** (`pkg/ast/ast.go`)
   ```go
   type ForOfStmt struct {
       Token    token.Token
       Variable *VarDecl
       Iterable Expression
       Body     *Block
   }
   ```

4. **Parser** (`pkg/parser/parser.go`)
   - Detect `for (let x of ...)` pattern

5. **Compiler** (`pkg/compiler/compiler.go`)
   - Desugar to index-based for loop:
     ```
     let __iter = iterable;
     let __i = 0;
     while (__i < len(__iter)) {
         let item = __iter[__i];
         // body
         __i = __i + 1;
     }
     ```

**Files to modify:**
- `pkg/token/token.go`
- `pkg/lexer/lexer.go`
- `pkg/ast/ast.go`
- `pkg/parser/parser.go`
- `pkg/types/checker.go`
- `pkg/compiler/compiler.go`

**Estimated effort:** 3-4 hours

---

## Low Priority Features

### 11. Try/Catch/Finally

**Syntax:**
```typescript
try {
    riskyOperation();
} catch (e: Error) {
    println(e.message);
} finally {
    cleanup();
}

throw new Error("message");
```

**Implementation Steps:**

1. **Tokens**: Add `TRY`, `CATCH`, `FINALLY`, `THROW`

2. **AST**: Add `TryStmt`, `ThrowStmt`

3. **Bytecode** (`pkg/bytecode/opcode.go`)
   - Add `OP_THROW`, `OP_SETUP_TRY`, `OP_POP_TRY`

4. **VM** (`pkg/vm/vm.go`)
   - Add exception frame stack
   - Implement unwinding on throw

**Estimated effort:** 8-12 hours

---

### 12. Template Literals

**Syntax:**
```typescript
let name = "world";
let greeting = `Hello, ${name}!`;
let multiline = `Line 1
Line 2`;
```

**Implementation Steps:**

1. **Lexer**: Handle backtick strings, parse `${...}` as embedded expressions

2. **AST**: Add `TemplateLiteral` with quasi strings and expressions

3. **Compiler**: Compile to string concatenation

**Estimated effort:** 4-6 hours

---

### 13. Static Class Members

**Syntax:**
```typescript
class Counter {
    static count: number = 0;

    static increment(): void {
        Counter.count = Counter.count + 1;
    }
}

Counter.increment();
println(Counter.count);
```

**Implementation Steps:**

1. **Tokens**: Add `STATIC`

2. **AST**: Add `IsStatic` flag to Field and Method

3. **Compiler**: Store static members on class object itself

4. **VM**: Handle property access on class values

**Estimated effort:** 4-5 hours

---

### 14. Interfaces

**Syntax:**
```typescript
interface Drawable {
    draw(): void;
    getArea(): number;
}

class Circle implements Drawable {
    radius: number

    constructor(radius: number) {
        this.radius = radius;
    }

    draw(): void {
        println("Drawing circle");
    }

    getArea(): number {
        return 3.14159 * this.radius * this.radius;
    }
}
```

**Implementation Steps:**

1. **Tokens**: Add `INTERFACE`, `IMPLEMENTS`

2. **AST**: Add `InterfaceDecl`

3. **Type Checker**: Structural subtyping check at class declaration

4. **Compiler**: Interfaces are compile-time only (no runtime representation)

**Estimated effort:** 6-8 hours

---

## Very Low Priority

### 15. Generics

**Syntax:**
```typescript
function identity<T>(value: T): T {
    return value;
}

class Box<T> {
    value: T

    constructor(value: T) {
        this.value = value;
    }

    get(): T {
        return this.value;
    }
}
```

**Implementation Steps:**

This is a complex feature requiring significant changes:

1. **Lexer/Parser**: Handle `<T>` in type positions
2. **Type System**: Type variables, constraints, instantiation
3. **Type Checker**: Type parameter substitution, inference
4. **Compiler**: Either monomorphization or type erasure

**Estimated effort:** 20-40 hours

---

## Implementation Order Recommendation

### Phase 1: Quick Wins (1-2 days)
1. **Type inference for variables** (compile-time only)
2. Compound assignment (`+=`, `-=`, etc.)
3. Arrow functions
4. Nullish coalescing (`??`)

### Phase 2: Core Improvements (3-5 days)
5. Increment/decrement (`++`, `--`)
6. String methods
7. Array methods
8. For-of loop

### Phase 3: Control Flow (2-3 days)
9. Switch statement
10. Optional chaining (`?.`)

### Phase 4: Advanced Features (1-2 weeks)
11. Static class members
12. Template literals
13. Interfaces
14. Try/catch/finally

### Phase 5: Future (TBD)
15. Generics

---

## Testing Strategy

For each feature:

1. **Unit tests** in corresponding `*_test.go` file
2. **Integration test** in `test/` directory with `.gts` file
3. **Error case tests** for type checker

Example test file structure:
```
test/
  arrow-functions.gts
  compound-assignment.gts
  string-methods.gts
  array-methods.gts
  ...
```

### Test Patterns

Based on the existing codebase, tests follow these patterns:

**Table-driven tests:**
```go
func TestVMFeature(t *testing.T) {
    tests := []struct {
        input    string
        expected float64  // or string, bool, etc.
    }{
        {"source code;", expectedValue},
        // ...
    }

    for _, tt := range tests {
        vm := runVM(t, tt.input)
        // assertions
    }
}
```

**Output capture tests:**
```go
func TestVMFeatureWithOutput(t *testing.T) {
    var buf bytes.Buffer
    vm := newVMWithOutput(t, `source code;`, &buf)
    err := vm.Run()
    // ...
    output := buf.String()
    // assertions on output
}
```

**Error tests:**
```go
func TestCheckerErrorCase(t *testing.T) {
    checker := newChecker(t, `invalid code;`)
    if len(checker.Errors()) == 0 {
        t.Error("expected type error")
    }
}
```

---

## Design Decisions & Rationale

### Why Explicit Function Signatures?

GoTS requires explicit type annotations for function parameters and return types, unlike TypeScript's inference. Rationale:

1. **Clarity at boundaries** - Functions are API contracts; explicit types document intent
2. **Simpler implementation** - No need for bidirectional type inference
3. **Faster type checking** - No need to analyze function body to determine parameter types
4. **Better error messages** - Errors point to explicit declarations, not inferred locations

### Why No `var` Keyword?

TypeScript supports `var`, `let`, and `const`. GoTS only supports `let` and `const`:

1. **Simplicity** - Two constructs sufficient for all use cases
2. **No hoisting** - `var` hoisting is confusing; `let`/`const` have cleaner semantics
3. **Block scoping only** - Consistent scoping rules throughout

### Why Value-Based Null (not Reference)?

GoTS uses `null` as a distinct value type, not a reference type:

1. **Explicit nullability** - `string | null` makes nullable types visible
2. **No billion-dollar mistake** - Can't accidentally pass null where not expected
3. **Simpler GC** - No special null reference handling

### Bytecode vs Tree-walking

GoTS compiles to bytecode rather than interpreting the AST directly:

1. **Performance** - Bytecode interpretation is ~10x faster than tree-walking
2. **Serialization** - Bytecode can be saved/loaded (`.gtsb` files)
3. **Separation of concerns** - Compiler and VM are independent

---

## Known Limitations

### Current Implementation Limits

| Limit | Value | Defined In |
|-------|-------|------------|
| Max local variables per scope | 256 | compiler.go |
| Max stack depth | 256 | vm.go |
| Max call frames | 64 | vm.go |
| GC initial threshold | 1MB | vm.go |

### Intentionally Unsupported Features

These TypeScript features are **not planned** for GoTS:

| Feature | Reason |
|---------|--------|
| `any` type | Defeats purpose of static typing |
| Type assertions (`as`) | Prefer explicit runtime checks |
| Non-null assertion (`!`) | Prefer explicit null handling |
| `typeof` type guards | Complex control flow analysis |
| Mapped types | Too complex for MVP |
| Conditional types | Too complex for MVP |
| `readonly` modifier | Can add later if needed |
| Enums | Use const objects or union types |
| Namespaces | Use modules when added |
| Decorators | Out of scope |

---

## Dependencies Between Features

Some features have implementation dependencies:

```
Arrow Functions
    └── (none, standalone)

Compound Assignment (+=, -=, etc.)
    └── (none, standalone)

Type Inference
    └── (none, compile-time only)

Array Methods
    └── Arrow Functions (for callbacks)
    └── Type Inference (for return types)

For-of Loop
    └── (none, standalone)

Optional Chaining (?.)
    └── Nullish Coalescing (??) (often used together)

Template Literals
    └── String Methods (for implementation)

Interfaces
    └── (none, compile-time only)

Try/Catch
    └── Classes (for Error class)

Generics
    └── Interfaces (for constraints)
    └── Type Inference (for type argument inference)
```

---

## Migration Notes

When implementing features, consider backwards compatibility:

### Safe to Add (Non-breaking)
- New syntax that was previously a parse error
- New built-in functions
- New operators
- New keywords (if not valid identifiers)

### Requires Care
- Changes to type checking rules
- Changes to operator semantics
- New reserved words that were valid identifiers

### Bytecode Versioning
- Current bytecode version: 1
- Increment version when adding new opcodes
- VM should reject incompatible versions

---

## Appendix: Quick Reference

### File Locations for Common Changes

| Change Type | Files to Modify |
|-------------|-----------------|
| New keyword | `token/token.go`, `lexer/lexer.go` |
| New operator | `token/token.go`, `lexer/lexer.go`, `parser/parser.go` |
| New AST node | `ast/ast.go` |
| New statement | `ast/ast.go`, `parser/parser.go`, `types/checker.go`, `compiler/compiler.go` |
| New expression | `ast/ast.go`, `parser/parser.go`, `types/checker.go`, `compiler/compiler.go` |
| New opcode | `bytecode/opcode.go`, `compiler/compiler.go`, `vm/vm.go` |
| New built-in | `vm/vm.go` (in `callNative`) |
| New type | `types/types.go`, `types/checker.go` |

### Parser Precedence Levels

```go
const (
    LOWEST      = 1
    ASSIGN      = 2   // =
    OR          = 3   // ||
    AND         = 4   // &&
    EQUALITY    = 5   // == !=
    COMPARISON  = 6   // < > <= >=
    SUM         = 7   // + -
    PRODUCT     = 8   // * / %
    PREFIX      = 9   // ! -
    CALL        = 10  // () [] .
)
```

When adding new operators:
- `??` should be between ASSIGN and OR (precedence ~2.5)
- `?.` should have same precedence as `.` (CALL level)

### Opcode Categories

| Category | Opcodes |
|----------|---------|
| Constants | `OP_CONSTANT`, `OP_TRUE`, `OP_FALSE`, `OP_NULL` |
| Arithmetic | `OP_ADD`, `OP_SUB`, `OP_MUL`, `OP_DIV`, `OP_MOD`, `OP_NEG` |
| Comparison | `OP_EQ`, `OP_NEQ`, `OP_LT`, `OP_LTE`, `OP_GT`, `OP_GTE` |
| Logical | `OP_NOT`, `OP_AND`, `OP_OR` |
| Variables | `OP_GET_GLOBAL`, `OP_SET_GLOBAL`, `OP_GET_LOCAL`, `OP_SET_LOCAL` |
| Control | `OP_JUMP`, `OP_JUMP_IF_FALSE`, `OP_LOOP` |
| Functions | `OP_CALL`, `OP_RETURN`, `OP_CLOSURE` |
| Objects | `OP_ARRAY`, `OP_OBJECT`, `OP_GET_PROP`, `OP_SET_PROP`, `OP_GET_INDEX`, `OP_SET_INDEX` |
| Classes | `OP_CLASS`, `OP_METHOD`, `OP_GET_SUPER`, `OP_INVOKE`, `OP_SUPER_INVOKE` |

---

## Changelog

| Date | Change |
|------|--------|
| 2024-01 | Initial plan document created |

---

*Last updated: January 2024*
