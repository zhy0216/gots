// Package types implements the type checker for goTS.
package types

import (
	"fmt"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/token"
)

// isNumericType checks if a type is int or float.
func isNumericType(t Type) bool {
	if p, ok := t.(*Primitive); ok {
		return p.Kind == KindInt || p.Kind == KindFloat
	}
	return false
}

// isNumericOrAny checks if a type is int, float, or any.
// This allows arithmetic operations on dynamic types.
func isNumericOrAny(t Type) bool {
	if p, ok := t.(*Primitive); ok {
		return p.Kind == KindInt || p.Kind == KindFloat || p.Kind == KindAny
	}
	return false
}

// numericResultType returns the result type for numeric operations.
// If either operand is any, result is any.
// If either operand is float, result is float. Otherwise int.
func numericResultType(left, right Type) Type {
	leftP, leftOk := left.(*Primitive)
	rightP, rightOk := right.(*Primitive)
	if leftOk && rightOk {
		// If either is any, result is any
		if leftP.Kind == KindAny || rightP.Kind == KindAny {
			return AnyType
		}
		if leftP.Kind == KindFloat || rightP.Kind == KindFloat {
			return FloatType
		}
		return IntType
	}
	return FloatType // default to float for safety
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

// Scope represents a lexical scope with variable bindings.
type Scope struct {
	parent    *Scope
	bindings  map[string]Type
	constVars map[string]bool // Track which variables are const
}

// NewScope creates a new scope with an optional parent.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		bindings:  make(map[string]Type),
		constVars: make(map[string]bool),
	}
}

// Define adds a binding to the current scope.
func (s *Scope) Define(name string, typ Type) {
	s.bindings[name] = typ
}

// DefineConst adds a const binding to the current scope.
func (s *Scope) DefineConst(name string, typ Type) {
	s.bindings[name] = typ
	s.constVars[name] = true
}

// IsConst checks if a variable is const, searching up the scope chain.
func (s *Scope) IsConst(name string) bool {
	if isConst, ok := s.constVars[name]; ok {
		return isConst
	}
	if s.parent != nil {
		return s.parent.IsConst(name)
	}
	return false
}

// Lookup finds a binding, searching up the scope chain.
func (s *Scope) Lookup(name string) (Type, bool) {
	if typ, ok := s.bindings[name]; ok {
		return typ, true
	}
	if s.parent != nil {
		return s.parent.Lookup(name)
	}
	return nil, false
}

// TypeNarrowing tracks narrowed types within a scope.
type TypeNarrowing struct {
	narrowedTypes map[string]Type
}

// NewTypeNarrowing creates a new type narrowing context.
func NewTypeNarrowing() *TypeNarrowing {
	return &TypeNarrowing{
		narrowedTypes: make(map[string]Type),
	}
}

// Narrow sets a narrowed type for a variable.
func (tn *TypeNarrowing) Narrow(name string, typ Type) {
	tn.narrowedTypes[name] = typ
}

// Get returns the narrowed type, or nil if not narrowed.
func (tn *TypeNarrowing) Get(name string) Type {
	return tn.narrowedTypes[name]
}

// Checker performs type checking on a goTS AST.
type Checker struct {
	errors       []*Error
	scope        *Scope
	typeAliases  map[string]Type
	classes      map[string]*Class
	currentFunc  *Function  // Current function for return type checking
	currentClass *Class     // Current class for 'this' type
	narrowing    *TypeNarrowing
	loopDepth    int        // Track if we're inside a loop for break/continue
}

// NewChecker creates a new type checker.
func NewChecker() *Checker {
	return &Checker{
		errors:      []*Error{},
		scope:       NewScope(nil),
		typeAliases: make(map[string]Type),
		classes:     make(map[string]*Class),
		narrowing:   NewTypeNarrowing(),
	}
}

// Errors returns the list of type checking errors.
func (c *Checker) Errors() []*Error {
	return c.errors
}

// HasErrors returns true if there are type checking errors.
func (c *Checker) HasErrors() bool {
	return len(c.errors) > 0
}

func (c *Checker) error(line, col int, format string, args ...interface{}) {
	c.errors = append(c.errors, &Error{
		Line:    line,
		Column:  col,
		Message: fmt.Sprintf(format, args...),
	})
}

func (c *Checker) pushScope() {
	c.scope = NewScope(c.scope)
}

func (c *Checker) popScope() {
	c.scope = c.scope.parent
}

// Check performs type checking on the program.
func (c *Checker) Check(program *ast.Program) {
	// First pass: collect type aliases and class declarations
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			c.collectTypeAlias(s)
		case *ast.ClassDecl:
			c.collectClass(s)
		}
	}

	// Second pass: resolve type aliases and class details
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			c.resolveTypeAlias(s)
		case *ast.ClassDecl:
			c.resolveClass(s)
		}
	}

	// Third pass: type check all statements
	for _, stmt := range program.Statements {
		c.checkStatement(stmt)
	}
}

// ----------------------------------------------------------------------------
// Type Resolution
// ----------------------------------------------------------------------------

// resolveType converts an AST type to a semantic type.
func (c *Checker) resolveType(astType ast.Type) Type {
	if astType == nil {
		return VoidType
	}

	switch t := astType.(type) {
	case *ast.PrimitiveType:
		switch t.Kind {
		case ast.TypeInt:
			return IntType
		case ast.TypeFloat:
			return FloatType
		case ast.TypeString:
			return StringType
		case ast.TypeBoolean:
			return BooleanType
		case ast.TypeVoid:
			return VoidType
		case ast.TypeNull:
			return NullType
		}

	case *ast.ArrayType:
		elem := c.resolveType(t.ElementType)
		return &Array{Element: elem}

	case *ast.ObjectType:
		props := make(map[string]*Property)
		for _, p := range t.Properties {
			props[p.Name] = &Property{
				Name: p.Name,
				Type: c.resolveType(p.PropType),
			}
		}
		return &Object{Properties: props}

	case *ast.FunctionType:
		params := make([]*Param, len(t.ParamTypes))
		for i, pt := range t.ParamTypes {
			params[i] = &Param{
				Name: fmt.Sprintf("arg%d", i),
				Type: c.resolveType(pt),
			}
		}
		return &Function{
			Params:     params,
			ReturnType: c.resolveType(t.ReturnType),
		}

	case *ast.NullableType:
		inner := c.resolveType(t.Inner)
		return &Nullable{Inner: inner}

	case *ast.NamedType:
		// Check type aliases first
		if alias, ok := c.typeAliases[t.Name]; ok {
			return alias
		}
		// Check classes
		if class, ok := c.classes[t.Name]; ok {
			return class
		}
		c.error(0, 0, "unknown type: %s", t.Name)
		return AnyType
	}

	return AnyType
}

// ----------------------------------------------------------------------------
// Collection Pass
// ----------------------------------------------------------------------------

func (c *Checker) collectTypeAlias(decl *ast.TypeAliasDecl) {
	// Create a placeholder that will be resolved later
	c.typeAliases[decl.Name] = &Alias{Name: decl.Name}
}

func (c *Checker) resolveTypeAlias(decl *ast.TypeAliasDecl) {
	alias := c.typeAliases[decl.Name].(*Alias)
	alias.Resolved = c.resolveType(decl.AliasType)
}

func (c *Checker) collectClass(decl *ast.ClassDecl) {
	c.classes[decl.Name] = &Class{
		Name:    decl.Name,
		Fields:  make(map[string]*Field),
		Methods: make(map[string]*Method),
	}
}

func (c *Checker) resolveClass(decl *ast.ClassDecl) {
	class := c.classes[decl.Name]

	// Resolve superclass
	if decl.SuperClass != "" {
		if super, ok := c.classes[decl.SuperClass]; ok {
			class.Super = super
		} else {
			c.error(decl.Token.Line, decl.Token.Column, "unknown superclass: %s", decl.SuperClass)
		}
	}

	// Resolve fields
	for _, f := range decl.Fields {
		class.Fields[f.Name] = &Field{
			Name: f.Name,
			Type: c.resolveType(f.FieldType),
		}
	}

	// Resolve constructor
	if decl.Constructor != nil {
		params := make([]*Param, len(decl.Constructor.Params))
		for i, p := range decl.Constructor.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: c.resolveType(p.ParamType),
			}
		}
		class.Constructor = &Function{
			Params:     params,
			ReturnType: VoidType,
		}
	}

	// Resolve methods
	for _, m := range decl.Methods {
		params := make([]*Param, len(m.Params))
		for i, p := range m.Params {
			params[i] = &Param{
				Name: p.Name,
				Type: c.resolveType(p.ParamType),
			}
		}
		class.Methods[m.Name] = &Method{
			Name:       m.Name,
			Params:     params,
			ReturnType: c.resolveType(m.ReturnType),
		}
	}
}

// ----------------------------------------------------------------------------
// Statement Type Checking
// ----------------------------------------------------------------------------

func (c *Checker) checkStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExprStmt:
		c.checkExpr(s.Expr)

	case *ast.VarDecl:
		c.checkVarDecl(s)

	case *ast.Block:
		c.checkBlock(s)

	case *ast.IfStmt:
		c.checkIfStmt(s)

	case *ast.WhileStmt:
		c.checkWhileStmt(s)

	case *ast.ForStmt:
		c.checkForStmt(s)

	case *ast.ForOfStmt:
		c.checkForOfStmt(s)

	case *ast.SwitchStmt:
		c.checkSwitchStmt(s)

	case *ast.ReturnStmt:
		c.checkReturnStmt(s)

	case *ast.BreakStmt:
		c.checkBreakStmt(s)

	case *ast.ContinueStmt:
		c.checkContinueStmt(s)

	case *ast.TryStmt:
		c.checkTryStmt(s)

	case *ast.ThrowStmt:
		c.checkThrowStmt(s)

	case *ast.FuncDecl:
		c.checkFuncDecl(s)

	case *ast.ClassDecl:
		c.checkClassDecl(s)

	case *ast.TypeAliasDecl:
		// Already processed in collection pass
	}
}

func (c *Checker) checkVarDecl(decl *ast.VarDecl) {
	var declaredType Type

	if decl.VarType != nil {
		// Explicit type annotation
		declaredType = c.resolveType(decl.VarType)
		if decl.Value != nil {
			initType := c.checkExpr(decl.Value)
			if !IsAssignableTo(initType, declaredType) {
				c.error(decl.Token.Line, decl.Token.Column,
					"cannot assign %s to %s", initType.String(), declaredType.String())
			}
		}
	} else {
		// Type inference
		if decl.Value == nil {
			c.error(decl.Token.Line, decl.Token.Column,
				"variable declaration requires type annotation or initializer")
			return
		}
		declaredType = c.inferType(decl.Value)
	}

	if decl.IsConst {
		c.scope.DefineConst(decl.Name, declaredType)
	} else {
		c.scope.Define(decl.Name, declaredType)
	}
}

// inferType infers the type of an expression for type inference.
func (c *Checker) inferType(expr ast.Expression) Type {
	return c.checkExpr(expr)
}

func (c *Checker) checkBlock(block *ast.Block) {
	c.pushScope()
	for _, stmt := range block.Statements {
		c.checkStatement(stmt)
	}
	c.popScope()
}

func (c *Checker) checkIfStmt(stmt *ast.IfStmt) {
	condType := c.checkExpr(stmt.Condition)

	// Condition must be boolean
	if !condType.Equals(BooleanType) {
		c.error(stmt.Token.Line, stmt.Token.Column,
			"condition must be boolean, got %s", condType.String())
	}

	// Check for null narrowing
	savedNarrowing := c.narrowing
	c.narrowing = NewTypeNarrowing()

	// Check if condition is a null check (e.g., x != null)
	if binary, ok := stmt.Condition.(*ast.BinaryExpr); ok {
		if binary.Op == token.NEQ {
			if ident, ok := binary.Left.(*ast.Identifier); ok {
				if _, ok := binary.Right.(*ast.NullLiteral); ok {
					// x != null - narrow x to non-nullable in then branch
					if varType, found := c.scope.Lookup(ident.Name); found {
						if nullable, ok := varType.(*Nullable); ok {
							c.narrowing.Narrow(ident.Name, nullable.Inner)
						}
					}
				}
			}
		}
	}

	c.checkBlock(stmt.Consequence)

	c.narrowing = savedNarrowing

	if stmt.Alternative != nil {
		c.checkStatement(stmt.Alternative)
	}
}

func (c *Checker) checkWhileStmt(stmt *ast.WhileStmt) {
	condType := c.checkExpr(stmt.Condition)

	if !condType.Equals(BooleanType) {
		c.error(stmt.Token.Line, stmt.Token.Column,
			"condition must be boolean, got %s", condType.String())
	}

	c.loopDepth++
	c.checkBlock(stmt.Body)
	c.loopDepth--
}

func (c *Checker) checkForStmt(stmt *ast.ForStmt) {
	c.pushScope()

	if stmt.Init != nil {
		c.checkVarDecl(stmt.Init)
	}

	if stmt.Condition != nil {
		condType := c.checkExpr(stmt.Condition)
		if !condType.Equals(BooleanType) {
			c.error(stmt.Token.Line, stmt.Token.Column,
				"condition must be boolean, got %s", condType.String())
		}
	}

	if stmt.Update != nil {
		c.checkExpr(stmt.Update)
	}

	c.loopDepth++
	c.checkBlock(stmt.Body)
	c.loopDepth--

	c.popScope()
}

func (c *Checker) checkReturnStmt(stmt *ast.ReturnStmt) {
	if c.currentFunc == nil {
		c.error(stmt.Token.Line, stmt.Token.Column, "return outside function")
		return
	}

	if stmt.Value == nil {
		if !c.currentFunc.ReturnType.Equals(VoidType) {
			c.error(stmt.Token.Line, stmt.Token.Column,
				"missing return value, expected %s", c.currentFunc.ReturnType.String())
		}
		return
	}

	returnType := c.checkExpr(stmt.Value)
	if !IsAssignableTo(returnType, c.currentFunc.ReturnType) {
		c.error(stmt.Token.Line, stmt.Token.Column,
			"cannot return %s, expected %s", returnType.String(), c.currentFunc.ReturnType.String())
	}
}

func (c *Checker) checkBreakStmt(stmt *ast.BreakStmt) {
	if c.loopDepth == 0 {
		c.error(stmt.Token.Line, stmt.Token.Column, "break outside loop")
	}
}

func (c *Checker) checkContinueStmt(stmt *ast.ContinueStmt) {
	if c.loopDepth == 0 {
		c.error(stmt.Token.Line, stmt.Token.Column, "continue outside loop")
	}
}

func (c *Checker) checkTryStmt(stmt *ast.TryStmt) {
	// Check the try block
	c.pushScope()
	for _, s := range stmt.TryBlock.Statements {
		c.checkStatement(s)
	}
	c.popScope()

	// Check the catch block with catch parameter in scope
	c.pushScope()
	c.scope.Define(stmt.CatchParam, AnyType)
	for _, s := range stmt.CatchBlock.Statements {
		c.checkStatement(s)
	}
	c.popScope()
}

func (c *Checker) checkThrowStmt(stmt *ast.ThrowStmt) {
	// Throw can throw any expression
	c.checkExpr(stmt.Value)
}

func (c *Checker) checkFuncDecl(decl *ast.FuncDecl) {
	// Build function type
	params := make([]*Param, len(decl.Params))
	for i, p := range decl.Params {
		params[i] = &Param{
			Name: p.Name,
			Type: c.resolveType(p.ParamType),
		}
	}

	funcType := &Function{
		Params:     params,
		ReturnType: c.resolveType(decl.ReturnType),
	}

	// Define function in current scope
	c.scope.Define(decl.Name, funcType)

	// Check function body
	savedFunc := c.currentFunc
	c.currentFunc = funcType

	c.pushScope()
	for _, p := range params {
		c.scope.Define(p.Name, p.Type)
	}

	for _, stmt := range decl.Body.Statements {
		c.checkStatement(stmt)
	}

	c.popScope()
	c.currentFunc = savedFunc
}

func (c *Checker) checkClassDecl(decl *ast.ClassDecl) {
	class := c.classes[decl.Name]

	savedClass := c.currentClass
	c.currentClass = class

	// Check constructor body
	if decl.Constructor != nil {
		c.pushScope()

		// Add constructor parameters to scope
		for _, p := range decl.Constructor.Params {
			c.scope.Define(p.Name, c.resolveType(p.ParamType))
		}

		for _, stmt := range decl.Constructor.Body.Statements {
			c.checkStatement(stmt)
		}

		c.popScope()
	}

	// Check method bodies
	for _, m := range decl.Methods {
		method := class.Methods[m.Name]

		savedFunc := c.currentFunc
		c.currentFunc = &Function{
			Params:     method.Params,
			ReturnType: method.ReturnType,
		}

		c.pushScope()

		for _, p := range m.Params {
			c.scope.Define(p.Name, c.resolveType(p.ParamType))
		}

		for _, stmt := range m.Body.Statements {
			c.checkStatement(stmt)
		}

		c.popScope()
		c.currentFunc = savedFunc
	}

	c.currentClass = savedClass
}

// ----------------------------------------------------------------------------
// Expression Type Checking
// ----------------------------------------------------------------------------

func (c *Checker) checkExpr(expr ast.Expression) Type {
	if expr == nil {
		return VoidType
	}

	switch e := expr.(type) {
	case *ast.NumberLiteral:
		// Infer int for integers, float for decimals
		if e.Value == float64(int64(e.Value)) {
			return IntType
		}
		return FloatType

	case *ast.StringLiteral:
		return StringType

	case *ast.BoolLiteral:
		return BooleanType

	case *ast.NullLiteral:
		return NullType

	case *ast.Identifier:
		return c.checkIdentifier(e)

	case *ast.BinaryExpr:
		return c.checkBinaryExpr(e)

	case *ast.UnaryExpr:
		return c.checkUnaryExpr(e)

	case *ast.CallExpr:
		return c.checkCallExpr(e)

	case *ast.IndexExpr:
		return c.checkIndexExpr(e)

	case *ast.PropertyExpr:
		return c.checkPropertyExpr(e)

	case *ast.ArrayLiteral:
		return c.checkArrayLiteral(e)

	case *ast.ObjectLiteral:
		return c.checkObjectLiteral(e)

	case *ast.FunctionExpr:
		return c.checkFunctionExpr(e)

	case *ast.ArrowFunctionExpr:
		return c.checkArrowFunctionExpr(e)

	case *ast.NewExpr:
		return c.checkNewExpr(e)

	case *ast.ThisExpr:
		return c.checkThisExpr(e)

	case *ast.SuperExpr:
		return c.checkSuperExpr(e)

	case *ast.AssignExpr:
		return c.checkAssignExpr(e)

	case *ast.CompoundAssignExpr:
		return c.checkCompoundAssignExpr(e)

	case *ast.UpdateExpr:
		return c.checkUpdateExpr(e)
	}

	return AnyType
}

func (c *Checker) checkIdentifier(ident *ast.Identifier) Type {
	// Check for type narrowing first
	if narrowed := c.narrowing.Get(ident.Name); narrowed != nil {
		return narrowed
	}

	if typ, found := c.scope.Lookup(ident.Name); found {
		return typ
	}

	c.error(ident.Token.Line, ident.Token.Column, "undefined variable: %s", ident.Name)
	return AnyType
}

func (c *Checker) checkBinaryExpr(expr *ast.BinaryExpr) Type {
	left := c.checkExpr(expr.Left)
	right := c.checkExpr(expr.Right)

	switch expr.Op {
	// Arithmetic operators
	case token.PLUS:
		// Allow string concatenation
		if left.Equals(StringType) && right.Equals(StringType) {
			return StringType
		}
		// Allow any type in arithmetic for dynamic typing support
		if left.Equals(AnyType) && right.Equals(AnyType) {
			return AnyType
		}
		if !isNumericOrAny(left) || !isNumericOrAny(right) {
			c.error(expr.Token.Line, expr.Token.Column,
				"operator + requires number or string, got %s and %s", left.String(), right.String())
		}
		return numericResultType(left, right)

	case token.MINUS, token.STAR, token.PERCENT:
		if !isNumericOrAny(left) || !isNumericOrAny(right) {
			c.error(expr.Token.Line, expr.Token.Column,
				"operator %s requires numbers, got %s and %s", expr.Token.Literal, left.String(), right.String())
		}
		return numericResultType(left, right)

	case token.SLASH:
		// Division always returns float (or any if either operand is any)
		if !isNumericOrAny(left) || !isNumericOrAny(right) {
			c.error(expr.Token.Line, expr.Token.Column,
				"operator / requires numbers, got %s and %s", left.String(), right.String())
		}
		if left.Equals(AnyType) || right.Equals(AnyType) {
			return AnyType
		}
		return FloatType

	// Comparison operators
	case token.LT, token.GT, token.LTE, token.GTE:
		if !isNumericOrAny(left) || !isNumericOrAny(right) {
			c.error(expr.Token.Line, expr.Token.Column,
				"comparison requires numbers, got %s and %s", left.String(), right.String())
		}
		return BooleanType

	// Equality operators
	case token.EQ, token.NEQ:
		// Allow comparison between compatible types
		if !IsAssignableTo(left, right) && !IsAssignableTo(right, left) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot compare %s and %s", left.String(), right.String())
		}
		return BooleanType

	// Logical operators
	case token.AND, token.OR:
		if !left.Equals(BooleanType) || !right.Equals(BooleanType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"logical operator requires booleans, got %s and %s", left.String(), right.String())
		}
		return BooleanType

	// Nullish coalescing
	case token.NULLISH_COALESCE:
		// left ?? right - returns right if left is null/undefined
		// The type is the union of the non-null type of left and right
		// For simplicity, return the type of right (or left if it's not nullable)
		if nullable, ok := left.(*Nullable); ok {
			// If left is nullable, result could be inner type or right type
			return LeastUpperBound(nullable.Inner, right)
		}
		return left
	}

	return AnyType
}

func (c *Checker) checkUnaryExpr(expr *ast.UnaryExpr) Type {
	operand := c.checkExpr(expr.Operand)

	switch expr.Op {
	case token.MINUS:
		if !isNumericType(operand) {
			c.error(expr.Token.Line, expr.Token.Column,
				"unary - requires number, got %s", operand.String())
		}
		return operand // preserve int/float type

	case token.NOT:
		if !operand.Equals(BooleanType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"unary ! requires boolean, got %s", operand.String())
		}
		return BooleanType
	}

	return AnyType
}

func (c *Checker) checkCallExpr(expr *ast.CallExpr) Type {
	calleeType := c.checkExpr(expr.Function)
	calleeType = Unwrap(calleeType)

	fn, ok := calleeType.(*Function)
	if !ok {
		c.error(expr.Token.Line, expr.Token.Column,
			"cannot call non-function type %s", calleeType.String())
		return AnyType
	}

	if len(expr.Arguments) != len(fn.Params) {
		c.error(expr.Token.Line, expr.Token.Column,
			"expected %d arguments, got %d", len(fn.Params), len(expr.Arguments))
		return fn.ReturnType
	}

	for i, arg := range expr.Arguments {
		argType := c.checkExpr(arg)
		if !IsAssignableTo(argType, fn.Params[i].Type) {
			c.error(expr.Token.Line, expr.Token.Column,
				"argument %d: cannot pass %s as %s", i+1, argType.String(), fn.Params[i].Type.String())
		}
	}

	return fn.ReturnType
}

func (c *Checker) checkIndexExpr(expr *ast.IndexExpr) Type {
	objectType := c.checkExpr(expr.Object)
	indexType := c.checkExpr(expr.Index)
	objectType = Unwrap(objectType)

	if arr, ok := objectType.(*Array); ok {
		if !indexType.Equals(IntType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"array index must be int, got %s", indexType.String())
		}
		return arr.Element
	}

	// String indexing
	if objectType.Equals(StringType) {
		if !indexType.Equals(IntType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"string index must be int, got %s", indexType.String())
		}
		return StringType
	}

	c.error(expr.Token.Line, expr.Token.Column,
		"cannot index type %s", objectType.String())
	return AnyType
}

func (c *Checker) checkPropertyExpr(expr *ast.PropertyExpr) Type {
	objectType := c.checkExpr(expr.Object)
	objectType = Unwrap(objectType)

	switch obj := objectType.(type) {
	case *Object:
		prop := obj.GetProperty(expr.Property)
		if prop == nil {
			c.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on %s", expr.Property, obj.String())
			return AnyType
		}
		return prop.Type

	case *Class:
		// Check for field
		if field := obj.GetField(expr.Property); field != nil {
			return field.Type
		}
		// Check for method
		if method := obj.GetMethod(expr.Property); method != nil {
			params := make([]*Param, len(method.Params))
			copy(params, method.Params)
			return &Function{
				Params:     params,
				ReturnType: method.ReturnType,
			}
		}
		c.error(expr.Token.Line, expr.Token.Column,
			"property %s does not exist on class %s", expr.Property, obj.Name)
		return AnyType

	case *Nullable:
		c.error(expr.Token.Line, expr.Token.Column,
			"cannot access property on potentially null value, use null check first")
		return AnyType
	}

	c.error(expr.Token.Line, expr.Token.Column,
		"cannot access property on type %s", objectType.String())
	return AnyType
}

func (c *Checker) checkArrayLiteral(expr *ast.ArrayLiteral) Type {
	if len(expr.Elements) == 0 {
		// Empty array - type will be inferred from context or default to any[]
		return &Array{Element: AnyType}
	}

	elemType := c.checkExpr(expr.Elements[0])
	for i := 1; i < len(expr.Elements); i++ {
		t := c.checkExpr(expr.Elements[i])
		elemType = LeastUpperBound(elemType, t)
	}

	return &Array{Element: elemType}
}

func (c *Checker) checkObjectLiteral(expr *ast.ObjectLiteral) Type {
	props := make(map[string]*Property)
	for _, p := range expr.Properties {
		propType := c.checkExpr(p.Value)
		props[p.Key] = &Property{
			Name: p.Key,
			Type: propType,
		}
	}
	return &Object{Properties: props}
}

func (c *Checker) checkFunctionExpr(expr *ast.FunctionExpr) Type {
	params := make([]*Param, len(expr.Params))
	for i, p := range expr.Params {
		params[i] = &Param{
			Name: p.Name,
			Type: c.resolveType(p.ParamType),
		}
	}

	funcType := &Function{
		Params:     params,
		ReturnType: c.resolveType(expr.ReturnType),
	}

	// Check body
	savedFunc := c.currentFunc
	c.currentFunc = funcType

	c.pushScope()
	for _, p := range params {
		c.scope.Define(p.Name, p.Type)
	}

	for _, stmt := range expr.Body.Statements {
		c.checkStatement(stmt)
	}

	c.popScope()
	c.currentFunc = savedFunc

	return funcType
}

func (c *Checker) checkNewExpr(expr *ast.NewExpr) Type {
	class, ok := c.classes[expr.ClassName]
	if !ok {
		c.error(expr.Token.Line, expr.Token.Column,
			"unknown class: %s", expr.ClassName)
		return AnyType
	}

	// Check constructor arguments
	if class.Constructor != nil {
		if len(expr.Arguments) != len(class.Constructor.Params) {
			c.error(expr.Token.Line, expr.Token.Column,
				"expected %d constructor arguments, got %d",
				len(class.Constructor.Params), len(expr.Arguments))
		} else {
			for i, arg := range expr.Arguments {
				argType := c.checkExpr(arg)
				if !IsAssignableTo(argType, class.Constructor.Params[i].Type) {
					c.error(expr.Token.Line, expr.Token.Column,
						"constructor argument %d: cannot pass %s as %s",
						i+1, argType.String(), class.Constructor.Params[i].Type.String())
				}
			}
		}
	} else if len(expr.Arguments) > 0 {
		c.error(expr.Token.Line, expr.Token.Column,
			"class %s has no constructor but was called with arguments", expr.ClassName)
	}

	return class
}

func (c *Checker) checkThisExpr(expr *ast.ThisExpr) Type {
	if c.currentClass == nil {
		c.error(expr.Token.Line, expr.Token.Column, "'this' outside of class")
		return AnyType
	}
	return c.currentClass
}

func (c *Checker) checkSuperExpr(expr *ast.SuperExpr) Type {
	if c.currentClass == nil {
		c.error(expr.Token.Line, expr.Token.Column, "'super' outside of class")
		return AnyType
	}
	if c.currentClass.Super == nil {
		c.error(expr.Token.Line, expr.Token.Column,
			"class %s has no superclass", c.currentClass.Name)
		return AnyType
	}

	// Check super constructor call arguments
	super := c.currentClass.Super
	if super.Constructor != nil {
		if len(expr.Arguments) != len(super.Constructor.Params) {
			c.error(expr.Token.Line, expr.Token.Column,
				"expected %d super arguments, got %d",
				len(super.Constructor.Params), len(expr.Arguments))
		} else {
			for i, arg := range expr.Arguments {
				argType := c.checkExpr(arg)
				if !IsAssignableTo(argType, super.Constructor.Params[i].Type) {
					c.error(expr.Token.Line, expr.Token.Column,
						"super argument %d: cannot pass %s as %s",
						i+1, argType.String(), super.Constructor.Params[i].Type.String())
				}
			}
		}
	}

	return VoidType
}

func (c *Checker) checkAssignExpr(expr *ast.AssignExpr) Type {
	valueType := c.checkExpr(expr.Value)

	switch target := expr.Target.(type) {
	case *ast.Identifier:
		// Check for const reassignment
		if c.scope.IsConst(target.Name) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", target.Name)
		}

		varType, found := c.scope.Lookup(target.Name)
		if !found {
			c.error(expr.Token.Line, expr.Token.Column,
				"undefined variable: %s", target.Name)
			return AnyType
		}
		if !IsAssignableTo(valueType, varType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot assign %s to %s", valueType.String(), varType.String())
		}
		return varType

	case *ast.IndexExpr:
		objectType := c.checkExpr(target.Object)
		objectType = Unwrap(objectType)

		if arr, ok := objectType.(*Array); ok {
			indexType := c.checkExpr(target.Index)
			if !indexType.Equals(IntType) {
				c.error(expr.Token.Line, expr.Token.Column,
					"array index must be int, got %s", indexType.String())
			}
			if !IsAssignableTo(valueType, arr.Element) {
				c.error(expr.Token.Line, expr.Token.Column,
					"cannot assign %s to array element of type %s", valueType.String(), arr.Element.String())
			}
			return arr.Element
		}
		c.error(expr.Token.Line, expr.Token.Column, "cannot index assign to non-array")
		return AnyType

	case *ast.PropertyExpr:
		objectType := c.checkExpr(target.Object)
		objectType = Unwrap(objectType)

		if obj, ok := objectType.(*Object); ok {
			prop := obj.GetProperty(target.Property)
			if prop == nil {
				c.error(expr.Token.Line, expr.Token.Column,
					"property %s does not exist", target.Property)
				return AnyType
			}
			if !IsAssignableTo(valueType, prop.Type) {
				c.error(expr.Token.Line, expr.Token.Column,
					"cannot assign %s to property of type %s", valueType.String(), prop.Type.String())
			}
			return prop.Type
		}

		if class, ok := objectType.(*Class); ok {
			field := class.GetField(target.Property)
			if field == nil {
				c.error(expr.Token.Line, expr.Token.Column,
					"field %s does not exist on class %s", target.Property, class.Name)
				return AnyType
			}
			if !IsAssignableTo(valueType, field.Type) {
				c.error(expr.Token.Line, expr.Token.Column,
					"cannot assign %s to field of type %s", valueType.String(), field.Type.String())
			}
			return field.Type
		}

		c.error(expr.Token.Line, expr.Token.Column, "cannot assign to property of non-object")
		return AnyType
	}

	c.error(expr.Token.Line, expr.Token.Column, "invalid assignment target")
	return AnyType
}

// checkForOfStmt checks a for-of statement.
func (c *Checker) checkForOfStmt(stmt *ast.ForOfStmt) {
	iterableType := c.checkExpr(stmt.Iterable)
	iterableType = Unwrap(iterableType)

	var elementType Type
	switch t := iterableType.(type) {
	case *Array:
		elementType = t.Element
	default:
		if iterableType.Equals(StringType) {
			elementType = StringType
		} else {
			c.error(stmt.Token.Line, stmt.Token.Column,
				"for-of requires array or string, got %s", iterableType.String())
			elementType = AnyType
		}
	}

	// Check variable type if explicitly declared
	if stmt.Variable.VarType != nil {
		declaredType := c.resolveType(stmt.Variable.VarType)
		if !IsAssignableTo(elementType, declaredType) {
			c.error(stmt.Token.Line, stmt.Token.Column,
				"cannot assign %s to %s", elementType.String(), declaredType.String())
		}
		elementType = declaredType
	}

	c.pushScope()
	c.scope.Define(stmt.Variable.Name, elementType)
	c.loopDepth++
	c.checkBlock(stmt.Body)
	c.loopDepth--
	c.popScope()
}

// checkSwitchStmt checks a switch statement.
func (c *Checker) checkSwitchStmt(stmt *ast.SwitchStmt) {
	discriminantType := c.checkExpr(stmt.Discriminant)

	for _, clause := range stmt.Cases {
		if clause.Test != nil {
			testType := c.checkExpr(clause.Test)
			if !IsAssignableTo(testType, discriminantType) && !IsAssignableTo(discriminantType, testType) {
				c.error(clause.Token.Line, clause.Token.Column,
					"case type %s is not comparable to switch type %s",
					testType.String(), discriminantType.String())
			}
		}

		// Check case body statements
		c.loopDepth++ // break is valid in switch
		for _, s := range clause.Consequent {
			c.checkStatement(s)
		}
		c.loopDepth--
	}
}

// checkCompoundAssignExpr checks compound assignment expressions (+=, -=, etc.).
func (c *Checker) checkCompoundAssignExpr(expr *ast.CompoundAssignExpr) Type {
	// Check for const reassignment
	if ident, ok := expr.Target.(*ast.Identifier); ok {
		if c.scope.IsConst(ident.Name) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", ident.Name)
		}
	}

	targetType := c.checkExpr(expr.Target)
	valueType := c.checkExpr(expr.Value)

	// For now, require both to be numbers for arithmetic compound assignment
	// String += string is also valid
	switch expr.Op {
	case token.PLUS_ASSIGN:
		if targetType.Equals(StringType) && valueType.Equals(StringType) {
			return StringType
		}
		if !isNumericType(targetType) || !isNumericType(valueType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"operator += requires numbers or strings, got %s and %s",
				targetType.String(), valueType.String())
		}
		return numericResultType(targetType, valueType)
	default:
		if !isNumericType(targetType) || !isNumericType(valueType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"operator %s requires numbers, got %s and %s",
				expr.Token.Literal, targetType.String(), valueType.String())
		}
		return numericResultType(targetType, valueType)
	}
}

// checkUpdateExpr checks increment/decrement expressions.
func (c *Checker) checkUpdateExpr(expr *ast.UpdateExpr) Type {
	// Check for const reassignment
	if ident, ok := expr.Operand.(*ast.Identifier); ok {
		if c.scope.IsConst(ident.Name) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot assign to const variable '%s'", ident.Name)
		}
	}

	operandType := c.checkExpr(expr.Operand)

	if !isNumericType(operandType) {
		c.error(expr.Token.Line, expr.Token.Column,
			"operator %s requires number, got %s",
			expr.Token.Literal, operandType.String())
	}

	// Verify operand is assignable (identifier, property, or index)
	switch expr.Operand.(type) {
	case *ast.Identifier, *ast.PropertyExpr, *ast.IndexExpr:
		// OK
	default:
		c.error(expr.Token.Line, expr.Token.Column,
			"invalid operand for %s", expr.Token.Literal)
	}

	return operandType // preserve int/float type
}

// checkArrowFunctionExpr checks arrow function expressions.
func (c *Checker) checkArrowFunctionExpr(expr *ast.ArrowFunctionExpr) Type {
	params := make([]*Param, len(expr.Params))
	for i, p := range expr.Params {
		params[i] = &Param{
			Name: p.Name,
			Type: c.resolveType(p.ParamType),
		}
	}

	funcType := &Function{
		Params:     params,
		ReturnType: c.resolveType(expr.ReturnType),
	}

	// Check body
	savedFunc := c.currentFunc
	c.currentFunc = funcType

	c.pushScope()
	for _, p := range params {
		c.scope.Define(p.Name, p.Type)
	}

	if expr.Body != nil {
		// Block body
		for _, stmt := range expr.Body.Statements {
			c.checkStatement(stmt)
		}
	} else if expr.Expression != nil {
		// Expression body - check that expression type matches return type
		exprType := c.checkExpr(expr.Expression)
		if !IsAssignableTo(exprType, funcType.ReturnType) {
			c.error(expr.Token.Line, expr.Token.Column,
				"cannot return %s, expected %s", exprType.String(), funcType.ReturnType.String())
		}
	}

	c.popScope()
	c.currentFunc = savedFunc

	return funcType
}
