package lexer

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

// Destructuring uses existing tokens (LBRACKET, RBRACKET, LBRACE, RBRACE)
// These tests verify the lexer correctly tokenizes destructuring patterns

func TestNextToken_ArrayDestructuring(t *testing.T) {
	input := `let [a, b] = arr`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.LBRACKET, "["},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RBRACKET, "]"},
		{token.ASSIGN, "="},
		{token.IDENT, "arr"},
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

func TestNextToken_ObjectDestructuring(t *testing.T) {
	input := `let {x, y} = point`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.ASSIGN, "="},
		{token.IDENT, "point"},
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

func TestNextToken_ObjectDestructuringWithRename(t *testing.T) {
	input := `let {x: newX, y: newY} = point`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.IDENT, "newX"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.COLON, ":"},
		{token.IDENT, "newY"},
		{token.RBRACE, "}"},
		{token.ASSIGN, "="},
		{token.IDENT, "point"},
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
