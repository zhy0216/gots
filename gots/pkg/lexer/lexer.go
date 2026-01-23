// Package lexer implements the lexical analyzer for GoTS.
package lexer

import (
	"github.com/zhy0216/quickts/gots/pkg/token"
)

// Lexer performs lexical analysis on GoTS source code.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
	lineStart    int // position where current line started

	peeked  *token.Token
	hasPeek bool
}

// New creates a new Lexer for the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input:     input,
		line:      1,
		column:    0,
		lineStart: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column = l.position - l.lineStart + 1
}

// peekChar returns the next character without advancing.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	// Return peeked token if available
	if l.hasPeek {
		l.hasPeek = false
		return *l.peeked
	}

	l.skipWhitespaceAndComments()

	tok := token.Token{
		Line:   l.line,
		Column: l.column,
	}

	switch l.ch {
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
	case '-':
		tok = l.newToken(token.MINUS, l.ch)
	case '*':
		tok = l.newToken(token.STAR, l.ch)
	case '/':
		tok = l.newToken(token.SLASH, l.ch)
	case '%':
		tok = l.newToken(token.PERCENT, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '.':
		tok = l.newToken(token.DOT, l.ch)
	case '|':
		if l.peekChar() == '|' {
			tok = l.makeTwoCharToken(token.OR)
		} else {
			tok = l.newToken(token.PIPE, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = l.makeTwoCharToken(token.AND)
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.EQ)
		} else if l.peekChar() == '>' {
			tok = l.makeTwoCharToken(token.ARROW)
		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.NEQ)
		} else {
			tok = l.newToken(token.NOT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.LTE)
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			tok = l.makeTwoCharToken(token.GTE)
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
		tok.Line = l.line
		tok.Column = l.column
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Column = l.column
		return tok
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = l.line
			tok.Column = l.column
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// PeekToken returns the next token without consuming it.
func (l *Lexer) PeekToken() token.Token {
	if l.hasPeek {
		return *l.peeked
	}

	tok := l.NextToken()
	l.peeked = &tok
	l.hasPeek = true
	return tok
}

// newToken creates a new token from the current character.
func (l *Lexer) newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}

// makeTwoCharToken creates a token from the current and next character.
func (l *Lexer) makeTwoCharToken(tokenType token.Type) token.Token {
	ch := l.ch
	col := l.column
	l.readChar()
	return token.Token{
		Type:    tokenType,
		Literal: string(ch) + string(l.ch),
		Line:    l.line,
		Column:  col,
	}
}

// skipWhitespaceAndComments skips whitespace and comments.
func (l *Lexer) skipWhitespaceAndComments() {
	for {
		switch l.ch {
		case ' ', '\t', '\r':
			l.readChar()
		case '\n':
			l.line++
			l.readChar()
			l.lineStart = l.position
			l.column = 1
		case '/':
			if l.peekChar() == '/' {
				l.skipLineComment()
			} else if l.peekChar() == '*' {
				l.skipBlockComment()
			} else {
				return
			}
		default:
			return
		}
	}
}

// skipLineComment skips a single-line comment.
func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// skipBlockComment skips a multi-line comment.
func (l *Lexer) skipBlockComment() {
	l.readChar() // skip /
	l.readChar() // skip *

	for {
		if l.ch == 0 {
			return
		}
		if l.ch == '\n' {
			l.line++
			l.readChar()
			l.lineStart = l.position
			l.column = 1
			continue
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip *
			l.readChar() // skip /
			return
		}
		l.readChar()
	}
}

// readIdentifier reads an identifier.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number (integer or float).
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume the dot
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

// readString reads a string literal.
func (l *Lexer) readString(quote byte) string {
	l.readChar() // skip opening quote

	var result []byte
	for {
		if l.ch == quote || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '\\':
				result = append(result, '\\')
			case '"':
				result = append(result, '"')
			case '\'':
				result = append(result, '\'')
			default:
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	l.readChar() // skip closing quote
	return string(result)
}

// isLetter returns true if the character is a letter or underscore.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns true if the character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
