package vm

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/pocketlang/gots/pkg/bytecode"
	"github.com/pocketlang/gots/pkg/compiler"
)

const (
	STACK_MAX = 256
)

// VM is the virtual machine that executes bytecode.
type VM struct {
	chunk      *bytecode.Chunk
	ip         int              // Instruction pointer
	stack      []Value          // Value stack
	sp         int              // Stack pointer (points to next free slot)
	globals    map[string]Value // Global variables
	output     io.Writer
	lastPopped Value            // Last value popped (for testing)
}

// New creates a new VM with the given bytecode chunk.
func New(chunk *bytecode.Chunk) *VM {
	return &VM{
		chunk:   chunk,
		ip:      0,
		stack:   make([]Value, STACK_MAX),
		sp:      0,
		globals: make(map[string]Value),
		output:  os.Stdout,
	}
}

// Run executes the bytecode.
func (vm *VM) Run() error {
	for {
		op := bytecode.OpCode(vm.readByte())

		switch op {
		case bytecode.OP_CONSTANT:
			idx := vm.readU16()
			constant := vm.chunk.Constants[idx]
			vm.push(vm.anyToValue(constant))

		case bytecode.OP_NULL:
			vm.push(NullValue())

		case bytecode.OP_TRUE:
			vm.push(BoolValue(true))

		case bytecode.OP_FALSE:
			vm.push(BoolValue(false))

		case bytecode.OP_ADD:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(NumberValue(a.AsNumber() + b.AsNumber()))

		case bytecode.OP_SUBTRACT:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(NumberValue(a.AsNumber() - b.AsNumber()))

		case bytecode.OP_MULTIPLY:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(NumberValue(a.AsNumber() * b.AsNumber()))

		case bytecode.OP_DIVIDE:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			if b.AsNumber() == 0 {
				return fmt.Errorf("division by zero")
			}
			vm.push(NumberValue(a.AsNumber() / b.AsNumber()))

		case bytecode.OP_MODULO:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(NumberValue(math.Mod(a.AsNumber(), b.AsNumber())))

		case bytecode.OP_NEGATE:
			a := vm.pop()
			if !a.IsNumber() {
				return fmt.Errorf("operand must be a number")
			}
			vm.push(NumberValue(-a.AsNumber()))

		case bytecode.OP_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(BoolValue(ValuesEqual(a, b)))

		case bytecode.OP_NOT_EQUAL:
			b := vm.pop()
			a := vm.pop()
			vm.push(BoolValue(!ValuesEqual(a, b)))

		case bytecode.OP_LESS:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(BoolValue(a.AsNumber() < b.AsNumber()))

		case bytecode.OP_LESS_EQUAL:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(BoolValue(a.AsNumber() <= b.AsNumber()))

		case bytecode.OP_GREATER:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(BoolValue(a.AsNumber() > b.AsNumber()))

		case bytecode.OP_GREATER_EQUAL:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return fmt.Errorf("operands must be numbers")
			}
			vm.push(BoolValue(a.AsNumber() >= b.AsNumber()))

		case bytecode.OP_NOT:
			a := vm.pop()
			vm.push(BoolValue(!IsTruthy(a)))

		case bytecode.OP_CONCAT:
			b := vm.pop()
			a := vm.pop()
			if !a.IsString() || !b.IsString() {
				return fmt.Errorf("operands must be strings")
			}
			result := a.AsString() + b.AsString()
			vm.push(ObjectValue(NewObjString(result)))

		case bytecode.OP_POP:
			vm.lastPopped = vm.pop()

		case bytecode.OP_GET_GLOBAL:
			nameIdx := vm.readU16()
			name := vm.chunk.Constants[nameIdx].(string)
			val, exists := vm.globals[name]
			if !exists {
				return fmt.Errorf("undefined variable: %s", name)
			}
			vm.push(val)

		case bytecode.OP_SET_GLOBAL:
			nameIdx := vm.readU16()
			name := vm.chunk.Constants[nameIdx].(string)
			vm.globals[name] = vm.peek(0)

		case bytecode.OP_JUMP:
			offset := vm.readU16()
			vm.ip += int(offset)

		case bytecode.OP_JUMP_BACK:
			offset := vm.readU16()
			vm.ip -= int(offset)

		case bytecode.OP_JUMP_IF_FALSE:
			offset := vm.readU16()
			if !IsTruthy(vm.peek(0)) {
				vm.ip += int(offset)
			}

		case bytecode.OP_JUMP_IF_TRUE:
			offset := vm.readU16()
			if IsTruthy(vm.peek(0)) {
				vm.ip += int(offset)
			}

		case bytecode.OP_BUILTIN:
			builtinID := vm.readByte()
			argCount := vm.readByte()
			if err := vm.callBuiltin(int(builtinID), int(argCount)); err != nil {
				return err
			}

		case bytecode.OP_RETURN:
			return nil

		default:
			return fmt.Errorf("unknown opcode: %v", op)
		}
	}
}

func (vm *VM) callBuiltin(builtinID int, argCount int) error {
	switch builtinID {
	case compiler.BUILTIN_PRINTLN:
		if argCount != 1 {
			return fmt.Errorf("println expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprintln(vm.output, val.String())
		return nil

	case compiler.BUILTIN_PRINT:
		if argCount != 1 {
			return fmt.Errorf("print expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprint(vm.output, val.String())
		return nil

	case compiler.BUILTIN_LEN:
		if argCount != 1 {
			return fmt.Errorf("len expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if val.IsString() {
			vm.push(NumberValue(float64(len(val.AsString()))))
		} else if val.IsArray() {
			vm.push(NumberValue(float64(len(val.AsArray().Elements))))
		} else {
			return fmt.Errorf("len: argument must be string or array")
		}
		return nil

	case compiler.BUILTIN_TYPEOF:
		if argCount != 1 {
			return fmt.Errorf("typeof expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		var typeName string
		switch val.Type {
		case VAL_NULL:
			typeName = "null"
		case VAL_BOOL:
			typeName = "boolean"
		case VAL_NUMBER:
			typeName = "number"
		case VAL_OBJECT:
			switch val.obj.Type() {
			case OBJ_STRING:
				typeName = "string"
			case OBJ_ARRAY:
				typeName = "array"
			case OBJ_FUNCTION, OBJ_CLOSURE:
				typeName = "function"
			case OBJ_CLASS:
				typeName = "class"
			case OBJ_INSTANCE:
				typeName = "object"
			default:
				typeName = "object"
			}
		}
		vm.push(ObjectValue(NewObjString(typeName)))
		return nil

	default:
		return fmt.Errorf("unknown builtin: %d", builtinID)
	}
}

func (vm *VM) anyToValue(v any) Value {
	switch val := v.(type) {
	case float64:
		return NumberValue(val)
	case int:
		return NumberValue(float64(val))
	case bool:
		return BoolValue(val)
	case string:
		return ObjectValue(NewObjString(val))
	case nil:
		return NullValue()
	default:
		// This shouldn't happen if the compiler is correct
		panic(fmt.Sprintf("unexpected constant type: %T", v))
	}
}

func (vm *VM) readByte() byte {
	b := vm.chunk.Code[vm.ip]
	vm.ip++
	return b
}

func (vm *VM) readU16() uint16 {
	val := bytecode.ReadU16(vm.chunk.Code, vm.ip)
	vm.ip += 2
	return val
}

func (vm *VM) push(val Value) {
	vm.stack[vm.sp] = val
	vm.sp++
}

func (vm *VM) pop() Value {
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *VM) peek(distance int) Value {
	return vm.stack[vm.sp-1-distance]
}
