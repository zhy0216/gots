---
title: "Type System"
description: "Understanding goTS types"
order: 1
category: "Core Concepts"
---

# Type System

goTS provides a static type system that compiles to Go's type system.

## Primitive Types

### int - Integer Type

Represents integer values. Maps to Go's `int` type. Used for whole numbers and array indexing.

```typescript
let count: int = 42
let index: int = 0
```

### float - Floating Point

Represents floating-point numbers. Maps to Go's `float64` type. Division (`/`) always returns float.

```typescript
let pi: float = 3.14159
let ratio: float = 22 / 7  // Division always returns float
```

### string - String Type

Sequences of characters. Maps to Go's `string` type. Supports concatenation with `+` operator.

```typescript
let name: string = "goTS"
let greeting: string = "Hello, " + name
```

### boolean - Boolean Type

Represents logical values true and false. Maps to Go's `bool` type.

```typescript
let isValid: boolean = true
let hasError: boolean = false
```

## Complex Types

### Arrays

```typescript
let numbers: int[] = [1, 2, 3, 4, 5]
let names: string[] = ["Alice", "Bob", "Charlie"]
```

### Nullable Types

```typescript
let value: string | null = null
value = "hello"
```

### Function Types

```typescript
let add: (a: int, b: int) => int = function(a: int, b: int): int {
  return a + b
}
```
