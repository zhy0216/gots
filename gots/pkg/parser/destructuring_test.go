package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseArrayDestructuring(t *testing.T) {
	input := `let [a, b]: int[] = arr`

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

	if varDecl.Pattern == nil {
		t.Fatal("expected Pattern to be set")
	}

	arrPat, ok := varDecl.Pattern.(*ast.ArrayPattern)
	if !ok {
		t.Fatalf("expected *ast.ArrayPattern, got %T", varDecl.Pattern)
	}

	if len(arrPat.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arrPat.Elements))
	}

	ident1, ok := arrPat.Elements[0].(*ast.IdentPattern)
	if !ok {
		t.Fatalf("expected *ast.IdentPattern, got %T", arrPat.Elements[0])
	}
	if ident1.Name != "a" {
		t.Errorf("expected 'a', got %q", ident1.Name)
	}

	ident2, ok := arrPat.Elements[1].(*ast.IdentPattern)
	if !ok {
		t.Fatalf("expected *ast.IdentPattern, got %T", arrPat.Elements[1])
	}
	if ident2.Name != "b" {
		t.Errorf("expected 'b', got %q", ident2.Name)
	}
}

func TestParseObjectDestructuring(t *testing.T) {
	input := `let {x, y}: {x: int, y: int} = point`

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

	if varDecl.Pattern == nil {
		t.Fatal("expected Pattern to be set")
	}

	objPat, ok := varDecl.Pattern.(*ast.ObjectPattern)
	if !ok {
		t.Fatalf("expected *ast.ObjectPattern, got %T", varDecl.Pattern)
	}

	if len(objPat.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(objPat.Properties))
	}

	if objPat.Properties[0].Key != "x" {
		t.Errorf("expected 'x', got %q", objPat.Properties[0].Key)
	}

	if objPat.Properties[1].Key != "y" {
		t.Errorf("expected 'y', got %q", objPat.Properties[1].Key)
	}
}

func TestParseObjectDestructuringWithRename(t *testing.T) {
	input := `let {x: newX, y: newY}: {x: int, y: int} = point`

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

	objPat, ok := varDecl.Pattern.(*ast.ObjectPattern)
	if !ok {
		t.Fatalf("expected *ast.ObjectPattern, got %T", varDecl.Pattern)
	}

	if len(objPat.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(objPat.Properties))
	}

	// First property: x: newX
	if objPat.Properties[0].Key != "x" {
		t.Errorf("expected key 'x', got %q", objPat.Properties[0].Key)
	}
	ident1, ok := objPat.Properties[0].Value.(*ast.IdentPattern)
	if !ok {
		t.Fatalf("expected *ast.IdentPattern, got %T", objPat.Properties[0].Value)
	}
	if ident1.Name != "newX" {
		t.Errorf("expected 'newX', got %q", ident1.Name)
	}

	// Second property: y: newY
	if objPat.Properties[1].Key != "y" {
		t.Errorf("expected key 'y', got %q", objPat.Properties[1].Key)
	}
	ident2, ok := objPat.Properties[1].Value.(*ast.IdentPattern)
	if !ok {
		t.Fatalf("expected *ast.IdentPattern, got %T", objPat.Properties[1].Value)
	}
	if ident2.Name != "newY" {
		t.Errorf("expected 'newY', got %q", ident2.Name)
	}
}

func TestParseNestedArrayDestructuring(t *testing.T) {
	input := `let [a, [b, c]]: int[][] = nested`

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

	arrPat, ok := varDecl.Pattern.(*ast.ArrayPattern)
	if !ok {
		t.Fatalf("expected *ast.ArrayPattern, got %T", varDecl.Pattern)
	}

	if len(arrPat.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arrPat.Elements))
	}

	// First element is simple ident
	ident, ok := arrPat.Elements[0].(*ast.IdentPattern)
	if !ok {
		t.Fatalf("expected *ast.IdentPattern, got %T", arrPat.Elements[0])
	}
	if ident.Name != "a" {
		t.Errorf("expected 'a', got %q", ident.Name)
	}

	// Second element is nested array pattern
	nestedArr, ok := arrPat.Elements[1].(*ast.ArrayPattern)
	if !ok {
		t.Fatalf("expected *ast.ArrayPattern, got %T", arrPat.Elements[1])
	}
	if len(nestedArr.Elements) != 2 {
		t.Fatalf("expected 2 nested elements, got %d", len(nestedArr.Elements))
	}
}
