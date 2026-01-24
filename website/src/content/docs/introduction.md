---
title: "Introduction"
description: "Welcome to goTS - A TypeScript-like language that compiles to Go"
order: 1
category: "Getting Started"
---

# Introduction to goTS

goTS (Go-TypeScript) is a statically-typed programming language with TypeScript-like syntax that compiles to Go. It provides a familiar development experience for TypeScript developers while leveraging the Go ecosystem and toolchain.

## What is goTS?

goTS is designed to be a strict subset of TypeScript with the following characteristics:

- **Valid goTS code** should be syntactically valid TypeScript
- **goTS enforces stricter rules** (e.g., all variables require explicit type annotations)
- **goTS uses int and float** as distinct numeric types

## Quick Example

```typescript
let x: int = 42
let pi: float = 3.14159
let name: string = "goTS"
println("Hello from " + name + "!")
```
