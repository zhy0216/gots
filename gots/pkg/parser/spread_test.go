package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseSpreadInArrayLiteral(t *testing.T) {
	input := `let arr: int[] = [...other, 1, 2]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	varDecl, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected *ast.VarDecl, got %T", program.Statements[0])
	}

	arrLit, ok := varDecl.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteral, got %T", varDecl.Value)
	}

	if len(arrLit.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arrLit.Elements))
	}

	spread, ok := arrLit.Elements[0].(*ast.SpreadExpr)
	if !ok {
		t.Fatalf("expected *ast.SpreadExpr, got %T", arrLit.Elements[0])
	}

	ident, ok := spread.Argument.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T", spread.Argument)
	}

	if ident.Name != "other" {
		t.Errorf("expected 'other', got %q", ident.Name)
	}
}

func TestParseSpreadInFunctionCall(t *testing.T) {
	input := `fn(...args)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected *ast.ExprStmt, got %T", program.Statements[0])
	}

	callExpr, ok := exprStmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expected *ast.CallExpr, got %T", exprStmt.Expr)
	}

	if len(callExpr.Arguments) != 1 {
		t.Fatalf("expected 1 argument, got %d", len(callExpr.Arguments))
	}

	spread, ok := callExpr.Arguments[0].(*ast.SpreadExpr)
	if !ok {
		t.Fatalf("expected *ast.SpreadExpr, got %T", callExpr.Arguments[0])
	}

	ident, ok := spread.Argument.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T", spread.Argument)
	}

	if ident.Name != "args" {
		t.Errorf("expected 'args', got %q", ident.Name)
	}
}

func TestParseMultipleSpreadInArray(t *testing.T) {
	input := `let arr: int[] = [...a, ...b, 3]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	varDecl, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected *ast.VarDecl, got %T", program.Statements[0])
	}

	arrLit, ok := varDecl.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected *ast.ArrayLiteral, got %T", varDecl.Value)
	}

	if len(arrLit.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(arrLit.Elements))
	}

	// First element is spread ...a
	spread1, ok := arrLit.Elements[0].(*ast.SpreadExpr)
	if !ok {
		t.Fatalf("expected *ast.SpreadExpr, got %T", arrLit.Elements[0])
	}
	ident1, _ := spread1.Argument.(*ast.Identifier)
	if ident1.Name != "a" {
		t.Errorf("expected 'a', got %q", ident1.Name)
	}

	// Second element is spread ...b
	spread2, ok := arrLit.Elements[1].(*ast.SpreadExpr)
	if !ok {
		t.Fatalf("expected *ast.SpreadExpr, got %T", arrLit.Elements[1])
	}
	ident2, _ := spread2.Argument.(*ast.Identifier)
	if ident2.Name != "b" {
		t.Errorf("expected 'b', got %q", ident2.Name)
	}

	// Third element is just 3
	_, ok = arrLit.Elements[2].(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("expected *ast.NumberLiteral, got %T", arrLit.Elements[2])
	}
}
