# GoTS Language Specification

**Version 1.0**
**January 2026**

---

## Table of Contents

1. [Introduction](#1-introduction)
   - 1.1 [Design Goals](#11-design-goals)
   - 1.2 [Relationship to TypeScript](#12-relationship-to-typescript)
   - 1.3 [Compilation Model](#13-compilation-model)
2. [Lexical Structure](#2-lexical-structure)
   - 2.1 [Programs](#21-programs)
   - 2.2 [Comments](#22-comments)
   - 2.3 [Identifiers](#23-identifiers)
   - 2.4 [Keywords](#24-keywords)
   - 2.5 [Literals](#25-literals)
   - 2.6 [Operators and Punctuation](#26-operators-and-punctuation)
3. [Types](#3-types)
   - 3.1 [Type System Overview](#31-type-system-overview)
   - 3.2 [Primitive Types](#32-primitive-types)
   - 3.3 [Array Types](#33-array-types)
   - 3.4 [Object Types](#34-object-types)
   - 3.5 [Function Types](#35-function-types)
   - 3.6 [Class Types](#36-class-types)
   - 3.7 [Union Types](#37-union-types)
   - 3.8 [Type Aliases](#38-type-aliases)
   - 3.9 [Type Annotations](#39-type-annotations)
   - 3.10 [Type Mapping to Go](#310-type-mapping-to-go)
4. [Variables and Constants](#4-variables-and-constants)
   - 4.1 [Variable Declarations](#41-variable-declarations)
   - 4.2 [Constant Declarations](#42-constant-declarations)
   - 4.3 [Scope Rules](#43-scope-rules)
5. [Expressions](#5-expressions)
   - 5.1 [Primary Expressions](#51-primary-expressions)
   - 5.2 [Arithmetic Expressions](#52-arithmetic-expressions)
   - 5.3 [Comparison Expressions](#53-comparison-expressions)
   - 5.4 [Logical Expressions](#54-logical-expressions)
   - 5.5 [Assignment Expressions](#55-assignment-expressions)
   - 5.6 [Property Access](#56-property-access)
   - 5.7 [Index Access](#57-index-access)
   - 5.8 [Function Calls](#58-function-calls)
   - 5.9 [Object Creation](#59-object-creation)
   - 5.10 [Optional Chaining](#510-optional-chaining)
   - 5.11 [Nullish Coalescing](#511-nullish-coalescing)
6. [Statements](#6-statements)
   - 6.1 [Expression Statements](#61-expression-statements)
   - 6.2 [Block Statements](#62-block-statements)
   - 6.3 [If Statements](#63-if-statements)
   - 6.4 [While Statements](#64-while-statements)
   - 6.5 [For Statements](#65-for-statements)
   - 6.6 [For-Of Statements](#66-for-of-statements)
   - 6.7 [Break and Continue](#67-break-and-continue)
   - 6.8 [Return Statements](#68-return-statements)
   - 6.9 [Try-Catch Statements](#69-try-catch-statements)
   - 6.10 [Throw Statements](#610-throw-statements)
7. [Functions](#7-functions)
   - 7.1 [Function Declarations](#71-function-declarations)
   - 7.2 [Function Expressions](#72-function-expressions)
   - 7.3 [Parameters and Return Types](#73-parameters-and-return-types)
   - 7.4 [Closures](#74-closures)
   - 7.5 [Higher-Order Functions](#75-higher-order-functions)
8. [Classes](#8-classes)
   - 8.1 [Class Declarations](#81-class-declarations)
   - 8.2 [Constructors](#82-constructors)
   - 8.3 [Properties](#83-properties)
   - 8.4 [Methods](#84-methods)
   - 8.5 [Inheritance](#85-inheritance)
   - 8.6 [The this Keyword](#86-the-this-keyword)
   - 8.7 [The super Keyword](#87-the-super-keyword)
9. [Built-in Functions](#9-built-in-functions)
   - 9.1 [I/O Functions](#91-io-functions)
   - 9.2 [Array Functions](#92-array-functions)
   - 9.3 [Type Conversion](#93-type-conversion)
   - 9.4 [Math Functions](#94-math-functions)
10. [Differences from TypeScript](#10-differences-from-typescript)
    - 10.1 [Supported Features](#101-supported-features)
    - 10.2 [Unsupported Features](#102-unsupported-features)
11. [Appendix](#11-appendix)
    - 11.1 [Grammar Reference](#111-grammar-reference)
    - 11.2 [Operator Precedence](#112-operator-precedence)
    - 11.3 [Reserved Words](#113-reserved-words)

---

## 1. Introduction

GoTS (Go-TypeScript) is a statically-typed programming language with TypeScript-like syntax that compiles to Go. It provides a familiar development experience for TypeScript developers while leveraging the Go ecosystem and toolchain.

### 1.1 Design Goals

The primary design goals of GoTS are:

1. **Static Type Safety**: All types are checked at compile time
2. **TypeScript Compatibility**: Use TypeScript-like syntax for familiarity
3. **Go Interoperability**: Compile to idiomatic Go code
4. **Simplicity**: Be a minimal, focused subset of TypeScript
5. **Performance**: Generate efficient Go code that compiles to native binaries

### 1.2 Relationship to TypeScript

GoTS is designed to be a **strict subset of TypeScript** with the following characteristics:

- Valid GoTS code should be syntactically valid TypeScript
- GoTS enforces stricter rules than TypeScript (e.g., all variables require explicit type annotations)
- Not all TypeScript features are supported (see [Section 10](#10-differences-from-typescript))
- GoTS uses `int` and `float` as distinct numeric types, whereas TypeScript uses `number`

### 1.3 Compilation Model

GoTS follows this compilation pipeline:

```
Source (.gts) → Lexer → Parser → Type Checker → Go Code Generator → go build → Native Binary
```

**Phases:**

1. **Lexical Analysis**: Source code is tokenized
2. **Parsing**: Tokens are parsed into an Abstract Syntax Tree (AST)
3. **Type Checking**: AST is transformed to a typed AST with full type annotations
4. **Code Generation**: Typed AST is translated to Go source code
5. **Go Compilation**: Generated Go code is compiled by the Go toolchain

---

## 2. Lexical Structure

### 2.1 Programs

A GoTS program consists of a sequence of statements and declarations in a source file with the `.gts` extension.

### 2.2 Comments

GoTS supports single-line comments:

```typescript
// This is a comment
let x: int = 42  // End-of-line comment
```

**Note**: Multi-line comments (`/* ... */`) are not currently supported.

### 2.3 Identifiers

An identifier is a sequence of characters used to name variables, functions, classes, and types.

**Syntax:**
- Must start with a letter (`a-z`, `A-Z`) or underscore (`_`)
- Subsequent characters can be letters, digits (`0-9`), or underscores
- Case-sensitive

**Examples:**
```typescript
x
myVariable
_private
counter2
MyClass
```

### 2.4 Keywords

The following identifiers are reserved as keywords and cannot be used as identifiers:

```
break       class       const       constructor continue
else        extends     false       for         function
if          let         new         null        of
return      super       this        throw       true
try         catch       type        void        while
```

### 2.5 Literals

#### 2.5.1 Integer Literals

Integer literals are sequences of digits representing integer values.

```typescript
0
42
1000
```

**Type**: Integer literals have type `int`.

#### 2.5.2 Float Literals

Float literals contain a decimal point.

```typescript
3.14
0.5
2.0
```

**Type**: Float literals have type `float`.

#### 2.5.3 String Literals

String literals are sequences of characters enclosed in double quotes.

```typescript
"hello"
"world"
"Hello, World!"
""  // empty string
```

**Escape sequences**: `\n` (newline), `\t` (tab), `\"` (quote), `\\` (backslash)

#### 2.5.4 Boolean Literals

```typescript
true
false
```

#### 2.5.5 Null Literal

```typescript
null
```

#### 2.5.6 Array Literals

```typescript
[1, 2, 3]
["a", "b", "c"]
[]  // empty array
```

#### 2.5.7 Object Literals

```typescript
{ x: 10, y: 20 }
{ name: "Alice", age: 30 }
{}  // empty object
```

### 2.6 Operators and Punctuation

**Operators:**
```
+  -  *  /  %  ==  !=  <  >  <=  >=  &&  ||  !  =  ?.  ??
```

**Punctuation:**
```
(  )  {  }  [  ]  ;  :  ,  .  =>  |
```

---

## 3. Types

### 3.1 Type System Overview

GoTS has a **static type system**. Every variable, parameter, and expression has a type that is determined at compile time. Type annotations are **mandatory** for variable declarations, function parameters, and function return types.

### 3.2 Primitive Types

#### 3.2.1 The int Type

The `int` type represents integer values.

```typescript
let count: int = 42
let negative: int = -10
```

**Mapping**: Maps to Go's `int` type.

#### 3.2.2 The float Type

The `float` type represents floating-point numbers.

```typescript
let pi: float = 3.14159
let half: float = 0.5
```

**Mapping**: Maps to Go's `float64` type.

#### 3.2.3 The string Type

The `string` type represents sequences of characters.

```typescript
let name: string = "Alice"
let greeting: string = "Hello, " + name
```

**Mapping**: Maps to Go's `string` type.

#### 3.2.4 The boolean Type

The `boolean` type represents logical values.

```typescript
let flag: boolean = true
let done: boolean = false
```

**Mapping**: Maps to Go's `bool` type.

#### 3.2.5 The void Type

The `void` type represents the absence of a value. It is used as the return type of functions that do not return a value.

```typescript
function greet(): void {
    println("Hello")
}
```

**Mapping**: Functions with `void` return type are translated to Go functions with no return value.

#### 3.2.6 The null Type

The `null` type has a single value: `null`.

```typescript
let x: null = null
```

**Mapping**: In union types like `T | null`, maps to Go pointer types (`*T`) where `nil` represents null.

### 3.3 Array Types

An array type is written as `ElementType[]`.

```typescript
let numbers: int[] = [1, 2, 3]
let names: string[] = ["Alice", "Bob"]
let matrix: int[][] = [[1, 2], [3, 4]]
```

**Mapping**: Maps to Go slice types (`[]ElementType`).

**Operations**:
- Index access: `arr[0]`
- Length: `len(arr)`
- Append: `push(arr, value)`
- Pop: `pop(arr)`

### 3.4 Object Types

Object types describe objects with named properties.

**Syntax:**
```typescript
{ property1: Type1, property2: Type2, ... }
```

**Example:**
```typescript
let point: { x: int, y: int } = { x: 10, y: 20 }
let person: { name: string, age: int } = { name: "Alice", age: 30 }
```

**Mapping**: Maps to Go struct types with corresponding fields.

**Property Access:**
```typescript
println(point.x)
person.age = 31
```

### 3.5 Function Types

Function types describe the signature of functions.

**Syntax:**
```typescript
(param1: Type1, param2: Type2, ...) => ReturnType
```

**Example:**
```typescript
let add: (a: int, b: int) => int = function(a: int, b: int): int {
    return a + b
}
```

**Dynamic Function Type:**

The special `Function` type represents a function with dynamic typing (any parameters, any return type):

```typescript
let callback: Function = function(x: int): int {
    return x * 2
}
```

**Mapping**:
- Typed function types map to Go function types
- `Function` maps to Go's `interface{}` and requires runtime type assertions

### 3.6 Class Types

Class types are introduced by class declarations (see [Section 8](#8-classes)).

```typescript
class Point {
    x: int
    y: int
    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }
}

let p: Point = new Point(10, 20)
```

**Mapping**: Class instances map to Go struct pointers (`*StructName`).

### 3.7 Union Types

Union types represent values that can be one of several types.

**Syntax:**
```typescript
Type1 | Type2
```

**Nullable Types:**

The most common use is nullable types:

```typescript
let name: string | null = null
name = "Alice"
```

**Mapping**: `T | null` maps to Go pointer type `*T`, where `nil` represents null.

**Type Guards:**

Use comparison with `null` for type narrowing:

```typescript
let x: string | null = getValue()
if (x != null) {
    // x is string here
    println(x)
}
```

### 3.8 Type Aliases

Type aliases create new names for existing types.

**Syntax:**
```typescript
type AliasName = Type
```

**Examples:**
```typescript
type Point = { x: int, y: int }
type Callback = (n: int) => void
type NullableString = string | null
```

### 3.9 Type Annotations

Type annotations are **required** in the following contexts:

1. **Variable declarations:**
   ```typescript
   let x: int = 42
   ```

2. **Function parameters:**
   ```typescript
   function add(a: int, b: int): int { ... }
   ```

3. **Function return types:**
   ```typescript
   function getValue(): string { ... }
   ```

4. **Class properties:**
   ```typescript
   class Person {
       name: string
       age: int
   }
   ```

**Type inference** is not supported for variable declarations. All declarations must have explicit type annotations.

### 3.10 Type Mapping to Go

| GoTS Type | Go Type | Notes |
|-----------|---------|-------|
| `int` | `int` | Integer type |
| `float` | `float64` | Floating-point type |
| `string` | `string` | String type |
| `boolean` | `bool` | Boolean type |
| `void` | (no return) | Used for function returns |
| `null` | `nil` / `interface{}` | Context-dependent |
| `T[]` | `[]T` | Slice type |
| `{ x: T, y: U }` | `struct { X T; Y U }` | Struct type (fields capitalized) |
| `(a: T) => U` | `func(T) U` | Function type |
| `Function` | `interface{}` | Dynamic function type |
| `T \| null` | `*T` | Pointer type (nil = null) |
| `class C` | `*C` | Struct pointer |

---

## 4. Variables and Constants

### 4.1 Variable Declarations

Variables are declared using the `let` keyword.

**Syntax:**
```typescript
let identifier: Type = initializer;
```

**Example:**
```typescript
let x: int = 10
let name: string = "Alice"
let items: int[] = [1, 2, 3]
```

**Rules:**
- Type annotation is **required**
- Initializer is **required**
- Variables can be reassigned

### 4.2 Constant Declarations

Constants are declared using the `const` keyword.

**Syntax:**
```typescript
const identifier: Type = initializer;
```

**Example:**
```typescript
const PI: float = 3.14159
const MAX_SIZE: int = 100
```

**Rules:**
- Type annotation is **required**
- Initializer is **required**
- Constants **cannot** be reassigned

### 4.3 Scope Rules

GoTS uses **lexical (block) scoping**:

1. Variables declared in a block are scoped to that block
2. Inner scopes can access outer scope variables
3. Inner scope variables shadow outer scope variables with the same name
4. Function parameters are scoped to the function body

**Example:**
```typescript
let x: int = 10

function foo(): void {
    let x: int = 20  // Shadows outer x
    println(x)       // Prints 20
}

foo()
println(x)           // Prints 10
```

---

## 5. Expressions

### 5.1 Primary Expressions

Primary expressions are the building blocks of more complex expressions:

- Literals: `42`, `"hello"`, `true`, `null`
- Identifiers: `x`, `myVar`
- Array literals: `[1, 2, 3]`
- Object literals: `{ x: 10, y: 20 }`
- Parenthesized expressions: `(x + y)`

### 5.2 Arithmetic Expressions

**Operators:** `+`, `-`, `*`, `/`, `%`

**Type Rules:**

| Left | Operator | Right | Result |
|------|----------|-------|--------|
| `int` | `+` `-` `*` | `int` | `int` |
| `int` | `+` `-` `*` | `float` | `float` |
| `float` | `+` `-` `*` | `int` | `float` |
| `float` | `+` `-` `*` | `float` | `float` |
| `int` | `/` | `int` | `float` |
| `int` | `/` | `float` | `float` |
| `float` | `/` | any | `float` |
| `int` | `%` | `int` | `int` |
| `string` | `+` | `string` | `string` |

**Special Rules:**
- Division (`/`) **always** returns `float`, even for integer operands
- Modulo (`%`) requires both operands to be `int`
- The `+` operator concatenates strings when both operands are strings
- Unary `-` negates numeric values

**Examples:**
```typescript
let a: int = 10 + 5      // 15 (int)
let b: float = 10 / 3    // 3.333... (float)
let c: int = 10 % 3      // 1 (int)
let d: float = 5 + 2.5   // 7.5 (float)
let s: string = "Hello" + " " + "World"  // "Hello World"
```

### 5.3 Comparison Expressions

**Operators:** `==`, `!=`, `<`, `>`, `<=`, `>=`

**Type Rules:**
- Operands must be comparable (same type or compatible types)
- Result is always `boolean`

**Examples:**
```typescript
let eq: boolean = 5 == 5        // true
let ne: boolean = 5 != 3        // true
let lt: boolean = 3 < 5         // true
let gte: boolean = 10 >= 10     // true
let strEq: boolean = "a" == "a" // true
```

### 5.4 Logical Expressions

**Operators:** `&&` (and), `||` (or), `!` (not)

**Type Rules:**
- Operands must be `boolean`
- Result is `boolean`

**Examples:**
```typescript
let a: boolean = true && false   // false
let b: boolean = true || false   // true
let c: boolean = !true           // false
```

### 5.5 Assignment Expressions

**Syntax:**
```typescript
target = value
```

**Targets:**
- Variable: `x = 10`
- Property: `obj.prop = value`
- Index: `arr[0] = value`

**Example:**
```typescript
let x: int = 5
x = x + 1

let arr: int[] = [1, 2, 3]
arr[0] = 10

let p: { x: int, y: int } = { x: 0, y: 0 }
p.x = 5
```

### 5.6 Property Access

**Syntax:**
```typescript
object.property
```

**Example:**
```typescript
let p: { x: int, y: int } = { x: 10, y: 20 }
let xVal: int = p.x
p.y = 30
```

### 5.7 Index Access

**Syntax:**
```typescript
array[index]
```

**Rules:**
- Index must be `int`
- Result type is the array element type

**Example:**
```typescript
let arr: int[] = [10, 20, 30]
let first: int = arr[0]
arr[1] = 25
```

### 5.8 Function Calls

**Syntax:**
```typescript
functionName(arg1, arg2, ...)
object.method(arg1, arg2, ...)
```

**Example:**
```typescript
function add(a: int, b: int): int {
    return a + b
}

let sum: int = add(5, 3)
println("Result: " + tostring(sum))
```

### 5.9 Object Creation

**Syntax:**
```typescript
new ClassName(arguments)
```

**Example:**
```typescript
class Point {
    x: int
    y: int
    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }
}

let p: Point = new Point(10, 20)
```

### 5.10 Optional Chaining

Optional chaining allows safe access to properties that might be null.

**Syntax:**
```typescript
object?.property
```

**Example:**
```typescript
let person: Person | null = getPerson()
if (person?.name != null) {
    println(person.name)
}
```

**Rules:**
- Can only be used with nullable types (`T | null`)
- If the object is `null`, the expression evaluates to `null`

### 5.11 Nullish Coalescing

The nullish coalescing operator provides a default value when a value is null.

**Syntax:**
```typescript
value ?? defaultValue
```

**Example:**
```typescript
let name: string | null = getName()
let displayName: string = name ?? "Anonymous"
```

**Rules:**
- Returns the left operand if it's not null
- Returns the right operand if the left is null

---

## 6. Statements

### 6.1 Expression Statements

Any expression can be used as a statement.

```typescript
println("Hello")
x = x + 1
add(5, 3)
```

### 6.2 Block Statements

A block statement groups multiple statements.

**Syntax:**
```typescript
{
    statement1
    statement2
    ...
}
```

**Example:**
```typescript
{
    let x: int = 10
    println(x)
}
```

### 6.3 If Statements

**Syntax:**
```typescript
if (condition) {
    // then branch
}

if (condition) {
    // then branch
} else {
    // else branch
}
```

**Rules:**
- Condition must be `boolean`
- Braces are required

**Example:**
```typescript
if (x > 10) {
    println("big")
} else {
    println("small")
}
```

### 6.4 While Statements

**Syntax:**
```typescript
while (condition) {
    // body
}
```

**Rules:**
- Condition must be `boolean`

**Example:**
```typescript
let i: int = 0
while (i < 10) {
    println(i)
    i = i + 1
}
```

### 6.5 For Statements

**Syntax:**
```typescript
for (initialization; condition; increment) {
    // body
}
```

**Example:**
```typescript
for (let i: int = 0; i < 10; i = i + 1) {
    println(i)
}
```

### 6.6 For-Of Statements

For-of loops iterate over array elements or string characters.

**Syntax:**
```typescript
for (let variable of iterable) {
    // body
}
```

**Example:**
```typescript
let numbers: int[] = [1, 2, 3, 4, 5]
for (let n of numbers) {
    println(n)
}

let text: string = "hello"
for (let ch of text) {
    println(ch)  // ch is string (single character)
}
```

### 6.7 Break and Continue

**break**: Exit the innermost loop

```typescript
while (true) {
    if (done) {
        break
    }
}
```

**continue**: Skip to the next iteration

```typescript
for (let i: int = 0; i < 10; i = i + 1) {
    if (i % 2 == 0) {
        continue
    }
    println(i)  // Only prints odd numbers
}
```

### 6.8 Return Statements

**Syntax:**
```typescript
return expression
return  // for void functions
```

**Example:**
```typescript
function add(a: int, b: int): int {
    return a + b
}

function greet(): void {
    println("Hello")
    return
}
```

### 6.9 Try-Catch Statements

Try-catch statements handle runtime errors.

**Syntax:**
```typescript
try {
    // code that might throw
} catch (errorVariable) {
    // error handling
}
```

**Example:**
```typescript
try {
    let result: int = riskyOperation()
    println(result)
} catch (e) {
    println("An error occurred")
}
```

### 6.10 Throw Statements

**Syntax:**
```typescript
throw expression
```

**Example:**
```typescript
function divide(a: int, b: int): float {
    if (b == 0) {
        throw "Division by zero"
    }
    return a / b
}
```

---

## 7. Functions

### 7.1 Function Declarations

**Syntax:**
```typescript
function functionName(param1: Type1, param2: Type2, ...): ReturnType {
    // body
}
```

**Example:**
```typescript
function add(a: int, b: int): int {
    return a + b
}

function greet(name: string): void {
    println("Hello, " + name)
}
```

### 7.2 Function Expressions

Functions can be assigned to variables.

**Syntax:**
```typescript
let variable: FunctionType = function(params): ReturnType {
    // body
}
```

**Example:**
```typescript
let add: (a: int, b: int) => int = function(a: int, b: int): int {
    return a + b
}

let greet: (name: string) => void = function(name: string): void {
    println("Hello, " + name)
}
```

### 7.3 Parameters and Return Types

**Rules:**
- All parameters must have type annotations
- Return type annotation is required
- Use `void` for functions that don't return a value

**Example:**
```typescript
function multiply(x: int, y: int): int {
    return x * y
}

function log(message: string): void {
    println(message)
}
```

### 7.4 Closures

Functions can capture variables from their enclosing scope.

**Example:**
```typescript
function makeCounter(): Function {
    let count: int = 0
    return function(): int {
        count = count + 1
        return count
    }
}

let counter: Function = makeCounter()
println(counter())  // 1
println(counter())  // 2
println(counter())  // 3
```

**Implementation:**
- Captured variables are lifted to heap-allocated structs in Go
- Each closure instance maintains its own copy of captured variables

### 7.5 Higher-Order Functions

Functions can accept and return other functions.

**Example:**
```typescript
function apply(f: Function, x: int): int {
    return f(x)
}

function double(x: int): int {
    return x * 2
}

let result: int = apply(double, 5)  // 10
```

**Currying Example:**
```typescript
function curry_add(a: int): Function {
    return function(b: int): int {
        return a + b
    }
}

let add5: Function = curry_add(5)
println(add5(3))   // 8
println(add5(10))  // 15
```

---

## 8. Classes

### 8.1 Class Declarations

**Syntax:**
```typescript
class ClassName {
    property1: Type1
    property2: Type2

    constructor(params) {
        // initialization
    }

    method1(params): ReturnType {
        // body
    }
}
```

**Example:**
```typescript
class Point {
    x: int
    y: int

    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }

    distance(): float {
        return sqrt(tofloat(this.x * this.x + this.y * this.y))
    }
}
```

### 8.2 Constructors

Every class must have exactly one constructor.

**Syntax:**
```typescript
constructor(param1: Type1, param2: Type2, ...) {
    // initialization
}
```

**Example:**
```typescript
class Person {
    name: string
    age: int

    constructor(name: string, age: int) {
        this.name = name
        this.age = age
    }
}

let p: Person = new Person("Alice", 30)
```

### 8.3 Properties

Class properties must have type annotations.

**Example:**
```typescript
class Rectangle {
    width: int
    height: int

    constructor(w: int, h: int) {
        this.width = w
        this.height = h
    }
}
```

**Access:**
```typescript
let r: Rectangle = new Rectangle(10, 20)
println(r.width)
r.height = 25
```

### 8.4 Methods

Methods are functions defined within a class.

**Example:**
```typescript
class Counter {
    count: int

    constructor() {
        this.count = 0
    }

    increment(): void {
        this.count = this.count + 1
    }

    getCount(): int {
        return this.count
    }
}
```

### 8.5 Inheritance

GoTS supports **single inheritance** using the `extends` keyword.

**Syntax:**
```typescript
class DerivedClass extends BaseClass {
    // additional members
}
```

**Example:**
```typescript
class Animal {
    name: string

    constructor(name: string) {
        this.name = name
    }

    speak(): void {
        println(this.name + " makes a sound")
    }
}

class Dog extends Animal {
    breed: string

    constructor(name: string, breed: string) {
        super(name)
        this.breed = breed
    }

    speak(): void {
        println(this.name + " barks")
    }
}
```

**Rules:**
- Derived class inherits all properties and methods from base class
- Methods can be overridden
- Constructor must call `super()` to initialize base class

### 8.6 The this Keyword

`this` refers to the current instance within class methods.

**Example:**
```typescript
class Calculator {
    value: int

    constructor(initial: int) {
        this.value = initial
    }

    add(n: int): void {
        this.value = this.value + n
    }
}
```

### 8.7 The super Keyword

`super` is used to call the base class constructor.

**Syntax:**
```typescript
super(arguments)
```

**Example:**
```typescript
class Vehicle {
    wheels: int

    constructor(wheels: int) {
        this.wheels = wheels
    }
}

class Car extends Vehicle {
    brand: string

    constructor(brand: string) {
        super(4)  // Cars have 4 wheels
        this.brand = brand
    }
}
```

**Rules:**
- Must be called in derived class constructor
- Must be called before accessing `this`

---

## 9. Built-in Functions

### 9.1 I/O Functions

#### println

Prints a value followed by a newline.

```typescript
println(value: any): void
```

**Example:**
```typescript
println("Hello, World!")
println(42)
println(true)
```

#### print

Prints a value without a newline.

```typescript
print(value: any): void
```

**Example:**
```typescript
print("Hello")
print(" ")
println("World")
```

### 9.2 Array Functions

#### len

Returns the length of an array or string.

```typescript
len(arr: T[]): int
len(str: string): int
```

**Example:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
println(len(arr))  // 5

let text: string = "hello"
println(len(text))  // 5
```

#### push

Appends an element to an array.

```typescript
push(arr: T[], value: T): void
```

**Example:**
```typescript
let arr: int[] = [1, 2, 3]
push(arr, 4)
println(len(arr))  // 4
```

#### pop

Removes and returns the last element from an array.

```typescript
pop(arr: T[]): T
```

**Example:**
```typescript
let arr: int[] = [1, 2, 3]
let last: int = pop(arr)
println(last)      // 3
println(len(arr))  // 2
```

### 9.3 Type Conversion

#### tostring

Converts a value to a string.

```typescript
tostring(value: any): string
```

**Example:**
```typescript
let n: int = 42
let s: string = tostring(n)  // "42"
println("The answer is " + s)
```

#### toint

Converts a value to an integer.

```typescript
toint(value: any): int
```

**Example:**
```typescript
let f: float = 3.14
let i: int = toint(f)  // 3
```

#### tofloat

Converts a value to a float.

```typescript
tofloat(value: any): float
```

**Example:**
```typescript
let i: int = 42
let f: float = tofloat(i)  // 42.0
```

#### typeof

Returns the type name of a value as a string.

```typescript
typeof(value: any): string
```

**Example:**
```typescript
println(typeof(42))        // "int"
println(typeof(3.14))      // "float"
println(typeof("hello"))   // "string"
println(typeof(true))      // "boolean"
```

### 9.4 Math Functions

#### sqrt

Returns the square root of a number.

```typescript
sqrt(n: float): float
```

**Example:**
```typescript
println(sqrt(16.0))  // 4.0
println(sqrt(2.0))   // 1.414...
```

#### floor

Returns the largest integer less than or equal to a number.

```typescript
floor(n: float): int
```

**Example:**
```typescript
println(floor(3.7))   // 3
println(floor(-2.3))  // -3
```

#### ceil

Returns the smallest integer greater than or equal to a number.

```typescript
ceil(n: float): int
```

**Example:**
```typescript
println(ceil(3.2))   // 4
println(ceil(-2.7))  // -2
```

#### abs

Returns the absolute value of a number.

```typescript
abs(n: float): float
abs(n: int): int
```

**Example:**
```typescript
println(abs(-5))     // 5
println(abs(-3.14))  // 3.14
```

---

## 10. Differences from TypeScript

### 10.1 Supported Features

GoTS supports the following TypeScript features:

✅ **Type Annotations**
- Primitive types: `int`, `float`, `string`, `boolean`, `void`, `null`
- Array types: `T[]`
- Object types: `{ x: T, y: U }`
- Function types: `(a: T) => U`
- Union types (limited to nullable: `T | null`)
- Type aliases: `type Name = Type`

✅ **Variables and Constants**
- `let` declarations
- `const` declarations

✅ **Functions**
- Function declarations
- Function expressions
- Closures
- Higher-order functions

✅ **Classes**
- Class declarations
- Constructors
- Properties
- Methods
- Single inheritance with `extends`
- `this` and `super` keywords

✅ **Control Flow**
- `if`/`else` statements
- `while` loops
- `for` loops
- `for-of` loops
- `break` and `continue`

✅ **Error Handling**
- `try`/`catch`/`throw`

✅ **Operators**
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `&&`, `||`, `!`
- Optional chaining: `?.`
- Nullish coalescing: `??`

### 10.2 Unsupported Features

The following TypeScript features are **not supported** in GoTS:

❌ **Type System Features**
- Type inference (all types must be explicit)
- `any` type (use `Function` for dynamic function types)
- `unknown` type
- `never` type
- Tuple types beyond arrays
- Enum types
- Intersection types (except in specific contexts)
- Generic types
- Conditional types
- Mapped types
- Template literal types

❌ **Language Features**
- Interfaces (use type aliases with object types)
- Namespaces/modules
- Decorators
- Async/await
- Generators
- Arrow functions (use function expressions)
- Destructuring
- Spread operator
- Rest parameters
- Default parameters
- Optional parameters
- Method overloading
- Getter/setter properties
- Static class members
- Private/protected/public modifiers
- Abstract classes
- Multiple inheritance
- Mixins

❌ **Advanced Operators**
- Ternary operator `? :`
- Increment/decrement: `++`, `--`
- Compound assignment: `+=`, `-=`, etc.
- Bitwise operators
- Type assertions/casting

❌ **Advanced Syntax**
- Multi-line comments `/* */`
- Template literals (use string concatenation)
- Regular expressions
- Symbol type
- BigInt type

---

## 11. Appendix

### 11.1 Grammar Reference

```
Program = Statement*

Statement =
    | VariableDeclaration
    | ConstDeclaration
    | FunctionDeclaration
    | ClassDeclaration
    | TypeAlias
    | IfStatement
    | WhileStatement
    | ForStatement
    | ForOfStatement
    | ReturnStatement
    | BreakStatement
    | ContinueStatement
    | ThrowStatement
    | TryStatement
    | BlockStatement
    | ExpressionStatement

VariableDeclaration = "let" Identifier ":" Type "=" Expression

ConstDeclaration = "const" Identifier ":" Type "=" Expression

FunctionDeclaration = "function" Identifier "(" Parameters ")" ":" Type Block

ClassDeclaration = "class" Identifier ("extends" Identifier)? "{" ClassMember* "}"

ClassMember =
    | PropertyDeclaration
    | Constructor
    | MethodDeclaration

TypeAlias = "type" Identifier "=" Type

Expression =
    | Literal
    | Identifier
    | BinaryExpression
    | UnaryExpression
    | CallExpression
    | PropertyAccess
    | IndexAccess
    | ObjectLiteral
    | ArrayLiteral
    | NewExpression
    | AssignmentExpression

Type =
    | "int" | "float" | "string" | "boolean" | "void" | "null"
    | Type "[]"
    | "(" ParameterTypes ")" "=>" Type
    | "{" PropertyTypes "}"
    | Identifier
    | Type "|" Type
```

### 11.2 Operator Precedence

| Precedence | Operators | Associativity |
|------------|-----------|---------------|
| 1 (highest) | `.`, `[]`, `()`, `new` | Left |
| 2 | Unary `-`, `!` | Right |
| 3 | `*`, `/`, `%` | Left |
| 4 | `+`, `-` | Left |
| 5 | `<`, `>`, `<=`, `>=` | Left |
| 6 | `==`, `!=` | Left |
| 7 | `&&` | Left |
| 8 | `\|\|` | Left |
| 9 | `??` | Left |
| 10 (lowest) | `=` | Right |

### 11.3 Reserved Words

The following words are reserved and cannot be used as identifiers:

```
break       case        catch       class       const
constructor continue    debugger    default     delete
do          else        enum        export      extends
false       finally     for         function    if
import      in          instanceof  let         new
null        of          return      super       switch
this        throw       true        try         typeof
var         void        while       with        yield
```

**Note**: While some keywords (like `case`, `switch`, `var`) are reserved, they are not currently used in GoTS but may be supported in future versions.

---

## Conclusion

GoTS provides a minimal, statically-typed language with TypeScript-like syntax that compiles to efficient Go code. By focusing on a core subset of TypeScript features, GoTS offers:

- **Simplicity**: Easy to learn and reason about
- **Safety**: Full static type checking
- **Performance**: Native code generation via Go
- **Familiarity**: Syntax compatible with TypeScript

For questions, issues, or contributions, please visit the GoTS repository.

---

**GoTS Language Specification v1.0**
**© 2026**
