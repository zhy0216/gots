package lexer

import (
	"testing"

	"github.com/pocketlang/gots/pkg/token"
)

func TestNextToken_SingleCharacters(t *testing.T) {
	input := `+-*/%(){}[];:,.`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.PERCENT, "%"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.COMMA, ","},
		{token.DOT, "."},
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

func TestNextToken_TwoCharOperators(t *testing.T) {
	input := `== != <= >= && || =>`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.ARROW, "=>"},
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

func TestNextToken_SingleVsTwoChar(t *testing.T) {
	input := `= == ! != < <= > >=`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.EQ, "=="},
		{token.NOT, "!"},
		{token.NEQ, "!="},
		{token.LT, "<"},
		{token.LTE, "<="},
		{token.GT, ">"},
		{token.GTE, ">="},
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

func TestNextToken_Numbers(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Type
		expectedLiteral string
	}{
		{"42", token.NUMBER, "42"},
		{"3.14", token.NUMBER, "3.14"},
		{"0", token.NUMBER, "0"},
		{"123456", token.NUMBER, "123456"},
		{"0.5", token.NUMBER, "0.5"},
		{"100.00", token.NUMBER, "100.00"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("type wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNextToken_Strings(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Type
		expectedLiteral string
	}{
		{`"hello"`, token.STRING, "hello"},
		{`"hello world"`, token.STRING, "hello world"},
		{`""`, token.STRING, ""},
		{`'single'`, token.STRING, "single"},
		{`'hello world'`, token.STRING, "hello world"},
		{`''`, token.STRING, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("type wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNextToken_StringEscapes(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
	}{
		{`"hello\nworld"`, "hello\nworld"},
		{`"tab\there"`, "tab\there"},
		{`"quote\""`, "quote\""},
		{`'single\'quote'`, "single'quote"},
		{`"backslash\\"`, "backslash\\"},
		{`"carriage\rreturn"`, "carriage\rreturn"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()

			if tok.Type != token.STRING {
				t.Errorf("type wrong. expected=STRING, got=%q", tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	tests := []struct {
		input           string
		expectedType    token.Type
		expectedLiteral string
	}{
		{"foo", token.IDENT, "foo"},
		{"bar", token.IDENT, "bar"},
		{"myVar", token.IDENT, "myVar"},
		{"_private", token.IDENT, "_private"},
		{"x", token.IDENT, "x"},
		{"abc123", token.IDENT, "abc123"},
		{"_123", token.IDENT, "_123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(tt.input)
			tok := l.NextToken()

			if tok.Type != tt.expectedType {
				t.Errorf("type wrong. expected=%q, got=%q", tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Errorf("literal wrong. expected=%q, got=%q", tt.expectedLiteral, tok.Literal)
			}
		})
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `let const function return if else while for break continue
		class extends new this super constructor type
		true false null
		number string boolean void`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.CONST, "const"},
		{token.FUNCTION, "function"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.WHILE, "while"},
		{token.FOR, "for"},
		{token.BREAK, "break"},
		{token.CONTINUE, "continue"},
		{token.CLASS, "class"},
		{token.EXTENDS, "extends"},
		{token.NEW, "new"},
		{token.THIS, "this"},
		{token.SUPER, "super"},
		{token.CONSTRUCTOR, "constructor"},
		{token.TYPE, "type"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NULL, "null"},
		{token.NUMBER_TYPE, "number"},
		{token.STRING_TYPE, "string"},
		{token.BOOLEAN_TYPE, "boolean"},
		{token.VOID_TYPE, "void"},
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

func TestNextToken_Comments(t *testing.T) {
	input := `
		// this is a comment
		let x = 5; // inline comment
		/* multi
		   line
		   comment */
		let y = 10;
	`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
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

func TestNextToken_LineAndColumn(t *testing.T) {
	input := `let x = 5;
let y = 10;`

	tests := []struct {
		expectedType   token.Type
		expectedLine   int
		expectedColumn int
	}{
		{token.LET, 1, 1},
		{token.IDENT, 1, 5},
		{token.ASSIGN, 1, 7},
		{token.NUMBER, 1, 9},
		{token.SEMICOLON, 1, 10},
		{token.LET, 2, 1},
		{token.IDENT, 2, 5},
		{token.ASSIGN, 2, 7},
		{token.NUMBER, 2, 9},
		{token.SEMICOLON, 2, 11},
		{token.EOF, 2, 12},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestNextToken_CompleteProgram(t *testing.T) {
	input := `function add(a: number, b: number): number {
	return a + b;
}

let result: number = add(1, 2);
println(result);`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.FUNCTION, "function"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.RPAREN, ")"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "result"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
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

func TestNextToken_Pipe(t *testing.T) {
	input := `string | null`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING_TYPE, "string"},
		{token.PIPE, "|"},
		{token.NULL, "null"},
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

func TestNextToken_ClassDeclaration(t *testing.T) {
	input := `class Point {
	x: number;
	y: number;

	constructor(x: number, y: number) {
		this.x = x;
		this.y = y;
	}
}`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.CLASS, "class"},
		{token.IDENT, "Point"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "y"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.SEMICOLON, ";"},
		{token.CONSTRUCTOR, "constructor"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.COLON, ":"},
		{token.NUMBER_TYPE, "number"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.THIS, "this"},
		{token.DOT, "."},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.THIS, "this"},
		{token.DOT, "."},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
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

func TestNextToken_IllegalCharacter(t *testing.T) {
	input := `@`

	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.ILLEGAL {
		t.Errorf("expected ILLEGAL token, got %q", tok.Type)
	}
}

func TestPeekToken(t *testing.T) {
	input := `let x = 5;`

	l := New(input)

	// Peek should not consume
	peeked := l.PeekToken()
	if peeked.Type != token.LET {
		t.Errorf("PeekToken() type = %q, want LET", peeked.Type)
	}

	// NextToken should return same token
	tok := l.NextToken()
	if tok.Type != token.LET {
		t.Errorf("NextToken() type = %q, want LET", tok.Type)
	}

	// Next peek should be next token
	peeked = l.PeekToken()
	if peeked.Type != token.IDENT {
		t.Errorf("PeekToken() type = %q, want IDENT", peeked.Type)
	}
}
