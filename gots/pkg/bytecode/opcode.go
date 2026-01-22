// Package bytecode defines the bytecode structures and opcodes for the GoTS VM.
package bytecode

import (
	"fmt"
	"strings"
)

// OpCode represents a VM instruction opcode.
type OpCode byte

const (
	// Constants & Literals
	OP_CONSTANT      OpCode = 0x01 // u16 index -> push constants[index]
	OP_NULL          OpCode = 0x02 // push null
	OP_TRUE          OpCode = 0x03 // push true
	OP_FALSE         OpCode = 0x04 // push false

	// Arithmetic
	OP_ADD           OpCode = 0x10 // pop b, pop a, push a + b
	OP_SUBTRACT      OpCode = 0x11 // pop b, pop a, push a - b
	OP_MULTIPLY      OpCode = 0x12 // pop b, pop a, push a * b
	OP_DIVIDE        OpCode = 0x13 // pop b, pop a, push a / b
	OP_MODULO        OpCode = 0x14 // pop b, pop a, push a % b
	OP_NEGATE        OpCode = 0x15 // pop a, push -a

	// Comparison
	OP_EQUAL         OpCode = 0x20 // pop b, pop a, push a == b
	OP_NOT_EQUAL     OpCode = 0x21 // pop b, pop a, push a != b
	OP_LESS          OpCode = 0x22 // pop b, pop a, push a < b
	OP_LESS_EQUAL    OpCode = 0x23 // pop b, pop a, push a <= b
	OP_GREATER       OpCode = 0x24 // pop b, pop a, push a > b
	OP_GREATER_EQUAL OpCode = 0x25 // pop b, pop a, push a >= b

	// Logical
	OP_NOT           OpCode = 0x30 // pop a, push !a

	// String
	OP_CONCAT        OpCode = 0x40 // pop b, pop a, push a + b (strings)

	// Variables
	OP_GET_LOCAL     OpCode = 0x50 // u8 slot -> push stack[frame.base + slot]
	OP_SET_LOCAL     OpCode = 0x51 // u8 slot -> stack[frame.base + slot] = peek()
	OP_GET_GLOBAL    OpCode = 0x52 // u16 index -> push globals[constants[index]]
	OP_SET_GLOBAL    OpCode = 0x53 // u16 index -> globals[constants[index]] = peek()
	OP_GET_UPVALUE   OpCode = 0x54 // u8 index -> push closure.upvalues[index]
	OP_SET_UPVALUE   OpCode = 0x55 // u8 index -> closure.upvalues[index] = peek()

	// Stack Operations
	OP_POP           OpCode = 0x60 // discard top of stack
	OP_POPN          OpCode = 0x61 // u8 n -> discard n values from stack
	OP_DUP           OpCode = 0x62 // duplicate top of stack

	// Control Flow
	OP_JUMP          OpCode = 0x70 // u16 offset -> ip += offset
	OP_JUMP_BACK     OpCode = 0x71 // u16 offset -> ip -= offset
	OP_JUMP_IF_FALSE OpCode = 0x72 // u16 offset -> if !pop() then ip += offset
	OP_JUMP_IF_TRUE  OpCode = 0x73 // u16 offset -> if pop() then ip += offset

	// Functions & Calls
	OP_CALL          OpCode = 0x80 // u8 argCount -> call function with args
	OP_RETURN        OpCode = 0x81 // return from function
	OP_CLOSURE       OpCode = 0x82 // u16 funcIndex, [u8 isLocal, u8 index]* -> create closure

	// Classes & Objects
	OP_CLASS         OpCode = 0x90 // u16 classIndex -> push class
	OP_GET_PROPERTY  OpCode = 0x91 // u16 nameIndex -> pop obj, push obj.property
	OP_SET_PROPERTY  OpCode = 0x92 // u16 nameIndex -> pop val, pop obj, obj.property = val, push val
	OP_METHOD        OpCode = 0x93 // u16 nameIndex -> pop closure, add method to class at stack top
	OP_INVOKE        OpCode = 0x94 // u16 nameIndex, u8 argCount -> invoke method directly
	OP_INHERIT       OpCode = 0x95 // pop super, pop sub, sub inherits from super
	OP_GET_SUPER     OpCode = 0x96 // u16 nameIndex -> lookup method in superclass
	OP_SUPER_INVOKE  OpCode = 0x97 // u16 nameIndex, u8 argCount -> invoke super method

	// Arrays
	OP_ARRAY         OpCode = 0xA0 // u16 count -> pop count values, push array
	OP_GET_INDEX     OpCode = 0xA1 // pop index, pop array, push array[index]
	OP_SET_INDEX     OpCode = 0xA2 // pop val, pop index, pop array, array[index] = val, push val

	// Objects (literals)
	OP_OBJECT        OpCode = 0xB0 // u16 count -> pop count key-value pairs, push object

	// Special
	OP_CLOSE_UPVALUE OpCode = 0xC0 // close upvalue at stack top
	OP_PRINT         OpCode = 0xD0 // pop value, print it (built-in)
	OP_PRINTLN       OpCode = 0xD1 // pop value, print it with newline

	// Built-in Functions
	OP_BUILTIN       OpCode = 0xE0 // u8 builtinId, u8 argCount -> call built-in
)

// String returns the name of the opcode.
func (op OpCode) String() string {
	if name, ok := opNames[op]; ok {
		return name
	}
	return fmt.Sprintf("OP_UNKNOWN(%d)", op)
}

var opNames = map[OpCode]string{
	OP_CONSTANT:      "OP_CONSTANT",
	OP_NULL:          "OP_NULL",
	OP_TRUE:          "OP_TRUE",
	OP_FALSE:         "OP_FALSE",
	OP_ADD:           "OP_ADD",
	OP_SUBTRACT:      "OP_SUBTRACT",
	OP_MULTIPLY:      "OP_MULTIPLY",
	OP_DIVIDE:        "OP_DIVIDE",
	OP_MODULO:        "OP_MODULO",
	OP_NEGATE:        "OP_NEGATE",
	OP_EQUAL:         "OP_EQUAL",
	OP_NOT_EQUAL:     "OP_NOT_EQUAL",
	OP_LESS:          "OP_LESS",
	OP_LESS_EQUAL:    "OP_LESS_EQUAL",
	OP_GREATER:       "OP_GREATER",
	OP_GREATER_EQUAL: "OP_GREATER_EQUAL",
	OP_NOT:           "OP_NOT",
	OP_CONCAT:        "OP_CONCAT",
	OP_GET_LOCAL:     "OP_GET_LOCAL",
	OP_SET_LOCAL:     "OP_SET_LOCAL",
	OP_GET_GLOBAL:    "OP_GET_GLOBAL",
	OP_SET_GLOBAL:    "OP_SET_GLOBAL",
	OP_GET_UPVALUE:   "OP_GET_UPVALUE",
	OP_SET_UPVALUE:   "OP_SET_UPVALUE",
	OP_POP:           "OP_POP",
	OP_POPN:          "OP_POPN",
	OP_DUP:           "OP_DUP",
	OP_JUMP:          "OP_JUMP",
	OP_JUMP_BACK:     "OP_JUMP_BACK",
	OP_JUMP_IF_FALSE: "OP_JUMP_IF_FALSE",
	OP_JUMP_IF_TRUE:  "OP_JUMP_IF_TRUE",
	OP_CALL:          "OP_CALL",
	OP_RETURN:        "OP_RETURN",
	OP_CLOSURE:       "OP_CLOSURE",
	OP_CLASS:         "OP_CLASS",
	OP_GET_PROPERTY:  "OP_GET_PROPERTY",
	OP_SET_PROPERTY:  "OP_SET_PROPERTY",
	OP_METHOD:        "OP_METHOD",
	OP_INVOKE:        "OP_INVOKE",
	OP_INHERIT:       "OP_INHERIT",
	OP_GET_SUPER:     "OP_GET_SUPER",
	OP_SUPER_INVOKE:  "OP_SUPER_INVOKE",
	OP_ARRAY:         "OP_ARRAY",
	OP_GET_INDEX:     "OP_GET_INDEX",
	OP_SET_INDEX:     "OP_SET_INDEX",
	OP_OBJECT:        "OP_OBJECT",
	OP_CLOSE_UPVALUE: "OP_CLOSE_UPVALUE",
	OP_PRINT:         "OP_PRINT",
	OP_PRINTLN:       "OP_PRINTLN",
	OP_BUILTIN:       "OP_BUILTIN",
}

// ReadU16 reads a 16-bit value from code at offset (big-endian).
func ReadU16(code []byte, offset int) uint16 {
	return uint16(code[offset])<<8 | uint16(code[offset+1])
}

// Disassemble returns a human-readable representation of the chunk.
func Disassemble(chunk *Chunk, name string) string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("== %s ==\n", name))

	for offset := 0; offset < len(chunk.Code); {
		offset = disassembleInstruction(chunk, offset, &out)
	}

	return out.String()
}

func disassembleInstruction(chunk *Chunk, offset int, out *strings.Builder) int {
	out.WriteString(fmt.Sprintf("%04d ", offset))

	// Show line number
	if offset > 0 && chunk.Lines[offset] == chunk.Lines[offset-1] {
		out.WriteString("   | ")
	} else {
		out.WriteString(fmt.Sprintf("%4d ", chunk.Lines[offset]))
	}

	instruction := OpCode(chunk.Code[offset])

	switch instruction {
	case OP_CONSTANT:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_NULL, OP_TRUE, OP_FALSE:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_ADD, OP_SUBTRACT, OP_MULTIPLY, OP_DIVIDE, OP_MODULO, OP_NEGATE:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_EQUAL, OP_NOT_EQUAL, OP_LESS, OP_LESS_EQUAL, OP_GREATER, OP_GREATER_EQUAL:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_NOT:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_CONCAT:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_GET_LOCAL, OP_SET_LOCAL:
		return byteInstruction(instruction.String(), chunk, offset, out)
	case OP_GET_GLOBAL, OP_SET_GLOBAL:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_GET_UPVALUE, OP_SET_UPVALUE:
		return byteInstruction(instruction.String(), chunk, offset, out)
	case OP_POP:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_POPN:
		return byteInstruction(instruction.String(), chunk, offset, out)
	case OP_DUP:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_JUMP, OP_JUMP_BACK, OP_JUMP_IF_FALSE, OP_JUMP_IF_TRUE:
		return jumpInstruction(instruction.String(), 1, chunk, offset, out)
	case OP_CALL:
		return byteInstruction(instruction.String(), chunk, offset, out)
	case OP_RETURN:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_CLOSURE:
		return closureInstruction(chunk, offset, out)
	case OP_CLASS:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_GET_PROPERTY, OP_SET_PROPERTY:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_METHOD:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_INVOKE, OP_SUPER_INVOKE:
		return invokeInstruction(instruction.String(), chunk, offset, out)
	case OP_INHERIT:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_GET_SUPER:
		return constantInstruction(instruction.String(), chunk, offset, out)
	case OP_ARRAY:
		return shortInstruction(instruction.String(), chunk, offset, out)
	case OP_GET_INDEX, OP_SET_INDEX:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_OBJECT:
		return shortInstruction(instruction.String(), chunk, offset, out)
	case OP_CLOSE_UPVALUE:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_PRINT, OP_PRINTLN:
		return simpleInstruction(instruction.String(), offset, out)
	case OP_BUILTIN:
		return builtinInstruction(instruction.String(), chunk, offset, out)
	default:
		out.WriteString(fmt.Sprintf("Unknown opcode %d\n", instruction))
		return offset + 1
	}
}

func simpleInstruction(name string, offset int, out *strings.Builder) int {
	out.WriteString(fmt.Sprintf("%s\n", name))
	return offset + 1
}

func constantInstruction(name string, chunk *Chunk, offset int, out *strings.Builder) int {
	constant := ReadU16(chunk.Code, offset+1)
	out.WriteString(fmt.Sprintf("%-16s %4d '", name, constant))
	if int(constant) < len(chunk.Constants) {
		out.WriteString(fmt.Sprintf("%v", chunk.Constants[constant]))
	}
	out.WriteString("'\n")
	return offset + 3
}

func byteInstruction(name string, chunk *Chunk, offset int, out *strings.Builder) int {
	slot := chunk.Code[offset+1]
	out.WriteString(fmt.Sprintf("%-16s %4d\n", name, slot))
	return offset + 2
}

func shortInstruction(name string, chunk *Chunk, offset int, out *strings.Builder) int {
	value := ReadU16(chunk.Code, offset+1)
	out.WriteString(fmt.Sprintf("%-16s %4d\n", name, value))
	return offset + 3
}

func jumpInstruction(name string, sign int, chunk *Chunk, offset int, out *strings.Builder) int {
	jump := int(ReadU16(chunk.Code, offset+1))
	target := offset + 3 + sign*jump
	out.WriteString(fmt.Sprintf("%-16s %4d -> %d\n", name, jump, target))
	return offset + 3
}

func closureInstruction(chunk *Chunk, offset int, out *strings.Builder) int {
	offset++
	constant := ReadU16(chunk.Code, offset)
	offset += 2
	out.WriteString(fmt.Sprintf("%-16s %4d ", "OP_CLOSURE", constant))
	if int(constant) < len(chunk.Constants) {
		out.WriteString(fmt.Sprintf("%v", chunk.Constants[constant]))
	}
	out.WriteString("\n")
	// Note: In a full implementation, we'd also print upvalue info here
	return offset
}

func invokeInstruction(name string, chunk *Chunk, offset int, out *strings.Builder) int {
	constant := ReadU16(chunk.Code, offset+1)
	argCount := chunk.Code[offset+3]
	out.WriteString(fmt.Sprintf("%-16s (%d args) %4d '", name, argCount, constant))
	if int(constant) < len(chunk.Constants) {
		out.WriteString(fmt.Sprintf("%v", chunk.Constants[constant]))
	}
	out.WriteString("'\n")
	return offset + 4
}

func builtinInstruction(name string, chunk *Chunk, offset int, out *strings.Builder) int {
	builtinId := chunk.Code[offset+1]
	argCount := chunk.Code[offset+2]
	out.WriteString(fmt.Sprintf("%-16s %4d (%d args)\n", name, builtinId, argCount))
	return offset + 3
}
