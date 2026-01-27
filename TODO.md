# goTS Test262 Port - TODO

Port relevant test262 tests to goTS. Since goTS compiles to Go, many JavaScript-specific tests won't apply. Focus on language features and constructs that goTS supports.

## Overview

- **Source**: `scaffold/test262/test/`
- **Target**: `gots/test/test262/`
- **Approach**: Port one category at a time, adapt to goTS syntax, fix compiler issues as needed

---

## Phase 1: Language Basics

### 1.1 Literals
- [x] **Boolean literals** (`test262/test/language/literals/boolean/`) ✓
  - `true`, `false` literal tests
  - Tests: `gots/test/test262/literals/boolean/`
- [ ] **Numeric literals** (`test262/test/language/literals/numeric/`)
  - Integer literals, float literals
  - Adapt: goTS has `number`, `int`, `float` types
- [ ] **String literals** (`test262/test/language/literals/string/`)
  - Basic strings, escape sequences
- [ ] **Null literal** (`test262/test/language/literals/null/`)
  - `null` behavior
- [ ] **RegExp literals** (`test262/test/language/literals/regexp/`)
  - Basic regex patterns, flags

### 1.2 Variables & Constants
- [ ] **let declarations** (`test262/test/language/statements/let/`)
  - Block scoping, initialization
- [ ] **const declarations** (`test262/test/language/statements/const/`)
  - Immutability, block scoping
- [ ] **Variable declarations** (`test262/test/language/statements/variable/`)
  - Basic `var` behavior (if supported)

### 1.3 Identifiers
- [ ] **Valid identifiers** (`test262/test/language/identifiers/`)
  - Valid/invalid identifier names
- [ ] **Reserved words** (`test262/test/language/reserved-words/`)
  - Keywords that can't be identifiers

---

## Phase 2: Expressions & Operators

### 2.1 Arithmetic Operators
- [ ] **Addition** (`+`)
- [ ] **Subtraction** (`-`)
- [ ] **Multiplication** (`*`)
- [ ] **Division** (`/`)
- [ ] **Modulo** (`%`)
- [ ] **Unary operators** (`+`, `-`, `!`)
- [ ] **Increment/Decrement** (`++`, `--`)

### 2.2 Comparison Operators
- [ ] **Equality** (`==`, `!=`)
- [ ] **Strict equality** (`===`, `!==`) - may map to `==`/`!=` in goTS
- [ ] **Relational** (`<`, `>`, `<=`, `>=`)

### 2.3 Logical Operators
- [ ] **AND** (`&&`)
- [ ] **OR** (`||`)
- [ ] **NOT** (`!`)

### 2.4 Assignment Operators
- [ ] **Basic assignment** (`=`)
- [ ] **Compound assignment** (`+=`, `-=`, `*=`, `/=`, `%=`)

### 2.5 Other Expressions
- [ ] **Ternary operator** (`? :`)
- [ ] **Grouping** (`()`)
- [ ] **Member access** (`.`, `[]`)
- [ ] **Function calls** (`()`)

---

## Phase 3: Statements

### 3.1 Control Flow
- [ ] **if statement** (`test262/test/language/statements/if/`)
  - if, if-else, if-else-if chains
- [ ] **switch statement** (`test262/test/language/statements/switch/`)
  - switch, case, default, break
- [ ] **block statement** (`test262/test/language/statements/block/`)
  - Scoping behavior

### 3.2 Loops
- [ ] **while statement** (`test262/test/language/statements/while/`)
- [ ] **do-while statement** (`test262/test/language/statements/do-while/`)
- [ ] **for statement** (`test262/test/language/statements/for/`)
- [ ] **for-of statement** (`test262/test/language/statements/for-of/`)
  - Array iteration
- [ ] **break statement** (`test262/test/language/statements/break/`)
- [ ] **continue statement** (`test262/test/language/statements/continue/`)

### 3.3 Functions
- [ ] **Function declarations** (`test262/test/language/statements/function/`)
  - Named functions, parameters, return
- [ ] **return statement** (`test262/test/language/statements/return/`)

### 3.4 Classes
- [ ] **Class declarations** (`test262/test/language/statements/class/`)
  - Class syntax, constructor, methods
  - Fields, inheritance (extends)
  - Method overriding, super calls

---

## Phase 4: Functions & Closures

### 4.1 Function Basics
- [ ] **Function parameters** (`test262/test/language/function-code/`)
  - Parameter handling, arity
- [ ] **Return values**
  - Explicit return, void functions
- [ ] **Recursion**
  - Recursive function calls

### 4.2 Closures
- [ ] **Variable capture**
  - Closures capturing outer scope
- [ ] **Higher-order functions**
  - Functions returning functions
  - Functions as parameters

### 4.3 Arrow Functions (if supported)
- [ ] **Arrow function syntax**
  - `(x) => x + 1`

---

## Phase 5: Types

### 5.1 Primitive Types
- [ ] **number type**
  - Default numeric type
- [ ] **int type**
  - Integer-specific operations
- [ ] **float type**
  - Float-specific operations
- [ ] **string type**
- [ ] **boolean type**
- [ ] **void type**

### 5.2 Complex Types
- [ ] **Arrays**
  - Array literals, indexing, length
  - `push`, `pop` operations
- [ ] **Objects**
  - Object literals, property access
- [ ] **Nullable types**
  - `T | null` handling

### 5.3 Type Operations
- [ ] **typeof operator**
- [ ] **Type coercion**
  - `tostring()`, `toint()`, `tofloat()`

---

## Phase 6: Built-in Functions

### 6.1 Output
- [ ] **println**
- [ ] **print**

### 6.2 Array Operations
- [ ] **len**
- [ ] **push**
- [ ] **pop**

### 6.3 Math Operations
- [ ] **sqrt**
- [ ] **floor**
- [ ] **ceil**
- [ ] **abs**

### 6.4 Type Conversions
- [ ] **tostring**
- [ ] **toint**
- [ ] **tofloat**
- [ ] **typeof**

---

## Phase 7: Regular Expressions

- [ ] **RegExp creation**
  - Literal syntax `/pattern/flags`
- [ ] **RegExp.test()**
  - Boolean match testing
- [ ] **RegExp.exec()**
  - Match extraction
- [ ] **Flags**
  - `i`, `m`, `s`, `g` flags

---

## Phase 8: Classes & OOP

### 8.1 Class Basics
- [ ] **Class declaration**
- [ ] **Constructor**
- [ ] **Instance methods**
- [ ] **Instance fields**
- [ ] **this keyword**

### 8.2 Inheritance
- [ ] **extends keyword**
- [ ] **super() calls**
- [ ] **Method overriding**
- [ ] **Inherited fields**

---

## Excluded Tests (Not Applicable to goTS)

The following JavaScript features are NOT supported in goTS and their tests should be skipped:

- **Async/Await** - goTS doesn't support async
- **Generators** - No generator functions
- **Promises** - No Promise API
- **Symbols** - No Symbol type
- **Proxy/Reflect** - No metaprogramming
- **WeakMap/WeakSet** - No weak references
- **Intl (i18n)** - No internationalization
- **eval()** - No runtime eval
- **with statement** - Not supported
- **Destructuring** - Not yet supported
- **Spread operator** - Not yet supported
- **Template literals** - Not yet supported
- **Optional chaining** - Not yet supported
- **Nullish coalescing** - Not yet supported
- **BigInt** - Not supported
- **SharedArrayBuffer/Atomics** - Not supported
- **Most built-in object methods** (Array.map, Array.filter, etc.)

---

## Progress Tracking

| Phase | Category | Tests Written | Tests Passing | Notes |
|-------|----------|---------------|---------------|-------|
| 1.1 | Literals | 4 | 4 | Boolean complete |
| 1.2 | Variables | 0 | 0 | |
| 1.3 | Identifiers | 0 | 0 | |
| 2.1 | Arithmetic | 0 | 0 | |
| 2.2 | Comparison | 0 | 0 | |
| 2.3 | Logical | 0 | 0 | |
| 2.4 | Assignment | 0 | 0 | |
| 2.5 | Expressions | 0 | 0 | |
| 3.1 | Control Flow | 0 | 0 | |
| 3.2 | Loops | 0 | 0 | |
| 3.3 | Functions | 0 | 0 | |
| 3.4 | Classes | 0 | 0 | |
| 4 | Closures | 0 | 0 | |
| 5 | Types | 0 | 0 | |
| 6 | Built-ins | 0 | 0 | |
| 7 | RegExp | 0 | 0 | |
| 8 | OOP | 0 | 0 | |

---

## Workflow

1. **Pick a category** from the TODO above
2. **Create test file** in `gots/test/test262/<category>/`
3. **Run test**: `gots run <test_file>.gts`
4. **If fails**: Fix goTS compiler/codegen
5. **Mark complete** in progress table
6. **Repeat**

---

## Current Focus

**Completed**: Phase 1.1 - Boolean literals (4/4 tests passing)

**Next up**: Phase 1.1 - Numeric literals

Start with the simplest tests and progressively tackle more complex features.
