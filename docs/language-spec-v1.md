# GoTS Language Specification v1.0

## Overview

GoTS is a statically-typed language with TypeScript-like syntax that compiles to bytecode and runs on a stack-based virtual machine. It features static typing, first-class functions, closures, classes with single inheritance, and automatic garbage collection.

---

## Table of Contents

1. [Lexical Elements](#1-lexical-elements)
2. [Types](#2-types)
3. [Expressions](#3-expressions)
4. [Statements](#4-statements)
5. [Functions](#5-functions)
6. [Classes](#6-classes)
7. [Built-in Functions](#7-built-in-functions)
8. [Bytecode Reference](#8-bytecode-reference)

---

## 1. Lexical Elements

### 1.1 Comments

```typescript
// Single-line comment
```

### 1.2 Identifiers

Identifiers start with a letter or underscore and can contain letters, digits, and underscores.

### 1.3 Keywords

```
let       const     function  return    if        else
while     for       break     continue  class     extends
new       this      super     constructor         type
true      false     null
```

### 1.4 Type Keywords

```
number    string    boolean   void      null
```

### 1.5 Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+` `-` `*` `/` `%` |
| Comparison | `==` `!=` `<` `>` `<=` `>=` |
| Logical | `&&` `\|\|` `!` |
| Assignment | `=` |

### 1.6 Delimiters

```
( ) { } [ ] ; : , . => |
```

---

## 2. Types

### 2.1 Primitive Types

| Type | Description | Example |
|------|-------------|---------|
| `number` | 64-bit floating point | `42`, `3.14` |
| `string` | UTF-8 string | `"hello"` |
| `boolean` | Boolean value | `true`, `false` |
| `void` | No value | Used for functions |
| `null` | Null value | `null` |

### 2.2 Array Types

```typescript
let arr: number[] = [1, 2, 3];
let names: string[] = ["a", "b"];
```

### 2.3 Object Types

```typescript
let point: { x: number, y: number } = { x: 10, y: 20 };
```

### 2.4 Function Types

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};
```

### 2.5 Nullable Types

```typescript
let name: string | null = null;
```

### 2.6 Type Aliases

```typescript
type Point = { x: number, y: number };
let p: Point = { x: 0, y: 0 };
```

### 2.7 Class Types

```typescript
class Animal { ... }
let a: Animal = new Animal();
```

---

## 3. Expressions

### 3.1 Literals

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

### 3.2 Binary Expressions

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

### 3.3 Unary Expressions

```typescript
-x       // negation
!flag    // logical not
```

### 3.4 Function Calls

```typescript
foo(1, 2)
obj.method(arg)
```

### 3.5 Property Access

```typescript
obj.property
```

### 3.6 Index Access

```typescript
arr[0]
str[i]
```

### 3.7 Assignment

```typescript
x = 10
arr[0] = 5
obj.prop = value
```

### 3.8 New Expression

```typescript
new ClassName(args)
```

### 3.9 This/Super

```typescript
this.property
super(args)
```

---

## 4. Statements

### 4.1 Variable Declaration

```typescript
let x: number = 10;
const name: string = "GoTS";
```

### 4.2 Block Statement

```typescript
{
    statement1;
    statement2;
}
```

### 4.3 If Statement

```typescript
if (condition) {
    // then branch
} else {
    // else branch
}
```

### 4.4 While Statement

```typescript
while (condition) {
    // body
}
```

### 4.5 For Statement

```typescript
for (let i: number = 0; i < 10; i = i + 1) {
    // body
}
```

### 4.6 Return Statement

```typescript
return value;
return;
```

### 4.7 Break/Continue

```typescript
break;
continue;
```

---

## 5. Functions

### 5.1 Function Declaration

```typescript
function add(a: number, b: number): number {
    return a + b;
}
```

### 5.2 Function Expression

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};
```

### 5.3 Closures

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

### 6.1 Class Declaration

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

### 6.2 Inheritance

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

### 7.1 I/O

| Function | Signature | Description |
|----------|-----------|-------------|
| `println` | `(value: any) => void` | Print with newline |
| `print` | `(value: any) => void` | Print without newline |

### 7.2 Array Operations

| Function | Signature | Description |
|----------|-----------|-------------|
| `len` | `(arr: T[]) => number` | Array/string length |
| `push` | `(arr: T[], val: T) => number` | Append to array |
| `pop` | `(arr: T[]) => T` | Remove last element |

### 7.3 Type Conversion

| Function | Signature | Description |
|----------|-----------|-------------|
| `typeof` | `(val: any) => string` | Get type name |
| `tostring` | `(val: any) => string` | Convert to string |
| `tonumber` | `(val: any) => number` | Convert to number |

### 7.4 Math Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `sqrt` | `(n: number) => number` | Square root |
| `floor` | `(n: number) => number` | Floor |
| `ceil` | `(n: number) => number` | Ceiling |
| `abs` | `(n: number) => number` | Absolute value |

---

## 8. Bytecode Reference

### Opcodes

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
