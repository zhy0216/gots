# goTS Built-in Functions Reference

This document describes all built-in functions available in goTS programs.

## Console Output

### println

Prints values to standard output with a newline.

**Signature:**
```typescript
function println(...args: any): void
```

**Parameters:**
- `args` - Variable number of arguments of any type

**Returns:** void

**Examples:**
```typescript
println("Hello, world!")
println(42)
println("x =", 10, "y =", 20)
```

**Go Mapping:** `fmt.Println()`

---

### print

Prints values to standard output without a newline.

**Signature:**
```typescript
function print(...args: any): void
```

**Parameters:**
- `args` - Variable number of arguments of any type

**Returns:** void

**Examples:**
```typescript
print("Hello")
print(" ")
print("world")
// Output: Hello world
```

**Go Mapping:** `fmt.Print()`

---

## Array Operations

### len

Returns the length of an array or string.

**Signature:**
```typescript
function len(arr: T[] | string): int
```

**Parameters:**
- `arr` - An array or string

**Returns:** int - The number of elements in the array or characters in the string

**Examples:**
```typescript
let arr = [1, 2, 3, 4, 5]
println(len(arr))  // 5

let str = "hello"
println(len(str))  // 5
```

**Go Mapping:** `len()`

---

### push

Appends an element to the end of an array.

**Signature:**
```typescript
function push<T>(arr: T[], element: T): void
```

**Parameters:**
- `arr` - The array to append to
- `element` - The element to append

**Returns:** void (modifies array in-place)

**Examples:**
```typescript
let numbers: int[] = [1, 2, 3]
push(numbers, 4)
println(numbers)  // [1, 2, 3, 4]

let names: string[] = ["Alice", "Bob"]
push(names, "Charlie")
println(names)  // ["Alice", "Bob", "Charlie"]
```

**Go Mapping:** `arr = append(arr, element)`

**Note:** In goTS, `push` modifies the array in-place. The array variable must be mutable (`let`, not `const`).

---

### pop

Removes and returns the last element from an array.

**Signature:**
```typescript
function pop<T>(arr: T[]): T
```

**Parameters:**
- `arr` - The array to pop from

**Returns:** T - The removed element

**Examples:**
```typescript
let numbers: int[] = [1, 2, 3, 4, 5]
let last = pop(numbers)
println(last)      // 5
println(numbers)   // [1, 2, 3, 4]

let names: string[] = ["Alice", "Bob", "Charlie"]
let name = pop(names)
println(name)   // "Charlie"
```

**Go Mapping:**
```go
last := arr[len(arr)-1]
arr = arr[:len(arr)-1]
```

**Note:** Behavior is undefined if the array is empty. Always check `len(arr) > 0` before calling `pop`.

---

## Type Operations

### typeof

Returns the type of a value as a string.

**Signature:**
```typescript
function typeof(value: any): string
```

**Parameters:**
- `value` - Any value

**Returns:** string - Type name: "int", "float", "string", "bool", "array", "object", "function", "null"

**Examples:**
```typescript
println(typeof(42))           // "int"
println(typeof(3.14))         // "float"
println(typeof("hello"))      // "string"
println(typeof(true))         // "bool"
println(typeof([1, 2, 3]))    // "array"
println(typeof(null))         // "null"

let f = function(): void {}
println(typeof(f))            // "function"
```

**Go Mapping:** Uses `reflect.TypeOf()` with custom type name mapping

---

## Type Conversion

### tostring

Converts a value to a string.

**Signature:**
```typescript
function tostring(value: any): string
```

**Parameters:**
- `value` - Value to convert (int, float, boolean, string, etc.)

**Returns:** string - String representation of the value

**Examples:**
```typescript
println(tostring(42))          // "42"
println(tostring(3.14))        // "3.14"
println(tostring(true))        // "true"
println(tostring(false))       // "false"

let x = 100
let msg = "The value is " + tostring(x)
println(msg)  // "The value is 100"
```

**Go Mapping:** `fmt.Sprint()` or `strconv` functions

---

### toint

Converts a value to an integer.

**Signature:**
```typescript
function toint(value: string | float | int): int
```

**Parameters:**
- `value` - String, float, or int to convert

**Returns:** int - Integer value

**Examples:**
```typescript
println(toint("42"))      // 42
println(toint(3.14))      // 3
println(toint(3.99))      // 3  (truncates decimal)
println(toint(100))       // 100

// Parse strings
let num = toint("123")
println(num + 1)  // 124
```

**Go Mapping:**
- String: `strconv.Atoi()`
- Float: `int(value)`
- Int: identity

**Note:** Returns 0 if string parsing fails.

---

### tofloat

Converts a value to a floating-point number.

**Signature:**
```typescript
function tofloat(value: string | int | float): float
```

**Parameters:**
- `value` - String, int, or float to convert

**Returns:** float - Floating-point value

**Examples:**
```typescript
println(tofloat("3.14"))   // 3.14
println(tofloat(42))       // 42.0
println(tofloat(10))       // 10.0

// Parse strings
let pi = tofloat("3.14159")
println(pi * 2)  // 6.28318
```

**Go Mapping:**
- String: `strconv.ParseFloat()`
- Int: `float64(value)`
- Float: identity

**Note:** Returns 0.0 if string parsing fails.

---

## Math Functions

### sqrt

Calculates the square root of a number.

**Signature:**
```typescript
function sqrt(x: float | int): float
```

**Parameters:**
- `x` - A non-negative number

**Returns:** float - Square root of x

**Examples:**
```typescript
println(sqrt(16))      // 4.0
println(sqrt(2))       // 1.414...
println(sqrt(100.0))   // 10.0

// Pythagorean theorem
let a = 3.0
let b = 4.0
let c = sqrt(a*a + b*b)
println(c)  // 5.0
```

**Go Mapping:** `math.Sqrt()`

---

### floor

Returns the largest integer less than or equal to a number.

**Signature:**
```typescript
function floor(x: float | int): float
```

**Parameters:**
- `x` - A number

**Returns:** float - Largest integer ≤ x (as float)

**Examples:**
```typescript
println(floor(3.7))    // 3.0
println(floor(3.2))    // 3.0
println(floor(-2.3))   // -3.0
println(floor(5))      // 5.0
```

**Go Mapping:** `math.Floor()`

---

### ceil

Returns the smallest integer greater than or equal to a number.

**Signature:**
```typescript
function ceil(x: float | int): float
```

**Parameters:**
- `x` - A number

**Returns:** float - Smallest integer ≥ x (as float)

**Examples:**
```typescript
println(ceil(3.2))    // 4.0
println(ceil(3.7))    // 4.0
println(ceil(-2.3))   // -2.0
println(ceil(5))      // 5.0
```

**Go Mapping:** `math.Ceil()`

---

### abs

Returns the absolute value of a number.

**Signature:**
```typescript
function abs(x: float | int): float
```

**Parameters:**
- `x` - A number

**Returns:** float - Absolute value of x

**Examples:**
```typescript
println(abs(-5))      // 5.0
println(abs(3.14))    // 3.14
println(abs(-2.5))    // 2.5
println(abs(0))       // 0.0

let distance = abs(x2 - x1)
```

**Go Mapping:** `math.Abs()`

---

## Function Reference Table

| Function | Purpose | Parameters | Return Type |
|----------|---------|------------|-------------|
| `println` | Print with newline | `...any` | `void` |
| `print` | Print without newline | `...any` | `void` |
| `len` | Get array/string length | `T[] \| string` | `int` |
| `push` | Append to array | `T[], T` | `void` |
| `pop` | Remove last element | `T[]` | `T` |
| `typeof` | Get type name | `any` | `string` |
| `tostring` | Convert to string | `any` | `string` |
| `toint` | Convert to int | `string \| float \| int` | `int` |
| `tofloat` | Convert to float | `string \| int \| float` | `float` |
| `sqrt` | Square root | `float \| int` | `float` |
| `floor` | Round down | `float \| int` | `float` |
| `ceil` | Round up | `float \| int` | `float` |
| `abs` | Absolute value | `float \| int` | `float` |

## Usage Notes

### Type Flexibility

Most built-in functions that accept numeric arguments work with both `int` and `float`:

```typescript
println(sqrt(16))      // int argument
println(sqrt(16.0))    // float argument

println(abs(-5))       // int argument
println(abs(-5.5))     // float argument
```

### Array Mutability

Functions that modify arrays (`push`, `pop`) require the array variable to be declared with `let`:

```typescript
let arr: int[] = [1, 2, 3]
push(arr, 4)  // OK

const fixed: int[] = [1, 2, 3]
// push(fixed, 4)  // Error: cannot modify const
```

### Type Conversion Chain

You can chain type conversions:

```typescript
let str = "42"
let num = toint(str)
let result = num * 2
println(tostring(result))  // "84"
```

### Dynamic Typing with `any`

When working with `Function` or `any` types, use type checking and conversion:

```typescript
function process(value: any): void {
    if (typeof(value) == "int") {
        println("Integer: " + tostring(value))
    } else if (typeof(value) == "string") {
        println("String: " + value)
    }
}
```

## Runtime Helpers

The following functions are generated by the goTS compiler for internal use:

- `gts_call` - Dynamic function invocation
- `gts_toint` - Runtime integer conversion
- `gts_tofloat` - Runtime float conversion
- `gts_tostring` - Runtime string conversion

These are not directly callable from goTS code but are used in generated Go code for type coercion and dynamic operations.
