// Package bytecode defines the bytecode structures for the GoTS VM.
package bytecode

// Chunk holds bytecode and associated data.
type Chunk struct {
	Code      []byte  // Bytecode instructions
	Constants []any   // Constant pool (will be replaced with Value later)
	Lines     []int   // Line numbers for debugging
}

// NewChunk creates a new empty chunk.
func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0),
		Constants: make([]any, 0),
		Lines:     make([]int, 0),
	}
}

// Write writes a byte to the chunk.
func (c *Chunk) Write(b byte, line int) {
	c.Code = append(c.Code, b)
	c.Lines = append(c.Lines, line)
}

// WriteU16 writes a 16-bit value to the chunk (big-endian).
func (c *Chunk) WriteU16(v uint16, line int) {
	c.Write(byte(v>>8), line)
	c.Write(byte(v), line)
}

// AddConstant adds a constant to the pool and returns its index.
func (c *Chunk) AddConstant(value any) int {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

// Count returns the number of bytes in the chunk.
func (c *Chunk) Count() int {
	return len(c.Code)
}
