// Package vm implements error handling for the GoTS VM.
package vm

import (
	"fmt"
	"strings"
)

// RuntimeError represents an error that occurred during VM execution.
type RuntimeError struct {
	Message    string
	Line       int
	StackTrace []StackFrame
}

// StackFrame represents a single frame in the stack trace.
type StackFrame struct {
	FunctionName string
	Line         int
}

func (e *RuntimeError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("RuntimeError: %s\n", e.Message))

	if e.Line > 0 {
		sb.WriteString(fmt.Sprintf("  at line %d\n", e.Line))
	}

	if len(e.StackTrace) > 0 {
		sb.WriteString("Stack trace:\n")
		for i, frame := range e.StackTrace {
			name := frame.FunctionName
			if name == "" {
				name = "<script>"
			}
			sb.WriteString(fmt.Sprintf("  %d: %s (line %d)\n", i, name, frame.Line))
		}
	}

	return sb.String()
}

// runtimeError creates a RuntimeError with the current stack trace.
func (vm *VM) runtimeError(format string, args ...any) *RuntimeError {
	message := fmt.Sprintf(format, args...)

	// Get current line from the current frame
	line := 0
	if vm.frameCount > 0 {
		frame := vm.frame()
		ip := frame.ip - 1 // Point to the instruction that caused the error
		if ip >= 0 && ip < len(frame.closure.Function.Chunk.Lines) {
			line = frame.closure.Function.Chunk.Lines[ip]
		}
	}

	// Build stack trace
	stackTrace := make([]StackFrame, 0, vm.frameCount)
	for i := vm.frameCount - 1; i >= 0; i-- {
		frame := &vm.frames[i]
		fn := frame.closure.Function

		// Get the line number for this frame
		frameLine := 0
		ip := frame.ip - 1
		if ip >= 0 && ip < len(fn.Chunk.Lines) {
			frameLine = fn.Chunk.Lines[ip]
		}

		stackTrace = append(stackTrace, StackFrame{
			FunctionName: fn.Name,
			Line:         frameLine,
		})
	}

	return &RuntimeError{
		Message:    message,
		Line:       line,
		StackTrace: stackTrace,
	}
}

// CompileError represents an error that occurred during compilation.
type CompileError struct {
	Message string
	Line    int
	Column  int
}

func (e *CompileError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return fmt.Sprintf("CompileError at line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	if e.Line > 0 {
		return fmt.Sprintf("CompileError at line %d: %s", e.Line, e.Message)
	}
	return fmt.Sprintf("CompileError: %s", e.Message)
}

// NewCompileError creates a new CompileError.
func NewCompileError(message string, line, column int) *CompileError {
	return &CompileError{
		Message: message,
		Line:    line,
		Column:  column,
	}
}
