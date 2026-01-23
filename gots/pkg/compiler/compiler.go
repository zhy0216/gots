// Package compiler compiles AST to bytecode.
package compiler

import (
	"fmt"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/bytecode"
	"github.com/zhy0216/quickts/gots/pkg/token"
)

// FunctionType indicates what type of function is being compiled.
type FunctionType int

const (
	TYPE_SCRIPT   FunctionType = iota // Top-level script
	TYPE_FUNCTION                     // Regular function
	TYPE_METHOD                       // Class method
)

// Local represents a local variable in the current scope.
type Local struct {
	Name       string
	Depth      int  // Scope depth (0 = global)
	IsCaptured bool
}

// Upvalue represents a captured variable from an enclosing scope.
type Upvalue struct {
	Index   int  // Index in enclosing function's locals or upvalues
	IsLocal bool // true if capturing from immediately enclosing scope
}

// ObjFunction represents a compiled function.
type ObjFunction struct {
	Name         string
	Arity        int
	UpvalueCount int
	Chunk        *bytecode.Chunk
}

// CompilerState holds the state for compiling a single function.
type CompilerState struct {
	function   *ObjFunction
	funcType   FunctionType
	locals     []Local
	upvalues   []Upvalue
	scopeDepth int
	enclosing  *CompilerState
}

// Compiler compiles AST to bytecode.
type Compiler struct {
	current *CompilerState
	globals map[string]int
}

// New creates a new compiler.
func New() *Compiler {
	c := &Compiler{
		globals: make(map[string]int),
	}
	c.initCompilerState(TYPE_SCRIPT, "")
	return c
}

// initCompilerState initializes a new compiler state for a function.
func (c *Compiler) initCompilerState(funcType FunctionType, name string) {
	state := &CompilerState{
		function: &ObjFunction{
			Name:  name,
			Arity: 0,
			Chunk: bytecode.NewChunk(),
		},
		funcType:   funcType,
		locals:     make([]Local, 0, 256),
		upvalues:   []Upvalue{},
		scopeDepth: 0,
		enclosing:  c.current,
	}
	c.current = state

	// Reserve slot 0 for the function itself (or "this" for methods)
	local := Local{
		Name:       "",
		Depth:      0,
		IsCaptured: false,
	}
	c.current.locals = append(c.current.locals, local)
}

// chunk returns the current chunk being compiled.
func (c *Compiler) chunk() *bytecode.Chunk {
	return c.current.function.Chunk
}

// endCompiler finishes compilation and returns the compiled function.
func (c *Compiler) endCompiler() *ObjFunction {
	c.emitReturn()
	fn := c.current.function
	fn.UpvalueCount = len(c.current.upvalues)

	// Restore enclosing compiler state
	if c.current.enclosing != nil {
		c.current = c.current.enclosing
	}

	return fn
}

// Compile compiles a program to bytecode.
func (c *Compiler) Compile(program *ast.Program) (*bytecode.Chunk, error) {
	for _, stmt := range program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return nil, err
		}
	}

	fn := c.endCompiler()
	return fn.Chunk, nil
}

// CompileFunction compiles a program and returns the function object.
func (c *Compiler) CompileFunction(program *ast.Program) (*ObjFunction, error) {
	for _, stmt := range program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return nil, err
		}
	}

	return c.endCompiler(), nil
}

// Scope management

// beginScope starts a new lexical scope.
func (c *Compiler) beginScope() {
	c.current.scopeDepth++
}

// endScope ends the current lexical scope.
func (c *Compiler) endScope(line int) {
	c.current.scopeDepth--

	// Pop all locals in the scope
	for len(c.current.locals) > 0 &&
		c.current.locals[len(c.current.locals)-1].Depth > c.current.scopeDepth {

		local := c.current.locals[len(c.current.locals)-1]
		if local.IsCaptured {
			c.emitByte(byte(bytecode.OP_CLOSE_UPVALUE), line)
		} else {
			c.emitByte(byte(bytecode.OP_POP), line)
		}
		c.current.locals = c.current.locals[:len(c.current.locals)-1]
	}
}

// Local variable management

// addLocal adds a local variable to the current scope.
func (c *Compiler) addLocal(name string) error {
	if len(c.current.locals) >= 256 {
		return fmt.Errorf("too many local variables in function")
	}

	local := Local{
		Name:       name,
		Depth:      -1,
		IsCaptured: false,
	}
	c.current.locals = append(c.current.locals, local)
	return nil
}

// markInitialized marks the most recent local as initialized.
func (c *Compiler) markInitialized() {
	if c.current.scopeDepth == 0 {
		return // Global variables don't need this
	}
	c.current.locals[len(c.current.locals)-1].Depth = c.current.scopeDepth
}

// resolveLocal resolves a local variable by name.
// Returns -1 if not found, -2 if used in own initializer.
func (c *Compiler) resolveLocal(name string) int {
	for i := len(c.current.locals) - 1; i >= 0; i-- {
		local := &c.current.locals[i]
		if local.Name == name {
			if local.Depth == -1 {
				return -2
			}
			return i
		}
	}
	return -1
}

// resolveUpvalue resolves an upvalue (captured variable from enclosing scope).
func (c *Compiler) resolveUpvalue(name string) int {
	if c.current.enclosing == nil {
		return -1
	}

	enclosing := c.current.enclosing

	// Try to find as a local in the immediately enclosing scope
	localIdx := c.resolveLocalIn(enclosing, name)
	if localIdx != -1 {
		enclosing.locals[localIdx].IsCaptured = true
		return c.addUpvalue(localIdx, true)
	}

	// Try to find as an upvalue in the enclosing scope
	upvalueIdx := c.resolveUpvalueIn(enclosing, name)
	if upvalueIdx != -1 {
		return c.addUpvalue(upvalueIdx, false)
	}

	return -1
}

// resolveLocalIn resolves a local variable in a specific compiler state.
func (c *Compiler) resolveLocalIn(state *CompilerState, name string) int {
	for i := len(state.locals) - 1; i >= 0; i-- {
		local := &state.locals[i]
		if local.Name == name {
			return i
		}
	}
	return -1
}

// resolveUpvalueIn resolves an upvalue in a specific compiler state.
func (c *Compiler) resolveUpvalueIn(state *CompilerState, name string) int {
	if state.enclosing == nil {
		return -1
	}

	// Check if it's a local in the enclosing scope
	enclosing := state.enclosing
	localIdx := c.resolveLocalIn(enclosing, name)
	if localIdx != -1 {
		enclosing.locals[localIdx].IsCaptured = true
		return c.addUpvalueIn(state, localIdx, true)
	}

	// Check if it's an upvalue in the enclosing scope
	upvalueIdx := c.resolveUpvalueIn(enclosing, name)
	if upvalueIdx != -1 {
		return c.addUpvalueIn(state, upvalueIdx, false)
	}

	return -1
}

// addUpvalue adds an upvalue to the current function.
func (c *Compiler) addUpvalue(index int, isLocal bool) int {
	return c.addUpvalueIn(c.current, index, isLocal)
}

// addUpvalueIn adds an upvalue to a specific compiler state.
func (c *Compiler) addUpvalueIn(state *CompilerState, index int, isLocal bool) int {
	// Check if we already have this upvalue
	for i, uv := range state.upvalues {
		if uv.Index == index && uv.IsLocal == isLocal {
			return i
		}
	}

	if len(state.upvalues) >= 256 {
		panic("too many closure variables in function")
	}

	state.upvalues = append(state.upvalues, Upvalue{Index: index, IsLocal: isLocal})
	return len(state.upvalues) - 1
}

// emitReturn emits a return instruction.
func (c *Compiler) emitReturn() {
	c.emitByte(byte(bytecode.OP_NULL), 0) // Return null by default
	c.emitByte(byte(bytecode.OP_RETURN), 0)
}

func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		if err := c.compileExpression(s.Expr); err != nil {
			return err
		}
		if !c.lastWasVoidBuiltin() {
			c.emitByte(byte(bytecode.OP_POP), s.Token.Line)
		}
		return nil

	case *ast.VarDecl:
		return c.compileVarDecl(s)

	case *ast.FuncDecl:
		return c.compileFuncDecl(s)

	case *ast.ReturnStmt:
		return c.compileReturnStmt(s)

	case *ast.IfStmt:
		return c.compileIfStmt(s)

	case *ast.WhileStmt:
		return c.compileWhileStmt(s)

	case *ast.Block:
		return c.compileBlock(s)

	case *ast.ClassDecl:
		return c.compileClassDecl(s)

	case *ast.TypeAliasDecl:
		return nil

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (c *Compiler) compileVarDecl(v *ast.VarDecl) error {
	// If we're in local scope, declare as local
	if c.current.scopeDepth > 0 {
		// Check for redeclaration in same scope
		for i := len(c.current.locals) - 1; i >= 0; i-- {
			local := &c.current.locals[i]
			if local.Depth != -1 && local.Depth < c.current.scopeDepth {
				break
			}
			if local.Name == v.Name {
				return fmt.Errorf("variable '%s' already declared in this scope", v.Name)
			}
		}

		// Add local variable (but don't mark as initialized yet)
		if err := c.addLocal(v.Name); err != nil {
			return err
		}
	}

	// Compile the initializer value
	if v.Value != nil {
		if err := c.compileExpression(v.Value); err != nil {
			return err
		}
	} else {
		// No initializer, default to null
		c.emitByte(byte(bytecode.OP_NULL), v.Token.Line)
	}

	if c.current.scopeDepth > 0 {
		// Local variable - just mark as initialized (value is already on stack)
		c.markInitialized()
	} else {
		// Global variable
		nameIdx := c.addGlobalVariable(v.Name)
		c.emitByte(byte(bytecode.OP_SET_GLOBAL), v.Token.Line)
		c.emitU16(uint16(nameIdx), v.Token.Line)
		// Pop the value (global variable declarations are statements)
		c.emitByte(byte(bytecode.OP_POP), v.Token.Line)
	}

	return nil
}

func (c *Compiler) compileFuncDecl(f *ast.FuncDecl) error {
	// Declare the function name (as a global or local)
	if c.current.scopeDepth > 0 {
		if err := c.addLocal(f.Name); err != nil {
			return err
		}
	}

	// Compile the function body
	if err := c.compileFunction(f.Name, f.Params, f.Body, TYPE_FUNCTION); err != nil {
		return err
	}

	if c.current.scopeDepth > 0 {
		// Local function - just mark as initialized
		c.markInitialized()
	} else {
		// Global function
		nameIdx := c.addGlobalVariable(f.Name)
		c.emitByte(byte(bytecode.OP_SET_GLOBAL), f.Token.Line)
		c.emitU16(uint16(nameIdx), f.Token.Line)
		c.emitByte(byte(bytecode.OP_POP), f.Token.Line)
	}

	return nil
}

func (c *Compiler) compileFunction(name string, params []*ast.Parameter, body *ast.Block, funcType FunctionType) error {
	// Start a new compiler state for this function
	c.initCompilerState(funcType, name)
	c.beginScope()

	// Bind parameters as local variables
	c.current.function.Arity = len(params)
	for _, param := range params {
		if err := c.addLocal(param.Name); err != nil {
			return err
		}
		c.markInitialized()
	}

	// Compile the function body
	for _, stmt := range body.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}

	// Save the upvalues before ending the compiler (we need them for OP_CLOSURE)
	upvalues := c.current.upvalues

	// End the function (this restores the enclosing state)
	fn := c.endCompiler()

	// Emit closure instruction in the enclosing function
	fnIdx := c.chunk().AddConstant(fn)
	c.emitByte(byte(bytecode.OP_CLOSURE), body.Token.Line)
	c.emitU16(uint16(fnIdx), body.Token.Line)

	// Emit upvalue descriptors (from the compiled function, not current state)
	for i := 0; i < fn.UpvalueCount; i++ {
		uv := upvalues[i]
		if uv.IsLocal {
			c.emitByte(1, body.Token.Line)
		} else {
			c.emitByte(0, body.Token.Line)
		}
		c.emitByte(byte(uv.Index), body.Token.Line)
	}

	return nil
}

func (c *Compiler) compileReturnStmt(r *ast.ReturnStmt) error {
	if c.current.funcType == TYPE_SCRIPT {
		return fmt.Errorf("cannot return from top-level code")
	}

	if r.Value != nil {
		if err := c.compileExpression(r.Value); err != nil {
			return err
		}
	} else {
		c.emitByte(byte(bytecode.OP_NULL), r.Token.Line)
	}

	c.emitByte(byte(bytecode.OP_RETURN), r.Token.Line)
	return nil
}

func (c *Compiler) addGlobalVariable(name string) int {
	idx := c.chunk().AddConstant(name)
	if c.current.funcType == TYPE_SCRIPT {
		c.globals[name] = idx
	}
	return idx
}

func (c *Compiler) getGlobalVariable(name string) (int, bool) {
	if _, exists := c.globals[name]; !exists {
		return -1, false
	}
	idx := c.chunk().AddConstant(name)
	return idx, true
}

func (c *Compiler) compileIfStmt(i *ast.IfStmt) error {
	line := i.Token.Line

	// Compile the condition
	if err := c.compileExpression(i.Condition); err != nil {
		return err
	}

	// Emit OP_JUMP_IF_FALSE with placeholder
	jumpIfFalse := c.emitJump(bytecode.OP_JUMP_IF_FALSE, line)

	// Pop the condition value (true case)
	c.emitByte(byte(bytecode.OP_POP), line)

	// Compile the consequence (then branch)
	if err := c.compileBlock(i.Consequence); err != nil {
		return err
	}

	// Always emit jump to skip the false-case pop (and else body if present)
	jumpOver := c.emitJump(bytecode.OP_JUMP, line)

	// Patch the jump-if-false to here
	c.patchJump(jumpIfFalse)

	// Pop the condition value (false case)
	c.emitByte(byte(bytecode.OP_POP), line)

	if i.Alternative != nil {
		// Compile the alternative (else branch)
		if err := c.compileStatement(i.Alternative); err != nil {
			return err
		}
	}

	// Patch the jump-over to here (end of if statement)
	c.patchJump(jumpOver)

	return nil
}

func (c *Compiler) compileBlock(b *ast.Block) error {
	c.beginScope()
	for _, stmt := range b.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	c.endScope(b.Token.Line)
	return nil
}

func (c *Compiler) compileWhileStmt(w *ast.WhileStmt) error {
	line := w.Token.Line

	// Remember start of loop (for jumping back)
	loopStart := c.chunk().Count()

	// Compile the condition
	if err := c.compileExpression(w.Condition); err != nil {
		return err
	}

	// Emit OP_JUMP_IF_FALSE to exit loop
	exitJump := c.emitJump(bytecode.OP_JUMP_IF_FALSE, line)

	// Pop the condition value (true case - entering loop body)
	c.emitByte(byte(bytecode.OP_POP), line)

	// Compile the body (with scope)
	c.beginScope()
	for _, stmt := range w.Body.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	c.endScope(line)

	// Emit jump back to loop start
	c.emitLoop(loopStart, line)

	// Patch the exit jump to here
	c.patchJump(exitJump)

	// Pop the condition value (false case - exiting loop)
	c.emitByte(byte(bytecode.OP_POP), line)

	return nil
}

// emitLoop emits a backward jump to loopStart.
func (c *Compiler) emitLoop(loopStart int, line int) {
	c.emitByte(byte(bytecode.OP_JUMP_BACK), line)

	offset := c.chunk().Count() - loopStart + 2
	if offset > 65535 {
		panic("loop body too large")
	}

	c.emitByte(byte(offset>>8), line)
	c.emitByte(byte(offset), line)
}

// emitJump emits a jump instruction with a placeholder offset.
// Returns the position to patch.
func (c *Compiler) emitJump(op bytecode.OpCode, line int) int {
	c.emitByte(byte(op), line)
	c.emitByte(0xff, line)
	c.emitByte(0xff, line)
	return c.chunk().Count() - 2
}

// patchJump patches a previously emitted jump to jump to the current position.
func (c *Compiler) patchJump(offset int) {
	jump := c.chunk().Count() - offset - 2
	if jump > 65535 {
		panic("jump too large")
	}

	c.chunk().Code[offset] = byte(jump >> 8)
	c.chunk().Code[offset+1] = byte(jump)
}

// lastWasVoidBuiltin checks if the last emitted opcode was a void builtin.
func (c *Compiler) lastWasVoidBuiltin() bool {
	code := c.chunk().Code
	if len(code) < 3 {
		return false
	}
	if bytecode.OpCode(code[len(code)-3]) != bytecode.OP_BUILTIN {
		return false
	}
	builtinID := code[len(code)-2]
	return builtinID == bytecode.BUILTIN_PRINTLN || builtinID == bytecode.BUILTIN_PRINT
}

func (c *Compiler) compileExpression(expr ast.Expression) error {
	switch e := expr.(type) {
	case *ast.NumberLiteral:
		return c.compileNumber(e)

	case *ast.StringLiteral:
		return c.compileString(e)

	case *ast.BoolLiteral:
		return c.compileBoolean(e)

	case *ast.NullLiteral:
		c.emitByte(byte(bytecode.OP_NULL), e.Token.Line)
		return nil

	case *ast.BinaryExpr:
		return c.compileBinary(e)

	case *ast.UnaryExpr:
		return c.compileUnary(e)

	case *ast.CallExpr:
		return c.compileCall(e)

	case *ast.Identifier:
		return c.compileIdentifier(e)

	case *ast.AssignExpr:
		return c.compileAssignment(e)

	case *ast.FunctionExpr:
		return c.compileFunctionExpr(e)

	case *ast.ArrayLiteral:
		return c.compileArrayLiteral(e)

	case *ast.ObjectLiteral:
		return c.compileObjectLiteral(e)

	case *ast.IndexExpr:
		return c.compileIndexExpr(e)

	case *ast.PropertyExpr:
		return c.compilePropertyExpr(e)

	case *ast.NewExpr:
		return c.compileNewExpr(e)

	case *ast.ThisExpr:
		return c.compileThisExpr(e)

	case *ast.SuperExpr:
		return c.compileSuperExpr(e)

	default:
		return fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (c *Compiler) compileFunctionExpr(f *ast.FunctionExpr) error {
	return c.compileFunction("", f.Params, f.Body, TYPE_FUNCTION)
}

func (c *Compiler) compileIdentifier(id *ast.Identifier) error {
	line := id.Token.Line

	// Try local variable first
	localIdx := c.resolveLocal(id.Name)
	if localIdx == -2 {
		return fmt.Errorf("cannot read local variable in its own initializer: %s", id.Name)
	}
	if localIdx != -1 {
		c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
		c.emitByte(byte(localIdx), line)
		return nil
	}

	// Try upvalue (captured variable)
	upvalueIdx := c.resolveUpvalue(id.Name)
	if upvalueIdx != -1 {
		c.emitByte(byte(bytecode.OP_GET_UPVALUE), line)
		c.emitByte(byte(upvalueIdx), line)
		return nil
	}

	// Must be a global variable
	idx, exists := c.getGlobalVariable(id.Name)
	if !exists {
		// For globals, we add them lazily (they might be defined later)
		idx = c.addGlobalVariable(id.Name)
	}

	c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
	c.emitU16(uint16(idx), line)
	return nil
}

func (c *Compiler) compileAssignment(a *ast.AssignExpr) error {
	line := a.Token.Line

	// Get the target variable
	switch target := a.Target.(type) {
	case *ast.Identifier:
		// Compile the value
		if err := c.compileExpression(a.Value); err != nil {
			return err
		}

		// Try local first
		localIdx := c.resolveLocal(target.Name)
		if localIdx != -1 && localIdx != -2 {
			c.emitByte(byte(bytecode.OP_SET_LOCAL), line)
			c.emitByte(byte(localIdx), line)
			return nil
		}

		// Try upvalue
		upvalueIdx := c.resolveUpvalue(target.Name)
		if upvalueIdx != -1 {
			c.emitByte(byte(bytecode.OP_SET_UPVALUE), line)
			c.emitByte(byte(upvalueIdx), line)
			return nil
		}

		// Must be global
		idx, exists := c.getGlobalVariable(target.Name)
		if !exists {
			idx = c.addGlobalVariable(target.Name)
		}

		c.emitByte(byte(bytecode.OP_SET_GLOBAL), line)
		c.emitU16(uint16(idx), line)
		return nil

	case *ast.IndexExpr:
		// Compile: array[index] = value
		// Stack order: array, index, value
		if err := c.compileExpression(target.Object); err != nil {
			return err
		}
		if err := c.compileExpression(target.Index); err != nil {
			return err
		}
		if err := c.compileExpression(a.Value); err != nil {
			return err
		}
		c.emitByte(byte(bytecode.OP_SET_INDEX), line)
		return nil

	case *ast.PropertyExpr:
		// Compile: obj.prop = value
		// Stack order: object, value
		if err := c.compileExpression(target.Object); err != nil {
			return err
		}
		if err := c.compileExpression(a.Value); err != nil {
			return err
		}
		nameIdx := c.chunk().AddConstant(target.Property)
		c.emitByte(byte(bytecode.OP_SET_PROPERTY), line)
		c.emitU16(uint16(nameIdx), line)
		return nil

	default:
		return fmt.Errorf("invalid assignment target: %T", a.Target)
	}
}

func (c *Compiler) compileNumber(n *ast.NumberLiteral) error {
	idx := c.chunk().AddConstant(n.Value)
	c.emitByte(byte(bytecode.OP_CONSTANT), n.Token.Line)
	c.emitU16(uint16(idx), n.Token.Line)
	return nil
}

func (c *Compiler) compileString(s *ast.StringLiteral) error {
	idx := c.chunk().AddConstant(s.Value)
	c.emitByte(byte(bytecode.OP_CONSTANT), s.Token.Line)
	c.emitU16(uint16(idx), s.Token.Line)
	return nil
}

func (c *Compiler) compileBoolean(b *ast.BoolLiteral) error {
	if b.Value {
		c.emitByte(byte(bytecode.OP_TRUE), b.Token.Line)
	} else {
		c.emitByte(byte(bytecode.OP_FALSE), b.Token.Line)
	}
	return nil
}

func (c *Compiler) compileBinary(b *ast.BinaryExpr) error {
	// For string concatenation, we use OP_CONCAT
	isStringConcat := isStringExpr(b.Left) && isStringExpr(b.Right) && b.Op == token.PLUS

	// Compile left operand
	if err := c.compileExpression(b.Left); err != nil {
		return err
	}

	// Compile right operand
	if err := c.compileExpression(b.Right); err != nil {
		return err
	}

	// Emit operator
	line := b.Token.Line
	switch b.Op {
	case token.PLUS:
		if isStringConcat {
			c.emitByte(byte(bytecode.OP_CONCAT), line)
		} else {
			c.emitByte(byte(bytecode.OP_ADD), line)
		}
	case token.MINUS:
		c.emitByte(byte(bytecode.OP_SUBTRACT), line)
	case token.STAR:
		c.emitByte(byte(bytecode.OP_MULTIPLY), line)
	case token.SLASH:
		c.emitByte(byte(bytecode.OP_DIVIDE), line)
	case token.PERCENT:
		c.emitByte(byte(bytecode.OP_MODULO), line)
	case token.EQ:
		c.emitByte(byte(bytecode.OP_EQUAL), line)
	case token.NEQ:
		c.emitByte(byte(bytecode.OP_NOT_EQUAL), line)
	case token.LT:
		c.emitByte(byte(bytecode.OP_LESS), line)
	case token.LTE:
		c.emitByte(byte(bytecode.OP_LESS_EQUAL), line)
	case token.GT:
		c.emitByte(byte(bytecode.OP_GREATER), line)
	case token.GTE:
		c.emitByte(byte(bytecode.OP_GREATER_EQUAL), line)
	default:
		return fmt.Errorf("unknown binary operator: %v", b.Op)
	}

	return nil
}

func (c *Compiler) compileUnary(u *ast.UnaryExpr) error {
	// Compile operand
	if err := c.compileExpression(u.Operand); err != nil {
		return err
	}

	// Emit operator
	line := u.Token.Line
	switch u.Op {
	case token.MINUS:
		c.emitByte(byte(bytecode.OP_NEGATE), line)
	case token.NOT:
		c.emitByte(byte(bytecode.OP_NOT), line)
	default:
		return fmt.Errorf("unknown unary operator: %v", u.Op)
	}

	return nil
}

func (c *Compiler) compileCall(call *ast.CallExpr) error {
	line := call.Token.Line

	// Check if it's a built-in function
	if ident, ok := call.Function.(*ast.Identifier); ok {
		builtinID, isBuiltin := builtinFunctions[ident.Name]
		if isBuiltin {
			return c.compileBuiltinCall(call, builtinID)
		}
	}

	// Compile the function expression (callee)
	if err := c.compileExpression(call.Function); err != nil {
		return err
	}

	// Compile the arguments
	for _, arg := range call.Arguments {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Emit the call instruction
	argCount := len(call.Arguments)
	if argCount > 255 {
		return fmt.Errorf("cannot have more than 255 arguments")
	}

	c.emitByte(byte(bytecode.OP_CALL), line)
	c.emitByte(byte(argCount), line)

	return nil
}

var builtinFunctions = map[string]int{
	"println":  bytecode.BUILTIN_PRINTLN,
	"print":    bytecode.BUILTIN_PRINT,
	"len":      bytecode.BUILTIN_LEN,
	"push":     bytecode.BUILTIN_PUSH,
	"pop":      bytecode.BUILTIN_POP,
	"typeof":   bytecode.BUILTIN_TYPEOF,
	"toString": bytecode.BUILTIN_TOSTRING,
	"toNumber": bytecode.BUILTIN_TONUMBER,
	"sqrt":     bytecode.BUILTIN_SQRT,
	"floor":    bytecode.BUILTIN_FLOOR,
	"ceil":     bytecode.BUILTIN_CEIL,
	"abs":      bytecode.BUILTIN_ABS,
}

func (c *Compiler) compileBuiltinCall(call *ast.CallExpr, builtinID int) error {
	line := call.Token.Line

	// Compile arguments
	for _, arg := range call.Arguments {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Emit OP_BUILTIN with builtin ID and arg count
	c.emitByte(byte(bytecode.OP_BUILTIN), line)
	c.emitByte(byte(builtinID), line)
	c.emitByte(byte(len(call.Arguments)), line)

	return nil
}

// Helper to check if expression is a string literal
func isStringExpr(expr ast.Expression) bool {
	_, ok := expr.(*ast.StringLiteral)
	return ok
}

// Bytecode emission helpers

func (c *Compiler) emitByte(b byte, line int) {
	c.chunk().Write(b, line)
}

func (c *Compiler) emitU16(v uint16, line int) {
	c.chunk().WriteU16(v, line)
}

func (c *Compiler) emitBytes(line int, bytes ...byte) {
	for _, b := range bytes {
		c.emitByte(b, line)
	}
}

// ----------------------------------------------------------------------------
// Array and Object Compilation
// ----------------------------------------------------------------------------

func (c *Compiler) compileArrayLiteral(arr *ast.ArrayLiteral) error {
	line := arr.Token.Line

	// Compile all elements
	for _, elem := range arr.Elements {
		if err := c.compileExpression(elem); err != nil {
			return err
		}
	}

	// Emit OP_ARRAY with element count
	count := len(arr.Elements)
	if count > 65535 {
		return fmt.Errorf("array literal too large (max 65535 elements)")
	}

	c.emitByte(byte(bytecode.OP_ARRAY), line)
	c.emitU16(uint16(count), line)

	return nil
}

func (c *Compiler) compileObjectLiteral(obj *ast.ObjectLiteral) error {
	line := obj.Token.Line

	// Compile all key-value pairs (key, then value)
	for _, prop := range obj.Properties {
		// Push key as string constant
		keyIdx := c.chunk().AddConstant(prop.Key)
		c.emitByte(byte(bytecode.OP_CONSTANT), line)
		c.emitU16(uint16(keyIdx), line)

		// Compile value
		if err := c.compileExpression(prop.Value); err != nil {
			return err
		}
	}

	// Emit OP_OBJECT with property count
	count := len(obj.Properties)
	if count > 65535 {
		return fmt.Errorf("object literal too large (max 65535 properties)")
	}

	c.emitByte(byte(bytecode.OP_OBJECT), line)
	c.emitU16(uint16(count), line)

	return nil
}

func (c *Compiler) compileIndexExpr(idx *ast.IndexExpr) error {
	line := idx.Token.Line

	// Compile the object/array
	if err := c.compileExpression(idx.Object); err != nil {
		return err
	}

	// Compile the index
	if err := c.compileExpression(idx.Index); err != nil {
		return err
	}

	c.emitByte(byte(bytecode.OP_GET_INDEX), line)
	return nil
}

func (c *Compiler) compilePropertyExpr(prop *ast.PropertyExpr) error {
	line := prop.Token.Line

	// Compile the object
	if err := c.compileExpression(prop.Object); err != nil {
		return err
	}

	// Emit OP_GET_PROPERTY with property name
	nameIdx := c.chunk().AddConstant(prop.Property)
	c.emitByte(byte(bytecode.OP_GET_PROPERTY), line)
	c.emitU16(uint16(nameIdx), line)

	return nil
}

// ----------------------------------------------------------------------------
// Class Compilation
// ----------------------------------------------------------------------------

func (c *Compiler) compileNewExpr(expr *ast.NewExpr) error {
	line := expr.Token.Line

	// Get the class
	classIdx := c.chunk().AddConstant(expr.ClassName)
	c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
	c.emitU16(uint16(classIdx), line)

	// Compile arguments
	for _, arg := range expr.Arguments {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Emit call instruction (the class acts as a constructor)
	argCount := len(expr.Arguments)
	if argCount > 255 {
		return fmt.Errorf("cannot have more than 255 constructor arguments")
	}

	c.emitByte(byte(bytecode.OP_CALL), line)
	c.emitByte(byte(argCount), line)

	return nil
}

func (c *Compiler) compileThisExpr(expr *ast.ThisExpr) error {
	line := expr.Token.Line

	if c.current.funcType != TYPE_METHOD {
		return fmt.Errorf("cannot use 'this' outside of a method")
	}

	// 'this' is always in slot 0 of the current frame
	c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
	c.emitByte(0, line)

	return nil
}

func (c *Compiler) compileSuperExpr(expr *ast.SuperExpr) error {
	line := expr.Token.Line

	if c.current.funcType != TYPE_METHOD {
		return fmt.Errorf("cannot use 'super' outside of a method")
	}

	// Get 'this' (slot 0)
	c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
	c.emitByte(0, line)

	// Compile arguments
	for _, arg := range expr.Arguments {
		if err := c.compileExpression(arg); err != nil {
			return err
		}
	}

	// Emit super invoke for constructor call
	initIdx := c.chunk().AddConstant("constructor")
	c.emitByte(byte(bytecode.OP_SUPER_INVOKE), line)
	c.emitU16(uint16(initIdx), line)
	c.emitByte(byte(len(expr.Arguments)), line)

	return nil
}

func (c *Compiler) compileClassDecl(decl *ast.ClassDecl) error {
	line := decl.Token.Line

	// Create the class
	nameIdx := c.chunk().AddConstant(decl.Name)
	c.emitByte(byte(bytecode.OP_CLASS), line)
	c.emitU16(uint16(nameIdx), line)

	// Define the class as a global
	if c.current.scopeDepth > 0 {
		if err := c.addLocal(decl.Name); err != nil {
			return err
		}
		c.markInitialized()
	} else {
		globalIdx := c.addGlobalVariable(decl.Name)
		c.emitByte(byte(bytecode.OP_SET_GLOBAL), line)
		c.emitU16(uint16(globalIdx), line)
		c.emitByte(byte(bytecode.OP_POP), line)
	}

	// Handle inheritance
	if decl.SuperClass != "" {
		// Get the superclass
		superIdx := c.chunk().AddConstant(decl.SuperClass)
		c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
		c.emitU16(uint16(superIdx), line)

		// Get the subclass
		if c.current.scopeDepth > 0 {
			localIdx := c.resolveLocal(decl.Name)
			c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
			c.emitByte(byte(localIdx), line)
		} else {
			c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
			c.emitU16(uint16(nameIdx), line)
		}

		// Emit inherit instruction
		c.emitByte(byte(bytecode.OP_INHERIT), line)
	}

	// Compile constructor
	if decl.Constructor != nil {
		// Get the class onto the stack
		if c.current.scopeDepth > 0 {
			localIdx := c.resolveLocal(decl.Name)
			c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
			c.emitByte(byte(localIdx), line)
		} else {
			c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
			c.emitU16(uint16(nameIdx), line)
		}

		// Compile constructor as a method
		if err := c.compileMethod("constructor", decl.Constructor.Params, decl.Constructor.Body); err != nil {
			return err
		}
	}

	// Compile methods
	for _, method := range decl.Methods {
		// Get the class onto the stack
		if c.current.scopeDepth > 0 {
			localIdx := c.resolveLocal(decl.Name)
			c.emitByte(byte(bytecode.OP_GET_LOCAL), line)
			c.emitByte(byte(localIdx), line)
		} else {
			c.emitByte(byte(bytecode.OP_GET_GLOBAL), line)
			c.emitU16(uint16(nameIdx), line)
		}

		if err := c.compileMethod(method.Name, method.Params, method.Body); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileMethod(name string, params []*ast.Parameter, body *ast.Block) error {
	line := body.Token.Line

	// Compile the method body
	if err := c.compileFunction(name, params, body, TYPE_METHOD); err != nil {
		return err
	}

	// Emit OP_METHOD with method name
	nameIdx := c.chunk().AddConstant(name)
	c.emitByte(byte(bytecode.OP_METHOD), line)
	c.emitU16(uint16(nameIdx), line)

	return nil
}
