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
	STACK_MAX  = 256
	FRAMES_MAX = 64
)

// CallFrame represents a function call frame.
type CallFrame struct {
	closure       *ObjClosure
	ip            int  // Instruction pointer (offset into closure.Function.Chunk.Code)
	slotBase      int  // Index in the stack where this frame's locals start
	isConstructor bool // True if this frame is a constructor call
}

// VM is the virtual machine that executes bytecode.
type VM struct {
	frames     [FRAMES_MAX]CallFrame
	frameCount int
	stack      []Value          // Value stack
	sp         int              // Stack pointer (points to next free slot)
	globals    map[string]Value // Global variables
	output     io.Writer
	lastPopped Value            // Last value popped (for testing)

	// Open upvalue linked list (sorted by stack slot, from top to bottom)
	openUpvalues *ObjUpvalue
}

// New creates a new VM with the given bytecode chunk.
// This creates a script-level function from the chunk.
func New(chunk *bytecode.Chunk) *VM {
	// Create a script function from the chunk
	fn := &ObjFunction{
		Name:  "",
		Arity: 0,
		Chunk: chunk,
	}
	closure := NewObjClosure(fn)

	vm := &VM{
		stack:   make([]Value, STACK_MAX),
		sp:      0,
		globals: make(map[string]Value),
		output:  os.Stdout,
	}

	// Push the closure as the first stack slot (slot 0 for the script)
	vm.push(ObjectValue(closure))

	// Set up the initial call frame
	vm.frames[0] = CallFrame{
		closure:  closure,
		ip:       0,
		slotBase: 0,
	}
	vm.frameCount = 1

	return vm
}

// NewWithClosure creates a new VM with a closure.
func NewWithClosure(closure *ObjClosure) *VM {
	vm := &VM{
		stack:   make([]Value, STACK_MAX),
		sp:      0,
		globals: make(map[string]Value),
		output:  os.Stdout,
	}

	// Push the closure as the first stack slot
	vm.push(ObjectValue(closure))

	// Set up the initial call frame
	vm.frames[0] = CallFrame{
		closure:  closure,
		ip:       0,
		slotBase: 0,
	}
	vm.frameCount = 1

	return vm
}

// frame returns the current call frame.
func (vm *VM) frame() *CallFrame {
	return &vm.frames[vm.frameCount-1]
}

// chunk returns the current function's bytecode chunk.
func (vm *VM) chunk() *bytecode.Chunk {
	return vm.frame().closure.Function.Chunk
}

// Run executes the bytecode.
func (vm *VM) Run() error {
	for {
		frame := vm.frame()
		op := bytecode.OpCode(vm.readByte())

		switch op {
		case bytecode.OP_CONSTANT:
			idx := vm.readU16()
			constant := vm.chunk().Constants[idx]
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

		case bytecode.OP_GET_LOCAL:
			slot := vm.readByte()
			vm.push(vm.stack[frame.slotBase+int(slot)])

		case bytecode.OP_SET_LOCAL:
			slot := vm.readByte()
			vm.stack[frame.slotBase+int(slot)] = vm.peek(0)

		case bytecode.OP_GET_GLOBAL:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			val, exists := vm.globals[name]
			if !exists {
				return fmt.Errorf("undefined variable: %s", name)
			}
			vm.push(val)

		case bytecode.OP_SET_GLOBAL:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			vm.globals[name] = vm.peek(0)

		case bytecode.OP_GET_UPVALUE:
			slot := vm.readByte()
			upvalue := frame.closure.Upvalues[slot]
			vm.push(*upvalue.Location)

		case bytecode.OP_SET_UPVALUE:
			slot := vm.readByte()
			upvalue := frame.closure.Upvalues[slot]
			*upvalue.Location = vm.peek(0)

		case bytecode.OP_CLOSE_UPVALUE:
			vm.closeUpvalues(vm.sp - 1)
			vm.pop()

		case bytecode.OP_JUMP:
			offset := vm.readU16()
			frame.ip += int(offset)

		case bytecode.OP_JUMP_BACK:
			offset := vm.readU16()
			frame.ip -= int(offset)

		case bytecode.OP_JUMP_IF_FALSE:
			offset := vm.readU16()
			if !IsTruthy(vm.peek(0)) {
				frame.ip += int(offset)
			}

		case bytecode.OP_JUMP_IF_TRUE:
			offset := vm.readU16()
			if IsTruthy(vm.peek(0)) {
				frame.ip += int(offset)
			}

		case bytecode.OP_CLOSURE:
			fnIdx := vm.readU16()
			// The compiler stores *compiler.ObjFunction, convert to *vm.ObjFunction
			compilerFn := vm.chunk().Constants[fnIdx].(*compiler.ObjFunction)
			fn := &ObjFunction{
				Name:         compilerFn.Name,
				Arity:        compilerFn.Arity,
				UpvalueCount: compilerFn.UpvalueCount,
				Chunk:        compilerFn.Chunk,
			}
			closure := NewObjClosure(fn)

			// Read upvalue descriptors
			for i := 0; i < fn.UpvalueCount; i++ {
				isLocal := vm.readByte() == 1
				index := vm.readByte()
				if isLocal {
					// Capture local from the enclosing function
					closure.Upvalues[i] = vm.captureUpvalue(frame.slotBase + int(index))
				} else {
					// Capture upvalue from the enclosing function
					closure.Upvalues[i] = frame.closure.Upvalues[index]
				}
			}

			vm.push(ObjectValue(closure))

		case bytecode.OP_CALL:
			argCount := int(vm.readByte())
			if err := vm.callValue(vm.peek(argCount), argCount); err != nil {
				return err
			}

		case bytecode.OP_RETURN:
			result := vm.pop()

			// If this is a constructor, return the instance (slot 0) instead
			if frame.isConstructor {
				result = vm.stack[frame.slotBase]
			}

			// Close any open upvalues in this frame
			vm.closeUpvalues(frame.slotBase)

			// Pop the frame
			vm.frameCount--
			if vm.frameCount == 0 {
				// We're returning from the script
				vm.pop() // Pop the script closure
				return nil
			}

			// Discard the called function and its locals
			vm.sp = frame.slotBase

			// Push the return value
			vm.push(result)

		case bytecode.OP_BUILTIN:
			builtinID := vm.readByte()
			argCount := vm.readByte()
			if err := vm.callBuiltin(int(builtinID), int(argCount)); err != nil {
				return err
			}

		case bytecode.OP_ARRAY:
			count := int(vm.readU16())
			arr := NewObjArray()
			arr.Elements = make([]Value, count)
			// Pop elements in reverse order (they were pushed left-to-right)
			for i := count - 1; i >= 0; i-- {
				arr.Elements[i] = vm.pop()
			}
			vm.push(ObjectValue(arr))

		case bytecode.OP_OBJECT:
			count := int(vm.readU16())
			obj := NewObjObject()
			// Pop key-value pairs in reverse order
			for i := 0; i < count; i++ {
				value := vm.pop()
				key := vm.pop()
				if !key.IsString() {
					return fmt.Errorf("object key must be a string")
				}
				obj.Fields[key.AsString()] = value
			}
			vm.push(ObjectValue(obj))

		case bytecode.OP_GET_INDEX:
			index := vm.pop()
			object := vm.pop()

			if object.IsArray() {
				if !index.IsNumber() {
					return fmt.Errorf("array index must be a number")
				}
				arr := object.AsArray()
				idx := int(index.AsNumber())
				if idx < 0 || idx >= len(arr.Elements) {
					return fmt.Errorf("array index out of bounds: %d", idx)
				}
				vm.push(arr.Elements[idx])
			} else if object.IsString() {
				if !index.IsNumber() {
					return fmt.Errorf("string index must be a number")
				}
				str := object.AsString()
				idx := int(index.AsNumber())
				if idx < 0 || idx >= len(str) {
					return fmt.Errorf("string index out of bounds: %d", idx)
				}
				vm.push(ObjectValue(NewObjString(string(str[idx]))))
			} else {
				return fmt.Errorf("cannot index type %T", object.obj)
			}

		case bytecode.OP_SET_INDEX:
			value := vm.pop()
			index := vm.pop()
			object := vm.pop()

			if !object.IsArray() {
				return fmt.Errorf("can only index-assign to arrays")
			}
			if !index.IsNumber() {
				return fmt.Errorf("array index must be a number")
			}
			arr := object.AsArray()
			idx := int(index.AsNumber())
			if idx < 0 || idx >= len(arr.Elements) {
				return fmt.Errorf("array index out of bounds: %d", idx)
			}
			arr.Elements[idx] = value
			vm.push(value) // Assignment expression returns the value

		case bytecode.OP_GET_PROPERTY:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			object := vm.pop()

			if object.IsInstance() {
				instance := object.AsInstance()
				// First check fields
				if value, ok := instance.Fields[name]; ok {
					vm.push(value)
				} else if method := instance.Class.Methods[name]; method != nil {
					// Bind method to instance
					bound := &ObjBoundMethod{
						Receiver: object,
						Method:   method,
					}
					vm.push(ObjectValue(bound))
				} else {
					return fmt.Errorf("undefined property: %s", name)
				}
			} else if obj, ok := object.AsObject().(*ObjObject); ok {
				if value, exists := obj.Fields[name]; exists {
					vm.push(value)
				} else {
					vm.push(NullValue())
				}
			} else {
				return fmt.Errorf("cannot access property on %T", object.obj)
			}

		case bytecode.OP_SET_PROPERTY:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			value := vm.pop()
			object := vm.pop()

			if object.IsInstance() {
				instance := object.AsInstance()
				instance.Fields[name] = value
				vm.push(value)
			} else if obj, ok := object.AsObject().(*ObjObject); ok {
				obj.Fields[name] = value
				vm.push(value)
			} else {
				return fmt.Errorf("cannot set property on %T", object.obj)
			}

		case bytecode.OP_CLASS:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			class := NewObjClass(name)
			vm.push(ObjectValue(class))

		case bytecode.OP_METHOD:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			method := vm.pop()
			class := vm.peek(0).AsClass()
			class.Methods[name] = method.AsClosure()

		case bytecode.OP_INHERIT:
			subclass := vm.pop().AsClass()
			superclass := vm.pop().AsClass()
			if superclass == nil {
				return fmt.Errorf("superclass must be a class")
			}
			subclass.Super = superclass
			// Copy methods from superclass
			for name, method := range superclass.Methods {
				if _, exists := subclass.Methods[name]; !exists {
					subclass.Methods[name] = method
				}
			}

		case bytecode.OP_INVOKE:
			nameIdx := vm.readU16()
			argCount := int(vm.readByte())
			name := vm.chunk().Constants[nameIdx].(string)

			receiver := vm.peek(argCount)
			if !receiver.IsInstance() {
				return fmt.Errorf("can only invoke methods on instances")
			}
			instance := receiver.AsInstance()

			// Check for field first (could be a closure stored as a field)
			if value, ok := instance.Fields[name]; ok {
				if value.IsClosure() {
					vm.stack[vm.sp-argCount-1] = value
					if err := vm.call(value.AsClosure(), argCount, false); err != nil {
						return err
					}
					continue
				}
			}

			// Look up method
			method := instance.Class.Methods[name]
			if method == nil {
				return fmt.Errorf("undefined method: %s", name)
			}

			if err := vm.call(method, argCount, false); err != nil {
				return err
			}

		case bytecode.OP_GET_SUPER:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			instance := vm.pop().AsInstance()
			superclass := instance.Class.Super
			if superclass == nil {
				return fmt.Errorf("no superclass")
			}
			method := superclass.Methods[name]
			if method == nil {
				return fmt.Errorf("undefined method in superclass: %s", name)
			}
			bound := &ObjBoundMethod{
				Receiver: ObjectValue(instance),
				Method:   method,
			}
			vm.push(ObjectValue(bound))

		case bytecode.OP_SUPER_INVOKE:
			nameIdx := vm.readU16()
			argCount := int(vm.readByte())
			name := vm.chunk().Constants[nameIdx].(string)

			receiver := vm.peek(argCount)
			instance := receiver.AsInstance()
			superclass := instance.Class.Super
			if superclass == nil {
				return fmt.Errorf("no superclass")
			}
			method := superclass.Methods[name]
			if method == nil {
				return fmt.Errorf("undefined method in superclass: %s", name)
			}
			// Super constructor calls are not marked as constructors - the child
			// constructor will return the instance
			if err := vm.call(method, argCount, false); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown opcode: %v", op)
		}
	}
}

// callValue calls a value as a function.
func (vm *VM) callValue(callee Value, argCount int) error {
	if callee.IsClosure() {
		return vm.call(callee.AsClosure(), argCount, false)
	}
	if callee.IsClass() {
		// Create a new instance
		class := callee.AsClass()
		instance := NewObjInstance(class)
		// Replace the class on the stack with the instance
		vm.stack[vm.sp-argCount-1] = ObjectValue(instance)

		// Call the constructor if it exists
		if constructor, ok := class.Methods["constructor"]; ok {
			return vm.call(constructor, argCount, true)
		} else if argCount != 0 {
			return fmt.Errorf("class %s has no constructor but was called with %d arguments", class.Name, argCount)
		}
		return nil
	}
	if callee.IsBoundMethod() {
		bound := callee.AsBoundMethod()
		// Put the receiver in slot 0
		vm.stack[vm.sp-argCount-1] = bound.Receiver
		return vm.call(bound.Method, argCount, false)
	}
	return fmt.Errorf("can only call functions and classes")
}

// call invokes a closure with the given arguments.
func (vm *VM) call(closure *ObjClosure, argCount int, isConstructor bool) error {
	if argCount != closure.Function.Arity {
		return fmt.Errorf("expected %d arguments but got %d", closure.Function.Arity, argCount)
	}

	if vm.frameCount >= FRAMES_MAX {
		return fmt.Errorf("stack overflow")
	}

	frame := &vm.frames[vm.frameCount]
	vm.frameCount++

	frame.closure = closure
	frame.ip = 0
	frame.isConstructor = isConstructor
	// The slot base is where the function sits on the stack
	// (args are above the function slot)
	frame.slotBase = vm.sp - argCount - 1

	return nil
}

// captureUpvalue creates or reuses an upvalue for the given stack slot.
func (vm *VM) captureUpvalue(stackIndex int) *ObjUpvalue {
	// Look for an existing open upvalue for this slot
	var prevUpvalue *ObjUpvalue
	upvalue := vm.openUpvalues

	// Walk the list to find the right position (sorted by slot, descending)
	for upvalue != nil && upvalue.stackIndex > stackIndex {
		prevUpvalue = upvalue
		upvalue = upvalue.Next
	}

	// If we found an upvalue for this slot, reuse it
	if upvalue != nil && upvalue.stackIndex == stackIndex {
		return upvalue
	}

	// Create a new upvalue
	newUpvalue := NewObjUpvalue(&vm.stack[stackIndex])
	newUpvalue.stackIndex = stackIndex
	newUpvalue.Next = upvalue

	// Insert into the linked list
	if prevUpvalue == nil {
		vm.openUpvalues = newUpvalue
	} else {
		prevUpvalue.Next = newUpvalue
	}

	return newUpvalue
}

// closeUpvalues closes all upvalues that refer to stack slots at or above the given index.
func (vm *VM) closeUpvalues(lastSlot int) {
	for vm.openUpvalues != nil && vm.openUpvalues.stackIndex >= lastSlot {
		upvalue := vm.openUpvalues

		// Copy the value from the stack to the upvalue's Closed field
		upvalue.Closed = *upvalue.Location
		// Point Location at the Closed field
		upvalue.Location = &upvalue.Closed

		// Remove from the open list
		vm.openUpvalues = upvalue.Next
	}
}

func (vm *VM) callBuiltin(builtinID int, argCount int) error {
	switch builtinID {
	case bytecode.BUILTIN_PRINTLN:
		if argCount != 1 {
			return fmt.Errorf("println expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprintln(vm.output, val.String())
		return nil

	case bytecode.BUILTIN_PRINT:
		if argCount != 1 {
			return fmt.Errorf("print expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprint(vm.output, val.String())
		return nil

	case bytecode.BUILTIN_LEN:
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

	case bytecode.BUILTIN_TYPEOF:
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
	case *ObjFunction:
		// Functions are stored as ObjFunction in constants
		// but are wrapped in closures at runtime
		return ObjectValue(val)
	case *compiler.ObjFunction:
		// The compiler stores its own ObjFunction type, convert it
		fn := &ObjFunction{
			Name:         val.Name,
			Arity:        val.Arity,
			UpvalueCount: val.UpvalueCount,
			Chunk:        val.Chunk,
		}
		return ObjectValue(fn)
	default:
		// This shouldn't happen if the compiler is correct
		panic(fmt.Sprintf("unexpected constant type: %T", v))
	}
}

func (vm *VM) readByte() byte {
	frame := vm.frame()
	b := frame.closure.Function.Chunk.Code[frame.ip]
	frame.ip++
	return b
}

func (vm *VM) readU16() uint16 {
	frame := vm.frame()
	val := bytecode.ReadU16(frame.closure.Function.Chunk.Code, frame.ip)
	frame.ip += 2
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
