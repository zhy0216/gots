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
	chunk *bytecode.Chunk
}

// New creates a new compiler.
func New() *Compiler {
	return &Compiler{
		chunk: bytecode.NewChunk(),
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

	default:
		return fmt.Errorf("unknown statement type: %T", stmt)
	}
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
		// Will be implemented when we add variables
		return fmt.Errorf("variables not yet implemented: %s", e.Name)

	default:
		return fmt.Errorf("unknown expression type: %T", expr)
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
