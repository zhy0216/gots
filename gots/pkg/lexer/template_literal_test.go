package lexer

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

func TestNextToken_TemplateLiteralSimple(t *testing.T) {
	input := "`hello world`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TEMPLATE_LITERAL, "hello world"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TemplateLiteralWithExpression(t *testing.T) {
	input := "`Hello, ${name}!`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TEMPLATE_HEAD, "Hello, "},
		{token.IDENT, "name"},
		{token.TEMPLATE_TAIL, "!"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TemplateLiteralWithMultipleExpressions(t *testing.T) {
	input := "`${a} + ${b} = ${c}`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TEMPLATE_HEAD, ""},
		{token.IDENT, "a"},
		{token.TEMPLATE_MIDDLE, " + "},
		{token.IDENT, "b"},
		{token.TEMPLATE_MIDDLE, " = "},
		{token.IDENT, "c"},
		{token.TEMPLATE_TAIL, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TemplateLiteralWithComplexExpression(t *testing.T) {
	input := "`Result: ${a + b}`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TEMPLATE_HEAD, "Result: "},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.TEMPLATE_TAIL, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TemplateLiteralWithNestedBraces(t *testing.T) {
	input := "`Value: ${obj.x}`"

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.TEMPLATE_HEAD, "Value: "},
		{token.IDENT, "obj"},
		{token.DOT, "."},
		{token.IDENT, "x"},
		{token.TEMPLATE_TAIL, ""},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_TemplateLiteralEscapes(t *testing.T) {
	input := "`Hello\\nWorld`"

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.TEMPLATE_LITERAL {
		t.Errorf("type wrong. expected=TEMPLATE_LITERAL, got=%q", tok.Type)
	}

	if tok.Literal != "Hello\nWorld" {
		t.Errorf("literal wrong. expected=%q, got=%q", "Hello\nWorld", tok.Literal)
	}
}

func TestNextToken_TemplateLiteralEmpty(t *testing.T) {
	input := "``"

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.TEMPLATE_LITERAL {
		t.Errorf("type wrong. expected=TEMPLATE_LITERAL, got=%q", tok.Type)
	}

	if tok.Literal != "" {
		t.Errorf("literal wrong. expected=%q, got=%q", "", tok.Literal)
	}
}
