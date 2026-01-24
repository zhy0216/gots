package typed

import (
	"fmt"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/token"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Builder transforms an AST into a TypedAST while performing type checking.
type Builder struct {
	errors       []*Error
	scope        *Scope
	typeAliases  map[string]types.Type
	classes      map[string]*types.Class
	currentFunc  *types.Function
	currentClass *types.Class
	narrowing    map[string]types.Type
	loopDepth    int
	scopeDepth   int
	constVars    map[string]bool // Track which variables are const (scoped)
}

// Error represents a type checking error.
type Error struct {
	Line    int
	Column  int
	Message string
}

func (e *Error) String() string {
	return fmt.Sprintf("line %d, col %d: %s", e.Line, e.Column, e.Message)
}

// Scope represents a lexical scope.
type Scope struct {
	parent    *Scope
	bindings  map[string]types.Type
	constVars map[string]bool
}

func newScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		bindings:  make(map[string]types.Type),
		constVars: make(map[string]bool),
	}
}

func (s *Scope) define(name string, typ types.Type) {
	s.bindings[name] = typ
}

func (s *Scope) defineConst(name string, typ types.Type) {
	s.bindings[name] = typ
	s.constVars[name] = true
}

func (s *Scope) isConst(name string) bool {
	if isConst, ok := s.constVars[name]; ok {
		return isConst
	}
	if s.parent != nil {
		return s.parent.isConst(name)
	}
	return false
}

func (s *Scope) lookup(name string) (types.Type, bool) {
	if typ, ok := s.bindings[name]; ok {
		return typ, true
	}
	if s.parent != nil {
		return s.parent.lookup(name)
	}
	return nil, false
}

// NewBuilder creates a new typed AST builder.
func NewBuilder() *Builder {
	b := &Builder{
		errors:      []*Error{},
		scope:       newScope(nil),
		typeAliases: make(map[string]types.Type),
		classes:     make(map[string]*types.Class),
		narrowing:   make(map[string]types.Type),
		constVars:   make(map[string]bool),
	}

	// Pre-define built-in type aliases
	// Function is a generic function type (any) => any for compatibility
	b.typeAliases["Function"] = &types.Function{
		Params:     []*types.Param{{Name: "x", Type: types.AnyType}},
		ReturnType: types.AnyType,
	}

	return b
}

// Errors returns the list of type errors.
func (b *Builder) Errors() []*Error {
	return b.errors
}

// HasErrors returns true if there are errors.
func (b *Builder) HasErrors() bool {
	return len(b.errors) > 0
}

func (b *Builder) error(line, col int, format string, args ...interface{}) {
	b.errors = append(b.errors, &Error{
		Line:    line,
		Column:  col,
		Message: fmt.Sprintf(format, args...),
	})
}

func (b *Builder) pushScope() {
	b.scope = newScope(b.scope)
	b.scopeDepth++
}

func (b *Builder) popScope() {
	b.scope = b.scope.parent
	b.scopeDepth--
}

// Build transforms an AST program into a typed program.
func (b *Builder) Build(program *ast.Program) *Program {
	// First pass: collect type aliases and classes
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			b.collectTypeAlias(s)
		case *ast.ClassDecl:
			b.collectClass(s)
		}
	}

	// Second pass: resolve types
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			b.resolveTypeAlias(s)
		case *ast.ClassDecl:
			b.resolveClass(s)
		}
	}

	// Third pass: build typed AST
	result := &Program{
		TypeAliases: make([]*TypeAlias, 0),
		Classes:     make([]*ClassDecl, 0),
		Functions:   make([]*FuncDecl, 0),
		TopLevel:    make([]Stmt, 0),
	}

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			alias := &TypeAlias{
				Name:     s.Name,
				Resolved: b.typeAliases[s.Name],
			}
			result.TypeAliases = append(result.TypeAliases, alias)

		case *ast.ClassDecl:
			classDecl := b.buildClassDecl(s)
			result.Classes = append(result.Classes, classDecl)

		case *ast.FuncDecl:
			funcDecl := b.buildFuncDecl(s)
			result.Functions = append(result.Functions, funcDecl)

		default:
			typedStmt := b.buildStmt(stmt)
			if typedStmt != nil {
				result.TopLevel = append(result.TopLevel, typedStmt)
			}
		}
	}

	return result
}

// ----------------------------------------------------------------------------
// Collection Pass
// ----------------------------------------------------------------------------

func (b *Builder) collectTypeAlias(decl *ast.TypeAliasDecl) {
	b.typeAliases[decl.Name] = &types.Alias{Name: decl.Name}
}

func (b *Builder) resolveTypeAlias(decl *ast.TypeAliasDecl) {
	alias := b.typeAliases[decl.Name].(*types.Alias)
	alias.Resolved = b.resolveType(decl.AliasType)
}

func (b *Builder) collectClass(decl *ast.ClassDecl) {
	b.classes[decl.Name] = &types.Class{
		Name:    decl.Name,
		Fields:  make(map[string]*types.Field),
		Methods: make(map[string]*types.Method),
	}
}

func (b *Builder) resolveClass(decl *ast.ClassDecl) {
	class := b.classes[decl.Name]

	if decl.SuperClass != "" {
		if super, ok := b.classes[decl.SuperClass]; ok {
			class.Super = super
		} else {
			b.error(decl.Token.Line, decl.Token.Column, "unknown superclass: %s", decl.SuperClass)
		}
	}

	for _, f := range decl.Fields {
		class.Fields[f.Name] = &types.Field{
			Name: f.Name,
			Type: b.resolveType(f.FieldType),
		}
	}

	if decl.Constructor != nil {
		params := make([]*types.Param, len(decl.Constructor.Params))
		for i, p := range decl.Constructor.Params {
			params[i] = &types.Param{
				Name: p.Name,
				Type: b.resolveType(p.ParamType),
			}
		}
		class.Constructor = &types.Function{
			Params:     params,
			ReturnType: types.VoidType,
		}
	}

	for _, m := range decl.Methods {
		params := make([]*types.Param, len(m.Params))
		for i, p := range m.Params {
			params[i] = &types.Param{
				Name: p.Name,
				Type: b.resolveType(p.ParamType),
			}
		}
		class.Methods[m.Name] = &types.Method{
			Name:       m.Name,
			Params:     params,
			ReturnType: b.resolveType(m.ReturnType),
		}
	}
}

// ----------------------------------------------------------------------------
// Type Resolution
// ----------------------------------------------------------------------------

func (b *Builder) resolveType(astType ast.Type) types.Type {
	if astType == nil {
		return types.VoidType
	}

	switch t := astType.(type) {
	case *ast.PrimitiveType:
		switch t.Kind {
		case ast.TypeInt:
			return types.IntType
		case ast.TypeFloat:
			return types.FloatType
		case ast.TypeString:
			return types.StringType
		case ast.TypeBoolean:
			return types.BooleanType
		case ast.TypeVoid:
			return types.VoidType
		case ast.TypeNull:
			return types.NullType
		}

	case *ast.ArrayType:
		return &types.Array{Element: b.resolveType(t.ElementType)}

	case *ast.ObjectType:
		props := make(map[string]*types.Property)
		for _, p := range t.Properties {
			props[p.Name] = &types.Property{
				Name: p.Name,
				Type: b.resolveType(p.PropType),
			}
		}
		return &types.Object{Properties: props}

	case *ast.FunctionType:
		params := make([]*types.Param, len(t.ParamTypes))
		for i, pt := range t.ParamTypes {
			params[i] = &types.Param{
				Name: fmt.Sprintf("arg%d", i),
				Type: b.resolveType(pt),
			}
		}
		return &types.Function{
			Params:     params,
			ReturnType: b.resolveType(t.ReturnType),
		}

	case *ast.NullableType:
		return &types.Nullable{Inner: b.resolveType(t.Inner)}

	case *ast.NamedType:
		if alias, ok := b.typeAliases[t.Name]; ok {
			return alias
		}
		if class, ok := b.classes[t.Name]; ok {
			return class
		}
		b.error(0, 0, "unknown type: %s", t.Name)
		return types.AnyType
	}

	return types.AnyType
}

// ----------------------------------------------------------------------------
// Statement Building
// ----------------------------------------------------------------------------

func (b *Builder) buildStmt(stmt ast.Statement) Stmt {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		return &ExprStmt{Expr: b.buildExpr(s.Expr)}

	case *ast.VarDecl:
		return b.buildVarDecl(s)

	case *ast.Block:
		return b.buildBlock(s)

	case *ast.IfStmt:
		return b.buildIfStmt(s)

	case *ast.WhileStmt:
		return b.buildWhileStmt(s)

	case *ast.ForStmt:
		return b.buildForStmt(s)

	case *ast.ForOfStmt:
		return b.buildForOfStmt(s)

	case *ast.SwitchStmt:
		return b.buildSwitchStmt(s)

	case *ast.ReturnStmt:
		return b.buildReturnStmt(s)

	case *ast.BreakStmt:
		return &BreakStmt{}

	case *ast.ContinueStmt:
		return &ContinueStmt{}

	case *ast.TryStmt:
		return b.buildTryStmt(s)

	case *ast.ThrowStmt:
		return b.buildThrowStmt(s)

	case *ast.FuncDecl:
		return b.buildFuncDecl(s)

	case *ast.ClassDecl:
		return b.buildClassDecl(s)

	case *ast.TypeAliasDecl:
		return nil // Already processed
	}

	return nil
}

func (b *Builder) buildVarDecl(decl *ast.VarDecl) *VarDecl {
	var varType types.Type
	var init Expr

	if decl.VarType != nil {
		varType = b.resolveType(decl.VarType)
	}

	if decl.Value != nil {
		init = b.buildExpr(decl.Value)
		if varType == nil {
			varType = init.Type()
		} else {
			if !types.IsAssignableTo(init.Type(), varType) {
				b.error(decl.Token.Line, decl.Token.Column,
					"cannot assign %s to %s", init.Type().String(), varType.String())
			}
		}
	}

	if varType == nil {
		b.error(decl.Token.Line, decl.Token.Column,
			"variable declaration requires type annotation or initializer")
		varType = types.AnyType
	}

	if decl.IsConst {
		b.scope.defineConst(decl.Name, varType)
	} else {
		b.scope.define(decl.Name, varType)
	}

	return &VarDecl{
		Name:    decl.Name,
		VarType: varType,
		Init:    init,
		IsConst: decl.IsConst,
	}
}

func (b *Builder) buildBlock(block *ast.Block) *BlockStmt {
	b.pushScope()
	stmts := make([]Stmt, 0, len(block.Statements))
	for _, s := range block.Statements {
		if typedStmt := b.buildStmt(s); typedStmt != nil {
			stmts = append(stmts, typedStmt)
		}
	}
	b.popScope()
	return &BlockStmt{Stmts: stmts}
}

func (b *Builder) buildIfStmt(stmt *ast.IfStmt) *IfStmt {
	cond := b.buildExpr(stmt.Condition)

	// Allow boolean or any type for conditions
	condType := types.Unwrap(cond.Type())
	if !condType.Equals(types.BooleanType) && !condType.Equals(types.AnyType) {
		b.error(stmt.Token.Line, stmt.Token.Column,
			"condition must be boolean, got %s", cond.Type().String())
	}

	// Handle type narrowing for null checks
	savedNarrowing := b.narrowing
	b.narrowing = make(map[string]types.Type)

	if binary, ok := stmt.Condition.(*ast.BinaryExpr); ok {
		if binary.Op == token.NEQ {
			if ident, ok := binary.Left.(*ast.Identifier); ok {
				if _, ok := binary.Right.(*ast.NullLiteral); ok {
					if varType, found := b.scope.lookup(ident.Name); found {
						if nullable, ok := varType.(*types.Nullable); ok {
							b.narrowing[ident.Name] = nullable.Inner
						}
					}
				}
			}
		}
	}

	then := b.buildBlock(stmt.Consequence)
	b.narrowing = savedNarrowing

	var elseStmt Stmt
	if stmt.Alternative != nil {
		elseStmt = b.buildStmt(stmt.Alternative)
	}

	return &IfStmt{
		Condition: cond,
		Then:      then,
		Else:      elseStmt,
	}
}

func (b *Builder) buildWhileStmt(stmt *ast.WhileStmt) *WhileStmt {
	cond := b.buildExpr(stmt.Condition)

	if !cond.Type().Equals(types.BooleanType) {
		b.error(stmt.Token.Line, stmt.Token.Column,
			"condition must be boolean, got %s", cond.Type().String())
	}

	b.loopDepth++
	body := b.buildBlock(stmt.Body)
	b.loopDepth--

	return &WhileStmt{
		Condition: cond,
		Body:      body,
	}
}

func (b *Builder) buildForStmt(stmt *ast.ForStmt) *ForStmt {
	b.pushScope()

	var init *VarDecl
	if stmt.Init != nil {
		init = b.buildVarDecl(stmt.Init)
	}

	var cond Expr
	if stmt.Condition != nil {
		cond = b.buildExpr(stmt.Condition)
		if !cond.Type().Equals(types.BooleanType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"condition must be boolean, got %s", cond.Type().String())
		}
	}

	var update Expr
	if stmt.Update != nil {
		update = b.buildExpr(stmt.Update)
	}

	b.loopDepth++
	body := b.buildBlock(stmt.Body)
	b.loopDepth--

	b.popScope()

	return &ForStmt{
		Init:      init,
		Condition: cond,
		Update:    update,
		Body:      body,
	}
}

func (b *Builder) buildForOfStmt(stmt *ast.ForOfStmt) *ForOfStmt {
	iterable := b.buildExpr(stmt.Iterable)
	iterType := types.Unwrap(iterable.Type())

	var elemType types.Type
	switch t := iterType.(type) {
	case *types.Array:
		elemType = t.Element
	default:
		if iterType.Equals(types.StringType) {
			elemType = types.StringType
		} else {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"for-of requires array or string, got %s", iterType.String())
			elemType = types.AnyType
		}
	}

	b.pushScope()

	if stmt.Variable.VarType != nil {
		declType := b.resolveType(stmt.Variable.VarType)
		if !types.IsAssignableTo(elemType, declType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"cannot assign %s to %s", elemType.String(), declType.String())
		}
		elemType = declType
	}

	b.scope.define(stmt.Variable.Name, elemType)

	varDecl := &VarDecl{
		Name:    stmt.Variable.Name,
		VarType: elemType,
		IsConst: false,
	}

	b.loopDepth++
	body := b.buildBlock(stmt.Body)
	b.loopDepth--

	b.popScope()

	return &ForOfStmt{
		Variable:    varDecl,
		Iterable:    iterable,
		ElementType: elemType,
		Body:        body,
	}
}

func (b *Builder) buildSwitchStmt(stmt *ast.SwitchStmt) *SwitchStmt {
	discriminant := b.buildExpr(stmt.Discriminant)

	cases := make([]*CaseClause, len(stmt.Cases))
	for i, c := range stmt.Cases {
		var test Expr
		if c.Test != nil {
			test = b.buildExpr(c.Test)
		}

		stmts := make([]Stmt, 0, len(c.Consequent))
		b.loopDepth++ // break is valid in switch
		for _, s := range c.Consequent {
			if typedStmt := b.buildStmt(s); typedStmt != nil {
				stmts = append(stmts, typedStmt)
			}
		}
		b.loopDepth--

		cases[i] = &CaseClause{
			Test:  test,
			Stmts: stmts,
		}
	}

	return &SwitchStmt{
		Discriminant: discriminant,
		Cases:        cases,
	}
}

func (b *Builder) buildTryStmt(stmt *ast.TryStmt) *TryStmt {
	// Build the try block
	tryBlock := b.buildBlock(stmt.TryBlock)

	// Build the catch block with the catch parameter in scope
	b.pushScope()
	// The catch parameter has type 'any' since any error can be thrown
	b.scope.define(stmt.CatchParam, types.AnyType)
	catchBlock := &BlockStmt{}
	stmts := make([]Stmt, 0, len(stmt.CatchBlock.Statements))
	for _, s := range stmt.CatchBlock.Statements {
		if typedStmt := b.buildStmt(s); typedStmt != nil {
			stmts = append(stmts, typedStmt)
		}
	}
	catchBlock.Stmts = stmts
	b.popScope()

	catchParam := &VarDecl{
		Name:    stmt.CatchParam,
		VarType: types.AnyType,
		IsConst: false,
	}

	return &TryStmt{
		TryBlock:   tryBlock,
		CatchParam: catchParam,
		CatchBlock: catchBlock,
	}
}

func (b *Builder) buildThrowStmt(stmt *ast.ThrowStmt) *ThrowStmt {
	value := b.buildExpr(stmt.Value)
	return &ThrowStmt{Value: value}
}

func (b *Builder) buildReturnStmt(stmt *ast.ReturnStmt) *ReturnStmt {
	if b.currentFunc == nil {
		b.error(stmt.Token.Line, stmt.Token.Column, "return outside function")
		return &ReturnStmt{}
	}

	var value Expr
	if stmt.Value != nil {
		value = b.buildExpr(stmt.Value)
		if !types.IsAssignableTo(value.Type(), b.currentFunc.ReturnType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"cannot return %s, expected %s", value.Type().String(), b.currentFunc.ReturnType.String())
		}
	} else {
		if !b.currentFunc.ReturnType.Equals(types.VoidType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"missing return value, expected %s", b.currentFunc.ReturnType.String())
		}
	}

	return &ReturnStmt{Value: value}
}

func (b *Builder) buildFuncDecl(decl *ast.FuncDecl) *FuncDecl {
	params := make([]*Param, len(decl.Params))
	typeParams := make([]*types.Param, len(decl.Params))
	for i, p := range decl.Params {
		paramType := b.resolveType(p.ParamType)
		params[i] = &Param{
			Name: p.Name,
			Type: paramType,
		}
		typeParams[i] = &types.Param{
			Name: p.Name,
			Type: paramType,
		}
	}

	returnType := b.resolveType(decl.ReturnType)
	funcType := &types.Function{
		Params:     typeParams,
		ReturnType: returnType,
	}

	b.scope.define(decl.Name, funcType)

	savedFunc := b.currentFunc
	b.currentFunc = funcType

	b.pushScope()
	for _, p := range params {
		b.scope.define(p.Name, p.Type)
	}

	body := b.buildBlock(decl.Body)

	b.popScope()
	b.currentFunc = savedFunc

	return &FuncDecl{
		Name:       decl.Name,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
	}
}

func (b *Builder) buildClassDecl(decl *ast.ClassDecl) *ClassDecl {
	class := b.classes[decl.Name]

	savedClass := b.currentClass
	b.currentClass = class

	fields := make([]*FieldDecl, len(decl.Fields))
	for i, f := range decl.Fields {
		fields[i] = &FieldDecl{
			Name: f.Name,
			Type: b.resolveType(f.FieldType),
		}
	}

	var constructor *ConstructorDecl
	if decl.Constructor != nil {
		params := make([]*Param, len(decl.Constructor.Params))
		for i, p := range decl.Constructor.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: b.resolveType(p.ParamType),
			}
		}

		b.pushScope()
		for _, p := range params {
			b.scope.define(p.Name, p.Type)
		}
		body := b.buildBlock(decl.Constructor.Body)
		b.popScope()

		constructor = &ConstructorDecl{
			Params: params,
			Body:   body,
		}
	}

	methods := make([]*MethodDecl, len(decl.Methods))
	for i, m := range decl.Methods {
		method := class.Methods[m.Name]

		params := make([]*Param, len(m.Params))
		for j, p := range m.Params {
			params[j] = &Param{
				Name: p.Name,
				Type: b.resolveType(p.ParamType),
			}
		}

		savedFunc := b.currentFunc
		b.currentFunc = &types.Function{
			Params:     method.Params,
			ReturnType: method.ReturnType,
		}

		b.pushScope()
		for _, p := range params {
			b.scope.define(p.Name, p.Type)
		}
		body := b.buildBlock(m.Body)
		b.popScope()

		b.currentFunc = savedFunc

		methods[i] = &MethodDecl{
			Name:       m.Name,
			Params:     params,
			ReturnType: method.ReturnType,
			Body:       body,
		}
	}

	b.currentClass = savedClass

	return &ClassDecl{
		Name:        decl.Name,
		Super:       decl.SuperClass,
		SuperClass:  class.Super,
		Fields:      fields,
		Constructor: constructor,
		Methods:     methods,
	}
}

// ----------------------------------------------------------------------------
// Expression Building
// ----------------------------------------------------------------------------

func (b *Builder) buildExpr(expr ast.Expression) Expr {
	if expr == nil {
		return &NullLit{ExprType: types.NullType}
	}

	switch e := expr.(type) {
	case *ast.NumberLiteral:
		// Check if number is an integer (no decimal part)
		if e.Value == float64(int64(e.Value)) {
			return &NumberLit{Value: e.Value, ExprType: types.IntType}
		}
		return &NumberLit{Value: e.Value, ExprType: types.FloatType}

	case *ast.StringLiteral:
		return &StringLit{Value: e.Value, ExprType: types.StringType}

	case *ast.BoolLiteral:
		return &BoolLit{Value: e.Value, ExprType: types.BooleanType}

	case *ast.NullLiteral:
		return &NullLit{ExprType: types.NullType}

	case *ast.Identifier:
		return b.buildIdent(e)

	case *ast.BinaryExpr:
		return b.buildBinaryExpr(e)

	case *ast.UnaryExpr:
		return b.buildUnaryExpr(e)

	case *ast.CallExpr:
		return b.buildCallExpr(e)

	case *ast.IndexExpr:
		return b.buildIndexExpr(e)

	case *ast.PropertyExpr:
		return b.buildPropertyExpr(e)

	case *ast.ArrayLiteral:
		return b.buildArrayLit(e)

	case *ast.ObjectLiteral:
		return b.buildObjectLit(e)

	case *ast.FunctionExpr:
		return b.buildFuncExpr(e)

	case *ast.ArrowFunctionExpr:
		return b.buildArrowFuncExpr(e)

	case *ast.NewExpr:
		return b.buildNewExpr(e)

	case *ast.ThisExpr:
		return b.buildThisExpr(e)

	case *ast.SuperExpr:
		return b.buildSuperExpr(e)

	case *ast.AssignExpr:
		return b.buildAssignExpr(e)

	case *ast.CompoundAssignExpr:
		return b.buildCompoundAssignExpr(e)

	case *ast.UpdateExpr:
		return b.buildUpdateExpr(e)
	}

	return &NullLit{ExprType: types.AnyType}
}

func (b *Builder) buildIdent(ident *ast.Identifier) Expr {
	// Check for type narrowing
	if narrowed, ok := b.narrowing[ident.Name]; ok {
		return &Ident{Name: ident.Name, ExprType: narrowed}
	}

	if typ, found := b.scope.lookup(ident.Name); found {
		return &Ident{Name: ident.Name, ExprType: typ}
	}

	b.error(ident.Token.Line, ident.Token.Column, "undefined variable: %s", ident.Name)
	return &Ident{Name: ident.Name, ExprType: types.AnyType}
}

func (b *Builder) buildBinaryExpr(expr *ast.BinaryExpr) Expr {
	left := b.buildExpr(expr.Left)
	right := b.buildExpr(expr.Right)
	op := expr.Token.Literal

	var resultType types.Type

	switch expr.Op {
	case token.PLUS:
		if left.Type().Equals(types.StringType) && right.Type().Equals(types.StringType) {
			resultType = types.StringType
		} else {
			if !types.IsNumeric(left.Type()) || !types.IsNumeric(right.Type()) {
				b.error(expr.Token.Line, expr.Token.Column,
					"operator + requires number or string, got %s and %s",
					left.Type().String(), right.Type().String())
			}
			resultType = types.NumericResultType(left.Type(), right.Type())
		}

	case token.MINUS, token.STAR:
		if !types.IsNumeric(left.Type()) || !types.IsNumeric(right.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"operator %s requires numbers, got %s and %s",
				op, left.Type().String(), right.Type().String())
		}
		resultType = types.NumericResultType(left.Type(), right.Type())

	case token.SLASH:
		// Division always produces float (or any if either operand is any)
		if !types.IsNumeric(left.Type()) || !types.IsNumeric(right.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"operator / requires numbers, got %s and %s",
				left.Type().String(), right.Type().String())
		}
		if left.Type().Equals(types.AnyType) || right.Type().Equals(types.AnyType) {
			resultType = types.AnyType
		} else {
			resultType = types.FloatType
		}

	case token.PERCENT:
		// Modulo requires int operands
		if !types.IsNumeric(left.Type()) || !types.IsNumeric(right.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"operator %% requires numbers, got %s and %s",
				left.Type().String(), right.Type().String())
		}
		resultType = types.NumericResultType(left.Type(), right.Type())

	case token.LT, token.GT, token.LTE, token.GTE:
		if !types.IsNumeric(left.Type()) || !types.IsNumeric(right.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"comparison requires numbers, got %s and %s",
				left.Type().String(), right.Type().String())
		}
		resultType = types.BooleanType

	case token.EQ, token.NEQ:
		resultType = types.BooleanType

	case token.AND, token.OR:
		if !left.Type().Equals(types.BooleanType) || !right.Type().Equals(types.BooleanType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"logical operator requires booleans, got %s and %s",
				left.Type().String(), right.Type().String())
		}
		resultType = types.BooleanType

	case token.NULLISH_COALESCE:
		if nullable, ok := left.Type().(*types.Nullable); ok {
			resultType = types.LeastUpperBound(nullable.Inner, right.Type())
		} else {
			resultType = left.Type()
		}

	default:
		resultType = types.AnyType
	}

	return &BinaryExpr{
		Left:     left,
		Op:       op,
		Right:    right,
		ExprType: resultType,
	}
}

func (b *Builder) buildUnaryExpr(expr *ast.UnaryExpr) Expr {
	operand := b.buildExpr(expr.Operand)
	op := expr.Token.Literal

	var resultType types.Type

	switch expr.Op {
	case token.MINUS:
		if !types.IsNumeric(operand.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"unary - requires number, got %s", operand.Type().String())
		}
		resultType = operand.Type() // Preserve int/float

	case token.NOT:
		if !operand.Type().Equals(types.BooleanType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"unary ! requires boolean, got %s", operand.Type().String())
		}
		resultType = types.BooleanType

	default:
		resultType = types.AnyType
	}

	return &UnaryExpr{
		Op:       op,
		Operand:  operand,
		ExprType: resultType,
	}
}

// isBuiltin checks if a name is a built-in function.
var builtins = map[string]types.Type{
	"println":  &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.VoidType},
	"print":    &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.VoidType},
	"len":      &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.IntType},
	"push":     &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "val", Type: types.AnyType}}, ReturnType: types.VoidType},
	"pop":      &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}}, ReturnType: types.AnyType},
	"typeof":   &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.StringType},
	"tostring": &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.StringType},
	"toint":    &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.IntType},
	"tofloat":  &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.FloatType},
	"sqrt":     &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"floor":    &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"ceil":     &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"abs":      &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
}

func (b *Builder) buildCallExpr(expr *ast.CallExpr) Expr {
	// Check for built-in function calls
	if ident, ok := expr.Function.(*ast.Identifier); ok {
		if builtinType, isBuiltin := builtins[ident.Name]; isBuiltin {
			fn := builtinType.(*types.Function)
			args := make([]Expr, len(expr.Arguments))
			for i, arg := range expr.Arguments {
				args[i] = b.buildExpr(arg)
			}

			// Special case: pop returns the element type of the array
			returnType := fn.ReturnType
			if ident.Name == "pop" && len(args) > 0 {
				argType := types.Unwrap(args[0].Type())
				if arrType, ok := argType.(*types.Array); ok {
					returnType = arrType.Element
				}
			}

			return &BuiltinCall{
				Name:     ident.Name,
				Args:     args,
				ExprType: returnType,
			}
		}
	}

	callee := b.buildExpr(expr.Function)
	calleeType := types.Unwrap(callee.Type())

	// Allow calling values of type 'any' - returns 'any'
	if prim, ok := calleeType.(*types.Primitive); ok && prim.Kind == types.KindAny {
		args := make([]Expr, len(expr.Arguments))
		for i, arg := range expr.Arguments {
			args[i] = b.buildExpr(arg)
		}
		return &CallExpr{
			Callee:   callee,
			Args:     args,
			Optional: expr.Optional,
			ExprType: types.AnyType,
		}
	}

	fn, ok := calleeType.(*types.Function)
	if !ok {
		b.error(expr.Token.Line, expr.Token.Column,
			"cannot call non-function type %s", calleeType.String())
		return &CallExpr{
			Callee:   callee,
			Args:     []Expr{},
			Optional: expr.Optional,
			ExprType: types.AnyType,
		}
	}

	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
		if i < len(fn.Params) {
			if !types.IsAssignableTo(args[i].Type(), fn.Params[i].Type) {
				b.error(expr.Token.Line, expr.Token.Column,
					"argument %d: cannot pass %s as %s",
					i+1, args[i].Type().String(), fn.Params[i].Type.String())
			}
		}
	}

	// Skip argument count check for generic function types (all any params)
	isGenericFunc := len(fn.Params) > 0 && fn.Params[0].Type.Equals(types.AnyType) && fn.ReturnType.Equals(types.AnyType)
	if !isGenericFunc && len(args) != len(fn.Params) {
		b.error(expr.Token.Line, expr.Token.Column,
			"expected %d arguments, got %d", len(fn.Params), len(args))
	}

	return &CallExpr{
		Callee:   callee,
		Args:     args,
		Optional: expr.Optional,
		ExprType: fn.ReturnType,
	}
}

func (b *Builder) buildIndexExpr(expr *ast.IndexExpr) Expr {
	object := b.buildExpr(expr.Object)
	index := b.buildExpr(expr.Index)
	objectType := types.Unwrap(object.Type())

	var resultType types.Type

	if arr, ok := objectType.(*types.Array); ok {
		if !index.Type().Equals(types.IntType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"array index must be int, got %s", index.Type().String())
		}
		resultType = arr.Element
	} else if objectType.Equals(types.StringType) {
		if !index.Type().Equals(types.IntType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"string index must be int, got %s", index.Type().String())
		}
		resultType = types.StringType
	} else {
		b.error(expr.Token.Line, expr.Token.Column,
			"cannot index type %s", objectType.String())
		resultType = types.AnyType
	}

	return &IndexExpr{
		Object:   object,
		Index:    index,
		Optional: expr.Optional,
		ExprType: resultType,
	}
}

func (b *Builder) buildPropertyExpr(expr *ast.PropertyExpr) Expr {
	object := b.buildExpr(expr.Object)
	objectType := object.Type()

	// Handle nullable types with optional chaining
	isNullable := false
	if nullable, ok := objectType.(*types.Nullable); ok {
		if expr.Optional {
			// With optional chaining, unwrap the nullable type
			objectType = nullable.Inner
			isNullable = true
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot access property on type %s", objectType.String())
			return &PropertyExpr{
				Object:   object,
				Property: expr.Property,
				Optional: expr.Optional,
				ExprType: types.AnyType,
			}
		}
	}

	objectType = types.Unwrap(objectType)

	var resultType types.Type

	switch obj := objectType.(type) {
	case *types.Object:
		prop := obj.GetProperty(expr.Property)
		if prop == nil {
			b.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on %s", expr.Property, obj.String())
			resultType = types.AnyType
		} else {
			resultType = prop.Type
		}

	case *types.Class:
		if field := obj.GetField(expr.Property); field != nil {
			resultType = field.Type
		} else if method := obj.GetMethod(expr.Property); method != nil {
			params := make([]*types.Param, len(method.Params))
			copy(params, method.Params)
			resultType = &types.Function{
				Params:     params,
				ReturnType: method.ReturnType,
			}
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on class %s", expr.Property, obj.Name)
			resultType = types.AnyType
		}

	default:
		b.error(expr.Token.Line, expr.Token.Column,
			"cannot access property on type %s", objectType.String())
		resultType = types.AnyType
	}

	// If original was nullable and using optional chaining, result is also nullable
	if isNullable && expr.Optional {
		resultType = &types.Nullable{Inner: resultType}
	}

	return &PropertyExpr{
		Object:   object,
		Property: expr.Property,
		Optional: expr.Optional,
		ExprType: resultType,
	}
}

func (b *Builder) buildArrayLit(expr *ast.ArrayLiteral) Expr {
	elements := make([]Expr, len(expr.Elements))
	var elemType types.Type = types.AnyType

	if len(expr.Elements) > 0 {
		elements[0] = b.buildExpr(expr.Elements[0])
		elemType = elements[0].Type()

		for i := 1; i < len(expr.Elements); i++ {
			elements[i] = b.buildExpr(expr.Elements[i])
			elemType = types.LeastUpperBound(elemType, elements[i].Type())
		}
	}

	return &ArrayLit{
		Elements: elements,
		ExprType: &types.Array{Element: elemType},
	}
}

func (b *Builder) buildObjectLit(expr *ast.ObjectLiteral) Expr {
	props := make([]*PropertyInit, len(expr.Properties))
	typeProps := make(map[string]*types.Property)

	for i, p := range expr.Properties {
		value := b.buildExpr(p.Value)
		props[i] = &PropertyInit{
			Key:   p.Key,
			Value: value,
		}
		typeProps[p.Key] = &types.Property{
			Name: p.Key,
			Type: value.Type(),
		}
	}

	return &ObjectLit{
		Properties: props,
		ExprType:   &types.Object{Properties: typeProps},
	}
}

func (b *Builder) buildFuncExpr(expr *ast.FunctionExpr) Expr {
	params := make([]*Param, len(expr.Params))
	typeParams := make([]*types.Param, len(expr.Params))

	for i, p := range expr.Params {
		paramType := b.resolveType(p.ParamType)
		params[i] = &Param{Name: p.Name, Type: paramType}
		typeParams[i] = &types.Param{Name: p.Name, Type: paramType}
	}

	returnType := b.resolveType(expr.ReturnType)
	funcType := &types.Function{
		Params:     typeParams,
		ReturnType: returnType,
	}

	savedFunc := b.currentFunc
	b.currentFunc = funcType

	b.pushScope()
	for _, p := range params {
		b.scope.define(p.Name, p.Type)
	}

	body := b.buildBlock(expr.Body)

	b.popScope()
	b.currentFunc = savedFunc

	return &FuncExpr{
		Params:   params,
		Body:     body,
		ExprType: funcType,
	}
}

func (b *Builder) buildArrowFuncExpr(expr *ast.ArrowFunctionExpr) Expr {
	params := make([]*Param, len(expr.Params))
	typeParams := make([]*types.Param, len(expr.Params))

	for i, p := range expr.Params {
		paramType := b.resolveType(p.ParamType)
		params[i] = &Param{Name: p.Name, Type: paramType}
		typeParams[i] = &types.Param{Name: p.Name, Type: paramType}
	}

	returnType := b.resolveType(expr.ReturnType)
	funcType := &types.Function{
		Params:     typeParams,
		ReturnType: returnType,
	}

	savedFunc := b.currentFunc
	b.currentFunc = funcType

	b.pushScope()
	for _, p := range params {
		b.scope.define(p.Name, p.Type)
	}

	var body *BlockStmt
	var bodyExpr Expr

	if expr.Body != nil {
		body = b.buildBlock(expr.Body)
	} else if expr.Expression != nil {
		bodyExpr = b.buildExpr(expr.Expression)
		if !types.IsAssignableTo(bodyExpr.Type(), returnType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot return %s, expected %s", bodyExpr.Type().String(), returnType.String())
		}
	}

	b.popScope()
	b.currentFunc = savedFunc

	return &FuncExpr{
		Params:   params,
		Body:     body,
		BodyExpr: bodyExpr,
		ExprType: funcType,
	}
}

func (b *Builder) buildNewExpr(expr *ast.NewExpr) Expr {
	class, ok := b.classes[expr.ClassName]
	if !ok {
		b.error(expr.Token.Line, expr.Token.Column,
			"unknown class: %s", expr.ClassName)
		return &NewExpr{
			ClassName: expr.ClassName,
			Args:      []Expr{},
			ExprType:  types.AnyType,
		}
	}

	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	if class.Constructor != nil {
		if len(args) != len(class.Constructor.Params) {
			b.error(expr.Token.Line, expr.Token.Column,
				"expected %d constructor arguments, got %d",
				len(class.Constructor.Params), len(args))
		} else {
			for i, arg := range args {
				if !types.IsAssignableTo(arg.Type(), class.Constructor.Params[i].Type) {
					b.error(expr.Token.Line, expr.Token.Column,
						"constructor argument %d: cannot pass %s as %s",
						i+1, arg.Type().String(), class.Constructor.Params[i].Type.String())
				}
			}
		}
	} else if len(args) > 0 {
		b.error(expr.Token.Line, expr.Token.Column,
			"class %s has no constructor but was called with arguments", expr.ClassName)
	}

	return &NewExpr{
		ClassName: expr.ClassName,
		Args:      args,
		ExprType:  class,
	}
}

func (b *Builder) buildThisExpr(expr *ast.ThisExpr) Expr {
	if b.currentClass == nil {
		b.error(expr.Token.Line, expr.Token.Column, "'this' outside of class")
		return &ThisExpr{ExprType: types.AnyType}
	}
	return &ThisExpr{ExprType: b.currentClass}
}

func (b *Builder) buildSuperExpr(expr *ast.SuperExpr) Expr {
	if b.currentClass == nil {
		b.error(expr.Token.Line, expr.Token.Column, "'super' outside of class")
		return &SuperExpr{Args: []Expr{}, ExprType: types.VoidType}
	}
	if b.currentClass.Super == nil {
		b.error(expr.Token.Line, expr.Token.Column,
			"class %s has no superclass", b.currentClass.Name)
		return &SuperExpr{Args: []Expr{}, ExprType: types.VoidType}
	}

	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	super := b.currentClass.Super
	if super.Constructor != nil {
		if len(args) != len(super.Constructor.Params) {
			b.error(expr.Token.Line, expr.Token.Column,
				"expected %d super arguments, got %d",
				len(super.Constructor.Params), len(args))
		} else {
			for i, arg := range args {
				if !types.IsAssignableTo(arg.Type(), super.Constructor.Params[i].Type) {
					b.error(expr.Token.Line, expr.Token.Column,
						"super argument %d: cannot pass %s as %s",
						i+1, arg.Type().String(), super.Constructor.Params[i].Type.String())
				}
			}
		}
	}

	return &SuperExpr{Args: args, ExprType: types.VoidType}
}

func (b *Builder) buildAssignExpr(expr *ast.AssignExpr) Expr {
	value := b.buildExpr(expr.Value)
	target := b.buildExpr(expr.Target)

	// Check for const reassignment
	if ident, ok := expr.Target.(*ast.Identifier); ok {
		if b.scope.isConst(ident.Name) {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", ident.Name)
		}
	}

	if !types.IsAssignableTo(value.Type(), target.Type()) {
		b.error(expr.Token.Line, expr.Token.Column,
			"cannot assign %s to %s", value.Type().String(), target.Type().String())
	}

	return &AssignExpr{
		Target:   target,
		Value:    value,
		ExprType: target.Type(),
	}
}

func (b *Builder) buildCompoundAssignExpr(expr *ast.CompoundAssignExpr) Expr {
	target := b.buildExpr(expr.Target)
	value := b.buildExpr(expr.Value)
	op := expr.Token.Literal

	// Check for const reassignment
	if ident, ok := expr.Target.(*ast.Identifier); ok {
		if b.scope.isConst(ident.Name) {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", ident.Name)
		}
	}

	var resultType types.Type

	switch expr.Op {
	case token.PLUS_ASSIGN:
		if target.Type().Equals(types.StringType) && value.Type().Equals(types.StringType) {
			resultType = types.StringType
		} else {
			if !types.IsNumeric(target.Type()) || !types.IsNumeric(value.Type()) {
				b.error(expr.Token.Line, expr.Token.Column,
					"operator += requires numbers or strings, got %s and %s",
					target.Type().String(), value.Type().String())
			}
			resultType = target.Type() // Preserve original type
		}
	default:
		if !types.IsNumeric(target.Type()) || !types.IsNumeric(value.Type()) {
			b.error(expr.Token.Line, expr.Token.Column,
				"operator %s requires numbers, got %s and %s",
				op, target.Type().String(), value.Type().String())
		}
		resultType = target.Type() // Preserve original type
	}

	return &CompoundAssignExpr{
		Target:   target,
		Op:       op,
		Value:    value,
		ExprType: resultType,
	}
}

func (b *Builder) buildUpdateExpr(expr *ast.UpdateExpr) Expr {
	operand := b.buildExpr(expr.Operand)
	op := expr.Token.Literal

	// Check for const reassignment
	if ident, ok := expr.Operand.(*ast.Identifier); ok {
		if b.scope.isConst(ident.Name) {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", ident.Name)
		}
	}

	if !types.IsNumeric(operand.Type()) {
		b.error(expr.Token.Line, expr.Token.Column,
			"operator %s requires number, got %s", op, operand.Type().String())
	}

	return &UpdateExpr{
		Op:       op,
		Operand:  operand,
		Prefix:   expr.Prefix,
		ExprType: operand.Type(), // Preserve type
	}
}
