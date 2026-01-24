# goTS Development Plan

## Phase 1: Core Language Extensions (COMPLETED)

### Template Literals ✅
- Lexer: Added `BACKTICK` token and template literal scanning
- Parser: Added `TemplateLiteral` AST node with parts and expressions
- Typed: Added `TemplateLit` expression
- Codegen: Generates `fmt.Sprintf()` calls

```typescript
let name = "World"
let msg = `Hello ${name}!`
```

### Destructuring ✅
- Parser: Added `ArrayPattern` and `ObjectPattern` AST nodes
- Typed: Added destructuring pattern handling in variable declarations
- Codegen: Generates individual variable assignments from array/object access

```typescript
let [a, b] = [1, 2]
let {x, y} = {x: 10, y: 20}
```

### Spread Operator ✅
- Lexer: Added `ELLIPSIS` token (`...`)
- Parser: Added `SpreadExpr` AST node
- Typed: Added `SpreadExpr` expression
- Codegen: Generates `append()` for array literals, `args...` for function calls

```typescript
let arr1 = [1, 2, 3]
let arr2 = [...arr1, 4, 5]
```

### Enums ✅
- Token: Added `ENUM` keyword
- AST: Added `EnumDecl` and `EnumMember` nodes
- Types: Added `Enum` type with member lookup
- Typed: Added `EnumMemberExpr` for member access
- Codegen: Generates Go `type` and `const` block

```typescript
enum Color { Red, Green, Blue }
enum Status { Pending = 1, Active = 2 }
let c = Color.Red
```

---

## Phase 2: Advanced Types (COMPLETED)

### Union Types ✅
- Token: `|` (PIPE) already exists
- AST: Added `UnionType` node
- Types: Added `Union` type with helper methods
- Parser: Parse union types with proper precedence
- Builder: Resolve union types with `MakeUnion` helper
- Codegen: Generate `interface{}` for union types
- Tests: Parser, types, and codegen tests

```typescript
type StringOrNumber = string | int
type Value = string | int | boolean
let x: string | int = "hello"
```

### Intersection Types ✅
- Token: Added `AMPERSAND` (&) token
- Lexer: Recognize `&` for intersection types
- AST: Added `IntersectionType` node
- Types: Added `Intersection` type with object merging
- Parser: Parse intersection types with higher precedence than union
- Builder: Resolve intersection types with `MakeIntersection` helper
- Codegen: Merge object properties or use `interface{}`
- Tests: Parser, types, and codegen tests

```typescript
type A = { x: int }
type B = { y: int }
type AB = A & B  // { x: int, y: int }
type Combined = { x: int } & { y: int }
```

### Literal Types ✅
- AST: Added `LiteralType` node
- Types: Added `Literal` type for singleton values
- Parser: Parse string, number, and boolean literals in type positions
- Builder: Resolve literal types
- Codegen: Generate base primitive types
- Tests: Parser tests

```typescript
type Direction = "north" | "south" | "east" | "west"
type One = 1
type Status = "active" | "inactive" | 0 | 1
let x: "hello" = "hello"
let y: 42 = 42
let z: true = true
```

### Tuple Types ✅
- Token: Uses existing `[` and `]` brackets, `...` for rest elements
- AST: Added `TupleType` node with Elements and RestElement
- Types: Added `Tuple` type with element types and optional rest
- Parser: Parse tuple types in type positions with rest element support
- Builder: Resolve tuple types with element type mapping
- Codegen: Generate Go structs with numbered fields (T0, T1, etc.)
- Tests: Parser, types, and codegen tests

```typescript
type Pair = [string, int]
type Triple = [string, int, boolean]
type VarArgs = [string, ...int[]]
let point: [int, int] = [1, 2]
```

---

## Phase 3: Control Flow & Error Handling (TODO)

### Optional Chaining Improvements
- Deep optional chaining: `obj?.a?.b?.c`
- Optional method calls: `obj?.method?.()`

### Nullish Coalescing
- `??` operator for null/undefined fallback
- `??=` assignment operator

### Pattern Matching (stretch goal)
- Match expressions with type guards
- Exhaustiveness checking

---

## Phase 4: Module System Enhancements (TODO)

### Re-exports
- `export { foo } from "./module"`
- `export * from "./module"`

### Default Exports
- `export default class Foo {}`
- `import Foo from "./module"`

### Namespace Imports
- `import * as utils from "./utils"`

---


## Phase 5: Standard Library (TODO)

### String Methods
- `split()`, `join()`, `replace()`, `trim()`
- `startsWith()`, `endsWith()`, `includes()`

### Array Methods
- `map()`, `filter()`, `reduce()`
- `find()`, `findIndex()`, `some()`, `every()`

### Object Utilities
- `Object.keys()`, `Object.values()`, `Object.entries()`
- `Object.assign()`, spread in object literals

### Date/Time
- Basic date handling
- Formatting utilities

---

## Implementation Notes

### Adding a New Feature Checklist
1. **Token** - Add to `pkg/token/token.go` if new syntax
2. **Lexer** - Update `pkg/lexer/lexer.go` to recognize tokens
3. **AST** - Add node types to `pkg/ast/ast.go`
4. **Parser** - Add parsing logic to `pkg/parser/parser.go`
5. **Types** - Add type definitions to `pkg/types/types.go` if needed
6. **Typed AST** - Add typed nodes to `pkg/typed/`
7. **Builder** - Update `pkg/typed/builder.go` for type checking
8. **Codegen** - Update `pkg/codegen/codegen.go` for Go generation
9. **Tests** - Add tests at each layer

### Testing Strategy
- Unit tests for each package
- Integration tests with full compilation
- End-to-end tests running generated Go code
