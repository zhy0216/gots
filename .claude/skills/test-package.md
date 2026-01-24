# Test Single GoTS Package

Run tests for a specific package with verbose output.

## Usage
Test individual packages when working on specific components.

## Commands
```bash
# Lexer tests
cd /Users/yang/workspace/quickts/gots && go test -v ./pkg/lexer

# Parser tests
cd /Users/yang/workspace/quickts/gots && go test -v ./pkg/parser

# Type checker tests
cd /Users/yang/workspace/quickts/gots && go test -v ./pkg/types

# Compiler tests
cd /Users/yang/workspace/quickts/gots && go test -v ./pkg/compiler

# VM tests
cd /Users/yang/workspace/quickts/gots && go test -v ./pkg/vm
```

## When to use
- When modifying a specific package
- To get detailed test output
- To debug failing tests
