---
layout: home

hero:
  name: 'goTS'
  text: 'TypeScript-like, compiles to Go'
  tagline: Write TypeScript-flavored code. Get a native, statically typed Go binary — with generics, async/await, SQL, and full access to the Go ecosystem.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: Language Spec
      link: /language-spec
    - theme: alt
      text: View on GitHub
      link: https://github.com/zhy0216/quickts

features:
  - icon: ⚡
    title: Compiles to Go
    details: The full pipeline — lexer, parser, type checker, and Go code generator — produces idiomatic Go source, then hands off to the Go toolchain for native binaries.
    link: /development/architecture
  - icon: 🧬
    title: Modern Type System
    details: Generics, interfaces with structural typing, union and intersection types, literal types, tuples, null safety, and type inference.
    link: /language-spec#3-types
  - icon: ⏳
    title: async/await & Event Loop
    details: Promises, async functions, setTimeout, setInterval, and queueMicrotask — an event loop running on the Go scheduler.
    link: /language-spec#18-async-await-and-the-event-loop
  - icon: 🗄️
    title: Built-in SQL
    details: 'Typed SQL via tagged template literals and generics, backed by SQLite. db.sql<User[]>`SELECT ...` returns typed rows.'
    link: /built-in-reference#sql-database-connect
  - icon: 🔗
    title: Go Interop
    details: 'Import functions straight from the Go standard library: import { ToUpper } from "go:strings".'
    link: /language-spec#20-go-interop
  - icon: 🧰
    title: Rich Standard Library
    details: Math, Number, JSON, Date, Object, Map, Set, RegExp, plus a deep set of String and Array methods.
    link: /built-in-reference
---

<div style="text-align: center; margin: 3rem 0;">

## Write TypeScript. Ship Go.

</div>

```typescript
// hello.gts — compile with: gots run hello.gts

function factorial(n: int): int {
    if (n <= 1) { return 1 }
    return n * factorial(n - 1)
}

class Animal {
    name: string
    constructor(name: string) { this.name = name }
    speak(): void { println(this.name) }
}

class Dog extends Animal {
    speak(): void { println(this.name + " barks") }
}

async function compute(): Promise<int> {
    let a: int = await Promise.resolve(42)
    return a * 2
}

let dog = new Dog("Rex")
dog.speak()
println(factorial(5))

const db = connect("./app.db")
db.sql`CREATE TABLE IF NOT EXISTS users (name TEXT NOT NULL)`
```

```go
// The generated Go — idiomatic, static, fast

package main

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return (n * factorial(n-1))
}

type Animal struct {
	Name string
}

func NewAnimal(name string) *Animal {
	this := &Animal{}
	this.Name = name
	return this
}
```

## Quick Start

```bash
git clone https://github.com/zhy0216/quickts
cd quickts/gots
make build          # build the gots binary

./gots run test/example.gts
./gots repl
```

## Documentation Map

| Page | Contents |
|------|----------|
| [Getting Started](/guide/getting-started) | Install, build, run your first program |
| [CLI Reference](/guide/cli) | All `gots` commands and flags |
| [Language Spec](/language-spec) | Complete language reference: types, generics, async, SQL, Go interop |
| [Built-in Reference](/built-in-reference) | String, Array, Map, Set, Date, Math, Number, JSON, RegExp, Promise |
| [Architecture](/development/architecture) | Compiler pipeline and package layout |
| [Testing](/development/testing) | Make targets, unit/integration/test262 tests |
| [FAQ](/faq) | Differences from TypeScript/JavaScript |
