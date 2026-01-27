## String Methods

String values in goTS have built-in methods for common string operations.

### split

Splits a string into an array of substrings using a delimiter.

**Signature:**
```typescript
str.split(separator: string): string[]
```

**Parameters:**
- `separator` - The delimiter to split on

**Returns:** string[] - Array of substrings

**Examples:**
```typescript
let str: string = "a,b,c"
let parts: string[] = str.split(",")
println(len(parts))  // 3

let path = "folder/subfolder/file.txt"
let segments = path.split("/")
println(segments[0])  // "folder"
```

**Go Mapping:** `strings.Split()`

---

### replace

Replaces all occurrences of a substring with another string.

**Signature:**
```typescript
str.replace(old: string, new: string): string
```

**Parameters:**
- `old` - The substring to replace
- `new` - The replacement string

**Returns:** string - The string with replacements made

**Examples:**
```typescript
let str: string = "hello world"
let newStr: string = str.replace("world", "goTS")
println(newStr)  // "hello goTS"

let text = "I like apples and apples"
println(text.replace("apples", "oranges"))  // "I like oranges and oranges"
```

**Go Mapping:** `strings.Replace()`

---

### trim

Removes leading and trailing whitespace from a string.

**Signature:**
```typescript
str.trim(): string
```

**Returns:** string - The trimmed string

**Examples:**
```typescript
let str: string = "  hello  "
let trimmed: string = str.trim()
println(trimmed)  // "hello"

let input = "\n\ttext\t\n"
println(input.trim())  // "text"
```

**Go Mapping:** `strings.TrimSpace()`

---

### startsWith

Checks if a string starts with a given prefix.

**Signature:**
```typescript
str.startsWith(prefix: string): boolean
```

**Parameters:**
- `prefix` - The prefix to look for

**Returns:** boolean - true if the string starts with the prefix

**Examples:**
```typescript
let str: string = "hello world"
let starts: boolean = str.startsWith("hello")
println(starts)  // true

if (filename.startsWith("test_")) {
    println("Test file detected")
}
```

**Go Mapping:** `strings.HasPrefix()`

---

### endsWith

Checks if a string ends with a given suffix.

**Signature:**
```typescript
str.endsWith(suffix: string): boolean
```

**Parameters:**
- `suffix` - The suffix to look for

**Returns:** boolean - true if the string ends with the suffix

**Examples:**
```typescript
let str: string = "hello world"
let ends: boolean = str.endsWith("world")
println(ends)  // true

if (filename.endsWith(".gts")) {
    println("goTS source file")
}
```

**Go Mapping:** `strings.HasSuffix()`

---

### includes

Checks if a string contains a substring.

**Signature:**
```typescript
str.includes(substr: string): boolean
```

**Parameters:**
- `substr` - The substring to search for

**Returns:** boolean - true if the substring is found

**Examples:**
```typescript
let str: string = "hello world"
let has: boolean = str.includes("world")
println(has)  // true

if (email.includes("@")) {
    println("Valid email format")
}
```

**Go Mapping:** `strings.Contains()`

---

### toLowerCase

Converts a string to lowercase.

**Signature:**
```typescript
str.toLowerCase(): string
```

**Returns:** string - The lowercase string

**Examples:**
```typescript
let str: string = "HELLO WORLD"
let lower: string = str.toLowerCase()
println(lower)  // "hello world"
```

**Go Mapping:** `strings.ToLower()`

---

### toUpperCase

Converts a string to uppercase.

**Signature:**
```typescript
str.toUpperCase(): string
```

**Returns:** string - The uppercase string

**Examples:**
```typescript
let str: string = "hello world"
let upper: string = str.toUpperCase()
println(upper)  // "HELLO WORLD"
```

**Go Mapping:** `strings.ToUpper()`

---

### indexOf

Finds the index of the first occurrence of a substring.

**Signature:**
```typescript
str.indexOf(substring: string): int
```

**Parameters:**
- `substring` - The substring to search for

**Returns:** int - Index of first occurrence, or -1 if not found

**Examples:**
```typescript
let str: string = "hello world"
let idx: int = str.indexOf("world")
println(idx)  // 6
```

**Go Mapping:** `strings.Index()`

---
## Array Methods

Array values in goTS have built-in methods for functional programming operations.

### map

Transforms each element of an array using a callback function.

**Signature:**
```typescript
arr.map<T, U>(fn: (x: T) => U): U[]
```

**Parameters:**
- `fn` - Function to apply to each element

**Returns:** Array of transformed elements

**Examples:**
```typescript
let arr: int[] = [1, 2, 3]
let doubled: int[] = arr.map((x: int): int => x * 2)
println(doubled)  // [2, 4, 6]

let nums = [1, 2, 3, 4, 5]
let squares = nums.map((n: int): int => n * n)
```

**Go Mapping:** Inline loop with result array

---

### filter

Filters array elements based on a predicate function.

**Signature:**
```typescript
arr.filter<T>(fn: (x: T) => boolean): T[]
```

**Parameters:**
- `fn` - Predicate function returning true to keep element

**Returns:** Array containing only elements that pass the test

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let evens: int[] = arr.filter((x: int): boolean => x % 2 == 0)
println(evens)  // [2, 4]

let names = ["Alice", "Bob", "Anna"]
let aNames = names.filter((n: string): boolean => n.startsWith("A"))
```

**Go Mapping:** Inline loop with conditional append

---

### reduce

Reduces an array to a single value using an accumulator function.

**Signature:**
```typescript
arr.reduce<T, U>(fn: (acc: U, x: T) => U, initial: U): U
```

**Parameters:**
- `fn` - Reducer function taking accumulator and current element
- `initial` - Initial accumulator value

**Returns:** Final accumulated value

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let sum: int = arr.reduce((acc: int, x: int): int => acc + x, 0)
println(sum)  // 15

let product = [1, 2, 3, 4].reduce((acc: int, x: int): int => acc * x, 1)
println(product)  // 24
```

**Go Mapping:** Inline loop with accumulator

---

### find

Finds the first element that satisfies a predicate.

**Signature:**
```typescript
arr.find<T>(fn: (x: T) => boolean): T | null
```

**Parameters:**
- `fn` - Predicate function

**Returns:** First element that matches, or null if none found

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let found: int | null = arr.find((x: int): boolean => x > 3)
println(found)  // 4

let names = ["Alice", "Bob", "Charlie"]
let longName = names.find((n: string): boolean => len(n) > 5)
```

**Go Mapping:** Inline loop with early return

---

### findIndex

Finds the index of the first element that satisfies a predicate.

**Signature:**
```typescript
arr.findIndex<T>(fn: (x: T) => boolean): int
```

**Parameters:**
- `fn` - Predicate function

**Returns:** Index of first matching element, or -1 if none found

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let idx: int = arr.findIndex((x: int): boolean => x > 3)
println(idx)  // 3

if (arr.findIndex((x: int): boolean => x < 0) == -1) {
    println("No negative numbers")
}
```

**Go Mapping:** Inline loop with index tracking

---

### some

Tests whether at least one element passes the predicate.

**Signature:**
```typescript
arr.some<T>(fn: (x: T) => boolean): boolean
```

**Parameters:**
- `fn` - Predicate function

**Returns:** true if any element passes the test

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let hasEven: boolean = arr.some((x: int): boolean => x % 2 == 0)
println(hasEven)  // true

if (scores.some((s: int): boolean => s < 60)) {
    println("Some students failed")
}
```

**Go Mapping:** Inline loop with early return

---

### every

Tests whether all elements pass the predicate.

**Signature:**
```typescript
arr.every<T>(fn: (x: T) => boolean): boolean
```

**Parameters:**
- `fn` - Predicate function

**Returns:** true if all elements pass the test

**Examples:**
```typescript
let arr: int[] = [2, 4, 6, 8]
let allEven: boolean = arr.every((x: int): boolean => x % 2 == 0)
println(allEven)  // true

if (ages.every((a: int): boolean => a >= 18)) {
    println("All adults")
}
```

**Go Mapping:** Inline loop with early return

---

### join

Joins an array of strings into a single string using a separator.

**Signature:**
```typescript
arr.join(separator: string): string
```

**Parameters:**
- `separator` - The string to insert between elements

**Returns:** string - The joined string

**Examples:**
```typescript
let parts: string[] = ["a", "b", "c"]
let str: string = parts.join(",")
println(str)  // "a,b,c"

let words = ["Hello", "World"]
println(words.join(" "))  // "Hello World"
```

**Go Mapping:** `strings.Join()`

---

## Math Object

The `Math` object provides mathematical constants and functions, similar to JavaScript's Math object.

### Constants

#### Math.PI

The ratio of a circle's circumference to its diameter (approximately 3.14159).

```typescript
let pi: number = Math.PI
println(pi)  // 3.141592653589793
```

**Go Mapping:** `math.Pi`

---

#### Math.E

Euler's number, the base of natural logarithms (approximately 2.71828).

```typescript
let e: number = Math.E
println(e)  // 2.718281828459045
```

**Go Mapping:** `math.E`

---

### Rounding Methods

#### Math.round

Rounds a number to the nearest integer.

**Signature:**
```typescript
Math.round(x: number): number
```

**Note:** For half values, Go rounds away from zero while JavaScript rounds toward positive infinity.
- `Math.round(4.5)` returns `5`
- `Math.round(-4.5)` returns `-5` (in goTS/Go) vs `-4` (in JavaScript)

**Go Mapping:** `math.Round()`

---

#### Math.floor

Returns the largest integer less than or equal to a number.

**Signature:**
```typescript
Math.floor(x: number): number
```

**Examples:**
```typescript
Math.floor(4.7)   // 4
Math.floor(-4.7)  // -5
```

**Go Mapping:** `math.Floor()`

---

#### Math.ceil

Returns the smallest integer greater than or equal to a number.

**Signature:**
```typescript
Math.ceil(x: number): number
```

**Examples:**
```typescript
Math.ceil(4.3)   // 5
Math.ceil(-4.3)  // -4
```

**Go Mapping:** `math.Ceil()`

---

#### Math.trunc

Returns the integer part of a number (truncates toward zero).

**Signature:**
```typescript
Math.trunc(x: number): number
```

**Examples:**
```typescript
Math.trunc(4.7)   // 4
Math.trunc(-4.7)  // -4
```

**Go Mapping:** `math.Trunc()`

---

### Power and Root Methods

#### Math.pow

Returns the base raised to the exponent power.

**Signature:**
```typescript
Math.pow(base: number, exponent: number): number
```

**Examples:**
```typescript
Math.pow(2, 3)    // 8
Math.pow(2, 0.5)  // 1.41421... (square root of 2)
```

**Go Mapping:** `math.Pow()`

---

#### Math.sqrt

Returns the square root of a number.

**Signature:**
```typescript
Math.sqrt(x: number): number
```

**Go Mapping:** `math.Sqrt()`

---

#### Math.cbrt

Returns the cube root of a number.

**Signature:**
```typescript
Math.cbrt(x: number): number
```

**Go Mapping:** `math.Cbrt()`

---

#### Math.exp

Returns e raised to the power of x.

**Signature:**
```typescript
Math.exp(x: number): number
```

**Go Mapping:** `math.Exp()`

---

### Logarithmic Methods

#### Math.log

Returns the natural logarithm (base e) of a number.

**Signature:**
```typescript
Math.log(x: number): number
```

**Go Mapping:** `math.Log()`

---

#### Math.log10

Returns the base-10 logarithm of a number.

**Signature:**
```typescript
Math.log10(x: number): number
```

**Go Mapping:** `math.Log10()`

---

#### Math.log2

Returns the base-2 logarithm of a number.

**Signature:**
```typescript
Math.log2(x: number): number
```

**Go Mapping:** `math.Log2()`

---

### Absolute Value and Sign

#### Math.abs

Returns the absolute value of a number.

**Signature:**
```typescript
Math.abs(x: number): number
```

**Go Mapping:** `math.Abs()`

---

#### Math.sign

Returns the sign of a number (-1, 0, or 1).

**Signature:**
```typescript
Math.sign(x: number): number
```

**Examples:**
```typescript
Math.sign(5)   // 1
Math.sign(-5)  // -1
Math.sign(0)   // 0
```

**Go Mapping:** Inline function (Go doesn't have `math.Sign`)

---

### Min/Max (Variadic)

#### Math.min

Returns the smallest of the given numbers.

**Signature:**
```typescript
Math.min(...values: number[]): number
```

**Examples:**
```typescript
Math.min(1, 2)        // 1
Math.min(5, 3, 8)     // 3
Math.min(-1, -5)      // -5
```

**Go Mapping:** Chained `math.Min()` calls

---

#### Math.max

Returns the largest of the given numbers.

**Signature:**
```typescript
Math.max(...values: number[]): number
```

**Examples:**
```typescript
Math.max(1, 2)        // 2
Math.max(5, 3, 8)     // 8
Math.max(-1, -5)      // -1
```

**Go Mapping:** Chained `math.Max()` calls

---

### Trigonometric Methods

#### Math.sin

Returns the sine of an angle (in radians).

**Signature:**
```typescript
Math.sin(x: number): number
```

**Go Mapping:** `math.Sin()`

---

#### Math.cos

Returns the cosine of an angle (in radians).

**Signature:**
```typescript
Math.cos(x: number): number
```

**Go Mapping:** `math.Cos()`

---

#### Math.tan

Returns the tangent of an angle (in radians).

**Signature:**
```typescript
Math.tan(x: number): number
```

**Go Mapping:** `math.Tan()`

---

#### Math.asin

Returns the arcsine of a number (result in radians).

**Signature:**
```typescript
Math.asin(x: number): number
```

**Go Mapping:** `math.Asin()`

---

#### Math.acos

Returns the arccosine of a number (result in radians).

**Signature:**
```typescript
Math.acos(x: number): number
```

**Go Mapping:** `math.Acos()`

---

#### Math.atan

Returns the arctangent of a number (result in radians).

**Signature:**
```typescript
Math.atan(x: number): number
```

**Go Mapping:** `math.Atan()`

---

#### Math.atan2

Returns the arctangent of the quotient of its arguments.

**Signature:**
```typescript
Math.atan2(y: number, x: number): number
```

**Go Mapping:** `math.Atan2()`

---

### Random

#### Math.random

Returns a pseudo-random number between 0 (inclusive) and 1 (exclusive).

**Signature:**
```typescript
Math.random(): number
```

**Examples:**
```typescript
let r: number = Math.random()
// r is between 0.0 and 0.999...

// Generate random integer between 0 and 9
let n: int = toint(Math.random() * 10)
```

**Go Mapping:** `rand.Float64()`

---

## Number Object

The `Number` object provides methods for working with numbers, similar to JavaScript's Number object.

### Constants

#### Number.MAX_SAFE_INTEGER

The maximum safe integer in JavaScript (2^53 - 1).

```typescript
let max: number = Number.MAX_SAFE_INTEGER  // 9007199254740991
```

**Go Mapping:** `float64(9007199254740991)`

---

#### Number.MIN_SAFE_INTEGER

The minimum safe integer in JavaScript (-(2^53 - 1)).

```typescript
let min: number = Number.MIN_SAFE_INTEGER  // -9007199254740991
```

**Go Mapping:** `float64(-9007199254740991)`

---

#### Number.MAX_VALUE

The largest positive representable number.

```typescript
let maxVal: number = Number.MAX_VALUE
```

**Go Mapping:** `math.MaxFloat64`

---

#### Number.MIN_VALUE

The smallest positive representable number (closest to zero).

```typescript
let minVal: number = Number.MIN_VALUE
```

**Go Mapping:** `math.SmallestNonzeroFloat64`

---

#### Number.POSITIVE_INFINITY

Positive infinity.

```typescript
let inf: number = Number.POSITIVE_INFINITY
```

**Go Mapping:** `math.Inf(1)`

---

#### Number.NEGATIVE_INFINITY

Negative infinity.

```typescript
let negInf: number = Number.NEGATIVE_INFINITY
```

**Go Mapping:** `math.Inf(-1)`

---

#### Number.NaN

Not-a-Number value.

```typescript
let nan: number = Number.NaN
```

**Go Mapping:** `math.NaN()`

---

### Static Methods

#### Number.isFinite

Determines whether the passed value is a finite number.

**Signature:**
```typescript
Number.isFinite(x: number): boolean
```

**Examples:**
```typescript
Number.isFinite(42)        // true
Number.isFinite(3.14)      // true
Number.isFinite(Infinity)  // false
```

**Go Mapping:** `!math.IsInf(x, 0) && !math.IsNaN(x)`

---

#### Number.isNaN

Determines whether the passed value is NaN.

**Signature:**
```typescript
Number.isNaN(x: number): boolean
```

**Examples:**
```typescript
Number.isNaN(NaN)   // true
Number.isNaN(42)    // false
```

**Go Mapping:** `math.IsNaN()`

---

#### Number.isInteger

Determines whether the passed value is an integer.

**Signature:**
```typescript
Number.isInteger(x: number): boolean
```

**Examples:**
```typescript
Number.isInteger(42)     // true
Number.isInteger(42.0)   // true
Number.isInteger(3.14)   // false
```

**Go Mapping:** `math.Trunc(x) == x && !math.IsInf(x, 0)`

---

#### Number.isSafeInteger

Determines whether the passed value is a safe integer.

**Signature:**
```typescript
Number.isSafeInteger(x: number): boolean
```

**Examples:**
```typescript
Number.isSafeInteger(42)                     // true
Number.isSafeInteger(9007199254740992)       // false (too large)
```

**Go Mapping:** `math.Trunc(x) == x && math.Abs(x) <= 9007199254740991`

---

#### Number.parseFloat

Parses a string argument and returns a floating point number.

**Signature:**
```typescript
Number.parseFloat(s: string): number
```

**Examples:**
```typescript
Number.parseFloat("3.14")   // 3.14
Number.parseFloat("-5.5")   // -5.5
```

**Go Mapping:** `strconv.ParseFloat(s, 64)`

---

#### Number.parseInt

Parses a string argument and returns an integer of the specified radix.

**Signature:**
```typescript
Number.parseInt(s: string, radix?: int): int
```

**Parameters:**
- `s` - The string to parse
- `radix` - (Optional) The base for conversion, defaults to 10

**Examples:**
```typescript
Number.parseInt("42")         // 42
Number.parseInt("ff", 16)     // 255
Number.parseInt("1010", 2)    // 10
```

**Go Mapping:** `strconv.ParseInt(s, radix, 64)`

---

## Global Number Functions

### isNaN

Global function to check if a value is NaN.

**Signature:**
```typescript
isNaN(x: number): boolean
```

**Go Mapping:** `math.IsNaN()`

---

### isFinite

Global function to check if a value is finite.

**Signature:**
```typescript
isFinite(x: number): boolean
```

**Go Mapping:** `!math.IsInf(x, 0) && !math.IsNaN(x)`

---

### parseFloat

Global function to parse a string as a floating-point number.

**Signature:**
```typescript
parseFloat(s: string): number
```

**Go Mapping:** `strconv.ParseFloat(s, 64)`

---
