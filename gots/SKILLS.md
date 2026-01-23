# GoTS Development Skills

Common patterns and commands for working with the GoTS codebase.

## Testing Skills

### Run all tests
```bash
cd gots && go test ./...
```

### Run specific package tests with output
```bash
cd gots && go test -v ./pkg/lexer
cd gots && go test -v ./pkg/parser
cd gots && go test -v ./pkg/types
cd gots && go test -v ./pkg/compiler
cd gots && go test -v ./pkg/vm
```

### Run single test function
```bash
cd gots && go test -v ./pkg/parser -run TestParseExpression
cd gots && go test -v ./pkg/compiler -run TestCompileFunction
```

### Test with coverage
```bash
cd gots && go test -cover ./...
cd gots && go test -coverprofile=coverage.out ./pkg/vm
cd gots && go tool cover -html=coverage.out
```

## Building Skills

### Build CLI
```bash
cd gots && go build -o gots ./cmd/gots
```

### Build and run example
```bash
cd gots && go build -o gots ./cmd/gots && ./gots run test/example.gts
```

### Install globally
```bash
cd gots && go install ./cmd/gots
```

## Debugging Skills

### Disassemble bytecode
```bash
cd gots && ./gots disasm test/example.gts
cd gots && ./gots disasm program.gtsb
```

### REPL for interactive testing
```bash
cd gots && ./gots repl
```

### Print debug info during compilation
Add to compiler code:
```go
fmt.Printf("Compiling: %T\n", node)
fmt.Printf("Locals: %v\n", c.locals)
fmt.Printf("Emitting: %s\n", bytecode.OpCode(op))
```

### Print debug info during VM execution
Add to vm code:
```go
fmt.Printf("IP: %d, OP: %s, Stack: %v\n", vm.ip, bytecode.OpCode(op), vm.stack[:vm.sp])
```

## Code Modification Skills

### Adding a new token type
1. Add to `pkg/token/token.go` constant definitions
2. Add to `tokens` map for string representation
3. Update lexer in `pkg/lexer/lexer.go` to recognize it

### Adding a new AST node
1. Define struct in `pkg/ast/ast.go`
2. Implement `Node`, `Statement` or `Expression` interface
3. Add to parser in `pkg/parser/parser.go`
4. Update type checker if needed in `pkg/types/checker.go`
5. Add compiler support in `pkg/compiler/compiler.go`

### Adding a new opcode
1. Define in `pkg/bytecode/bytecode.go`
2. Add to `instructionNames` map
3. Emit in compiler `pkg/compiler/compiler.go`
4. Handle in VM switch statement `pkg/vm/vm.go`

### Adding a new built-in function
1. Define constant in `pkg/vm/builtin.go` (if exists) or `pkg/vm/vm.go`
2. Add to `builtins` map initialization
3. Implement the function
4. Test in REPL or example

## Common Code Patterns

### Parser: Adding prefix parser
```go
p.prefixParseFns[token.SOMETOKEN] = p.parseSomething
```

### Parser: Adding infix parser
```go
p.infixParseFns[token.SOMEOP] = p.parseInfixExpression
```

### Compiler: Emit instruction
```go
c.emit(bytecode.OP_SOMETHING, operand)
```

### Compiler: Emit with jump placeholder
```go
jump := c.emitJump(bytecode.OP_JUMP_IF_FALSE)
// ... emit body ...
c.patchJump(jump)
```

### VM: Push/Pop stack
```go
vm.push(value)
value := vm.pop()
value := vm.peek(0) // peek without removing
```

### VM: Call frame management
```go
frame := &CallFrame{
    fn: closure,
    ip: 0,
    bp: vm.sp - argCount - 1,
}
vm.pushFrame(frame)
```

## Testing Patterns

### Lexer test pattern
```go
tests := []struct {
    expectedType token.Type
    expectedLiteral string
}{
    {token.LET, "let"},
    {token.IDENT, "x"},
}
```

### Parser test pattern
```go
input := `let x: number = 5;`
l := lexer.New(input)
p := parser.New(l)
program := p.ParseProgram()
checkParserErrors(t, p)
```

### Compiler test pattern
```go
testCompile(t, input, expectedInstructions, expectedConstants)
```

### VM test pattern
```go
testVM(t, input, expectedStackTop)
```

## Performance Investigation

### Benchmark a package
```bash
cd gots && go test -bench=. ./pkg/vm
cd gots && go test -bench=. -benchmem ./pkg/compiler
```

### Profile CPU usage
```bash
cd gots && go test -cpuprofile=cpu.prof ./pkg/vm
cd gots && go tool pprof cpu.prof
```

### Profile memory
```bash
cd gots && go test -memprofile=mem.prof ./pkg/vm
cd gots && go tool pprof mem.prof
```

## Git Workflow

### Check current changes
```bash
git status
git diff
```

### Stage and commit
```bash
git add gots/pkg/...
git commit -m "feat: add support for X"
```

### Common commit prefixes
- `feat:` - new feature
- `fix:` - bug fix
- `refactor:` - code restructuring
- `test:` - add/update tests
- `docs:` - documentation
- `perf:` - performance improvement

## Quick Reference

### File locations
- Tokens: `pkg/token/token.go`
- Lexer: `pkg/lexer/lexer.go`
- AST nodes: `pkg/ast/ast.go`
- Parser: `pkg/parser/parser.go`
- Type checker: `pkg/types/checker.go`
- Opcodes: `pkg/bytecode/bytecode.go`
- Compiler: `pkg/compiler/compiler.go`
- VM: `pkg/vm/vm.go`
- Values: `pkg/vm/value.go`
- CLI: `cmd/gots/main.go`

### Important constants
- Max locals: 256
- Max call frames: 64
- Max stack size: 256
- Initial GC threshold: 1MB
- GC growth factor: 2x

### Debug flags (add as needed)
```go
const DEBUG_TRACE_EXECUTION = true
const DEBUG_PRINT_CODE = true
const DEBUG_STRESS_GC = true
```
