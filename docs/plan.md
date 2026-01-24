# goTS Implementation Plan (Next)

## Current Status
- End-to-end pipeline exists: lexer -> parser -> typed AST builder -> Go codegen -> CLI run/build/emit-go/repl.
- Core language features implemented: literals, variables, functions/closures, classes/inheritance, arrays/objects, control flow, switch, for-of, builtins.
- Test coverage: lexer/parser/type system unit tests plus example programs in `gots/test`.

## Next Implementation Plan

### 1) Correctness and Semantics Gaps
- Optional chaining: type-check nullable/optional cases and add codegen for `?.`/`?.[]`/`?.()` with nil checks.
- `super()` handling: generate parent constructor call and initialize embedded fields (see Decisions for codegen strategy).
- Nullish coalescing: ensure codegen preserves static type (avoid interface-only result) and add tests.
- Const and assignment validation: enforce `const` immutability, validate assignment targets, and ensure ++/-- and compound assigns emit Go-legal code.
- For-of on strings: ensure element type is string and codegen performs proper conversion.
- Error handling: implement `try`/`catch`/`throw` with codegen mapping to Go's `panic`/`recover`.

### 2) Type Checking Parity and Determinism
- Unify `typed.Builder` and `types.Checker` into a single source of truth, consolidating missing checks (break/continue, switch case compatibility, builtin arity).
- Stabilize object type codegen (deterministic field ordering or map representation) to avoid Go type mismatches.
- Improve builtin typing: `pop` returns element type, `push` returns length, and add stricter type/arity checks.

### 3) Tests and Documentation Alignment
- Add tests for optional chaining, nullish coalescing, super, const reassignment, for-of strings, try/catch/throw, and error cases.
- Update `docs/language-spec-v1.md` and `docs/starter-guide.md` to reflect the Go transpiler and actual runtime semantics.
- Add example snippets for new/changed features and document any known limitations.

## Decisions to Make
- Nullable representation (pointer vs interface) and how it affects optional chaining and object literals.
- Object literal representation (anonymous struct vs map) for predictable typing and interop.
- Whether to allow update/compound assignments as expressions or restrict them to statement contexts.
- `super()` codegen strategy: call parent's `NewParent()` constructor when available, fall back to inline field initialization when parent has no explicit constructor.
