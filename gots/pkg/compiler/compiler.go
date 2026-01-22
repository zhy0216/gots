// Package compiler compiles AST to bytecode.
package compiler

import (
	"fmt"

	"github.com/pocketlang/gots/pkg/ast"
	"github.com/pocketlang/gots/pkg/bytecode"
	"github.com/pocketlang/gots/pkg/token"
)

// Built-in function IDs
const (
	BUILTIN_PRINTLN = iota
	BUILTIN_PRINT
	BUILTIN_LEN
	BUILTIN_PUSH
	BUILTIN_POP
	BUILTIN_TYPEOF
)

// Compiler compiles AST to bytecode.
type Compiler struct {
	chunk   *bytecode.Chunk
	globals map[string]int // Maps global variable names to constant pool indices
}

// New creates a new compiler.
func New() *Compiler {
	return &Compiler{
		chunk:   bytecode.NewChunk(),
		globals: make(map[string]int),
	}
}

// Compile compiles a program to bytecode.
func (c *Compiler) Compile(program *ast.Program) (*bytecode.Chunk, error) {
	for _, stmt := range program.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return nil, err
		}
	}

	// Emit return at end
	c.emitByte(byte(bytecode.OP_RETURN), 0)

	return c.chunk, nil
}

func (c *Compiler) compileStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		if err := c.compileExpression(s.Expr); err != nil {
			return err
		}
		// Pop the result of expression statement (unless it's a builtin that doesn't push)
		if !c.lastWasBuiltin() {
			c.emitByte(byte(bytecode.OP_POP), s.Token.Line)
		}
		return nil

	case *ast.VarDecl:
		return c.compileVarDecl(s)

	case *ast.IfStmt:
		return c.compileIfStmt(s)

	case *ast.WhileStmt:
		return c.compileWhileStmt(s)

	case *ast.Block:
		return c.compileBlock(s)

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
}

func (c *Compiler) compileVarDecl(v *ast.VarDecl) error {
	// Compile the initializer value
	if v.Value != nil {
		if err := c.compileExpression(v.Value); err != nil {
			return err
		}
	} else {
		// No initializer, default to null
		c.emitByte(byte(bytecode.OP_NULL), v.Token.Line)
	}

	// Add variable name to constant pool and track it
	nameIdx := c.addGlobalVariable(v.Name)

	// Emit OP_SET_GLOBAL to define the variable
	c.emitByte(byte(bytecode.OP_SET_GLOBAL), v.Token.Line)
	c.emitU16(uint16(nameIdx), v.Token.Line)

	// Pop the value (variable declarations are statements, not expressions)
	c.emitByte(byte(bytecode.OP_POP), v.Token.Line)

	return nil
}

func (c *Compiler) addGlobalVariable(name string) int {
	// Check if variable already exists
	if idx, exists := c.globals[name]; exists {
		return idx
	}

	// Add name to constant pool
	idx := c.chunk.AddConstant(name)
	c.globals[name] = idx
	return idx
}

func (c *Compiler) getGlobalVariable(name string) (int, bool) {
	idx, exists := c.globals[name]
	return idx, exists
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
	for _, stmt := range b.Statements {
		if err := c.compileStatement(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileWhileStmt(w *ast.WhileStmt) error {
	line := w.Token.Line

	// Remember start of loop (for jumping back)
	loopStart := c.chunk.Count()

	// Compile the condition
	if err := c.compileExpression(w.Condition); err != nil {
		return err
	}

	// Emit OP_JUMP_IF_FALSE to exit loop
	exitJump := c.emitJump(bytecode.OP_JUMP_IF_FALSE, line)

	// Pop the condition value (true case - entering loop body)
	c.emitByte(byte(bytecode.OP_POP), line)

	// Compile the body
	if err := c.compileBlock(w.Body); err != nil {
		return err
	}

	// Emit jump back to loop start
	c.emitLoop(loopStart, line)

	// Patch the exit jump to here
	c.patchJump(exitJump)

	// Pop the condition value (false case - exiting loop)
	c.emitByte(byte(bytecode.OP_POP), line)

	return nil
}

// emitLoop emits a backward jump to loopStart
func (c *Compiler) emitLoop(loopStart int, line int) {
	c.emitByte(byte(bytecode.OP_JUMP_BACK), line)

	// Calculate offset (from current position to loop start)
	offset := c.chunk.Count() - loopStart + 2 // +2 for the offset bytes we're about to emit

	if offset > 65535 {
		panic("loop body too large")
	}

	c.emitByte(byte(offset>>8), line)
	c.emitByte(byte(offset), line)
}

// emitJump emits a jump instruction with a placeholder offset and returns the position to patch
func (c *Compiler) emitJump(op bytecode.OpCode, line int) int {
	c.emitByte(byte(op), line)
	c.emitByte(0xff, line) // Placeholder high byte
	c.emitByte(0xff, line) // Placeholder low byte
	return c.chunk.Count() - 2 // Return position of the offset bytes
}

// patchJump patches a previously emitted jump to jump to the current position
func (c *Compiler) patchJump(offset int) {
	// Calculate the jump distance (from after the jump instruction to current position)
	jump := c.chunk.Count() - offset - 2 // -2 because offset points to the high byte

	if jump > 65535 {
		panic("jump too large")
	}

	c.chunk.Code[offset] = byte(jump >> 8)
	c.chunk.Code[offset+1] = byte(jump)
}

// lastWasBuiltin checks if the last emitted opcode was OP_BUILTIN
func (c *Compiler) lastWasBuiltin() bool {
	if len(c.chunk.Code) < 3 {
		return false
	}
	// OP_BUILTIN is followed by 2 bytes (builtin_id, arg_count)
	// So we check 3 positions back
	return bytecode.OpCode(c.chunk.Code[len(c.chunk.Code)-3]) == bytecode.OP_BUILTIN
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

	default:
		return fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (c *Compiler) compileIdentifier(id *ast.Identifier) error {
	// Look up the variable
	idx, exists := c.getGlobalVariable(id.Name)
	if !exists {
		return fmt.Errorf("undefined variable: %s", id.Name)
	}

	c.emitByte(byte(bytecode.OP_GET_GLOBAL), id.Token.Line)
	c.emitU16(uint16(idx), id.Token.Line)
	return nil
}

func (c *Compiler) compileAssignment(a *ast.AssignExpr) error {
	// Compile the value
	if err := c.compileExpression(a.Value); err != nil {
		return err
	}

	// Get the target variable
	switch target := a.Target.(type) {
	case *ast.Identifier:
		idx, exists := c.getGlobalVariable(target.Name)
		if !exists {
			return fmt.Errorf("undefined variable: %s", target.Name)
		}

		c.emitByte(byte(bytecode.OP_SET_GLOBAL), a.Token.Line)
		c.emitU16(uint16(idx), a.Token.Line)
		return nil

	default:
		return fmt.Errorf("invalid assignment target: %T", a.Target)
	}
}

func (c *Compiler) compileNumber(n *ast.NumberLiteral) error {
	idx := c.chunk.AddConstant(n.Value)
	c.emitByte(byte(bytecode.OP_CONSTANT), n.Token.Line)
	c.emitU16(uint16(idx), n.Token.Line)
	return nil
}

func (c *Compiler) compileString(s *ast.StringLiteral) error {
	idx := c.chunk.AddConstant(s.Value)
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
	// Check if it's a built-in function
	if ident, ok := call.Function.(*ast.Identifier); ok {
		builtinID, isBuiltin := builtinFunctions[ident.Name]
		if isBuiltin {
			return c.compileBuiltinCall(call, builtinID)
		}
	}

	// Regular function call - will be implemented later
	return fmt.Errorf("user function calls not yet implemented")
}

var builtinFunctions = map[string]int{
	"println": BUILTIN_PRINTLN,
	"print":   BUILTIN_PRINT,
	"len":     BUILTIN_LEN,
	"push":    BUILTIN_PUSH,
	"pop":     BUILTIN_POP,
	"typeof":  BUILTIN_TYPEOF,
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
	c.chunk.Write(b, line)
}

func (c *Compiler) emitU16(v uint16, line int) {
	c.chunk.WriteU16(v, line)
}

func (c *Compiler) emitBytes(line int, bytes ...byte) {
	for _, b := range bytes {
		c.emitByte(b, line)
	}
}
