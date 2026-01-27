# goTS Language Specification

## 1. Introduction

goTS (GoTypeScript) is a TypeScript-like language that compiles to Go. It provides static typing, classes, functions, and modern language features while leveraging the Go ecosystem.

## 2. Lexical Structure

### 2.1 Comments

```typescript
// Single-line comment

/* Multi-line
   comment */
```

### 2.2 Identifiers

Identifiers start with a letter or underscore, followed by letters, digits, or underscores.

```typescript
validName
_private
name123
```

### 2.3 Keywords

```
break      case       class      const      continue
default    else       enum       export     extends
false      for        function   if         import
let        new        null       return     super
switch     this       true       type       typeof
void       while
```

### 2.4 Literals

```typescript
// Number literals
42          // int
3.14        // float
0.5         // float

// String literals
"hello"
'world'
`template literal with ${expr}`

// Boolean literals
true
false

// Null literal
null
```

## 3. Types

### 3.1 Primitive Types

| Type | Description | Go Equivalent |
|------|-------------|---------------|
| `number` | Default numeric type (like TypeScript) | `float64` |
| `int` | Integer number | `int` |
| `float` | Floating-point number | `float64` |
| `string` | String of characters | `string` |
| `boolean` | Boolean value | `bool` |
| `void` | No value (function returns) | (no return) |
| `null` | Null value | `nil` |

**Note:** Numeric literals default to `number` type. Use `int` or `float` for explicit Go type mapping.

### 3.2 Type Annotations

```typescript
let x: int = 42
let name: string = "Alice"
let pi: float = 3.14159
let flag: boolean = true
```

### 3.3 Array Types

```typescript
let numbers: int[] = [1, 2, 3]
let names: string[] = ["Alice", "Bob"]
let matrix: int[][] = [[1, 2], [3, 4]]
```

### 3.4 Object Types

```typescript
type Point = {
    x: int
    y: int
}

type Person = {
    name: string
    age: int
    address: {
        street: string
        city: string
    }
}

let p: Point = {x: 10, y: 20}
```

### 3.5 Function Types

```typescript
// Function type annotation
let add: (a: int, b: int) => int

// Function with function parameter
function apply(f: (x: int) => int, value: int): int {
    return f(value)
}

// Dynamic function type
let dynamicFn: Function = function(x: int): int {
    return x * 2
}
```

### 3.6 Union Types

Union types allow a value to be one of several types using the `|` operator.

```typescript
// Basic union
type StringOrNumber = string | int
let value: StringOrNumber = "hello"
value = 42  // Also valid

// Multiple types in union
type Value = string | int | boolean
let x: Value = "text"
x = 100
x = true

// Nullable types (special case of union)
type NullableString = string | null
let name: string | null = null
name = "Alice"

// Union in function parameters
function print(value: string | int): void {
    println(value)
}
```

**Type Checking Rules:**
- A value of type `T` is assignable to `T1 | T2 | ... | Tn` if `T` is assignable to any `Ti`
- A value of union type can be used where any member type is expected (with runtime type checking)

**Code Generation:**
- Union types compile to `interface{}` in Go
- Type assertions are used at runtime when specific types are needed

### 3.7 Intersection Types

Intersection types combine multiple types using the `&` operator. For object types, this merges their properties.

```typescript
type HasName = { name: string }
type HasAge = { age: int }

// Intersection merges properties
type Person = HasName & HasAge
// Equivalent to: { name: string, age: int }

let p: Person = { name: "Alice", age: 30 }

// Multiple intersections
type A = { x: int }
type B = { y: int }
type C = { z: int }
type ABC = A & B & C
// Equivalent to: { x: int, y: int, z: int }

// Inline intersection
let point: { x: int } & { y: int } = { x: 10, y: 20 }
```

**Type Checking Rules:**
- For object types, intersection creates a new object type with all properties from all intersected types
- A value is assignable to `T1 & T2` only if it's assignable to both `T1` and `T2`
- Intersection has higher precedence than union in type expressions

**Code Generation:**
- Object intersections compile to Go structs with merged fields
- Non-object intersections compile to `interface{}`

### 3.8 Literal Types

Literal types allow you to specify exact values that a type can have.

```typescript
// String literals
type Direction = "north" | "south" | "east" | "west"
let dir: Direction = "north"
// dir = "up"  // Error: "up" is not assignable to Direction

type Status = "active" | "inactive"
let status: Status = "active"

// Number literals
type Zero = 0
type One = 1
type Port = 80 | 443 | 8080
let port: Port = 443

// Boolean literals
type AlwaysTrue = true
type AlwaysFalse = false
let t: AlwaysTrue = true

// Mixed literal union
type StatusCode = "success" | "error" | 0 | 1
let code: StatusCode = "success"
code = 0

// Literal types in function parameters
function setDirection(dir: "up" | "down" | "left" | "right"): void {
    println(dir)
}
setDirection("up")
```

**Type Checking Rules:**
- A literal value is assignable to its corresponding literal type
- A literal type is assignable to its base primitive type (`"hello"` → `string`, `42` → `int`)
- Literal types enable exhaustiveness checking and precise type constraints

**Code Generation:**
- Literal types compile to their base primitive types in Go
- `"hello"` → `string`, `42` → `int`, `3.14` → `float64`, `true` → `bool`

### 3.9 Tuple Types

Tuple types represent fixed-length arrays where each position has a specific type.

```typescript
// Basic tuple types
type Pair = [string, int]
let pair: Pair = ["hello", 42]

type Triple = [string, int, boolean]
let triple: Triple = ["test", 100, true]

// Inline tuple types
let point: [int, int] = [10, 20]
let entry: [string, int] = ["count", 5]

// Nested tuples
type Nested = [[int, int], string]
let nested: Nested = [[1, 2], "label"]

// Rest elements (variable length tail)
type StringAndNumbers = [string, ...int[]]
// First element must be string, followed by any number of ints
```

**Type Checking Rules:**
- Tuple types have a fixed number of element positions
- Each position has its own type that values must match
- Rest elements (`...T[]`) allow variable-length tails of a single type
- An array literal is assignable to a tuple if element types are compatible

**Code Generation:**
- Tuple types compile to Go structs with numbered fields
- `[string, int]` → `struct{ T0 string; T1 int }`
- Field access uses `.T0`, `.T1`, etc.

### 3.10 Type Aliases

```typescript
type Name = string
type Age = int
type Point = {x: int, y: int}
type Callback = (result: int) => void

// Using union and intersection with aliases
type ID = string | int
type Coordinates = { x: int } & { y: int }
type Mixed = { name: string } & { id: ID }
```

### 3.11 Type Inference

goTS infers types when not explicitly specified:

```typescript
let x = 42          // inferred as int
let pi = 3.14       // inferred as float
let name = "Alice"  // inferred as string
let flag = true     // inferred as boolean

// Arrays
let nums = [1, 2, 3]           // inferred as int[]
let names = ["a", "b"]         // inferred as string[]

// Functions
function add(a: int, b: int) {
    return a + b  // return type inferred as int
}
```

### 3.12 Type Mapping to Go

| goTS Type | Go Type |
|-----------|---------|
| `int` | `int` |
| `float` | `float64` |
| `string` | `string` |
| `boolean` | `bool` |
| `void` | (no return) |
| `null` | `nil` / `interface{}` |
| `Function` | `interface{}` |
| `T[]` | `[]T` |
| `T \| null` | `*T` (pointer) |
| `T1 \| T2` | `interface{}` |
| `T1 & T2` | merged struct or `interface{}` |
| `[T1, T2]` | `struct{ T0 T1; T1 T2 }` |
| `"literal"` | `string` |
| `42` | `int` |
| `class C` | `*C` (struct pointer) |

## 4. Variables and Constants

### 4.1 Variable Declarations

```typescript
// With type annotation
let x: int = 42
let name: string = "Alice"

// Type inference
let y = 100        // int
let z = 3.14       // float

// Without initialization (requires type)
let count: int
count = 0
```

### 4.2 Constants

```typescript
const PI: float = 3.14159
const MAX_SIZE: int = 100
const GREETING: string = "Hello"
```

Constants must be initialized at declaration and cannot be reassigned.

## 5. Operators

### 5.1 Arithmetic Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `+` | Addition | `a + b` |
| `-` | Subtraction | `a - b` |
| `*` | Multiplication | `a * b` |
| `/` | Division | `a / b` |
| `%` | Modulo | `a % b` |

**Type Rules:**
- `int + int = int`
- `float + float = float`
- `int + float = float`
- Division `/` always returns `float`
- Modulo `%` requires `int` operands

### 5.2 Comparison Operators

| Operator | Description |
|----------|-------------|
| `==` | Equal |
| `!=` | Not equal |
| `<` | Less than |
| `<=` | Less than or equal |
| `>` | Greater than |
| `>=` | Greater than or equal |

### 5.3 Logical Operators

| Operator | Description |
|----------|-------------|
| `&&` | Logical AND |
| `\|\|` | Logical OR |
| `!` | Logical NOT |

### 5.4 Assignment Operators

| Operator | Description |
|----------|-------------|
| `=` | Assignment |
| `+=` | Add and assign |
| `-=` | Subtract and assign |
| `*=` | Multiply and assign |
| `/=` | Divide and assign |

## 6. Control Flow

### 6.1 If Statements

```typescript
if (x > 0) {
    println("positive")
} else if (x < 0) {
    println("negative")
} else {
    println("zero")
}
```

### 6.2 While Loops

```typescript
let i = 0
while (i < 10) {
    println(i)
    i = i + 1
}
```

### 6.3 For Loops

```typescript
// Traditional for loop
for (let i = 0; i < 10; i = i + 1) {
    println(i)
}

// For-of loop (arrays)
let arr = [1, 2, 3]
for (let x of arr) {
    println(x)
}
```

### 6.4 Switch Statements

```typescript
let x: int = 2
switch (x) {
    case 1:
        println("one")
        break
    case 2:
        println("two")
        break
    default:
        println("other")
}

// Switch with string
let s: string = "hello"
switch (s) {
    case "hello":
        println("greeting")
        break
    case "bye":
        println("farewell")
        break
    default:
        println("unknown")
}
```

**Note:** Unlike JavaScript, goTS follows Go semantics where each case block is independent. Fallthrough is not automatic. Always use `break` or the case will fall through to the next.

### 6.5 Break and Continue

```typescript
while (true) {
    if (condition) {
        break
    }
    if (skip) {
        continue
    }
}
```

## 7. Functions

### 7.1 Function Declarations

```typescript
function add(a: int, b: int): int {
    return a + b
}

function greet(name: string): void {
    println("Hello, " + name)
}

// No return type (inferred as void)
function log(msg: string) {
    println(msg)
}
```

### 7.2 Function Expressions

```typescript
let add = function(a: int, b: int): int {
    return a + b
}

let greet = function(name: string): void {
    println(name)
}
```

### 7.3 Arrow Functions

```typescript
let add = (a: int, b: int): int => {
    return a + b
}

// Single expression (implicit return)
let square = (x: int): int => x * x
```

### 7.4 Higher-Order Functions

```typescript
// Function returning function
function makeAdder(x: int): Function {
    return function(y: int): int {
        return x + y
    }
}

let add5 = makeAdder(5)
println(add5(10))  // 15

// Function taking function
function apply(f: (x: int) => int, value: int): int {
    return f(value)
}

println(apply((x: int): int => x * 2, 5))  // 10
```

### 7.5 Closures

```typescript
function counter(): Function {
    let count = 0
    return function(): int {
        count = count + 1
        return count
    }
}

let c = counter()
println(c())  // 1
println(c())  // 2
```

## 8. Classes

### 8.1 Class Declarations

```typescript
class Point {
    x: int
    y: int

    constructor(x: int, y: int) {
        this.x = x
        this.y = y
    }

    distance(): float {
        return sqrt(this.x * this.x + this.y * this.y)
    }
}

let p = new Point(3, 4)
println(p.distance())  // 5.0
```

### 8.2 Inheritance

```typescript
class Animal {
    name: string

    constructor(name: string) {
        this.name = name
    }

    speak(): void {
        println(this.name)
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

let dog = new Dog("Rex", "Labrador")
dog.speak()  // "Rex barks"
```

### 8.3 Method Overriding

Subclasses can override parent class methods:

```typescript
class Base {
    method(): void {
        println("Base")
    }
}

class Derived extends Base {
    method(): void {
        println("Derived")
    }
}
```

## 9. Enums

### 9.1 Numeric Enums

```typescript
enum Color {
    Red,      // 0
    Green,    // 1
    Blue      // 2
}

let c: Color = Color.Red
println(c)  // 0

enum Status {
    Pending = 1,
    Active = 2,
    Closed = 3
}

let s: Status = Status.Active
println(s)  // 2
```

### 9.2 Enum Usage

```typescript
enum Direction {
    North,
    South,
    East,
    West
}

function move(dir: Direction): void {
    if (dir == Direction.North) {
        println("Moving north")
    }
}

move(Direction.North)
```

## 10. Template Literals

Template literals allow string interpolation using backticks:

```typescript
let name = "Alice"
let age = 30
let msg = `Hello, ${name}! You are ${age} years old.`
println(msg)  // "Hello, Alice! You are 30 years old."

// Expressions in interpolation
let x = 10
let y = 20
println(`Sum: ${x + y}`)  // "Sum: 30"

// Nested interpolation
let greeting = `Welcome, ${`Mr. ${name}`}!`
```

## 11. Destructuring

### 11.1 Array Destructuring

```typescript
let arr = [1, 2, 3]
let [a, b, c] = arr
println(a)  // 1
println(b)  // 2

// Partial destructuring
let [first, second] = [10, 20, 30]
println(first)   // 10
println(second)  // 20
```

### 11.2 Object Destructuring

```typescript
let obj = {x: 10, y: 20}
let {x, y} = obj
println(x)  // 10
println(y)  // 20

// Nested objects
let person = {
    name: "Alice",
    address: {
        city: "NYC"
    }
}
let {name, address} = person
println(address.city)  // "NYC"
```

## 12. Spread Operator

### 12.1 Array Spread

```typescript
let arr1 = [1, 2, 3]
let arr2 = [...arr1, 4, 5]
println(arr2)  // [1, 2, 3, 4, 5]

let combined = [...arr1, ...arr2]
```

### 12.2 Function Call Spread

```typescript
function sum(a: int, b: int, c: int): int {
    return a + b + c
}

let nums = [1, 2, 3]
println(sum(...nums))  // 6
```

## 13. Modules

### 13.1 Exports

```typescript
// Named exports
export function add(a: int, b: int): int {
    return a + b
}

export class Point {
    x: int
    y: int
}

export type Coord = {x: int, y: int}

// Default exports
export default class Calculator {
    // ...
}

export default function multiply(a: int, b: int): int {
    return a * b
}
```

### 13.2 Imports

```typescript
// Import specific items
import { add, Point } from "./math"

// Import with renaming
import { add as sum } from "./math"

// Type-only imports
import type { Coord } from "./types"

// Default imports
import Calculator from "./calculator"
import multiply from "./math"

// Namespace imports
import * as utils from "./utils"
// Use as: utils.add(1, 2)
```

### 13.3 Re-exports

```typescript
// Re-export specific items
export { foo, bar } from "./module"

// Re-export all
export * from "./module"
```

## 14. Type System Features

### 14.1 Type Checking

goTS performs static type checking at compile time:

```typescript
let x: int = 42
x = "hello"  // Error: string not assignable to int

function f(n: int): void {}
f("text")  // Error: argument type mismatch
```

### 14.2 Type Assertions

```typescript
let value: any = "hello"
let length: int = (value as string).length
```

### 14.3 Null Safety

```typescript
let name: string | null = null

// Must check before use
if (name != null) {
    println(name.length)
}
```

## 15. Operator Precedence

From highest to lowest:

1. Member access (`.`), function call `()`
2. Unary operators (`!`, `-`, `+`)
3. Multiplicative (`*`, `/`, `%`)
4. Additive (`+`, `-`)
5. Relational (`<`, `<=`, `>`, `>=`)
6. Equality (`==`, `!=`)
7. Logical AND (`&&`)
8. Logical OR (`||`)
9. Assignment (`=`, `+=`, `-=`, etc.)

## 16. Reserved Words and Naming

### 16.1 Go Reserved Words

When compiling to Go, these identifiers get a `_` suffix:
- Go keywords: `chan`, `defer`, `fallthrough`, `go`, `interface`, `map`, `package`, `range`, `select`, `struct`, `var`
- Go built-ins: `append`, `cap`, `close`, `complex`, `copy`, `delete`, `imag`, `make`, `panic`, `real`, `recover`

### 16.2 Naming Conventions

- Exported names (public) are capitalized in generated Go code
- Constructor functions become `NewClassName`
- Method receivers use `this` pointer

## 17. Built-in Objects

### 17.1 Math Object

The `Math` object provides mathematical constants and functions.

**Constants:**
- `Math.PI` - The ratio of a circle's circumference to its diameter (~3.14159)
- `Math.E` - Euler's number, base of natural logarithms (~2.71828)

**Rounding:**
- `Math.round(x)` - Round to nearest integer (rounds half away from zero)
- `Math.floor(x)` - Round down to integer
- `Math.ceil(x)` - Round up to integer
- `Math.trunc(x)` - Truncate to integer

**Power and Roots:**
- `Math.pow(x, y)` - x raised to the power y
- `Math.sqrt(x)` - Square root
- `Math.cbrt(x)` - Cube root
- `Math.exp(x)` - e^x

**Logarithms:**
- `Math.log(x)` - Natural logarithm
- `Math.log10(x)` - Base-10 logarithm
- `Math.log2(x)` - Base-2 logarithm

**Absolute Value and Sign:**
- `Math.abs(x)` - Absolute value
- `Math.sign(x)` - Sign of x (-1, 0, or 1)

**Min/Max:**
- `Math.min(...values)` - Minimum of values
- `Math.max(...values)` - Maximum of values

**Trigonometry:**
- `Math.sin(x)`, `Math.cos(x)`, `Math.tan(x)` - Trig functions (radians)
- `Math.asin(x)`, `Math.acos(x)`, `Math.atan(x)` - Inverse trig
- `Math.atan2(y, x)` - Angle from x-axis to point (y, x)

**Random:**
- `Math.random()` - Random number between 0 (inclusive) and 1 (exclusive)

```typescript
// Example usage
let radius: number = 5
let area: number = Math.PI * Math.pow(radius, 2)
println(area)  // ~78.54

let angle: number = Math.PI / 4  // 45 degrees
println(Math.sin(angle))  // ~0.707

let randomInt: int = toint(Math.random() * 100)  // 0-99
```

**Note:** `Math.round()` uses Go's rounding semantics (round half away from zero), which differs from JavaScript's (round half toward positive infinity) for negative half values.

### 17.2 Number Object

The `Number` object provides constants and methods for working with numbers.

**Constants:**
- `Number.MAX_SAFE_INTEGER` - Maximum safe integer (9007199254740991)
- `Number.MIN_SAFE_INTEGER` - Minimum safe integer (-9007199254740991)
- `Number.MAX_VALUE` - Largest positive number
- `Number.MIN_VALUE` - Smallest positive number (closest to zero)
- `Number.POSITIVE_INFINITY` - Positive infinity
- `Number.NEGATIVE_INFINITY` - Negative infinity
- `Number.NaN` - Not-a-Number value

**Static Methods:**
- `Number.isFinite(x)` - Check if finite number
- `Number.isNaN(x)` - Check if NaN
- `Number.isInteger(x)` - Check if integer
- `Number.isSafeInteger(x)` - Check if safe integer
- `Number.parseFloat(s)` - Parse float from string
- `Number.parseInt(s, radix?)` - Parse int from string with optional radix

**Global Functions:**
- `isNaN(x)` - Global NaN check
- `isFinite(x)` - Global finite check
- `parseFloat(s)` - Global float parser

```typescript
// Example usage
println(Number.isInteger(42))      // true
println(Number.isInteger(3.14))    // false
println(Number.parseFloat("3.14")) // 3.14
println(Number.parseInt("ff", 16)) // 255
println(isFinite(42))              // true
```

### 17.3 JSON Object

The `JSON` object provides methods for parsing and serializing JSON data.

**Methods:**
- `JSON.stringify(value)` - Convert value to JSON string
- `JSON.parse(text)` - Parse JSON string to value

```typescript
// Stringify
let json: string = JSON.stringify([1, 2, 3])  // "[1,2,3]"
println(JSON.stringify(42))                    // "42"
println(JSON.stringify(true))                  // "true"

// Parse (provide type annotation)
let num: number = JSON.parse("42")
let arr: number[] = JSON.parse("[1,2,3]")
```

**Note:** `JSON.parse` returns `any` type. Provide type annotations for proper type checking.

## 18. Unsupported Features

The following JavaScript/TypeScript features are **not supported** in goTS:

| Feature | Reason |
|---------|--------|
| `do-while` loops | Not implemented |
| `async/await` | No async support |
| Generators | No generator functions |
| Promises | No Promise API |
| Symbols | No Symbol type |
| Proxy/Reflect | No metaprogramming |
| WeakMap/WeakSet | No weak references |
| `eval()` | No runtime eval |
| `with` statement | Not supported |
| Destructuring | Not yet implemented |
| Spread operator | Not yet implemented |
| Optional chaining (`?.`) | Not yet implemented |
| Nullish coalescing (`??`) | Not yet implemented |
| BigInt | Not supported |

### 18.1 Differences from JavaScript

1. **Switch fallthrough**: Unlike JavaScript, switch cases don't automatically fall through. Use explicit `break` statements.

2. **typeof for numbers**: `typeof` returns `"number"` for both `int` and `float` types (consistent with TypeScript semantics).

3. **Strict typing**: goTS enforces static types at compile time. Dynamic typing patterns common in JavaScript may not work.

4. **No hoisting**: Variables must be declared before use. Function hoisting is not supported.
