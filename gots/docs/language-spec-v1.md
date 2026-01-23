# GoTS Language Specification v1.0

This document describes the GoTS language specification and serves as an implementation plan for features not yet implemented.

## Overview

GoTS is a TypeScript-like language that compiles to bytecode and runs on a stack-based virtual machine. It supports static typing, functions, closures, classes with inheritance, and garbage collection.

---

## Table of Contents

1. [Lexical Elements](#1-lexical-elements)
2. [Types](#2-types)
3. [Expressions](#3-expressions)
4. [Statements](#4-statements)
5. [Functions](#5-functions)
6. [Classes](#6-classes)
7. [Built-in Functions](#7-built-in-functions)
8. [Implementation Plan for New Features](#8-implementation-plan-for-new-features)

---

## 1. Lexical Elements

### 1.1 Comments (✅ Implemented)

```typescript
// Single-line comment
```

### 1.2 Identifiers (✅ Implemented)

Identifiers start with a letter or underscore and can contain letters, digits, and underscores.

### 1.3 Keywords (✅ Implemented)

```
let       const     function  return    if        else
while     for       break     continue  class     extends
new       this      super     constructor         type
true      false     null
```

### 1.4 Type Keywords (✅ Implemented)

```
number    string    boolean   void      null
```

### 1.5 Operators (✅ Implemented)

| Category | Operators |
|----------|-----------|
| Arithmetic | `+` `-` `*` `/` `%` |
| Comparison | `==` `!=` `<` `>` `<=` `>=` |
| Logical | `&&` `\|\|` `!` |
| Assignment | `=` |

### 1.6 Delimiters (✅ Implemented)

```
( ) { } [ ] ; : , . => |
```

---

## 2. Types

### 2.1 Primitive Types (✅ Implemented)

| Type | Description | Example |
|------|-------------|---------|
| `number` | 64-bit floating point | `42`, `3.14` |
| `string` | UTF-8 string | `"hello"` |
| `boolean` | Boolean value | `true`, `false` |
| `void` | No value | Used for functions |
| `null` | Null value | `null` |

### 2.2 Array Types (✅ Implemented)

```typescript
let arr: number[] = [1, 2, 3];
let names: string[] = ["a", "b"];
```

### 2.3 Object Types (✅ Implemented)

```typescript
let point: { x: number, y: number } = { x: 10, y: 20 };
```

### 2.4 Function Types (✅ Implemented)

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};
```

### 2.5 Nullable Types (✅ Implemented)

```typescript
let name: string | null = null;
```

### 2.6 Type Aliases (✅ Implemented)

```typescript
type Point = { x: number, y: number };
let p: Point = { x: 0, y: 0 };
```

### 2.7 Class Types (✅ Implemented)

```typescript
class Animal { ... }
let a: Animal = new Animal();
```

---

## 3. Expressions

### 3.1 Literals (✅ Implemented)

```typescript
42              // number
3.14            // number
"hello"         // string
true            // boolean
false           // boolean
null            // null
[1, 2, 3]       // array
{ x: 1, y: 2 }  // object
```

### 3.2 Binary Expressions (✅ Implemented)

```typescript
a + b    // addition, string concatenation
a - b    // subtraction
a * b    // multiplication
a / b    // division
a % b    // modulo
a == b   // equality
a != b   // inequality
a < b    // less than
a > b    // greater than
a <= b   // less or equal
a >= b   // greater or equal
a && b   // logical and
a || b   // logical or
```

### 3.3 Unary Expressions (✅ Implemented)

```typescript
-x       // negation
!flag    // logical not
```

### 3.4 Function Calls (✅ Implemented)

```typescript
foo(1, 2)
obj.method(arg)
```

### 3.5 Property Access (✅ Implemented)

```typescript
obj.property
```

### 3.6 Index Access (✅ Implemented)

```typescript
arr[0]
str[i]
```

### 3.7 Assignment (✅ Implemented)

```typescript
x = 10
arr[0] = 5
obj.prop = value
```

### 3.8 New Expression (✅ Implemented)

```typescript
new ClassName(args)
```

### 3.9 This/Super (✅ Implemented)

```typescript
this.property
super(args)
```

---

## 4. Statements

### 4.1 Variable Declaration (✅ Implemented)

```typescript
let x: number = 10;
const name: string = "GoTS";
```

### 4.2 Block Statement (✅ Implemented)

```typescript
{
    statement1;
    statement2;
}
```

### 4.3 If Statement (✅ Implemented)

```typescript
if (condition) {
    // then branch
} else {
    // else branch
}
```

### 4.4 While Statement (✅ Implemented)

```typescript
while (condition) {
    // body
}
```

### 4.5 For Statement (✅ Implemented)

```typescript
for (let i: number = 0; i < 10; i = i + 1) {
    // body
}
```

### 4.6 Return Statement (✅ Implemented)

```typescript
return value;
return;
```

### 4.7 Break/Continue (✅ Implemented)

```typescript
break;
continue;
```

---

## 5. Functions

### 5.1 Function Declaration (✅ Implemented)

```typescript
function add(a: number, b: number): number {
    return a + b;
}
```

### 5.2 Function Expression (✅ Implemented)

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};
```

### 5.3 Closures (✅ Implemented)

```typescript
function makeCounter(): () => number {
    let count: number = 0;
    return function(): number {
        count = count + 1;
        return count;
    };
}
```

---

## 6. Classes

### 6.1 Class Declaration (✅ Implemented)

```typescript
class Point {
    x: number
    y: number

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    distance(): number {
        return sqrt(this.x * this.x + this.y * this.y);
    }
}
```

### 6.2 Inheritance (✅ Implemented)

```typescript
class Animal {
    name: string
    constructor(name: string) {
        this.name = name;
    }
}

class Dog extends Animal {
    constructor(name: string) {
        super(name);
    }
}
```

---

## 7. Built-in Functions

### 7.1 I/O (✅ Implemented)

| Function | Signature | Description |
|----------|-----------|-------------|
| `println` | `(value: any) => void` | Print with newline |
| `print` | `(value: any) => void` | Print without newline |

### 7.2 Array Operations (✅ Implemented)

| Function | Signature | Description |
|----------|-----------|-------------|
| `len` | `(arr: T[]) => number` | Array/string length |
| `push` | `(arr: T[], val: T) => number` | Append to array |
| `pop` | `(arr: T[]) => T` | Remove last element |

### 7.3 Type Conversion (✅ Implemented)

| Function | Signature | Description |
|----------|-----------|-------------|
| `typeof` | `(val: any) => string` | Get type name |
| `tostring` | `(val: any) => string` | Convert to string |
| `tonumber` | `(val: any) => number` | Convert to number |

### 7.4 Math Functions (✅ Implemented)

| Function | Signature | Description |
|----------|-----------|-------------|
| `sqrt` | `(n: number) => number` | Square root |
| `floor` | `(n: number) => number` | Floor |
| `ceil` | `(n: number) => number` | Ceiling |
| `abs` | `(n: number) => number` | Absolute value |

---

## 8. Implementation Plan for New Features

The following features are **not yet implemented** and are candidates for future development.

### 8.1 Arrow Functions (🔴 Not Implemented)

**Priority: High**

```typescript
// Desired syntax
let add = (a: number, b: number): number => a + b;
let double = (x: number): number => { return x * 2; };
```

**Implementation Plan:**
1. **Lexer**: Already has `ARROW` token (`=>`)
2. **Parser**: Add `parseArrowFunction()` - detect `(params) =>` pattern
3. **AST**: Create `ArrowFunctionExpr` node or reuse `FunctionExpr`
4. **Compiler**: Same as function expression

**Files to modify:**
- `pkg/parser/parser.go`: Add arrow function parsing in expression parser
- `pkg/ast/ast.go`: Optional - can reuse FunctionExpr

---

### 8.2 Compound Assignment Operators (🔴 Not Implemented)

**Priority: High**

```typescript
x += 5;    // x = x + 5
x -= 3;    // x = x - 3
x *= 2;    // x = x * 2
x /= 4;    // x = x / 4
x %= 3;    // x = x % 3
```

**Implementation Plan:**
1. **Lexer**: Add tokens `PLUS_ASSIGN`, `MINUS_ASSIGN`, `STAR_ASSIGN`, `SLASH_ASSIGN`, `PERCENT_ASSIGN`
2. **Parser**: Parse as syntactic sugar, desugar to `x = x + value`
3. **Type Checker**: Already handles assignment
4. **Compiler**: Same as desugared form

**Files to modify:**
- `pkg/token/token.go`: Add 5 new tokens
- `pkg/lexer/lexer.go`: Recognize `+=`, `-=`, `*=`, `/=`, `%=`
- `pkg/parser/parser.go`: Desugar in assignment parsing

---

### 8.3 Increment/Decrement Operators (🔴 Not Implemented)

**Priority: Medium**

```typescript
i++;    // post-increment
++i;    // pre-increment
i--;    // post-decrement
--i;    // pre-decrement
```

**Implementation Plan:**
1. **Lexer**: Add tokens `PLUS_PLUS`, `MINUS_MINUS`
2. **AST**: Add `UpdateExpr` node with `prefix` boolean
3. **Parser**: Handle prefix in unary, postfix in postfix position
4. **Compiler**: Generate appropriate load/add/store sequence

**Files to modify:**
- `pkg/token/token.go`: Add `PLUS_PLUS`, `MINUS_MINUS`
- `pkg/lexer/lexer.go`: Recognize `++`, `--`
- `pkg/ast/ast.go`: Add `UpdateExpr`
- `pkg/parser/parser.go`: Parse pre/post increment
- `pkg/compiler/compiler.go`: Generate bytecode

---

### 8.4 Optional Chaining (🔴 Not Implemented)

**Priority: Medium**

```typescript
let name = user?.profile?.name;
let result = obj?.method?.(arg);
```

**Implementation Plan:**
1. **Lexer**: Add token `QUESTION_DOT` for `?.`
2. **AST**: Add `OptionalChainExpr` or flag on `PropertyExpr`
3. **Parser**: Handle `?.` in property/call parsing
4. **Compiler**: Generate null check + conditional jump

**Files to modify:**
- `pkg/token/token.go`: Add `QUESTION_DOT`
- `pkg/lexer/lexer.go`: Recognize `?.`
- `pkg/ast/ast.go`: Modify `PropertyExpr` or add new node
- `pkg/compiler/compiler.go`: Generate null-check code

---

### 8.5 Nullish Coalescing (🔴 Not Implemented)

**Priority: Medium**

```typescript
let value = maybeNull ?? defaultValue;
```

**Implementation Plan:**
1. **Lexer**: Add token `NULLISH_COALESCE` for `??`
2. **Parser**: Add precedence between `||` and `?:`
3. **Compiler**: Generate null check, short-circuit evaluation

**Files to modify:**
- `pkg/token/token.go`: Add `NULLISH_COALESCE`
- `pkg/lexer/lexer.go`: Recognize `??`
- `pkg/parser/parser.go`: Add to binary expression parsing
- `pkg/compiler/compiler.go`: Generate bytecode

---

### 8.6 Switch Statement (🔴 Not Implemented)

**Priority: Medium**

```typescript
switch (value) {
    case 1:
        println("one");
        break;
    case 2:
        println("two");
        break;
    default:
        println("other");
}
```

**Implementation Plan:**
1. **Lexer**: Add tokens `SWITCH`, `CASE`, `DEFAULT`
2. **AST**: Add `SwitchStmt`, `CaseClause`
3. **Parser**: Parse switch structure
4. **Type Checker**: Check case values match switch expression type
5. **Compiler**: Generate jump table or if-else chain

**Files to modify:**
- `pkg/token/token.go`: Add `SWITCH`, `CASE`, `DEFAULT`
- `pkg/lexer/lexer.go`: Recognize keywords
- `pkg/ast/ast.go`: Add `SwitchStmt`, `CaseClause`
- `pkg/parser/parser.go`: Add `parseSwitchStatement()`
- `pkg/types/checker.go`: Add case type checking
- `pkg/compiler/compiler.go`: Generate switch bytecode

---

### 8.7 Try/Catch/Finally (🔴 Not Implemented)

**Priority: Low**

```typescript
try {
    riskyOperation();
} catch (e: Error) {
    println(e.message);
} finally {
    cleanup();
}
```

**Implementation Plan:**
1. **Lexer**: Add tokens `TRY`, `CATCH`, `FINALLY`, `THROW`
2. **AST**: Add `TryStmt`, `ThrowStmt`
3. **VM**: Add exception handling infrastructure
4. **Bytecode**: Add `OP_THROW`, `OP_SETUP_EXCEPT`, `OP_POP_EXCEPT`

**Files to modify:**
- `pkg/token/token.go`: Add exception tokens
- `pkg/ast/ast.go`: Add `TryStmt`, `ThrowStmt`
- `pkg/parser/parser.go`: Parse try/catch/finally
- `pkg/bytecode/opcode.go`: Add exception opcodes
- `pkg/vm/vm.go`: Add exception frame handling

---

### 8.8 Template Literals (🔴 Not Implemented)

**Priority: Low**

```typescript
let name = "world";
let greeting = `Hello, ${name}!`;
```

**Implementation Plan:**
1. **Lexer**: Handle backtick strings with `${...}` interpolation
2. **AST**: Add `TemplateLiteral` with parts and expressions
3. **Compiler**: Compile to string concatenation

**Files to modify:**
- `pkg/lexer/lexer.go`: Parse template literals
- `pkg/ast/ast.go`: Add `TemplateLiteral`
- `pkg/parser/parser.go`: Parse template parts
- `pkg/compiler/compiler.go`: Generate string concat code

---

### 8.9 Static Class Members (🔴 Not Implemented)

**Priority: Low**

```typescript
class Math {
    static PI: number = 3.14159;
    static max(a: number, b: number): number {
        if (a > b) { return a; }
        return b;
    }
}
```

**Implementation Plan:**
1. **Lexer**: Add `STATIC` token
2. **AST**: Add `IsStatic` flag to `Field` and `Method`
3. **Type Checker**: Handle static member access
4. **Compiler**: Store static members on class object

**Files to modify:**
- `pkg/token/token.go`: Add `STATIC`
- `pkg/ast/ast.go`: Add `IsStatic` to Field/Method
- `pkg/types/checker.go`: Handle `ClassName.staticMember`
- `pkg/compiler/compiler.go`: Generate static member code

---

### 8.10 Interfaces (🔴 Not Implemented)

**Priority: Low**

```typescript
interface Drawable {
    draw(): void;
}

class Circle implements Drawable {
    draw(): void { ... }
}
```

**Implementation Plan:**
1. **Lexer**: Add `INTERFACE`, `IMPLEMENTS` tokens
2. **AST**: Add `InterfaceDecl`
3. **Type Checker**: Structural subtyping check
4. **Compiler**: Interfaces are compile-time only

---

### 8.11 Generics (🔴 Not Implemented)

**Priority: Very Low** (Complex feature)

```typescript
function identity<T>(value: T): T {
    return value;
}
```

**Implementation Plan:**
1. **Lexer**: Handle `<>` in type contexts
2. **AST**: Add type parameters to functions/classes
3. **Type Checker**: Implement type parameter substitution
4. **Compiler**: Monomorphization or type erasure

---

### 8.12 For-of Loop (🔴 Not Implemented)

**Priority: Medium**

```typescript
for (let item of array) {
    println(item);
}
```

**Implementation Plan:**
1. **Lexer**: Add `OF` token
2. **AST**: Add `ForOfStmt`
3. **Compiler**: Desugar to index-based for loop

---

### 8.13 More Array Methods (🔴 Not Implemented)

**Priority: High**

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
```

**Implementation Plan:**
1. Add as built-in methods on arrays
2. **VM**: Implement method dispatch for arrays
3. **Type Checker**: Add array method types

---

### 8.14 String Methods (🔴 Not Implemented)

**Priority: High**

```typescript
str.toUpperCase();
str.toLowerCase();
str.split(delimiter);
str.trim();
str.substring(start, end);
str.indexOf(substr);
str.replace(old, new);
str.startsWith(prefix);
str.endsWith(suffix);
```

**Implementation Plan:**
1. Add as built-in methods on strings
2. **VM**: Implement method dispatch for strings
3. **Type Checker**: Add string method types

---

## Implementation Priority Summary

| Priority | Feature | Complexity |
|----------|---------|------------|
| 🔴 High | Arrow functions | Low |
| 🔴 High | Compound assignment (`+=`, `-=`) | Low |
| 🔴 High | Array methods (map, filter, etc.) | Medium |
| 🔴 High | String methods | Medium |
| 🟡 Medium | Increment/decrement (`++`, `--`) | Low |
| 🟡 Medium | Optional chaining (`?.`) | Medium |
| 🟡 Medium | Nullish coalescing (`??`) | Low |
| 🟡 Medium | Switch statement | Medium |
| 🟡 Medium | For-of loop | Low |
| 🟢 Low | Try/catch/finally | High |
| 🟢 Low | Template literals | Medium |
| 🟢 Low | Static class members | Medium |
| 🟢 Low | Interfaces | Medium |
| ⚪ Very Low | Generics | Very High |

---

## Appendix: Bytecode Reference

### Current Opcodes

| Opcode | Hex | Description |
|--------|-----|-------------|
| OP_CONSTANT | 0x01 | Push constant |
| OP_NULL | 0x02 | Push null |
| OP_TRUE | 0x03 | Push true |
| OP_FALSE | 0x04 | Push false |
| OP_ADD | 0x10 | Add |
| OP_SUBTRACT | 0x11 | Subtract |
| OP_MULTIPLY | 0x12 | Multiply |
| OP_DIVIDE | 0x13 | Divide |
| OP_MODULO | 0x14 | Modulo |
| OP_NEGATE | 0x15 | Negate |
| OP_EQUAL | 0x20 | Equal |
| OP_NOT_EQUAL | 0x21 | Not equal |
| OP_LESS | 0x22 | Less than |
| OP_LESS_EQUAL | 0x23 | Less or equal |
| OP_GREATER | 0x24 | Greater than |
| OP_GREATER_EQUAL | 0x25 | Greater or equal |
| OP_NOT | 0x30 | Logical not |
| OP_CONCAT | 0x40 | String concat |
| OP_GET_LOCAL | 0x50 | Get local |
| OP_SET_LOCAL | 0x51 | Set local |
| OP_GET_GLOBAL | 0x52 | Get global |
| OP_SET_GLOBAL | 0x53 | Set global |
| OP_GET_UPVALUE | 0x54 | Get upvalue |
| OP_SET_UPVALUE | 0x55 | Set upvalue |
| OP_POP | 0x60 | Pop stack |
| OP_POPN | 0x61 | Pop n values |
| OP_DUP | 0x62 | Duplicate |
| OP_JUMP | 0x70 | Jump |
| OP_JUMP_BACK | 0x71 | Jump back |
| OP_JUMP_IF_FALSE | 0x72 | Conditional jump |
| OP_JUMP_IF_TRUE | 0x73 | Conditional jump |
| OP_CALL | 0x80 | Call function |
| OP_RETURN | 0x81 | Return |
| OP_CLOSURE | 0x82 | Create closure |
| OP_CLASS | 0x90 | Define class |
| OP_GET_PROPERTY | 0x91 | Get property |
| OP_SET_PROPERTY | 0x92 | Set property |
| OP_METHOD | 0x93 | Define method |
| OP_INVOKE | 0x94 | Invoke method |
| OP_INHERIT | 0x95 | Set inheritance |
| OP_GET_SUPER | 0x96 | Get super method |
| OP_SUPER_INVOKE | 0x97 | Invoke super |
| OP_ARRAY | 0xA0 | Create array |
| OP_GET_INDEX | 0xA1 | Get index |
| OP_SET_INDEX | 0xA2 | Set index |
| OP_OBJECT | 0xB0 | Create object |
| OP_CLOSE_UPVALUE | 0xC0 | Close upvalue |
| OP_PRINT | 0xD0 | Print |
| OP_PRINTLN | 0xD1 | Print line |
| OP_BUILTIN | 0xE0 | Call builtin |
