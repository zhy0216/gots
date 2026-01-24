---
title: "Functions & Closures"
description: "Working with functions and closures in goTS"
order: 2
category: "Core Concepts"
---

# Functions & Closures

goTS supports first-class functions and closures, similar to TypeScript.

## Function Declarations

```typescript
function add(a: int, b: int): int {
  return a + b
}

function greet(name: string): void {
  println("Hello, " + name)
}
```

## Function Expressions

```typescript
let multiply: (a: int, b: int) => int = function(a: int, b: int): int {
  return a * b
}
```

## Closures

Functions can capture variables from their enclosing scope:

```typescript
function makeCounter(): Function {
  let count: int = 0
  return function(): int {
    count = count + 1
    return count
  }
}

let counter: Function = makeCounter()
```

## Higher-Order Functions

Functions that take or return other functions:

```typescript
function curry_add(a: int): Function {
  return function(b: int): int {
    return a + b
  }
}

let add5: Function = curry_add(5)
println(add5(3))  // Prints: 8
```
