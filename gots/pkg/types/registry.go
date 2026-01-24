package types

import (
	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/declaration"
)

// GoPackageRegistry contains type information for Go standard library packages.
// This allows the type checker to resolve imported Go functions.
// NOTE: This registry is a fallback. The preferred way is to use .d.gts declaration files.
var GoPackageRegistry = map[string]map[string]Type{
	"fmt": {
		"Println":  &Function{Params: []*Param{{Name: "args", Type: AnyType}}, ReturnType: VoidType},
		"Print":    &Function{Params: []*Param{{Name: "args", Type: AnyType}}, ReturnType: VoidType},
		"Printf":   &Function{Params: []*Param{{Name: "format", Type: StringType}, {Name: "args", Type: AnyType}}, ReturnType: VoidType},
		"Sprintf":  &Function{Params: []*Param{{Name: "format", Type: StringType}, {Name: "args", Type: AnyType}}, ReturnType: StringType},
		"Errorf":   &Function{Params: []*Param{{Name: "format", Type: StringType}, {Name: "args", Type: AnyType}}, ReturnType: AnyType},
		"Sscanf":   &Function{Params: []*Param{{Name: "str", Type: StringType}, {Name: "format", Type: StringType}, {Name: "args", Type: AnyType}}, ReturnType: IntType},
		"Fscanf":   &Function{Params: []*Param{{Name: "r", Type: AnyType}, {Name: "format", Type: StringType}, {Name: "args", Type: AnyType}}, ReturnType: IntType},
	},
	"strings": {
		"Join":       &Function{Params: []*Param{{Name: "a", Type: &Array{Element: StringType}}, {Name: "sep", Type: StringType}}, ReturnType: StringType},
		"Split":      &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "sep", Type: StringType}}, ReturnType: &Array{Element: StringType}},
		"Contains":   &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "substr", Type: StringType}}, ReturnType: BooleanType},
		"HasPrefix":  &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "prefix", Type: StringType}}, ReturnType: BooleanType},
		"HasSuffix":  &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "suffix", Type: StringType}}, ReturnType: BooleanType},
		"ToUpper":    &Function{Params: []*Param{{Name: "s", Type: StringType}}, ReturnType: StringType},
		"ToLower":    &Function{Params: []*Param{{Name: "s", Type: StringType}}, ReturnType: StringType},
		"TrimSpace":  &Function{Params: []*Param{{Name: "s", Type: StringType}}, ReturnType: StringType},
		"Replace":    &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "old", Type: StringType}, {Name: "new", Type: StringType}, {Name: "n", Type: IntType}}, ReturnType: StringType},
		"ReplaceAll": &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "old", Type: StringType}, {Name: "new", Type: StringType}}, ReturnType: StringType},
		"Index":      &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "substr", Type: StringType}}, ReturnType: IntType},
		"Count":      &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "substr", Type: StringType}}, ReturnType: IntType},
	},
	"strconv": {
		"Itoa":      &Function{Params: []*Param{{Name: "i", Type: IntType}}, ReturnType: StringType},
		"Atoi":      &Function{Params: []*Param{{Name: "s", Type: StringType}}, ReturnType: IntType},
		"ParseInt":  &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "base", Type: IntType}, {Name: "bitSize", Type: IntType}}, ReturnType: IntType},
		"ParseFloat": &Function{Params: []*Param{{Name: "s", Type: StringType}, {Name: "bitSize", Type: IntType}}, ReturnType: FloatType},
		"FormatInt": &Function{Params: []*Param{{Name: "i", Type: IntType}, {Name: "base", Type: IntType}}, ReturnType: StringType},
		"FormatFloat": &Function{Params: []*Param{{Name: "f", Type: FloatType}, {Name: "fmt", Type: IntType}, {Name: "prec", Type: IntType}, {Name: "bitSize", Type: IntType}}, ReturnType: StringType},
	},
	"math": {
		"Sqrt":  &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Abs":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Floor": &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Ceil":  &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Round": &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Sin":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Cos":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Tan":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Log":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Log10": &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Exp":   &Function{Params: []*Param{{Name: "x", Type: FloatType}}, ReturnType: FloatType},
		"Pow":   &Function{Params: []*Param{{Name: "x", Type: FloatType}, {Name: "y", Type: FloatType}}, ReturnType: FloatType},
		"Max":   &Function{Params: []*Param{{Name: "x", Type: FloatType}, {Name: "y", Type: FloatType}}, ReturnType: FloatType},
		"Min":   &Function{Params: []*Param{{Name: "x", Type: FloatType}, {Name: "y", Type: FloatType}}, ReturnType: FloatType},
		"Mod":   &Function{Params: []*Param{{Name: "x", Type: FloatType}, {Name: "y", Type: FloatType}}, ReturnType: FloatType},
	},
	"time": {
		"Now":   &Function{Params: []*Param{}, ReturnType: AnyType},
		"Sleep": &Function{Params: []*Param{{Name: "d", Type: IntType}}, ReturnType: VoidType},
		"Since": &Function{Params: []*Param{{Name: "t", Type: AnyType}}, ReturnType: IntType},
	},
	"os": {
		"Exit":    &Function{Params: []*Param{{Name: "code", Type: IntType}}, ReturnType: VoidType},
		"Getenv":  &Function{Params: []*Param{{Name: "key", Type: StringType}}, ReturnType: StringType},
		"Setenv":  &Function{Params: []*Param{{Name: "key", Type: StringType}, {Name: "value", Type: StringType}}, ReturnType: AnyType},
		"Getwd":   &Function{Params: []*Param{}, ReturnType: StringType},
		"Chdir":   &Function{Params: []*Param{{Name: "dir", Type: StringType}}, ReturnType: AnyType},
		"Mkdir":   &Function{Params: []*Param{{Name: "name", Type: StringType}, {Name: "perm", Type: IntType}}, ReturnType: AnyType},
		"Remove":  &Function{Params: []*Param{{Name: "name", Type: StringType}}, ReturnType: AnyType},
		"ReadFile": &Function{Params: []*Param{{Name: "name", Type: StringType}}, ReturnType: &Array{Element: IntType}},
	},
	"regexp": {
		"MatchString": &Function{Params: []*Param{{Name: "pattern", Type: StringType}, {Name: "s", Type: StringType}}, ReturnType: BooleanType},
		"Compile":     &Function{Params: []*Param{{Name: "expr", Type: StringType}}, ReturnType: AnyType},
		"MustCompile": &Function{Params: []*Param{{Name: "expr", Type: StringType}}, ReturnType: AnyType},
	},
	"json": {
		"Marshal":   &Function{Params: []*Param{{Name: "v", Type: AnyType}}, ReturnType: &Array{Element: IntType}},
		"Unmarshal": &Function{Params: []*Param{{Name: "data", Type: &Array{Element: IntType}}, {Name: "v", Type: AnyType}}, ReturnType: AnyType},
	},
}

// GoPackageConstants contains constants from Go standard library packages.
var GoPackageConstants = map[string]map[string]Type{
	"math": {
		"Pi": FloatType,
		"E":  FloatType,
	},
}

// GetGoPackageFunction returns the type of a function from a Go package.
// It first tries to load from .d.gts declaration files, then falls back to the hardcoded registry.
// Returns nil if the package or function is not found.
func GetGoPackageFunction(pkg, name string) Type {
	// Try declaration loader first
	if fn, err := declaration.DefaultLoader.GetFunction("go:"+pkg, name); err == nil {
		return convertDeclFunctionToType(fn)
	}

	// Fall back to hardcoded registry
	if pkgFuncs, ok := GoPackageRegistry[pkg]; ok {
		if fn, ok := pkgFuncs[name]; ok {
			return fn
		}
	}
	// Check constants as well
	if pkgConsts, ok := GoPackageConstants[pkg]; ok {
		if c, ok := pkgConsts[name]; ok {
			return c
		}
	}
	return nil
}

// GetGoPackageConstant returns the type of a constant from a Go package.
func GetGoPackageConstant(pkg, name string) Type {
	// Try declaration loader first
	if t, err := declaration.DefaultLoader.GetConstant("go:"+pkg, name); err == nil {
		return convertAstTypeToType(t)
	}

	// Fall back to hardcoded registry
	if pkgConsts, ok := GoPackageConstants[pkg]; ok {
		if c, ok := pkgConsts[name]; ok {
			return c
		}
	}
	return nil
}

// convertDeclFunctionToType converts a declaration.FunctionInfo to types.Type
func convertDeclFunctionToType(fn *declaration.FunctionInfo) Type {
	params := make([]*Param, len(fn.Params))
	for i, p := range fn.Params {
		params[i] = &Param{
			Name: p.Name,
			Type: convertAstTypeToType(p.ParamType),
		}
	}
	return &Function{
		Params:     params,
		ReturnType: convertAstTypeToType(fn.ReturnType),
	}
}

// convertAstTypeToType converts an ast.Type to types.Type
func convertAstTypeToType(t ast.Type) Type {
	if t == nil {
		return AnyType
	}
	switch at := t.(type) {
	case *ast.PrimitiveType:
		switch at.Kind {
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
		return &Array{Element: convertAstTypeToType(at.ElementType)}
	case *ast.NullableType:
		return &Nullable{Inner: convertAstTypeToType(at.Inner)}
	case *ast.AnyType:
		return AnyType
	case *ast.ByteType:
		return IntType // byte maps to int in GoTS
	case *ast.NamedType:
		// For now, return any for named types (interfaces, classes, etc.)
		return AnyType
	case *ast.ObjectType:
		props := make(map[string]*Property)
		for _, p := range at.Properties {
			props[p.Name] = &Property{
				Name: p.Name,
				Type: convertAstTypeToType(p.PropType),
			}
		}
		return &Object{Properties: props}
	case *ast.FunctionType:
		params := make([]*Param, len(at.ParamTypes))
		for i, pt := range at.ParamTypes {
			params[i] = &Param{
				Name: "",
				Type: convertAstTypeToType(pt),
			}
		}
		return &Function{
			Params:     params,
			ReturnType: convertAstTypeToType(at.ReturnType),
		}
	case *ast.TupleType:
		// Tuples from declarations represent Go's multiple return values
		// For now, we'll treat them as any since GoTS doesn't fully support tuples
		return AnyType
	}
	return AnyType
}
