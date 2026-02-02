// Package parser implements the parser for goTS.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/token"
)

// Precedence levels for Pratt parsing
const (
	_ int = iota
	LOWEST
	ASSIGN          // =
	NULLISH         // ??
	OR              // ||
	AND             // &&
	EQUALS          // == !=
	LESSGREATER     // > < >= <=
	SUM             // + -
	PRODUCT         // * / %
	PREFIX          // -x !x ++x --x
	POSTFIX         // x++ x--
	CALL            // function() array[index] obj.property
)

// precedences maps token types to their precedence levels
var precedences = map[token.Type]int{
	token.ASSIGN:           ASSIGN,
	token.PLUS_ASSIGN:      ASSIGN,
	token.MINUS_ASSIGN:     ASSIGN,
	token.STAR_ASSIGN:      ASSIGN,
	token.SLASH_ASSIGN:     ASSIGN,
	token.PERCENT_ASSIGN:   ASSIGN,
	token.NULLISH_COALESCE: NULLISH,
	token.OR:               OR,
	token.AND:              AND,
	token.EQ:               EQUALS,
	token.NEQ:              EQUALS,
	token.LT:               LESSGREATER,
	token.GT:               LESSGREATER,
	token.LTE:              LESSGREATER,
	token.GTE:              LESSGREATER,
	token.PLUS:             SUM,
	token.MINUS:            SUM,
	token.STAR:             PRODUCT,
	token.SLASH:            PRODUCT,
	token.PERCENT:          PRODUCT,
	token.INCREMENT:        POSTFIX,
	token.DECREMENT:        POSTFIX,
	token.LPAREN:           CALL,
	token.LBRACKET:         CALL,
	token.DOT:              CALL,
	token.QUESTION_DOT:     CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses goTS source code into an AST.
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

// New creates a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[token.Type]prefixParseFn),
		infixParseFns:  make(map[token.Type]infixParseFn),
	}
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.NULL, p.parseNullLiteral)
	p.registerPrefix(token.MINUS, p.parseUnaryExpression)
	p.registerPrefix(token.NOT, p.parseUnaryExpression)
	p.registerPrefix(token.INCREMENT, p.parsePrefixUpdateExpression)
	p.registerPrefix(token.DECREMENT, p.parsePrefixUpdateExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedOrArrowFunction)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseObjectLiteral)
	p.registerPrefix(token.THIS, p.parseThisExpression)
	p.registerPrefix(token.NEW, p.parseNewExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.SUPER, p.parseSuperExpression)
	p.registerPrefix(token.TEMPLATE_LITERAL, p.parseTemplateLiteral)
	p.registerPrefix(token.TEMPLATE_HEAD, p.parseTemplateLiteral)
	p.registerPrefix(token.ELLIPSIS, p.parseSpreadExpression)
	p.registerPrefix(token.AWAIT, p.parseAwaitExpression)
	p.registerPrefix(token.ASYNC, p.parseAsyncExpression)
	p.registerPrefix(token.SLASH, p.parseRegexLiteral)

	p.registerInfix(token.PLUS, p.parseBinaryExpression)
	p.registerInfix(token.MINUS, p.parseBinaryExpression)
	p.registerInfix(token.STAR, p.parseBinaryExpression)
	p.registerInfix(token.SLASH, p.parseBinaryExpression)
	p.registerInfix(token.PERCENT, p.parseBinaryExpression)
	p.registerInfix(token.EQ, p.parseBinaryExpression)
	p.registerInfix(token.NEQ, p.parseBinaryExpression)
	p.registerInfix(token.LT, p.parseBinaryExpression)
	p.registerInfix(token.GT, p.parseBinaryExpression)
	p.registerInfix(token.LTE, p.parseBinaryExpression)
	p.registerInfix(token.GTE, p.parseBinaryExpression)
	p.registerInfix(token.AND, p.parseBinaryExpression)
	p.registerInfix(token.OR, p.parseBinaryExpression)
	p.registerInfix(token.NULLISH_COALESCE, p.parseBinaryExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parsePropertyExpression)
	p.registerInfix(token.QUESTION_DOT, p.parseOptionalChainExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.PLUS_ASSIGN, p.parseCompoundAssignExpression)
	p.registerInfix(token.MINUS_ASSIGN, p.parseCompoundAssignExpression)
	p.registerInfix(token.STAR_ASSIGN, p.parseCompoundAssignExpression)
	p.registerInfix(token.SLASH_ASSIGN, p.parseCompoundAssignExpression)
	p.registerInfix(token.PERCENT_ASSIGN, p.parseCompoundAssignExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixUpdateExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixUpdateExpression)

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("line %d: expected next token to be %s, got %s instead",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("line %d: no prefix parse function for %s found",
		p.curToken.Line, t)
	p.errors = append(p.errors, msg)
}

// Errors returns the list of parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

// ParseProgram parses the entire program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Check for decorators before function/class declarations
	if p.curTokenIs(token.AT) {
		return p.parseDecoratedDeclaration()
	}

	switch p.curToken.Type {
	case token.LET:
		return p.parseVarDeclaration(false)
	case token.CONST:
		return p.parseVarDeclaration(true)
	case token.FUNCTION:
		return p.parseFunctionDeclaration()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.BREAK:
		return p.parseBreakStatement()
	case token.CONTINUE:
		return p.parseContinueStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	case token.CLASS:
		return p.parseClassDeclaration()
	case token.TYPE:
		return p.parseTypeAlias()
	case token.INTERFACE:
		return p.parseInterfaceDeclaration()
	case token.IMPORT:
		return p.parseImportDeclaration()
	case token.EXPORT:
		return p.parseExportDeclaration()
	case token.SWITCH:
		return p.parseSwitchStatement()
	case token.TRY:
		return p.parseTryStatement()
	case token.THROW:
		return p.parseThrowStatement()
	case token.ENUM:
		return p.parseEnumDeclaration()
	case token.ASYNC:
		return p.parseAsyncFunctionDeclaration()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.curToken}
	stmt.Expr = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && (precedence < p.peekPrecedence() || p.peekTokenIs(token.TEMPLATE_LITERAL) || p.peekTokenIs(token.TEMPLATE_HEAD)) {
		// Try to parse explicit type arguments for generic function calls: f<T>(...)
		if p.peekTokenIs(token.LT) {
			if _, isIdent := leftExp.(*ast.Identifier); isIdent {
				p.nextToken() // move to LT
				if result := p.tryParseExplicitTypeArgs(leftExp); result != nil {
					leftExp = result
					continue
				}
				// Failed - LT is already curToken, fall through to normal infix
				leftExp = p.parseBinaryExpression(leftExp)
				continue
			}
			if _, isProp := leftExp.(*ast.PropertyExpr); isProp {
				p.nextToken() // move to LT
				if result := p.tryParseExplicitTypeArgs(leftExp); result != nil {
					leftExp = result
					continue
				}
				leftExp = p.parseBinaryExpression(leftExp)
				continue
			}
		}

		// Check for tagged template literal: expr`...`
		if p.peekTokenIs(token.TEMPLATE_LITERAL) || p.peekTokenIs(token.TEMPLATE_HEAD) {
			p.nextToken()
			tmpl := p.parseTemplateLiteral().(*ast.TemplateLiteral)
			leftExp = &ast.TaggedTemplateLiteral{
				Token:       tmpl.Token,
				Tag:         leftExp,
				Parts:       tmpl.Parts,
				Expressions: tmpl.Expressions,
			}
			continue
		}

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// Prefix parse functions

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("line %d: could not parse %q as number",
			p.curToken.Line, p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseTemplateLiteral() ast.Expression {
	lit := &ast.TemplateLiteral{Token: p.curToken}

	// Simple template literal with no interpolations
	if p.curToken.Type == token.TEMPLATE_LITERAL {
		lit.Parts = []string{p.curToken.Literal}
		lit.Expressions = []ast.Expression{}
		return lit
	}

	// Template literal with interpolations
	// Starts with TEMPLATE_HEAD
	lit.Parts = []string{p.curToken.Literal}
	lit.Expressions = []ast.Expression{}

	for {
		// Move past the TEMPLATE_HEAD or TEMPLATE_MIDDLE
		p.nextToken()

		// Parse the interpolated expression
		expr := p.parseExpression(LOWEST)
		if expr == nil {
			return nil
		}
		lit.Expressions = append(lit.Expressions, expr)

		// After the expression, we expect TEMPLATE_MIDDLE or TEMPLATE_TAIL
		// The lexer handles the } and returns the appropriate token
		p.nextToken()

		if p.curToken.Type == token.TEMPLATE_TAIL {
			lit.Parts = append(lit.Parts, p.curToken.Literal)
			break
		} else if p.curToken.Type == token.TEMPLATE_MIDDLE {
			lit.Parts = append(lit.Parts, p.curToken.Literal)
			// Continue parsing more expressions
		} else {
			p.errors = append(p.errors, fmt.Sprintf("expected TEMPLATE_MIDDLE or TEMPLATE_TAIL, got %s", p.curToken.Type))
			return nil
		}
	}

	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BoolLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseRegexLiteral() ast.Expression {
	// The current token is SLASH. We need to seek the lexer back to just after
	// the slash and re-read the input as a regex literal.
	slashPos := p.curToken.Position

	// Seek lexer to just after the slash
	p.l.SeekTo(slashPos + 1)

	// Read the regex literal
	regexTok := p.l.ReadRegexLiteral()

	if regexTok.Type == token.ILLEGAL {
		p.errors = append(p.errors, fmt.Sprintf("line %d: unterminated regex literal", regexTok.Line))
		return nil
	}

	// Parse "pattern\x00flags" format
	parts := strings.SplitN(regexTok.Literal, "\x00", 2)
	pattern := parts[0]
	flags := ""
	if len(parts) > 1 {
		flags = parts[1]
	}

	// Update parser's peek token since we re-read from the lexer
	p.peekToken = p.l.NextToken()

	return &ast.RegexLiteral{
		Token:   regexTok,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expr := &ast.UnaryExpr{
		Token: p.curToken,
		Op:    p.curToken.Type,
	}

	p.nextToken()
	expr.Operand = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseSpreadExpression() ast.Expression {
	expr := &ast.SpreadExpr{Token: p.curToken}

	p.nextToken()
	expr.Argument = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{Token: p.curToken}
	arr.Elements = p.parseExpressionList(token.RBRACKET)
	return arr
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{Token: p.curToken}
	obj.Properties = []*ast.PropertyDef{}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return obj
	}

	p.nextToken()

	// Parse first property
	prop := p.parsePropertyDef()
	if prop != nil {
		obj.Properties = append(obj.Properties, prop)
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip comma
		p.nextToken() // move to key
		prop := p.parsePropertyDef()
		if prop != nil {
			obj.Properties = append(obj.Properties, prop)
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return obj
}

func (p *Parser) parsePropertyDef() *ast.PropertyDef {
	prop := &ast.PropertyDef{}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	prop.Key = p.curToken.Literal

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	prop.Value = p.parseExpression(LOWEST)

	return prop
}

func (p *Parser) parseThisExpression() ast.Expression {
	return &ast.ThisExpr{Token: p.curToken}
}

func (p *Parser) parseSuperExpression() ast.Expression {
	expr := &ast.SuperExpr{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr.Arguments = p.parseExpressionList(token.RPAREN)

	return expr
}

func (p *Parser) parseNewExpression() ast.Expression {
	expr := &ast.NewExpr{Token: p.curToken}

	// Accept IDENT, MAP, or SET as class name
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		expr.ClassName = p.curToken.Literal
	} else if p.peekTokenIs(token.MAP) {
		p.nextToken()
		expr.ClassName = "Map"
	} else if p.peekTokenIs(token.SET) {
		p.nextToken()
		expr.ClassName = "Set"
	} else {
		msg := fmt.Sprintf("line %d: expected class name after 'new', got %s", p.peekToken.Line, p.peekToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Check for type arguments: new Map<K, V>() or new Set<T>()
	if p.peekTokenIs(token.LT) {
		expr.TypeArgs = p.parseTypeArguments()
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	expr.Arguments = p.parseExpressionList(token.RPAREN)

	return expr
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	fn := &ast.FunctionExpr{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Params = p.parseParameterList()

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	fn.ReturnType = p.parseType()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

// Infix parse functions

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpr{
		Token: p.curToken,
		Left:  left,
		Op:    p.curToken.Type,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

// parserState holds the parser state for speculative parsing.
type parserState struct {
	curToken  token.Token
	peekToken token.Token
	errors    []string
	lexState  lexer.LexerState
}

func (p *Parser) saveState() parserState {
	errorsCopy := make([]string, len(p.errors))
	copy(errorsCopy, p.errors)
	return parserState{
		curToken:  p.curToken,
		peekToken: p.peekToken,
		errors:    errorsCopy,
		lexState:  p.l.SaveState(),
	}
}

func (p *Parser) restoreState(state parserState) {
	p.curToken = state.curToken
	p.peekToken = state.peekToken
	p.errors = state.errors
	p.l.RestoreState(state.lexState)
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpr{Token: p.curToken, Function: function}
	call.Arguments = p.parseExpressionList(token.RPAREN)
	return call
}

// tryParseExplicitTypeArgs attempts to parse explicit type arguments on a function call.
// Called when we see LT after an identifier-like expression in the expression loop.
// Returns the call expression with type args, or nil if this is not a generic call.
func (p *Parser) tryParseExplicitTypeArgs(left ast.Expression) ast.Expression {
	// Save state for backtracking
	state := p.saveState()

	// curToken is LT, try to parse type arguments
	// We're positioned at LT (consumed by expression loop moving to infix)
	typeArgs := []ast.Type{}

	// Parse first type argument
	p.nextToken()
	arg := p.parseType()
	if arg == nil {
		p.restoreState(state)
		return nil
	}
	typeArgs = append(typeArgs, arg)

	// Parse remaining type arguments
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next type arg
		arg := p.parseType()
		if arg == nil {
			p.restoreState(state)
			return nil
		}
		typeArgs = append(typeArgs, arg)
	}

	// Expect >
	if !p.peekTokenIs(token.GT) {
		p.restoreState(state)
		return nil
	}
	p.nextToken() // consume >

	// Must be followed by ( for a generic call, or template literal for tagged template
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // consume (
		call := &ast.CallExpr{
			Token:    p.curToken,
			Function: left,
			TypeArgs: typeArgs,
		}
		call.Arguments = p.parseExpressionList(token.RPAREN)
		return call
	}

	if p.peekTokenIs(token.TEMPLATE_LITERAL) || p.peekTokenIs(token.TEMPLATE_HEAD) {
		p.nextToken()
		tmpl := p.parseTemplateLiteral().(*ast.TemplateLiteral)
		return &ast.TaggedTemplateLiteral{
			Token:       tmpl.Token,
			Tag:         left,
			TypeArgs:    typeArgs,
			Parts:       tmpl.Parts,
			Expressions: tmpl.Expressions,
		}
	}

	p.restoreState(state)
	return nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpr{Token: p.curToken, Object: left}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parsePropertyExpression(left ast.Expression) ast.Expression {
	expr := &ast.PropertyExpr{Token: p.curToken, Object: left}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	expr.Property = p.curToken.Literal
	return expr
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expr := &ast.AssignExpr{
		Token:  p.curToken,
		Target: left,
	}

	p.nextToken()
	expr.Value = p.parseExpression(ASSIGN - 1) // Right associative

	return expr
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// Statement parsers

func (p *Parser) parseVarDeclaration(isConst bool) *ast.VarDecl {
	decl := &ast.VarDecl{Token: p.curToken, IsConst: isConst}

	// Check for destructuring pattern
	if p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		decl.Pattern = p.parseArrayPattern()
	} else if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		decl.Pattern = p.parseObjectPattern()
	} else {
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		decl.Name = p.curToken.Literal
	}

	// Type annotation is optional for type inference
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		p.nextToken()
		decl.VarType = p.parseType()
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	decl.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return decl
}

func (p *Parser) parseArrayPattern() *ast.ArrayPattern {
	pattern := &ast.ArrayPattern{Token: p.curToken}
	pattern.Elements = []ast.Pattern{}

	// Empty array pattern
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return pattern
	}

	p.nextToken()
	pattern.Elements = append(pattern.Elements, p.parsePatternElement())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ','
		p.nextToken()
		pattern.Elements = append(pattern.Elements, p.parsePatternElement())
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return pattern
}

func (p *Parser) parseObjectPattern() *ast.ObjectPattern {
	pattern := &ast.ObjectPattern{Token: p.curToken}
	pattern.Properties = []*ast.PropertyPattern{}

	// Empty object pattern
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return pattern
	}

	p.nextToken()
	pattern.Properties = append(pattern.Properties, p.parsePropertyPattern())

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume ','
		p.nextToken()
		pattern.Properties = append(pattern.Properties, p.parsePropertyPattern())
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return pattern
}

func (p *Parser) parsePropertyPattern() *ast.PropertyPattern {
	prop := &ast.PropertyPattern{}

	if !p.curTokenIs(token.IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("expected identifier in object pattern, got %s", p.curToken.Type))
		return nil
	}

	prop.Key = p.curToken.Literal

	// Check for rename pattern: {x: newX}
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		p.nextToken()
		prop.Value = p.parsePatternElement()
	} else {
		// Shorthand: {x} is equivalent to {x: x}
		prop.Value = &ast.IdentPattern{Token: p.curToken, Name: p.curToken.Literal}
	}

	return prop
}

func (p *Parser) parsePatternElement() ast.Pattern {
	switch p.curToken.Type {
	case token.LBRACKET:
		return p.parseArrayPattern()
	case token.LBRACE:
		return p.parseObjectPattern()
	case token.IDENT:
		return &ast.IdentPattern{Token: p.curToken, Name: p.curToken.Literal}
	default:
		p.errors = append(p.errors, fmt.Sprintf("unexpected token in pattern: %s", p.curToken.Type))
		return nil
	}
}

func (p *Parser) parseFunctionDeclaration() *ast.FuncDecl {
	decl := &ast.FuncDecl{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	decl.Name = p.curToken.Literal

	// Check for type parameters <T, U>
	if p.peekTokenIs(token.LT) {
		decl.TypeParams = p.parseTypeParameters()
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	decl.Params = p.parseParameterList()

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	decl.ReturnType = p.parseType()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	decl.Body = p.parseBlockStatement()

	return decl
}

// parseAsyncFunctionDeclaration parses: async function name(...): Promise<T> { ... }
func (p *Parser) parseAsyncFunctionDeclaration() *ast.FuncDecl {
	asyncToken := p.curToken // Save 'async' token

	if !p.expectPeek(token.FUNCTION) {
		return nil
	}

	decl := p.parseFunctionDeclaration()
	if decl != nil {
		decl.IsAsync = true
		decl.Token = asyncToken // Use async token for error reporting
	}
	return decl
}

// parseAsyncExpression handles async arrow functions as expressions: async () => ...
func (p *Parser) parseAsyncExpression() ast.Expression {
	asyncToken := p.curToken

	// Could be: async function() {} or async () => {}
	if p.peekTokenIs(token.FUNCTION) {
		p.nextToken()
		fn := p.parseFunctionExpression()
		if fnExpr, ok := fn.(*ast.FunctionExpr); ok {
			fnExpr.IsAsync = true
			fnExpr.Token = asyncToken
			return fnExpr
		}
		return fn
	}

	// Async arrow function: async () => {} or async (x) => {}
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
		arrow := p.parseGroupedOrArrowFunction()
		if arrowFn, ok := arrow.(*ast.ArrowFunctionExpr); ok {
			arrowFn.IsAsync = true
			arrowFn.Token = asyncToken
			return arrowFn
		}
		return arrow
	}

	// Single param async arrow: async x => ...
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		paramName := p.curToken.Literal

		if !p.expectPeek(token.ARROW) {
			return nil
		}

		arrowToken := p.curToken
		p.nextToken()

		arrow := &ast.ArrowFunctionExpr{
			Token:   arrowToken,
			Params:  []*ast.Parameter{{Name: paramName}},
			IsAsync: true,
		}

		if p.curTokenIs(token.LBRACE) {
			arrow.Body = p.parseBlockStatement()
		} else {
			arrow.Expression = p.parseExpression(LOWEST)
		}

		return arrow
	}

	p.errors = append(p.errors, fmt.Sprintf("line %d: unexpected token after async: %s", p.peekToken.Line, p.peekToken.Type))
	return nil
}

// parseAwaitExpression parses: await <expression>
func (p *Parser) parseAwaitExpression() ast.Expression {
	expr := &ast.AwaitExpr{Token: p.curToken}

	p.nextToken()
	expr.Argument = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseParameterList() []*ast.Parameter {
	params := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()
	param := p.parseParameter()
	if param != nil {
		params = append(params, param)
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := p.parseParameter()
		if param != nil {
			params = append(params, param)
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseParameter() *ast.Parameter {
	param := &ast.Parameter{}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	param.Name = p.curToken.Literal

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	param.ParamType = p.parseType()

	return param
}

func (p *Parser) parseReturnStatement() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) {
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStmt {
	stmt := &ast.IfStmt{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			stmt.Alternative = p.parseIfStatement()
		} else if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStmt {
	stmt := &ast.WhileStmt{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	forToken := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Check if this is a for-of loop
	// Pattern: for (let/const ident of ...)
	if p.peekTokenIs(token.LET) || p.peekTokenIs(token.CONST) {
		// Look ahead to detect for-of
		// Save current position conceptually (we'll parse and decide)
		p.nextToken() // move to let/const
		isConst := p.curTokenIs(token.CONST)

		if !p.expectPeek(token.IDENT) {
			return nil
		}

		varName := p.curToken.Literal
		varToken := p.curToken

		// Check for optional type annotation
		var varType ast.Type
		if p.peekTokenIs(token.COLON) {
			p.nextToken()
			p.nextToken()
			varType = p.parseType()
		}

		// Check if next is 'of' (for-of) or '=' (regular for)
		if p.peekTokenIs(token.OF) {
			// This is a for-of loop
			stmt := &ast.ForOfStmt{Token: forToken}
			stmt.Variable = &ast.VarDecl{
				Token:   varToken,
				Name:    varName,
				VarType: varType,
				IsConst: isConst,
			}

			p.nextToken() // consume 'of'
			p.nextToken() // move to iterable expression
			stmt.Iterable = p.parseExpression(LOWEST)

			if !p.expectPeek(token.RPAREN) {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			stmt.Body = p.parseBlockStatement()
			return stmt
		}

		// Regular for loop - continue parsing
		stmt := &ast.ForStmt{Token: forToken}
		stmt.Init = &ast.VarDecl{
			Token:   varToken,
			Name:    varName,
			VarType: varType,
			IsConst: isConst,
		}

		if !p.expectPeek(token.ASSIGN) {
			return nil
		}

		p.nextToken()
		stmt.Init.Value = p.parseExpression(LOWEST)

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)

		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		p.nextToken()
		stmt.Update = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		stmt.Body = p.parseBlockStatement()

		return stmt
	}

	// Fallback for other for loop patterns
	stmt := &ast.ForStmt{Token: forToken}

	p.nextToken()
	stmt.Init = p.parseVarDeclaration(false)

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	p.nextToken()
	stmt.Update = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseBreakStatement() *ast.BreakStmt {
	stmt := &ast.BreakStmt{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseContinueStatement() *ast.ContinueStmt {
	stmt := &ast.ContinueStmt{Token: p.curToken}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.Block {
	block := &ast.Block{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseClassDeclaration() *ast.ClassDecl {
	decl := &ast.ClassDecl{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	decl.Name = p.curToken.Literal

	// Check for type parameters <T, U>
	if p.peekTokenIs(token.LT) {
		decl.TypeParams = p.parseTypeParameters()
	}

	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		decl.SuperClass = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.parseClassBody(decl)

	return decl
}

func (p *Parser) parseClassBody(decl *ast.ClassDecl) {
	decl.Fields = []*ast.Field{}
	decl.Methods = []*ast.Method{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.CONSTRUCTOR) {
			decl.Constructor = p.parseConstructor()
		} else if p.curTokenIs(token.IDENT) {
			// Could be field or method
			if p.peekTokenIs(token.COLON) {
				// Field
				field := p.parseField()
				if field != nil {
					decl.Fields = append(decl.Fields, field)
				}
			} else if p.peekTokenIs(token.LPAREN) {
				// Method
				method := p.parseMethod()
				if method != nil {
					decl.Methods = append(decl.Methods, method)
				}
			}
		}
		p.nextToken()
	}
}

func (p *Parser) parseField() *ast.Field {
	field := &ast.Field{Name: p.curToken.Literal}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	field.FieldType = p.parseType()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return field
}

func (p *Parser) parseConstructor() *ast.Constructor {
	cons := &ast.Constructor{}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	cons.Params = p.parseParameterList()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	cons.Body = p.parseBlockStatement()

	return cons
}

func (p *Parser) parseMethod() *ast.Method {
	method := &ast.Method{Name: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	method.Params = p.parseParameterList()

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	method.ReturnType = p.parseType()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	method.Body = p.parseBlockStatement()

	return method
}

func (p *Parser) parseTypeAlias() *ast.TypeAliasDecl {
	decl := &ast.TypeAliasDecl{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	decl.Name = p.curToken.Literal

	// Check for type parameters <T, U>
	if p.peekTokenIs(token.LT) {
		decl.TypeParams = p.parseTypeParameters()
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	decl.AliasType = p.parseType()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return decl
}

func (p *Parser) parseEnumDeclaration() *ast.EnumDecl {
	decl := &ast.EnumDecl{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	decl.Name = p.curToken.Literal

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	decl.Members = p.parseEnumMembers()

	return decl
}

func (p *Parser) parseEnumMembers() []*ast.EnumMember {
	members := []*ast.EnumMember{}

	// Handle empty enum
	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return members
	}

	p.nextToken()

	for {
		member := &ast.EnumMember{}

		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("expected enum member name, got %s", p.curToken.Type))
			return members
		}

		member.Name = p.curToken.Literal

		// Check for optional value assignment
		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken() // consume =
			p.nextToken() // move to value
			member.Value = p.parseExpression(LOWEST)
		}

		members = append(members, member)

		// Check for comma or end
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next member
		} else if p.peekTokenIs(token.RBRACE) {
			p.nextToken()
			break
		} else {
			p.errors = append(p.errors, fmt.Sprintf("expected ',' or '}' in enum, got %s", p.peekToken.Type))
			return members
		}
	}

	return members
}

// Type parsing

func (p *Parser) parseType() ast.Type {
	// Parse base type (possibly with array suffix)
	typ := p.parseSingleType()
	if typ == nil {
		return nil
	}

	// Handle intersection types (higher precedence) - A & B & C
	if p.peekTokenIs(token.AMPERSAND) {
		types := []ast.Type{typ}
		for p.peekTokenIs(token.AMPERSAND) {
			p.nextToken() // consume &
			p.nextToken() // move to next type
			nextType := p.parseSingleType()
			if nextType == nil {
				return nil
			}
			types = append(types, nextType)
		}
		typ = &ast.IntersectionType{Types: types}
	}

	// Handle union types (lower precedence) - A | B | C or (A & B) | C
	if p.peekTokenIs(token.PIPE) {
		types := []ast.Type{typ}
		for p.peekTokenIs(token.PIPE) {
			p.nextToken() // consume |
			p.nextToken() // move to next type

			// Parse next type which could be an intersection
			nextType := p.parseSingleType()
			if nextType == nil {
				return nil
			}

			// Check for intersection at this level too
			if p.peekTokenIs(token.AMPERSAND) {
				intersectionTypes := []ast.Type{nextType}
				for p.peekTokenIs(token.AMPERSAND) {
					p.nextToken() // consume &
					p.nextToken() // move to next type
					interType := p.parseSingleType()
					if interType == nil {
						return nil
					}
					intersectionTypes = append(intersectionTypes, interType)
				}
				nextType = &ast.IntersectionType{Types: intersectionTypes}
			}

			types = append(types, nextType)
		}

		// Special case: if it's just T | null, use NullableType for backward compatibility
		if len(types) == 2 {
			if prim, ok := types[1].(*ast.PrimitiveType); ok && prim.Kind == ast.TypeNull {
				return &ast.NullableType{Inner: types[0]}
			}
		}

		return &ast.UnionType{Types: types}
	}

	return typ
}

// parseSingleType parses a single type without union or intersection handling
func (p *Parser) parseSingleType() ast.Type {
	var typ ast.Type

	switch p.curToken.Type {
	case token.INT_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeInt}
	case token.FLOAT_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeFloat}
	case token.NUMBER_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeNumber}
	case token.STRING_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeString}
	case token.BOOLEAN_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeBoolean}
	case token.VOID_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeVoid}
	case token.NULL:
		typ = &ast.PrimitiveType{Kind: ast.TypeNull}
	case token.MAP:
		typ = p.parseMapType()
	case token.SET:
		typ = p.parseSetType()
	case token.IDENT:
		name := p.curToken.Literal
		if name == "RegExp" {
			typ = &ast.RegExpType{}
		} else {
			namedType := &ast.NamedType{Name: name}
			if p.peekTokenIs(token.LT) {
				namedType.TypeArgs = p.parseTypeArguments()
			}
			typ = namedType
		}
	case token.LBRACE:
		typ = p.parseObjectType()
	case token.LPAREN:
		typ = p.parseFunctionType()
	case token.LBRACKET:
		typ = p.parseTupleType()
	// Handle literal types
	case token.STRING:
		typ = &ast.LiteralType{
			Kind:  ast.TypeString,
			Value: p.curToken.Literal,
		}
	case token.NUMBER:
		// Determine if it's int or float
		kind := ast.TypeInt
		if strings.Contains(p.curToken.Literal, ".") {
			kind = ast.TypeFloat
		}
		typ = &ast.LiteralType{
			Kind:  kind,
			Value: p.curToken.Literal,
		}
	case token.TRUE, token.FALSE:
		typ = &ast.LiteralType{
			Kind:  ast.TypeBoolean,
			Value: p.curToken.Literal,
		}
	default:
		msg := fmt.Sprintf("line %d: unexpected token %s in type", p.curToken.Line, p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	// Check for array type
	for p.peekTokenIs(token.LBRACKET) {
		p.nextToken()
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
		typ = &ast.ArrayType{ElementType: typ}
	}

	return typ
}

func (p *Parser) parseObjectType() *ast.ObjectType {
	objType := &ast.ObjectType{Properties: []*ast.ObjectTypeProperty{}}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return objType
	}

	p.nextToken()

	prop := p.parseObjectTypeProperty()
	if prop != nil {
		objType.Properties = append(objType.Properties, prop)
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		prop := p.parseObjectTypeProperty()
		if prop != nil {
			objType.Properties = append(objType.Properties, prop)
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return objType
}

func (p *Parser) parseObjectTypeProperty() *ast.ObjectTypeProperty {
	prop := &ast.ObjectTypeProperty{}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	prop.Name = p.curToken.Literal

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	prop.PropType = p.parseType()

	return prop
}

func (p *Parser) parseTupleType() *ast.TupleType {
	tupleType := &ast.TupleType{
		Token:    p.curToken, // The '['
		Elements: []ast.Type{},
	}

	// Empty tuple: []
	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return tupleType
	}

	p.nextToken() // move to first element

	// Check for rest element at the start: [...int[]]
	if p.curTokenIs(token.ELLIPSIS) {
		p.nextToken() // move to type after ...
		restType := p.parseType()
		if restType == nil {
			return nil
		}
		// Extract the element type from the array type
		if arrayType, ok := restType.(*ast.ArrayType); ok {
			tupleType.RestElement = arrayType
		} else {
			p.errors = append(p.errors, "rest element must be an array type")
			return nil
		}
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
		return tupleType
	}

	// Parse first element
	firstElem := p.parseType()
	if firstElem == nil {
		return nil
	}
	tupleType.Elements = append(tupleType.Elements, firstElem)

	// Parse remaining elements
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next element

		// Check for rest element: [string, int, ...number[]]
		if p.curTokenIs(token.ELLIPSIS) {
			p.nextToken() // move to type after ...
			restType := p.parseType()
			if restType == nil {
				return nil
			}
			// Extract the element type from the array type
			if arrayType, ok := restType.(*ast.ArrayType); ok {
				tupleType.RestElement = arrayType
			} else {
				p.errors = append(p.errors, "rest element must be an array type")
				return nil
			}
			break // Rest element must be last
		}

		// Parse normal element
		elem := p.parseType()
		if elem == nil {
			return nil
		}
		tupleType.Elements = append(tupleType.Elements, elem)
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return tupleType
}

func (p *Parser) parseFunctionType() *ast.FunctionType {
	funcType := &ast.FunctionType{ParamTypes: []ast.Type{}}

	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		funcType.ParamTypes = append(funcType.ParamTypes, p.parseType())

		for p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
			funcType.ParamTypes = append(funcType.ParamTypes, p.parseType())
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.ARROW) {
		return nil
	}

	p.nextToken()
	funcType.ReturnType = p.parseType()

	return funcType
}

// ============================================================
// New Feature Parsers
// ============================================================

// parseGroupedOrArrowFunction handles both grouped expressions and arrow functions.
// It needs to look ahead to determine which one it is.
func (p *Parser) parseGroupedOrArrowFunction() ast.Expression {
	startToken := p.curToken

	// Check for empty params: ()
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() // consume )

		// Check if followed by : (arrow function)
		if p.peekTokenIs(token.COLON) {
			// Empty params arrow function
			return p.parseArrowFunctionWithParams(startToken, []*ast.Parameter{})
		}

		// Empty parentheses - this is actually a syntax error in most cases
		// but we'll return nil for now
		return nil
	}

	// Check if next token is IDENT followed by COLON (arrow function param pattern)
	if p.peekTokenIs(token.IDENT) {
		// Save current position by noting peek token
		p.nextToken() // move to ident

		if p.peekTokenIs(token.COLON) {
			// This looks like arrow function params
			// Parse first parameter
			params := []*ast.Parameter{}
			param := p.parseParameter()
			if param != nil {
				params = append(params, param)
			}

			// Parse remaining parameters
			for p.peekTokenIs(token.COMMA) {
				p.nextToken() // consume comma
				p.nextToken() // move to next param
				param := p.parseParameter()
				if param != nil {
					params = append(params, param)
				}
			}

			if !p.expectPeek(token.RPAREN) {
				return nil
			}

			// Check if followed by : type =>
			if p.peekTokenIs(token.COLON) {
				return p.parseArrowFunctionWithParams(startToken, params)
			}

			// Not an arrow function after all - this is an error
			// because we've consumed ident: type which isn't a valid expression
			return nil
		}

		// Not an arrow function, parse as grouped expression
		// curToken is now the identifier, we need to parse from here
		exp := p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		return exp
	}

	// Not an identifier, so definitely not an arrow function
	// Parse as grouped expression
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseArrowFunctionWithParams parses an arrow function after parameters have been parsed.
func (p *Parser) parseArrowFunctionWithParams(startToken token.Token, params []*ast.Parameter) ast.Expression {
	arrow := &ast.ArrowFunctionExpr{
		Token:  startToken,
		Params: params,
	}

	// Parse return type
	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	arrow.ReturnType = p.parseType()

	// Expect =>
	if !p.expectPeek(token.ARROW) {
		return nil
	}

	p.nextToken()

	// Check if body is a block or expression
	if p.curTokenIs(token.LBRACE) {
		arrow.Body = p.parseBlockStatement()
	} else {
		arrow.Expression = p.parseExpression(LOWEST)
	}

	return arrow
}

// parseCompoundAssignExpression parses compound assignment (+=, -=, etc.)
func (p *Parser) parseCompoundAssignExpression(left ast.Expression) ast.Expression {
	expr := &ast.CompoundAssignExpr{
		Token:  p.curToken,
		Target: left,
		Op:     p.curToken.Type,
	}

	p.nextToken()
	expr.Value = p.parseExpression(ASSIGN - 1) // Right associative

	return expr
}

// parsePrefixUpdateExpression parses prefix increment/decrement (++x, --x)
func (p *Parser) parsePrefixUpdateExpression() ast.Expression {
	expr := &ast.UpdateExpr{
		Token:  p.curToken,
		Op:     p.curToken.Type,
		Prefix: true,
	}

	p.nextToken()
	expr.Operand = p.parseExpression(PREFIX)

	return expr
}

// parsePostfixUpdateExpression parses postfix increment/decrement (x++, x--)
func (p *Parser) parsePostfixUpdateExpression(left ast.Expression) ast.Expression {
	return &ast.UpdateExpr{
		Token:   p.curToken,
		Op:      p.curToken.Type,
		Operand: left,
		Prefix:  false,
	}
}

// parseOptionalChainExpression parses optional chaining (?.)
func (p *Parser) parseOptionalChainExpression(left ast.Expression) ast.Expression {
	tok := p.curToken

	// Check what follows ?.
	if p.peekTokenIs(token.LPAREN) {
		// Optional call: fn?.()
		p.nextToken()
		call := &ast.CallExpr{Token: p.curToken, Function: left, Optional: true}
		call.Arguments = p.parseExpressionList(token.RPAREN)
		return call
	} else if p.peekTokenIs(token.LBRACKET) {
		// Optional index: arr?.[0]
		p.nextToken()
		expr := &ast.IndexExpr{Token: p.curToken, Object: left, Optional: true}
		p.nextToken()
		expr.Index = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
		return expr
	} else {
		// Optional property: obj?.prop
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		return &ast.PropertyExpr{
			Token:    tok,
			Object:   left,
			Property: p.curToken.Literal,
			Optional: true,
		}
	}
}

// parseForOfStatement parses for-of loops
func (p *Parser) parseForOfStatement() *ast.ForOfStmt {
	stmt := &ast.ForOfStmt{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Expect 'let' or 'const'
	if !p.expectPeek(token.LET) && !p.curTokenIs(token.CONST) {
		return nil
	}

	isConst := p.curTokenIs(token.CONST)

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	varName := p.curToken.Literal

	// Optional type annotation
	var varType ast.Type
	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		varType = p.parseType()
	}

	stmt.Variable = &ast.VarDecl{
		Token:   p.curToken,
		Name:    varName,
		VarType: varType,
		IsConst: isConst,
	}

	if !p.expectPeek(token.OF) {
		return nil
	}

	p.nextToken()
	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseSwitchStatement parses switch statements
func (p *Parser) parseSwitchStatement() *ast.SwitchStmt {
	stmt := &ast.SwitchStmt{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	stmt.Discriminant = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Cases = []*ast.CaseClause{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		clause := p.parseCaseClause()
		if clause != nil {
			stmt.Cases = append(stmt.Cases, clause)
		}
		p.nextToken()
	}

	return stmt
}

// parseCaseClause parses a case or default clause
func (p *Parser) parseCaseClause() *ast.CaseClause {
	clause := &ast.CaseClause{Token: p.curToken}

	if p.curTokenIs(token.CASE) {
		p.nextToken()
		clause.Test = p.parseExpression(LOWEST)
	} else if p.curTokenIs(token.DEFAULT) {
		clause.Test = nil
	} else {
		return nil
	}

	if !p.expectPeek(token.COLON) {
		return nil
	}

	clause.Consequent = []ast.Statement{}

	// Parse statements until we hit case, default, or }
	for !p.peekTokenIs(token.CASE) && !p.peekTokenIs(token.DEFAULT) &&
		!p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil {
			clause.Consequent = append(clause.Consequent, stmt)
		}
	}

	return clause
}

// parseTryStatement parses a try/catch statement
func (p *Parser) parseTryStatement() *ast.TryStmt {
	stmt := &ast.TryStmt{Token: p.curToken}

	// Parse the try block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.TryBlock = p.parseBlockStatement()

	// Expect 'catch'
	if !p.expectPeek(token.CATCH) {
		p.errors = append(p.errors, "expected 'catch' after try block")
		return nil
	}

	// Parse catch parameter: catch (e)
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.CatchParam = p.curToken.Literal

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Parse the catch block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.CatchBlock = p.parseBlockStatement()

	return stmt
}

// parseThrowStatement parses a throw statement
func (p *Parser) parseThrowStatement() *ast.ThrowStmt {
	stmt := &ast.ThrowStmt{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseMapType parses Map<K, V> type syntax
func (p *Parser) parseMapType() *ast.MapType {
	// curToken is MAP, expect <
	if !p.expectPeek(token.LT) {
		return nil
	}

	p.nextToken() // move to key type
	keyType := p.parseType()

	if !p.expectPeek(token.COMMA) {
		return nil
	}

	p.nextToken() // move to value type
	valueType := p.parseType()

	if !p.expectPeek(token.GT) {
		return nil
	}

	return &ast.MapType{
		KeyType:   keyType,
		ValueType: valueType,
	}
}

// parseSetType parses Set<T> type syntax
func (p *Parser) parseSetType() *ast.SetType {
	// curToken is SET, expect <
	if !p.expectPeek(token.LT) {
		return nil
	}

	p.nextToken() // move to element type
	elementType := p.parseType()

	if !p.expectPeek(token.GT) {
		return nil
	}

	return &ast.SetType{
		ElementType: elementType,
	}
}

// parseInterfaceDeclaration parses interface declarations
func (p *Parser) parseInterfaceDeclaration() *ast.InterfaceDecl {
	decl := &ast.InterfaceDecl{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	decl.Name = p.curToken.Literal

	// Check for type parameters <T, U>
	if p.peekTokenIs(token.LT) {
		decl.TypeParams = p.parseTypeParameters()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	decl.Fields = []*ast.InterfaceField{}
	decl.Methods = []*ast.InterfaceMethod{}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		p.nextToken() // move to member name
		if p.curToken.Type != token.IDENT {
			continue
		}

		name := p.curToken.Literal

		if p.peekTokenIs(token.LPAREN) {
			// Method
			p.nextToken() // consume (
			params := p.parseParameterList()
			var retType ast.Type
			if p.peekTokenIs(token.COLON) {
				p.nextToken()
				p.nextToken()
				retType = p.parseType()
			}
			decl.Methods = append(decl.Methods, &ast.InterfaceMethod{
				Name:       name,
				Params:     params,
				ReturnType: retType,
			})
		} else if p.peekTokenIs(token.COLON) {
			// Field
			p.nextToken() // consume :
			p.nextToken() // move to type
			fieldType := p.parseType()
			decl.Fields = append(decl.Fields, &ast.InterfaceField{
				Name:      name,
				FieldType: fieldType,
			})
		}

		// Handle optional separator (semicolon or newline)
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return decl
}

// parseImportDeclaration parses an import statement.
// import { Name1, Name2 } from "go:package"  - Go package import
// import { Name1, Name2 } from "./module"    - Module import
// import Foo from "./module"                  - Default import
// import * as utils from "./utils"            - Namespace import
func (p *Parser) parseImportDeclaration() ast.Statement {
	importToken := p.curToken

	// Check for default import: import Foo from "./module"
	if p.peekTokenIs(token.IDENT) {
		p.nextToken() // move to identifier
		name := p.curToken.Literal

		// Expect 'from'
		if !p.expectPeek(token.FROM) {
			return nil
		}

		// Expect module string
		if !p.expectPeek(token.STRING) {
			return nil
		}

		return &ast.DefaultImport{
			Token: importToken,
			Name:  name,
			Path:  p.curToken.Literal,
		}
	}

	// Check for namespace import: import * as utils from "./utils"
	if p.peekTokenIs(token.STAR) {
		p.nextToken() // move to *

		// Expect 'as'
		if !p.expectPeek(token.AS) {
			return nil
		}

		// Expect alias identifier
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		alias := p.curToken.Literal

		// Expect 'from'
		if !p.expectPeek(token.FROM) {
			return nil
		}

		// Expect module string
		if !p.expectPeek(token.STRING) {
			return nil
		}

		return &ast.NamespaceImport{
			Token: importToken,
			Alias: alias,
			Path:  p.curToken.Literal,
		}
	}

	// Named imports: import { Name1, Name2 } from "./module"
	// Expect '{'
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// Parse imported names
	names := []string{}
	for {
		p.nextToken()
		if p.curToken.Type == token.RBRACE {
			break
		}

		if p.curToken.Type != token.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("expected identifier in import, got %s", p.curToken.Type))
			return nil
		}
		names = append(names, p.curToken.Literal)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
		} else if p.peekTokenIs(token.RBRACE) {
			p.nextToken() // consume closing brace
			break
		}
	}

	// Expect 'from'
	if !p.expectPeek(token.FROM) {
		return nil
	}

	// Expect package/module string
	if !p.expectPeek(token.STRING) {
		return nil
	}

	pathStr := p.curToken.Literal

	// Distinguish between Go package imports and module imports
	if strings.HasPrefix(pathStr, "go:") {
		// Go package import
		return &ast.GoImportDecl{
			Token:   importToken,
			Names:   names,
			Package: strings.TrimPrefix(pathStr, "go:"),
		}
	}

	// Module import (local file)
	return &ast.ModuleImportDecl{
		Token: importToken,
		Names: names,
		Path:  pathStr,
	}
}

// parseExportDeclaration parses an export statement.
// export function foo() { ... }           - Named export
// export class Foo { ... }                 - Named export
// export let x: int = 42                   - Named export
// export default class Foo {}              - Default export
// export default function foo() {}         - Default export
// export { foo, bar } from "./module"      - Re-export
// export * from "./module"                 - Re-export all
func (p *Parser) parseExportDeclaration() ast.Statement {
	exportToken := p.curToken

	// Move to the next token
	p.nextToken()

	// Check for re-export: export * from "./module"
	if p.curToken.Type == token.STAR {
		// Expect 'from'
		if !p.expectPeek(token.FROM) {
			return nil
		}

		// Expect module string
		if !p.expectPeek(token.STRING) {
			return nil
		}

		return &ast.ReExportDecl{
			Token:      exportToken,
			Names:      nil,
			Path:       p.curToken.Literal,
			IsWildcard: true,
		}
	}

	// Check for re-export: export { foo, bar } from "./module"
	if p.curToken.Type == token.LBRACE {
		// Parse exported names
		names := []string{}
		for {
			p.nextToken()
			if p.curToken.Type == token.RBRACE {
				break
			}

			if p.curToken.Type != token.IDENT {
				p.errors = append(p.errors, fmt.Sprintf("expected identifier in export, got %s", p.curToken.Type))
				return nil
			}
			names = append(names, p.curToken.Literal)

			if p.peekTokenIs(token.COMMA) {
				p.nextToken() // consume comma
			} else if p.peekTokenIs(token.RBRACE) {
				p.nextToken() // consume closing brace
				break
			}
		}

		// Check if followed by 'from' for re-export
		if p.peekTokenIs(token.FROM) {
			p.nextToken() // consume 'from'

			// Expect module string
			if !p.expectPeek(token.STRING) {
				return nil
			}

			return &ast.ReExportDecl{
				Token:      exportToken,
				Names:      names,
				Path:       p.curToken.Literal,
				IsWildcard: false,
			}
		}

		// Regular named export (not implemented yet, would need export { foo } syntax)
		p.errors = append(p.errors, "export { ... } without 'from' is not yet supported")
		return nil
	}

	// Check for default export
	if p.curToken.Type == token.DEFAULT {
		p.nextToken() // move to the declaration

		var decl ast.Statement
		switch p.curToken.Type {
		case token.FUNCTION:
			decl = p.parseFunctionDeclaration()
		case token.CLASS:
			decl = p.parseClassDeclaration()
		default:
			p.errors = append(p.errors, fmt.Sprintf("expected function or class after export default, got %s", p.curToken.Type))
			return nil
		}

		if decl == nil {
			return nil
		}

		return &ast.DefaultExport{
			Token: exportToken,
			Decl:  decl,
		}
	}

	// Named export
	var decl ast.Statement
	switch p.curToken.Type {
	case token.FUNCTION:
		decl = p.parseFunctionDeclaration()
	case token.CLASS:
		decl = p.parseClassDeclaration()
	case token.LET:
		decl = p.parseVarDeclaration(false)
	case token.CONST:
		decl = p.parseVarDeclaration(true)
	case token.TYPE:
		decl = p.parseTypeAlias()
	case token.INTERFACE:
		decl = p.parseInterfaceDeclaration()
	default:
		p.errors = append(p.errors, fmt.Sprintf("expected function, class, let, const, type, or interface after export, got %s", p.curToken.Type))
		return nil
	}

	if decl == nil {
		return nil
	}

	return &ast.ExportModifier{
		Token: exportToken,
		Decl:  decl,
	}
}

// parseTypeParameters parses generic type parameters <T, U extends V>
func (p *Parser) parseTypeParameters() []*ast.TypeParam {
	if !p.expectPeek(token.LT) {
		return nil
	}

	params := []*ast.TypeParam{}

	// Parse first type parameter
	p.nextToken()
	param := p.parseTypeParameter()
	if param != nil {
		params = append(params, param)
	}

	// Parse remaining type parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next type param
		param := p.parseTypeParameter()
		if param != nil {
			params = append(params, param)
		}
	}

	if !p.expectPeek(token.GT) {
		return nil
	}

	return params
}

// parseTypeParameter parses a single type parameter (e.g., T or T extends Comparable or T = string)
func (p *Parser) parseTypeParameter() *ast.TypeParam {
	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	param := &ast.TypeParam{
		Name: p.curToken.Literal,
	}

	// Check for constraint: T extends SomeType
	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken() // consume 'extends'
		p.nextToken() // move to constraint type
		param.Constraint = p.parseType()
	}

	// Check for default: T = SomeType
	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume '='
		p.nextToken() // move to default type
		param.Default = p.parseType()
	}

	return param
}

// parseTypeArguments parses type arguments <int, string>
func (p *Parser) parseTypeArguments() []ast.Type {
	if !p.expectPeek(token.LT) {
		return nil
	}

	args := []ast.Type{}

	// Parse first type argument
	p.nextToken()
	arg := p.parseType()
	if arg != nil {
		args = append(args, arg)
	}

	// Parse remaining type arguments
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next type arg
		arg := p.parseType()
		if arg != nil {
			args = append(args, arg)
		}
	}

	if !p.expectPeek(token.GT) {
		return nil
	}

	return args
}

// parseDecorator parses a single decorator: @name, @name(args), @obj.method(args)
func (p *Parser) parseDecorator() *ast.Decorator {
	decorator := &ast.Decorator{Token: p.curToken}

	// Expect identifier after @
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	name := p.curToken.Literal

	// Check for member access: @obj.method
	if p.peekTokenIs(token.DOT) {
		p.nextToken() // consume '.'
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		decorator.Object = name
		decorator.Property = p.curToken.Literal
	} else {
		decorator.Name = name
	}

	// Check for parameterized args: @name(args) or @obj.method(args)
	if p.peekTokenIs(token.LPAREN) {
		p.nextToken() // consume '('
		decorator.Arguments = p.parseExpressionList(token.RPAREN)
		// parseExpressionList already consumed the closing ')'
	}

	return decorator
}

// parseDecoratedDeclaration parses decorators followed by a function declaration.
// @decorator1
// @decorator2
// function name() { ... }
func (p *Parser) parseDecoratedDeclaration() ast.Statement {
	// Collect all decorators
	decorators := []*ast.Decorator{}
	for p.curTokenIs(token.AT) {
		decorator := p.parseDecorator()
		if decorator != nil {
			decorators = append(decorators, decorator)
		}
		p.nextToken() // Move to next token (could be another @ or function/async)
	}

	// Now parse the declaration that follows
	var decl ast.Statement
	if p.curTokenIs(token.FUNCTION) {
		decl = p.parseFunctionDeclaration()
	} else if p.curTokenIs(token.ASYNC) {
		decl = p.parseAsyncFunctionDeclaration()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("line %d: decorator can only be applied to function declarations, got %s", p.curToken.Line, p.curToken.Type))
		return nil
	}

	// Attach decorators to the function declaration
	if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl != nil {
		funcDecl.Decorators = decorators
	}

	return decl
}
