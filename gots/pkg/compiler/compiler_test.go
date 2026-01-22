package compiler

import (
	"testing"

	"github.com/pocketlang/gots/pkg/bytecode"
	"github.com/pocketlang/gots/pkg/lexer"
	"github.com/pocketlang/gots/pkg/parser"
)

func TestCompileNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42;", 42},
		{"3.14;", 3.14},
		{"0;", 0},
	}

	for _, tt := range tests {
		chunk := compileSource(t, tt.input)

		// Should have: OP_CONSTANT, u16 index, OP_POP, OP_RETURN
		if len(chunk.Code) < 4 {
			t.Fatalf("expected at least 4 bytes, got %d", len(chunk.Code))
		}

		if bytecode.OpCode(chunk.Code[0]) != bytecode.OP_CONSTANT {
			t.Errorf("expected OP_CONSTANT, got %v", bytecode.OpCode(chunk.Code[0]))
		}

		// Check constant value
		constIdx := bytecode.ReadU16(chunk.Code, 1)
		if int(constIdx) >= len(chunk.Constants) {
			t.Fatalf("constant index %d out of range", constIdx)
		}

		val, ok := chunk.Constants[constIdx].(float64)
		if !ok {
			t.Fatalf("constant is not float64: %T", chunk.Constants[constIdx])
		}
		if val != tt.expected {
			t.Errorf("constant = %v, want %v", val, tt.expected)
		}
	}
}

func TestCompileStringLiteral(t *testing.T) {
	input := `"hello";`
	chunk := compileSource(t, input)

	if bytecode.OpCode(chunk.Code[0]) != bytecode.OP_CONSTANT {
		t.Errorf("expected OP_CONSTANT, got %v", bytecode.OpCode(chunk.Code[0]))
	}

	constIdx := bytecode.ReadU16(chunk.Code, 1)
	val, ok := chunk.Constants[constIdx].(string)
	if !ok {
		t.Fatalf("constant is not string: %T", chunk.Constants[constIdx])
	}
	if val != "hello" {
		t.Errorf("constant = %q, want %q", val, "hello")
	}
}

func TestCompileBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bytecode.OpCode
	}{
		{"true;", bytecode.OP_TRUE},
		{"false;", bytecode.OP_FALSE},
	}

	for _, tt := range tests {
		chunk := compileSource(t, tt.input)

		if bytecode.OpCode(chunk.Code[0]) != tt.expected {
			t.Errorf("input %q: expected %v, got %v", tt.input, tt.expected, bytecode.OpCode(chunk.Code[0]))
		}
	}
}

func TestCompileNullLiteral(t *testing.T) {
	input := "null;"
	chunk := compileSource(t, input)

	if bytecode.OpCode(chunk.Code[0]) != bytecode.OP_NULL {
		t.Errorf("expected OP_NULL, got %v", bytecode.OpCode(chunk.Code[0]))
	}
}

func TestCompileBinaryArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected bytecode.OpCode
	}{
		{"1 + 2;", bytecode.OP_ADD},
		{"3 - 4;", bytecode.OP_SUBTRACT},
		{"5 * 6;", bytecode.OP_MULTIPLY},
		{"8 / 2;", bytecode.OP_DIVIDE},
		{"10 % 3;", bytecode.OP_MODULO},
	}

	for _, tt := range tests {
		chunk := compileSource(t, tt.input)

		// Should have: OP_CONSTANT, u16, OP_CONSTANT, u16, OP_<op>, OP_POP, OP_RETURN
		found := false
		for i := 0; i < len(chunk.Code); i++ {
			if bytecode.OpCode(chunk.Code[i]) == tt.expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("input %q: expected %v in bytecode", tt.input, tt.expected)
		}
	}
}

func TestCompileBinaryComparison(t *testing.T) {
	tests := []struct {
		input    string
		expected bytecode.OpCode
	}{
		{"1 == 2;", bytecode.OP_EQUAL},
		{"1 != 2;", bytecode.OP_NOT_EQUAL},
		{"1 < 2;", bytecode.OP_LESS},
		{"1 <= 2;", bytecode.OP_LESS_EQUAL},
		{"1 > 2;", bytecode.OP_GREATER},
		{"1 >= 2;", bytecode.OP_GREATER_EQUAL},
	}

	for _, tt := range tests {
		chunk := compileSource(t, tt.input)

		found := false
		for i := 0; i < len(chunk.Code); i++ {
			if bytecode.OpCode(chunk.Code[i]) == tt.expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("input %q: expected %v in bytecode", tt.input, tt.expected)
		}
	}
}

func TestCompileUnaryNegate(t *testing.T) {
	input := "-42;"
	chunk := compileSource(t, input)

	// Should have: OP_CONSTANT, u16, OP_NEGATE
	found := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_NEGATE {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected OP_NEGATE in bytecode")
	}
}

func TestCompileUnaryNot(t *testing.T) {
	input := "!true;"
	chunk := compileSource(t, input)

	found := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_NOT {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected OP_NOT in bytecode")
	}
}

func TestCompileStringConcat(t *testing.T) {
	input := `"hello" + "world";`
	chunk := compileSource(t, input)

	// String concatenation uses OP_CONCAT
	found := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_CONCAT {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected OP_CONCAT in bytecode")
	}
}

func TestCompileGroupedExpression(t *testing.T) {
	input := "(1 + 2) * 3;"
	chunk := compileSource(t, input)

	// The bytecode should evaluate (1+2) first, then multiply by 3
	// Order should be: const 1, const 2, add, const 3, multiply
	ops := extractOpcodes(chunk)

	expectedOrder := []bytecode.OpCode{
		bytecode.OP_CONSTANT, // 1
		bytecode.OP_CONSTANT, // 2
		bytecode.OP_ADD,
		bytecode.OP_CONSTANT, // 3
		bytecode.OP_MULTIPLY,
	}

	opIdx := 0
	for _, op := range ops {
		if opIdx < len(expectedOrder) && op == expectedOrder[opIdx] {
			opIdx++
		}
	}

	if opIdx != len(expectedOrder) {
		t.Errorf("expected ops in order %v, got %v", expectedOrder, ops)
	}
}

func TestCompilePrintlnBuiltin(t *testing.T) {
	input := "println(42);"
	chunk := compileSource(t, input)

	// println is compiled as OP_BUILTIN with builtin ID 0
	foundBuiltin := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_BUILTIN {
			foundBuiltin = true
			break
		}
	}
	if !foundBuiltin {
		t.Error("expected OP_BUILTIN in bytecode")
	}
}

func TestCompileExpressionStatement(t *testing.T) {
	input := "42;"
	chunk := compileSource(t, input)

	// Expression statement should pop the result
	found := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_POP {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected OP_POP in bytecode for expression statement")
	}
}

func TestCompileReturn(t *testing.T) {
	input := "42;"
	chunk := compileSource(t, input)

	// Should end with OP_RETURN
	lastOp := bytecode.OpCode(chunk.Code[len(chunk.Code)-1])
	if lastOp != bytecode.OP_RETURN {
		t.Errorf("expected last op to be OP_RETURN, got %v", lastOp)
	}
}

func TestCompileGlobalVariable(t *testing.T) {
	input := `let x: number = 42;`
	chunk := compileSource(t, input)

	// Should have: OP_CONSTANT (42), OP_SET_GLOBAL
	foundSetGlobal := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_SET_GLOBAL {
			foundSetGlobal = true
			break
		}
	}
	if !foundSetGlobal {
		t.Error("expected OP_SET_GLOBAL in bytecode")
	}
}

func TestCompileGlobalVariableAccess(t *testing.T) {
	input := `let x: number = 42; x;`
	chunk := compileSource(t, input)

	// Should have OP_GET_GLOBAL when accessing x
	foundGetGlobal := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_GET_GLOBAL {
			foundGetGlobal = true
			break
		}
	}
	if !foundGetGlobal {
		t.Error("expected OP_GET_GLOBAL in bytecode")
	}
}

func TestCompileGlobalVariableAssignment(t *testing.T) {
	input := `let x: number = 1; x = 2;`
	chunk := compileSource(t, input)

	// Should have two OP_SET_GLOBAL (declaration and assignment)
	setGlobalCount := 0
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_SET_GLOBAL {
			setGlobalCount++
		}
	}
	if setGlobalCount < 2 {
		t.Errorf("expected at least 2 OP_SET_GLOBAL, got %d", setGlobalCount)
	}
}

func TestCompileIfStatement(t *testing.T) {
	input := `if (true) { println(1); }`
	chunk := compileSource(t, input)

	// Should have OP_JUMP_IF_FALSE for the condition
	found := false
	for i := 0; i < len(chunk.Code); i++ {
		if bytecode.OpCode(chunk.Code[i]) == bytecode.OP_JUMP_IF_FALSE {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected OP_JUMP_IF_FALSE in bytecode")
	}
}

func TestCompileIfElseStatement(t *testing.T) {
	input := `if (false) { println(1); } else { println(2); }`
	chunk := compileSource(t, input)

	// Should have OP_JUMP_IF_FALSE and OP_JUMP
	foundJumpIfFalse := false
	foundJump := false
	for i := 0; i < len(chunk.Code); i++ {
		op := bytecode.OpCode(chunk.Code[i])
		if op == bytecode.OP_JUMP_IF_FALSE {
			foundJumpIfFalse = true
		}
		if op == bytecode.OP_JUMP {
			foundJump = true
		}
	}
	if !foundJumpIfFalse {
		t.Error("expected OP_JUMP_IF_FALSE in bytecode")
	}
	if !foundJump {
		t.Error("expected OP_JUMP in bytecode for else branch")
	}
}

func TestCompileWhileStatement(t *testing.T) {
	input := `let x: number = 0; while (x < 3) { x = x + 1; }`
	chunk := compileSource(t, input)

	// Should have OP_JUMP_IF_FALSE and OP_JUMP_BACK
	foundJumpIfFalse := false
	foundJumpBack := false
	for i := 0; i < len(chunk.Code); i++ {
		op := bytecode.OpCode(chunk.Code[i])
		if op == bytecode.OP_JUMP_IF_FALSE {
			foundJumpIfFalse = true
		}
		if op == bytecode.OP_JUMP_BACK {
			foundJumpBack = true
		}
	}
	if !foundJumpIfFalse {
		t.Error("expected OP_JUMP_IF_FALSE in bytecode")
	}
	if !foundJumpBack {
		t.Error("expected OP_JUMP_BACK in bytecode for loop")
	}
}

// Helper functions

func compileSource(t *testing.T, source string) *bytecode.Chunk {
	t.Helper()

	l := lexer.New(source)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser errors: %v", p.Errors())
	}

	compiler := New()
	chunk, err := compiler.Compile(program)
	if err != nil {
		t.Fatalf("compiler error: %v", err)
	}

	return chunk
}

func extractOpcodes(chunk *bytecode.Chunk) []bytecode.OpCode {
	var ops []bytecode.OpCode
	for i := 0; i < len(chunk.Code); {
		op := bytecode.OpCode(chunk.Code[i])
		ops = append(ops, op)
		i += opcodeSize(op)
	}
	return ops
}

func opcodeSize(op bytecode.OpCode) int {
	switch op {
	case bytecode.OP_CONSTANT, bytecode.OP_GET_GLOBAL, bytecode.OP_SET_GLOBAL,
		bytecode.OP_JUMP, bytecode.OP_JUMP_BACK, bytecode.OP_JUMP_IF_FALSE, bytecode.OP_JUMP_IF_TRUE,
		bytecode.OP_CLASS, bytecode.OP_GET_PROPERTY, bytecode.OP_SET_PROPERTY, bytecode.OP_METHOD,
		bytecode.OP_GET_SUPER, bytecode.OP_ARRAY, bytecode.OP_OBJECT:
		return 3 // op + u16
	case bytecode.OP_GET_LOCAL, bytecode.OP_SET_LOCAL, bytecode.OP_GET_UPVALUE, bytecode.OP_SET_UPVALUE,
		bytecode.OP_CALL, bytecode.OP_POPN:
		return 2 // op + u8
	case bytecode.OP_INVOKE, bytecode.OP_SUPER_INVOKE:
		return 4 // op + u16 + u8
	case bytecode.OP_BUILTIN:
		return 3 // op + u8 + u8
	default:
		return 1 // simple op
	}
}
