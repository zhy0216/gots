package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseTemplateLiteralSimple(t *testing.T) {
	input := "let x: string = `hello world`"

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

	tmpl, ok := varDecl.Value.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("expected *ast.TemplateLiteral, got %T", varDecl.Value)
	}

	if len(tmpl.Parts) != 1 {
		t.Errorf("expected 1 part, got %d", len(tmpl.Parts))
	}

	if tmpl.Parts[0] != "hello world" {
		t.Errorf("expected 'hello world', got %q", tmpl.Parts[0])
	}

	if len(tmpl.Expressions) != 0 {
		t.Errorf("expected 0 expressions, got %d", len(tmpl.Expressions))
	}
}

func TestParseTemplateLiteralWithExpression(t *testing.T) {
	input := "let x: string = `Hello, ${name}!`"

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

	tmpl, ok := varDecl.Value.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("expected *ast.TemplateLiteral, got %T", varDecl.Value)
	}

	if len(tmpl.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(tmpl.Parts))
	}

	if tmpl.Parts[0] != "Hello, " {
		t.Errorf("expected 'Hello, ', got %q", tmpl.Parts[0])
	}

	if tmpl.Parts[1] != "!" {
		t.Errorf("expected '!', got %q", tmpl.Parts[1])
	}

	if len(tmpl.Expressions) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(tmpl.Expressions))
	}

	ident, ok := tmpl.Expressions[0].(*ast.Identifier)
	if !ok {
		t.Fatalf("expected *ast.Identifier, got %T", tmpl.Expressions[0])
	}

	if ident.Name != "name" {
		t.Errorf("expected 'name', got %q", ident.Name)
	}
}

func TestParseTemplateLiteralWithMultipleExpressions(t *testing.T) {
	input := "let x: string = `${a} + ${b} = ${c}`"

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

	tmpl, ok := varDecl.Value.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("expected *ast.TemplateLiteral, got %T", varDecl.Value)
	}

	if len(tmpl.Parts) != 4 {
		t.Fatalf("expected 4 parts, got %d", len(tmpl.Parts))
	}

	expectedParts := []string{"", " + ", " = ", ""}
	for i, exp := range expectedParts {
		if tmpl.Parts[i] != exp {
			t.Errorf("parts[%d]: expected %q, got %q", i, exp, tmpl.Parts[i])
		}
	}

	if len(tmpl.Expressions) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(tmpl.Expressions))
	}

	expectedExprs := []string{"a", "b", "c"}
	for i, exp := range expectedExprs {
		ident, ok := tmpl.Expressions[i].(*ast.Identifier)
		if !ok {
			t.Fatalf("expressions[%d]: expected *ast.Identifier, got %T", i, tmpl.Expressions[i])
		}
		if ident.Name != exp {
			t.Errorf("expressions[%d]: expected %q, got %q", i, exp, ident.Name)
		}
	}
}

func TestParseTemplateLiteralWithComplexExpression(t *testing.T) {
	input := "let x: string = `Result: ${a + b}`"

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

	tmpl, ok := varDecl.Value.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("expected *ast.TemplateLiteral, got %T", varDecl.Value)
	}

	if len(tmpl.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(tmpl.Parts))
	}

	if len(tmpl.Expressions) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(tmpl.Expressions))
	}

	binExpr, ok := tmpl.Expressions[0].(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expected *ast.BinaryExpr, got %T", tmpl.Expressions[0])
	}

	if binExpr.Token.Literal != "+" {
		t.Errorf("expected '+', got %q", binExpr.Token.Literal)
	}
}

func TestParseTemplateLiteralEmpty(t *testing.T) {
	input := "let x: string = ``"

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

	tmpl, ok := varDecl.Value.(*ast.TemplateLiteral)
	if !ok {
		t.Fatalf("expected *ast.TemplateLiteral, got %T", varDecl.Value)
	}

	if len(tmpl.Parts) != 1 {
		t.Fatalf("expected 1 part, got %d", len(tmpl.Parts))
	}

	if tmpl.Parts[0] != "" {
		t.Errorf("expected empty string, got %q", tmpl.Parts[0])
	}

	if len(tmpl.Expressions) != 0 {
		t.Errorf("expected 0 expressions, got %d", len(tmpl.Expressions))
	}
}
