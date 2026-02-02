package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
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

// ============================================================
// Type Inference Tests
// ============================================================

func TestParseVarDeclWithTypeInference(t *testing.T) {
	tests := []struct {
		input    string
		name     string
		hasType  bool
		isConst  bool
	}{
		{"let x = 10;", "x", false, false},
		{"let name = \"hello\";", "name", false, false},
		{"let flag = true;", "flag", false, false},
		{"const PI = 3.14;", "PI", false, true},
		{"let arr = [1, 2, 3];", "arr", false, false},
		{"let obj = { x: 1, y: 2 };", "obj", false, false},
		// Explicit type still works
		{"let y: number = 20;", "y", true, false},
		{"const NAME: string = \"test\";", "NAME", true, true},
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

			stmt, ok := program.Statements[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("statement is not *ast.VarDecl, got %T", program.Statements[0])
			}

			if stmt.Name != tt.name {
				t.Errorf("stmt.Name = %q, want %q", stmt.Name, tt.name)
			}

			if tt.hasType && stmt.VarType == nil {
				t.Error("expected VarType to be non-nil")
			}
			if !tt.hasType && stmt.VarType != nil {
				t.Error("expected VarType to be nil for inferred type")
			}

			if stmt.IsConst != tt.isConst {
				t.Errorf("stmt.IsConst = %v, want %v", stmt.IsConst, tt.isConst)
			}
		})
	}
}

// ============================================================
// Compound Assignment Tests
// ============================================================

func TestParseCompoundAssignment(t *testing.T) {
	tests := []struct {
		input    string
		target   string
		operator string
		value    float64
	}{
		{"x += 5;", "x", "+=", 5},
		{"x -= 3;", "x", "-=", 3},
		{"x *= 2;", "x", "*=", 2},
		{"x /= 4;", "x", "/=", 4},
		{"x %= 3;", "x", "%=", 3},
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

			assign, ok := stmt.Expr.(*ast.CompoundAssignExpr)
			if !ok {
				t.Fatalf("expr is not *ast.CompoundAssignExpr, got %T", stmt.Expr)
			}

			testIdentifier(t, assign.Target, tt.target)
			if assign.Token.Literal != tt.operator {
				t.Errorf("operator = %q, want %q", assign.Token.Literal, tt.operator)
			}
			testNumberLiteral(t, assign.Value, tt.value)
		})
	}
}

// ============================================================
// Arrow Function Tests
// ============================================================

func TestParseArrowFunctionExpression(t *testing.T) {
	input := "let add = (a: number, b: number): number => a + b;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarDecl)
	fn, ok := stmt.Value.(*ast.ArrowFunctionExpr)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ArrowFunctionExpr, got %T", stmt.Value)
	}

	if len(fn.Params) != 2 {
		t.Fatalf("wrong number of params. got=%d, want=2", len(fn.Params))
	}

	if fn.Params[0].Name != "a" {
		t.Errorf("param[0].Name = %q, want %q", fn.Params[0].Name, "a")
	}

	if fn.Params[1].Name != "b" {
		t.Errorf("param[1].Name = %q, want %q", fn.Params[1].Name, "b")
	}

	// Body should be expression, not block
	if fn.Body != nil {
		t.Error("expected Body to be nil for expression arrow function")
	}

	if fn.Expression == nil {
		t.Fatal("expected Expression to be non-nil")
	}
}

func TestParseArrowFunctionWithBlock(t *testing.T) {
	input := "let double = (x: number): number => { return x * 2; };"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarDecl)
	fn, ok := stmt.Value.(*ast.ArrowFunctionExpr)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ArrowFunctionExpr, got %T", stmt.Value)
	}

	if fn.Body == nil {
		t.Fatal("expected Body to be non-nil for block arrow function")
	}

	if fn.Expression != nil {
		t.Error("expected Expression to be nil for block arrow function")
	}
}

func TestParseArrowFunctionNoParams(t *testing.T) {
	input := "let getZero = (): number => 0;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.VarDecl)
	fn, ok := stmt.Value.(*ast.ArrowFunctionExpr)
	if !ok {
		t.Fatalf("stmt.Value is not *ast.ArrowFunctionExpr, got %T", stmt.Value)
	}

	if len(fn.Params) != 0 {
		t.Errorf("wrong number of params. got=%d, want=0", len(fn.Params))
	}
}

// ============================================================
// Nullish Coalescing Tests
// ============================================================

func TestParseNullishCoalescing(t *testing.T) {
	input := "x ?? defaultValue;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	binary, ok := stmt.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expr is not *ast.BinaryExpr, got %T", stmt.Expr)
	}

	testIdentifier(t, binary.Left, "x")
	if binary.Token.Literal != "??" {
		t.Errorf("operator = %q, want %q", binary.Token.Literal, "??")
	}
	testIdentifier(t, binary.Right, "defaultValue")
}

func TestNullishCoalescingPrecedence(t *testing.T) {
	// ?? should have lower precedence than ||
	input := "a || b ?? c;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	// Should parse as (a || b) ?? c
	binary, ok := stmt.Expr.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("expr is not *ast.BinaryExpr, got %T", stmt.Expr)
	}

	if binary.Token.Literal != "??" {
		t.Errorf("outer operator = %q, want %q", binary.Token.Literal, "??")
	}

	inner, ok := binary.Left.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("left is not *ast.BinaryExpr, got %T", binary.Left)
	}

	if inner.Token.Literal != "||" {
		t.Errorf("inner operator = %q, want %q", inner.Token.Literal, "||")
	}
}

// ============================================================
// Increment/Decrement Tests
// ============================================================

func TestParseIncrementDecrement(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		prefix   bool
	}{
		{"++x;", "++", true},
		{"--x;", "--", true},
		{"x++;", "++", false},
		{"x--;", "--", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExprStmt)
			update, ok := stmt.Expr.(*ast.UpdateExpr)
			if !ok {
				t.Fatalf("expr is not *ast.UpdateExpr, got %T", stmt.Expr)
			}

			if update.Token.Literal != tt.operator {
				t.Errorf("operator = %q, want %q", update.Token.Literal, tt.operator)
			}

			if update.Prefix != tt.prefix {
				t.Errorf("prefix = %v, want %v", update.Prefix, tt.prefix)
			}

			testIdentifier(t, update.Operand, "x")
		})
	}
}

// ============================================================
// For-of Loop Tests
// ============================================================

func TestParseForOfLoop(t *testing.T) {
	input := "for (let item of items) { println(item); }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ForOfStmt)
	if !ok {
		t.Fatalf("statement is not *ast.ForOfStmt, got %T", program.Statements[0])
	}

	if stmt.Variable.Name != "item" {
		t.Errorf("variable name = %q, want %q", stmt.Variable.Name, "item")
	}

	testIdentifier(t, stmt.Iterable, "items")

	if stmt.Body == nil {
		t.Fatal("expected Body to be non-nil")
	}
}

// ============================================================
// Switch Statement Tests
// ============================================================

func TestParseSwitchStatement(t *testing.T) {
	input := `switch (x) {
		case 1:
			println("one");
			break;
		case 2:
		case 3:
			println("two or three");
			break;
		default:
			println("other");
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.SwitchStmt)
	if !ok {
		t.Fatalf("statement is not *ast.SwitchStmt, got %T", program.Statements[0])
	}

	testIdentifier(t, stmt.Discriminant, "x")

	if len(stmt.Cases) != 4 {
		t.Fatalf("wrong number of cases. got=%d, want=4", len(stmt.Cases))
	}

	// First case: 1
	if stmt.Cases[0].Test == nil {
		t.Error("case 0 test should not be nil")
	}

	// Last case should be default
	if stmt.Cases[3].Test != nil {
		t.Error("default case test should be nil")
	}
}

// ============================================================
// Optional Chaining Tests
// ============================================================

func TestParseOptionalChaining(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"obj?.property;"},
		{"arr?.[0];"},
		{"fn?.();"},
		{"a?.b?.c;"},
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
		})
	}
}

func TestParseOptionalPropertyAccess(t *testing.T) {
	input := "obj?.property;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	prop, ok := stmt.Expr.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expr is not *ast.PropertyExpr, got %T", stmt.Expr)
	}

	if !prop.Optional {
		t.Error("expected Optional to be true")
	}

	testIdentifier(t, prop.Object, "obj")
	if prop.Property != "property" {
		t.Errorf("prop.Property = %q, want %q", prop.Property, "property")
	}
}

func TestInterfaceWithFields(t *testing.T) {
	input := `interface User {
		id: int
		name: string
		active: boolean
	}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	iface, ok := program.Statements[0].(*ast.InterfaceDecl)
	if !ok {
		t.Fatalf("expected InterfaceDecl, got %T", program.Statements[0])
	}
	if len(iface.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(iface.Fields))
	}
	if iface.Fields[0].Name != "id" {
		t.Errorf("expected field name 'id', got '%s'", iface.Fields[0].Name)
	}
	if iface.Fields[1].Name != "name" {
		t.Errorf("expected field name 'name', got '%s'", iface.Fields[1].Name)
	}
	if iface.Fields[2].Name != "active" {
		t.Errorf("expected field name 'active', got '%s'", iface.Fields[2].Name)
	}
	if len(iface.Methods) != 0 {
		t.Errorf("expected 0 methods, got %d", len(iface.Methods))
	}
}

func TestInterfaceWithFieldsAndMethods(t *testing.T) {
	input := `interface Shape {
		x: int
		y: int
		area(): float
	}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	iface := program.Statements[0].(*ast.InterfaceDecl)
	if len(iface.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(iface.Fields))
	}
	if len(iface.Methods) != 1 {
		t.Fatalf("expected 1 method, got %d", len(iface.Methods))
	}
	if iface.Fields[0].Name != "x" {
		t.Errorf("expected field 'x', got '%s'", iface.Fields[0].Name)
	}
	if iface.Methods[0].Name != "area" {
		t.Errorf("expected method 'area', got '%s'", iface.Methods[0].Name)
	}
}

func TestTaggedTemplateLiteral(t *testing.T) {
	input := "html`<h1>Hello</h1>`"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected ExprStmt, got %T", program.Statements[0])
	}
	tagged, ok := stmt.Expr.(*ast.TaggedTemplateLiteral)
	if !ok {
		t.Fatalf("expected TaggedTemplateLiteral, got %T", stmt.Expr)
	}
	tagIdent, ok := tagged.Tag.(*ast.Identifier)
	if !ok {
		t.Fatalf("expected Identifier tag, got %T", tagged.Tag)
	}
	if tagIdent.Name != "html" {
		t.Errorf("expected tag 'html', got '%s'", tagIdent.Name)
	}
	if len(tagged.Parts) != 1 || tagged.Parts[0] != "<h1>Hello</h1>" {
		t.Errorf("unexpected parts: %v", tagged.Parts)
	}
	if len(tagged.Expressions) != 0 {
		t.Errorf("expected 0 expressions, got %d", len(tagged.Expressions))
	}
}

func TestTaggedTemplateLiteralWithInterpolation(t *testing.T) {
	input := "sql`SELECT * FROM users WHERE id = ${id}`"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	tagged, ok := stmt.Expr.(*ast.TaggedTemplateLiteral)
	if !ok {
		t.Fatalf("expected TaggedTemplateLiteral, got %T", stmt.Expr)
	}
	if len(tagged.Parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(tagged.Parts))
	}
	if tagged.Parts[0] != "SELECT * FROM users WHERE id = " {
		t.Errorf("unexpected first part: %q", tagged.Parts[0])
	}
	if len(tagged.Expressions) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(tagged.Expressions))
	}
}

func TestTaggedTemplateLiteralWithTypeArgs(t *testing.T) {
	input := "db.sql<User>`SELECT * FROM users`"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	tagged, ok := stmt.Expr.(*ast.TaggedTemplateLiteral)
	if !ok {
		t.Fatalf("expected TaggedTemplateLiteral, got %T", stmt.Expr)
	}
	if len(tagged.TypeArgs) != 1 {
		t.Fatalf("expected 1 type arg, got %d", len(tagged.TypeArgs))
	}
	// Tag should be db.sql (PropertyExpr)
	prop, ok := tagged.Tag.(*ast.PropertyExpr)
	if !ok {
		t.Fatalf("expected PropertyExpr tag, got %T", tagged.Tag)
	}
	if prop.Property != "sql" {
		t.Errorf("expected property 'sql', got '%s'", prop.Property)
	}
}

func TestTaggedTemplateLiteralArrayTypeArg(t *testing.T) {
	input := "db.sql<User[]>`SELECT * FROM users`"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExprStmt)
	tagged, ok := stmt.Expr.(*ast.TaggedTemplateLiteral)
	if !ok {
		t.Fatalf("expected TaggedTemplateLiteral, got %T", stmt.Expr)
	}
	if len(tagged.TypeArgs) != 1 {
		t.Fatalf("expected 1 type arg, got %d", len(tagged.TypeArgs))
	}
	_, isArray := tagged.TypeArgs[0].(*ast.ArrayType)
	if !isArray {
		t.Fatalf("expected ArrayType type arg, got %T", tagged.TypeArgs[0])
	}
}
