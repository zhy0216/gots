package ast

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

func TestNodeInterface(t *testing.T) {
	// Helper nodes for composition
	num := &NumberLiteral{Value: 42, Token: token.Token{Literal: "42"}}
	str := &StringLiteral{Value: "hello", Token: token.Token{Literal: "hello"}}
	id := &Identifier{Name: "foo", Token: token.Token{Literal: "foo"}}
	primType := &PrimitiveType{Kind: TypeNumber}

	// Test that all node types implement the Node interface
	var nodes []Node = []Node{
		// Expressions
		num,
		str,
		&BoolLiteral{Value: true, Token: token.Token{Literal: "true"}},
		&NullLiteral{},
		id,
		&BinaryExpr{Op: token.PLUS, Left: num, Right: num, Token: token.Token{Literal: "+"}},
		&UnaryExpr{Op: token.MINUS, Operand: num, Token: token.Token{Literal: "-"}},
		&CallExpr{Function: id, Arguments: []Expression{}},
		&IndexExpr{Object: id, Index: num},
		&PropertyExpr{Object: id, Property: "x"},
		&ArrayLiteral{Elements: []Expression{}},
		&ObjectLiteral{Properties: []*PropertyDef{}},
		&FunctionExpr{Params: []*Parameter{}, ReturnType: primType, Body: &Block{}},
		&NewExpr{ClassName: "Point", Arguments: []Expression{}},
		&ThisExpr{},
		&AssignExpr{Target: id, Value: num},

		// Statements
		&ExprStmt{Expr: num},
		&VarDecl{Name: "x", VarType: primType, Value: num},
		&Block{Statements: []Statement{}},
		&IfStmt{Condition: num, Consequence: &Block{}},
		&WhileStmt{Condition: num, Body: &Block{}},
		&ForStmt{Init: &VarDecl{Name: "i", VarType: primType, Value: num}, Condition: num, Update: num, Body: &Block{}},
		&ReturnStmt{},
		&BreakStmt{},
		&ContinueStmt{},

		// Declarations
		&FuncDecl{Name: "foo", Params: []*Parameter{}, ReturnType: primType, Body: &Block{}},
		&ClassDecl{Name: "Point"},
		&TypeAliasDecl{Name: "MyType", AliasType: primType},

		// Program
		&Program{Statements: []Statement{}},
	}

	for _, n := range nodes {
		// Just verify they have TokenLiteral and String methods
		_ = n.TokenLiteral()
		_ = n.String()
	}
}

func TestExpressionInterface(t *testing.T) {
	// Helper nodes
	num := &NumberLiteral{Value: 42, Token: token.Token{Literal: "42"}}
	id := &Identifier{Name: "foo", Token: token.Token{Literal: "foo"}}
	primType := &PrimitiveType{Kind: TypeNumber}

	// Test that expression types implement Expression interface
	var exprs []Expression = []Expression{
		num,
		&StringLiteral{Value: "hello", Token: token.Token{Literal: "hello"}},
		&BoolLiteral{Value: true, Token: token.Token{Literal: "true"}},
		&NullLiteral{},
		id,
		&BinaryExpr{Op: token.PLUS, Left: num, Right: num, Token: token.Token{Literal: "+"}},
		&UnaryExpr{Op: token.MINUS, Operand: num, Token: token.Token{Literal: "-"}},
		&CallExpr{Function: id, Arguments: []Expression{}},
		&IndexExpr{Object: id, Index: num},
		&PropertyExpr{Object: id, Property: "x"},
		&ArrayLiteral{Elements: []Expression{}},
		&ObjectLiteral{Properties: []*PropertyDef{}},
		&FunctionExpr{Params: []*Parameter{}, ReturnType: primType, Body: &Block{}},
		&NewExpr{ClassName: "Point", Arguments: []Expression{}},
		&ThisExpr{},
		&AssignExpr{Target: id, Value: num},
	}

	for _, e := range exprs {
		e.expressionNode() // Should not panic
	}
}

func TestStatementInterface(t *testing.T) {
	// Helper nodes
	num := &NumberLiteral{Value: 42, Token: token.Token{Literal: "42"}}
	primType := &PrimitiveType{Kind: TypeNumber}

	// Test that statement types implement Statement interface
	var stmts []Statement = []Statement{
		&ExprStmt{Expr: num},
		&VarDecl{Name: "x", VarType: primType, Value: num},
		&Block{Statements: []Statement{}},
		&IfStmt{Condition: num, Consequence: &Block{}},
		&WhileStmt{Condition: num, Body: &Block{}},
		&ForStmt{Init: &VarDecl{Name: "i", VarType: primType, Value: num}, Condition: num, Update: num, Body: &Block{}},
		&ReturnStmt{},
		&BreakStmt{},
		&ContinueStmt{},
		&FuncDecl{Name: "foo", Params: []*Parameter{}, ReturnType: primType, Body: &Block{}},
		&ClassDecl{Name: "Point"},
		&TypeAliasDecl{Name: "MyType", AliasType: primType},
	}

	for _, s := range stmts {
		s.statementNode() // Should not panic
	}
}

func TestNumberLiteral(t *testing.T) {
	n := &NumberLiteral{
		Token: token.Token{Type: token.NUMBER, Literal: "42"},
		Value: 42,
	}

	if n.TokenLiteral() != "42" {
		t.Errorf("TokenLiteral() = %q, want %q", n.TokenLiteral(), "42")
	}

	if n.String() != "42" {
		t.Errorf("String() = %q, want %q", n.String(), "42")
	}
}

func TestStringLiteral(t *testing.T) {
	s := &StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: "hello"},
		Value: "hello",
	}

	if s.TokenLiteral() != "hello" {
		t.Errorf("TokenLiteral() = %q, want %q", s.TokenLiteral(), "hello")
	}

	if s.String() != `"hello"` {
		t.Errorf("String() = %q, want %q", s.String(), `"hello"`)
	}
}

func TestBinaryExpr(t *testing.T) {
	expr := &BinaryExpr{
		Token: token.Token{Type: token.PLUS, Literal: "+"},
		Left:  &NumberLiteral{Value: 1, Token: token.Token{Literal: "1"}},
		Op:    token.PLUS,
		Right: &NumberLiteral{Value: 2, Token: token.Token{Literal: "2"}},
	}

	if expr.TokenLiteral() != "+" {
		t.Errorf("TokenLiteral() = %q, want %q", expr.TokenLiteral(), "+")
	}

	expected := "(1 + 2)"
	if expr.String() != expected {
		t.Errorf("String() = %q, want %q", expr.String(), expected)
	}
}

func TestUnaryExpr(t *testing.T) {
	expr := &UnaryExpr{
		Token:   token.Token{Type: token.MINUS, Literal: "-"},
		Op:      token.MINUS,
		Operand: &NumberLiteral{Value: 5, Token: token.Token{Literal: "5"}},
	}

	if expr.TokenLiteral() != "-" {
		t.Errorf("TokenLiteral() = %q, want %q", expr.TokenLiteral(), "-")
	}

	expected := "(-5)"
	if expr.String() != expected {
		t.Errorf("String() = %q, want %q", expr.String(), expected)
	}
}

func TestIdentifier(t *testing.T) {
	id := &Identifier{
		Token: token.Token{Type: token.IDENT, Literal: "myVar"},
		Name:  "myVar",
	}

	if id.TokenLiteral() != "myVar" {
		t.Errorf("TokenLiteral() = %q, want %q", id.TokenLiteral(), "myVar")
	}

	if id.String() != "myVar" {
		t.Errorf("String() = %q, want %q", id.String(), "myVar")
	}
}

func TestVarDecl(t *testing.T) {
	decl := &VarDecl{
		Token:   token.Token{Type: token.LET, Literal: "let"},
		Name:    "x",
		VarType: &PrimitiveType{Kind: TypeNumber},
		Value:   &NumberLiteral{Value: 42, Token: token.Token{Literal: "42"}},
		IsConst: false,
	}

	if decl.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral() = %q, want %q", decl.TokenLiteral(), "let")
	}

	expected := "let x: number = 42"
	if decl.String() != expected {
		t.Errorf("String() = %q, want %q", decl.String(), expected)
	}
}

func TestFuncDecl(t *testing.T) {
	decl := &FuncDecl{
		Token: token.Token{Type: token.FUNCTION, Literal: "function"},
		Name:  "add",
		Params: []*Parameter{
			{Name: "a", ParamType: &PrimitiveType{Kind: TypeNumber}},
			{Name: "b", ParamType: &PrimitiveType{Kind: TypeNumber}},
		},
		ReturnType: &PrimitiveType{Kind: TypeNumber},
		Body:       &Block{Statements: []Statement{}},
	}

	if decl.TokenLiteral() != "function" {
		t.Errorf("TokenLiteral() = %q, want %q", decl.TokenLiteral(), "function")
	}

	if decl.Name != "add" {
		t.Errorf("Name = %q, want %q", decl.Name, "add")
	}

	if len(decl.Params) != 2 {
		t.Errorf("len(Params) = %d, want %d", len(decl.Params), 2)
	}
}

func TestTypeNodes(t *testing.T) {
	// Primitive type
	prim := &PrimitiveType{Kind: TypeNumber}
	if prim.String() != "number" {
		t.Errorf("PrimitiveType.String() = %q, want %q", prim.String(), "number")
	}

	// Array type
	arr := &ArrayType{ElementType: &PrimitiveType{Kind: TypeString}}
	if arr.String() != "string[]" {
		t.Errorf("ArrayType.String() = %q, want %q", arr.String(), "string[]")
	}

	// Nullable type
	nullable := &NullableType{Inner: &PrimitiveType{Kind: TypeString}}
	if nullable.String() != "string | null" {
		t.Errorf("NullableType.String() = %q, want %q", nullable.String(), "string | null")
	}

	// Function type
	funcType := &FunctionType{
		ParamTypes: []Type{&PrimitiveType{Kind: TypeNumber}},
		ReturnType: &PrimitiveType{Kind: TypeString},
	}
	if funcType.String() != "(number) => string" {
		t.Errorf("FunctionType.String() = %q, want %q", funcType.String(), "(number) => string")
	}
}

func TestProgram(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&VarDecl{
				Token:   token.Token{Type: token.LET, Literal: "let"},
				Name:    "x",
				VarType: &PrimitiveType{Kind: TypeNumber},
				Value:   &NumberLiteral{Value: 5, Token: token.Token{Literal: "5"}},
			},
		},
	}

	if len(prog.Statements) != 1 {
		t.Errorf("len(Statements) = %d, want %d", len(prog.Statements), 1)
	}
}

func TestCallExpr(t *testing.T) {
	call := &CallExpr{
		Token: token.Token{Type: token.LPAREN, Literal: "("},
		Function: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "add"},
			Name:  "add",
		},
		Arguments: []Expression{
			&NumberLiteral{Value: 1, Token: token.Token{Literal: "1"}},
			&NumberLiteral{Value: 2, Token: token.Token{Literal: "2"}},
		},
	}

	expected := "add(1, 2)"
	if call.String() != expected {
		t.Errorf("String() = %q, want %q", call.String(), expected)
	}
}

func TestArrayLiteral(t *testing.T) {
	arr := &ArrayLiteral{
		Token: token.Token{Type: token.LBRACKET, Literal: "["},
		Elements: []Expression{
			&NumberLiteral{Value: 1, Token: token.Token{Literal: "1"}},
			&NumberLiteral{Value: 2, Token: token.Token{Literal: "2"}},
			&NumberLiteral{Value: 3, Token: token.Token{Literal: "3"}},
		},
	}

	expected := "[1, 2, 3]"
	if arr.String() != expected {
		t.Errorf("String() = %q, want %q", arr.String(), expected)
	}
}

func TestObjectLiteral(t *testing.T) {
	obj := &ObjectLiteral{
		Token: token.Token{Type: token.LBRACE, Literal: "{"},
		Properties: []*PropertyDef{
			{
				Key:   "x",
				Value: &NumberLiteral{Value: 10, Token: token.Token{Literal: "10"}},
			},
			{
				Key:   "y",
				Value: &NumberLiteral{Value: 20, Token: token.Token{Literal: "20"}},
			},
		},
	}

	expected := "{x: 10, y: 20}"
	if obj.String() != expected {
		t.Errorf("String() = %q, want %q", obj.String(), expected)
	}
}
