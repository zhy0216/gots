package bytecode

import "testing"

func TestOpcodeString(t *testing.T) {
	tests := []struct {
		op   OpCode
		want string
	}{
		{OP_CONSTANT, "OP_CONSTANT"},
		{OP_NULL, "OP_NULL"},
		{OP_TRUE, "OP_TRUE"},
		{OP_FALSE, "OP_FALSE"},
		{OP_ADD, "OP_ADD"},
		{OP_SUBTRACT, "OP_SUBTRACT"},
		{OP_MULTIPLY, "OP_MULTIPLY"},
		{OP_DIVIDE, "OP_DIVIDE"},
		{OP_MODULO, "OP_MODULO"},
		{OP_NEGATE, "OP_NEGATE"},
		{OP_NOT, "OP_NOT"},
		{OP_EQUAL, "OP_EQUAL"},
		{OP_NOT_EQUAL, "OP_NOT_EQUAL"},
		{OP_LESS, "OP_LESS"},
		{OP_LESS_EQUAL, "OP_LESS_EQUAL"},
		{OP_GREATER, "OP_GREATER"},
		{OP_GREATER_EQUAL, "OP_GREATER_EQUAL"},
		{OP_PRINT, "OP_PRINT"},
		{OP_PRINTLN, "OP_PRINTLN"},
		{OP_POP, "OP_POP"},
		{OP_RETURN, "OP_RETURN"},
	}

	for _, tt := range tests {
		got := tt.op.String()
		if got != tt.want {
			t.Errorf("OpCode(%d).String() = %q, want %q", tt.op, got, tt.want)
		}
	}
}

func TestChunk(t *testing.T) {
	chunk := NewChunk()

	// Test Write
	chunk.Write(byte(OP_CONSTANT), 1)
	if chunk.Count() != 1 {
		t.Errorf("chunk.Count() = %d, want 1", chunk.Count())
	}

	// Test WriteU16
	chunk.WriteU16(256, 1)
	if chunk.Count() != 3 {
		t.Errorf("chunk.Count() = %d, want 3", chunk.Count())
	}

	// Verify big-endian encoding
	if chunk.Code[1] != 1 || chunk.Code[2] != 0 {
		t.Errorf("WriteU16(256) = [%d, %d], want [1, 0]", chunk.Code[1], chunk.Code[2])
	}

	// Test AddConstant
	idx := chunk.AddConstant(42.0)
	if idx != 0 {
		t.Errorf("AddConstant returned %d, want 0", idx)
	}

	idx = chunk.AddConstant("hello")
	if idx != 1 {
		t.Errorf("AddConstant returned %d, want 1", idx)
	}

	if len(chunk.Constants) != 2 {
		t.Errorf("len(chunk.Constants) = %d, want 2", len(chunk.Constants))
	}

	// Test line numbers
	if len(chunk.Lines) != 3 {
		t.Errorf("len(chunk.Lines) = %d, want 3", len(chunk.Lines))
	}
	for i, line := range chunk.Lines {
		if line != 1 {
			t.Errorf("chunk.Lines[%d] = %d, want 1", i, line)
		}
	}
}

func TestChunkReadU16(t *testing.T) {
	chunk := NewChunk()
	chunk.WriteU16(0x1234, 1)

	val := ReadU16(chunk.Code, 0)
	if val != 0x1234 {
		t.Errorf("ReadU16() = %x, want %x", val, 0x1234)
	}
}

func TestDisassemble(t *testing.T) {
	chunk := NewChunk()

	// OP_CONSTANT 0 (42)
	idx := chunk.AddConstant(42.0)
	chunk.Write(byte(OP_CONSTANT), 1)
	chunk.WriteU16(uint16(idx), 1)

	// OP_CONSTANT 1 (3.14)
	idx = chunk.AddConstant(3.14)
	chunk.Write(byte(OP_CONSTANT), 1)
	chunk.WriteU16(uint16(idx), 1)

	// OP_ADD
	chunk.Write(byte(OP_ADD), 1)

	// OP_PRINTLN
	chunk.Write(byte(OP_PRINTLN), 1)

	// OP_RETURN
	chunk.Write(byte(OP_RETURN), 2)

	// Just verify it doesn't panic
	output := Disassemble(chunk, "test")
	if output == "" {
		t.Error("Disassemble returned empty string")
	}
}
