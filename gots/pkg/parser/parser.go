// Package parser implements the parser for GoTS.
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

// Parser parses GoTS source code into an AST.
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

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BoolLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{Token: p.curToken}
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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	call := &ast.CallExpr{Token: p.curToken, Function: function}
	call.Arguments = p.parseExpressionList(token.RPAREN)
	return call
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

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	decl.Name = p.curToken.Literal

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

// Type parsing

func (p *Parser) parseType() ast.Type {
	var typ ast.Type

	switch p.curToken.Type {
	case token.INT_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeInt}
	case token.FLOAT_TYPE:
		typ = &ast.PrimitiveType{Kind: ast.TypeFloat}
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
		namedType := &ast.NamedType{Name: p.curToken.Literal}
		// Check for type arguments: Stack<int>
		if p.peekTokenIs(token.LT) {
			namedType.TypeArgs = p.parseTypeArguments()
		}
		typ = namedType
	case token.LBRACE:
		typ = p.parseObjectType()
	case token.LPAREN:
		typ = p.parseFunctionType()
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

	// Check for nullable type
	if p.peekTokenIs(token.PIPE) {
		p.nextToken()
		if !p.expectPeek(token.NULL) {
			return nil
		}
		typ = &ast.NullableType{Inner: typ}
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

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	decl.Methods = []*ast.InterfaceMethod{}

	for !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.EOF) {
		method := p.parseInterfaceMethod()
		if method != nil {
			decl.Methods = append(decl.Methods, method)
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

// parseInterfaceMethod parses a method signature in an interface
func (p *Parser) parseInterfaceMethod() *ast.InterfaceMethod {
	p.nextToken() // move to method name

	if p.curToken.Type != token.IDENT {
		p.peekError(token.IDENT)
		return nil
	}

	method := &ast.InterfaceMethod{
		Name: p.curToken.Literal,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	method.Params = p.parseParameterList()

	// Parse return type
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // consume ':'
		p.nextToken() // move to type
		method.ReturnType = p.parseType()
	}

	return method
}

// parseImportDeclaration parses an import statement.
// import { Name1, Name2 } from "go:package"  - Go package import
// import { Name1, Name2 } from "./module"    - Module import
func (p *Parser) parseImportDeclaration() ast.Statement {
	importToken := p.curToken

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
// export function foo() { ... }
// export class Foo { ... }
// export let x: int = 42
func (p *Parser) parseExportDeclaration() *ast.ExportModifier {
	exportToken := p.curToken

	// Move to the next token which should be the declaration
	p.nextToken()

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

// parseTypeParameter parses a single type parameter (e.g., T or T extends Comparable)
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
