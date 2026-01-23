// Package vm implements garbage collection for the GoTS VM.
package vm

// GC implements a simple mark-and-sweep garbage collector.
type GC struct {
	objects       []Object // All allocated objects
	bytesAllocated int     // Current bytes allocated
	nextGC         int     // Threshold for next GC
	vm             *VM     // Reference to the VM for marking roots
}

const (
	GC_HEAP_GROW_FACTOR = 2
	GC_INITIAL_THRESHOLD = 1024 * 1024 // 1 MB
)

// NewGC creates a new garbage collector.
func NewGC(vm *VM) *GC {
	return &GC{
		objects:       make([]Object, 0),
		bytesAllocated: 0,
		nextGC:         GC_INITIAL_THRESHOLD,
		vm:             vm,
	}
}

// Track registers an object with the GC.
func (gc *GC) Track(obj Object) {
	gc.objects = append(gc.objects, obj)
	gc.bytesAllocated += gc.objectSize(obj)

	if gc.bytesAllocated > gc.nextGC {
		gc.Collect()
	}
}

// objectSize estimates the memory size of an object.
func (gc *GC) objectSize(obj Object) int {
	switch o := obj.(type) {
	case *ObjString:
		return 24 + len(o.Value)
	case *ObjArray:
		return 24 + len(o.Elements)*16
	case *ObjObject:
		return 24 + len(o.Fields)*32
	case *ObjFunction:
		return 48 + len(o.Chunk.Code) + len(o.Chunk.Constants)*16
	case *ObjClosure:
		return 24 + len(o.Upvalues)*8
	case *ObjUpvalue:
		return 40
	case *ObjClass:
		return 24 + len(o.Methods)*32
	case *ObjInstance:
		return 24 + len(o.Fields)*32
	case *ObjBoundMethod:
		return 24
	default:
		return 16
	}
}

// Collect performs a garbage collection cycle.
func (gc *GC) Collect() {
	beforeSize := gc.bytesAllocated

	// Mark phase
	gc.markRoots()

	// Sweep phase
	gc.sweep()

	// Adjust threshold
	gc.nextGC = gc.bytesAllocated * GC_HEAP_GROW_FACTOR
	if gc.nextGC < GC_INITIAL_THRESHOLD {
		gc.nextGC = GC_INITIAL_THRESHOLD
	}

	_ = beforeSize // Could be used for debug logging
}

// markRoots marks all root objects.
func (gc *GC) markRoots() {
	// Mark values on the stack
	for i := 0; i < gc.vm.sp; i++ {
		gc.markValue(gc.vm.stack[i])
	}

	// Mark global variables
	for _, v := range gc.vm.globals {
		gc.markValue(v)
	}

	// Mark call frame closures
	for i := 0; i < gc.vm.frameCount; i++ {
		gc.markObject(gc.vm.frames[i].closure)
	}

	// Mark open upvalues
	for uv := gc.vm.openUpvalues; uv != nil; uv = uv.Next {
		gc.markObject(uv)
	}
}

// markValue marks a value if it contains an object.
func (gc *GC) markValue(v Value) {
	if v.IsObject() && v.obj != nil {
		gc.markObject(v.obj)
	}
}

// Markable objects have a marked flag.
type Markable interface {
	IsMarked() bool
	SetMarked(bool)
}

// markObject marks an object and traces its references.
func (gc *GC) markObject(obj Object) {
	if obj == nil {
		return
	}

	// Check if already marked using type assertion
	if m, ok := obj.(Markable); ok && m.IsMarked() {
		return
	}

	// Mark the object
	if m, ok := obj.(Markable); ok {
		m.SetMarked(true)
	}

	// Trace references based on object type
	switch o := obj.(type) {
	case *ObjArray:
		for _, elem := range o.Elements {
			gc.markValue(elem)
		}
	case *ObjObject:
		for _, v := range o.Fields {
			gc.markValue(v)
		}
	case *ObjClosure:
		gc.markObject(o.Function)
		for _, uv := range o.Upvalues {
			gc.markObject(uv)
		}
	case *ObjUpvalue:
		gc.markValue(*o.Location)
	case *ObjClass:
		for _, m := range o.Methods {
			gc.markObject(m)
		}
		if o.Super != nil {
			gc.markObject(o.Super)
		}
	case *ObjInstance:
		gc.markObject(o.Class)
		for _, v := range o.Fields {
			gc.markValue(v)
		}
	case *ObjBoundMethod:
		gc.markValue(o.Receiver)
		gc.markObject(o.Method)
	}
}

// sweep removes unmarked objects.
func (gc *GC) sweep() {
	// Remove unmarked objects
	newObjects := make([]Object, 0, len(gc.objects))
	gc.bytesAllocated = 0

	for _, obj := range gc.objects {
		if m, ok := obj.(Markable); ok && m.IsMarked() {
			// Unmark for next cycle
			m.SetMarked(false)
			newObjects = append(newObjects, obj)
			gc.bytesAllocated += gc.objectSize(obj)
		}
		// Unmarked objects are implicitly freed by Go's GC
	}

	gc.objects = newObjects
}

// Stats returns GC statistics.
func (gc *GC) Stats() (allocated, threshold int) {
	return gc.bytesAllocated, gc.nextGC
}
