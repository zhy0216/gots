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
