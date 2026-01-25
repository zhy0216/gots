package typed

import (
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/token"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Builder transforms an AST into a TypedAST while performing type checking.
type Builder struct {
	errors          []*Error
	scope           *Scope
	typeAliases     map[string]types.Type
	classes         map[string]*types.Class
	genericClasses  map[string]*types.GenericClass
	interfaces      map[string]*types.Interface
	enums           map[string]*types.Enum
	goImports       []*GoImportDecl
	moduleImports   []*ModuleImportDecl
	exports         []string
	currentFunc     *types.Function
	currentClass    *types.Class
	narrowing       map[string]types.Type
	loopDepth       int
	scopeDepth      int
	constVars       map[string]bool                 // Track which variables are const (scoped)
	typeParamScope  map[string]*types.TypeParameter // Type parameters currently in scope
	inAsyncFunc     bool                            // Track if we're inside an async function
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
		errors:         []*Error{},
		scope:          newScope(nil),
		typeAliases:    make(map[string]types.Type),
		classes:        make(map[string]*types.Class),
		genericClasses: make(map[string]*types.GenericClass),
		interfaces:     make(map[string]*types.Interface),
		enums:          make(map[string]*types.Enum),
		narrowing:      make(map[string]types.Type),
		constVars:      make(map[string]bool),
		typeParamScope: make(map[string]*types.TypeParameter),
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
	// Initialize imports and exports
	b.goImports = make([]*GoImportDecl, 0)
	b.moduleImports = make([]*ModuleImportDecl, 0)
	b.exports = make([]string, 0)

	// First pass: collect type aliases, classes, interfaces, and imports
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			b.collectTypeAlias(s)
		case *ast.ClassDecl:
			b.collectClass(s)
		case *ast.InterfaceDecl:
			b.collectInterface(s)
		case *ast.GoImportDecl:
			b.collectGoImport(s)
		case *ast.ModuleImportDecl:
			b.collectModuleImport(s)
		case *ast.ExportModifier:
			b.collectExport(s)
		case *ast.DefaultExport:
			// Collect types from default exports
			switch d := s.Decl.(type) {
			case *ast.TypeAliasDecl:
				b.collectTypeAlias(d)
			case *ast.ClassDecl:
				b.collectClass(d)
			case *ast.InterfaceDecl:
				b.collectInterface(d)
			}
		}
	}

	// Second pass: resolve types
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.TypeAliasDecl:
			b.resolveTypeAlias(s)
		case *ast.ClassDecl:
			b.resolveClass(s)
		case *ast.InterfaceDecl:
			b.resolveInterface(s)
		case *ast.ExportModifier:
			// Resolve types for exported declarations
			switch d := s.Decl.(type) {
			case *ast.TypeAliasDecl:
				b.resolveTypeAlias(d)
			case *ast.ClassDecl:
				b.resolveClass(d)
			case *ast.InterfaceDecl:
				b.resolveInterface(d)
			}
		case *ast.DefaultExport:
			// Resolve types for default exports
			switch d := s.Decl.(type) {
			case *ast.TypeAliasDecl:
				b.resolveTypeAlias(d)
			case *ast.ClassDecl:
				b.resolveClass(d)
			case *ast.InterfaceDecl:
				b.resolveInterface(d)
			}
		}
	}

	// Third pass: build typed AST
	result := &Program{
		GoImports:        b.goImports,
		ModuleImports:    b.moduleImports,
		DefaultImports:   make([]*DefaultImport, 0),
		NamespaceImports: make([]*NamespaceImport, 0),
		TypeAliases:      make([]*TypeAlias, 0),
		Enums:            make([]*EnumDecl, 0),
		Classes:          make([]*ClassDecl, 0),
		Interfaces:       make([]*InterfaceDecl, 0),
		Functions:        make([]*FuncDecl, 0),
		TopLevel:         make([]Stmt, 0),
		Exports:          b.exports,
		ReExports:        make([]*ReExportDecl, 0),
		DefaultExports:   make([]*DefaultExport, 0),
	}

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.GoImportDecl:
			// Already handled in collectGoImport, skip
			continue

		case *ast.ModuleImportDecl:
			// Already handled in collectModuleImport, skip
			continue

		case *ast.DefaultImport:
			result.DefaultImports = append(result.DefaultImports, &DefaultImport{
				Name: s.Name,
				Path: s.Path,
			})
			continue

		case *ast.NamespaceImport:
			result.NamespaceImports = append(result.NamespaceImports, &NamespaceImport{
				Alias: s.Alias,
				Path:  s.Path,
			})
			continue

		case *ast.ReExportDecl:
			result.ReExports = append(result.ReExports, &ReExportDecl{
				Names:      s.Names,
				Path:       s.Path,
				IsWildcard: s.IsWildcard,
			})
			continue

		case *ast.DefaultExport:
			// Build the default export
			typedDecl := b.buildDecl(s.Decl)
			result.DefaultExports = append(result.DefaultExports, &DefaultExport{
				Decl: typedDecl,
			})
			continue

		case *ast.ExportModifier:
			// Build the inner declaration
			b.buildExportedDecl(s.Decl, result)
			continue

		case *ast.TypeAliasDecl:
			alias := &TypeAlias{
				Name:     s.Name,
				Resolved: b.typeAliases[s.Name],
			}
			result.TypeAliases = append(result.TypeAliases, alias)

		case *ast.EnumDecl:
			enumDecl := b.buildEnumDecl(s)
			result.Enums = append(result.Enums, enumDecl)

		case *ast.ClassDecl:
			classDecl := b.buildClassDecl(s)
			result.Classes = append(result.Classes, classDecl)

		case *ast.InterfaceDecl:
			ifaceDecl := b.buildInterfaceDecl(s)
			result.Interfaces = append(result.Interfaces, ifaceDecl)

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

func (b *Builder) buildEnumDecl(decl *ast.EnumDecl) *EnumDecl {
	enumDecl := &EnumDecl{
		Name:    decl.Name,
		Members: make([]*EnumMember, len(decl.Members)),
	}

	// Create the types.Enum for type checking
	enumType := &types.Enum{
		Name:    decl.Name,
		Members: make([]*types.EnumMember, len(decl.Members)),
	}

	// Assign values to enum members
	// Start from 0, or use explicit values if provided
	nextValue := 0
	for i, member := range decl.Members {
		value := nextValue
		if member.Value != nil {
			// If there's an explicit value, evaluate it
			if numLit, ok := member.Value.(*ast.NumberLiteral); ok {
				value = int(numLit.Value)
			}
		}
		enumDecl.Members[i] = &EnumMember{
			Name:  member.Name,
			Value: value,
		}
		enumType.Members[i] = &types.EnumMember{
			Name:  member.Name,
			Value: value,
		}
		nextValue = value + 1
	}

	// Register the enum type
	b.enums[decl.Name] = enumType

	return enumDecl
}

func (b *Builder) collectClass(decl *ast.ClassDecl) {
	// Check if this is a generic class
	if len(decl.TypeParams) > 0 {
		typeParams := make([]*types.TypeParameter, len(decl.TypeParams))
		for i, tp := range decl.TypeParams {
			var constraint types.Type
			if tp.Constraint != nil {
				constraint = b.resolveType(tp.Constraint)
			}
			typeParams[i] = &types.TypeParameter{
				Name:       tp.Name,
				Constraint: constraint,
			}
		}
		b.genericClasses[decl.Name] = &types.GenericClass{
			Name:       decl.Name,
			TypeParams: typeParams,
			Fields:     make(map[string]*types.Field),
			Methods:    make(map[string]*types.Method),
		}
	} else {
		b.classes[decl.Name] = &types.Class{
			Name:    decl.Name,
			Fields:  make(map[string]*types.Field),
			Methods: make(map[string]*types.Method),
		}
	}
}

func (b *Builder) resolveClass(decl *ast.ClassDecl) {
	// Handle generic classes separately
	if len(decl.TypeParams) > 0 {
		genericClass := b.genericClasses[decl.Name]

		// Put type parameters in scope for resolving field/method types
		savedTypeParams := b.typeParamScope
		b.typeParamScope = make(map[string]*types.TypeParameter)
		for _, tp := range genericClass.TypeParams {
			b.typeParamScope[tp.Name] = tp
		}

		for _, f := range decl.Fields {
			genericClass.Fields[f.Name] = &types.Field{
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
			genericClass.Constructor = &types.Function{
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
			genericClass.Methods[m.Name] = &types.Method{
				Name:       m.Name,
				Params:     params,
				ReturnType: b.resolveType(m.ReturnType),
			}
		}

		b.typeParamScope = savedTypeParams
		return
	}

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

func (b *Builder) collectInterface(decl *ast.InterfaceDecl) {
	b.interfaces[decl.Name] = &types.Interface{
		Name:    decl.Name,
		Methods: make(map[string]*types.InterfaceMethod),
	}
}

func (b *Builder) collectGoImport(decl *ast.GoImportDecl) {
	// Register imported names in the global scope
	for _, name := range decl.Names {
		fnType := types.GetGoPackageFunction(decl.Package, name)
		if fnType == nil {
			b.error(decl.Token.Line, decl.Token.Column,
				"unknown import %s from package %s", name, decl.Package)
			continue
		}
		b.scope.define(name, fnType)
	}

	// Add to goImports for codegen
	b.goImports = append(b.goImports, &GoImportDecl{
		Names:   decl.Names,
		Package: decl.Package,
	})
}

func (b *Builder) collectModuleImport(decl *ast.ModuleImportDecl) {
	// For now, register imported names with any type
	// In a full implementation, we would load the module and get the actual types
	for _, name := range decl.Names {
		// Register with any type for now - full implementation would resolve actual types
		b.scope.define(name, types.AnyType)
	}

	// Add to moduleImports for codegen
	b.moduleImports = append(b.moduleImports, &ModuleImportDecl{
		Names: decl.Names,
		Path:  decl.Path,
	})
}

func (b *Builder) collectExport(export *ast.ExportModifier) {
	// Extract the name of the exported declaration
	switch d := export.Decl.(type) {
	case *ast.FuncDecl:
		b.exports = append(b.exports, d.Name)
		// Also collect the function for type checking
		b.scope.define(d.Name, types.AnyType)
	case *ast.ClassDecl:
		b.exports = append(b.exports, d.Name)
		b.collectClass(d)
	case *ast.VarDecl:
		b.exports = append(b.exports, d.Name)
	case *ast.TypeAliasDecl:
		b.exports = append(b.exports, d.Name)
		b.collectTypeAlias(d)
	case *ast.InterfaceDecl:
		b.exports = append(b.exports, d.Name)
		b.collectInterface(d)
	}
}

func (b *Builder) buildExportedDecl(decl ast.Statement, result *Program) {
	// Build the declaration and add to the appropriate list in result
	switch d := decl.(type) {
	case *ast.FuncDecl:
		funcDecl := b.buildFuncDecl(d)
		result.Functions = append(result.Functions, funcDecl)
	case *ast.ClassDecl:
		classDecl := b.buildClassDecl(d)
		result.Classes = append(result.Classes, classDecl)
	case *ast.VarDecl:
		varDecl := b.buildVarDecl(d)
		if varDecl != nil {
			result.TopLevel = append(result.TopLevel, varDecl)
		}
	case *ast.TypeAliasDecl:
		alias := &TypeAlias{
			Name:     d.Name,
			Resolved: b.typeAliases[d.Name],
		}
		result.TypeAliases = append(result.TypeAliases, alias)
	case *ast.InterfaceDecl:
		ifaceDecl := b.buildInterfaceDecl(d)
		result.Interfaces = append(result.Interfaces, ifaceDecl)
	}
}

// buildDecl builds a typed declaration from an AST declaration.
func (b *Builder) buildDecl(decl ast.Statement) Stmt {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return b.buildFuncDecl(d)
	case *ast.ClassDecl:
		return b.buildClassDecl(d)
	case *ast.VarDecl:
		return b.buildVarDecl(d)
	default:
		return nil
	}
}

func (b *Builder) resolveInterface(decl *ast.InterfaceDecl) {
	iface := b.interfaces[decl.Name]

	for _, m := range decl.Methods {
		params := make([]*types.Param, len(m.Params))
		for i, p := range m.Params {
			params[i] = &types.Param{
				Name: p.Name,
				Type: b.resolveType(p.ParamType),
			}
		}
		iface.Methods[m.Name] = &types.InterfaceMethod{
			Name:       m.Name,
			Params:     params,
			ReturnType: b.resolveType(m.ReturnType),
		}
	}
}

func (b *Builder) buildInterfaceDecl(decl *ast.InterfaceDecl) *InterfaceDecl {
	iface := b.interfaces[decl.Name]

	methods := make([]*InterfaceMethodDecl, len(decl.Methods))
	for i, m := range decl.Methods {
		ifaceMethod := iface.Methods[m.Name]

		params := make([]*Param, len(m.Params))
		for j, p := range m.Params {
			params[j] = &Param{
				Name: p.Name,
				Type: ifaceMethod.Params[j].Type,
			}
		}

		methods[i] = &InterfaceMethodDecl{
			Name:       m.Name,
			Params:     params,
			ReturnType: ifaceMethod.ReturnType,
		}
	}

	return &InterfaceDecl{
		Name:    decl.Name,
		Methods: methods,
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

	case *ast.UnionType:
		unionTypes := make([]types.Type, len(t.Types))
		for i, ut := range t.Types {
			unionTypes[i] = b.resolveType(ut)
		}
		return types.MakeUnion(unionTypes...)

	case *ast.IntersectionType:
		intersectionTypes := make([]types.Type, len(t.Types))
		for i, it := range t.Types {
			intersectionTypes[i] = b.resolveType(it)
		}
		return types.MakeIntersection(intersectionTypes...)

	case *ast.LiteralType:
		// Convert AST literal type to types.Literal
		kind := types.KindAny
		switch t.Kind {
		case ast.TypeInt:
			kind = types.KindInt
		case ast.TypeFloat:
			kind = types.KindFloat
		case ast.TypeString:
			kind = types.KindString
		case ast.TypeBoolean:
			kind = types.KindBoolean
		}
		return &types.Literal{
			Kind:  kind,
			Value: t.Value,
		}

	case *ast.TupleType:
		// Convert AST tuple type to types.Tuple
		elements := make([]types.Type, len(t.Elements))
		for i, e := range t.Elements {
			elements[i] = b.resolveType(e)
		}
		var restType types.Type
		if t.RestElement != nil {
			// Extract element type from the array type
			if arrayType, ok := t.RestElement.(*ast.ArrayType); ok {
				restType = b.resolveType(arrayType.ElementType)
			}
		}
		return &types.Tuple{
			Elements: elements,
			Rest:     restType,
		}

	case *ast.MapType:
		return &types.Map{
			Key:   b.resolveType(t.KeyType),
			Value: b.resolveType(t.ValueType),
		}

	case *ast.SetType:
		return &types.Set{
			Element: b.resolveType(t.ElementType),
		}

	case *ast.PromiseType:
		return &types.Promise{
			Value: b.resolveType(t.ResultType),
		}

	case *ast.NamedType:
		// First check if it's a type parameter in scope
		if tp, ok := b.typeParamScope[t.Name]; ok {
			return tp
		}
		// Check for Promise<T>
		if t.Name == "Promise" && len(t.TypeArgs) == 1 {
			return &types.Promise{
				Value: b.resolveType(t.TypeArgs[0]),
			}
		}
		// Check for generic class instantiation (e.g., Stack<int>)
		if len(t.TypeArgs) > 0 {
			if genericClass, ok := b.genericClasses[t.Name]; ok {
				typeArgs := make([]types.Type, len(t.TypeArgs))
				for i, ta := range t.TypeArgs {
					typeArgs[i] = b.resolveType(ta)
				}
				instantiated, err := genericClass.Instantiate(typeArgs)
				if err != nil {
					b.error(0, 0, "cannot instantiate %s: %s", t.Name, err.Error())
					return types.AnyType
				}
				return instantiated
			}
		}
		if alias, ok := b.typeAliases[t.Name]; ok {
			return alias
		}
		if class, ok := b.classes[t.Name]; ok {
			return class
		}
		if genericClass, ok := b.genericClasses[t.Name]; ok {
			// Using generic class without type arguments - return the generic type itself
			return genericClass
		}
		if iface, ok := b.interfaces[t.Name]; ok {
			return iface
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
		// Special handling for empty object literal assigned to Map type
		if objLit, ok := decl.Value.(*ast.ObjectLiteral); ok && len(objLit.Properties) == 0 {
			if mapType, ok := varType.(*types.Map); ok {
				init = &MapLit{
					Entries:  []*MapEntry{},
					ExprType: mapType,
				}
			} else {
				init = b.buildExpr(decl.Value)
			}
		} else {
			init = b.buildExpr(decl.Value)
		}

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

	// Handle destructuring patterns
	if decl.Pattern != nil {
		pattern := b.buildPattern(decl.Pattern, varType, decl.IsConst)
		return &VarDecl{
			VarType: varType,
			Init:    init,
			IsConst: decl.IsConst,
			Pattern: pattern,
		}
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

func (b *Builder) buildPattern(pattern ast.Pattern, sourceType types.Type, isConst bool) Pattern {
	switch p := pattern.(type) {
	case *ast.ArrayPattern:
		return b.buildArrayPattern(p, sourceType, isConst)
	case *ast.ObjectPattern:
		return b.buildObjectPattern(p, sourceType, isConst)
	case *ast.IdentPattern:
		return b.buildIdentPattern(p, sourceType, isConst)
	default:
		return nil
	}
}

func (b *Builder) buildArrayPattern(pattern *ast.ArrayPattern, sourceType types.Type, isConst bool) *ArrayPattern {
	var elemType types.Type = types.AnyType

	if arrType, ok := sourceType.(*types.Array); ok {
		elemType = arrType.Element
	}

	elements := make([]Pattern, len(pattern.Elements))
	for i, elem := range pattern.Elements {
		if elem != nil {
			elements[i] = b.buildPattern(elem, elemType, isConst)
		}
	}

	return &ArrayPattern{
		Elements:    elements,
		PatternType: sourceType,
	}
}

func (b *Builder) buildObjectPattern(pattern *ast.ObjectPattern, sourceType types.Type, isConst bool) *ObjectPattern {
	properties := make([]*PropertyPattern, len(pattern.Properties))

	for i, prop := range pattern.Properties {
		var propType types.Type = types.AnyType

		// Try to get the property type from the source object type
		if objType, ok := sourceType.(*types.Object); ok {
			if p, exists := objType.Properties[prop.Key]; exists {
				propType = p.Type
			}
		}

		properties[i] = &PropertyPattern{
			Key:   prop.Key,
			Value: b.buildPattern(prop.Value, propType, isConst),
		}
	}

	return &ObjectPattern{
		Properties:  properties,
		PatternType: sourceType,
	}
}

func (b *Builder) buildIdentPattern(pattern *ast.IdentPattern, varType types.Type, isConst bool) *IdentPattern {
	if isConst {
		b.scope.defineConst(pattern.Name, varType)
	} else {
		b.scope.define(pattern.Name, varType)
	}

	return &IdentPattern{
		Name:        pattern.Name,
		PatternType: varType,
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
		// In async functions returning Promise<T>, we return T, not Promise<T>
		expectedType := b.currentFunc.ReturnType
		if b.inAsyncFunc {
			if promiseType, ok := expectedType.(*types.Promise); ok {
				expectedType = promiseType.Value
			}
		}
		if !types.IsAssignableTo(value.Type(), expectedType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"cannot return %s, expected %s", value.Type().String(), expectedType.String())
		}
	} else {
		// In async functions returning Promise<void>, empty return is OK
		expectedType := b.currentFunc.ReturnType
		if b.inAsyncFunc {
			if promiseType, ok := expectedType.(*types.Promise); ok {
				expectedType = promiseType.Value
			}
		}
		if !expectedType.Equals(types.VoidType) {
			b.error(stmt.Token.Line, stmt.Token.Column,
				"missing return value, expected %s", expectedType.String())
		}
	}

	return &ReturnStmt{Value: value}
}

func (b *Builder) buildFuncDecl(decl *ast.FuncDecl) *FuncDecl {
	// Handle type parameters
	var typeParams []*types.TypeParameter
	savedTypeParamScope := b.typeParamScope
	if len(decl.TypeParams) > 0 {
		b.typeParamScope = make(map[string]*types.TypeParameter)
		typeParams = make([]*types.TypeParameter, len(decl.TypeParams))
		for i, tp := range decl.TypeParams {
			var constraint types.Type
			if tp.Constraint != nil {
				constraint = b.resolveType(tp.Constraint)
			}
			typeParams[i] = &types.TypeParameter{
				Name:       tp.Name,
				Constraint: constraint,
			}
			b.typeParamScope[tp.Name] = typeParams[i]
		}
	}

	params := make([]*Param, len(decl.Params))
	typeParamsForFunc := make([]*types.Param, len(decl.Params))
	for i, p := range decl.Params {
		paramType := b.resolveType(p.ParamType)
		params[i] = &Param{
			Name: p.Name,
			Type: paramType,
		}
		typeParamsForFunc[i] = &types.Param{
			Name: p.Name,
			Type: paramType,
		}
	}

	returnType := b.resolveType(decl.ReturnType)

	// Register function type
	// If the function has decorators, register as any type since decorator can return anything
	if len(decl.Decorators) > 0 {
		b.scope.define(decl.Name, types.AnyType)
	} else if len(typeParams) > 0 {
		// Generic function
		b.scope.define(decl.Name, &types.GenericFunction{
			TypeParams: typeParams,
			Params:     typeParamsForFunc,
			ReturnType: returnType,
		})
	} else {
		// Regular function
		b.scope.define(decl.Name, &types.Function{
			Params:     typeParamsForFunc,
			ReturnType: returnType,
		})
	}

	funcType := &types.Function{
		Params:     typeParamsForFunc,
		ReturnType: returnType,
	}

	savedFunc := b.currentFunc
	savedInAsync := b.inAsyncFunc
	b.currentFunc = funcType
	b.inAsyncFunc = decl.IsAsync

	b.pushScope()
	for _, p := range params {
		b.scope.define(p.Name, p.Type)
	}

	body := b.buildBlock(decl.Body)

	b.popScope()
	b.currentFunc = savedFunc
	b.inAsyncFunc = savedInAsync
	b.typeParamScope = savedTypeParamScope

	return &FuncDecl{
		Name:       decl.Name,
		TypeParams: typeParams,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
		IsAsync:    decl.IsAsync,
		Decorators: b.buildDecorators(decl.Decorators),
	}
}

// buildDecorators builds typed decorators from AST decorators.
func (b *Builder) buildDecorators(decorators []*ast.Decorator) []*Decorator {
	if len(decorators) == 0 {
		return nil
	}

	result := make([]*Decorator, len(decorators))
	for i, d := range decorators {
		// Look up the decorator function in scope
		typ, ok := b.scope.lookup(d.Name)
		if !ok {
			b.error(d.Token.Line, d.Token.Column, "undefined decorator: %s", d.Name)
			typ = types.AnyType
		}
		result[i] = &Decorator{
			Name: d.Name,
			Type: typ,
		}
	}
	return result
}

func (b *Builder) buildClassDecl(decl *ast.ClassDecl) *ClassDecl {
	// Handle type parameters for generic classes
	var typeParams []*types.TypeParameter
	savedTypeParamScope := b.typeParamScope

	if len(decl.TypeParams) > 0 {
		genericClass := b.genericClasses[decl.Name]
		typeParams = genericClass.TypeParams

		// Put type parameters in scope
		b.typeParamScope = make(map[string]*types.TypeParameter)
		for _, tp := range typeParams {
			b.typeParamScope[tp.Name] = tp
		}

		savedClass := b.currentClass
		// For generic classes, we use a placeholder class during building
		b.currentClass = &types.Class{
			Name:    decl.Name,
			Fields:  genericClass.Fields,
			Methods: genericClass.Methods,
		}

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
			method := genericClass.Methods[m.Name]

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
		b.typeParamScope = savedTypeParamScope

		return &ClassDecl{
			Name:        decl.Name,
			TypeParams:  typeParams,
			Super:       decl.SuperClass,
			Fields:      fields,
			Constructor: constructor,
			Methods:     methods,
		}
	}

	// Non-generic class handling
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
		// Check if the original literal contains a decimal point to determine type
		// 5.0 should be float, 5 should be int
		if strings.Contains(e.Token.Literal, ".") {
			return &NumberLit{Value: e.Value, ExprType: types.FloatType}
		}
		return &NumberLit{Value: e.Value, ExprType: types.IntType}

	case *ast.StringLiteral:
		return &StringLit{Value: e.Value, ExprType: types.StringType}

	case *ast.TemplateLiteral:
		return b.buildTemplateLiteral(e)

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

	case *ast.SpreadExpr:
		return b.buildSpreadExpr(e)

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

	case *ast.AwaitExpr:
		return b.buildAwaitExpr(e)
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

func (b *Builder) buildSpreadExpr(expr *ast.SpreadExpr) Expr {
	argument := b.buildExpr(expr.Argument)

	// The argument should be an array type
	// The spread expression type is the same as the argument type
	return &SpreadExpr{
		Argument: argument,
		ExprType: argument.Type(),
	}
}

// isBuiltin checks if a name is a built-in function.
var builtins = map[string]types.Type{
	"println":    &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.VoidType},
	"print":      &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.VoidType},
	"len":        &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.IntType},
	"push":       &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "val", Type: types.AnyType}}, ReturnType: types.VoidType},
	"pop":        &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}}, ReturnType: types.AnyType},
	"typeof":     &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.StringType},
	"tostring":   &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.StringType},
	"toint":      &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.IntType},
	"tofloat":    &types.Function{Params: []*types.Param{{Name: "x", Type: types.AnyType}}, ReturnType: types.FloatType},
	"sqrt":       &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"floor":      &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"ceil":       &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	"abs":        &types.Function{Params: []*types.Param{{Name: "x", Type: types.FloatType}}, ReturnType: types.FloatType},
	// String methods
	"split":      &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}, {Name: "sep", Type: types.StringType}}, ReturnType: &types.Array{Element: types.StringType}},
	"join":       &types.Function{Params: []*types.Param{{Name: "arr", Type: &types.Array{Element: types.StringType}}, {Name: "sep", Type: types.StringType}}, ReturnType: types.StringType},
	"replace":    &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}, {Name: "old", Type: types.StringType}, {Name: "new", Type: types.StringType}}, ReturnType: types.StringType},
	"trim":       &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}}, ReturnType: types.StringType},
	"startsWith": &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}, {Name: "prefix", Type: types.StringType}}, ReturnType: types.BooleanType},
	"endsWith":   &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}, {Name: "suffix", Type: types.StringType}}, ReturnType: types.BooleanType},
	"includes":   &types.Function{Params: []*types.Param{{Name: "str", Type: types.StringType}, {Name: "substr", Type: types.StringType}}, ReturnType: types.BooleanType},
	// Array methods (higher-order functions)
	"map":       &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.AnyType},
	"filter":    &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.AnyType},
	"reduce":    &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "initial", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.AnyType},
	"find":      &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.AnyType},
	"findIndex": &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.IntType},
	"some":      &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.BooleanType},
	"every":     &types.Function{Params: []*types.Param{{Name: "arr", Type: types.AnyType}, {Name: "fn", Type: types.AnyType}}, ReturnType: types.BooleanType},
}

func (b *Builder) buildCallExpr(expr *ast.CallExpr) Expr {
	// Check for built-in function calls (only if not shadowed by user-defined function)
	if ident, ok := expr.Function.(*ast.Identifier); ok {
		// First check if there's a user-defined function with this name
		if varType, exists := b.scope.lookup(ident.Name); !exists {
			// Not in scope, check if it's a builtin
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

				// Special case: map returns array of callback return type
				if ident.Name == "map" && len(args) > 0 {
					argType := types.Unwrap(args[0].Type())
					if arrType, ok := argType.(*types.Array); ok {
						returnType = arrType // Default to same array type
					}
				}

				// Special case: filter returns same array type as input
				if ident.Name == "filter" && len(args) > 0 {
					argType := types.Unwrap(args[0].Type())
					if arrType, ok := argType.(*types.Array); ok {
						returnType = arrType
					}
				}

				// Special case: find returns element type of array
				if ident.Name == "find" && len(args) > 0 {
					argType := types.Unwrap(args[0].Type())
					if arrType, ok := argType.(*types.Array); ok {
						returnType = arrType.Element
					}
				}

				// Special case: reduce returns the type of initial value
				if ident.Name == "reduce" && len(args) > 1 {
					returnType = args[1].Type()
				}

				return &BuiltinCall{
					Name:     ident.Name,
					Args:     args,
					ExprType: returnType,
				}
			}
		} else {
			// User-defined function exists in scope, let it be handled by normal call logic below
			_ = varType // Suppress unused warning
		}
	}

	// Check for map method calls (e.g., m.get("key"), m.set("key", value))
	if propExpr, ok := expr.Function.(*ast.PropertyExpr); ok {
		objExpr := b.buildExpr(propExpr.Object)
		objType := types.Unwrap(objExpr.Type())
		if mapType, ok := objType.(*types.Map); ok {
			return b.buildMapMethodCall(objExpr, mapType, propExpr.Property, expr)
		}
		if setType, ok := objType.(*types.Set); ok {
			return b.buildSetMethodCall(objExpr, setType, propExpr.Property, expr)
		}
		if arrType, ok := objType.(*types.Array); ok {
			return b.buildArrayMethodCall(objExpr, arrType, propExpr.Property, expr)
		}
		if prim, ok := objType.(*types.Primitive); ok && prim.Kind == types.KindString {
			return b.buildStringMethodCall(objExpr, propExpr.Property, expr)
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

	// Handle generic function calls
	if gfn, ok := calleeType.(*types.GenericFunction); ok {
		args := make([]Expr, len(expr.Arguments))
		for i, arg := range expr.Arguments {
			args[i] = b.buildExpr(arg)
		}

		// Infer type arguments from actual arguments
		typeArgs := make([]types.Type, len(gfn.TypeParams))
		for i, tp := range gfn.TypeParams {
			// Find the first parameter that uses this type parameter and infer from argument
			for j, param := range gfn.Params {
				if typeParam, ok := param.Type.(*types.TypeParameter); ok && typeParam.Name == tp.Name {
					if j < len(args) {
						typeArgs[i] = args[j].Type()
						break
					}
				}
			}
			if typeArgs[i] == nil {
				typeArgs[i] = types.AnyType
			}
		}

		// Instantiate the generic function
		fn, err := gfn.Instantiate(typeArgs)
		if err != nil {
			b.error(expr.Token.Line, expr.Token.Column, "cannot instantiate generic function: %s", err.Error())
			return &CallExpr{
				Callee:   callee,
				Args:     args,
				Optional: expr.Optional,
				ExprType: types.AnyType,
			}
		}

		return &CallExpr{
			Callee:   callee,
			Args:     args,
			Optional: expr.Optional,
			ExprType: fn.ReturnType,
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
	// Check if this is enum member access (e.g., Color.Red)
	if ident, ok := expr.Object.(*ast.Identifier); ok {
		if enumType, ok := b.enums[ident.Name]; ok {
			member := enumType.GetMember(expr.Property)
			if member == nil {
				b.error(expr.Token.Line, expr.Token.Column,
					"enum %s does not have member %s", ident.Name, expr.Property)
				return &EnumMemberExpr{
					EnumName:   ident.Name,
					MemberName: expr.Property,
					ExprType:   enumType,
				}
			}
			return &EnumMemberExpr{
				EnumName:   ident.Name,
				MemberName: expr.Property,
				ExprType:   enumType,
			}
		}
	}

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

	case *types.Interface:
		// Access method on interface type
		if method, ok := obj.Methods[expr.Property]; ok {
			params := make([]*types.Param, len(method.Params))
			copy(params, method.Params)
			resultType = &types.Function{
				Params:     params,
				ReturnType: method.ReturnType,
			}
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"method %s does not exist on interface %s", expr.Property, obj.Name)
			resultType = types.AnyType
		}

	case *types.Array:
		// Arrays have a .length property
		if expr.Property == "length" {
			resultType = types.IntType
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on array type", expr.Property)
			resultType = types.AnyType
		}

	case *types.Map:
		// Maps have a .size property
		if expr.Property == "size" {
			resultType = types.IntType
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on Map type", expr.Property)
			resultType = types.AnyType
		}

	case *types.Set:
		// Sets have a .size property
		if expr.Property == "size" {
			resultType = types.IntType
		} else {
			b.error(expr.Token.Line, expr.Token.Column,
				"property %s does not exist on Set type", expr.Property)
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

func (b *Builder) buildTemplateLiteral(expr *ast.TemplateLiteral) Expr {
	expressions := make([]Expr, len(expr.Expressions))
	for i, e := range expr.Expressions {
		expressions[i] = b.buildExpr(e)
	}

	return &TemplateLit{
		Parts:       expr.Parts,
		Expressions: expressions,
		ExprType:    types.StringType,
	}
}

func (b *Builder) buildArrayLit(expr *ast.ArrayLiteral) Expr {
	elements := make([]Expr, len(expr.Elements))
	var elemType types.Type = types.AnyType

	if len(expr.Elements) > 0 {
		elements[0] = b.buildExpr(expr.Elements[0])
		// For spread expressions, get the element type of the spread array
		elemType = getElementTypeForArrayLit(elements[0])

		for i := 1; i < len(expr.Elements); i++ {
			elements[i] = b.buildExpr(expr.Elements[i])
			nextElemType := getElementTypeForArrayLit(elements[i])
			elemType = types.LeastUpperBound(elemType, nextElemType)
		}
	}

	return &ArrayLit{
		Elements: elements,
		ExprType: &types.Array{Element: elemType},
	}
}

// getElementTypeForArrayLit extracts the appropriate element type from an expression
// For spread expressions, it returns the element type of the spread array
// For regular expressions, it returns the expression's type
func getElementTypeForArrayLit(expr Expr) types.Type {
	if spread, ok := expr.(*SpreadExpr); ok {
		// If spreading an array, get its element type
		if arrType, ok := types.Unwrap(spread.Argument.Type()).(*types.Array); ok {
			return arrType.Element
		}
		// Fallback: return the spread's type
		return spread.Argument.Type()
	}
	return expr.Type()
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
	savedInAsync := b.inAsyncFunc
	b.currentFunc = funcType
	b.inAsyncFunc = expr.IsAsync

	b.pushScope()
	for _, p := range params {
		b.scope.define(p.Name, p.Type)
	}

	body := b.buildBlock(expr.Body)

	b.popScope()
	b.currentFunc = savedFunc
	b.inAsyncFunc = savedInAsync

	return &FuncExpr{
		Params:   params,
		Body:     body,
		ExprType: funcType,
		IsAsync:  expr.IsAsync,
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
	savedInAsync := b.inAsyncFunc
	b.currentFunc = funcType
	b.inAsyncFunc = expr.IsAsync

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
	b.inAsyncFunc = savedInAsync

	return &FuncExpr{
		Params:   params,
		Body:     body,
		BodyExpr: bodyExpr,
		ExprType: funcType,
		IsAsync:  expr.IsAsync,
	}
}

func (b *Builder) buildNewExpr(expr *ast.NewExpr) Expr {
	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	// Handle new Map<K, V>() - creates an empty map
	if expr.ClassName == "Map" && len(expr.TypeArgs) == 2 {
		keyType := b.resolveType(expr.TypeArgs[0])
		valType := b.resolveType(expr.TypeArgs[1])
		mapType := &types.Map{Key: keyType, Value: valType}
		return &MapLit{
			Entries:  []*MapEntry{},
			ExprType: mapType,
		}
	}

	// Handle new Set<T>() - creates an empty set
	if expr.ClassName == "Set" && len(expr.TypeArgs) == 1 {
		elemType := b.resolveType(expr.TypeArgs[0])
		setType := &types.Set{Element: elemType}
		return &SetLit{
			ExprType: setType,
		}
	}

	// Check for generic class
	if genericClass, ok := b.genericClasses[expr.ClassName]; ok {
		// Infer type arguments from constructor arguments
		typeArgs := make([]types.Type, len(genericClass.TypeParams))
		if genericClass.Constructor != nil {
			for i, tp := range genericClass.TypeParams {
				for j, param := range genericClass.Constructor.Params {
					if typeParam, ok := param.Type.(*types.TypeParameter); ok && typeParam.Name == tp.Name {
						if j < len(args) {
							typeArgs[i] = args[j].Type()
							break
						}
					}
				}
				if typeArgs[i] == nil {
					typeArgs[i] = types.AnyType
				}
			}
		} else {
			// No constructor, can't infer - use any
			for i := range typeArgs {
				typeArgs[i] = types.AnyType
			}
		}

		class, err := genericClass.Instantiate(typeArgs)
		if err != nil {
			b.error(expr.Token.Line, expr.Token.Column,
				"cannot instantiate generic class: %s", err.Error())
			return &NewExpr{
				ClassName: expr.ClassName,
				Args:      args,
				ExprType:  types.AnyType,
			}
		}

		if class.Constructor != nil {
			if len(args) != len(class.Constructor.Params) {
				b.error(expr.Token.Line, expr.Token.Column,
					"expected %d constructor arguments, got %d",
					len(class.Constructor.Params), len(args))
			}
		}

		return &NewExpr{
			ClassName: expr.ClassName,
			TypeArgs:  typeArgs,
			Args:      args,
			ExprType:  class,
		}
	}

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

// buildAwaitExpr handles await expressions.
func (b *Builder) buildAwaitExpr(expr *ast.AwaitExpr) Expr {
	// Verify we're in an async context
	if !b.inAsyncFunc {
		b.error(expr.Token.Line, expr.Token.Column,
			"'await' is only valid inside async functions")
	}

	arg := b.buildExpr(expr.Argument)
	argType := types.Unwrap(arg.Type())

	// Check if argument is a Promise
	if promiseType, ok := argType.(*types.Promise); ok {
		return &AwaitExpr{
			Argument: arg,
			ExprType: promiseType.Value, // Unwrap Promise<T> to T
		}
	}

	// If not a Promise, the await expression returns the value as-is (TypeScript behavior)
	return &AwaitExpr{
		Argument: arg,
		ExprType: argType,
	}
}

// buildMapMethodCall handles method calls on Map types (get, set, has, delete, keys, values)
func (b *Builder) buildMapMethodCall(obj Expr, mapType *types.Map, method string, expr *ast.CallExpr) Expr {
	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	var resultType types.Type

	switch method {
	case "get":
		// get(key) returns value type
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.get expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), mapType.Key) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.get key type mismatch: expected %s, got %s",
				mapType.Key.String(), args[0].Type().String())
		}
		resultType = mapType.Value

	case "set":
		// set(key, value) returns void
		if len(args) != 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.set expects 2 arguments, got %d", len(args))
		} else {
			if !types.IsAssignableTo(args[0].Type(), mapType.Key) {
				b.error(expr.Token.Line, expr.Token.Column,
					"Map.set key type mismatch: expected %s, got %s",
					mapType.Key.String(), args[0].Type().String())
			}
			if !types.IsAssignableTo(args[1].Type(), mapType.Value) {
				b.error(expr.Token.Line, expr.Token.Column,
					"Map.set value type mismatch: expected %s, got %s",
					mapType.Value.String(), args[1].Type().String())
			}
		}
		resultType = types.VoidType

	case "has":
		// has(key) returns boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.has expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), mapType.Key) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.has key type mismatch: expected %s, got %s",
				mapType.Key.String(), args[0].Type().String())
		}
		resultType = types.BooleanType

	case "delete":
		// delete(key) returns void
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.delete expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), mapType.Key) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.delete key type mismatch: expected %s, got %s",
				mapType.Key.String(), args[0].Type().String())
		}
		resultType = types.VoidType

	case "keys":
		// keys() returns key[]
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.keys expects 0 arguments, got %d", len(args))
		}
		resultType = &types.Array{Element: mapType.Key}

	case "values":
		// values() returns value[]
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.values expects 0 arguments, got %d", len(args))
		}
		resultType = &types.Array{Element: mapType.Value}

	case "entries":
		// entries() returns [K, V][] - approximated as any[] for now
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.entries expects 0 arguments, got %d", len(args))
		}
		// Return type is an array of tuples - we'll use any[] as approximation
		resultType = &types.Array{Element: types.AnyType}

	case "clear":
		// clear() returns void
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.clear expects 0 arguments, got %d", len(args))
		}
		resultType = types.VoidType

	case "forEach":
		// forEach(callback: (value: V, key: K) => void) returns void
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Map.forEach expects 1 argument, got %d", len(args))
		}
		resultType = types.VoidType

	default:
		b.error(expr.Token.Line, expr.Token.Column,
			"unknown Map method: %s", method)
		resultType = types.AnyType
	}

	return &MethodCallExpr{
		Object:   obj,
		Method:   method,
		Args:     args,
		ExprType: resultType,
	}
}

// buildSetMethodCall handles method calls on Set types (add, has, delete, clear, values, forEach)
func (b *Builder) buildSetMethodCall(obj Expr, setType *types.Set, method string, expr *ast.CallExpr) Expr {
	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	var resultType types.Type

	switch method {
	case "add":
		// add(value) returns Set<T>
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.add expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), setType.Element) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.add value type mismatch: expected %s, got %s",
				setType.Element.String(), args[0].Type().String())
		}
		resultType = setType

	case "has":
		// has(value) returns boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.has expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), setType.Element) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.has value type mismatch: expected %s, got %s",
				setType.Element.String(), args[0].Type().String())
		}
		resultType = types.BooleanType

	case "delete":
		// delete(value) returns boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.delete expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), setType.Element) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.delete value type mismatch: expected %s, got %s",
				setType.Element.String(), args[0].Type().String())
		}
		resultType = types.BooleanType

	case "clear":
		// clear() returns void
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.clear expects 0 arguments, got %d", len(args))
		}
		resultType = types.VoidType

	case "values":
		// values() returns T[]
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.values expects 0 arguments, got %d", len(args))
		}
		resultType = &types.Array{Element: setType.Element}

	case "forEach":
		// forEach(callback: (value: T) => void) returns void
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Set.forEach expects 1 argument, got %d", len(args))
		}
		resultType = types.VoidType

	default:
		b.error(expr.Token.Line, expr.Token.Column,
			"unknown Set method: %s", method)
		resultType = types.AnyType
	}

	return &MethodCallExpr{
		Object:   obj,
		Method:   method,
		Args:     args,
		ExprType: resultType,
	}
}

// buildArrayMethodCall handles method calls on Array types
func (b *Builder) buildArrayMethodCall(obj Expr, arrType *types.Array, method string, expr *ast.CallExpr) Expr {
	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	var resultType types.Type

	switch method {
	case "push":
		// push(value: T, ...): int
		if len(args) < 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.push expects at least 1 argument, got %d", len(args))
		} else {
			for i, arg := range args {
				if !types.IsAssignableTo(arg.Type(), arrType.Element) {
					b.error(expr.Token.Line, expr.Token.Column,
						"Array.push argument %d type mismatch: expected %s, got %s",
						i+1, arrType.Element.String(), arg.Type().String())
				}
			}
		}
		resultType = types.IntType

	case "pop":
		// pop(): T | null
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.pop expects 0 arguments, got %d", len(args))
		}
		resultType = &types.Nullable{Inner: arrType.Element}

	case "shift":
		// shift(): T | null
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.shift expects 0 arguments, got %d", len(args))
		}
		resultType = &types.Nullable{Inner: arrType.Element}

	case "unshift":
		// unshift(value: T, ...): int
		if len(args) < 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.unshift expects at least 1 argument, got %d", len(args))
		} else {
			for i, arg := range args {
				if !types.IsAssignableTo(arg.Type(), arrType.Element) {
					b.error(expr.Token.Line, expr.Token.Column,
						"Array.unshift argument %d type mismatch: expected %s, got %s",
						i+1, arrType.Element.String(), arg.Type().String())
				}
			}
		}
		resultType = types.IntType

	case "slice":
		// slice(start?: int, end?: int): T[]
		if len(args) > 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.slice expects 0-2 arguments, got %d", len(args))
		}
		resultType = arrType

	case "splice":
		// splice(start: int, deleteCount?: int, ...items: T[]): T[]
		if len(args) < 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.splice expects at least 1 argument, got %d", len(args))
		}
		resultType = arrType

	case "concat":
		// concat(...arrays: T[][]): T[]
		resultType = arrType

	case "indexOf":
		// indexOf(value: T, fromIndex?: int): int
		if len(args) < 1 || len(args) > 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.indexOf expects 1-2 arguments, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), arrType.Element) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.indexOf value type mismatch: expected %s, got %s",
				arrType.Element.String(), args[0].Type().String())
		}
		resultType = types.IntType

	case "includes":
		// includes(value: T, fromIndex?: int): boolean
		if len(args) < 1 || len(args) > 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.includes expects 1-2 arguments, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), arrType.Element) {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.includes value type mismatch: expected %s, got %s",
				arrType.Element.String(), args[0].Type().String())
		}
		resultType = types.BooleanType

	case "join":
		// join(separator?: string): string
		if len(args) > 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.join expects 0-1 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "reverse":
		// reverse(): T[]
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.reverse expects 0 arguments, got %d", len(args))
		}
		resultType = arrType

	case "sort":
		// sort(compareFn?: (a: T, b: T) => int): T[]
		if len(args) > 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.sort expects 0-1 arguments, got %d", len(args))
		}
		resultType = arrType

	case "map":
		// map(callback: (value: T, index: int) => U): U[]
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.map expects 1 argument, got %d", len(args))
			resultType = &types.Array{Element: types.AnyType}
		} else {
			// Try to get return type from callback
			cbType := types.Unwrap(args[0].Type())
			if fn, ok := cbType.(*types.Function); ok {
				resultType = &types.Array{Element: fn.ReturnType}
			} else {
				resultType = &types.Array{Element: types.AnyType}
			}
		}

	case "filter":
		// filter(callback: (value: T, index: int) => boolean): T[]
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.filter expects 1 argument, got %d", len(args))
		}
		resultType = arrType

	case "reduce":
		// reduce(callback: (acc: U, value: T, index: int) => U, initialValue: U): U
		if len(args) != 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.reduce expects 2 arguments, got %d", len(args))
			resultType = types.AnyType
		} else {
			// Return type is the type of initialValue
			resultType = args[1].Type()
		}

	case "forEach":
		// forEach(callback: (value: T, index: int) => void): void
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.forEach expects 1 argument, got %d", len(args))
		}
		resultType = types.VoidType

	case "find":
		// find(callback: (value: T, index: int) => boolean): T | null
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.find expects 1 argument, got %d", len(args))
		}
		resultType = &types.Nullable{Inner: arrType.Element}

	case "findIndex":
		// findIndex(callback: (value: T, index: int) => boolean): int
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.findIndex expects 1 argument, got %d", len(args))
		}
		resultType = types.IntType

	case "some":
		// some(callback: (value: T, index: int) => boolean): boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.some expects 1 argument, got %d", len(args))
		}
		resultType = types.BooleanType

	case "every":
		// every(callback: (value: T, index: int) => boolean): boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"Array.every expects 1 argument, got %d", len(args))
		}
		resultType = types.BooleanType

	default:
		b.error(expr.Token.Line, expr.Token.Column,
			"unknown Array method: %s", method)
		resultType = types.AnyType
	}

	return &MethodCallExpr{
		Object:   obj,
		Method:   method,
		Args:     args,
		ExprType: resultType,
	}
}

// buildStringMethodCall handles method calls on string types
func (b *Builder) buildStringMethodCall(obj Expr, method string, expr *ast.CallExpr) Expr {
	args := make([]Expr, len(expr.Arguments))
	for i, arg := range expr.Arguments {
		args[i] = b.buildExpr(arg)
	}

	var resultType types.Type

	switch method {
	case "split":
		// split(separator: string): string[]
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.split expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), types.StringType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.split separator must be string, got %s", args[0].Type().String())
		}
		resultType = &types.Array{Element: types.StringType}

	case "replace":
		// replace(old: string, new: string): string
		if len(args) != 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.replace expects 2 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "trim":
		// trim(): string
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.trim expects 0 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "startsWith":
		// startsWith(prefix: string): boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.startsWith expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), types.StringType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.startsWith prefix must be string, got %s", args[0].Type().String())
		}
		resultType = types.BooleanType

	case "endsWith":
		// endsWith(suffix: string): boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.endsWith expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), types.StringType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.endsWith suffix must be string, got %s", args[0].Type().String())
		}
		resultType = types.BooleanType

	case "includes":
		// includes(substring: string): boolean
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.includes expects 1 argument, got %d", len(args))
		} else if !types.IsAssignableTo(args[0].Type(), types.StringType) {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.includes substring must be string, got %s", args[0].Type().String())
		}
		resultType = types.BooleanType

	case "toLowerCase":
		// toLowerCase(): string
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.toLowerCase expects 0 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "toUpperCase":
		// toUpperCase(): string
		if len(args) != 0 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.toUpperCase expects 0 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "substring":
		// substring(start: int, end?: int): string
		if len(args) < 1 || len(args) > 2 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.substring expects 1-2 arguments, got %d", len(args))
		}
		resultType = types.StringType

	case "charAt":
		// charAt(index: int): string
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.charAt expects 1 argument, got %d", len(args))
		}
		resultType = types.StringType

	case "indexOf":
		// indexOf(substring: string): int
		if len(args) != 1 {
			b.error(expr.Token.Line, expr.Token.Column,
				"String.indexOf expects 1 argument, got %d", len(args))
		}
		resultType = types.IntType

	default:
		b.error(expr.Token.Line, expr.Token.Column,
			"unknown String method: %s", method)
		resultType = types.AnyType
	}

	return &MethodCallExpr{
		Object:   obj,
		Method:   method,
		Args:     args,
		ExprType: resultType,
	}
}
