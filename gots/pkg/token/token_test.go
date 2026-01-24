package token

import "testing"

func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		typ  Type
		want string
	}{
		// Literals
		{NUMBER, "NUMBER"},
		{STRING, "STRING"},
		{TRUE, "TRUE"},
		{FALSE, "FALSE"},
		{NULL, "NULL"},
		{IDENT, "IDENT"},

		// Operators
		{PLUS, "+"},
		{MINUS, "-"},
		{STAR, "*"},
		{SLASH, "/"},
		{PERCENT, "%"},

		// Comparison
		{EQ, "=="},
		{NEQ, "!="},
		{LT, "<"},
		{GT, ">"},
		{LTE, "<="},
		{GTE, ">="},

		// Logical
		{AND, "&&"},
		{OR, "||"},
		{NOT, "!"},

		// Assignment
		{ASSIGN, "="},

		// Delimiters
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{SEMICOLON, ";"},
		{COLON, ":"},
		{COMMA, ","},
		{DOT, "."},
		{ARROW, "=>"},

		// Keywords
		{LET, "let"},
		{CONST, "const"},
		{FUNCTION, "function"},
		{RETURN, "return"},
		{IF, "if"},
		{ELSE, "else"},
		{WHILE, "while"},
		{FOR, "for"},
		{BREAK, "break"},
		{CONTINUE, "continue"},
		{CLASS, "class"},
		{EXTENDS, "extends"},
		{NEW, "new"},
		{THIS, "this"},
		{SUPER, "super"},
		{CONSTRUCTOR, "constructor"},
		{TYPE, "type"},

		// Type keywords
		{INT_TYPE, "int"},
		{FLOAT_TYPE, "float"},
		{STRING_TYPE, "string"},
		{BOOLEAN_TYPE, "boolean"},
		{VOID_TYPE, "void"},
		{NULL_TYPE, "null"},

		// Special
		{EOF, "EOF"},
		{ILLEGAL, "ILLEGAL"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.typ.String()
			if got != tt.want {
				t.Errorf("Type.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident string
		want  Type
	}{
		// Keywords
		{"let", LET},
		{"const", CONST},
		{"function", FUNCTION},
		{"return", RETURN},
		{"if", IF},
		{"else", ELSE},
		{"while", WHILE},
		{"for", FOR},
		{"break", BREAK},
		{"continue", CONTINUE},
		{"class", CLASS},
		{"extends", EXTENDS},
		{"new", NEW},
		{"this", THIS},
		{"super", SUPER},
		{"constructor", CONSTRUCTOR},
		{"type", TYPE},

		// Literals
		{"true", TRUE},
		{"false", FALSE},
		{"null", NULL},

		// Type keywords
		{"int", INT_TYPE},
		{"float", FLOAT_TYPE},
		{"string", STRING_TYPE},
		{"boolean", BOOLEAN_TYPE},
		{"void", VOID_TYPE},

		// Identifiers (not keywords)
		{"foo", IDENT},
		{"bar", IDENT},
		{"myVar", IDENT},
		{"_private", IDENT},
		{"x", IDENT},
		{"letx", IDENT}, // not "let"
		{"ifelse", IDENT}, // not "if" or "else"
	}

	for _, tt := range tests {
		t.Run(tt.ident, func(t *testing.T) {
			got := LookupIdent(tt.ident)
			if got != tt.want {
				t.Errorf("LookupIdent(%q) = %v, want %v", tt.ident, got, tt.want)
			}
		})
	}
}

func TestTokenCreation(t *testing.T) {
	tok := Token{
		Type:    NUMBER,
		Literal: "42",
		Line:    1,
		Column:  5,
	}

	if tok.Type != NUMBER {
		t.Errorf("tok.Type = %v, want %v", tok.Type, NUMBER)
	}
	if tok.Literal != "42" {
		t.Errorf("tok.Literal = %q, want %q", tok.Literal, "42")
	}
	if tok.Line != 1 {
		t.Errorf("tok.Line = %d, want %d", tok.Line, 1)
	}
	if tok.Column != 5 {
		t.Errorf("tok.Column = %d, want %d", tok.Column, 5)
	}
}

func TestIsKeyword(t *testing.T) {
	keywords := []Type{
		LET, CONST, FUNCTION, RETURN, IF, ELSE, WHILE, FOR,
		BREAK, CONTINUE, CLASS, EXTENDS, NEW, THIS, SUPER,
		CONSTRUCTOR, TYPE, TRUE, FALSE, NULL,
		INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOLEAN_TYPE, VOID_TYPE,
	}

	nonKeywords := []Type{
		NUMBER, STRING, IDENT, PLUS, MINUS, LPAREN, EOF, ILLEGAL,
	}

	for _, kw := range keywords {
		if !IsKeyword(kw) {
			t.Errorf("IsKeyword(%v) = false, want true", kw)
		}
	}

	for _, nk := range nonKeywords {
		if IsKeyword(nk) {
			t.Errorf("IsKeyword(%v) = true, want false", nk)
		}
	}
}

func TestIsOperator(t *testing.T) {
	operators := []Type{
		PLUS, MINUS, STAR, SLASH, PERCENT,
		EQ, NEQ, LT, GT, LTE, GTE,
		AND, OR, NOT, ASSIGN,
	}

	nonOperators := []Type{
		NUMBER, STRING, IDENT, LET, LPAREN, EOF,
	}

	for _, op := range operators {
		if !IsOperator(op) {
			t.Errorf("IsOperator(%v) = false, want true", op)
		}
	}

	for _, nop := range nonOperators {
		if IsOperator(nop) {
			t.Errorf("IsOperator(%v) = true, want false", nop)
		}
	}
}
