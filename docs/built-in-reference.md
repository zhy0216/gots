# goTS Built-in Functions and Objects Reference

**Version 2.0**
**January 2026**

This document describes the built-in functions and objects available in goTS, designed to closely mirror JavaScript's standard library while being adapted for the goTS type system and compilation to Go.

---

## Table of Contents

1. [Overview](#overview)
2. [Console I/O](#console-io)
3. [Array](#array)
4. [Map](#map)
5. [Set](#set)
6. [String Functions](#string-functions)
7. [Type Utilities](#type-utilities)
8. [Math Functions](#math-functions)
9. [Comparison with JavaScript](#comparison-with-javascript)

---

## Overview

goTS provides built-in objects and functions that closely follow JavaScript conventions:

- **Console I/O**: Print output to the console
- **Array**: Full-featured array object with methods like JavaScript
- **Map**: Key-value collection with any key type
- **Set**: Collection of unique values
- **Type conversion**: Convert between types
- **Type introspection**: Check runtime types
- **Math operations**: Common mathematical functions

---

## Console I/O

### println

Prints a value to the console followed by a newline.

**Signature:**
```typescript
println(value: any): void
```

**Parameters:**
- `value`: The value to print (can be of any type)

**Return Value:**
- `void`

**Examples:**
```typescript
println("Hello, World!")           // Prints: Hello, World!
println(42)                        // Prints: 42
println(3.14)                      // Prints: 3.14
println(true)                      // Prints: true
println([1, 2, 3])                 // Prints: [1, 2, 3]
```

**Notes:**
- Automatically converts values to string representation
- Adds a newline character after the output
- Similar to JavaScript's `console.log()`

---

### print

Prints a value to the console without a newline.

**Signature:**
```typescript
print(value: any): void
```

**Parameters:**
- `value`: The value to print (can be of any type)

**Return Value:**
- `void`

**Examples:**
```typescript
print("Hello")
print(" ")
print("World")
println("!")
// Output: Hello World!
```

---

## Array

Arrays in goTS are objects with methods, similar to JavaScript arrays. They are generic and type-safe.

### Creating Arrays

```typescript
// Array literal
let numbers: int[] = [1, 2, 3, 4, 5]
let names: string[] = ["Alice", "Bob", "Charlie"]
let empty: int[] = []

// Mixed initialization
let items: int[] = [1, 2]
items.push(3)
```

### Properties

#### length

Returns the number of elements in the array.

**Type:** `int` (read-only)

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
println(arr.length)                // Prints: 5

let empty: string[] = []
println(empty.length)              // Prints: 0
```

---

### Methods

#### push

Appends one or more elements to the end of the array and returns the new length.

**Signature:**
```typescript
arr.push(value: T): int
arr.push(value1: T, value2: T, ...): int
```

**Parameters:**
- `value`: The value(s) to append

**Return Value:**
- `int`: The new length of the array

**Examples:**
```typescript
let numbers: int[] = [1, 2, 3]
let newLen: int = numbers.push(4)
println(newLen)                    // Prints: 4
println(numbers)                   // Prints: [1, 2, 3, 4]

numbers.push(5, 6, 7)
println(numbers)                   // Prints: [1, 2, 3, 4, 5, 6, 7]
```

---

#### pop

Removes the last element from the array and returns it.

**Signature:**
```typescript
arr.pop(): T | null
```

**Return Value:**
- `T | null`: The removed element, or `null` if the array is empty

**Examples:**
```typescript
let numbers: int[] = [1, 2, 3, 4, 5]
let last: int | null = numbers.pop()
println(last)                      // Prints: 5
println(numbers)                   // Prints: [1, 2, 3, 4]

let empty: int[] = []
let result: int | null = empty.pop()
println(result)                    // Prints: null
```

---

#### shift

Removes the first element from the array and returns it.

**Signature:**
```typescript
arr.shift(): T | null
```

**Return Value:**
- `T | null`: The removed element, or `null` if the array is empty

**Examples:**
```typescript
let numbers: int[] = [1, 2, 3, 4, 5]
let first: int | null = numbers.shift()
println(first)                     // Prints: 1
println(numbers)                   // Prints: [2, 3, 4, 5]
```

---

#### unshift

Adds one or more elements to the beginning of the array and returns the new length.

**Signature:**
```typescript
arr.unshift(value: T): int
arr.unshift(value1: T, value2: T, ...): int
```

**Return Value:**
- `int`: The new length of the array

**Examples:**
```typescript
let numbers: int[] = [3, 4, 5]
let newLen: int = numbers.unshift(1, 2)
println(newLen)                    // Prints: 5
println(numbers)                   // Prints: [1, 2, 3, 4, 5]
```

---

#### slice

Returns a shallow copy of a portion of the array.

**Signature:**
```typescript
arr.slice(): T[]
arr.slice(start: int): T[]
arr.slice(start: int, end: int): T[]
```

**Parameters:**
- `start`: Starting index (inclusive). If negative, counts from end. Defaults to 0.
- `end`: Ending index (exclusive). If negative, counts from end. Defaults to array length.

**Return Value:**
- `T[]`: A new array containing the extracted elements

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
println(arr.slice(1, 4))           // Prints: [2, 3, 4]
println(arr.slice(2))              // Prints: [3, 4, 5]
println(arr.slice(-2))             // Prints: [4, 5]
println(arr.slice())               // Prints: [1, 2, 3, 4, 5] (copy)
```

---

#### splice

Changes the contents of an array by removing or replacing existing elements and/or adding new elements.

**Signature:**
```typescript
arr.splice(start: int): T[]
arr.splice(start: int, deleteCount: int): T[]
arr.splice(start: int, deleteCount: int, item1: T, item2: T, ...): T[]
```

**Parameters:**
- `start`: The index at which to start changing the array
- `deleteCount`: The number of elements to remove (defaults to all remaining)
- `items`: Elements to add to the array

**Return Value:**
- `T[]`: An array containing the deleted elements

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
let removed: int[] = arr.splice(2, 2)
println(removed)                   // Prints: [3, 4]
println(arr)                       // Prints: [1, 2, 5]

arr.splice(1, 0, 10, 20)           // Insert without removing
println(arr)                       // Prints: [1, 10, 20, 2, 5]
```

---

#### concat

Merges two or more arrays into a new array.

**Signature:**
```typescript
arr.concat(arr2: T[]): T[]
arr.concat(arr2: T[], arr3: T[], ...): T[]
```

**Return Value:**
- `T[]`: A new array with all elements combined

**Examples:**
```typescript
let a: int[] = [1, 2]
let b: int[] = [3, 4]
let c: int[] = [5, 6]
let combined: int[] = a.concat(b, c)
println(combined)                  // Prints: [1, 2, 3, 4, 5, 6]
```

---

#### indexOf

Returns the first index at which a given element is found.

**Signature:**
```typescript
arr.indexOf(value: T): int
arr.indexOf(value: T, fromIndex: int): int
```

**Parameters:**
- `value`: The element to search for
- `fromIndex`: The index to start searching from (defaults to 0)

**Return Value:**
- `int`: The index of the element, or `-1` if not found

**Examples:**
```typescript
let arr: string[] = ["apple", "banana", "cherry", "banana"]
println(arr.indexOf("banana"))     // Prints: 1
println(arr.indexOf("banana", 2))  // Prints: 3
println(arr.indexOf("grape"))      // Prints: -1
```

---

#### includes

Determines whether the array includes a certain element.

**Signature:**
```typescript
arr.includes(value: T): boolean
arr.includes(value: T, fromIndex: int): boolean
```

**Return Value:**
- `boolean`: `true` if the element is found, `false` otherwise

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
println(arr.includes(3))           // Prints: true
println(arr.includes(6))           // Prints: false
```

---

#### join

Creates a string by concatenating all elements with a separator.

**Signature:**
```typescript
arr.join(): string
arr.join(separator: string): string
```

**Parameters:**
- `separator`: The string to separate elements (defaults to `","`)

**Return Value:**
- `string`: The joined string

**Examples:**
```typescript
let words: string[] = ["Hello", "World"]
println(words.join(" "))           // Prints: Hello World
println(words.join("-"))           // Prints: Hello-World
println(words.join())              // Prints: Hello,World

let nums: int[] = [1, 2, 3]
println(nums.join(" + "))          // Prints: 1 + 2 + 3
```

---

#### reverse

Reverses the array in place and returns it.

**Signature:**
```typescript
arr.reverse(): T[]
```

**Return Value:**
- `T[]`: The reversed array (same reference)

**Examples:**
```typescript
let arr: int[] = [1, 2, 3, 4, 5]
arr.reverse()
println(arr)                       // Prints: [5, 4, 3, 2, 1]
```

---

#### sort

Sorts the array in place and returns it.

**Signature:**
```typescript
arr.sort(): T[]
arr.sort(compareFn: (a: T, b: T) => int): T[]
```

**Parameters:**
- `compareFn`: Optional comparison function. Should return:
  - Negative number if `a` should come before `b`
  - Positive number if `a` should come after `b`
  - Zero if they are equal

**Return Value:**
- `T[]`: The sorted array (same reference)

**Examples:**
```typescript
let nums: int[] = [3, 1, 4, 1, 5, 9, 2, 6]
nums.sort()
println(nums)                      // Prints: [1, 1, 2, 3, 4, 5, 6, 9]

// Custom comparator (descending)
nums.sort(function(a: int, b: int): int {
    return b - a
})
println(nums)                      // Prints: [9, 6, 5, 4, 3, 2, 1, 1]
```

---

#### map

Creates a new array with the results of calling a function on every element.

**Signature:**
```typescript
arr.map(callback: (value: T, index: int) => U): U[]
arr.map(callback: (value: T) => U): U[]
```

**Parameters:**
- `callback`: Function that produces an element of the new array

**Return Value:**
- `U[]`: A new array with transformed elements

**Examples:**
```typescript
let nums: int[] = [1, 2, 3, 4, 5]
let doubled: int[] = nums.map(function(x: int): int {
    return x * 2
})
println(doubled)                   // Prints: [2, 4, 6, 8, 10]

let strs: string[] = nums.map(function(x: int): string {
    return "num:" + tostring(x)
})
println(strs)                      // Prints: [num:1, num:2, num:3, num:4, num:5]
```

---

#### filter

Creates a new array with all elements that pass the test.

**Signature:**
```typescript
arr.filter(callback: (value: T, index: int) => boolean): T[]
arr.filter(callback: (value: T) => boolean): T[]
```

**Parameters:**
- `callback`: Function that tests each element

**Return Value:**
- `T[]`: A new array with elements that pass the test

**Examples:**
```typescript
let nums: int[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
let evens: int[] = nums.filter(function(x: int): boolean {
    return x % 2 == 0
})
println(evens)                     // Prints: [2, 4, 6, 8, 10]
```

---

#### reduce

Executes a reducer function on each element, resulting in a single value.

**Signature:**
```typescript
arr.reduce(callback: (accumulator: U, value: T, index: int) => U, initialValue: U): U
arr.reduce(callback: (accumulator: U, value: T) => U, initialValue: U): U
```

**Parameters:**
- `callback`: Function that combines accumulator and current value
- `initialValue`: Initial value for the accumulator

**Return Value:**
- `U`: The final accumulated value

**Examples:**
```typescript
let nums: int[] = [1, 2, 3, 4, 5]
let sum: int = nums.reduce(function(acc: int, x: int): int {
    return acc + x
}, 0)
println(sum)                       // Prints: 15

let product: int = nums.reduce(function(acc: int, x: int): int {
    return acc * x
}, 1)
println(product)                   // Prints: 120
```

---

#### forEach

Executes a function for each array element.

**Signature:**
```typescript
arr.forEach(callback: (value: T, index: int) => void): void
arr.forEach(callback: (value: T) => void): void
```

**Parameters:**
- `callback`: Function to execute for each element

**Examples:**
```typescript
let names: string[] = ["Alice", "Bob", "Charlie"]
names.forEach(function(name: string, i: int): void {
    println(tostring(i) + ": " + name)
})
// Output:
// 0: Alice
// 1: Bob
// 2: Charlie
```

---

#### find

Returns the first element that satisfies the testing function.

**Signature:**
```typescript
arr.find(callback: (value: T, index: int) => boolean): T | null
arr.find(callback: (value: T) => boolean): T | null
```

**Return Value:**
- `T | null`: The first matching element, or `null` if not found

**Examples:**
```typescript
let nums: int[] = [1, 5, 10, 15, 20]
let found: int | null = nums.find(function(x: int): boolean {
    return x > 8
})
println(found)                     // Prints: 10
```

---

#### findIndex

Returns the index of the first element that satisfies the testing function.

**Signature:**
```typescript
arr.findIndex(callback: (value: T, index: int) => boolean): int
arr.findIndex(callback: (value: T) => boolean): int
```

**Return Value:**
- `int`: The index of the first matching element, or `-1` if not found

**Examples:**
```typescript
let nums: int[] = [1, 5, 10, 15, 20]
let idx: int = nums.findIndex(function(x: int): boolean {
    return x > 8
})
println(idx)                       // Prints: 2
```

---

#### some

Tests whether at least one element passes the test.

**Signature:**
```typescript
arr.some(callback: (value: T, index: int) => boolean): boolean
arr.some(callback: (value: T) => boolean): boolean
```

**Return Value:**
- `boolean`: `true` if at least one element passes, `false` otherwise

**Examples:**
```typescript
let nums: int[] = [1, 2, 3, 4, 5]
let hasEven: boolean = nums.some(function(x: int): boolean {
    return x % 2 == 0
})
println(hasEven)                   // Prints: true
```

---

#### every

Tests whether all elements pass the test.

**Signature:**
```typescript
arr.every(callback: (value: T, index: int) => boolean): boolean
arr.every(callback: (value: T) => boolean): boolean
```

**Return Value:**
- `boolean`: `true` if all elements pass, `false` otherwise

**Examples:**
```typescript
let nums: int[] = [2, 4, 6, 8, 10]
let allEven: boolean = nums.every(function(x: int): boolean {
    return x % 2 == 0
})
println(allEven)                   // Prints: true
```

---

## Map

Map is a collection of key-value pairs where keys can be of any type. It maintains insertion order.

### Creating Maps

```typescript
// Empty map
let scores: Map<string, int> = new Map<string, int>()

// Map with initial values
let ages: Map<string, int> = new Map<string, int>([
    ["Alice", 30],
    ["Bob", 25]
])
```

### Properties

#### size

Returns the number of key-value pairs in the map.

**Type:** `int` (read-only)

**Examples:**
```typescript
let map: Map<string, int> = new Map<string, int>()
map.set("a", 1)
map.set("b", 2)
println(map.size)                  // Prints: 2
```

---

### Methods

#### set

Adds or updates a key-value pair and returns the map.

**Signature:**
```typescript
map.set(key: K, value: V): Map<K, V>
```

**Return Value:**
- `Map<K, V>`: The map itself (for chaining)

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
       .set("Bob", 95)
       .set("Charlie", 88)

println(scores.size)               // Prints: 3
```

---

#### get

Returns the value associated with a key.

**Signature:**
```typescript
map.get(key: K): V | null
```

**Return Value:**
- `V | null`: The value if found, or `null` if the key doesn't exist

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)

let score: int | null = scores.get("Alice")
println(score)                     // Prints: 100

let missing: int | null = scores.get("Dave")
println(missing)                   // Prints: null
```

---

#### has

Returns whether a key exists in the map.

**Signature:**
```typescript
map.has(key: K): boolean
```

**Return Value:**
- `boolean`: `true` if the key exists, `false` otherwise

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)

println(scores.has("Alice"))       // Prints: true
println(scores.has("Bob"))         // Prints: false
```

---

#### delete

Removes a key-value pair from the map.

**Signature:**
```typescript
map.delete(key: K): boolean
```

**Return Value:**
- `boolean`: `true` if the key was found and removed, `false` otherwise

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)

let deleted: boolean = scores.delete("Alice")
println(deleted)                   // Prints: true
println(scores.has("Alice"))       // Prints: false

let notFound: boolean = scores.delete("Bob")
println(notFound)                  // Prints: false
```

---

#### clear

Removes all key-value pairs from the map.

**Signature:**
```typescript
map.clear(): void
```

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
scores.set("Bob", 95)

scores.clear()
println(scores.size)               // Prints: 0
```

---

#### keys

Returns an array of all keys in the map.

**Signature:**
```typescript
map.keys(): K[]
```

**Return Value:**
- `K[]`: An array containing all keys

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
scores.set("Bob", 95)

let names: string[] = scores.keys()
println(names)                     // Prints: [Alice, Bob]
```

---

#### values

Returns an array of all values in the map.

**Signature:**
```typescript
map.values(): V[]
```

**Return Value:**
- `V[]`: An array containing all values

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
scores.set("Bob", 95)

let allScores: int[] = scores.values()
println(allScores)                 // Prints: [100, 95]
```

---

#### entries

Returns an array of all key-value pairs as tuples.

**Signature:**
```typescript
map.entries(): [K, V][]
```

**Return Value:**
- `[K, V][]`: An array of key-value tuples

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
scores.set("Bob", 95)

let pairs: [string, int][] = scores.entries()
for (let pair of pairs) {
    println(pair[0] + ": " + tostring(pair[1]))
}
// Output:
// Alice: 100
// Bob: 95
```

---

#### forEach

Executes a function for each key-value pair.

**Signature:**
```typescript
map.forEach(callback: (value: V, key: K) => void): void
```

**Examples:**
```typescript
let scores: Map<string, int> = new Map<string, int>()
scores.set("Alice", 100)
scores.set("Bob", 95)

scores.forEach(function(value: int, key: string): void {
    println(key + " scored " + tostring(value))
})
// Output:
// Alice scored 100
// Bob scored 95
```

---

## Set

Set is a collection of unique values of any type. It maintains insertion order.

### Creating Sets

```typescript
// Empty set
let numbers: Set<int> = new Set<int>()

// Set with initial values
let names: Set<string> = new Set<string>(["Alice", "Bob", "Charlie"])

// Duplicates are automatically removed
let unique: Set<int> = new Set<int>([1, 2, 2, 3, 3, 3])
println(unique.size)               // Prints: 3
```

### Properties

#### size

Returns the number of values in the set.

**Type:** `int` (read-only)

**Examples:**
```typescript
let set: Set<int> = new Set<int>([1, 2, 3])
println(set.size)                  // Prints: 3
```

---

### Methods

#### add

Adds a value to the set and returns the set.

**Signature:**
```typescript
set.add(value: T): Set<T>
```

**Return Value:**
- `Set<T>`: The set itself (for chaining)

**Examples:**
```typescript
let numbers: Set<int> = new Set<int>()
numbers.add(1)
       .add(2)
       .add(3)
       .add(2)  // Duplicate, ignored

println(numbers.size)              // Prints: 3
```

---

#### has

Returns whether a value exists in the set.

**Signature:**
```typescript
set.has(value: T): boolean
```

**Return Value:**
- `boolean`: `true` if the value exists, `false` otherwise

**Examples:**
```typescript
let numbers: Set<int> = new Set<int>([1, 2, 3])
println(numbers.has(2))            // Prints: true
println(numbers.has(4))            // Prints: false
```

---

#### delete

Removes a value from the set.

**Signature:**
```typescript
set.delete(value: T): boolean
```

**Return Value:**
- `boolean`: `true` if the value was found and removed, `false` otherwise

**Examples:**
```typescript
let numbers: Set<int> = new Set<int>([1, 2, 3])
let deleted: boolean = numbers.delete(2)
println(deleted)                   // Prints: true
println(numbers.has(2))            // Prints: false
```

---

#### clear

Removes all values from the set.

**Signature:**
```typescript
set.clear(): void
```

**Examples:**
```typescript
let numbers: Set<int> = new Set<int>([1, 2, 3])
numbers.clear()
println(numbers.size)              // Prints: 0
```

---

#### values

Returns an array of all values in the set.

**Signature:**
```typescript
set.values(): T[]
```

**Return Value:**
- `T[]`: An array containing all values

**Examples:**
```typescript
let numbers: Set<int> = new Set<int>([3, 1, 2])
let arr: int[] = numbers.values()
println(arr)                       // Prints: [3, 1, 2] (insertion order)
```

---

#### forEach

Executes a function for each value.

**Signature:**
```typescript
set.forEach(callback: (value: T) => void): void
```

**Examples:**
```typescript
let names: Set<string> = new Set<string>(["Alice", "Bob", "Charlie"])
names.forEach(function(name: string): void {
    println("Hello, " + name + "!")
})
// Output:
// Hello, Alice!
// Hello, Bob!
// Hello, Charlie!
```

---

### Set Operations

#### Union (combining sets)

```typescript
let a: Set<int> = new Set<int>([1, 2, 3])
let b: Set<int> = new Set<int>([3, 4, 5])

let union: Set<int> = new Set<int>(a.values().concat(b.values()))
println(union.values())            // Prints: [1, 2, 3, 4, 5]
```

#### Intersection (common elements)

```typescript
let a: Set<int> = new Set<int>([1, 2, 3, 4])
let b: Set<int> = new Set<int>([3, 4, 5, 6])

let intersection: Set<int> = new Set<int>(
    a.values().filter(function(x: int): boolean {
        return b.has(x)
    })
)
println(intersection.values())     // Prints: [3, 4]
```

#### Difference (elements in A but not in B)

```typescript
let a: Set<int> = new Set<int>([1, 2, 3, 4])
let b: Set<int> = new Set<int>([3, 4, 5, 6])

let difference: Set<int> = new Set<int>(
    a.values().filter(function(x: int): boolean {
        return !b.has(x)
    })
)
println(difference.values())       // Prints: [1, 2]
```

---

## String Functions

### len (for strings)

Returns the number of bytes in a string.

**Signature:**
```typescript
len(str: string): int
```

**Examples:**
```typescript
println(len("hello"))              // Prints: 5
println(len(""))                   // Prints: 0
```

**Notes:**
- Returns byte count, not character count
- For UTF-8 strings with multi-byte characters, byte count may be larger

---

### String Concatenation

Strings can be concatenated using the `+` operator.

**Examples:**
```typescript
let greeting: string = "Hello, " + "World!"
println(greeting)                  // Prints: Hello, World!

let name: string = "Alice"
let message: string = "Hello, " + name + "!"
println(message)                   // Prints: Hello, Alice!
```

---

### String Indexing

Access individual characters using bracket notation.

**Examples:**
```typescript
let text: string = "hello"
println(text[0])                   // Prints: h
println(text[4])                   // Prints: o
```

---

## Type Utilities

### typeof

Returns the runtime type of a value as a string.

**Signature:**
```typescript
typeof(value: any): string
```

**Return Values:**
| Value Type | Result |
|------------|--------|
| Integer | `"int"` |
| Float | `"float"` |
| String | `"string"` |
| Boolean | `"boolean"` |
| `null` | `"null"` |
| Array | `"array"` |
| Map | `"map"` |
| Set | `"set"` |
| Object | `"object"` |
| Function | `"function"` |

**Examples:**
```typescript
println(typeof(42))                // Prints: int
println(typeof(3.14))              // Prints: float
println(typeof("hello"))           // Prints: string
println(typeof([1, 2, 3]))         // Prints: array
println(typeof(new Map<string, int>())) // Prints: map
println(typeof(new Set<int>()))    // Prints: set
```

---

### tostring

Converts a value to its string representation.

**Signature:**
```typescript
tostring(value: any): string
```

**Examples:**
```typescript
println(tostring(42))              // Prints: 42
println(tostring(3.14))            // Prints: 3.14
println(tostring(true))            // Prints: true
```

---

### toint

Converts a value to an integer.

**Signature:**
```typescript
toint(value: any): int
```

**Examples:**
```typescript
println(toint(3.14))               // Prints: 3
println(toint("42"))               // Prints: 42
println(toint(true))               // Prints: 1
```

---

### tofloat

Converts a value to a floating-point number.

**Signature:**
```typescript
tofloat(value: any): float
```

**Examples:**
```typescript
println(tofloat(42))               // Prints: 42.0
println(tofloat("3.14"))           // Prints: 3.14
```

---

## Math Functions

### sqrt

Returns the square root of a number.

**Signature:**
```typescript
sqrt(x: float): float
```

**Examples:**
```typescript
println(sqrt(16.0))                // Prints: 4
println(sqrt(2.0))                 // Prints: 1.4142135623730951
```

---

### floor

Returns the largest integer less than or equal to a number.

**Signature:**
```typescript
floor(x: float): float
```

**Examples:**
```typescript
println(floor(3.7))                // Prints: 3
println(floor(-2.3))               // Prints: -3
```

---

### ceil

Returns the smallest integer greater than or equal to a number.

**Signature:**
```typescript
ceil(x: float): float
```

**Examples:**
```typescript
println(ceil(3.2))                 // Prints: 4
println(ceil(-2.7))                // Prints: -2
```

---

### abs

Returns the absolute value of a number.

**Signature:**
```typescript
abs(x: float): float
```

**Examples:**
```typescript
println(abs(-5.0))                 // Prints: 5
println(abs(3.14))                 // Prints: 3.14
```

---

## Comparison with JavaScript

### Available in goTS

| goTS | JavaScript Equivalent |
|------|----------------------|
| `arr.length` | `arr.length` |
| `arr.push(x)` | `arr.push(x)` |
| `arr.pop()` | `arr.pop()` |
| `arr.shift()` | `arr.shift()` |
| `arr.unshift(x)` | `arr.unshift(x)` |
| `arr.slice()` | `arr.slice()` |
| `arr.splice()` | `arr.splice()` |
| `arr.concat()` | `arr.concat()` |
| `arr.indexOf(x)` | `arr.indexOf(x)` |
| `arr.includes(x)` | `arr.includes(x)` |
| `arr.join()` | `arr.join()` |
| `arr.reverse()` | `arr.reverse()` |
| `arr.sort()` | `arr.sort()` |
| `arr.map(fn)` | `arr.map(fn)` |
| `arr.filter(fn)` | `arr.filter(fn)` |
| `arr.reduce(fn, init)` | `arr.reduce(fn, init)` |
| `arr.forEach(fn)` | `arr.forEach(fn)` |
| `arr.find(fn)` | `arr.find(fn)` |
| `arr.findIndex(fn)` | `arr.findIndex(fn)` |
| `arr.some(fn)` | `arr.some(fn)` |
| `arr.every(fn)` | `arr.every(fn)` |
| `new Map<K,V>()` | `new Map()` |
| `map.set(k, v)` | `map.set(k, v)` |
| `map.get(k)` | `map.get(k)` |
| `map.has(k)` | `map.has(k)` |
| `map.delete(k)` | `map.delete(k)` |
| `map.clear()` | `map.clear()` |
| `map.keys()` | `[...map.keys()]` |
| `map.values()` | `[...map.values()]` |
| `map.entries()` | `[...map.entries()]` |
| `new Set<T>()` | `new Set()` |
| `set.add(v)` | `set.add(v)` |
| `set.has(v)` | `set.has(v)` |
| `set.delete(v)` | `set.delete(v)` |
| `set.clear()` | `set.clear()` |
| `set.values()` | `[...set.values()]` |

### Not Available in goTS

- `Array.from()`, `Array.isArray()`
- `arr.flat()`, `arr.flatMap()`
- `arr.reduceRight()`
- `arr.fill()`, `arr.copyWithin()`
- `arr.lastIndexOf()`
- String methods (`.split()`, `.replace()`, etc.)
- `Object.keys()`, `Object.values()`, `Object.entries()`
- `JSON.parse()`, `JSON.stringify()`
- `Promise`, `async/await`
- `RegExp`

---

**goTS Built-in Reference v2.0**
**© 2026**
