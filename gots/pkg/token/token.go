// Package token defines the token types and utilities for the GoTS lexer.
package token

// Type represents the type of a token.
type Type int

// Token types enumeration
const (
	// Special tokens
	ILLEGAL Type = iota
	EOF

	// Literals
	IDENT  // identifier
	NUMBER // numeric literal
	STRING // string literal

	// Literal keywords
	TRUE  // true
	FALSE // false
	NULL  // null

	// Operators
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %

	// Compound assignment operators
	PLUS_ASSIGN    // +=
	MINUS_ASSIGN   // -=
	STAR_ASSIGN    // *=
	SLASH_ASSIGN   // /=
	PERCENT_ASSIGN // %=

	// Increment/decrement operators
	INCREMENT // ++
	DECREMENT // --

	// Comparison operators
	EQ  // ==
	NEQ // !=
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=

	// Logical operators
	AND // &&
	OR  // ||
	NOT // !

	// Nullish coalescing
	NULLISH_COALESCE // ??

	// Optional chaining
	QUESTION_DOT // ?.

	// Assignment
	ASSIGN // =

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	SEMICOLON // ;
	COLON     // :
	COMMA     // ,
	DOT       // .
	ARROW     // =>
	PIPE      // |
	QUESTION  // ?

	// Keywords
	LET         // let
	CONST       // const
	FUNCTION    // function
	RETURN      // return
	IF          // if
	ELSE        // else
	WHILE       // while
	FOR         // for
	BREAK       // break
	CONTINUE    // continue
	CLASS       // class
	EXTENDS     // extends
	NEW         // new
	THIS        // this
	SUPER       // super
	CONSTRUCTOR // constructor
	TYPE        // type
	SWITCH      // switch
	CASE        // case
	DEFAULT     // default
	OF          // of
	TRY         // try
	CATCH       // catch
	THROW       // throw

	// Type keywords
	INT_TYPE     // int
	FLOAT_TYPE   // float
	STRING_TYPE  // string
	BOOLEAN_TYPE // boolean
	VOID_TYPE    // void
	NULL_TYPE    // null (as type)

	// Advanced type keywords
	MAP       // Map
	SET       // Set
	INTERFACE // interface
	IMPORT    // import
	FROM      // from
	EXPORT    // export
)

// Token represents a lexical token with its metadata.
type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

// String returns the string representation of a token type.
func (t Type) String() string {
	if s, ok := typeStrings[t]; ok {
		return s
	}
	return "UNKNOWN"
}

// typeStrings maps token types to their string representations.
var typeStrings = map[Type]string{
	// Special
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	// Literals
	IDENT:  "IDENT",
	NUMBER: "NUMBER",
	STRING: "STRING",
	TRUE:   "TRUE",
	FALSE:  "FALSE",
	NULL:   "NULL",

	// Operators
	PLUS:    "+",
	MINUS:   "-",
	STAR:    "*",
	SLASH:   "/",
	PERCENT: "%",

	// Compound assignment
	PLUS_ASSIGN:    "+=",
	MINUS_ASSIGN:   "-=",
	STAR_ASSIGN:    "*=",
	SLASH_ASSIGN:   "/=",
	PERCENT_ASSIGN: "%=",

	// Increment/decrement
	INCREMENT: "++",
	DECREMENT: "--",

	// Comparison
	EQ:  "==",
	NEQ: "!=",
	LT:  "<",
	GT:  ">",
	LTE: "<=",
	GTE: ">=",

	// Logical
	AND: "&&",
	OR:  "||",
	NOT: "!",

	// Nullish coalescing
	NULLISH_COALESCE: "??",

	// Optional chaining
	QUESTION_DOT: "?.",

	// Assignment
	ASSIGN: "=",

	// Delimiters
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	SEMICOLON: ";",
	COLON:     ":",
	COMMA:     ",",
	DOT:       ".",
	ARROW:     "=>",
	PIPE:      "|",
	QUESTION:  "?",

	// Keywords
	LET:         "let",
	CONST:       "const",
	FUNCTION:    "function",
	RETURN:      "return",
	IF:          "if",
	ELSE:        "else",
	WHILE:       "while",
	FOR:         "for",
	BREAK:       "break",
	CONTINUE:    "continue",
	CLASS:       "class",
	EXTENDS:     "extends",
	NEW:         "new",
	THIS:        "this",
	SUPER:       "super",
	CONSTRUCTOR: "constructor",
	TYPE:        "type",
	SWITCH:      "switch",
	CASE:        "case",
	DEFAULT:     "default",
	OF:          "of",
	TRY:         "try",
	CATCH:       "catch",
	THROW:       "throw",

	// Type keywords
	INT_TYPE:     "int",
	FLOAT_TYPE:   "float",
	STRING_TYPE:  "string",
	BOOLEAN_TYPE: "boolean",
	VOID_TYPE:    "void",
	NULL_TYPE:    "null",

	// Advanced type keywords
	MAP:       "Map",
	SET:       "Set",
	INTERFACE: "interface",
	IMPORT:    "import",
	FROM:      "from",
	EXPORT:    "export",
}

// keywords maps keyword strings to their token types.
var keywords = map[string]Type{
	"let":         LET,
	"const":       CONST,
	"function":    FUNCTION,
	"return":      RETURN,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"for":         FOR,
	"break":       BREAK,
	"continue":    CONTINUE,
	"class":       CLASS,
	"extends":     EXTENDS,
	"new":         NEW,
	"this":        THIS,
	"super":       SUPER,
	"constructor": CONSTRUCTOR,
	"type":        TYPE,
	"switch":      SWITCH,
	"case":        CASE,
	"default":     DEFAULT,
	"of":          OF,
	"try":         TRY,
	"catch":       CATCH,
	"throw":       THROW,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"int":         INT_TYPE,
	"float":       FLOAT_TYPE,
	"string":      STRING_TYPE,
	"boolean":     BOOLEAN_TYPE,
	"void":        VOID_TYPE,
	"Map":         MAP,
	"Set":         SET,
	"interface":   INTERFACE,
	"import":      IMPORT,
	"from":        FROM,
	"export":      EXPORT,
}

// LookupIdent checks if an identifier is a keyword and returns the appropriate token type.
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// IsKeyword returns true if the token type is a keyword.
func IsKeyword(t Type) bool {
	switch t {
	case LET, CONST, FUNCTION, RETURN, IF, ELSE, WHILE, FOR,
		BREAK, CONTINUE, CLASS, EXTENDS, NEW, THIS, SUPER,
		CONSTRUCTOR, TYPE, TRUE, FALSE, NULL, SWITCH, CASE, DEFAULT, OF,
		TRY, CATCH, THROW,
		INT_TYPE, FLOAT_TYPE, STRING_TYPE, BOOLEAN_TYPE, VOID_TYPE,
		MAP, SET, INTERFACE, IMPORT, FROM, EXPORT:
		return true
	}
	return false
}

// IsOperator returns true if the token type is an operator.
func IsOperator(t Type) bool {
	switch t {
	case PLUS, MINUS, STAR, SLASH, PERCENT,
		EQ, NEQ, LT, GT, LTE, GTE,
		AND, OR, NOT, ASSIGN:
		return true
	}
	return false
}

// IsLiteral returns true if the token type is a literal.
func IsLiteral(t Type) bool {
	switch t {
	case NUMBER, STRING, TRUE, FALSE, NULL:
		return true
	}
	return false
}
