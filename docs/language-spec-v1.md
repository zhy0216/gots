# GoTS Language Specification v1.0

A restricted TypeScript subset compiled to bytecode and executed on a Go VM.

---

## 1. Design Goals

- **Minimal**: Smallest feature set for practical programming
- **Static typing**: All types known at compile time
- **TypeScript syntax**: Valid GoTS is valid TypeScript (subset)
- **Simple implementation**: Easy to parse, compile, and execute

---

## 2. Types

### 2.1 Primitive Types

| Type      | Description                     | Example Literals        |
|-----------|---------------------------------|-------------------------|
| `number`  | 64-bit floating point           | `42`, `3.14`, `-7`      |
| `string`  | UTF-8 string                    | `"hello"`, `'world'`    |
| `boolean` | Boolean value                   | `true`, `false`         |
| `null`    | Absence of value                | `null`                  |
| `void`    | No value (function return only) | -                       |

### 2.2 Array Types

```typescript
let numbers: number[] = [1, 2, 3];
let names: string[] = ["alice", "bob"];
let matrix: number[][] = [[1, 2], [3, 4]];
```

### 2.3 Object Types

```typescript
let point: { x: number, y: number } = { x: 10, y: 20 };
let person: { name: string, age: number } = { name: "Alice", age: 30 };
```

### 2.4 Function Types

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};

let callback: () => void = function(): void {
    println("called");
};
```

### 2.5 Nullable Types

Use `| null` to allow null values:

```typescript
let name: string | null = null;
let value: number | null = 42;
```

### 2.6 Type Aliases

```typescript
type Point = { x: number, y: number };
type Callback = (value: number) => void;
type MaybeString = string | null;

let p: Point = { x: 1, y: 2 };
```

### 2.7 Type Annotations

```typescript
let x: number = 10;
let name: string = "alice";
let flag: boolean = true;
let items: number[] = [1, 2, 3];
```

**Note**: Type inference is NOT supported in v1. All declarations require explicit types.

---

## 3. Variables

### 3.1 Let Declarations

```typescript
let x: number = 0;
let message: string = "hello";
let data: number[] = [];
```

- `let` declares a mutable variable
- Initializer is **required**
- Type annotation is **required**

### 3.2 Const Declarations

```typescript
const PI: number = 3.14159;
const GREETING: string = "Hello";
const EMPTY: number[] = [];
```

- `const` declares an immutable binding
- Cannot be reassigned after declaration
- Note: Object/array contents can still be mutated

---

## 4. Expressions

### 4.1 Arithmetic Operators

| Operator | Description    | Example     |
|----------|----------------|-------------|
| `+`      | Addition       | `a + b`     |
| `-`      | Subtraction    | `a - b`     |
| `*`      | Multiplication | `a * b`     |
| `/`      | Division       | `a / b`     |
| `%`      | Modulo         | `a % b`     |
| `-`      | Unary negation | `-x`        |

**String concatenation**: `+` also concatenates strings.

```typescript
let result: string = "Hello, " + "World";
```

### 4.2 Comparison Operators

| Operator | Description              |
|----------|--------------------------|
| `==`     | Equal                    |
| `!=`     | Not equal                |
| `<`      | Less than                |
| `>`      | Greater than             |
| `<=`     | Less than or equal       |
| `>=`     | Greater than or equal    |

All comparisons return `boolean`.

### 4.3 Logical Operators

| Operator | Description | Example      |
|----------|-------------|--------------|
| `&&`     | Logical AND | `a && b`     |
| `\|\|`   | Logical OR  | `a \|\| b`   |
| `!`      | Logical NOT | `!flag`      |

### 4.4 Assignment

```typescript
x = 10;          // Simple assignment
x = x + 1;       // Compound (no += in v1)
```

### 4.5 Array Operations

```typescript
let arr: number[] = [1, 2, 3];

// Indexing (0-based)
let first: number = arr[0];
arr[1] = 10;

// Nested arrays
let matrix: number[][] = [[1, 2], [3, 4]];
let val: number = matrix[0][1];  // 2
```

### 4.6 Object Operations

```typescript
type Person = { name: string, age: number };
let person: Person = { name: "Alice", age: 30 };

// Property access
let n: string = person.name;
person.age = 31;
```

### 4.7 Grouping

```typescript
let result: number = (a + b) * c;
```

---

## 5. Statements

### 5.1 Expression Statement

Any expression followed by semicolon:

```typescript
x = x + 1;
print("hello");
```

### 5.2 Block Statement

```typescript
{
    let x: number = 10;
    print(x);
}
```

Blocks create new lexical scopes.

### 5.3 If Statement

```typescript
if (condition) {
    // then branch
}

if (condition) {
    // then branch
} else {
    // else branch
}

if (condition1) {
    // ...
} else if (condition2) {
    // ...
} else {
    // ...
}
```

**Note**: Braces are **required** (no single-statement bodies).

### 5.4 While Statement

```typescript
while (condition) {
    // body
}
```

### 5.5 For Statement

```typescript
for (let i: number = 0; i < 10; i = i + 1) {
    // body
}
```

All three clauses (init, condition, update) are **required**.

### 5.6 Break and Continue

```typescript
while (true) {
    if (done) {
        break;
    }
    if (skip) {
        continue;
    }
}
```

### 5.7 Return Statement

```typescript
return;          // void function
return value;    // returning a value
```

---

## 6. Functions

### 6.1 Function Declaration

```typescript
function add(a: number, b: number): number {
    return a + b;
}

function greet(name: string): void {
    println("Hello, " + name);
}

function noParams(): number {
    return 42;
}
```

- Return type annotation is **required**
- All parameters require type annotations
- Functions are hoisted (can be called before declaration)

### 6.2 Function Expressions

```typescript
let add: (a: number, b: number) => number = function(a: number, b: number): number {
    return a + b;
};

// Shorter with type alias
type BinaryOp = (a: number, b: number) => number;
let multiply: BinaryOp = function(a: number, b: number): number {
    return a * b;
};
```

### 6.3 Function Calls

```typescript
let sum: number = add(1, 2);
greet("World");

// Calling function expressions
let result: number = multiply(3, 4);
```

### 6.4 Closures

Functions capture variables from their enclosing scope:

```typescript
function makeCounter(): () => number {
    let count: number = 0;
    return function(): number {
        count = count + 1;
        return count;
    };
}

let counter: () => number = makeCounter();
println(counter());  // 1
println(counter());  // 2
println(counter());  // 3
```

### 6.5 Higher-Order Functions

Functions can take functions as parameters and return functions:

```typescript
function apply(f: (x: number) => number, value: number): number {
    return f(value);
}

function double(x: number): number {
    return x * 2;
}

let result: number = apply(double, 5);  // 10
```

---

## 7. Classes

### 7.1 Class Declaration

```typescript
class Point {
    x: number;
    y: number;

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    distance(other: Point): number {
        let dx: number = this.x - other.x;
        let dy: number = this.y - other.y;
        return sqrt(dx * dx + dy * dy);
    }
}
```

### 7.2 Instance Creation

```typescript
let p1: Point = new Point(0, 0);
let p2: Point = new Point(3, 4);
println(p1.distance(p2));  // 5
```

### 7.3 Class Features

- **Fields**: Declared at class level with types
- **Constructor**: Special method for initialization
- **Methods**: Functions bound to instances
- **this**: Reference to current instance

### 7.4 Inheritance

```typescript
class Animal {
    name: string;

    constructor(name: string) {
        this.name = name;
    }

    speak(): void {
        println(this.name + " makes a sound");
    }
}

class Dog extends Animal {
    constructor(name: string) {
        super(name);
    }

    speak(): void {
        println(this.name + " barks");
    }
}

let dog: Dog = new Dog("Rex");
dog.speak();  // "Rex barks"
```

- Single inheritance only
- `super(...)` must be first statement in derived constructor
- Methods can be overridden

### 7.5 Class as Type

```typescript
let animal: Animal = new Dog("Rex");  // Polymorphism
animal.speak();  // "Rex barks" (dynamic dispatch)
```

---

## 8. Null Handling

### 8.1 Nullable Types

```typescript
let name: string | null = null;
let value: number | null = 42;
```

### 8.2 Null Checks

```typescript
let name: string | null = getName();

if (name != null) {
    println(name);  // name is string here
}
```

### 8.3 Null Assignment

```typescript
let x: string | null = "hello";
x = null;  // OK

let y: string = "world";
y = null;  // ERROR: null not assignable to string
```

---

## 9. Built-in Functions

| Function    | Signature                        | Description                    |
|-------------|----------------------------------|--------------------------------|
| `print`     | `(value: any): void`             | Print value to stdout          |
| `println`   | `(value: any): void`             | Print value with newline       |
| `len`       | `(s: string \| T[]): number`     | Get string/array length        |
| `toString`  | `(n: number): string`            | Convert number to string       |
| `toNumber`  | `(s: string): number`            | Parse string as number         |
| `push`      | `(arr: T[], value: T): void`     | Append to array                |
| `pop`       | `(arr: T[]): T \| null`          | Remove and return last element |
| `sqrt`      | `(n: number): number`            | Square root                    |
| `floor`     | `(n: number): number`            | Floor to integer               |
| `ceil`      | `(n: number): number`            | Ceiling to integer             |
| `abs`       | `(n: number): number`            | Absolute value                 |

---

## 10. Comments

```typescript
// Single line comment

/*
   Multi-line
   comment
*/
```

---

## 11. Program Structure

A program is a sequence of:
1. Type aliases
2. Class declarations
3. Function declarations
4. Variable declarations (global scope)
5. Statements

Execution begins at the first statement in global scope (top-to-bottom).

```typescript
// Type alias
type Counter = () => number;

// Class
class Box {
    value: number;
    constructor(v: number) {
        this.value = v;
    }
}

// Function
function increment(b: Box): void {
    b.value = b.value + 1;
}

// Global variable
let box: Box = new Box(0);

// Main execution
increment(box);
increment(box);
println(box.value);  // 2
```

---

## 12. Lexical Structure

### 12.1 Keywords

```
let, const, function, return, if, else, while, for, break, continue,
true, false, null, number, string, boolean, void,
class, constructor, this, new, extends, super, type
```

### 12.2 Identifiers

- Start with letter or underscore
- Followed by letters, digits, or underscores
- Case-sensitive

### 12.3 Semicolons

Semicolons are **required** (no automatic semicolon insertion).

---

## 13. Operator Precedence (highest to lowest)

| Precedence | Operators                   | Associativity |
|------------|-----------------------------|---------------|
| 1          | `()` `[]` `.` (access)      | Left          |
| 2          | `new`                       | Right         |
| 3          | `!`, `-` (unary)            | Right         |
| 4          | `*`, `/`, `%`               | Left          |
| 5          | `+`, `-`                    | Left          |
| 6          | `<`, `>`, `<=`, `>=`        | Left          |
| 7          | `==`, `!=`                  | Left          |
| 8          | `&&`                        | Left          |
| 9          | `\|\|`                      | Left          |
| 10         | `=`                         | Right         |

---

## 14. Scoping Rules

- **Global scope**: Top-level declarations
- **Class scope**: Fields and methods
- **Function scope**: Parameters and local variables
- **Block scope**: Variables declared in `{}` blocks
- **Closure scope**: Captured variables from enclosing functions
- **Shadowing**: Inner scope can shadow outer scope names

```typescript
let x: number = 1;        // global

function foo(): void {
    let x: number = 2;    // shadows global x
    {
        let x: number = 3; // shadows function x
        print(x);          // 3
    }
    print(x);              // 2
}

print(x);                  // 1
```

---

## 15. Type Checking Rules

### 15.1 Assignment Compatibility

- Value type must match variable type exactly
- Subclass instances assignable to superclass variables
- `null` only assignable to nullable types (`T | null`)
- No implicit type conversions

### 15.2 Operator Type Rules

| Operation           | Operand Types        | Result Type |
|---------------------|----------------------|-------------|
| `+` (arithmetic)    | `number`, `number`   | `number`    |
| `+` (concat)        | `string`, `string`   | `string`    |
| `-`, `*`, `/`, `%`  | `number`, `number`   | `number`    |
| `<`, `>`, `<=`, `>=`| `number`, `number`   | `boolean`   |
| `==`, `!=`          | same type, same type | `boolean`   |
| `&&`, `\|\|`        | `boolean`, `boolean` | `boolean`   |
| `!`                 | `boolean`            | `boolean`   |
| `-` (unary)         | `number`             | `number`    |
| `[]` (index)        | `T[]`, `number`      | `T`         |
| `.` (property)      | `{...}` or class     | property type |

### 15.3 Null Safety

```typescript
let x: string | null = maybeGetString();

// Error: x might be null
println(x.length);  // COMPILE ERROR

// OK: null check first
if (x != null) {
    println(len(x));  // x is string here
}
```

---

## 16. Example Programs

### 16.1 Hello World

```typescript
println("Hello, World!");
```

### 16.2 Factorial

```typescript
function factorial(n: number): number {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}

println(factorial(5));  // 120
```

### 16.3 FizzBuzz

```typescript
function fizzbuzz(n: number): void {
    for (let i: number = 1; i <= n; i = i + 1) {
        if (i % 15 == 0) {
            println("FizzBuzz");
        } else if (i % 3 == 0) {
            println("Fizz");
        } else if (i % 5 == 0) {
            println("Buzz");
        } else {
            println(i);
        }
    }
}

fizzbuzz(20);
```

### 16.4 Linked List

```typescript
class Node {
    value: number;
    next: Node | null;

    constructor(value: number) {
        this.value = value;
        this.next = null;
    }
}

class LinkedList {
    head: Node | null;

    constructor() {
        this.head = null;
    }

    append(value: number): void {
        let node: Node = new Node(value);
        if (this.head == null) {
            this.head = node;
        } else {
            let current: Node | null = this.head;
            while (current != null && current.next != null) {
                current = current.next;
            }
            if (current != null) {
                current.next = node;
            }
        }
    }

    print(): void {
        let current: Node | null = this.head;
        while (current != null) {
            println(current.value);
            current = current.next;
        }
    }
}

let list: LinkedList = new LinkedList();
list.append(1);
list.append(2);
list.append(3);
list.print();
```

### 16.5 Higher-Order Functions

```typescript
type Predicate = (n: number) => boolean;

function filter(arr: number[], pred: Predicate): number[] {
    let result: number[] = [];
    for (let i: number = 0; i < len(arr); i = i + 1) {
        if (pred(arr[i])) {
            push(result, arr[i]);
        }
    }
    return result;
}

function isEven(n: number): boolean {
    return n % 2 == 0;
}

let numbers: number[] = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
let evens: number[] = filter(numbers, isEven);

for (let i: number = 0; i < len(evens); i = i + 1) {
    println(evens[i]);  // 2, 4, 6, 8, 10
}
```

### 16.6 Closure Counter

```typescript
type Counter = { increment: () => number, decrement: () => number, get: () => number };

function createCounter(initial: number): Counter {
    let count: number = initial;

    return {
        increment: function(): number {
            count = count + 1;
            return count;
        },
        decrement: function(): number {
            count = count - 1;
            return count;
        },
        get: function(): number {
            return count;
        }
    };
}

let counter: Counter = createCounter(10);
println(counter.increment());  // 11
println(counter.increment());  // 12
println(counter.decrement());  // 11
println(counter.get());        // 11
```

---

## 17. What's NOT in v1

Explicitly excluded for simplicity (candidates for v2+):

- Type inference (`let x = 10`)
- Arrow functions (`=>` in expressions, only in types)
- Compound assignment (`+=`, `-=`, etc.)
- Increment/decrement (`++`, `--`)
- Ternary operator (`? :`)
- Switch statement
- Do-while loop
- String interpolation (template literals)
- Modules/imports
- Interfaces (use type aliases with object types)
- Access modifiers (public/private/protected)
- Static members
- Getters/setters
- Optional properties
- Spread operator
- Destructuring
- Generics
- Union types beyond `T | null`
- undefined (only null)

---

## 18. Grammar (EBNF)

```ebnf
program        = { declaration } ;

declaration    = typeAlias | classDecl | funcDecl | varDecl | statement ;

typeAlias      = "type" IDENTIFIER "=" type ";" ;

classDecl      = "class" IDENTIFIER [ "extends" IDENTIFIER ] "{" { classMember } "}" ;

classMember    = fieldDecl | constructorDecl | methodDecl ;

fieldDecl      = IDENTIFIER ":" type ";" ;

constructorDecl = "constructor" "(" [ paramList ] ")" block ;

methodDecl     = IDENTIFIER "(" [ paramList ] ")" ":" type block ;

varDecl        = ( "let" | "const" ) IDENTIFIER ":" type "=" expression ";" ;

funcDecl       = "function" IDENTIFIER "(" [ paramList ] ")" ":" type block ;

paramList      = param { "," param } ;
param          = IDENTIFIER ":" type ;

type           = primaryType [ "|" "null" ] ;

primaryType    = "number" | "string" | "boolean" | "void" | "null"
               | IDENTIFIER
               | primaryType "[" "]"
               | "{" [ propTypeList ] "}"
               | "(" [ paramTypeList ] ")" "=>" type
               | "(" type ")" ;

propTypeList   = propType { "," propType } ;
propType       = IDENTIFIER ":" type ;

paramTypeList  = type { "," type } ;

statement      = exprStmt | block | ifStmt | whileStmt | forStmt
               | returnStmt | breakStmt | continueStmt ;

exprStmt       = expression ";" ;

block          = "{" { declaration } "}" ;

ifStmt         = "if" "(" expression ")" block [ "else" ( ifStmt | block ) ] ;

whileStmt      = "while" "(" expression ")" block ;

forStmt        = "for" "(" varDecl expression ";" expression ")" block ;

returnStmt     = "return" [ expression ] ";" ;

breakStmt      = "break" ";" ;

continueStmt   = "continue" ";" ;

expression     = assignment ;

assignment     = ( call "." IDENTIFIER | call "[" expression "]" | IDENTIFIER ) "=" assignment
               | logicOr ;

logicOr        = logicAnd { "||" logicAnd } ;

logicAnd       = equality { "&&" equality } ;

equality       = comparison { ( "==" | "!=" ) comparison } ;

comparison     = term { ( "<" | ">" | "<=" | ">=" ) term } ;

term           = factor { ( "+" | "-" ) factor } ;

factor         = unary { ( "*" | "/" | "%" ) unary } ;

unary          = ( "!" | "-" ) unary | call ;

call           = primary { "(" [ argList ] ")" | "." IDENTIFIER | "[" expression "]" } ;

argList        = expression { "," expression } ;

primary        = NUMBER | STRING | "true" | "false" | "null"
               | "this"
               | IDENTIFIER
               | "(" expression ")"
               | arrayLiteral
               | objectLiteral
               | functionExpr
               | "new" IDENTIFIER "(" [ argList ] ")"
               | "super" "(" [ argList ] ")" ;

arrayLiteral   = "[" [ expression { "," expression } ] "]" ;

objectLiteral  = "{" [ propList ] "}" ;

propList       = property { "," property } ;

property       = IDENTIFIER ":" expression ;

functionExpr   = "function" "(" [ paramList ] ")" ":" type block ;
```

---

## 19. Summary

**v1 provides:**
- 4 primitive types + void + null
- Arrays with indexing and mutation
- Object literals with typed properties
- Classes with inheritance and polymorphism
- First-class functions and closures
- Type aliases for complex types
- Nullable types (`T | null`)
- Variables (let/const)
- Arithmetic, comparison, logical operators
- Control flow (if/else, while, for, break, continue)
- Functions with parameters and return values
- Recursion
- Block scoping

This is sufficient to write any practical program while remaining feasible to implement.
