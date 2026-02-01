// Package codegen generates Go source code from a typed AST.
package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/typed"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Generator transforms typed AST into Go source code.
type Generator struct {
	buf             *bytes.Buffer
	indent          int
	imports         map[string]bool
	goImportedNames map[string]string // maps imported name to its package
	currentRetType  types.Type        // Track return type for type assertions
	currentClass    *typed.ClassDecl  // Track current class for super() handling
}

// Generate produces Go source code from a typed program.
func Generate(prog *typed.Program) ([]byte, error) {
	g := &Generator{
		buf:             new(bytes.Buffer),
		imports:         make(map[string]bool),
		goImportedNames: make(map[string]string),
	}

	// Build map of imported names to their packages
	for _, imp := range prog.GoImports {
		for _, name := range imp.Names {
			g.goImportedNames[name] = imp.Package
		}
	}

	// First pass: collect required imports
	g.collectImports(prog)

	// Generate the code
	g.genProgram(prog)

	// Format the output
	src := g.buf.Bytes()
	formatted, err := format.Source(src)
	if err != nil {
		// Return unformatted source for debugging
		return src, fmt.Errorf("format error: %v\n%s", err, src)
	}

	return formatted, nil
}

func (g *Generator) collectImports(prog *typed.Program) {
	// fmt is always needed for print
	g.imports["fmt"] = true
	// reflect is needed for dynamic function calls
	g.imports["reflect"] = true
	// sync is needed for event loop and Promise
	g.imports["sync"] = true
	// time is needed for event loop timers
	g.imports["time"] = true

	// Add Go package imports from the source
	for _, imp := range prog.GoImports {
		g.imports[imp.Package] = true
	}

	// Check for math functions
	for _, fn := range prog.Functions {
		g.collectImportsFromBlock(fn.Body)
	}
	for _, stmt := range prog.TopLevel {
		g.collectImportsFromStmt(stmt)
	}
	for _, class := range prog.Classes {
		if class.Constructor != nil {
			g.collectImportsFromBlock(class.Constructor.Body)
		}
		for _, method := range class.Methods {
			g.collectImportsFromBlock(method.Body)
		}
	}
}

func (g *Generator) collectImportsFromBlock(block *typed.BlockStmt) {
	if block == nil {
		return
	}
	for _, stmt := range block.Stmts {
		g.collectImportsFromStmt(stmt)
	}
}

func (g *Generator) collectImportsFromStmt(stmt typed.Stmt) {
	switch s := stmt.(type) {
	case *typed.ExprStmt:
		g.collectImportsFromExpr(s.Expr)
	case *typed.VarDecl:
		if s.Init != nil {
			g.collectImportsFromExpr(s.Init)
		}
	case *typed.BlockStmt:
		g.collectImportsFromBlock(s)
	case *typed.IfStmt:
		g.collectImportsFromExpr(s.Condition)
		g.collectImportsFromBlock(s.Then)
		if s.Else != nil {
			g.collectImportsFromStmt(s.Else)
		}
	case *typed.WhileStmt:
		g.collectImportsFromExpr(s.Condition)
		g.collectImportsFromBlock(s.Body)
	case *typed.ForStmt:
		if s.Condition != nil {
			g.collectImportsFromExpr(s.Condition)
		}
		if s.Update != nil {
			g.collectImportsFromExpr(s.Update)
		}
		g.collectImportsFromBlock(s.Body)
	case *typed.ForOfStmt:
		g.collectImportsFromExpr(s.Iterable)
		g.collectImportsFromBlock(s.Body)
	case *typed.SwitchStmt:
		g.collectImportsFromExpr(s.Discriminant)
		for _, c := range s.Cases {
			if c.Test != nil {
				g.collectImportsFromExpr(c.Test)
			}
			for _, cs := range c.Stmts {
				g.collectImportsFromStmt(cs)
			}
		}
	case *typed.ReturnStmt:
		if s.Value != nil {
			g.collectImportsFromExpr(s.Value)
		}
	case *typed.TryStmt:
		g.collectImportsFromBlock(s.TryBlock)
		g.collectImportsFromBlock(s.CatchBlock)
	case *typed.ThrowStmt:
		g.collectImportsFromExpr(s.Value)
	case *typed.FuncDecl:
		g.collectImportsFromBlock(s.Body)
	case *typed.ClassDecl:
		if s.Constructor != nil {
			g.collectImportsFromBlock(s.Constructor.Body)
		}
		for _, m := range s.Methods {
			g.collectImportsFromBlock(m.Body)
		}
	}
}

func (g *Generator) collectImportsFromExpr(expr typed.Expr) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *typed.BuiltinCall:
		switch e.Name {
		case "sqrt", "floor", "ceil", "abs", "isNaN", "isFinite":
			g.imports["math"] = true
		case "split", "join", "replace", "trim", "startsWith", "endsWith", "includes":
			g.imports["strings"] = true
		case "parseFloat":
			g.imports["strconv"] = true
		}
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.BuiltinObjectCall:
		// Add imports for built-in object method calls (e.g., Math.round)
		for _, imp := range typed.GetBuiltinImports(e.Object) {
			g.imports[imp] = true
		}
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.BuiltinObjectConstant:
		// Add imports for built-in object constants (e.g., Math.PI)
		for _, imp := range typed.GetBuiltinImports(e.Object) {
			g.imports[imp] = true
		}
	case *typed.DateNewExpr:
		g.imports["time"] = true
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.DateMethodCall:
		g.imports["time"] = true
		g.collectImportsFromExpr(e.Object)
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.BinaryExpr:
		// Modulo for float64 uses math.Mod
		if e.Op == "%" {
			leftType := types.Unwrap(e.Left.Type())
			if lp, ok := leftType.(*types.Primitive); ok && lp.Kind == types.KindFloat {
				g.imports["math"] = true
			}
		}
		g.collectImportsFromExpr(e.Left)
		g.collectImportsFromExpr(e.Right)
	case *typed.UnaryExpr:
		g.collectImportsFromExpr(e.Operand)
	case *typed.CallExpr:
		g.collectImportsFromExpr(e.Callee)
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.IndexExpr:
		g.collectImportsFromExpr(e.Object)
		g.collectImportsFromExpr(e.Index)
	case *typed.PropertyExpr:
		g.collectImportsFromExpr(e.Object)
	case *typed.ArrayLit:
		for _, elem := range e.Elements {
			g.collectImportsFromExpr(elem)
		}
	case *typed.ObjectLit:
		for _, prop := range e.Properties {
			g.collectImportsFromExpr(prop.Value)
		}
	case *typed.FuncExpr:
		g.collectImportsFromBlock(e.Body)
		if e.BodyExpr != nil {
			g.collectImportsFromExpr(e.BodyExpr)
		}
	case *typed.NewExpr:
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.SuperExpr:
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
	case *typed.AssignExpr:
		g.collectImportsFromExpr(e.Target)
		g.collectImportsFromExpr(e.Value)
	case *typed.CompoundAssignExpr:
		g.collectImportsFromExpr(e.Target)
		g.collectImportsFromExpr(e.Value)
	case *typed.UpdateExpr:
		g.collectImportsFromExpr(e.Operand)
	case *typed.MapLit:
		for _, entry := range e.Entries {
			g.collectImportsFromExpr(entry.Key)
			g.collectImportsFromExpr(entry.Value)
		}
	case *typed.SetLit:
		// Empty set literal - nothing to collect
	case *typed.RegexLit:
		g.imports["regexp"] = true
	case *typed.MethodCallExpr:
		g.collectImportsFromExpr(e.Object)
		for _, arg := range e.Args {
			g.collectImportsFromExpr(arg)
		}
		// Check if this is a string method call that requires strings import
		objType := types.Unwrap(e.Object.Type())
		if prim, ok := objType.(*types.Primitive); ok && prim.Kind == types.KindString {
			g.imports["strings"] = true
		}
	case *typed.PromiseMethodCall:
		g.collectImportsFromExpr(e.Object)
		if e.Callback != nil {
			g.collectImportsFromExpr(e.Callback)
		}
	}
}

func (g *Generator) genProgram(prog *typed.Program) {
	g.writeln("package main")
	g.writeln("")

	// Write imports
	if len(g.imports) > 0 {
		if len(g.imports) == 1 {
			for pkg := range g.imports {
				g.writeln("import %q", pkg)
			}
		} else {
			g.writeln("import (")
			g.indent++
			for pkg := range g.imports {
				g.writeln("%q", pkg)
			}
			g.indent--
			g.writeln(")")
		}
		g.writeln("")
	}

	// Generate runtime helpers
	g.genRuntime()
	g.writeln("")

	// Generate enums
	for _, enum := range prog.Enums {
		g.genEnum(enum)
		g.writeln("")
	}

	// Generate type aliases
	for _, alias := range prog.TypeAliases {
		g.genTypeAlias(alias)
		g.writeln("")
	}

	// Generate interfaces
	for _, iface := range prog.Interfaces {
		g.genInterface(iface)
		g.writeln("")
	}

	// Generate classes as structs
	for _, class := range prog.Classes {
		g.genClass(class)
		g.writeln("")
	}

	// Generate top-level functions
	for _, fn := range prog.Functions {
		g.genFuncDecl(fn)
		g.writeln("")
	}

	// Generate main function
	g.writeln("func main() {")
	g.indent++
	// Initialize event loop
	g.writeln("// Initialize event loop")
	g.writeln("gts_eventLoop = &GTS_EventLoop{")
	g.indent++
	g.writeln("microtasks:    make([]func(), 0),")
	g.writeln("macrotasks:    make([]func(), 0),")
	g.writeln("timers:        make(map[int64]*time.Timer),")
	g.writeln("pendingTimers: make(chan func(), 1000),")
	g.indent--
	g.writeln("}")
	g.writeln("")
	// Execute top-level code
	for _, stmt := range prog.TopLevel {
		g.genStmt(stmt)
	}
	g.writeln("")
	// Run event loop
	g.writeln("// Run event loop until all tasks complete")
	g.writeln("gts_runEventLoop()")
	g.indent--
	g.writeln("}")
}

func (g *Generator) genRuntime() {
	// Generate runtime helper functions
	g.writeln("// Runtime helpers")
	g.writeln("")

	// len helper - returns int
	g.writeln("func gts_len(v interface{}) int {")
	g.indent++
	g.writeln("switch x := v.(type) {")
	g.writeln("case string:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("case []interface{}:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("case []int:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("case []float64:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("case []string:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("case []bool:")
	g.indent++
	g.writeln("return len(x)")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("return 0")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// typeof helper
	g.writeln("func gts_typeof(v interface{}) string {")
	g.indent++
	g.writeln("if v == nil {")
	g.indent++
	g.writeln("return \"null\"")
	g.indent--
	g.writeln("}")
	g.writeln("switch v.(type) {")
	g.writeln("case int:")
	g.indent++
	g.writeln("return \"number\"")
	g.indent--
	g.writeln("case float64:")
	g.indent++
	g.writeln("return \"number\"")
	g.indent--
	g.writeln("case string:")
	g.indent++
	g.writeln("return \"string\"")
	g.indent--
	g.writeln("case bool:")
	g.indent++
	g.writeln("return \"boolean\"")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("return \"object\"")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// tostring helper
	g.writeln("func gts_tostring(v interface{}) string {")
	g.indent++
	g.writeln("return fmt.Sprintf(\"%%v\", v)")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// toint helper
	g.writeln("func gts_toint(v interface{}) int {")
	g.indent++
	g.writeln("switch x := v.(type) {")
	g.writeln("case int:")
	g.indent++
	g.writeln("return x")
	g.indent--
	g.writeln("case float64:")
	g.indent++
	g.writeln("return int(x)")
	g.indent--
	g.writeln("case string:")
	g.indent++
	g.writeln("var n int")
	g.writeln("fmt.Sscanf(x, \"%%d\", &n)")
	g.writeln("return n")
	g.indent--
	g.writeln("case bool:")
	g.indent++
	g.writeln("if x { return 1 }")
	g.writeln("return 0")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("return 0")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// tofloat helper
	g.writeln("func gts_tofloat(v interface{}) float64 {")
	g.indent++
	g.writeln("switch x := v.(type) {")
	g.writeln("case float64:")
	g.indent++
	g.writeln("return x")
	g.indent--
	g.writeln("case int:")
	g.indent++
	g.writeln("return float64(x)")
	g.indent--
	g.writeln("case string:")
	g.indent++
	g.writeln("var n float64")
	g.writeln("fmt.Sscanf(x, \"%%f\", &n)")
	g.writeln("return n")
	g.indent--
	g.writeln("case bool:")
	g.indent++
	g.writeln("if x { return 1 }")
	g.writeln("return 0")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("return 0")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_call - dynamic function caller using reflection
	g.writeln("func gts_call(fn interface{}, args ...interface{}) interface{} {")
	g.indent++
	g.writeln("v := reflect.ValueOf(fn)")
	g.writeln("in := make([]reflect.Value, len(args))")
	g.writeln("fnType := v.Type()")
	g.writeln("for i, arg := range args {")
	g.indent++
	g.writeln("if i < fnType.NumIn() {")
	g.indent++
	g.writeln("// Convert argument to expected type")
	g.writeln("expectedType := fnType.In(i)")
	g.writeln("argVal := reflect.ValueOf(arg)")
	g.writeln("if argVal.Type().ConvertibleTo(expectedType) {")
	g.indent++
	g.writeln("in[i] = argVal.Convert(expectedType)")
	g.indent--
	g.writeln("} else {")
	g.indent++
	g.writeln("in[i] = argVal")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("} else {")
	g.indent++
	g.writeln("in[i] = reflect.ValueOf(arg)")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("out := v.Call(in)")
	g.writeln("if len(out) > 0 {")
	g.indent++
	g.writeln("return out[0].Interface()")
	g.indent--
	g.writeln("}")
	g.writeln("return nil")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_tobool - convert any value to boolean for conditions
	g.writeln("func gts_tobool(v interface{}) bool {")
	g.indent++
	g.writeln("if v == nil {")
	g.indent++
	g.writeln("return false")
	g.indent--
	g.writeln("}")
	g.writeln("switch x := v.(type) {")
	g.writeln("case bool:")
	g.indent++
	g.writeln("return x")
	g.indent--
	g.writeln("case float64:")
	g.indent++
	g.writeln("return x != 0")
	g.indent--
	g.writeln("case string:")
	g.indent++
	g.writeln("return x != \"\"")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("return true")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_toarr_float - convert []interface{} to []float64
	g.writeln("func gts_toarr_float(v []interface{}) []float64 {")
	g.indent++
	g.writeln("result := make([]float64, len(v))")
	g.writeln("for i, x := range v {")
	g.indent++
	g.writeln("result[i] = gts_tofloat(x)")
	g.indent--
	g.writeln("}")
	g.writeln("return result")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_toarr_int - convert []interface{} to []int
	g.writeln("func gts_toarr_int(v []interface{}) []int {")
	g.indent++
	g.writeln("result := make([]int, len(v))")
	g.writeln("for i, x := range v {")
	g.indent++
	g.writeln("result[i] = gts_toint(x)")
	g.indent--
	g.writeln("}")
	g.writeln("return result")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Array method helpers
	// gts_map
	g.writeln("func gts_map(arr interface{}, fn interface{}) interface{} {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("if v.Kind() != reflect.Slice {")
	g.indent++
	g.writeln("return arr")
	g.indent--
	g.writeln("}")
	g.writeln("resultType := reflect.SliceOf(v.Type().Elem())")
	g.writeln("result := reflect.MakeSlice(resultType, 0, v.Len())")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("result = reflect.Append(result, out[0])")
	g.indent--
	g.writeln("}")
	g.writeln("return result.Interface()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_filter
	g.writeln("func gts_filter(arr interface{}, fn interface{}) interface{} {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("if v.Kind() != reflect.Slice {")
	g.indent++
	g.writeln("return arr")
	g.indent--
	g.writeln("}")
	g.writeln("resultType := reflect.SliceOf(v.Type().Elem())")
	g.writeln("result := reflect.MakeSlice(resultType, 0, v.Len())")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("if out[0].Bool() {")
	g.indent++
	g.writeln("result = reflect.Append(result, elem)")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("return result.Interface()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_reduce
	g.writeln("func gts_reduce(arr interface{}, initial interface{}, fn interface{}) interface{} {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("acc := reflect.ValueOf(initial)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{acc, elem})")
	g.writeln("acc = out[0]")
	g.indent--
	g.writeln("}")
	g.writeln("return acc.Interface()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_find
	g.writeln("func gts_find(arr interface{}, fn interface{}) interface{} {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("if out[0].Bool() {")
	g.indent++
	g.writeln("return elem.Interface()")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("var zero interface{}")
	g.writeln("return zero")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_findIndex
	g.writeln("func gts_findIndex(arr interface{}, fn interface{}) int {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("if out[0].Bool() {")
	g.indent++
	g.writeln("return i")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("return -1")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_some
	g.writeln("func gts_some(arr interface{}, fn interface{}) bool {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("if out[0].Bool() {")
	g.indent++
	g.writeln("return true")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("return false")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_every
	g.writeln("func gts_every(arr interface{}, fn interface{}) bool {")
	g.indent++
	g.writeln("v := reflect.ValueOf(arr)")
	g.writeln("f := reflect.ValueOf(fn)")
	g.writeln("for i := 0; i < v.Len(); i++ {")
	g.indent++
	g.writeln("elem := v.Index(i)")
	g.writeln("out := f.Call([]reflect.Value{elem})")
	g.writeln("if !out[0].Bool() {")
	g.indent++
	g.writeln("return false")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("return true")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Event loop runtime (must be before Promise runtime)
	g.genEventLoopRuntime()
	g.writeln("")

	// Promise runtime type and helpers
	g.genPromiseRuntime()
}

// genEventLoopRuntime generates the JavaScript-compatible event loop.
func (g *Generator) genEventLoopRuntime() {
	g.writeln("// GTS_EventLoop implements JavaScript-compatible event loop")
	g.writeln("type GTS_EventLoop struct {")
	g.indent++
	g.writeln("microtasks    []func()")
	g.writeln("macrotasks    []func()")
	g.writeln("timers        map[int64]*time.Timer")
	g.writeln("timerIDGen    int64")
	g.writeln("pendingTimers chan func()")
	g.writeln("mu            sync.Mutex")
	g.indent--
	g.writeln("}")
	g.writeln("")

	g.writeln("var gts_eventLoop *GTS_EventLoop")
	g.writeln("")

	// Queue microtask
	g.writeln("func gts_queueMicrotask(fn func()) {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("gts_eventLoop.microtasks = append(gts_eventLoop.microtasks, fn)")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Queue macrotask
	g.writeln("func gts_queueMacrotask(fn func()) {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("gts_eventLoop.macrotasks = append(gts_eventLoop.macrotasks, fn)")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// setTimeout
	g.writeln("func gts_setTimeout(callback func(), delay int) int64 {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("gts_eventLoop.timerIDGen++")
	g.writeln("id := gts_eventLoop.timerIDGen")
	g.writeln("timerID := id")
	g.writeln("timer := time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("delete(gts_eventLoop.timers, timerID)")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("gts_eventLoop.pendingTimers <- callback")
	g.indent--
	g.writeln("})")
	g.writeln("gts_eventLoop.timers[id] = timer")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("return id")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// setInterval
	g.writeln("func gts_setInterval(callback func(), delay int) int64 {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("gts_eventLoop.timerIDGen++")
	g.writeln("id := gts_eventLoop.timerIDGen")
	g.writeln("var timer *time.Timer")
	g.writeln("var tick func()")
	g.writeln("tick = func() {")
	g.indent++
	g.writeln("gts_eventLoop.pendingTimers <- callback")
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("if _, exists := gts_eventLoop.timers[id]; exists {")
	g.indent++
	g.writeln("timer = time.AfterFunc(time.Duration(delay)*time.Millisecond, tick)")
	g.writeln("gts_eventLoop.timers[id] = timer")
	g.indent--
	g.writeln("}")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.writeln("timer = time.AfterFunc(time.Duration(delay)*time.Millisecond, tick)")
	g.writeln("gts_eventLoop.timers[id] = timer")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("return id")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// clearTimeout / clearInterval
	g.writeln("func gts_clearTimeout(id int64) {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("if timer, exists := gts_eventLoop.timers[id]; exists {")
	g.indent++
	g.writeln("timer.Stop()")
	g.writeln("delete(gts_eventLoop.timers, id)")
	g.indent--
	g.writeln("}")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Run event loop
	g.writeln("func gts_runEventLoop() {")
	g.indent++
	g.writeln("for {")
	g.indent++
	// Drain all microtasks
	g.writeln("for {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("if len(gts_eventLoop.microtasks) == 0 {")
	g.indent++
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("break")
	g.indent--
	g.writeln("}")
	g.writeln("task := gts_eventLoop.microtasks[0]")
	g.writeln("gts_eventLoop.microtasks = gts_eventLoop.microtasks[1:]")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("task()")
	g.indent--
	g.writeln("}")
	// Execute one macrotask
	g.writeln("gts_eventLoop.mu.Lock()")
	g.writeln("if len(gts_eventLoop.macrotasks) > 0 {")
	g.indent++
	g.writeln("task := gts_eventLoop.macrotasks[0]")
	g.writeln("gts_eventLoop.macrotasks = gts_eventLoop.macrotasks[1:]")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("task()")
	g.writeln("continue")
	g.indent--
	g.writeln("}")
	// Wait for timer or exit
	g.writeln("hasTimers := len(gts_eventLoop.timers) > 0")
	g.writeln("gts_eventLoop.mu.Unlock()")
	g.writeln("if !hasTimers {")
	g.indent++
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("select {")
	g.writeln("case task := <-gts_eventLoop.pendingTimers:")
	g.indent++
	g.writeln("gts_queueMacrotask(task)")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("")
}

// genPromiseRuntime generates the Promise type and related helpers.
func (g *Generator) genPromiseRuntime() {
	// Promise struct with event loop support
	g.writeln("// GTS_Promise represents a JavaScript-style Promise with event loop support")
	g.writeln("type GTS_Promise[T any] struct {")
	g.indent++
	g.writeln("state     int // 0=pending, 1=fulfilled, 2=rejected")
	g.writeln("value     T")
	g.writeln("err       error")
	g.writeln("onFulfill []func(T)")
	g.writeln("onReject  []func(error)")
	g.writeln("onFinally []func()")
	g.writeln("mu        sync.Mutex")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise constructor - runs executor synchronously
	g.writeln("func GTS_NewPromise[T any](executor func(resolve func(T), reject func(error))) *GTS_Promise[T] {")
	g.indent++
	g.writeln("p := &GTS_Promise[T]{state: 0}")
	g.writeln("resolve := func(v T) {")
	g.indent++
	g.writeln("p.mu.Lock()")
	g.writeln("if p.state != 0 {")
	g.indent++
	g.writeln("p.mu.Unlock()")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("p.state = 1")
	g.writeln("p.value = v")
	g.writeln("handlers := p.onFulfill")
	g.writeln("finalizers := p.onFinally")
	g.writeln("p.mu.Unlock()")
	g.writeln("for _, h := range handlers {")
	g.indent++
	g.writeln("handler := h")
	g.writeln("val := v")
	g.writeln("gts_queueMicrotask(func() { handler(val) })")
	g.indent--
	g.writeln("}")
	g.writeln("for _, f := range finalizers {")
	g.indent++
	g.writeln("finalizer := f")
	g.writeln("gts_queueMicrotask(func() { finalizer() })")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.writeln("reject := func(e error) {")
	g.indent++
	g.writeln("p.mu.Lock()")
	g.writeln("if p.state != 0 {")
	g.indent++
	g.writeln("p.mu.Unlock()")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("p.state = 2")
	g.writeln("p.err = e")
	g.writeln("handlers := p.onReject")
	g.writeln("finalizers := p.onFinally")
	g.writeln("p.mu.Unlock()")
	g.writeln("for _, h := range handlers {")
	g.indent++
	g.writeln("handler := h")
	g.writeln("err := e")
	g.writeln("gts_queueMicrotask(func() { handler(err) })")
	g.indent--
	g.writeln("}")
	g.writeln("for _, f := range finalizers {")
	g.indent++
	g.writeln("finalizer := f")
	g.writeln("gts_queueMicrotask(func() { finalizer() })")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	// Execute synchronously (not in goroutine)
	g.writeln("defer func() {")
	g.indent++
	g.writeln("if r := recover(); r != nil {")
	g.indent++
	g.writeln("if err, ok := r.(error); ok {")
	g.indent++
	g.writeln("reject(err)")
	g.indent--
	g.writeln("} else {")
	g.indent++
	g.writeln("reject(fmt.Errorf(\"%%v\", r))")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}()")
	g.writeln("executor(resolve, reject)")
	g.writeln("return p")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.Then method
	g.writeln("func (p *GTS_Promise[T]) Then(onFulfill func(T) interface{}) *GTS_Promise[interface{}] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(interface{}), reject func(error)) {")
	g.indent++
	g.writeln("handler := func(v T) {")
	g.indent++
	g.writeln("defer func() {")
	g.indent++
	g.writeln("if r := recover(); r != nil {")
	g.indent++
	g.writeln("if err, ok := r.(error); ok {")
	g.indent++
	g.writeln("reject(err)")
	g.indent--
	g.writeln("} else {")
	g.indent++
	g.writeln("reject(fmt.Errorf(\"%%v\", r))")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}()")
	g.writeln("result := onFulfill(v)")
	g.writeln("resolve(result)")
	g.indent--
	g.writeln("}")
	g.writeln("errHandler := func(e error) { reject(e) }")
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 0:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, handler)")
	g.writeln("p.onReject = append(p.onReject, errHandler)")
	g.indent--
	g.writeln("case 1:")
	g.indent++
	g.writeln("val := p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { handler(val) })")
	g.writeln("return")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err := p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { errHandler(err) })")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.Catch method
	g.writeln("func (p *GTS_Promise[T]) Catch(onReject func(error) interface{}) *GTS_Promise[interface{}] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(interface{}), reject func(error)) {")
	g.indent++
	g.writeln("fulfillHandler := func(v T) { resolve(v) }")
	g.writeln("errHandler := func(e error) {")
	g.indent++
	g.writeln("defer func() {")
	g.indent++
	g.writeln("if r := recover(); r != nil {")
	g.indent++
	g.writeln("if err, ok := r.(error); ok {")
	g.indent++
	g.writeln("reject(err)")
	g.indent--
	g.writeln("} else {")
	g.indent++
	g.writeln("reject(fmt.Errorf(\"%%v\", r))")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}()")
	g.writeln("result := onReject(e)")
	g.writeln("resolve(result)")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 0:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, fulfillHandler)")
	g.writeln("p.onReject = append(p.onReject, errHandler)")
	g.indent--
	g.writeln("case 1:")
	g.indent++
	g.writeln("val := p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { fulfillHandler(val) })")
	g.writeln("return")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err := p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { errHandler(err) })")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.Finally method
	g.writeln("func (p *GTS_Promise[T]) Finally(onFinally func()) *GTS_Promise[T] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(T), reject func(error)) {")
	g.indent++
	g.writeln("fulfillHandler := func(v T) {")
	g.indent++
	g.writeln("onFinally()")
	g.writeln("resolve(v)")
	g.indent--
	g.writeln("}")
	g.writeln("errHandler := func(e error) {")
	g.indent++
	g.writeln("onFinally()")
	g.writeln("reject(e)")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 0:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, fulfillHandler)")
	g.writeln("p.onReject = append(p.onReject, errHandler)")
	g.indent--
	g.writeln("case 1:")
	g.indent++
	g.writeln("val := p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { fulfillHandler(val) })")
	g.writeln("return")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err := p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("gts_queueMicrotask(func() { errHandler(err) })")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// gts_await helper - uses event loop for waiting
	g.writeln("func gts_await[T any](p *GTS_Promise[T]) T {")
	g.indent++
	g.writeln("done := make(chan struct{})")
	g.writeln("var result T")
	g.writeln("var err error")
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 1:")
	g.indent++
	g.writeln("result = p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("return result")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err = p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("panic(err)")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, func(v T) {")
	g.indent++
	g.writeln("result = v")
	g.writeln("close(done)")
	g.indent--
	g.writeln("})")
	g.writeln("p.onReject = append(p.onReject, func(e error) {")
	g.indent++
	g.writeln("err = e")
	g.writeln("close(done)")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("p.mu.Unlock()")
	g.writeln("<-done")
	g.writeln("if err != nil {")
	g.indent++
	g.writeln("panic(err)")
	g.indent--
	g.writeln("}")
	g.writeln("return result")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.resolve static method
	g.writeln("func GTS_Promise_Resolve[T any](value T) *GTS_Promise[T] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(T), reject func(error)) {")
	g.indent++
	g.writeln("resolve(value)")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.reject static method
	g.writeln("func GTS_Promise_Reject[T any](err error) *GTS_Promise[T] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(T), reject func(error)) {")
	g.indent++
	g.writeln("reject(err)")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.all
	g.writeln("func GTS_Promise_All[T any](promises []*GTS_Promise[T]) *GTS_Promise[[]T] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func([]T), reject func(error)) {")
	g.indent++
	g.writeln("if len(promises) == 0 {")
	g.indent++
	g.writeln("resolve([]T{})")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("results := make([]T, len(promises))")
	g.writeln("remaining := len(promises)")
	g.writeln("rejected := false")
	g.writeln("var mu sync.Mutex")
	g.writeln("for i, p := range promises {")
	g.indent++
	g.writeln("idx := i")
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 1:")
	g.indent++
	g.writeln("results[idx] = p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("mu.Lock()")
	g.writeln("remaining--")
	g.writeln("if remaining == 0 && !rejected {")
	g.indent++
	g.writeln("mu.Unlock()")
	g.writeln("resolve(results)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err := p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("mu.Lock()")
	g.writeln("if !rejected {")
	g.indent++
	g.writeln("rejected = true")
	g.writeln("mu.Unlock()")
	g.writeln("reject(err)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, func(v T) {")
	g.indent++
	g.writeln("mu.Lock()")
	g.writeln("results[idx] = v")
	g.writeln("remaining--")
	g.writeln("if remaining == 0 && !rejected {")
	g.indent++
	g.writeln("mu.Unlock()")
	g.writeln("resolve(results)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.writeln("p.onReject = append(p.onReject, func(e error) {")
	g.indent++
	g.writeln("mu.Lock()")
	g.writeln("if !rejected {")
	g.indent++
	g.writeln("rejected = true")
	g.writeln("mu.Unlock()")
	g.writeln("reject(e)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.writeln("p.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
	g.writeln("")

	// Promise.race
	g.writeln("func GTS_Promise_Race[T any](promises []*GTS_Promise[T]) *GTS_Promise[T] {")
	g.indent++
	g.writeln("return GTS_NewPromise(func(resolve func(T), reject func(error)) {")
	g.indent++
	g.writeln("settled := false")
	g.writeln("var mu sync.Mutex")
	g.writeln("for _, p := range promises {")
	g.indent++
	g.writeln("p.mu.Lock()")
	g.writeln("switch p.state {")
	g.writeln("case 1:")
	g.indent++
	g.writeln("val := p.value")
	g.writeln("p.mu.Unlock()")
	g.writeln("mu.Lock()")
	g.writeln("if !settled {")
	g.indent++
	g.writeln("settled = true")
	g.writeln("mu.Unlock()")
	g.writeln("resolve(val)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.writeln("return")
	g.indent--
	g.writeln("case 2:")
	g.indent++
	g.writeln("err := p.err")
	g.writeln("p.mu.Unlock()")
	g.writeln("mu.Lock()")
	g.writeln("if !settled {")
	g.indent++
	g.writeln("settled = true")
	g.writeln("mu.Unlock()")
	g.writeln("reject(err)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.writeln("return")
	g.indent--
	g.writeln("default:")
	g.indent++
	g.writeln("p.onFulfill = append(p.onFulfill, func(v T) {")
	g.indent++
	g.writeln("mu.Lock()")
	g.writeln("if !settled {")
	g.indent++
	g.writeln("settled = true")
	g.writeln("mu.Unlock()")
	g.writeln("resolve(v)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.writeln("p.onReject = append(p.onReject, func(e error) {")
	g.indent++
	g.writeln("mu.Lock()")
	g.writeln("if !settled {")
	g.indent++
	g.writeln("settled = true")
	g.writeln("mu.Unlock()")
	g.writeln("reject(e)")
	g.writeln("return")
	g.indent--
	g.writeln("}")
	g.writeln("mu.Unlock()")
	g.indent--
	g.writeln("})")
	g.writeln("p.mu.Unlock()")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("})")
	g.indent--
	g.writeln("}")
}

// genEnum generates a Go type and const block for an enum.
func (g *Generator) genEnum(enum *typed.EnumDecl) {
	name := exportName(enum.Name)

	// Generate the type definition
	g.writeln("type %s int", name)
	g.writeln("")

	// Generate const block for enum members
	if len(enum.Members) > 0 {
		g.writeln("const (")
		g.indent++
		for _, member := range enum.Members {
			memberName := name + exportName(member.Name)
			g.writeln("%s %s = %d", memberName, name, member.Value)
		}
		g.indent--
		g.writeln(")")
	}
}

// genTypeAlias generates a Go type alias from a TypeAlias.
func (g *Generator) genTypeAlias(alias *typed.TypeAlias) {
	name := exportName(alias.Name)
	goTypeName := g.goType(alias.Resolved)
	g.writeln("type %s %s", name, goTypeName)
}

// genInterface generates a Go interface type from an InterfaceDecl.
func (g *Generator) genInterface(iface *typed.InterfaceDecl) {
	name := exportName(iface.Name)

	g.writeln("type %s interface {", name)
	g.indent++

	for _, method := range iface.Methods {
		params := make([]string, len(method.Params))
		for i, p := range method.Params {
			params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
		}

		if method.ReturnType == nil || types.VoidType.Equals(method.ReturnType) {
			g.writeln("%s(%s)", exportName(method.Name), strings.Join(params, ", "))
		} else {
			g.writeln("%s(%s) %s", exportName(method.Name), strings.Join(params, ", "), g.goType(method.ReturnType))
		}
	}

	g.indent--
	g.writeln("}")
}

func (g *Generator) genClass(class *typed.ClassDecl) {
	name := exportName(class.Name)

	// Generate type parameters for generic classes
	typeParamStr := ""
	typeParamNames := ""
	if len(class.TypeParams) > 0 {
		typeParams := make([]string, len(class.TypeParams))
		typeNames := make([]string, len(class.TypeParams))
		for i, tp := range class.TypeParams {
			typeNames[i] = tp.Name
			if tp.Constraint != nil {
				typeParams[i] = fmt.Sprintf("%s %s", tp.Name, g.goType(tp.Constraint))
			} else {
				typeParams[i] = fmt.Sprintf("%s any", tp.Name)
			}
		}
		typeParamStr = fmt.Sprintf("[%s]", strings.Join(typeParams, ", "))
		typeParamNames = fmt.Sprintf("[%s]", strings.Join(typeNames, ", "))
	}

	// Generate struct
	g.writeln("type %s%s struct {", name, typeParamStr)
	g.indent++

	// Embed super class if present
	if class.Super != "" {
		g.writeln("%s", exportName(class.Super))
	}

	// Generate fields
	for _, field := range class.Fields {
		g.writeln("%s %s", exportName(field.Name), g.goType(field.Type))
	}

	g.indent--
	g.writeln("}")
	g.writeln("")

	// Generate constructor function
	if class.Constructor != nil {
		params := make([]string, len(class.Constructor.Params))
		for i, p := range class.Constructor.Params {
			params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
		}

		g.writeln("func New%s%s(%s) *%s%s {", name, typeParamStr, strings.Join(params, ", "), name, typeParamNames)
		g.indent++
		g.writeln("this := &%s%s{}", name, typeParamNames)

		// Set current class for super() handling
		prevClass := g.currentClass
		g.currentClass = class

		// Generate constructor body
		for _, stmt := range class.Constructor.Body.Stmts {
			g.genStmt(stmt)
		}

		g.currentClass = prevClass

		g.writeln("return this")
		g.indent--
		g.writeln("}")
		g.writeln("")
	} else {
		// Default constructor
		g.writeln("func New%s%s() *%s%s {", name, typeParamStr, name, typeParamNames)
		g.indent++
		g.writeln("return &%s%s{}", name, typeParamNames)
		g.indent--
		g.writeln("}")
		g.writeln("")
	}

	// Generate methods
	for _, method := range class.Methods {
		params := make([]string, len(method.Params))
		for i, p := range method.Params {
			params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
		}

		returnType := g.goType(method.ReturnType)
		if returnType == "" {
			g.writeln("func (this *%s%s) %s(%s) {", name, typeParamNames, exportName(method.Name), strings.Join(params, ", "))
		} else {
			g.writeln("func (this *%s%s) %s(%s) %s {", name, typeParamNames, exportName(method.Name), strings.Join(params, ", "), returnType)
		}
		g.indent++

		for _, stmt := range method.Body.Stmts {
			g.genStmt(stmt)
		}

		g.indent--
		g.writeln("}")
		g.writeln("")
	}
}

func (g *Generator) genFuncDecl(fn *typed.FuncDecl) {
	// Handle decorated functions specially
	if len(fn.Decorators) > 0 {
		g.genDecoratedFuncDecl(fn)
		return
	}

	params := make([]string, len(fn.Params))
	for i, p := range fn.Params {
		params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
	}

	returnType := g.goType(fn.ReturnType)

	// Track current return type for type assertions
	savedRetType := g.currentRetType
	g.currentRetType = fn.ReturnType

	// Generate type parameters for generic functions
	typeParamStr := ""
	if len(fn.TypeParams) > 0 {
		typeParams := make([]string, len(fn.TypeParams))
		for i, tp := range fn.TypeParams {
			if tp.Constraint != nil {
				typeParams[i] = fmt.Sprintf("%s %s", tp.Name, g.goType(tp.Constraint))
			} else {
				typeParams[i] = fmt.Sprintf("%s any", tp.Name)
			}
		}
		typeParamStr = fmt.Sprintf("[%s]", strings.Join(typeParams, ", "))
	}

	if fn.IsAsync {
		// Async function: returns Promise
		// Extract the inner type from Promise<T>
		promiseType, ok := fn.ReturnType.(*types.Promise)
		if !ok {
			// Fallback if type is not Promise
			promiseType = &types.Promise{Value: types.VoidType}
		}
		innerType := g.goType(promiseType.Value)
		if innerType == "" {
			innerType = "interface{}"
		}

		g.writeln("func %s%s(%s) *GTS_Promise[%s] {", goName(fn.Name), typeParamStr, strings.Join(params, ", "), innerType)
		g.indent++
		g.writeln("return GTS_NewPromise(func(__resolve func(%s), __reject func(error)) {", innerType)
		g.indent++

		// Generate body with transformed returns
		g.genAsyncBody(fn.Body, innerType)

		g.indent--
		g.writeln("})")
		g.indent--
		g.writeln("}")
	} else {
		// Regular function
		if returnType == "" {
			g.writeln("func %s%s(%s) {", goName(fn.Name), typeParamStr, strings.Join(params, ", "))
		} else {
			g.writeln("func %s%s(%s) %s {", goName(fn.Name), typeParamStr, strings.Join(params, ", "), returnType)
		}
		g.indent++

		for _, stmt := range fn.Body.Stmts {
			g.genStmt(stmt)
		}

		g.indent--
		g.writeln("}")
	}

	g.currentRetType = savedRetType
}

// genDecoratedFuncDecl generates a decorated function.
// @decorator
// function a(): int { return 1 }
// becomes:
// func _a() int { return 1 }
// var A = decorator(_a)
func (g *Generator) genDecoratedFuncDecl(fn *typed.FuncDecl) {
	params := make([]string, len(fn.Params))
	for i, p := range fn.Params {
		params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
	}

	returnType := g.goType(fn.ReturnType)

	// Track current return type for type assertions
	savedRetType := g.currentRetType
	g.currentRetType = fn.ReturnType

	// Internal function name (with underscore prefix)
	internalName := "_" + fn.Name

	// Generate the original function with internal name
	if returnType == "" {
		g.writeln("func %s(%s) {", internalName, strings.Join(params, ", "))
	} else {
		g.writeln("func %s(%s) %s {", internalName, strings.Join(params, ", "), returnType)
	}
	g.indent++

	for _, stmt := range fn.Body.Stmts {
		g.genStmt(stmt)
	}

	g.indent--
	g.writeln("}")
	g.writeln("")

	// Generate the decorated variable
	// Chain decorators: @d1 @d2 function a() {} => d1(d2(_a))
	// Decorators are applied from bottom to top (closest to function first)
	decoratedExpr := internalName
	for i := len(fn.Decorators) - 1; i >= 0; i-- {
		d := fn.Decorators[i]
		if d.Object != "" {
			// Member-access decorator with args: @obj.method(args)
			// becomes: gts_call(obj.Method(args...), _fn)
			args := make([]string, len(d.Arguments))
			for j, arg := range d.Arguments {
				args[j] = g.genExpr(arg)
			}
			decoratedExpr = fmt.Sprintf("gts_call(%s.%s(%s), %s)", goName(d.Object), exportName(d.Property), strings.Join(args, ", "), decoratedExpr)
		} else if len(d.Arguments) > 0 {
			// Simple parameterized decorator: @name(args)
			// becomes: gts_call(Name(args...), _fn)
			args := make([]string, len(d.Arguments))
			for j, arg := range d.Arguments {
				args[j] = g.genExpr(arg)
			}
			decoratedExpr = fmt.Sprintf("gts_call(%s(%s), %s)", goName(d.Name), strings.Join(args, ", "), decoratedExpr)
		} else {
			// Simple decorator: @name => Name(_fn)
			decoratedExpr = fmt.Sprintf("%s(%s)", goName(d.Name), decoratedExpr)
		}
	}

	g.writeln("var %s = %s", goName(fn.Name), decoratedExpr)
	g.writeln("")

	g.currentRetType = savedRetType
}

// genAsyncBody generates the body of an async function, transforming return statements.
func (g *Generator) genAsyncBody(body *typed.BlockStmt, innerType string) {
	for _, stmt := range body.Stmts {
		if ret, ok := stmt.(*typed.ReturnStmt); ok {
			// Transform return x into __resolve(x); return
			if ret.Value != nil {
				g.writeln("__resolve(%s)", g.genExpr(ret.Value))
			} else {
				// For void returns in Promise<void>
				if innerType == "interface{}" {
					g.writeln("__resolve(nil)")
				}
			}
			g.writeln("return")
		} else {
			g.genStmt(stmt)
		}
	}
}

// ----------------------------------------------------------------------------
// Statement Generation
// ----------------------------------------------------------------------------

func (g *Generator) genStmt(stmt typed.Stmt) {
	switch s := stmt.(type) {
	case *typed.ExprStmt:
		exprCode := g.genExpr(s.Expr)
		// Check if the expression has a non-void return type that needs to be discarded
		exprType := types.Unwrap(s.Expr.Type())
		if prim, ok := exprType.(*types.Primitive); ok && prim.Kind == types.KindVoid {
			g.writeln("%s", exprCode)
		} else if exprType.Equals(types.AnyType) {
			// For any type, check if it's a timer function that returns a value
			if builtin, ok := s.Expr.(*typed.BuiltinCall); ok {
				if builtin.Name == "setTimeout" || builtin.Name == "setInterval" {
					g.writeln("_ = %s", exprCode)
				} else {
					g.writeln("%s", exprCode)
				}
			} else {
				g.writeln("%s", exprCode)
			}
		} else if _, ok := exprType.(*types.Promise); ok {
			// Promises don't need to be discarded (they're used for side effects)
			g.writeln("%s", exprCode)
		} else if prim, ok := exprType.(*types.Primitive); ok && prim.Kind != types.KindVoid {
			// Non-void primitives (like int from setTimeout) need to be discarded
			if builtin, ok := s.Expr.(*typed.BuiltinCall); ok {
				if builtin.Name == "setTimeout" || builtin.Name == "setInterval" {
					g.writeln("_ = %s", exprCode)
				} else {
					g.writeln("%s", exprCode)
				}
			} else {
				g.writeln("%s", exprCode)
			}
		} else {
			g.writeln("%s", exprCode)
		}

	case *typed.VarDecl:
		g.genVarDecl(s)

	case *typed.BlockStmt:
		g.writeln("{")
		g.indent++
		for _, inner := range s.Stmts {
			g.genStmt(inner)
		}
		g.indent--
		g.writeln("}")

	case *typed.IfStmt:
		g.genIfStmt(s)

	case *typed.WhileStmt:
		// Check if condition is any type (interface{}) and needs gts_tobool
		condType := types.Unwrap(s.Condition.Type())
		condExpr := g.genExpr(s.Condition)
		if prim, ok := condType.(*types.Primitive); ok && prim.Kind == types.KindAny {
			condExpr = fmt.Sprintf("gts_tobool(%s)", condExpr)
		}
		g.writeln("for %s {", condExpr)
		g.indent++
		for _, inner := range s.Body.Stmts {
			g.genStmt(inner)
		}
		g.indent--
		g.writeln("}")

	case *typed.ForStmt:
		g.genForStmt(s)

	case *typed.ForOfStmt:
		g.genForOfStmt(s)

	case *typed.SwitchStmt:
		g.genSwitchStmt(s)

	case *typed.ReturnStmt:
		if s.Value != nil {
			retExpr := g.genExpr(s.Value)
			// Check if we need type assertion/conversion
			if g.currentRetType != nil {
				valType := types.Unwrap(s.Value.Type())
				expectedType := types.Unwrap(g.currentRetType)

				// If value is any but expected is concrete
				if prim, ok := valType.(*types.Primitive); ok && prim.Kind == types.KindAny {
					expectedGoType := g.goType(expectedType)
					if expectedGoType != "interface{}" {
						// Check if this is a binary arithmetic expression - if so, we've already
						// handled type assertions for the operands, so the result is already typed
						needsAssertion := true
						if binExpr, ok := s.Value.(*typed.BinaryExpr); ok {
							switch binExpr.Op {
							case "+", "-", "*", "/", "%", "<", ">", "<=", ">=":
								// Skip - we've already asserted operands in genBinaryExpr
								needsAssertion = false
							}
						}
						if needsAssertion {
							retExpr = fmt.Sprintf("%s.(%s)", retExpr, expectedGoType)
						}
					}
				}
				// If value is []any but expected is []concrete
				if valArr, ok := valType.(*types.Array); ok {
					if valElem, ok := valArr.Element.(*types.Primitive); ok && valElem.Kind == types.KindAny {
						if expArr, ok := expectedType.(*types.Array); ok {
							if expElem, ok := expArr.Element.(*types.Primitive); ok {
								if expElem.Kind == types.KindInt {
									retExpr = fmt.Sprintf("gts_toarr_int(%s)", retExpr)
								} else if expElem.Kind == types.KindFloat {
									retExpr = fmt.Sprintf("gts_toarr_float(%s)", retExpr)
								}
							}
						}
					}
				}
			}
			g.writeln("return %s", retExpr)
		} else {
			g.writeln("return")
		}

	case *typed.BreakStmt:
		g.writeln("break")

	case *typed.ContinueStmt:
		g.writeln("continue")

	case *typed.TryStmt:
		g.genTryStmt(s)

	case *typed.ThrowStmt:
		g.genThrowStmt(s)

	case *typed.FuncDecl:
		g.genFuncDecl(s)

	case *typed.ClassDecl:
		g.genClass(s)
	}
}

func (g *Generator) genVarDecl(decl *typed.VarDecl) {
	// Handle destructuring patterns
	if decl.Pattern != nil {
		g.genDestructuringDecl(decl)
		return
	}

	if decl.Init != nil {
		// Use explicit type declaration for clarity
		goType := g.goType(decl.VarType)

		// Handle empty array literals - use the declared type's element type
		initExpr := g.genExprWithContext(decl.Init, decl.VarType)

		if goType != "" {
			// Always use explicit type to ensure correct typing
			// This is especially important for interface{} types
			g.writeln("var %s %s = %s", goName(decl.Name), goType, initExpr)
		} else {
			g.writeln("%s := %s", goName(decl.Name), initExpr)
		}
	} else {
		g.writeln("var %s %s", goName(decl.Name), g.goType(decl.VarType))
	}
}

func (g *Generator) genDestructuringDecl(decl *typed.VarDecl) {
	// Generate a temporary variable to hold the source value
	tempVar := "_destructure_temp"
	initExpr := g.genExprWithContext(decl.Init, decl.VarType)
	g.writeln("%s := %s", tempVar, initExpr)

	// Generate assignments for each pattern element
	g.genPatternAssignments(decl.Pattern, tempVar)
}

func (g *Generator) genPatternAssignments(pattern typed.Pattern, source string) {
	switch p := pattern.(type) {
	case *typed.ArrayPattern:
		for i, elem := range p.Elements {
			if elem != nil {
				elemSource := fmt.Sprintf("%s[%d]", source, i)
				g.genPatternAssignments(elem, elemSource)
			}
		}
	case *typed.ObjectPattern:
		for _, prop := range p.Properties {
			propSource := fmt.Sprintf("%s.%s", source, exportName(prop.Key))
			g.genPatternAssignments(prop.Value, propSource)
		}
	case *typed.IdentPattern:
		goType := g.goType(p.PatternType)
		if goType != "" {
			g.writeln("var %s %s = %s", goName(p.Name), goType, source)
		} else {
			g.writeln("%s := %s", goName(p.Name), source)
		}
	}
}

// genExprWithContext generates an expression with knowledge of the target type context
func (g *Generator) genExprWithContext(expr typed.Expr, targetType types.Type) string {
	// Don't wrap null literals - they should just be nil
	if _, isNull := expr.(*typed.NullLit); isNull {
		return "nil"
	}

	// Handle assigning any type to a primitive type (e.g., from gts_call result)
	exprType := types.Unwrap(expr.Type())
	targetTypeUnwrapped := types.Unwrap(targetType)
	if exprPrim, ok := exprType.(*types.Primitive); ok && exprPrim.Kind == types.KindAny {
		if targetPrim, ok := targetTypeUnwrapped.(*types.Primitive); ok {
			value := g.genExpr(expr)
			switch targetPrim.Kind {
			case types.KindInt:
				return fmt.Sprintf("gts_toint(%s)", value)
			case types.KindFloat, types.KindNumber:
				return fmt.Sprintf("gts_tofloat(%s)", value)
			case types.KindString:
				return fmt.Sprintf("gts_tostring(%s)", value)
			case types.KindBoolean:
				return fmt.Sprintf("gts_tobool(%s)", value)
			}
		}
	}

	// Handle numeric type conversions: int <-> float64 (number)
	if exprPrim, ok := exprType.(*types.Primitive); ok {
		if targetPrim, ok := targetTypeUnwrapped.(*types.Primitive); ok {
			// int to number/float: needs float64() conversion
			if exprPrim.Kind == types.KindInt && (targetPrim.Kind == types.KindNumber || targetPrim.Kind == types.KindFloat) {
				value := g.genExpr(expr)
				return fmt.Sprintf("float64(%s)", value)
			}
			// number/float to int: needs int() conversion (should have been caught by type checker, but handle gracefully)
			if (exprPrim.Kind == types.KindNumber || exprPrim.Kind == types.KindFloat) && targetPrim.Kind == types.KindInt {
				value := g.genExpr(expr)
				return fmt.Sprintf("int(%s)", value)
			}
		}
	}

	// Handle assigning a non-pointer value to a nullable (pointer) type
	// e.g., `var name: string | null = "Alice"` needs to become `func() *string { v := "Alice"; return &v }()`
	if nullable, ok := targetTypeUnwrapped.(*types.Nullable); ok {
		// Check if expr type is not already nullable (i.e., not a pointer)
		if _, exprIsNullable := exprType.(*types.Nullable); !exprIsNullable {
			// For class types, we don't need pointer wrapping since classes are already pointers
			if _, isClass := nullable.Inner.(*types.Class); !isClass {
				ptrGoType := g.goType(targetType)
				value := g.genExpr(expr)
				return fmt.Sprintf("func() %s { v := %s; return &v }()", ptrGoType, value)
			}
		}
	}

	// Handle array literals where the expression type is any[] but target is concrete[]
	if arrLit, ok := expr.(*typed.ArrayLit); ok {
		if targetArr, ok := targetType.(*types.Array); ok {
			exprArr := arrLit.ExprType.(*types.Array)
			if exprElem, ok := exprArr.Element.(*types.Primitive); ok && exprElem.Kind == types.KindAny {
				// Check for spread expressions - if present, delegate to genArrayLit
				hasSpread := false
				for _, elem := range arrLit.Elements {
					if _, ok := elem.(*typed.SpreadExpr); ok {
						hasSpread = true
						break
					}
				}
				if hasSpread {
					// Create a temporary ArrayLit with the target type for proper spread handling
					tempArr := &typed.ArrayLit{
						Elements: arrLit.Elements,
						ExprType: targetArr,
					}
					return g.genArrayLit(tempArr)
				}

				// Use target element type instead (no spreads)
				elements := make([]string, len(arrLit.Elements))
				for i, elem := range arrLit.Elements {
					elements[i] = g.genExpr(elem)
				}
				return fmt.Sprintf("[]%s{%s}", g.goType(targetArr.Element), strings.Join(elements, ", "))
			}
		}
	}
	return g.genExpr(expr)
}

func (g *Generator) genIfStmt(stmt *typed.IfStmt) {
	// Check if condition is any type (interface{}) and needs gts_tobool
	condType := types.Unwrap(stmt.Condition.Type())
	condExpr := g.genExpr(stmt.Condition)
	if prim, ok := condType.(*types.Primitive); ok && prim.Kind == types.KindAny {
		condExpr = fmt.Sprintf("gts_tobool(%s)", condExpr)
	}
	g.writeln("if %s {", condExpr)
	g.indent++
	for _, inner := range stmt.Then.Stmts {
		g.genStmt(inner)
	}
	g.indent--

	if stmt.Else != nil {
		switch elseStmt := stmt.Else.(type) {
		case *typed.IfStmt:
			g.write("} else ")
			g.genIfStmt(elseStmt)
			return
		case *typed.BlockStmt:
			g.writeln("} else {")
			g.indent++
			for _, inner := range elseStmt.Stmts {
				g.genStmt(inner)
			}
			g.indent--
		}
	}
	g.writeln("}")
}

func (g *Generator) genForStmt(stmt *typed.ForStmt) {
	// Go-style for loop
	var init, cond, update string

	if stmt.Init != nil {
		init = fmt.Sprintf("%s := %s", goName(stmt.Init.Name), g.genExpr(stmt.Init.Init))
	}
	if stmt.Condition != nil {
		cond = g.genExpr(stmt.Condition)
	}
	if stmt.Update != nil {
		update = g.genExpr(stmt.Update)
	}

	g.writeln("for %s; %s; %s {", init, cond, update)
	g.indent++
	for _, inner := range stmt.Body.Stmts {
		g.genStmt(inner)
	}
	g.indent--
	g.writeln("}")
}

func (g *Generator) genForOfStmt(stmt *typed.ForOfStmt) {
	varName := goName(stmt.Variable.Name)
	iterable := g.genExpr(stmt.Iterable)

	// Check if we're iterating over a string (range returns runes, need to convert to string)
	iterType := types.Unwrap(stmt.Iterable.Type())
	if prim, ok := iterType.(*types.Primitive); ok && prim.Kind == types.KindString {
		// For string iteration, convert rune to string
		g.writeln("for _, _r := range %s {", iterable)
		g.indent++
		g.writeln("%s := string(_r)", varName)
		for _, inner := range stmt.Body.Stmts {
			g.genStmt(inner)
		}
		g.indent--
		g.writeln("}")
	} else {
		// Standard array iteration
		g.writeln("for _, %s := range %s {", varName, iterable)
		g.indent++
		for _, inner := range stmt.Body.Stmts {
			g.genStmt(inner)
		}
		g.indent--
		g.writeln("}")
	}
}

func (g *Generator) genSwitchStmt(stmt *typed.SwitchStmt) {
	g.writeln("switch %s {", g.genExpr(stmt.Discriminant))
	for _, c := range stmt.Cases {
		if c.Test != nil {
			g.writeln("case %s:", g.genExpr(c.Test))
		} else {
			g.writeln("default:")
		}
		g.indent++
		for _, inner := range c.Stmts {
			g.genStmt(inner)
		}
		g.indent--
	}
	g.writeln("}")
}

func (g *Generator) genTryStmt(stmt *typed.TryStmt) {
	// Generate try/catch using Go's defer/recover pattern
	// try { ... } catch (e) { ... }
	// becomes:
	// func() {
	//     defer func() {
	//         if r := recover(); r != nil {
	//             e := r
	//             // catch block
	//         }
	//     }()
	//     // try block
	// }()

	g.writeln("func() {")
	g.indent++
	g.writeln("defer func() {")
	g.indent++
	g.writeln("if r := recover(); r != nil {")
	g.indent++
	g.writeln("%s := r", goName(stmt.CatchParam.Name))
	g.writeln("_ = %s // prevent unused variable error", goName(stmt.CatchParam.Name))
	for _, inner := range stmt.CatchBlock.Stmts {
		g.genStmt(inner)
	}
	g.indent--
	g.writeln("}")
	g.indent--
	g.writeln("}()")
	for _, inner := range stmt.TryBlock.Stmts {
		g.genStmt(inner)
	}
	g.indent--
	g.writeln("}()")
}

func (g *Generator) genThrowStmt(stmt *typed.ThrowStmt) {
	g.writeln("panic(%s)", g.genExpr(stmt.Value))
}

// ----------------------------------------------------------------------------
// Expression Generation
// ----------------------------------------------------------------------------

func (g *Generator) genExpr(expr typed.Expr) string {
	if expr == nil {
		return "nil"
	}

	switch e := expr.(type) {
	case *typed.NumberLit:
		return fmt.Sprintf("%v", e.Value)

	case *typed.StringLit:
		return fmt.Sprintf("%q", e.Value)

	case *typed.TemplateLit:
		return g.genTemplateLit(e)

	case *typed.BoolLit:
		if e.Value {
			return "true"
		}
		return "false"

	case *typed.NullLit:
		return "nil"

	case *typed.RegexLit:
		return g.genRegexLit(e)

	case *typed.Ident:
		// Check if this is an imported Go package function
		if pkg, ok := g.goImportedNames[e.Name]; ok {
			return pkg + "." + e.Name
		}
		return goName(e.Name)

	case *typed.BinaryExpr:
		return g.genBinaryExpr(e)

	case *typed.UnaryExpr:
		return g.genUnaryExpr(e)

	case *typed.SpreadExpr:
		// For function call contexts, this returns arg...
		// For array literals, genArrayLit handles it specially
		return g.genExpr(e.Argument) + "..."

	case *typed.EnumMemberExpr:
		// Generate enum member constant name: EnumNameMemberName
		return exportName(e.EnumName) + exportName(e.MemberName)

	case *typed.CallExpr:
		return g.genCallExpr(e)

	case *typed.BuiltinCall:
		return g.genBuiltinCall(e)

	case *typed.IndexExpr:
		return g.genIndexExpr(e)

	case *typed.PropertyExpr:
		return g.genPropertyExpr(e)

	case *typed.ArrayLit:
		return g.genArrayLit(e)

	case *typed.ObjectLit:
		return g.genObjectLit(e)

	case *typed.FuncExpr:
		return g.genFuncExpr(e)

	case *typed.NewExpr:
		return g.genNewExpr(e)

	case *typed.ThisExpr:
		return "this"

	case *typed.SuperExpr:
		return g.genSuperExpr(e)

	case *typed.AssignExpr:
		return g.genAssignExpr(e)

	case *typed.CompoundAssignExpr:
		return g.genCompoundAssignExpr(e)

	case *typed.UpdateExpr:
		return g.genUpdateExpr(e)

	case *typed.MapLit:
		return g.genMapLit(e)

	case *typed.SetLit:
		return g.genSetLit(e)

	case *typed.MethodCallExpr:
		return g.genMethodCallExpr(e)

	case *typed.ConsoleCall:
		return g.genConsoleCall(e)

	case *typed.DateNewExpr:
		return g.genDateNewExpr(e)

	case *typed.DateMethodCall:
		return g.genDateMethodCall(e)

	case *typed.BuiltinObjectCall:
		return g.genBuiltinObjectCall(e)

	case *typed.BuiltinObjectConstant:
		return g.genBuiltinObjectConstant(e)

	case *typed.AwaitExpr:
		return g.genAwaitExpr(e)

	case *typed.PromiseMethodCall:
		return g.genPromiseMethodCall(e)
	}

	return "nil"
}

func (g *Generator) genBinaryExpr(expr *typed.BinaryExpr) string {
	left := g.genExpr(expr.Left)
	right := g.genExpr(expr.Right)

	// For arithmetic operations, add type assertions if operands are any
	leftType := types.Unwrap(expr.Left.Type())
	rightType := types.Unwrap(expr.Right.Type())
	leftIsAny := leftType.Equals(types.AnyType)
	rightIsAny := rightType.Equals(types.AnyType)

	// Helper to add type assertion for any type in arithmetic
	assertNumeric := func(operandExpr string, operandType types.Type, otherType types.Type) string {
		if !operandType.Equals(types.AnyType) {
			return operandExpr
		}
		// If the other operand is int, assert to int; if float, assert to float64
		if p, ok := otherType.(*types.Primitive); ok && p.Kind == types.KindInt {
			return fmt.Sprintf("%s.(int)", operandExpr)
		} else if p, ok := otherType.(*types.Primitive); ok && p.Kind == types.KindFloat {
			return fmt.Sprintf("%s.(float64)", operandExpr)
		}
		// Both are any - default to int for Y combinator and similar patterns
		return fmt.Sprintf("%s.(int)", operandExpr)
	}

	switch expr.Op {
	case "??":
		// Nullish coalescing - return proper typed result
		resultGoType := g.goType(expr.ExprType)
		resultType := types.Unwrap(expr.ExprType)

		// Check if left operand is a nullable type
		if nullable, ok := types.Unwrap(leftType).(*types.Nullable); ok {
			// Check if result type is also nullable (we keep the pointer)
			if _, resultIsNullable := resultType.(*types.Nullable); resultIsNullable {
				// Result is nullable, don't dereference
				return fmt.Sprintf("func() %s { if %s != nil { return %s }; return %s }()", resultGoType, left, left, right)
			}
			// For class types, don't dereference since classes are already pointers
			if _, isClass := nullable.Inner.(*types.Class); isClass {
				return fmt.Sprintf("func() %s { if %s != nil { return %s }; return %s }()", resultGoType, left, left, right)
			}
			// For primitive nullable types, dereference when non-nil
			return fmt.Sprintf("func() %s { if %s != nil { return *%s }; return %s }()", resultGoType, left, left, right)
		}
		// Left is already non-nullable (shouldn't happen in practice for ??)
		return fmt.Sprintf("func() %s { if %s != nil { return %s }; return %s }()", resultGoType, left, left, right)
	case "%":
		// Modulo - works directly on int in Go, need math.Mod for float64
		if lp, ok := leftType.(*types.Primitive); ok && lp.Kind == types.KindFloat {
			return fmt.Sprintf("math.Mod(%s, %s)", left, right)
		}
		if leftIsAny || rightIsAny {
			left = assertNumeric(left, leftType, rightType)
			right = assertNumeric(right, rightType, leftType)
		}
		return fmt.Sprintf("(%s %% %s)", left, right)
	case "+", "-", "*", "/":
		// For arithmetic, assert any types to numeric
		if leftIsAny || rightIsAny {
			left = assertNumeric(left, leftType, rightType)
			right = assertNumeric(right, rightType, leftType)
		}
		return fmt.Sprintf("(%s %s %s)", left, expr.Op, right)
	case "<", ">", "<=", ">=":
		// For comparisons, assert any types to numeric
		if leftIsAny || rightIsAny {
			left = assertNumeric(left, leftType, rightType)
			right = assertNumeric(right, rightType, leftType)
		}
		return fmt.Sprintf("(%s %s %s)", left, expr.Op, right)
	default:
		return fmt.Sprintf("(%s %s %s)", left, expr.Op, right)
	}
}

func (g *Generator) genUnaryExpr(expr *typed.UnaryExpr) string {
	operand := g.genExpr(expr.Operand)
	return fmt.Sprintf("(%s%s)", expr.Op, operand)
}

func (g *Generator) genCallExpr(expr *typed.CallExpr) string {
	callee := g.genExpr(expr.Callee)
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	// Handle optional chaining: fn?.()
	if expr.Optional {
		resultType := g.goType(expr.ExprType)
		argsStr := strings.Join(args, ", ")
		if resultType == "" || resultType == "interface{}" {
			return fmt.Sprintf("func() interface{} { if %s != nil { return gts_call(%s, %s) }; return nil }()",
				callee, callee, argsStr)
		}
		return fmt.Sprintf("func() %s { if %s != nil { return %s(%s) }; var zero %s; return zero }()",
			resultType, callee, callee, argsStr, resultType)
	}

	// Check if the callee is a generic function (interface{} type)
	// In that case, use gts_call for dynamic dispatch
	calleeType := types.Unwrap(expr.Callee.Type())
	if fn, ok := calleeType.(*types.Function); ok {
		isGeneric := len(fn.Params) > 0 && fn.Params[0].Type.Equals(types.AnyType) && fn.ReturnType.Equals(types.AnyType)
		if isGeneric {
			call := fmt.Sprintf("gts_call(%s, %s)", callee, strings.Join(args, ", "))
			// Add type assertion if the call result is used as a specific type
			return g.addTypeAssertion(call, expr.ExprType)
		}
	}
	// Also handle when the callee is directly an any type
	if prim, ok := calleeType.(*types.Primitive); ok && prim.Kind == types.KindAny {
		call := fmt.Sprintf("gts_call(%s, %s)", callee, strings.Join(args, ", "))
		// Add type assertion if the call result is used as a specific type
		return g.addTypeAssertion(call, expr.ExprType)
	}

	return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", "))
}

// addTypeAssertion wraps an interface{} expression with a type assertion if needed
func (g *Generator) addTypeAssertion(expr string, targetType types.Type) string {
	targetType = types.Unwrap(targetType)

	// If target type is any/interface{}, no assertion needed
	if prim, ok := targetType.(*types.Primitive); ok && prim.Kind == types.KindAny {
		return expr
	}

	// Add type assertion for specific types
	goType := g.goType(targetType)
	if goType != "interface{}" {
		return fmt.Sprintf("%s.(%s)", expr, goType)
	}
	return expr
}

func (g *Generator) genBuiltinCall(expr *typed.BuiltinCall) string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	switch expr.Name {
	case "println":
		return fmt.Sprintf("fmt.Println(%s)", strings.Join(args, ", "))
	case "print":
		return fmt.Sprintf("fmt.Print(%s)", strings.Join(args, ", "))
	case "len":
		return fmt.Sprintf("gts_len(%s)", args[0])
	case "push":
		// Push modifies the first argument (slice)
		// Get element type from the array to ensure proper type assertion
		arrType := types.Unwrap(expr.Args[0].Type())
		valueExpr := args[1]
		if arr, ok := arrType.(*types.Array); ok {
			// Check if value type is any but array element type is concrete
			valType := types.Unwrap(expr.Args[1].Type())
			if prim, ok := valType.(*types.Primitive); ok && prim.Kind == types.KindAny {
				elemGoType := g.goType(arr.Element)
				if elemGoType != "interface{}" {
					valueExpr = fmt.Sprintf("%s.(%s)", valueExpr, elemGoType)
				}
			}
		}
		return fmt.Sprintf("%s = append(%s, %s)", args[0], args[0], valueExpr)
	case "pop":
		// Pop returns and removes the last element
		// Return the proper element type
		elemType := "interface{}"
		arrType := types.Unwrap(expr.Args[0].Type())
		if arr, ok := arrType.(*types.Array); ok {
			elemType = g.goType(arr.Element)
		}
		return fmt.Sprintf("func() %s { n := len(%s); v := %s[n-1]; %s = %s[:n-1]; return v }()", elemType, args[0], args[0], args[0], args[0])
	case "typeof":
		return fmt.Sprintf("gts_typeof(%s)", args[0])
	case "tostring":
		return fmt.Sprintf("gts_tostring(%s)", args[0])
	case "toint":
		return fmt.Sprintf("gts_toint(%s)", args[0])
	case "parseInt":
		return fmt.Sprintf("gts_toint(%s)", args[0])
	case "tofloat":
		return fmt.Sprintf("gts_tofloat(%s)", args[0])
	case "sqrt":
		return fmt.Sprintf("math.Sqrt(%s)", args[0])
	case "floor":
		return fmt.Sprintf("math.Floor(%s)", args[0])
	case "ceil":
		return fmt.Sprintf("math.Ceil(%s)", args[0])
	case "abs":
		return fmt.Sprintf("math.Abs(%s)", args[0])
	// String methods
	case "split":
		return fmt.Sprintf("strings.Split(%s, %s)", args[0], args[1])
	case "join":
		return fmt.Sprintf("strings.Join(%s, %s)", args[0], args[1])
	case "replace":
		return fmt.Sprintf("strings.ReplaceAll(%s, %s, %s)", args[0], args[1], args[2])
	case "trim":
		return fmt.Sprintf("strings.TrimSpace(%s)", args[0])
	case "startsWith":
		return fmt.Sprintf("strings.HasPrefix(%s, %s)", args[0], args[1])
	case "endsWith":
		return fmt.Sprintf("strings.HasSuffix(%s, %s)", args[0], args[1])
	case "includes":
		return fmt.Sprintf("strings.Contains(%s, %s)", args[0], args[1])
	// Array methods
	case "map":
		goType := g.goType(expr.ExprType)
		return fmt.Sprintf("gts_map(%s, %s).(%s)", args[0], args[1], goType)
	case "filter":
		goType := g.goType(expr.ExprType)
		return fmt.Sprintf("gts_filter(%s, %s).(%s)", args[0], args[1], goType)
	case "reduce":
		goType := g.goType(expr.ExprType)
		return fmt.Sprintf("gts_reduce(%s, %s, %s).(%s)", args[0], args[1], args[2], goType)
	case "find":
		goType := g.goType(expr.ExprType)
		return fmt.Sprintf("gts_find(%s, %s).(%s)", args[0], args[1], goType)
	case "findIndex":
		return fmt.Sprintf("gts_findIndex(%s, %s)", args[0], args[1])
	case "some":
		return fmt.Sprintf("gts_some(%s, %s)", args[0], args[1])
	case "every":
		return fmt.Sprintf("gts_every(%s, %s)", args[0], args[1])
	// Global number functions
	case "isNaN":
		return fmt.Sprintf("math.IsNaN(%s)", args[0])
	case "isFinite":
		return fmt.Sprintf("(!math.IsInf(%s, 0) && !math.IsNaN(%s))", args[0], args[0])
	case "parseFloat":
		return fmt.Sprintf("func() float64 { v, _ := strconv.ParseFloat(%s, 64); return v }()", args[0])
	// Timer functions (event loop)
	case "setTimeout":
		return fmt.Sprintf("int(gts_setTimeout(func() { gts_call(%s) }, %s))", args[0], args[1])
	case "setInterval":
		return fmt.Sprintf("int(gts_setInterval(func() { gts_call(%s) }, %s))", args[0], args[1])
	case "clearTimeout":
		return fmt.Sprintf("gts_clearTimeout(int64(%s))", args[0])
	case "clearInterval":
		return fmt.Sprintf("gts_clearTimeout(int64(%s))", args[0]) // Same implementation as clearTimeout
	case "queueMicrotask":
		return fmt.Sprintf("gts_queueMicrotask(func() { gts_call(%s) })", args[0])
	default:
		return fmt.Sprintf("%s(%s)", expr.Name, strings.Join(args, ", "))
	}
}

func (g *Generator) genIndexExpr(expr *typed.IndexExpr) string {
	object := g.genExpr(expr.Object)
	index := g.genExpr(expr.Index)

	if expr.Optional {
		// Optional chaining: obj?.[index]
		// Generate: func() T { if obj != nil { return obj[int(index)] }; return nil }()
		resultType := g.goType(expr.ExprType)
		return fmt.Sprintf("func() %s { if %s != nil { return %s[int(%s)] }; var zero %s; return zero }()",
			resultType, object, object, index, resultType)
	}

	// Convert float64 index to int
	return fmt.Sprintf("%s[int(%s)]", object, index)
}

func (g *Generator) genPropertyExpr(expr *typed.PropertyExpr) string {
	object := g.genExpr(expr.Object)
	objType := types.Unwrap(expr.Object.Type())

	// Handle special properties for collection types
	switch objType.(type) {
	case *types.Array:
		if expr.Property == "length" {
			return fmt.Sprintf("len(%s)", object)
		}
	case *types.Map:
		if expr.Property == "size" {
			return fmt.Sprintf("len(%s)", object)
		}
	case *types.Set:
		if expr.Property == "size" {
			return fmt.Sprintf("len(%s)", object)
		}
	}

	if expr.Optional {
		// Optional chaining: obj?.prop
		// Generate: func() T { if obj != nil { return obj.Prop }; return nil }()
		resultType := g.goType(expr.ExprType)
		propAccess := fmt.Sprintf("%s.%s", object, exportName(expr.Property))

		// Check if we need to take address of the property (for nullable primitive types)
		// The result type will be a pointer for nullable types
		if strings.HasPrefix(resultType, "*") {
			// Check if it's a pointer to a class (e.g., *Person) vs pointer to primitive (e.g., *string)
			innerType := resultType[1:]
			// If the inner type starts with uppercase, it's likely a class
			if len(innerType) > 0 && innerType[0] >= 'A' && innerType[0] <= 'Z' {
				// It's a class pointer, can return directly
				return fmt.Sprintf("func() %s { if %s != nil { return %s }; return nil }()",
					resultType, object, propAccess)
			}
			// It's a pointer to a primitive type, need to take address
			return fmt.Sprintf("func() %s { if %s != nil { v := %s; return &v }; return nil }()",
				resultType, object, propAccess)
		}

		return fmt.Sprintf("func() %s { if %s != nil { return %s }; var zero %s; return zero }()",
			resultType, object, propAccess, resultType)
	}

	return fmt.Sprintf("%s.%s", object, exportName(expr.Property))
}

func (g *Generator) genArrayLit(expr *typed.ArrayLit) string {
	arrType := expr.ExprType.(*types.Array)
	goElemType := g.goType(arrType.Element)

	// Check if there are any spread expressions
	hasSpread := false
	for _, elem := range expr.Elements {
		if _, ok := elem.(*typed.SpreadExpr); ok {
			hasSpread = true
			break
		}
	}

	// Simple case: no spread expressions
	if !hasSpread {
		elements := make([]string, len(expr.Elements))
		for i, elem := range expr.Elements {
			elements[i] = g.genExpr(elem)
		}
		return fmt.Sprintf("[]%s{%s}", goElemType, strings.Join(elements, ", "))
	}

	// Complex case: handle spread expressions using append
	// Strategy: Build groups of consecutive non-spread elements, then use append
	var result string
	var currentGroup []string

	flushGroup := func() {
		if len(currentGroup) == 0 {
			return
		}
		groupLit := fmt.Sprintf("[]%s{%s}", goElemType, strings.Join(currentGroup, ", "))
		if result == "" {
			result = groupLit
		} else {
			// Append the group as individual elements
			result = fmt.Sprintf("append(%s, %s...)", result, groupLit)
		}
		currentGroup = nil
	}

	for _, elem := range expr.Elements {
		if spread, ok := elem.(*typed.SpreadExpr); ok {
			// Flush any accumulated non-spread elements first
			flushGroup()

			spreadArg := g.genExpr(spread.Argument)
			if result == "" {
				// First element is a spread - use append to create a new array (proper copy semantics)
				result = fmt.Sprintf("append([]%s{}, %s...)", goElemType, spreadArg)
			} else {
				// Append the spread
				result = fmt.Sprintf("append(%s, %s...)", result, spreadArg)
			}
		} else {
			// Accumulate non-spread elements
			currentGroup = append(currentGroup, g.genExpr(elem))
		}
	}

	// Flush any remaining non-spread elements
	if len(currentGroup) > 0 {
		if result == "" {
			// All elements are non-spread (shouldn't happen due to hasSpread check)
			result = fmt.Sprintf("[]%s{%s}", goElemType, strings.Join(currentGroup, ", "))
		} else {
			// Append remaining elements
			result = fmt.Sprintf("append(%s, %s)", result, strings.Join(currentGroup, ", "))
		}
	}

	// Handle empty array with only spreads
	if result == "" {
		result = fmt.Sprintf("[]%s{}", goElemType)
	}

	return result
}

func (g *Generator) genObjectLit(expr *typed.ObjectLit) string {
	// Generate as a struct literal or map
	// For now, use anonymous struct
	if len(expr.Properties) == 0 {
		return "struct{}{}"
	}

	var fields []string
	var values []string

	for _, prop := range expr.Properties {
		fields = append(fields, fmt.Sprintf("%s %s", exportName(prop.Key), g.goType(prop.Value.Type())))
		values = append(values, fmt.Sprintf("%s: %s", exportName(prop.Key), g.genExpr(prop.Value)))
	}

	return fmt.Sprintf("struct{%s}{%s}", strings.Join(fields, "; "), strings.Join(values, ", "))
}

func (g *Generator) genFuncExpr(expr *typed.FuncExpr) string {
	params := make([]string, len(expr.Params))
	for i, p := range expr.Params {
		params[i] = fmt.Sprintf("%s %s", goName(p.Name), g.goType(p.Type))
	}

	funcType := expr.ExprType.(*types.Function)
	returnType := g.goType(funcType.ReturnType)

	var buf bytes.Buffer
	if returnType == "" {
		buf.WriteString(fmt.Sprintf("func(%s) {\n", strings.Join(params, ", ")))
	} else {
		buf.WriteString(fmt.Sprintf("func(%s) %s {\n", strings.Join(params, ", "), returnType))
	}

	// Save and restore current return type for nested function expressions
	savedRetType := g.currentRetType
	g.currentRetType = funcType.ReturnType

	if expr.Body != nil {
		for _, stmt := range expr.Body.Stmts {
			buf.WriteString(g.genStmtToString(stmt))
		}
	} else if expr.BodyExpr != nil {
		retExpr := g.genExpr(expr.BodyExpr)
		// Add type assertion if body expression is any but return type is concrete
		// However, skip this for binary arithmetic expressions where we've already
		// asserted the operands (the generated code already produces the correct type)
		valType := types.Unwrap(expr.BodyExpr.Type())
		expectedType := types.Unwrap(funcType.ReturnType)
		needsAssertion := false
		if prim, ok := valType.(*types.Primitive); ok && prim.Kind == types.KindAny {
			expectedGoType := g.goType(expectedType)
			if expectedGoType != "interface{}" {
				// Check if this is a binary arithmetic expression - if so, we've already
				// handled type assertions for the operands, so the result is already typed
				if binExpr, ok := expr.BodyExpr.(*typed.BinaryExpr); ok {
					switch binExpr.Op {
					case "+", "-", "*", "/", "%", "<", ">", "<=", ">=":
						// Skip - we've already asserted operands in genBinaryExpr
						needsAssertion = false
					default:
						needsAssertion = true
					}
				} else {
					needsAssertion = true
				}
				if needsAssertion {
					retExpr = fmt.Sprintf("%s.(%s)", retExpr, expectedGoType)
				}
			}
		}
		buf.WriteString(fmt.Sprintf("return %s\n", retExpr))
	}

	g.currentRetType = savedRetType

	buf.WriteString("}")
	return buf.String()
}

func (g *Generator) genStmtToString(stmt typed.Stmt) string {
	var buf bytes.Buffer
	saved := g.buf
	g.buf = &buf
	g.genStmt(stmt)
	g.buf = saved
	return buf.String()
}

func (g *Generator) genNewExpr(expr *typed.NewExpr) string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}
	return fmt.Sprintf("New%s(%s)", exportName(expr.ClassName), strings.Join(args, ", "))
}

func (g *Generator) genSuperExpr(expr *typed.SuperExpr) string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}
	// Generate parent struct initialization: this.ParentName = *NewParent(args...)
	if g.currentClass != nil && g.currentClass.Super != "" {
		parentName := exportName(g.currentClass.Super)
		return fmt.Sprintf("this.%s = *New%s(%s)", parentName, parentName, strings.Join(args, ", "))
	}
	// Fallback - should not happen if type checker is working
	return fmt.Sprintf("/* super(%s) */", strings.Join(args, ", "))
}

func (g *Generator) genAssignExpr(expr *typed.AssignExpr) string {
	target := g.genExpr(expr.Target)
	value := g.genExpr(expr.Value)

	// Add type assertion if value type is any but target type is concrete
	valType := types.Unwrap(expr.Value.Type())
	targetType := types.Unwrap(expr.Target.Type())
	if prim, ok := valType.(*types.Primitive); ok && prim.Kind == types.KindAny {
		goType := g.goType(targetType)
		if goType != "interface{}" {
			value = fmt.Sprintf("%s.(%s)", value, goType)
		}
	}

	return fmt.Sprintf("%s = %s", target, value)
}

func (g *Generator) genCompoundAssignExpr(expr *typed.CompoundAssignExpr) string {
	target := g.genExpr(expr.Target)
	value := g.genExpr(expr.Value)
	return fmt.Sprintf("%s %s %s", target, expr.Op, value)
}

func (g *Generator) genUpdateExpr(expr *typed.UpdateExpr) string {
	operand := g.genExpr(expr.Operand)
	if expr.Prefix {
		return fmt.Sprintf("%s%s", expr.Op, operand)
	}
	return fmt.Sprintf("%s%s", operand, expr.Op)
}

func (g *Generator) genAwaitExpr(expr *typed.AwaitExpr) string {
	arg := g.genExpr(expr.Argument)
	resultType := g.goType(expr.ExprType)
	if resultType == "" {
		resultType = "interface{}"
	}
	return fmt.Sprintf("gts_await[%s](%s)", resultType, arg)
}

func (g *Generator) genPromiseMethodCall(expr *typed.PromiseMethodCall) string {
	obj := g.genExpr(expr.Object)
	callback := g.genExpr(expr.Callback)

	switch expr.Method {
	case "then":
		// p.Then(func(v T) interface{} { return callback(v) })
		return fmt.Sprintf("%s.Then(func(__v interface{}) interface{} { return gts_call(%s, __v) })", obj, callback)
	case "catch":
		// p.Catch(func(e error) interface{} { return callback(e) })
		return fmt.Sprintf("%s.Catch(func(__e error) interface{} { return gts_call(%s, __e) })", obj, callback)
	case "finally":
		// p.Finally(func() { callback() })
		return fmt.Sprintf("%s.Finally(func() { gts_call(%s) })", obj, callback)
	default:
		return fmt.Sprintf("%s.%s(%s)", obj, expr.Method, callback)
	}
}

func (g *Generator) genMapLit(expr *typed.MapLit) string {
	mapType := expr.ExprType.(*types.Map)
	goKeyType := g.goType(mapType.Key)
	goValType := g.goType(mapType.Value)

	if len(expr.Entries) == 0 {
		return fmt.Sprintf("make(map[%s]%s)", goKeyType, goValType)
	}

	entries := make([]string, len(expr.Entries))
	for i, entry := range expr.Entries {
		entries[i] = fmt.Sprintf("%s: %s", g.genExpr(entry.Key), g.genExpr(entry.Value))
	}

	return fmt.Sprintf("map[%s]%s{%s}", goKeyType, goValType, strings.Join(entries, ", "))
}

func (g *Generator) genSetLit(expr *typed.SetLit) string {
	setType := expr.ExprType.(*types.Set)
	goElemType := g.goType(setType.Element)
	// Sets are represented as map[T]struct{} in Go
	return fmt.Sprintf("make(map[%s]struct{})", goElemType)
}

func (g *Generator) genMethodCallExpr(expr *typed.MethodCallExpr) string {
	obj := g.genExpr(expr.Object)
	objType := types.Unwrap(expr.Object.Type())

	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	// Handle Map method calls
	if mapType, ok := objType.(*types.Map); ok {
		switch expr.Method {
		case "get":
			// m.get(key) => m[key]
			return fmt.Sprintf("%s[%s]", obj, args[0])

		case "set":
			// m.set(key, value) => func() map[K]V { m[key] = value; return m }()
			keyType := g.goType(mapType.Key)
			valType := g.goType(mapType.Value)
			return fmt.Sprintf("func() map[%s]%s { %s[%s] = %s; return %s }()", keyType, valType, obj, args[0], args[1], obj)

		case "has":
			// m.has(key) => func() bool { _, ok := m[key]; return ok }()
			return fmt.Sprintf("func() bool { _, ok := %s[%s]; return ok }()", obj, args[0])

		case "delete":
			// m.delete(key) => func() bool { _, ok := m[key]; if ok { delete(m, key) }; return ok }()
			return fmt.Sprintf("func() bool { _, ok := %s[%s]; if ok { delete(%s, %s) }; return ok }()", obj, args[0], obj, args[0])

		case "keys":
			keyType := g.goType(mapType.Key)
			return fmt.Sprintf("func() []%s { keys := make([]%s, 0, len(%s)); for k := range %s { keys = append(keys, k) }; return keys }()", keyType, keyType, obj, obj)

		case "values":
			valType := g.goType(mapType.Value)
			return fmt.Sprintf("func() []%s { vals := make([]%s, 0, len(%s)); for _, v := range %s { vals = append(vals, v) }; return vals }()", valType, valType, obj, obj)

		case "entries":
			// entries() => func() [][2]interface{} { ... }
			return fmt.Sprintf("func() []interface{} { entries := make([]interface{}, 0, len(%s)); for k, v := range %s { entries = append(entries, []interface{}{k, v}) }; return entries }()", obj, obj)

		case "clear":
			// clear() => func() { for k := range m { delete(m, k) } }()
			return fmt.Sprintf("func() { for k := range %s { delete(%s, k) } }()", obj, obj)

		case "forEach":
			// forEach(callback) => func() { for k, v := range m { callback(v, k) } }()
			return fmt.Sprintf("func() { for k, v := range %s { %s(v, k) } }()", obj, args[0])
		}
	}

	// Handle Set method calls (implemented as map[T]struct{})
	if setType, ok := objType.(*types.Set); ok {
		elemType := g.goType(setType.Element)
		switch expr.Method {
		case "add":
			// s.add(value) => func() map[T]struct{} { s[value] = struct{}{}; return s }()
			return fmt.Sprintf("func() map[%s]struct{} { %s[%s] = struct{}{}; return %s }()", elemType, obj, args[0], obj)

		case "has":
			// s.has(value) => func() bool { _, ok := s[value]; return ok }()
			return fmt.Sprintf("func() bool { _, ok := %s[%s]; return ok }()", obj, args[0])

		case "delete":
			// s.delete(value) => func() bool { _, ok := s[value]; if ok { delete(s, value) }; return ok }()
			return fmt.Sprintf("func() bool { _, ok := %s[%s]; if ok { delete(%s, %s) }; return ok }()", obj, args[0], obj, args[0])

		case "clear":
			// clear() => func() { for k := range s { delete(s, k) } }()
			return fmt.Sprintf("func() { for k := range %s { delete(%s, k) } }()", obj, obj)

		case "values":
			// values() => func() []T { vals := make([]T, 0, len(s)); for v := range s { vals = append(vals, v) }; return vals }()
			return fmt.Sprintf("func() []%s { vals := make([]%s, 0, len(%s)); for v := range %s { vals = append(vals, v) }; return vals }()", elemType, elemType, obj, obj)

		case "forEach":
			// forEach(callback) => func() { for v := range s { callback(v) } }()
			return fmt.Sprintf("func() { for v := range %s { %s(v) } }()", obj, args[0])
		}
	}

	// Handle Array method calls
	if arrType, ok := objType.(*types.Array); ok {
		elemType := g.goType(arrType.Element)
		switch expr.Method {
		case "push":
			// arr.push(values...) => func() int { arr = append(arr, values...); return len(arr) }()
			// Note: This requires the array to be a pointer or passed by reference for mutation
			// For now, generate code that works with slice semantics
			if len(args) == 1 {
				return fmt.Sprintf("func() int { %s = append(%s, %s); return len(%s) }()", obj, obj, args[0], obj)
			}
			return fmt.Sprintf("func() int { %s = append(%s, %s); return len(%s) }()", obj, obj, strings.Join(args, ", "), obj)

		case "pop":
			// arr.pop() => func() *T { if len(arr) == 0 { return nil }; v := arr[len(arr)-1]; arr = arr[:len(arr)-1]; return &v }()
			return fmt.Sprintf("func() *%s { if len(%s) == 0 { return nil }; v := %s[len(%s)-1]; %s = %s[:len(%s)-1]; return &v }()", elemType, obj, obj, obj, obj, obj, obj)

		case "shift":
			// arr.shift() => func() *T { if len(arr) == 0 { return nil }; v := arr[0]; arr = arr[1:]; return &v }()
			return fmt.Sprintf("func() *%s { if len(%s) == 0 { return nil }; v := %s[0]; %s = %s[1:]; return &v }()", elemType, obj, obj, obj, obj)

		case "unshift":
			// arr.unshift(values...) => func() int { arr = append([]T{values...}, arr...); return len(arr) }()
			if len(args) == 1 {
				return fmt.Sprintf("func() int { %s = append([]%s{%s}, %s...); return len(%s) }()", obj, elemType, args[0], obj, obj)
			}
			return fmt.Sprintf("func() int { %s = append([]%s{%s}, %s...); return len(%s) }()", obj, elemType, strings.Join(args, ", "), obj, obj)

		case "slice":
			// arr.slice(start?, end?) => arr[start:end]
			if len(args) == 0 {
				return fmt.Sprintf("append([]%s{}, %s...)", elemType, obj)
			} else if len(args) == 1 {
				return fmt.Sprintf("%s[%s:]", obj, args[0])
			}
			return fmt.Sprintf("%s[%s:%s]", obj, args[0], args[1])

		case "splice":
			// splice is complex - simplified implementation
			// For now, return a simple version that just returns the deleted elements
			if len(args) >= 2 {
				return fmt.Sprintf("func() []%s { start := %s; deleteCount := %s; deleted := %s[start:start+deleteCount]; %s = append(%s[:start], %s[start+deleteCount:]...); return deleted }()", elemType, args[0], args[1], obj, obj, obj, obj)
			}
			return fmt.Sprintf("%s[%s:]", obj, args[0])

		case "concat":
			// arr.concat(other...) => append(arr, other...)
			if len(args) == 1 {
				return fmt.Sprintf("append(%s, %s...)", obj, args[0])
			}
			// Multiple arrays
			result := obj
			for _, arg := range args {
				result = fmt.Sprintf("append(%s, %s...)", result, arg)
			}
			return result

		case "indexOf":
			// arr.indexOf(value, fromIndex?) => func() int { for i := start; i < len(arr); i++ { if arr[i] == value { return i } }; return -1 }()
			start := "0"
			if len(args) > 1 {
				start = args[1]
			}
			return fmt.Sprintf("func() int { for i := %s; i < len(%s); i++ { if %s[i] == %s { return i } }; return -1 }()", start, obj, obj, args[0])

		case "includes":
			// arr.includes(value, fromIndex?) => func() bool { for i := start; i < len(arr); i++ { if arr[i] == value { return true } }; return false }()
			start := "0"
			if len(args) > 1 {
				start = args[1]
			}
			return fmt.Sprintf("func() bool { for i := %s; i < len(%s); i++ { if %s[i] == %s { return true } }; return false }()", start, obj, obj, args[0])

		case "join":
			// arr.join(separator?) => strings.Join(...)  -- for string arrays only, otherwise needs conversion
			g.imports["strings"] = true
			sep := `","`
			if len(args) > 0 {
				sep = args[0]
			}
			// Need to handle different element types
			if elemType == "string" {
				return fmt.Sprintf("strings.Join(%s, %s)", obj, sep)
			}
			// For other types, need to convert to strings
			g.imports["fmt"] = true
			return fmt.Sprintf("func() string { strs := make([]string, len(%s)); for i, v := range %s { strs[i] = fmt.Sprint(v) }; return strings.Join(strs, %s) }()", obj, obj, sep)

		case "reverse":
			// arr.reverse() => func() []T { for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 { arr[i], arr[j] = arr[j], arr[i] }; return arr }()
			return fmt.Sprintf("func() []%s { for i, j := 0, len(%s)-1; i < j; i, j = i+1, j-1 { %s[i], %s[j] = %s[j], %s[i] }; return %s }()", elemType, obj, obj, obj, obj, obj, obj)

		case "sort":
			// arr.sort(compareFn?) => requires sort package
			g.imports["sort"] = true
			if len(args) == 0 {
				// Default sort based on type
				return fmt.Sprintf("func() []%s { sort.Slice(%s, func(i, j int) bool { return %s[i] < %s[j] }); return %s }()", elemType, obj, obj, obj, obj)
			}
			// Custom comparator
			return fmt.Sprintf("func() []%s { sort.Slice(%s, func(i, j int) bool { return %s(%s[i], %s[j]) < 0 }); return %s }()", elemType, obj, args[0], obj, obj, obj)

		case "map":
			// arr.map(callback) => func() []U { result := make([]U, 0); for _, v := range arr { result = append(result, callback(v)) }; return result }()
			resultElemType := "interface{}"
			if exprType, ok := expr.ExprType.(*types.Array); ok {
				resultElemType = g.goType(exprType.Element)
			}
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() []%s { result := make([]%s, 0); %s %s { result = append(result, %s) }; return result }()", resultElemType, resultElemType, forClause, obj, callbackCall)

		case "filter":
			// arr.filter(callback) => func() []T { result := make([]T, 0); for i, v := range arr { if callback(v, i) { result = append(result, v) } }; return result }()
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() []%s { result := make([]%s, 0); %s %s { if %s { result = append(result, v) } }; return result }()", elemType, elemType, forClause, obj, callbackCall)

		case "reduce":
			// Determine the accumulator type from the callback's return type (preferred)
			// or fall back to the initial value's type
			accType := "interface{}"
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && callbackType.ReturnType != nil {
				accType = g.goType(callbackType.ReturnType)
			} else if len(args) > 1 {
				accType = g.goType(expr.Args[1].Type())
			}
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 3 {
				callbackCall = fmt.Sprintf("%s(acc, v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(acc, v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() %s { acc := %s; %s %s { acc = %s }; return acc }()", accType, args[1], forClause, obj, callbackCall)

		case "forEach":
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() { %s %s { %s } }()", forClause, obj, callbackCall)

		case "find":
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() *%s { %s %s { if %s { return &v } }; return nil }()", elemType, forClause, obj, callbackCall)

		case "findIndex":
			// Check how many parameters the callback expects
			var callbackCall string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
			}
			// findIndex always needs the index, so always use 'for i, v := range'
			return fmt.Sprintf("func() int { for i, v := range %s { if %s { return i } }; return -1 }()", obj, callbackCall)

		case "some":
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() bool { %s %s { if %s { return true } }; return false }()", forClause, obj, callbackCall)

		case "every":
			// Check how many parameters the callback expects
			var callbackCall, forClause string
			if callbackType, ok := expr.Args[0].Type().(*types.Function); ok && len(callbackType.Params) >= 2 {
				callbackCall = fmt.Sprintf("%s(v, i)", args[0])
				forClause = "for i, v := range"
			} else {
				callbackCall = fmt.Sprintf("%s(v)", args[0])
				forClause = "for _, v := range"
			}
			return fmt.Sprintf("func() bool { %s %s { if !%s { return false } }; return true }()", forClause, obj, callbackCall)

		case "at":
			// arr.at(index) => supports negative indices
			return fmt.Sprintf("func() %s { arr := %s; i := %s; if i < 0 { i = len(arr) + i }; return arr[i] }()", elemType, obj, args[0])

		case "lastIndexOf":
			// arr.lastIndexOf(value) => search from end
			return fmt.Sprintf("func() int { for i := len(%s) - 1; i >= 0; i-- { if %s[i] == %s { return i } }; return -1 }()", obj, obj, args[0])

		case "fill":
			// arr.fill(value, start?, end?) => fill with value
			if len(args) == 1 {
				return fmt.Sprintf("func() []%s { for i := range %s { %s[i] = %s }; return %s }()", elemType, obj, obj, args[0], obj)
			} else if len(args) == 2 {
				return fmt.Sprintf("func() []%s { for i := %s; i < len(%s); i++ { %s[i] = %s }; return %s }()", elemType, args[1], obj, obj, args[0], obj)
			}
			return fmt.Sprintf("func() []%s { for i := %s; i < %s; i++ { %s[i] = %s }; return %s }()", elemType, args[1], args[2], obj, args[0], obj)

		case "copyWithin":
			// arr.copyWithin(target, start, end?) => copy within array
			if len(args) == 2 {
				return fmt.Sprintf("func() []%s { copy(%s[%s:], %s[%s:]); return %s }()", elemType, obj, args[0], obj, args[1], obj)
			}
			return fmt.Sprintf("func() []%s { copy(%s[%s:], %s[%s:%s]); return %s }()", elemType, obj, args[0], obj, args[1], args[2], obj)
		}
	}

	// Handle Number method calls
	if prim, ok := objType.(*types.Primitive); ok && (prim.Kind == types.KindInt || prim.Kind == types.KindFloat || prim.Kind == types.KindNumber) {
		switch expr.Method {
		case "toString":
			// num.toString() => gts_tostring(num)
			return fmt.Sprintf("gts_tostring(%s)", obj)
		}
	}

	// Handle String method calls
	if prim, ok := objType.(*types.Primitive); ok && prim.Kind == types.KindString {
		g.imports["strings"] = true
		switch expr.Method {
		case "split":
			// str.split(separator) => strings.Split(str, separator)
			return fmt.Sprintf("strings.Split(%s, %s)", obj, args[0])

		case "replace":
			// str.replace(old, new) => strings.Replace(str, old, new, -1)
			return fmt.Sprintf("strings.Replace(%s, %s, %s, -1)", obj, args[0], args[1])

		case "trim":
			// str.trim() => strings.TrimSpace(str)
			return fmt.Sprintf("strings.TrimSpace(%s)", obj)

		case "startsWith":
			// str.startsWith(prefix) => strings.HasPrefix(str, prefix)
			return fmt.Sprintf("strings.HasPrefix(%s, %s)", obj, args[0])

		case "endsWith":
			// str.endsWith(suffix) => strings.HasSuffix(str, suffix)
			return fmt.Sprintf("strings.HasSuffix(%s, %s)", obj, args[0])

		case "includes":
			// str.includes(substring) => strings.Contains(str, substring)
			return fmt.Sprintf("strings.Contains(%s, %s)", obj, args[0])

		case "toLowerCase":
			// str.toLowerCase() => strings.ToLower(str)
			return fmt.Sprintf("strings.ToLower(%s)", obj)

		case "toUpperCase":
			// str.toUpperCase() => strings.ToUpper(str)
			return fmt.Sprintf("strings.ToUpper(%s)", obj)

		case "substring":
			// str.substring(start, end?) => str[start:end]
			if len(args) == 1 {
				return fmt.Sprintf("%s[%s:]", obj, args[0])
			}
			return fmt.Sprintf("%s[%s:%s]", obj, args[0], args[1])

		case "charAt":
			// str.charAt(index) => string(str[index])
			return fmt.Sprintf("string(%s[%s])", obj, args[0])

		case "indexOf":
			// str.indexOf(substring) => strings.Index(str, substring)
			return fmt.Sprintf("strings.Index(%s, %s)", obj, args[0])

		case "charCodeAt":
			// str.charCodeAt(index) => int(str[index])
			return fmt.Sprintf("int(%s[%s])", obj, args[0])

		case "at":
			// str.at(index) => supports negative indices
			return fmt.Sprintf("func() string { s := %s; i := %s; if i < 0 { i = len(s) + i }; return string(s[i]) }()", obj, args[0])

		case "slice":
			// str.slice(start, end?) => with negative index support
			if len(args) == 1 {
				return fmt.Sprintf("func() string { s := %s; start := %s; if start < 0 { start = len(s) + start }; return s[start:] }()", obj, args[0])
			}
			return fmt.Sprintf("func() string { s := %s; start, end := %s, %s; if start < 0 { start = len(s) + start }; if end < 0 { end = len(s) + end }; return s[start:end] }()", obj, args[0], args[1])

		case "repeat":
			// str.repeat(count) => strings.Repeat(str, count)
			return fmt.Sprintf("strings.Repeat(%s, %s)", obj, args[0])

		case "padStart":
			// str.padStart(targetLength, padString)
			return fmt.Sprintf("func() string { s := %s; n := %s; pad := %s; for len(s) < n { s = pad + s }; return s[len(s)-n:] }()", obj, args[0], args[1])

		case "padEnd":
			// str.padEnd(targetLength, padString)
			return fmt.Sprintf("func() string { s := %s; n := %s; pad := %s; for len(s) < n { s = s + pad }; return s[:n] }()", obj, args[0], args[1])

		case "trimStart":
			// str.trimStart() => strings.TrimLeft(str, " \t\n\r")
			return fmt.Sprintf("strings.TrimLeft(%s, \" \\t\\n\\r\")", obj)

		case "trimEnd":
			// str.trimEnd() => strings.TrimRight(str, " \t\n\r")
			return fmt.Sprintf("strings.TrimRight(%s, \" \\t\\n\\r\")", obj)

		case "replaceAll":
			// str.replaceAll(old, new) => strings.ReplaceAll(str, old, new)
			return fmt.Sprintf("strings.ReplaceAll(%s, %s, %s)", obj, args[0], args[1])
		}
	}

	// Handle RegExp method calls
	if _, ok := objType.(*types.RegExp); ok {
		g.imports["regexp"] = true
		switch expr.Method {
		case "test":
			// re.test(str) => re.MatchString(str)
			return fmt.Sprintf("%s.MatchString(%s)", obj, args[0])

		case "exec":
			// re.exec(str) => re.FindStringSubmatch(str) (returns nil if no match)
			return fmt.Sprintf("%s.FindStringSubmatch(%s)", obj, args[0])
		}
	}

	// Fallback for other method calls (class methods, etc.)
	return fmt.Sprintf("%s.%s(%s)", obj, exportName(expr.Method), strings.Join(args, ", "))
}

func (g *Generator) genConsoleCall(expr *typed.ConsoleCall) string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	switch expr.Method {
	case "log":
		// console.log(...) => fmt.Println(...)
		return fmt.Sprintf("fmt.Println(%s)", strings.Join(args, ", "))
	default:
		// Default to Println for unknown methods
		return fmt.Sprintf("fmt.Println(%s)", strings.Join(args, ", "))
	}
}

func (g *Generator) genDateNewExpr(expr *typed.DateNewExpr) string {
	g.imports["time"] = true

	if len(expr.Args) == 0 {
		// new Date() => time.Now()
		return "time.Now()"
	} else if len(expr.Args) == 1 {
		// new Date(timestamp) => time.UnixMilli(timestamp)
		arg := g.genExpr(expr.Args[0])
		return fmt.Sprintf("time.UnixMilli(int64(%s))", arg)
	}

	// For other argument combinations, use current time as fallback
	return "time.Now()"
}

func (g *Generator) genDateMethodCall(expr *typed.DateMethodCall) string {
	g.imports["time"] = true

	obj := g.genExpr(expr.Object)
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	switch expr.Method {
	// Getter methods
	case "getTime":
		return fmt.Sprintf("float64(%s.UnixMilli())", obj)
	case "getFullYear":
		return fmt.Sprintf("%s.Year()", obj)
	case "getMonth":
		return fmt.Sprintf("(int(%s.Month()) - 1)", obj) // JS months are 0-indexed
	case "getDate":
		return fmt.Sprintf("%s.Day()", obj)
	case "getDay":
		return fmt.Sprintf("int(%s.Weekday())", obj)
	case "getHours":
		return fmt.Sprintf("%s.Hour()", obj)
	case "getMinutes":
		return fmt.Sprintf("%s.Minute()", obj)
	case "getSeconds":
		return fmt.Sprintf("%s.Second()", obj)
	case "getMilliseconds":
		return fmt.Sprintf("(%s.Nanosecond() / 1000000)", obj)

	// UTC getter methods
	case "getUTCFullYear":
		return fmt.Sprintf("%s.UTC().Year()", obj)
	case "getUTCMonth":
		return fmt.Sprintf("(int(%s.UTC().Month()) - 1)", obj)
	case "getUTCDate":
		return fmt.Sprintf("%s.UTC().Day()", obj)
	case "getUTCDay":
		return fmt.Sprintf("int(%s.UTC().Weekday())", obj)
	case "getUTCHours":
		return fmt.Sprintf("%s.UTC().Hour()", obj)
	case "getUTCMinutes":
		return fmt.Sprintf("%s.UTC().Minute()", obj)
	case "getUTCSeconds":
		return fmt.Sprintf("%s.UTC().Second()", obj)
	case "getUTCMilliseconds":
		return fmt.Sprintf("(%s.UTC().Nanosecond() / 1000000)", obj)

	// Setter methods - these mutate the time and return timestamp
	case "setTime":
		return fmt.Sprintf("func() float64 { %s = time.UnixMilli(int64(%s)); return float64(%s.UnixMilli()) }()", obj, args[0], obj)
	case "setFullYear":
		return g.genDateSetter(obj, args[0], obj+".Month()", obj+".Day()", obj+".Hour()", obj+".Minute()", obj+".Second()", obj+".Nanosecond()")
	case "setMonth":
		return g.genDateSetter(obj, obj+".Year()", "time.Month("+args[0]+"+1)", obj+".Day()", obj+".Hour()", obj+".Minute()", obj+".Second()", obj+".Nanosecond()")
	case "setDate":
		return g.genDateSetter(obj, obj+".Year()", obj+".Month()", args[0], obj+".Hour()", obj+".Minute()", obj+".Second()", obj+".Nanosecond()")
	case "setHours":
		return g.genDateSetter(obj, obj+".Year()", obj+".Month()", obj+".Day()", args[0], obj+".Minute()", obj+".Second()", obj+".Nanosecond()")
	case "setMinutes":
		return g.genDateSetter(obj, obj+".Year()", obj+".Month()", obj+".Day()", obj+".Hour()", args[0], obj+".Second()", obj+".Nanosecond()")
	case "setSeconds":
		return g.genDateSetter(obj, obj+".Year()", obj+".Month()", obj+".Day()", obj+".Hour()", obj+".Minute()", args[0], obj+".Nanosecond()")
	case "setMilliseconds":
		return g.genDateSetter(obj, obj+".Year()", obj+".Month()", obj+".Day()", obj+".Hour()", obj+".Minute()", obj+".Second()", args[0]+"*1000000")

	// String methods
	case "toString":
		return fmt.Sprintf("%s.String()", obj)
	case "toDateString":
		return fmt.Sprintf("%s.Format(\"Mon Jan 02 2006\")", obj)
	case "toTimeString":
		return fmt.Sprintf("%s.Format(\"15:04:05 MST\")", obj)
	case "toISOString":
		return fmt.Sprintf("%s.UTC().Format(time.RFC3339Nano)", obj)
	case "toJSON":
		return fmt.Sprintf("%s.UTC().Format(time.RFC3339Nano)", obj)
	case "toLocaleString":
		return fmt.Sprintf("%s.Format(\"1/2/2006, 3:04:05 PM\")", obj)
	case "toLocaleDateString":
		return fmt.Sprintf("%s.Format(\"1/2/2006\")", obj)
	case "toLocaleTimeString":
		return fmt.Sprintf("%s.Format(\"3:04:05 PM\")", obj)
	case "valueOf":
		return fmt.Sprintf("float64(%s.UnixMilli())", obj)

	default:
		return fmt.Sprintf("/* unknown Date method %s */ nil", expr.Method)
	}
}

// genDateSetter generates Go code for a Date setter that reconstructs a time.Date
// with one field replaced, then returns the timestamp.
func (g *Generator) genDateSetter(obj, year, month, day, hour, min, sec, nsec string) string {
	return fmt.Sprintf("func() float64 { %s = time.Date(%s, %s, %s, %s, %s, %s, %s, %s.Location()); return float64(%s.UnixMilli()) }()",
		obj, year, month, day, hour, min, sec, nsec, obj, obj)
}

func (g *Generator) genBuiltinObjectCall(expr *typed.BuiltinObjectCall) string {
	args := make([]string, len(expr.Args))
	for i, arg := range expr.Args {
		args[i] = g.genExpr(arg)
	}

	// Add required imports for this built-in object
	for _, imp := range typed.GetBuiltinImports(expr.Object) {
		g.imports[imp] = true
	}

	// Special handling for Object methods that need typed code generation
	if expr.Object == "Object" {
		return g.genObjectMethodCall(expr.Method, args, expr)
	}

	code, err := typed.GenerateBuiltinCall(expr.Object, expr.Method, args)
	if err != nil {
		// Fallback - shouldn't happen if type checking passed
		return fmt.Sprintf("/* %v */ nil", err)
	}
	return code
}

func (g *Generator) genObjectMethodCall(method string, args []string, expr *typed.BuiltinObjectCall) string {
	switch method {
	case "keys":
		// Object.keys returns []string with keys from map
		return fmt.Sprintf("func() []string { keys := make([]string, 0, len(%s)); for k := range %s { keys = append(keys, k) }; return keys }()", args[0], args[0])
	case "values":
		// Object.values returns slice of value type
		resultType := g.goType(expr.ExprType)
		// Get the element type for the slice
		if _, ok := expr.ExprType.(*types.Array); ok {
			return fmt.Sprintf("func() %s { vals := make(%s, 0, len(%s)); for _, v := range %s { vals = append(vals, v) }; return vals }()", resultType, resultType, args[0], args[0])
		}
		return fmt.Sprintf("func() %s { var vals %s; for _, v := range %s { vals = append(vals, v) }; return vals }()", resultType, resultType, args[0])
	case "assign":
		// Object.assign merges source into target and returns target
		resultType := g.goType(expr.ExprType)
		return fmt.Sprintf("func() %s { for k, v := range %s { %s[k] = v }; return %s }()", resultType, args[1], args[0], args[0])
	case "hasOwn":
		// Object.hasOwn checks if key exists
		return fmt.Sprintf("func() bool { _, ok := %s[%s]; return ok }()", args[0], args[1])
	default:
		return fmt.Sprintf("/* unknown Object.%s */ nil", method)
	}
}

func (g *Generator) genBuiltinObjectConstant(expr *typed.BuiltinObjectConstant) string {
	// Add required imports for this built-in object
	for _, imp := range typed.GetBuiltinImports(expr.Object) {
		g.imports[imp] = true
	}

	code, err := typed.GenerateBuiltinConstant(expr.Object, expr.Name)
	if err != nil {
		// Fallback - shouldn't happen if type checking passed
		return fmt.Sprintf("/* %v */ 0", err)
	}
	return code
}

// ----------------------------------------------------------------------------
// Type Mapping
// ----------------------------------------------------------------------------

func (g *Generator) goType(t types.Type) string {
	if t == nil {
		return ""
	}

	t = types.Unwrap(t)

	switch typ := t.(type) {
	case *types.Primitive:
		switch typ.Kind {
		case types.KindInt:
			return "int"
		case types.KindFloat:
			return "float64"
		case types.KindNumber:
			return "float64"
		case types.KindString:
			return "string"
		case types.KindBoolean:
			return "bool"
		case types.KindVoid:
			return ""
		case types.KindNull:
			return "interface{}"
		case types.KindAny:
			return "interface{}"
		case types.KindNever:
			return ""
		}

	case *types.Literal:
		// For literal types, use the base primitive type
		switch typ.Kind {
		case types.KindInt:
			return "int"
		case types.KindFloat:
			return "float64"
		case types.KindNumber:
			return "float64"
		case types.KindString:
			return "string"
		case types.KindBoolean:
			return "bool"
		default:
			return "interface{}"
		}

	case *types.Tuple:
		// Generate an anonymous struct with numbered fields
		var fields []string
		for i, elem := range typ.Elements {
			fields = append(fields, fmt.Sprintf("T%d %s", i, g.goType(elem)))
		}
		// Note: Rest elements are not yet fully supported in codegen
		// For now, we ignore the rest element
		return fmt.Sprintf("struct{%s}", strings.Join(fields, "; "))

	case *types.Array:
		return "[]" + g.goType(typ.Element)

	case *types.Map:
		return fmt.Sprintf("map[%s]%s", g.goType(typ.Key), g.goType(typ.Value))

	case *types.Set:
		return fmt.Sprintf("map[%s]struct{}", g.goType(typ.Element))

	case *types.Promise:
		innerType := g.goType(typ.Value)
		if innerType == "" {
			innerType = "interface{}" // void becomes interface{}
		}
		return fmt.Sprintf("*GTS_Promise[%s]", innerType)

	case *types.Enum:
		return exportName(typ.Name)

	case *types.RegExp:
		g.imports["regexp"] = true
		return "*regexp.Regexp"

	case *types.Date:
		g.imports["time"] = true
		return "time.Time"

	case *types.Object:
		// Anonymous struct
		var fields []string
		for name, prop := range typ.Properties {
			fields = append(fields, fmt.Sprintf("%s %s", exportName(name), g.goType(prop.Type)))
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, "; "))

	case *types.Function:
		// Check if this is a generic function type (all any params)
		// If so, use interface{} to hold any function
		isGeneric := len(typ.Params) > 0 && typ.Params[0].Type.Equals(types.AnyType) && typ.ReturnType.Equals(types.AnyType)
		if isGeneric {
			return "interface{}"
		}
		params := make([]string, len(typ.Params))
		for i, p := range typ.Params {
			params[i] = g.goType(p.Type)
		}
		retType := g.goType(typ.ReturnType)
		if retType == "" {
			return fmt.Sprintf("func(%s)", strings.Join(params, ", "))
		}
		return fmt.Sprintf("func(%s) %s", strings.Join(params, ", "), retType)

	case *types.Nullable:
		// Use pointer for nullable, but class types and arrays are already nullable in Go
		if _, isClass := typ.Inner.(*types.Class); isClass {
			// Class types are already *ClassName, so nullable class is just *ClassName (can be nil)
			return g.goType(typ.Inner)
		}
		if _, isArray := typ.Inner.(*types.Array); isArray {
			// Slices in Go can already be nil, no need for pointer
			return g.goType(typ.Inner)
		}
		if _, isMap := typ.Inner.(*types.Map); isMap {
			// Maps in Go can already be nil, no need for pointer
			return g.goType(typ.Inner)
		}
		inner := g.goType(typ.Inner)
		return "*" + inner

	case *types.Union:
		// Go doesn't have native union types, so we use interface{}
		// At runtime, the value can be any of the union member types
		return "interface{}"

	case *types.Intersection:
		// Try to merge as object first
		if merged := typ.MergeAsObject(); merged != nil {
			return g.goType(merged)
		}
		// Otherwise use interface{} for non-object intersections
		return "interface{}"

	case *types.Class:
		// Check if this class was instantiated from a generic class
		if typ.GenericBaseName != "" && len(typ.TypeArgs) > 0 {
			typeArgs := make([]string, len(typ.TypeArgs))
			for i, ta := range typ.TypeArgs {
				typeArgs[i] = g.goType(ta)
			}
			return fmt.Sprintf("*%s[%s]", exportName(typ.GenericBaseName), strings.Join(typeArgs, ", "))
		}
		return "*" + exportName(typ.Name)

	case *types.Interface:
		return exportName(typ.Name)

	case *types.TypeParameter:
		// Type parameters are used directly by name in Go generics
		return typ.Name

	case *types.GenericFunction:
		// For generic function types, use interface{} as they need special handling
		return "interface{}"

	case *types.GenericClass:
		// For generic class types without type arguments, use the name with type parameter placeholder
		return "*" + exportName(typ.Name)
	}

	return "interface{}"
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func (g *Generator) write(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) writeln(format string, args ...interface{}) {
	for i := 0; i < g.indent; i++ {
		g.buf.WriteString("\t")
	}
	fmt.Fprintf(g.buf, format, args...)
	g.buf.WriteString("\n")
}

// genTemplateLit generates Go code for a template literal.
// Template literals like `Hello, ${name}!` become fmt.Sprintf("Hello, %v!", name)
func (g *Generator) genTemplateLit(e *typed.TemplateLit) string {
	// If there are no expressions, it's just a simple string
	if len(e.Expressions) == 0 {
		return fmt.Sprintf("%q", e.Parts[0])
	}

	// Build the format string with %v placeholders
	var formatParts []string
	for i, part := range e.Parts {
		// Escape % characters in the static parts
		escaped := strings.ReplaceAll(part, "%", "%%")
		formatParts = append(formatParts, escaped)
		if i < len(e.Expressions) {
			formatParts = append(formatParts, "%v")
		}
	}
	formatStr := strings.Join(formatParts, "")

	// Generate the expression arguments
	args := make([]string, len(e.Expressions))
	for i, expr := range e.Expressions {
		args[i] = g.genExpr(expr)
	}

	// Mark fmt as used
	g.imports["fmt"] = true

	return fmt.Sprintf("fmt.Sprintf(%q, %s)", formatStr, strings.Join(args, ", "))
}

func (g *Generator) genRegexLit(e *typed.RegexLit) string {
	g.imports["regexp"] = true

	// Convert TypeScript flags to Go embedded flags
	// i -> (?i), m -> (?m), s -> (?s)
	pattern := e.Pattern
	goFlags := ""

	for _, f := range e.Flags {
		switch f {
		case 'i':
			goFlags += "i"
		case 'm':
			goFlags += "m"
		case 's':
			goFlags += "s"
		// 'g' (global) is handled at match time in Go, not in pattern
		// 'u' (unicode) is default in Go RE2
		// 'y' (sticky) has no direct equivalent
		}
	}

	if goFlags != "" {
		pattern = "(?" + goFlags + ")" + pattern
	}

	return fmt.Sprintf("regexp.MustCompile(%q)", pattern)
}

// goName converts a GTS name to a valid Go identifier.
func goName(name string) string {
	// Avoid Go reserved words
	switch name {
	case "type", "func", "var", "const", "package", "import", "return",
		"if", "else", "for", "range", "switch", "case", "default",
		"break", "continue", "goto", "fallthrough", "defer", "go",
		"chan", "map", "struct", "interface", "select":
		return name + "_"
	}
	return name
}

// exportName converts to Go exported name (capitalize first letter).
func exportName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
