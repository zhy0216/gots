// Package declaration provides parsing for .d.gts declaration files.
package declaration

import (
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/token"
)

// Parser parses .d.gts declaration files.
type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

// New creates a new declaration parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens to initialize curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// NewFromSource creates a parser from source code.
func NewFromSource(source string) *Parser {
	l := lexer.New(source)
	return New(l)
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

// Errors returns the list of parsing errors.
func (p *Parser) Errors() []string {
	return p.errors
}

// Parse parses a declaration file.
func (p *Parser) Parse() *ast.DeclarationFile {
	file := &ast.DeclarationFile{
		Modules: []*ast.DeclareModule{},
	}

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.DECLARE) {
			module := p.parseDeclareModule()
			if module != nil {
				file.Modules = append(file.Modules, module)
			}
		} else {
			p.nextToken()
		}
	}

	return file
}

// parseDeclareModule parses: declare module "name" { ... }
func (p *Parser) parseDeclareModule() *ast.DeclareModule {
	module := &ast.DeclareModule{Token: p.curToken}

	// Expect 'module'
	if !p.expectPeek(token.MODULE) {
		return nil
	}

	// Expect module name (string)
	if !p.expectPeek(token.STRING) {
		return nil
	}
	module.Name = p.curToken.Literal

	// Expect '{'
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	// Parse members until '}'
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		decl := p.parseDeclaration()
		if decl != nil {
			module.Members = append(module.Members, decl)
		}
		p.nextToken()
	}

	return module
}

// parseDeclaration parses a single declaration (function, type, interface, const, class)
func (p *Parser) parseDeclaration() ast.Declaration {
	switch p.curToken.Type {
	case token.FUNCTION:
		return p.parseDeclareFunction()
	case token.TYPE:
		return p.parseDeclareType()
	case token.INTERFACE:
		return p.parseDeclareInterface()
	case token.CONST:
		return p.parseDeclareConst()
	case token.CLASS:
		return p.parseDeclareClass()
	default:
		return nil
	}
}

// parseDeclareFunction parses: function name(params): returnType
func (p *Parser) parseDeclareFunction() *ast.DeclareFunction {
	fn := &ast.DeclareFunction{Token: p.curToken}

	// Expect function name
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	fn.Name = p.curToken.Literal

	// Expect '('
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Parse parameters
	fn.Params, fn.Variadic = p.parseParameters()

	// Expect ':'
	if !p.expectPeek(token.COLON) {
		return nil
	}

	// Parse return type
	p.nextToken()
	fn.ReturnType = p.parseType()

	return fn
}

// parseParameters parses function parameters
func (p *Parser) parseParameters() ([]*ast.Parameter, bool) {
	params := []*ast.Parameter{}
	variadic := false

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params, false
	}

	p.nextToken()

	for {
		// Check for variadic parameter (...)
		if p.curTokenIs(token.DOT) && p.peekTokenIs(token.DOT) {
			p.nextToken() // skip first dot
			p.nextToken() // skip second dot
			if p.curTokenIs(token.DOT) {
				p.nextToken() // skip third dot
			}
			variadic = true
		}

		param := &ast.Parameter{}

		if !p.curTokenIs(token.IDENT) {
			p.errors = append(p.errors, fmt.Sprintf("line %d: expected parameter name, got %s", p.curToken.Line, p.curToken.Type))
			return params, variadic
		}
		param.Name = p.curToken.Literal

		if !p.expectPeek(token.COLON) {
			return params, variadic
		}

		p.nextToken()
		param.ParamType = p.parseType()
		params = append(params, param)

		if !p.peekTokenIs(token.COMMA) {
			break
		}
		p.nextToken() // skip comma
		p.nextToken() // move to next param
	}

	if !p.expectPeek(token.RPAREN) {
		return params, variadic
	}

	return params, variadic
}

// parseType parses a type expression
func (p *Parser) parseType() ast.Type {
	var baseType ast.Type

	switch p.curToken.Type {
	case token.INT_TYPE:
		baseType = &ast.PrimitiveType{Kind: ast.TypeInt}
	case token.FLOAT_TYPE:
		baseType = &ast.PrimitiveType{Kind: ast.TypeFloat}
	case token.STRING_TYPE:
		baseType = &ast.PrimitiveType{Kind: ast.TypeString}
	case token.BOOLEAN_TYPE:
		baseType = &ast.PrimitiveType{Kind: ast.TypeBoolean}
	case token.VOID_TYPE:
		baseType = &ast.PrimitiveType{Kind: ast.TypeVoid}
	case token.NULL:
		baseType = &ast.PrimitiveType{Kind: ast.TypeNull}
	case token.IDENT:
		name := p.curToken.Literal
		if name == "any" {
			baseType = &ast.AnyType{}
		} else if name == "byte" {
			baseType = &ast.ByteType{}
		} else {
			baseType = &ast.NamedType{Name: name}
			// Check for type arguments <T, U>
			if p.peekTokenIs(token.LT) {
				p.nextToken() // consume <
				typeArgs := p.parseTypeArguments()
				baseType.(*ast.NamedType).TypeArgs = typeArgs
			}
		}
	case token.LBRACE:
		baseType = p.parseObjectType()
	case token.LPAREN:
		baseType = p.parseTupleOrFunctionType()
	default:
		p.errors = append(p.errors, fmt.Sprintf("line %d: unexpected type token %s", p.curToken.Line, p.curToken.Type))
		return nil
	}

	// Check for array type []
	for p.peekTokenIs(token.LBRACKET) {
		p.nextToken() // consume [
		if !p.expectPeek(token.RBRACKET) {
			return baseType
		}
		baseType = &ast.ArrayType{ElementType: baseType}
	}

	// Check for nullable type | null
	if p.peekTokenIs(token.PIPE) {
		p.nextToken() // consume |
		if p.peekTokenIs(token.NULL) {
			p.nextToken() // consume null
			baseType = &ast.NullableType{Inner: baseType}
		}
	}

	return baseType
}

// parseTypeArguments parses <T, U, ...>
func (p *Parser) parseTypeArguments() []ast.Type {
	args := []ast.Type{}

	p.nextToken() // move past <
	for !p.curTokenIs(token.GT) && !p.curTokenIs(token.EOF) {
		t := p.parseType()
		if t != nil {
			args = append(args, t)
		}
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // skip comma
			p.nextToken() // move to next type
		} else if p.peekTokenIs(token.GT) {
			p.nextToken() // move to >
		} else {
			p.nextToken()
		}
	}

	return args
}

// parseObjectType parses { prop: type, ... }
func (p *Parser) parseObjectType() ast.Type {
	obj := &ast.ObjectType{Properties: []*ast.ObjectTypeProperty{}}

	p.nextToken() // skip {

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if !p.curTokenIs(token.IDENT) {
			p.nextToken()
			continue
		}

		prop := &ast.ObjectTypeProperty{Name: p.curToken.Literal}

		// Check if it's a method (has parentheses)
		if p.peekTokenIs(token.LPAREN) {
			// It's a method - parse as function type
			p.nextToken() // move to (
			params, _ := p.parseParameters()

			var returnType ast.Type = &ast.PrimitiveType{Kind: ast.TypeVoid}
			if p.peekTokenIs(token.COLON) {
				p.nextToken() // consume :
				p.nextToken()
				returnType = p.parseType()
			}

			paramTypes := make([]ast.Type, len(params))
			for i, param := range params {
				paramTypes[i] = param.ParamType
			}
			prop.PropType = &ast.FunctionType{
				ParamTypes: paramTypes,
				ReturnType: returnType,
			}
		} else {
			// Regular property
			if !p.expectPeek(token.COLON) {
				return obj
			}
			p.nextToken()
			prop.PropType = p.parseType()
		}

		obj.Properties = append(obj.Properties, prop)

		// Skip comma or semicolon
		if p.peekTokenIs(token.COMMA) || p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		p.nextToken()
	}

	return obj
}

// parseTupleOrFunctionType parses (T, U) for tuples or (a: T) => U for functions
func (p *Parser) parseTupleOrFunctionType() ast.Type {
	p.nextToken() // skip (

	// Empty parens
	if p.curTokenIs(token.RPAREN) {
		if p.peekTokenIs(token.ARROW) {
			p.nextToken() // consume =>
			p.nextToken()
			returnType := p.parseType()
			return &ast.FunctionType{ParamTypes: []ast.Type{}, ReturnType: returnType}
		}
		return &ast.TupleType{Elements: []ast.Type{}}
	}

	// Check if it looks like function params (name: type) or tuple (type, type)
	// If first token is IDENT followed by COLON, it's function params
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON) {
		// Parse as function type
		var paramTypes []ast.Type
		for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
			// Skip parameter name
			if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON) {
				p.nextToken() // skip name
				p.nextToken() // skip :
			}
			t := p.parseType()
			if t != nil {
				paramTypes = append(paramTypes, t)
			}
			if p.peekTokenIs(token.COMMA) {
				p.nextToken() // skip comma
				p.nextToken()
			} else if p.peekTokenIs(token.RPAREN) {
				p.nextToken()
			} else {
				p.nextToken()
			}
		}

		if p.peekTokenIs(token.ARROW) {
			p.nextToken() // consume =>
			p.nextToken()
			returnType := p.parseType()
			return &ast.FunctionType{ParamTypes: paramTypes, ReturnType: returnType}
		}
		return &ast.TupleType{Elements: paramTypes}
	}

	// Parse as tuple
	var elements []ast.Type
	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		t := p.parseType()
		if t != nil {
			elements = append(elements, t)
		}
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
		} else if p.peekTokenIs(token.RPAREN) {
			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	return &ast.TupleType{Elements: elements}
}

// parseDeclareType parses: type Name = Type
func (p *Parser) parseDeclareType() *ast.DeclareType {
	dt := &ast.DeclareType{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	dt.Name = p.curToken.Literal

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	dt.AliasType = p.parseType()

	return dt
}

// parseDeclareInterface parses: interface Name { methods }
func (p *Parser) parseDeclareInterface() *ast.DeclareInterface {
	di := &ast.DeclareInterface{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	di.Name = p.curToken.Literal

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	di.Methods = []*ast.InterfaceMethod{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) {
			method := p.parseInterfaceMethod()
			if method != nil {
				di.Methods = append(di.Methods, method)
			}
		}
		p.nextToken()
	}

	return di
}

// parseInterfaceMethod parses a method signature
func (p *Parser) parseInterfaceMethod() *ast.InterfaceMethod {
	method := &ast.InterfaceMethod{Name: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	method.Params, _ = p.parseParameters()

	if p.peekTokenIs(token.COLON) {
		p.nextToken()
		p.nextToken()
		method.ReturnType = p.parseType()
	} else {
		method.ReturnType = &ast.PrimitiveType{Kind: ast.TypeVoid}
	}

	return method
}

// parseDeclareConst parses: const Name: Type
func (p *Parser) parseDeclareConst() *ast.DeclareConst {
	dc := &ast.DeclareConst{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	dc.Name = p.curToken.Literal

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()
	dc.ConstType = p.parseType()

	return dc
}

// parseDeclareClass parses: class Name { fields and methods }
func (p *Parser) parseDeclareClass() *ast.DeclareClass {
	dc := &ast.DeclareClass{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	dc.Name = p.curToken.Literal

	// Check for extends
	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken()
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		dc.SuperClass = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	dc.Fields = []*ast.Field{}
	dc.Methods = []*ast.Method{}
	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.IDENT) {
			name := p.curToken.Literal
			if p.peekTokenIs(token.LPAREN) {
				// It's a method
				method := &ast.Method{Name: name}
				p.nextToken() // move to (
				params, _ := p.parseParameters()
				method.Params = params

				if p.peekTokenIs(token.COLON) {
					p.nextToken()
					p.nextToken()
					method.ReturnType = p.parseType()
				} else {
					method.ReturnType = &ast.PrimitiveType{Kind: ast.TypeVoid}
				}
				dc.Methods = append(dc.Methods, method)
			} else if p.peekTokenIs(token.COLON) {
				// It's a field
				p.nextToken()
				p.nextToken()
				fieldType := p.parseType()
				dc.Fields = append(dc.Fields, &ast.Field{Name: name, FieldType: fieldType})
			}
		}
		p.nextToken()
	}

	return dc
}

// ModuleInfo contains the parsed type information for a module
type ModuleInfo struct {
	Name      string
	Functions map[string]*FunctionInfo
	Types     map[string]ast.Type
	Constants map[string]ast.Type
	Interfaces map[string]*InterfaceInfo
	Classes   map[string]*ClassInfo
}

// FunctionInfo contains function signature information
type FunctionInfo struct {
	Name       string
	Params     []*ast.Parameter
	ReturnType ast.Type
	Variadic   bool
}

// InterfaceInfo contains interface information
type InterfaceInfo struct {
	Name    string
	Methods []*ast.InterfaceMethod
}

// ClassInfo contains class information
type ClassInfo struct {
	Name       string
	SuperClass string
	Fields     []*ast.Field
	Methods    []*ast.Method
}

// ExtractModuleInfo extracts type information from a parsed declaration file
func ExtractModuleInfo(file *ast.DeclarationFile) map[string]*ModuleInfo {
	modules := make(map[string]*ModuleInfo)

	for _, mod := range file.Modules {
		info := &ModuleInfo{
			Name:       mod.Name,
			Functions:  make(map[string]*FunctionInfo),
			Types:      make(map[string]ast.Type),
			Constants:  make(map[string]ast.Type),
			Interfaces: make(map[string]*InterfaceInfo),
			Classes:    make(map[string]*ClassInfo),
		}

		for _, decl := range mod.Members {
			switch d := decl.(type) {
			case *ast.DeclareFunction:
				info.Functions[d.Name] = &FunctionInfo{
					Name:       d.Name,
					Params:     d.Params,
					ReturnType: d.ReturnType,
					Variadic:   d.Variadic,
				}
			case *ast.DeclareType:
				info.Types[d.Name] = d.AliasType
			case *ast.DeclareConst:
				info.Constants[d.Name] = d.ConstType
			case *ast.DeclareInterface:
				info.Interfaces[d.Name] = &InterfaceInfo{
					Name:    d.Name,
					Methods: d.Methods,
				}
			case *ast.DeclareClass:
				info.Classes[d.Name] = &ClassInfo{
					Name:       d.Name,
					SuperClass: d.SuperClass,
					Fields:     d.Fields,
					Methods:    d.Methods,
				}
			}
		}

		// Strip "go:" prefix for easier lookup
		name := mod.Name
		if strings.HasPrefix(name, "go:") {
			name = strings.TrimPrefix(name, "go:")
		}
		modules[name] = info
	}

	return modules
}
