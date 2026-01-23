// Package bytecode implements binary format read/write for compiled bytecode.
package bytecode

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Binary format constants
const (
	MagicNumber    uint32 = 0x47545342 // "GTSB" in hex
	FormatVersion  uint16 = 1
)

// Constant type tags
const (
	CONST_NUMBER   byte = 0x01
	CONST_STRING   byte = 0x02
	CONST_BOOL     byte = 0x03
	CONST_NULL     byte = 0x04
	CONST_FUNCTION byte = 0x05
)

// BinaryFunction represents a function in binary format
type BinaryFunction struct {
	Name         string
	Arity        int
	UpvalueCount int
	Chunk        *Chunk
}

// WriteBinary writes a chunk to a binary format.
func WriteBinary(w io.Writer, chunk *Chunk) error {
	// Write magic number
	if err := binary.Write(w, binary.BigEndian, MagicNumber); err != nil {
		return fmt.Errorf("write magic: %w", err)
	}

	// Write version
	if err := binary.Write(w, binary.BigEndian, FormatVersion); err != nil {
		return fmt.Errorf("write version: %w", err)
	}

	// Write the chunk
	return writeChunk(w, chunk)
}

func writeChunk(w io.Writer, chunk *Chunk) error {
	// Write code length and code
	codeLen := uint32(len(chunk.Code))
	if err := binary.Write(w, binary.BigEndian, codeLen); err != nil {
		return fmt.Errorf("write code length: %w", err)
	}
	if _, err := w.Write(chunk.Code); err != nil {
		return fmt.Errorf("write code: %w", err)
	}

	// Write line info (compressed as run-length encoding)
	if err := writeLineInfo(w, chunk.Lines); err != nil {
		return err
	}

	// Write constants
	constLen := uint32(len(chunk.Constants))
	if err := binary.Write(w, binary.BigEndian, constLen); err != nil {
		return fmt.Errorf("write constants length: %w", err)
	}
	for i, c := range chunk.Constants {
		if err := writeConstant(w, c); err != nil {
			return fmt.Errorf("write constant %d: %w", i, err)
		}
	}

	return nil
}

func writeLineInfo(w io.Writer, lines []int) error {
	if len(lines) == 0 {
		return binary.Write(w, binary.BigEndian, uint32(0))
	}

	// Run-length encode lines
	var encoded []struct {
		line  int32
		count uint32
	}

	currentLine := int32(lines[0])
	count := uint32(1)

	for i := 1; i < len(lines); i++ {
		if int32(lines[i]) == currentLine {
			count++
		} else {
			encoded = append(encoded, struct {
				line  int32
				count uint32
			}{currentLine, count})
			currentLine = int32(lines[i])
			count = 1
		}
	}
	encoded = append(encoded, struct {
		line  int32
		count uint32
	}{currentLine, count})

	// Write number of run-length entries
	if err := binary.Write(w, binary.BigEndian, uint32(len(encoded))); err != nil {
		return fmt.Errorf("write line info length: %w", err)
	}

	// Write each entry
	for _, e := range encoded {
		if err := binary.Write(w, binary.BigEndian, e.line); err != nil {
			return fmt.Errorf("write line: %w", err)
		}
		if err := binary.Write(w, binary.BigEndian, e.count); err != nil {
			return fmt.Errorf("write count: %w", err)
		}
	}

	return nil
}

func writeConstant(w io.Writer, c any) error {
	switch v := c.(type) {
	case float64:
		if _, err := w.Write([]byte{CONST_NUMBER}); err != nil {
			return err
		}
		return binary.Write(w, binary.BigEndian, math.Float64bits(v))

	case int:
		if _, err := w.Write([]byte{CONST_NUMBER}); err != nil {
			return err
		}
		return binary.Write(w, binary.BigEndian, math.Float64bits(float64(v)))

	case string:
		if _, err := w.Write([]byte{CONST_STRING}); err != nil {
			return err
		}
		strLen := uint32(len(v))
		if err := binary.Write(w, binary.BigEndian, strLen); err != nil {
			return err
		}
		_, err := w.Write([]byte(v))
		return err

	case bool:
		if _, err := w.Write([]byte{CONST_BOOL}); err != nil {
			return err
		}
		b := byte(0)
		if v {
			b = 1
		}
		_, err := w.Write([]byte{b})
		return err

	case nil:
		_, err := w.Write([]byte{CONST_NULL})
		return err

	case *BinaryFunction:
		if _, err := w.Write([]byte{CONST_FUNCTION}); err != nil {
			return err
		}
		return writeFunction(w, v)

	default:
		// Try to handle function types via interface
		if fn, ok := c.(interface {
			GetName() string
			GetArity() int
			GetUpvalueCount() int
			GetChunk() *Chunk
		}); ok {
			if _, err := w.Write([]byte{CONST_FUNCTION}); err != nil {
				return err
			}
			bfn := &BinaryFunction{
				Name:         fn.GetName(),
				Arity:        fn.GetArity(),
				UpvalueCount: fn.GetUpvalueCount(),
				Chunk:        fn.GetChunk(),
			}
			return writeFunction(w, bfn)
		}

		// Try using reflection to access fields
		return writeConstantReflect(w, c)
	}
}

func writeFunction(w io.Writer, fn *BinaryFunction) error {
	// Write name
	nameLen := uint32(len(fn.Name))
	if err := binary.Write(w, binary.BigEndian, nameLen); err != nil {
		return err
	}
	if _, err := w.Write([]byte(fn.Name)); err != nil {
		return err
	}

	// Write arity
	if err := binary.Write(w, binary.BigEndian, int32(fn.Arity)); err != nil {
		return err
	}

	// Write upvalue count
	if err := binary.Write(w, binary.BigEndian, int32(fn.UpvalueCount)); err != nil {
		return err
	}

	// Write chunk
	return writeChunk(w, fn.Chunk)
}

// writeConstantReflect uses reflection to write function types from different packages
func writeConstantReflect(w io.Writer, c any) error {
	v := reflect.ValueOf(c)

	// Handle pointer to struct
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("unsupported constant type: %T", c)
	}

	// Check if this looks like a function type (has Name, Arity, UpvalueCount, Chunk fields)
	nameField := v.FieldByName("Name")
	arityField := v.FieldByName("Arity")
	upvalueCountField := v.FieldByName("UpvalueCount")
	chunkField := v.FieldByName("Chunk")

	if !nameField.IsValid() || !arityField.IsValid() || !upvalueCountField.IsValid() || !chunkField.IsValid() {
		return fmt.Errorf("unsupported constant type: %T", c)
	}

	// Extract the chunk
	chunkVal := chunkField.Interface()
	chunk, ok := chunkVal.(*Chunk)
	if !ok {
		return fmt.Errorf("unsupported chunk type in function: %T", chunkVal)
	}

	fn := &BinaryFunction{
		Name:         nameField.String(),
		Arity:        int(arityField.Int()),
		UpvalueCount: int(upvalueCountField.Int()),
		Chunk:        chunk,
	}

	if _, err := w.Write([]byte{CONST_FUNCTION}); err != nil {
		return err
	}
	return writeFunction(w, fn)
}

// ReadBinary reads a chunk from binary format.
func ReadBinary(r io.Reader) (*Chunk, error) {
	// Read magic number
	var magic uint32
	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	if magic != MagicNumber {
		return nil, fmt.Errorf("invalid magic number: expected %X, got %X", MagicNumber, magic)
	}

	// Read version
	var version uint16
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return nil, fmt.Errorf("read version: %w", err)
	}
	if version != FormatVersion {
		return nil, fmt.Errorf("unsupported format version: %d (expected %d)", version, FormatVersion)
	}

	// Read the chunk
	return readChunk(r)
}

func readChunk(r io.Reader) (*Chunk, error) {
	chunk := NewChunk()

	// Read code length and code
	var codeLen uint32
	if err := binary.Read(r, binary.BigEndian, &codeLen); err != nil {
		return nil, fmt.Errorf("read code length: %w", err)
	}
	chunk.Code = make([]byte, codeLen)
	if _, err := io.ReadFull(r, chunk.Code); err != nil {
		return nil, fmt.Errorf("read code: %w", err)
	}

	// Read line info
	lines, err := readLineInfo(r)
	if err != nil {
		return nil, err
	}
	chunk.Lines = lines

	// Read constants
	var constLen uint32
	if err := binary.Read(r, binary.BigEndian, &constLen); err != nil {
		return nil, fmt.Errorf("read constants length: %w", err)
	}
	chunk.Constants = make([]any, constLen)
	for i := uint32(0); i < constLen; i++ {
		c, err := readConstant(r)
		if err != nil {
			return nil, fmt.Errorf("read constant %d: %w", i, err)
		}
		chunk.Constants[i] = c
	}

	return chunk, nil
}

func readLineInfo(r io.Reader) ([]int, error) {
	var numEntries uint32
	if err := binary.Read(r, binary.BigEndian, &numEntries); err != nil {
		return nil, fmt.Errorf("read line info length: %w", err)
	}

	if numEntries == 0 {
		return []int{}, nil
	}

	var lines []int
	for i := uint32(0); i < numEntries; i++ {
		var line int32
		var count uint32
		if err := binary.Read(r, binary.BigEndian, &line); err != nil {
			return nil, fmt.Errorf("read line: %w", err)
		}
		if err := binary.Read(r, binary.BigEndian, &count); err != nil {
			return nil, fmt.Errorf("read count: %w", err)
		}
		for j := uint32(0); j < count; j++ {
			lines = append(lines, int(line))
		}
	}

	return lines, nil
}

func readConstant(r io.Reader) (any, error) {
	tag := make([]byte, 1)
	if _, err := io.ReadFull(r, tag); err != nil {
		return nil, fmt.Errorf("read constant tag: %w", err)
	}

	switch tag[0] {
	case CONST_NUMBER:
		var bits uint64
		if err := binary.Read(r, binary.BigEndian, &bits); err != nil {
			return nil, err
		}
		return math.Float64frombits(bits), nil

	case CONST_STRING:
		var strLen uint32
		if err := binary.Read(r, binary.BigEndian, &strLen); err != nil {
			return nil, err
		}
		str := make([]byte, strLen)
		if _, err := io.ReadFull(r, str); err != nil {
			return nil, err
		}
		return string(str), nil

	case CONST_BOOL:
		b := make([]byte, 1)
		if _, err := io.ReadFull(r, b); err != nil {
			return nil, err
		}
		return b[0] == 1, nil

	case CONST_NULL:
		return nil, nil

	case CONST_FUNCTION:
		return readFunction(r)

	default:
		return nil, fmt.Errorf("unknown constant tag: %d", tag[0])
	}
}

func readFunction(r io.Reader) (*BinaryFunction, error) {
	fn := &BinaryFunction{}

	// Read name
	var nameLen uint32
	if err := binary.Read(r, binary.BigEndian, &nameLen); err != nil {
		return nil, err
	}
	name := make([]byte, nameLen)
	if _, err := io.ReadFull(r, name); err != nil {
		return nil, err
	}
	fn.Name = string(name)

	// Read arity
	var arity int32
	if err := binary.Read(r, binary.BigEndian, &arity); err != nil {
		return nil, err
	}
	fn.Arity = int(arity)

	// Read upvalue count
	var upvalueCount int32
	if err := binary.Read(r, binary.BigEndian, &upvalueCount); err != nil {
		return nil, err
	}
	fn.UpvalueCount = int(upvalueCount)

	// Read chunk
	chunk, err := readChunk(r)
	if err != nil {
		return nil, err
	}
	fn.Chunk = chunk

	return fn, nil
}
