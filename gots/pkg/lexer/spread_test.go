package lexer

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/token"
)

func TestNextToken_Ellipsis(t *testing.T) {
	input := `...arr`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ELLIPSIS, "..."},
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

func TestNextToken_SpreadInArray(t *testing.T) {
	input := `[...arr, 1, 2]`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LBRACKET, "["},
		{token.ELLIPSIS, "..."},
		{token.IDENT, "arr"},
		{token.COMMA, ","},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RBRACKET, "]"},
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

func TestNextToken_SpreadInFunctionCall(t *testing.T) {
	input := `fn(...args)`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "fn"},
		{token.LPAREN, "("},
		{token.ELLIPSIS, "..."},
		{token.IDENT, "args"},
		{token.RPAREN, ")"},
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
