# Testing

goTS has four layers of tests: unit tests per package, integration tests in the repo root, a conformance suite under `test/test262/`, and runnable example programs under `test/`.

## Quick Commands

```bash
make test             # All unit + integration tests (verbose)
make test-unit        # pkg/... unit tests only
make test-integration # Repo-root integration test
make test-test262     # Conformance suite
make test-cover       # Coverage summary
make test-cover-html  # HTML coverage report
make check            # fmt + vet + test
```

## Unit Tests

Each package has `*_test.go` files next to the implementation:

```bash
go test -v ./pkg/lexer
go test -v ./pkg/parser
go test -v ./pkg/types
go test -v ./pkg/codegen
```

The codegen tests are golden-style: they assert on emitted Go snippets (e.g. `pkg/codegen/generics_test.go`, `pkg/codegen/tuple_test.go`, `pkg/codegen/sql_test.go`).

Run a single test:

```bash
go test -v ./pkg/codegen -run TestEnumCodegen
```

## Integration Tests

`integration_test.go` (repo root) compiles and runs every program in `test/*.gts`, asserting on the output. Individual programs can be targeted:

```bash
make test-file FILE=example
go test -v -run "TestIntegration/example"
```

## Conformance Suite

`test/test262/` contains a curated subset of ECMAScript test262 cases organized by feature area:

```
test/test262/
├── builtins/     # built-in object behavior
├── classes/      # classes & inheritance
├── closures/     # closures & captures
├── functions/    # functions & arrows
├── identifiers/  # naming rules
├── literals/     # literal forms
├── operators/    # operator semantics
├── regexp/       # regex behavior
├── statements/   # control flow
├── types/        # type rules
└── variables/    # declaration rules
```

```bash
make test-test262
# prints PASS/FAIL per file and a summary
```

## Example Programs

`test/` doubles as living documentation — each program exercises one feature:

| File | Demonstrates |
|------|--------------|
| `example.gts` | Core language: functions, arrays, classes, builtins |
| `generics.gts` | Generic functions & classes, inference, defaults |
| `interfaces.gts` | Interfaces, structural typing |
| `higher_order.gts` | Higher-order functions, closures |
| `y_combinator.gts` | The Y combinator (closures) |
| `church_encoding.gts` | Church-encoded numerals |
| `inheritance_super.gts` | Class inheritance, `super` |
| `eventloop.gts` | async/await, timers, microtasks |
| `promise.gts` | Promise chaining |
| `sql_test.gts`, `sql_tx_test.gts` | SQL + transactions |
| `go_imports.gts` | Go stdlib imports |
| `maps.gts` | `Map<K, V>` |
| `optional_chaining.gts` | `?.` |
| `nullish_coalescing.gts` | `??` |
| `try_catch.gts` | try/catch/finally |
| `modules/` | Local module imports |
| `methods_demo.gts` | String & array methods |
| `new_features.gts` | Array methods, `Set`, `Map` |

Run one:

```bash
gots run test/generics.gts
```

## Coverage

```bash
make test-cover-html
# opens coverage.html
```

## CI Checklist

```bash
make check          # fmt, vet, tests
make test-test262   # conformance
go build ./cmd/gots # compiler builds
```
