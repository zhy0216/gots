# GoTS Bytecode & VM Specification v1.0

A stack-based virtual machine for executing GoTS programs.

---

## 1. Design Overview

### 1.1 Architecture

- **Stack-based**: Operands pushed/popped from evaluation stack
- **Bytecode**: Variable-length instructions
- **Memory**: Garbage-collected heap for objects
- **Closures**: Upvalue mechanism for captured variables

### 1.2 Execution Model

```
┌─────────────────────────────────────────────────────┐
│                      VM State                        │
├─────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐  │
│  │  Stack  │  │  Frames │  │       Globals       │  │
│  │ [value] │  │ [frame] │  │  map[string]Value   │  │
│  │ [value] │  │ [frame] │  └─────────────────────┘  │
│  │ [value] │  │   ...   │                           │
│  │   ...   │  └─────────┘  ┌─────────────────────┐  │
│  └─────────┘               │        Heap         │  │
│                            │ (GC-managed objects)│  │
│                            └─────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

---

## 2. Value Representation

### 2.1 Value Types

All values in the VM are tagged unions:

```go
type ValueType byte

const (
    VAL_NULL    ValueType = 0x00
    VAL_BOOL    ValueType = 0x01
    VAL_NUMBER  ValueType = 0x02
    VAL_OBJECT  ValueType = 0x03  // Pointer to heap object
)

type Value struct {
    Type ValueType
    Data uint64  // Holds bool, float64 bits, or object pointer
}
```

### 2.2 Object Types

Heap-allocated objects:

```go
type ObjectType byte

const (
    OBJ_STRING   ObjectType = 0x01
    OBJ_ARRAY    ObjectType = 0x02
    OBJ_OBJECT   ObjectType = 0x03  // Plain object literal
    OBJ_FUNCTION ObjectType = 0x04  // Function object
    OBJ_CLOSURE  ObjectType = 0x05  // Closure (function + upvalues)
    OBJ_CLASS    ObjectType = 0x06  // Class definition
    OBJ_INSTANCE ObjectType = 0x07  // Class instance
    OBJ_UPVALUE  ObjectType = 0x08  // Captured variable
)
```

### 2.3 Object Structures

```go
// Base object header
type Object struct {
    Type   ObjectType
    Marked bool       // For GC
    Next   *Object    // Intrusive list for GC
}

// String object
type ObjString struct {
    Object
    Value string
    Hash  uint32  // Cached hash for interning
}

// Array object
type ObjArray struct {
    Object
    Elements []Value
}

// Plain object (object literal)
type ObjObject struct {
    Object
    Fields map[string]Value
}

// Compiled function
type ObjFunction struct {
    Object
    Name       string
    Arity      int        // Parameter count
    Chunk      *Chunk     // Bytecode
    UpvalueCount int      // Number of upvalues
}

// Closure (function + captured environment)
type ObjClosure struct {
    Object
    Function *ObjFunction
    Upvalues []*ObjUpvalue
}

// Upvalue (captured variable)
type ObjUpvalue struct {
    Object
    Location *Value  // Points to stack slot or Closed
    Closed   Value   // Holds value after variable goes out of scope
    Next     *ObjUpvalue  // Linked list of open upvalues
}

// Class definition
type ObjClass struct {
    Object
    Name    string
    Super   *ObjClass           // Parent class (or nil)
    Methods map[string]*ObjClosure
    Fields  []string            // Field names for initialization
}

// Class instance
type ObjInstance struct {
    Object
    Class  *ObjClass
    Fields map[string]Value
}
```

---

## 3. Bytecode Format

### 3.1 Chunk Structure

A chunk holds bytecode and associated data:

```go
type Chunk struct {
    Code      []byte    // Bytecode instructions
    Constants []Value   // Constant pool
    Lines     []int     // Line numbers for debugging (parallel to Code)
}
```

### 3.2 Module Structure

A compiled module:

```go
type Module struct {
    Name      string
    Functions []*ObjFunction  // All functions (index 0 = top-level)
    Classes   []*ObjClass     // All class definitions
    Globals   []string        // Global variable names
}
```

---

## 4. Instruction Set

### 4.1 Instruction Format

Instructions are 1-byte opcodes, optionally followed by operands:

| Format | Description |
|--------|-------------|
| `OP` | Single byte, no operands |
| `OP u8` | Opcode + 1-byte operand |
| `OP u16` | Opcode + 2-byte operand (big-endian) |
| `OP u8 u8` | Opcode + two 1-byte operands |

### 4.2 Opcodes

```go
type OpCode byte

const (
    // ============ Constants & Literals ============
    OP_CONSTANT      OpCode = 0x01  // u16 index -> push constants[index]
    OP_NULL          OpCode = 0x02  // push null
    OP_TRUE          OpCode = 0x03  // push true
    OP_FALSE         OpCode = 0x04  // push false

    // ============ Arithmetic ============
    OP_ADD           OpCode = 0x10  // pop b, pop a, push a + b
    OP_SUBTRACT      OpCode = 0x11  // pop b, pop a, push a - b
    OP_MULTIPLY      OpCode = 0x12  // pop b, pop a, push a * b
    OP_DIVIDE        OpCode = 0x13  // pop b, pop a, push a / b
    OP_MODULO        OpCode = 0x14  // pop b, pop a, push a % b
    OP_NEGATE        OpCode = 0x15  // pop a, push -a

    // ============ Comparison ============
    OP_EQUAL         OpCode = 0x20  // pop b, pop a, push a == b
    OP_NOT_EQUAL     OpCode = 0x21  // pop b, pop a, push a != b
    OP_LESS          OpCode = 0x22  // pop b, pop a, push a < b
    OP_LESS_EQUAL    OpCode = 0x23  // pop b, pop a, push a <= b
    OP_GREATER       OpCode = 0x24  // pop b, pop a, push a > b
    OP_GREATER_EQUAL OpCode = 0x25  // pop b, pop a, push a >= b

    // ============ Logical ============
    OP_NOT           OpCode = 0x30  // pop a, push !a

    // ============ String ============
    OP_CONCAT        OpCode = 0x40  // pop b, pop a, push a + b (strings)

    // ============ Variables ============
    OP_GET_LOCAL     OpCode = 0x50  // u8 slot -> push stack[frame.base + slot]
    OP_SET_LOCAL     OpCode = 0x51  // u8 slot -> stack[frame.base + slot] = peek()
    OP_GET_GLOBAL    OpCode = 0x52  // u16 index -> push globals[constants[index]]
    OP_SET_GLOBAL    OpCode = 0x53  // u16 index -> globals[constants[index]] = peek()
    OP_GET_UPVALUE   OpCode = 0x54  // u8 index -> push closure.upvalues[index]
    OP_SET_UPVALUE   OpCode = 0x55  // u8 index -> closure.upvalues[index] = peek()

    // ============ Stack Operations ============
    OP_POP           OpCode = 0x60  // discard top of stack
    OP_POPN          OpCode = 0x61  // u8 n -> discard n values from stack
    OP_DUP           OpCode = 0x62  // duplicate top of stack

    // ============ Control Flow ============
    OP_JUMP          OpCode = 0x70  // u16 offset -> ip += offset
    OP_JUMP_BACK     OpCode = 0x71  // u16 offset -> ip -= offset
    OP_JUMP_IF_FALSE OpCode = 0x72  // u16 offset -> if !pop() then ip += offset
    OP_JUMP_IF_TRUE  OpCode = 0x73  // u16 offset -> if pop() then ip += offset

    // ============ Functions & Calls ============
    OP_CALL          OpCode = 0x80  // u8 argCount -> call function with args
    OP_RETURN        OpCode = 0x81  // return from function
    OP_CLOSURE       OpCode = 0x82  // u16 funcIndex, [u8 isLocal, u8 index]* -> create closure

    // ============ Classes & Objects ============
    OP_CLASS         OpCode = 0x90  // u16 classIndex -> push class
    OP_GET_PROPERTY  OpCode = 0x91  // u16 nameIndex -> pop obj, push obj.property
    OP_SET_PROPERTY  OpCode = 0x92  // u16 nameIndex -> pop val, pop obj, obj.property = val, push val
    OP_METHOD        OpCode = 0x93  // u16 nameIndex -> pop closure, add method to class at stack top
    OP_INVOKE        OpCode = 0x94  // u16 nameIndex, u8 argCount -> invoke method directly
    OP_INHERIT       OpCode = 0x95  // pop super, pop sub, sub inherits from super
    OP_GET_SUPER     OpCode = 0x96  // u16 nameIndex -> lookup method in superclass
    OP_SUPER_INVOKE  OpCode = 0x97  // u16 nameIndex, u8 argCount -> invoke super method

    // ============ Arrays ============
    OP_ARRAY         OpCode = 0xA0  // u16 count -> pop count values, push array
    OP_GET_INDEX     OpCode = 0xA1  // pop index, pop array, push array[index]
    OP_SET_INDEX     OpCode = 0xA2  // pop val, pop index, pop array, array[index] = val, push val

    // ============ Objects (literals) ============
    OP_OBJECT        OpCode = 0xB0  // u16 count -> pop count key-value pairs, push object

    // ============ Special ============
    OP_CLOSE_UPVALUE OpCode = 0xC0  // close upvalue at stack top
    OP_PRINT         OpCode = 0xD0  // pop value, print it (built-in)
    OP_PRINTLN       OpCode = 0xD1  // pop value, print it with newline

    // ============ Built-in Functions ============
    OP_BUILTIN       OpCode = 0xE0  // u8 builtinId, u8 argCount -> call built-in
)
```

### 4.3 Built-in Function IDs

```go
const (
    BUILTIN_LEN       = 0x01
    BUILTIN_TOSTRING  = 0x02
    BUILTIN_TONUMBER  = 0x03
    BUILTIN_PUSH      = 0x04
    BUILTIN_POP       = 0x05
    BUILTIN_SQRT      = 0x06
    BUILTIN_FLOOR     = 0x07
    BUILTIN_CEIL      = 0x08
    BUILTIN_ABS       = 0x09
)
```

---

## 5. Call Frame

### 5.1 Frame Structure

```go
type CallFrame struct {
    Closure   *ObjClosure  // Currently executing closure
    IP        int          // Instruction pointer (index into chunk.Code)
    BasePtr   int          // Base of this frame's stack window
}
```

### 5.2 Stack Layout During Call

```
Before call foo(a, b):
┌─────────────┐
│     b       │  <- sp
│     a       │
│   closure   │  <- function being called
│    ...      │
└─────────────┘

After call setup:
┌─────────────┐
│  (locals)   │  <- sp (grows as locals declared)
│     b       │  <- slot 2 (param)
│     a       │  <- slot 1 (param)
│   closure   │  <- slot 0 (for 'this' or just reserved)
│    ...      │  <- previous frame
└─────────────┘
     ^
     basePtr
```

---

## 6. Upvalues & Closures

### 6.1 Upvalue Mechanism

Closures capture variables from enclosing scopes using upvalues:

1. **Open upvalue**: Points to a stack slot (variable still in scope)
2. **Closed upvalue**: Holds the value itself (variable went out of scope)

### 6.2 OP_CLOSURE Instruction

Format: `OP_CLOSURE u16:funcIndex [u8:isLocal u8:index]*`

For each upvalue in the function:
- `isLocal=1`: Capture from current frame's local at `index`
- `isLocal=0`: Capture from enclosing closure's upvalue at `index`

### 6.3 Closing Upvalues

When a local variable goes out of scope, `OP_CLOSE_UPVALUE` moves its value from the stack into the upvalue's `Closed` field.

```
// Before closing (open upvalue)
ObjUpvalue {
    Location: &stack[slot]  // Points to stack
    Closed: (unused)
}

// After closing (closed upvalue)
ObjUpvalue {
    Location: &self.Closed  // Points to own Closed field
    Closed: capturedValue
}
```

---

## 7. Method Dispatch

### 7.1 Instance Method Call

`instance.method(args)` compiles to:

```
... push instance ...
... push args ...
OP_INVOKE nameIndex argCount
```

OP_INVOKE:
1. Peek instance from stack (below args)
2. Look up method in instance's class
3. Create call frame with instance as slot 0 (`this`)

### 7.2 Super Call

`super.method(args)` in a method compiles to:

```
OP_GET_LOCAL 0          // push 'this'
... push args ...
OP_GET_SUPER nameIndex  // lookup in superclass
OP_CALL argCount
```

Or optimized:
```
OP_GET_LOCAL 0          // push 'this'
... push args ...
OP_SUPER_INVOKE nameIndex argCount
```

---

## 8. VM State

### 8.1 VM Structure

```go
type VM struct {
    // Execution state
    Frames     [MAX_FRAMES]CallFrame
    FrameCount int
    Stack      [MAX_STACK]Value
    StackTop   int

    // Global state
    Globals    map[string]Value

    // Object management
    Objects    *Object       // Head of all objects list
    OpenUpvalues *ObjUpvalue // Head of open upvalues list

    // GC state
    BytesAllocated int
    NextGC         int
    GrayStack      []*Object  // For tri-color marking
}

const (
    MAX_FRAMES = 256
    MAX_STACK  = 65536
)
```

### 8.2 Execution Loop

```go
func (vm *VM) Run() error {
    frame := &vm.Frames[vm.FrameCount-1]

    for {
        instruction := frame.ReadByte()

        switch OpCode(instruction) {
        case OP_CONSTANT:
            index := frame.ReadU16()
            vm.Push(frame.Chunk().Constants[index])

        case OP_ADD:
            b := vm.Pop()
            a := vm.Pop()
            if a.IsNumber() && b.IsNumber() {
                vm.Push(NumberValue(a.AsNumber() + b.AsNumber()))
            } else if a.IsString() && b.IsString() {
                vm.Push(StringValue(a.AsString() + b.AsString()))
            } else {
                return vm.runtimeError("operands must be numbers or strings")
            }

        case OP_CALL:
            argCount := frame.ReadByte()
            if err := vm.callValue(vm.Peek(argCount), argCount); err != nil {
                return err
            }
            frame = &vm.Frames[vm.FrameCount-1]

        case OP_RETURN:
            result := vm.Pop()
            vm.closeUpvalues(frame.BasePtr)
            vm.FrameCount--
            if vm.FrameCount == 0 {
                return nil  // Program complete
            }
            vm.StackTop = frame.BasePtr
            vm.Push(result)
            frame = &vm.Frames[vm.FrameCount-1]

        // ... other cases ...
        }
    }
}
```

---

## 9. Garbage Collection

### 9.1 Algorithm

Tri-color mark-and-sweep:

1. **Mark roots**: Stack, globals, call frames, open upvalues
2. **Trace**: Follow references from gray objects to white objects
3. **Sweep**: Free all white (unreachable) objects

### 9.2 GC Triggers

GC runs when `BytesAllocated > NextGC`. After GC:
```go
vm.NextGC = vm.BytesAllocated * GC_HEAP_GROW_FACTOR
```

### 9.3 Write Barriers

Not needed for simple mark-and-sweep (only runs when VM is paused).

---

## 10. Binary Format

### 10.1 File Structure

```
GoTS Bytecode File (.gtsb)
┌──────────────────────────────────────┐
│ Magic: "GOTS" (4 bytes)              │
│ Version: u16                         │
├──────────────────────────────────────┤
│ Constant Pool                        │
│   count: u16                         │
│   [constants...]                     │
├──────────────────────────────────────┤
│ Global Names                         │
│   count: u16                         │
│   [strings...]                       │
├──────────────────────────────────────┤
│ Classes                              │
│   count: u16                         │
│   [class definitions...]             │
├──────────────────────────────────────┤
│ Functions                            │
│   count: u16                         │
│   [function definitions...]          │
│   (index 0 = top-level/main)         │
└──────────────────────────────────────┘
```

### 10.2 Constant Encoding

```
Constant:
  tag: u8
  data: (depends on tag)

Tags:
  0x01 = Null
  0x02 = Boolean (u8: 0 or 1)
  0x03 = Number (f64: 8 bytes, IEEE 754)
  0x04 = String (u16 length, UTF-8 bytes)
```

### 10.3 Function Encoding

```
Function:
  nameLength: u16
  name: [u8; nameLength]
  arity: u8
  upvalueCount: u8
  codeLength: u32
  code: [u8; codeLength]
  constantCount: u16
  constants: [Constant; constantCount]
  lineCount: u32
  lines: [u16; lineCount]  // Parallel to code, RLE encoded
```

### 10.4 Class Encoding

```
Class:
  nameLength: u16
  name: [u8; nameLength]
  superIndex: i16  // -1 if no superclass, else index into classes
  fieldCount: u16
  fields: [String; fieldCount]
  methodCount: u16
  methods: [
    nameLength: u16
    name: [u8; nameLength]
    functionIndex: u16
  ; methodCount]
```

---

## 11. Instruction Reference

### 11.1 Constants & Literals

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_CONSTANT` | u16 idx | → value | Push constant from pool |
| `OP_NULL` | - | → null | Push null |
| `OP_TRUE` | - | → true | Push true |
| `OP_FALSE` | - | → false | Push false |

### 11.2 Arithmetic

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_ADD` | - | a, b → result | Add (or concat strings) |
| `OP_SUBTRACT` | - | a, b → result | Subtract |
| `OP_MULTIPLY` | - | a, b → result | Multiply |
| `OP_DIVIDE` | - | a, b → result | Divide |
| `OP_MODULO` | - | a, b → result | Modulo |
| `OP_NEGATE` | - | a → result | Negate |

### 11.3 Comparison

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_EQUAL` | - | a, b → bool | Equal |
| `OP_NOT_EQUAL` | - | a, b → bool | Not equal |
| `OP_LESS` | - | a, b → bool | Less than |
| `OP_LESS_EQUAL` | - | a, b → bool | Less or equal |
| `OP_GREATER` | - | a, b → bool | Greater than |
| `OP_GREATER_EQUAL` | - | a, b → bool | Greater or equal |

### 11.4 Logical

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_NOT` | - | a → bool | Logical not |

### 11.5 Variables

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_GET_LOCAL` | u8 slot | → value | Get local variable |
| `OP_SET_LOCAL` | u8 slot | value → value | Set local variable |
| `OP_GET_GLOBAL` | u16 idx | → value | Get global variable |
| `OP_SET_GLOBAL` | u16 idx | value → value | Set global variable |
| `OP_GET_UPVALUE` | u8 idx | → value | Get captured variable |
| `OP_SET_UPVALUE` | u8 idx | value → value | Set captured variable |

### 11.6 Stack

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_POP` | - | value → | Discard top |
| `OP_POPN` | u8 n | n values → | Discard n values |
| `OP_DUP` | - | value → value, value | Duplicate top |

### 11.7 Control Flow

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_JUMP` | u16 off | - | Jump forward |
| `OP_JUMP_BACK` | u16 off | - | Jump backward |
| `OP_JUMP_IF_FALSE` | u16 off | cond → | Jump if false |
| `OP_JUMP_IF_TRUE` | u16 off | cond → | Jump if true |

### 11.8 Functions

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_CALL` | u8 argc | fn, args → result | Call function |
| `OP_RETURN` | - | result → (frame destroyed) | Return from function |
| `OP_CLOSURE` | u16 idx, upvals... | → closure | Create closure |

### 11.9 Classes & Objects

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_CLASS` | u16 idx | → class | Push class |
| `OP_GET_PROPERTY` | u16 name | obj → value | Get property |
| `OP_SET_PROPERTY` | u16 name | obj, val → val | Set property |
| `OP_METHOD` | u16 name | class, closure → class | Define method |
| `OP_INVOKE` | u16 name, u8 argc | obj, args → result | Invoke method |
| `OP_INHERIT` | - | sub, super → sub | Set up inheritance |
| `OP_GET_SUPER` | u16 name | instance → method | Get super method |
| `OP_SUPER_INVOKE` | u16 name, u8 argc | this, args → result | Invoke super method |

### 11.10 Arrays

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_ARRAY` | u16 count | values → array | Create array |
| `OP_GET_INDEX` | - | arr, idx → value | Get array element |
| `OP_SET_INDEX` | - | arr, idx, val → val | Set array element |

### 11.11 Object Literals

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_OBJECT` | u16 count | keys, values → object | Create object |

### 11.12 Special

| Opcode | Operands | Stack Effect | Description |
|--------|----------|--------------|-------------|
| `OP_CLOSE_UPVALUE` | - | value → | Close upvalue at top |
| `OP_PRINT` | - | value → | Print value |
| `OP_PRINTLN` | - | value → | Print with newline |
| `OP_BUILTIN` | u8 id, u8 argc | args → result | Call built-in |

---

## 12. Compilation Examples

### 12.1 Simple Expression

Source:
```typescript
let x: number = 1 + 2 * 3;
```

Bytecode:
```
OP_CONSTANT    0      // 1
OP_CONSTANT    1      // 2
OP_CONSTANT    2      // 3
OP_MULTIPLY           // 2 * 3 = 6
OP_ADD                // 1 + 6 = 7
OP_SET_LOCAL   0      // x = 7
OP_POP
```

### 12.2 If Statement

Source:
```typescript
if (x > 0) {
    println("positive");
} else {
    println("non-positive");
}
```

Bytecode:
```
OP_GET_LOCAL   0           // x
OP_CONSTANT    0           // 0
OP_GREATER                 // x > 0
OP_JUMP_IF_FALSE +12       // jump to else
OP_CONSTANT    1           // "positive"
OP_PRINTLN
OP_JUMP        +8          // jump past else
OP_CONSTANT    2           // "non-positive"  <- else
OP_PRINTLN
                           // <- end
```

### 12.3 While Loop

Source:
```typescript
while (i < 10) {
    i = i + 1;
}
```

Bytecode:
```
         <- loop start
OP_GET_LOCAL   0           // i
OP_CONSTANT    0           // 10
OP_LESS                    // i < 10
OP_JUMP_IF_FALSE +12       // exit loop
OP_GET_LOCAL   0           // i
OP_CONSTANT    1           // 1
OP_ADD                     // i + 1
OP_SET_LOCAL   0           // i = ...
OP_POP
OP_JUMP_BACK   -18         // back to loop start
                           // <- loop end
```

### 12.4 Function Call

Source:
```typescript
function add(a: number, b: number): number {
    return a + b;
}
let result: number = add(1, 2);
```

Bytecode (add function):
```
OP_GET_LOCAL   1      // a (slot 0 is reserved)
OP_GET_LOCAL   2      // b
OP_ADD
OP_RETURN
```

Bytecode (call site):
```
OP_CLOSURE     0      // create closure for add
OP_SET_GLOBAL  0      // store in global "add"
OP_POP
OP_GET_GLOBAL  0      // get "add"
OP_CONSTANT    0      // 1
OP_CONSTANT    1      // 2
OP_CALL        2      // call with 2 args
OP_SET_LOCAL   0      // result = ...
OP_POP
```

### 12.5 Closure

Source:
```typescript
function makeCounter(): () => number {
    let count: number = 0;
    return function(): number {
        count = count + 1;
        return count;
    };
}
```

Bytecode (outer function):
```
OP_CONSTANT      0           // 0
                             // count is local slot 1
OP_CLOSURE       1  1 1      // inner func, 1 upvalue, local slot 1
OP_RETURN
```

Bytecode (inner function):
```
OP_GET_UPVALUE   0           // count
OP_CONSTANT      0           // 1
OP_ADD
OP_SET_UPVALUE   0           // count = count + 1
OP_GET_UPVALUE   0
OP_RETURN
```

### 12.6 Class & Method

Source:
```typescript
class Point {
    x: number;
    y: number;

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
    }

    toString(): string {
        return "(" + toString(this.x) + ", " + toString(this.y) + ")";
    }
}

let p: Point = new Point(3, 4);
println(p.toString());
```

Bytecode (simplified):
```
// Class definition
OP_CLASS         0           // Point class
OP_CLOSURE       0           // constructor closure
OP_METHOD        0           // "constructor"
OP_CLOSURE       1           // toString closure
OP_METHOD        1           // "toString"
OP_SET_GLOBAL    0           // store class in global
OP_POP

// Instantiation: new Point(3, 4)
OP_GET_GLOBAL    0           // Point class
OP_CONSTANT      0           // 3
OP_CONSTANT      1           // 4
OP_CALL          2           // creates instance, calls constructor
OP_SET_LOCAL     0           // p = ...
OP_POP

// Method call: p.toString()
OP_GET_LOCAL     0           // p
OP_INVOKE        1  0        // invoke "toString" with 0 args
OP_PRINTLN
```

---

## 13. Error Handling

### 13.1 Runtime Errors

The VM produces runtime errors for:

- Type errors (wrong operand types)
- Null pointer dereference
- Array index out of bounds
- Stack overflow
- Division by zero
- Undefined variable access
- Property access on non-object
- Method not found

### 13.2 Error Format

```go
type RuntimeError struct {
    Message string
    Line    int
    Stack   []StackTraceEntry
}

type StackTraceEntry struct {
    Function string
    Line     int
}
```

### 13.3 Stack Trace

On error, unwind call frames to build stack trace:

```
RuntimeError: Cannot read property 'x' of null
  at getX (example.gts:15)
  at calculate (example.gts:23)
  at <main> (example.gts:30)
```

---

## 14. Summary

| Component | Description |
|-----------|-------------|
| **Values** | Tagged union: null, bool, number, object pointer |
| **Objects** | GC-managed: string, array, object, closure, class, instance |
| **Bytecode** | Variable-length instructions, constant pool |
| **Execution** | Stack-based, call frames, upvalues for closures |
| **GC** | Tri-color mark-and-sweep |
| **Binary** | Magic header, constants, globals, classes, functions |

This design balances simplicity with the features needed for GoTS v1.
