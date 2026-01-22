package parser

import (
	"testing"

	"github.com/pocketlang/gots/pkg/ast"
	"github.com/pocketlang/gots/pkg/lexer"
)

func TestParseNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"42;", 42},
		{"3.14;", 3.14},
		{"0;", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has %d statements, want 1", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExprStmt)
			if !ok {
				t.Fatalf("statement is not *ast.ExprStmt, got %T", program.Statements[0])
			}

			testNumberLiteral(t, stmt.Expr, tt.expected)
		})
	}
}

func TestParseStringLiteral(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	literal, ok := stmt.Expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.StringLiteral, got %T", stmt.Expr)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value = %q, want %q", literal.Value, "hello world")
	}
}

func TestParseBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExprStmt)
			testBooleanLiteral(t, stmt.Expr, tt.expected)
		})
	}
}

func TestParseNullLiteral(t *testing.T) {
	input := "null;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	_, ok := stmt.Expr.(*ast.NullLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.NullLiteral, got %T", stmt.Expr)
	}
}

func TestParseIdentifier(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	testIdentifier(t, stmt.Expr, "foobar")
}

func TestParseUnaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"-5;", "-", 5.0},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExprStmt)
			expr, ok := stmt.Expr.(*ast.UnaryExpr)
			if !ok {
				t.Fatalf("expr is not *ast.UnaryExpr, got %T", stmt.Expr)
			}

			if expr.Token.Literal != tt.operator {
				t.Errorf("operator = %q, want %q", expr.Token.Literal, tt.operator)
			}

			testLiteralExpression(t, expr.Operand, tt.value)
		})
	}
}

func TestParseBinaryExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5.0, "+", 5.0},
		{"5 - 5;", 5.0, "-", 5.0},
		{"5 * 5;", 5.0, "*", 5.0},
		{"5 / 5;", 5.0, "/", 5.0},
		{"5 % 5;", 5.0, "%", 5.0},
		{"5 > 5;", 5.0, ">", 5.0},
		{"5 < 5;", 5.0, "<", 5.0},
		{"5 >= 5;", 5.0, ">=", 5.0},
		{"5 <= 5;", 5.0, "<=", 5.0},
		{"5 == 5;", 5.0, "==", 5.0},
		{"5 != 5;", 5.0, "!=", 5.0},
		{"true && false;", true, "&&", false},
		{"true || false;", true, "||", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExprStmt)
			testBinaryExpression(t, stmt.Expr, tt.leftValue, tt.operator, tt.rightValue)
		})
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b;", "((-a) * b)"},
		{"!-a;", "(!(-a))"},
		{"a + b + c;", "((a + b) + c)"},
		{"a + b - c;", "((a + b) - c)"},
		{"a * b * c;", "((a * b) * c)"},
		{"a * b / c;", "((a * b) / c)"},
		{"a + b / c;", "(a + (b / c))"},
		{"a + b * c + d / e - f;", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 > 4 == 3 < 4;", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4;", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5;", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true && false || true;", "((true && false) || true)"},
		{"a || b && c;", "(a || (b && c))"},
		{"1 + (2 + 3) + 4;", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2;", "((5 + 5) * 2)"},
		{"2 / (5 + 5);", "(2 / (5 + 5))"},
		{"-(5 + 5);", "(-(5 + 5))"},
		{"!(true == true);", "(!(true == true))"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExprStmt)
			actual := stmt.Expr.String()
			if actual != tt.expected {
				t.Errorf("got %q, want %q", actual, tt.expected)
			}
		})
	}
}

func TestParseCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	call, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, call.Function, "add")

	if len(call.Arguments) != 3 {
		t.Fatalf("wrong number of arguments. got=%d, want=3", len(call.Arguments))
	}

	testLiteralExpression(t, call.Arguments[0], 1.0)
	testBinaryExpression(t, call.Arguments[1], 2.0, "*", 3.0)
	testBinaryExpression(t, call.Arguments[2], 4.0, "+", 5.0)
}

func TestParseCallExpressionNoArgs(t *testing.T) {
	input := "foo();"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	call, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, call.Function, "foo")

	if len(call.Arguments) != 0 {
		t.Errorf("wrong number of arguments. got=%d, want=0", len(call.Arguments))
	}
}

func TestParseArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	arr, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.ArrayLiteral, got %T", stmt.Expr)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("wrong number of elements. got=%d, want=3", len(arr.Elements))
	}

	testNumberLiteral(t, arr.Elements[0], 1)
	testBinaryExpression(t, arr.Elements[1], 2.0, "*", 2.0)
	testBinaryExpression(t, arr.Elements[2], 3.0, "+", 3.0)
}

func TestParseEmptyArrayLiteral(t *testing.T) {
	input := "[];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	arr, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.ArrayLiteral, got %T", stmt.Expr)
	}

	if len(arr.Elements) != 0 {
		t.Errorf("wrong number of elements. got=%d, want=0", len(arr.Elements))
	}
}

func TestParseIndexExpression(t *testing.T) {
	input := "myArray[1 + 1];"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	indexExpr, ok := stmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("expr is not *ast.IndexExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, indexExpr.Object, "myArray")
	testBinaryExpression(t, indexExpr.Index, 1.0, "+", 1.0)
}

func TestParseObjectLiteral(t *testing.T) {
	input := "let obj: { x: number, y: number } = { x: 1, y: 2 };"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarDecl)
	obj, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ObjectLiteral, got %T", stmt.Value)
	}

	if len(obj.Properties) != 2 {
		t.Fatalf("wrong number of properties. got=%d, want=2", len(obj.Properties))
	}

	expected := map[string]float64{"x": 1, "y": 2}
	for _, prop := range obj.Properties {
		expectedVal, ok := expected[prop.Key]
		if !ok {
			t.Errorf("unexpected key %q", prop.Key)
			continue
		}
		testNumberLiteral(t, prop.Value, expectedVal)
	}
}

func TestParseEmptyObjectLiteral(t *testing.T) {
	input := "let obj: {} = {};"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarDecl)
	obj, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ObjectLiteral, got %T", stmt.Value)
	}

	if len(obj.Properties) != 0 {
		t.Errorf("wrong number of properties. got=%d, want=0", len(obj.Properties))
	}
}

func TestParsePropertyAccess(t *testing.T) {
	input := "obj.property;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	prop, ok := stmt.Expr.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expr is not *ast.PropertyExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, prop.Object, "obj")
	if prop.Property != "property" {
		t.Errorf("prop.Property = %q, want %q", prop.Property, "property")
	}
}

func TestParseChainedPropertyAccess(t *testing.T) {
	input := "a.b.c;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	prop, ok := stmt.Expr.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expr is not *ast.PropertyExpr, got %T", stmt.Expr)
	}

	if prop.Property != "c" {
		t.Errorf("prop.Property = %q, want %q", prop.Property, "c")
	}

	inner, ok := prop.Object.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("prop.Object is not *ast.PropertyExpr, got %T", prop.Object)
	}

	if inner.Property != "b" {
		t.Errorf("inner.Property = %q, want %q", inner.Property, "b")
	}

	testIdentifier(t, inner.Object, "a")
}

func TestParseAssignExpression(t *testing.T) {
	input := "x = 5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	assign, ok := stmt.Expr.(*ast.AssignExpr)
	if !ok {
		t.Fatalf("expr is not *ast.AssignExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, assign.Target, "x")
	testNumberLiteral(t, assign.Value, 5)
}

func TestParseMethodCall(t *testing.T) {
	input := "obj.method(1, 2);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	call, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("expr is not *ast.CallExpr, got %T", stmt.Expr)
	}

	prop, ok := call.Function.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("call.Function is not *ast.PropertyExpr, got %T", call.Function)
	}

	testIdentifier(t, prop.Object, "obj")
	if prop.Property != "method" {
		t.Errorf("prop.Property = %q, want %q", prop.Property, "method")
	}

	if len(call.Arguments) != 2 {
		t.Errorf("wrong number of arguments. got=%d, want=2", len(call.Arguments))
	}
}

// Helper functions

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testNumberLiteral(t *testing.T, expr ast.Expression, expected float64) {
	num, ok := expr.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.NumberLiteral, got %T", expr)
	}

	if num.Value != expected {
		t.Errorf("num.Value = %f, want %f", num.Value, expected)
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, expected bool) {
	b, ok := expr.(*ast.BoolLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.BoolLiteral, got %T", expr)
	}

	if b.Value != expected {
		t.Errorf("b.Value = %t, want %t", b.Value, expected)
	}
}

func testIdentifier(t *testing.T, expr ast.Expression, expected string) {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr is not *ast.Identifier, got %T", expr)
	}

	if ident.Name != expected {
		t.Errorf("ident.Name = %q, want %q", ident.Name, expected)
	}
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case float64:
		testNumberLiteral(t, expr, v)
	case int:
		testNumberLiteral(t, expr, float64(v))
	case bool:
		testBooleanLiteral(t, expr, v)
	case string:
		testIdentifier(t, expr, v)
	default:
		t.Errorf("type of expected not handled. got=%T", expected)
	}
}

func testBinaryExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) {
	binary, ok := expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expr is not *ast.BinaryExpr, got %T", expr)
	}

	testLiteralExpression(t, binary.Left, left)

	if binary.Token.Literal != operator {
		t.Errorf("binary.Token.Literal = %q, want %q", binary.Token.Literal, operator)
	}

	testLiteralExpression(t, binary.Right, right)
}
