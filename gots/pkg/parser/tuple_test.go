package parser

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
)

func TestParseTupleTypeInTypeAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		numElems int
	}{
		{
			name:     "two element tuple",
			input:    `type Pair = [string, int]`,
			numElems: 2,
		},
		{
			name:     "three element tuple",
			input:    `type Triple = [string, int, boolean]`,
			numElems: 3,
		},
		{
			name:     "single element tuple",
			input:    `type Single = [string]`,
			numElems: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has %d statements, want 1", len(program.Statements))
			}

			typeAlias, ok := program.Statements[0].(*ast.TypeAliasDecl)
			if !ok {
				t.Fatalf("statement is not *ast.TypeAliasDecl, got %T", program.Statements[0])
			}

			tupleType, ok := typeAlias.AliasType.(*ast.TupleType)
			if !ok {
				t.Fatalf("aliasType is not *ast.TupleType, got %T", typeAlias.AliasType)
			}

			if len(tupleType.Elements) != tt.numElems {
				t.Errorf("tuple has %d elements, want %d", len(tupleType.Elements), tt.numElems)
			}
		})
	}
}

func TestParseTupleTypeInVarDecl(t *testing.T) {
	input := `let point: [int, int] = [1, 2]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	varDecl := program.Statements[0].(*ast.VarDecl)

	tupleType, ok := varDecl.VarType.(*ast.TupleType)
	if !ok {
		t.Fatalf("varType is not *ast.TupleType, got %T", varDecl.VarType)
	}

	if len(tupleType.Elements) != 2 {
		t.Errorf("tuple has %d elements, want 2", len(tupleType.Elements))
	}
}

func TestParseTupleTypeInFunctionParam(t *testing.T) {
	input := `function process(pair: [string, int]): void {
		println(pair)
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	funcDecl := program.Statements[0].(*ast.FuncDecl)
	if len(funcDecl.Params) != 1 {
		t.Fatalf("function has %d params, want 1", len(funcDecl.Params))
	}

	tupleType, ok := funcDecl.Params[0].ParamType.(*ast.TupleType)
	if !ok {
		t.Fatalf("param type is not *ast.TupleType, got %T", funcDecl.Params[0].ParamType)
	}

	if len(tupleType.Elements) != 2 {
		t.Errorf("tuple has %d elements, want 2", len(tupleType.Elements))
	}
}

func TestParseTupleWithRestElement(t *testing.T) {
	input := `type StringAndNumbers = [string, ...int[]]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	typeAlias := program.Statements[0].(*ast.TypeAliasDecl)
	tupleType, ok := typeAlias.AliasType.(*ast.TupleType)
	if !ok {
		t.Fatalf("aliasType is not *ast.TupleType, got %T", typeAlias.AliasType)
	}

	if len(tupleType.Elements) != 1 {
		t.Errorf("tuple has %d regular elements, want 1", len(tupleType.Elements))
	}

	if tupleType.RestElement == nil {
		t.Fatal("expected rest element but got nil")
	}

	arrayType, ok := tupleType.RestElement.(*ast.ArrayType)
	if !ok {
		t.Fatalf("rest element is not *ast.ArrayType, got %T", tupleType.RestElement)
	}

	primType, ok := arrayType.ElementType.(*ast.PrimitiveType)
	if !ok {
		t.Fatalf("array element is not *ast.PrimitiveType, got %T", arrayType.ElementType)
	}

	if primType.Kind != ast.TypeInt {
		t.Errorf("rest element type is %v, want int", primType.Kind)
	}
}

func TestParseTupleWithNestedTypes(t *testing.T) {
	input := `type Nested = [[int, int], string]`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	typeAlias := program.Statements[0].(*ast.TypeAliasDecl)
	tupleType, ok := typeAlias.AliasType.(*ast.TupleType)
	if !ok {
		t.Fatalf("aliasType is not *ast.TupleType, got %T", typeAlias.AliasType)
	}

	if len(tupleType.Elements) != 2 {
		t.Errorf("tuple has %d elements, want 2", len(tupleType.Elements))
	}

	// First element should be a nested tuple
	nestedTuple, ok := tupleType.Elements[0].(*ast.TupleType)
	if !ok {
		t.Fatalf("first element is not *ast.TupleType, got %T", tupleType.Elements[0])
	}

	if len(nestedTuple.Elements) != 2 {
		t.Errorf("nested tuple has %d elements, want 2", len(nestedTuple.Elements))
	}
}
