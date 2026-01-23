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
	ip            int
	slotBase      int
	isConstructor bool
}

// VM is the virtual machine that executes bytecode.
type VM struct {
	frames       [FRAMES_MAX]CallFrame
	frameCount   int
	stack        []Value
	sp           int
	globals      map[string]Value
	output       io.Writer
	lastPopped   Value
	openUpvalues *ObjUpvalue
	gc           *GC
}

// New creates a new VM with the given bytecode chunk.
func New(chunk *bytecode.Chunk) *VM {
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
	vm.gc = NewGC(vm)

	vm.gc.Track(fn)
	vm.gc.Track(closure)

	vm.push(ObjectValue(closure))

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
	vm.gc = NewGC(vm)
	vm.gc.Track(closure)
	vm.gc.Track(closure.Function)

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
				return vm.runtimeError("operands must be numbers for '+'")
			}
			vm.push(NumberValue(a.AsNumber() + b.AsNumber()))

		case bytecode.OP_SUBTRACT:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '-'")
			}
			vm.push(NumberValue(a.AsNumber() - b.AsNumber()))

		case bytecode.OP_MULTIPLY:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '*'")
			}
			vm.push(NumberValue(a.AsNumber() * b.AsNumber()))

		case bytecode.OP_DIVIDE:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '/'")
			}
			if b.AsNumber() == 0 {
				return vm.runtimeError("division by zero")
			}
			vm.push(NumberValue(a.AsNumber() / b.AsNumber()))

		case bytecode.OP_MODULO:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '%%'")
			}
			vm.push(NumberValue(math.Mod(a.AsNumber(), b.AsNumber())))

		case bytecode.OP_NEGATE:
			a := vm.pop()
			if !a.IsNumber() {
				return vm.runtimeError("operand must be a number for unary '-'")
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
				return vm.runtimeError("operands must be numbers for '<'")
			}
			vm.push(BoolValue(a.AsNumber() < b.AsNumber()))

		case bytecode.OP_LESS_EQUAL:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '<='")
			}
			vm.push(BoolValue(a.AsNumber() <= b.AsNumber()))

		case bytecode.OP_GREATER:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '>'")
			}
			vm.push(BoolValue(a.AsNumber() > b.AsNumber()))

		case bytecode.OP_GREATER_EQUAL:
			b := vm.pop()
			a := vm.pop()
			if !a.IsNumber() || !b.IsNumber() {
				return vm.runtimeError("operands must be numbers for '>='")
			}
			vm.push(BoolValue(a.AsNumber() >= b.AsNumber()))

		case bytecode.OP_NOT:
			a := vm.pop()
			vm.push(BoolValue(!IsTruthy(a)))

		case bytecode.OP_CONCAT:
			b := vm.pop()
			a := vm.pop()
			if !a.IsString() || !b.IsString() {
				return vm.runtimeError("operands must be strings for string concatenation")
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
				return vm.runtimeError("undefined variable '%s'", name)
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
			var fn *ObjFunction

			switch fnVal := vm.chunk().Constants[fnIdx].(type) {
			case *compiler.ObjFunction:
				fn = &ObjFunction{
					Name:         fnVal.Name,
					Arity:        fnVal.Arity,
					UpvalueCount: fnVal.UpvalueCount,
					Chunk:        fnVal.Chunk,
				}
			case *bytecode.BinaryFunction:
				fn = &ObjFunction{
					Name:         fnVal.Name,
					Arity:        fnVal.Arity,
					UpvalueCount: fnVal.UpvalueCount,
					Chunk:        fnVal.Chunk,
				}
			default:
				return vm.runtimeError("invalid function constant type: %T", fnVal)
			}

			closure := NewObjClosure(fn)

			for i := 0; i < fn.UpvalueCount; i++ {
				isLocal := vm.readByte() == 1
				index := vm.readByte()
				if isLocal {
					closure.Upvalues[i] = vm.captureUpvalue(frame.slotBase + int(index))
				} else {
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

			if frame.isConstructor {
				result = vm.stack[frame.slotBase]
			}

			vm.closeUpvalues(frame.slotBase)

			vm.frameCount--
			if vm.frameCount == 0 {
				vm.pop()
				return nil
			}

			vm.sp = frame.slotBase
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
			for i := count - 1; i >= 0; i-- {
				arr.Elements[i] = vm.pop()
			}
			vm.push(ObjectValue(arr))

		case bytecode.OP_OBJECT:
			count := int(vm.readU16())
			obj := NewObjObject()
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
					return vm.runtimeError("array index must be a number")
				}
				arr := object.AsArray()
				idx := int(index.AsNumber())
				if idx < 0 || idx >= len(arr.Elements) {
					return vm.runtimeError("array index out of bounds: %d (length: %d)", idx, len(arr.Elements))
				}
				vm.push(arr.Elements[idx])
			} else if object.IsString() {
				if !index.IsNumber() {
					return vm.runtimeError("string index must be a number")
				}
				str := object.AsString()
				idx := int(index.AsNumber())
				if idx < 0 || idx >= len(str) {
					return vm.runtimeError("string index out of bounds: %d (length: %d)", idx, len(str))
				}
				vm.push(ObjectValue(NewObjString(string(str[idx]))))
			} else {
				return vm.runtimeError("cannot index value of type %s", valueTypeName(object))
			}

		case bytecode.OP_SET_INDEX:
			value := vm.pop()
			index := vm.pop()
			object := vm.pop()

			if !object.IsArray() {
				return vm.runtimeError("can only index-assign to arrays, got %s", valueTypeName(object))
			}
			if !index.IsNumber() {
				return vm.runtimeError("array index must be a number")
			}

			arr := object.AsArray()
			idx := int(index.AsNumber())
			if idx < 0 || idx >= len(arr.Elements) {
				return vm.runtimeError("array index out of bounds: %d (length: %d)", idx, len(arr.Elements))
			}
			arr.Elements[idx] = value
			vm.push(value)

		case bytecode.OP_GET_PROPERTY:
			nameIdx := vm.readU16()
			name := vm.chunk().Constants[nameIdx].(string)
			object := vm.pop()

			if object.IsInstance() {
				instance := object.AsInstance()
				if value, ok := instance.Fields[name]; ok {
					vm.push(value)
				} else if method := instance.Class.Methods[name]; method != nil {
					bound := &ObjBoundMethod{
						Receiver: object,
						Method:   method,
					}
					vm.push(ObjectValue(bound))
				} else {
					return vm.runtimeError("undefined property '%s' on instance of '%s'", name, instance.Class.Name)
				}
			} else if obj, ok := object.AsObject().(*ObjObject); ok {
				if value, exists := obj.Fields[name]; exists {
					vm.push(value)
				} else {
					vm.push(NullValue())
				}
			} else {
				return vm.runtimeError("cannot access property '%s' on value of type %s", name, valueTypeName(object))
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
				return vm.runtimeError("cannot set property '%s' on value of type %s", name, valueTypeName(object))
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
				return vm.runtimeError("superclass must be a class")
			}
			subclass.Super = superclass
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
				return vm.runtimeError("can only invoke methods on instances, got %s", valueTypeName(receiver))
			}
			instance := receiver.AsInstance()

			if value, ok := instance.Fields[name]; ok && value.IsClosure() {
				vm.stack[vm.sp-argCount-1] = value
				if err := vm.call(value.AsClosure(), argCount, false); err != nil {
					return err
				}
				continue
			}

			method := instance.Class.Methods[name]
			if method == nil {
				return vm.runtimeError("undefined method '%s' on instance of '%s'", name, instance.Class.Name)
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
				return vm.runtimeError("class '%s' has no superclass", instance.Class.Name)
			}
			method := superclass.Methods[name]
			if method == nil {
				return vm.runtimeError("undefined method '%s' in superclass '%s'", name, superclass.Name)
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
				return vm.runtimeError("class '%s' has no superclass", instance.Class.Name)
			}
			method := superclass.Methods[name]
			if method == nil {
				return vm.runtimeError("undefined method '%s' in superclass '%s'", name, superclass.Name)
			}
			if err := vm.call(method, argCount, false); err != nil {
				return err
			}

		default:
			return vm.runtimeError("unknown opcode: %v", op)
		}
	}
}

// callValue calls a value as a function.
func (vm *VM) callValue(callee Value, argCount int) error {
	if callee.IsClosure() {
		return vm.call(callee.AsClosure(), argCount, false)
	}

	if callee.IsClass() {
		class := callee.AsClass()
		instance := NewObjInstance(class)
		vm.stack[vm.sp-argCount-1] = ObjectValue(instance)

		if constructor, ok := class.Methods["constructor"]; ok {
			return vm.call(constructor, argCount, true)
		}
		if argCount != 0 {
			return vm.runtimeError("class '%s' has no constructor but was called with %d arguments", class.Name, argCount)
		}
		return nil
	}

	if callee.IsBoundMethod() {
		bound := callee.AsBoundMethod()
		vm.stack[vm.sp-argCount-1] = bound.Receiver
		return vm.call(bound.Method, argCount, false)
	}

	return vm.runtimeError("cannot call value of type %s", valueTypeName(callee))
}

// call invokes a closure with the given arguments.
func (vm *VM) call(closure *ObjClosure, argCount int, isConstructor bool) error {
	if argCount != closure.Function.Arity {
		return vm.runtimeError("function '%s' expected %d arguments but got %d",
			closure.Function.Name, closure.Function.Arity, argCount)
	}

	if vm.frameCount >= FRAMES_MAX {
		return vm.runtimeError("stack overflow (max call depth: %d)", FRAMES_MAX)
	}

	frame := &vm.frames[vm.frameCount]
	vm.frameCount++

	frame.closure = closure
	frame.ip = 0
	frame.isConstructor = isConstructor
	frame.slotBase = vm.sp - argCount - 1

	return nil
}

// captureUpvalue creates or reuses an upvalue for the given stack slot.
func (vm *VM) captureUpvalue(stackIndex int) *ObjUpvalue {
	var prevUpvalue *ObjUpvalue
	upvalue := vm.openUpvalues

	for upvalue != nil && upvalue.stackIndex > stackIndex {
		prevUpvalue = upvalue
		upvalue = upvalue.Next
	}

	if upvalue != nil && upvalue.stackIndex == stackIndex {
		return upvalue
	}

	newUpvalue := NewObjUpvalue(&vm.stack[stackIndex])
	newUpvalue.stackIndex = stackIndex
	newUpvalue.Next = upvalue

	if prevUpvalue == nil {
		vm.openUpvalues = newUpvalue
	} else {
		prevUpvalue.Next = newUpvalue
	}

	return newUpvalue
}

// closeUpvalues closes all upvalues at or above the given stack index.
func (vm *VM) closeUpvalues(lastSlot int) {
	for vm.openUpvalues != nil && vm.openUpvalues.stackIndex >= lastSlot {
		upvalue := vm.openUpvalues
		upvalue.Closed = *upvalue.Location
		upvalue.Location = &upvalue.Closed
		vm.openUpvalues = upvalue.Next
	}
}

func (vm *VM) callBuiltin(builtinID int, argCount int) error {
	switch builtinID {
	case bytecode.BUILTIN_PRINTLN:
		if argCount != 1 {
			return vm.runtimeError("println expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprintln(vm.output, val.String())
		return nil

	case bytecode.BUILTIN_PRINT:
		if argCount != 1 {
			return vm.runtimeError("print expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		fmt.Fprint(vm.output, val.String())
		return nil

	case bytecode.BUILTIN_LEN:
		if argCount != 1 {
			return vm.runtimeError("len expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if val.IsString() {
			vm.push(NumberValue(float64(len(val.AsString()))))
		} else if val.IsArray() {
			vm.push(NumberValue(float64(len(val.AsArray().Elements))))
		} else {
			return vm.runtimeError("len: argument must be string or array, got %s", valueTypeName(val))
		}
		return nil

	case bytecode.BUILTIN_TYPEOF:
		if argCount != 1 {
			return vm.runtimeError("typeof expects 1 argument, got %d", argCount)
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

	case bytecode.BUILTIN_PUSH:
		if argCount != 2 {
			return vm.runtimeError("push expects 2 arguments, got %d", argCount)
		}
		val := vm.pop()
		arr := vm.pop()
		if !arr.IsArray() {
			return vm.runtimeError("push: first argument must be an array, got %s", valueTypeName(arr))
		}
		arr.AsArray().Elements = append(arr.AsArray().Elements, val)
		vm.push(NumberValue(float64(len(arr.AsArray().Elements))))
		return nil

	case bytecode.BUILTIN_POP:
		if argCount != 1 {
			return vm.runtimeError("pop expects 1 argument, got %d", argCount)
		}
		arr := vm.pop()
		if !arr.IsArray() {
			return vm.runtimeError("pop: argument must be an array, got %s", valueTypeName(arr))
		}
		elements := arr.AsArray().Elements
		if len(elements) == 0 {
			return vm.runtimeError("pop: cannot pop from empty array")
		}
		lastVal := elements[len(elements)-1]
		arr.AsArray().Elements = elements[:len(elements)-1]
		vm.push(lastVal)
		return nil

	case bytecode.BUILTIN_TOSTRING:
		if argCount != 1 {
			return vm.runtimeError("toString expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		vm.push(ObjectValue(NewObjString(val.String())))
		return nil

	case bytecode.BUILTIN_TONUMBER:
		if argCount != 1 {
			return vm.runtimeError("toNumber expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if val.IsNumber() {
			vm.push(val)
		} else if val.IsString() {
			var num float64
			_, err := fmt.Sscanf(val.AsString(), "%f", &num)
			if err != nil {
				vm.push(NumberValue(math.NaN()))
			} else {
				vm.push(NumberValue(num))
			}
		} else if val.IsBool() {
			if val.AsBool() {
				vm.push(NumberValue(1))
			} else {
				vm.push(NumberValue(0))
			}
		} else if val.IsNull() {
			vm.push(NumberValue(0))
		} else {
			vm.push(NumberValue(math.NaN()))
		}
		return nil

	case bytecode.BUILTIN_SQRT:
		if argCount != 1 {
			return vm.runtimeError("sqrt expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if !val.IsNumber() {
			return vm.runtimeError("sqrt: argument must be a number, got %s", valueTypeName(val))
		}
		vm.push(NumberValue(math.Sqrt(val.AsNumber())))
		return nil

	case bytecode.BUILTIN_FLOOR:
		if argCount != 1 {
			return vm.runtimeError("floor expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if !val.IsNumber() {
			return vm.runtimeError("floor: argument must be a number, got %s", valueTypeName(val))
		}
		vm.push(NumberValue(math.Floor(val.AsNumber())))
		return nil

	case bytecode.BUILTIN_CEIL:
		if argCount != 1 {
			return vm.runtimeError("ceil expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if !val.IsNumber() {
			return vm.runtimeError("ceil: argument must be a number, got %s", valueTypeName(val))
		}
		vm.push(NumberValue(math.Ceil(val.AsNumber())))
		return nil

	case bytecode.BUILTIN_ABS:
		if argCount != 1 {
			return vm.runtimeError("abs expects 1 argument, got %d", argCount)
		}
		val := vm.pop()
		if !val.IsNumber() {
			return vm.runtimeError("abs: argument must be a number, got %s", valueTypeName(val))
		}
		vm.push(NumberValue(math.Abs(val.AsNumber())))
		return nil

	default:
		return vm.runtimeError("unknown builtin: %d", builtinID)
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
		return ObjectValue(val)
	case *compiler.ObjFunction:
		fn := &ObjFunction{
			Name:         val.Name,
			Arity:        val.Arity,
			UpvalueCount: val.UpvalueCount,
			Chunk:        val.Chunk,
		}
		return ObjectValue(fn)
	default:
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

// SetOutput sets the output writer for the VM.
func (vm *VM) SetOutput(w io.Writer) {
	vm.output = w
}

// SetGlobal sets a global variable.
func (vm *VM) SetGlobal(name string, val Value) {
	vm.globals[name] = val
}

// GetGlobals returns a copy of the globals map.
func (vm *VM) GetGlobals() map[string]Value {
	result := make(map[string]Value)
	for k, v := range vm.globals {
		result[k] = v
	}
	return result
}

// LastPopped returns the last popped value.
func (vm *VM) LastPopped() Value {
	return vm.lastPopped
}

// valueTypeName returns a human-readable name for a value's type.
func valueTypeName(v Value) string {
	switch v.Type {
	case VAL_NULL:
		return "null"
	case VAL_BOOL:
		return "boolean"
	case VAL_NUMBER:
		return "number"
	case VAL_OBJECT:
		switch v.obj.Type() {
		case OBJ_STRING:
			return "string"
		case OBJ_ARRAY:
			return "array"
		case OBJ_OBJECT:
			return "object"
		case OBJ_FUNCTION, OBJ_CLOSURE:
			return "function"
		case OBJ_CLASS:
			return "class"
		case OBJ_INSTANCE:
			return "instance"
		case OBJ_BOUND_METHOD:
			return "method"
		default:
			return "object"
		}
	}
	return "unknown"
}
