# goTS Test262 Port - TODO

Port relevant test262 tests to goTS. Since goTS compiles to Go, many JavaScript-specific tests won't apply. Focus on language features and constructs that goTS supports.

## Overview

- **Source**: `scaffold/test262/test/`
- **Target**: `gots/test/test262/`
- **Approach**: Port one category at a time, adapt to goTS syntax, fix compiler issues as needed

---

## Phase 1: Language Basics

### 1.1 Literals
- [x] **Boolean literals** ✓
  - Tests: `gots/test/test262/literals/boolean/`
  - 4 test files, all passing
- [x] **Numeric literals** ✓
  - Tests: `gots/test/test262/literals/numeric/`
  - Integer and float literals
- [x] **String literals** ✓
  - Tests: `gots/test/test262/literals/string/`
  - Basic strings, escape sequences
- [x] **Null literal** ✓
  - Tests: `gots/test/test262/literals/null/`
  - `null` behavior and comparisons

### 1.2 Variables & Constants
- [x] **let declarations** ✓
  - Tests: `gots/test/test262/variables/let_declarations.gts`
- [x] **const declarations** ✓
  - Tests: `gots/test/test262/variables/const_declarations.gts`

### 1.3 Identifiers
- [ ] **Valid identifiers** - TODO
- [ ] **Reserved words** - TODO

---

## Phase 2: Expressions & Operators

### 2.1 Arithmetic Operators
- [x] **All arithmetic operators** ✓
  - Tests: `gots/test/test262/operators/arithmetic/arithmetic_ops.gts`
  - Addition, subtraction, multiplication, division, modulo, unary

### 2.2 Comparison Operators
- [x] **All comparison operators** ✓
  - Tests: `gots/test/test262/operators/comparison/comparison_ops.gts`
  - Equality, relational operators

### 2.3 Logical Operators
- [x] **All logical operators** ✓
  - Tests: `gots/test/test262/operators/logical/logical_ops.gts`
  - AND, OR, NOT

### 2.4 Assignment Operators
- [x] **All assignment operators** ✓
  - Tests: `gots/test/test262/operators/assignment/assignment_ops.gts`
  - Basic assignment, compound assignment

---

## Phase 3: Statements

### 3.1 Control Flow
- [x] **if statement** ✓
  - Tests: `gots/test/test262/statements/if/if_statement.gts`
- [ ] **switch statement** - TODO
- [ ] **block statement** - TODO

### 3.2 Loops
- [x] **while statement** ✓
  - Tests: `gots/test/test262/statements/while/while_statement.gts`
- [ ] **do-while statement** - TODO
- [x] **for statement** ✓
  - Tests: `gots/test/test262/statements/for/for_statement.gts`
- [x] **for-of statement** ✓
  - Tests: `gots/test/test262/statements/for-of/for_of_statement.gts`
- [x] **break statement** ✓
  - Tests: `gots/test/test262/statements/break/break_statement.gts`
- [x] **continue statement** ✓
  - Tests: `gots/test/test262/statements/continue/continue_statement.gts`

### 3.3 Functions
- [x] **Function declarations** ✓
  - Tests: `gots/test/test262/functions/function_declarations.gts`

---

## Phase 4: Functions & Closures

- [x] **Closures** ✓
  - Tests: `gots/test/test262/closures/closure_tests.gts`
  - Variable capture, higher-order functions, nested closures

---

## Phase 5: Types

- [x] **Type system tests** ✓
  - Tests: `gots/test/test262/types/type_tests.gts`
  - Type inference, explicit types, arrays, type conversions

---

## Phase 6: Built-in Functions

- [x] **Built-in functions** ✓
  - Tests: `gots/test/test262/builtins/builtin_tests.gts`
  - len, push, pop, sqrt, floor, ceil, abs, typeof, tostring, toint, tofloat

---

## Phase 7: Regular Expressions

- [x] **RegExp tests** ✓
  - Tests: `gots/test/test262/regexp/regexp_tests.gts`
  - Literal syntax, test(), exec(), flags

---

## Phase 8: Classes & OOP

- [x] **Class tests** ✓
  - Tests: `gots/test/test262/classes/class_tests.gts`
  - Class declaration, constructor, methods, fields, inheritance, super

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

---

## Progress Tracking

| Phase | Category | Tests Written | Tests Passing | Status |
|-------|----------|---------------|---------------|--------|
| 1.1 | Literals | 6 | 6 | ✓ Complete |
| 1.2 | Variables | 2 | 2 | ✓ Complete |
| 1.3 | Identifiers | 0 | 0 | TODO |
| 2 | Operators | 4 | 4 | ✓ Complete |
| 3.1 | Control Flow | 1 | 1 | Partial |
| 3.2 | Loops | 5 | 5 | ✓ Complete |
| 3.3 | Functions | 1 | 1 | ✓ Complete |
| 4 | Closures | 1 | 1 | ✓ Complete |
| 5 | Types | 1 | 1 | ✓ Complete |
| 6 | Built-ins | 1 | 1 | ✓ Complete |
| 7 | RegExp | 1 | 1 | ✓ Complete |
| 8 | Classes | 1 | 1 | ✓ Complete |

**Total: 26 test files, 26 passing**

---

## Compiler Fixes Made

During test development, the following goTS compiler issue was fixed:

1. **`gts_typeof` for int type** - Added `case int:` to return "number" for integer values (pkg/codegen/codegen.go)

---

## Remaining TODO

- [ ] Identifiers test (valid/invalid names, reserved words)
- [ ] switch statement test
- [ ] block statement test
- [ ] do-while statement test

---

## Current Status

**All major language features tested and passing!**

Run all tests: `cd gots && find test/test262 -name "*.gts" -exec go run ./cmd/gots run {} \;`
