package typed

import (
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/types"
)

// BuiltinConstant defines a constant on a built-in object (e.g., Math.PI).
type BuiltinConstant struct {
	Type   types.Type
	GoCode string // The Go code to generate (e.g., "math.Pi")
}

// BuiltinMethod defines a method on a built-in object (e.g., Math.round()).
type BuiltinMethod struct {
	Params     []*types.Param
	ReturnType types.Type
	Variadic   bool // If true, accepts variable number of arguments
	// GoCodeGen generates Go code for this method call.
	// args contains the generated Go code for each argument.
	GoCodeGen func(args []string) string
}

// BuiltinObject defines a global built-in object (e.g., Math, JSON).
type BuiltinObject struct {
	Name      string
	Constants map[string]*BuiltinConstant
	Methods   map[string]*BuiltinMethod
	// Imports lists Go packages that need to be imported when using this object.
	Imports []string
}

// BuiltinRegistry holds all built-in objects.
var BuiltinRegistry = map[string]*BuiltinObject{}

// RegisterBuiltin adds a built-in object to the registry.
func RegisterBuiltin(obj *BuiltinObject) {
	BuiltinRegistry[obj.Name] = obj
}

// GetBuiltin retrieves a built-in object by name.
func GetBuiltin(name string) (*BuiltinObject, bool) {
	obj, ok := BuiltinRegistry[name]
	return obj, ok
}

// IsBuiltinObject checks if a name is a registered built-in object.
func IsBuiltinObject(name string) bool {
	_, ok := BuiltinRegistry[name]
	return ok
}

// BuiltinObjectCall represents a method call on a built-in object (e.g., Math.round(x)).
type BuiltinObjectCall struct {
	Object   string     // "Math", "JSON", etc.
	Method   string     // "round", "parse", etc.
	Args     []Expr
	ExprType types.Type
}

func (b *BuiltinObjectCall) exprNode()        {}
func (b *BuiltinObjectCall) Type() types.Type { return b.ExprType }

// BuiltinObjectConstant represents access to a constant on a built-in object (e.g., Math.PI).
type BuiltinObjectConstant struct {
	Object   string // "Math", etc.
	Name     string // "PI", "E", etc.
	ExprType types.Type
}

func (b *BuiltinObjectConstant) exprNode()        {}
func (b *BuiltinObjectConstant) Type() types.Type { return b.ExprType }

// ----------------------------------------------------------------------------
// Math Built-in Object
// ----------------------------------------------------------------------------

// mathUnary creates a BuiltinMethod for a single-arg math function (e.g., math.Floor).
func mathUnary(goFunc string) *BuiltinMethod {
	return &BuiltinMethod{
		Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
		ReturnType: types.NumberType,
		GoCodeGen:  func(args []string) string { return fmt.Sprintf("%s(%s)", goFunc, args[0]) },
	}
}

// mathBinary creates a BuiltinMethod for a two-arg math function (e.g., math.Pow).
func mathBinary(paramA, paramB string, goFunc string) *BuiltinMethod {
	return &BuiltinMethod{
		Params:     []*types.Param{{Name: paramA, Type: types.NumberType}, {Name: paramB, Type: types.NumberType}},
		ReturnType: types.NumberType,
		GoCodeGen:  func(args []string) string { return fmt.Sprintf("%s(%s, %s)", goFunc, args[0], args[1]) },
	}
}

// chainedMinMax creates a variadic BuiltinMethod that chains calls to goFunc (e.g., math.Min).
func chainedMinMax(goFunc string) *BuiltinMethod {
	return &BuiltinMethod{
		Params:     []*types.Param{{Name: "values", Type: types.NumberType}},
		ReturnType: types.NumberType,
		Variadic:   true,
		GoCodeGen: func(args []string) string {
			if len(args) == 1 {
				return args[0]
			}
			result := fmt.Sprintf("%s(%s, %s)", goFunc, args[0], args[1])
			for i := 2; i < len(args); i++ {
				result = fmt.Sprintf("%s(%s, %s)", goFunc, result, args[i])
			}
			return result
		},
	}
}

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:    "Math",
		Imports: []string{"math", "math/rand"},
		Constants: map[string]*BuiltinConstant{
			"PI": {Type: types.NumberType, GoCode: "math.Pi"},
			"E":  {Type: types.NumberType, GoCode: "math.E"},
		},
		Methods: map[string]*BuiltinMethod{
			// Rounding
			"round": mathUnary("math.Round"),
			"floor": mathUnary("math.Floor"),
			"ceil":  mathUnary("math.Ceil"),
			"trunc": mathUnary("math.Trunc"),
			// Power and roots
			"sqrt": mathUnary("math.Sqrt"),
			"cbrt": mathUnary("math.Cbrt"),
			"pow":  mathBinary("x", "y", "math.Pow"),
			"exp":  mathUnary("math.Exp"),
			// Logarithms
			"log":   mathUnary("math.Log"),
			"log10": mathUnary("math.Log10"),
			"log2":  mathUnary("math.Log2"),
			// Absolute value and sign
			"abs": mathUnary("math.Abs"),
			"sign": {
				Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
				ReturnType: types.NumberType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() float64 { x := %s; if x > 0 { return 1 } else if x < 0 { return -1 }; return 0 }()", args[0])
				},
			},
			// Min/Max (variadic)
			"min": chainedMinMax("math.Min"),
			"max": chainedMinMax("math.Max"),
			// Trigonometric
			"sin":   mathUnary("math.Sin"),
			"cos":   mathUnary("math.Cos"),
			"tan":   mathUnary("math.Tan"),
			"asin":  mathUnary("math.Asin"),
			"acos":  mathUnary("math.Acos"),
			"atan":  mathUnary("math.Atan"),
			"atan2": mathBinary("y", "x", "math.Atan2"),
			// Random
			"random": {
				Params:     []*types.Param{},
				ReturnType: types.NumberType,
				GoCodeGen:  func(args []string) string { return "rand.Float64()" },
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Helper functions for builder
// ----------------------------------------------------------------------------

// BuildBuiltinMethodCall creates a BuiltinCall expression for a method call.
func BuildBuiltinMethodCall(objName, methodName string, args []Expr, line, col int) (Expr, error) {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return nil, fmt.Errorf("unknown built-in object: %s", objName)
	}

	method, ok := obj.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("unknown method %s.%s", objName, methodName)
	}

	// Validate argument count
	expectedArgs := len(method.Params)
	if method.Variadic {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s.%s expects at least 1 argument, got %d", objName, methodName, len(args))
		}
	} else if len(args) != expectedArgs {
		return nil, fmt.Errorf("%s.%s expects %d arguments, got %d", objName, methodName, expectedArgs, len(args))
	}

	return &BuiltinObjectCall{
		Object:   objName,
		Method:   methodName,
		Args:     args,
		ExprType: method.ReturnType,
	}, nil
}

// BuildBuiltinConstant creates a BuiltinObjectConstant for a constant access.
func BuildBuiltinConstant(objName, constName string) (Expr, error) {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return nil, fmt.Errorf("unknown built-in object: %s", objName)
	}

	constant, ok := obj.Constants[constName]
	if !ok {
		return nil, fmt.Errorf("unknown constant %s.%s", objName, constName)
	}

	return &BuiltinObjectConstant{
		Object:   objName,
		Name:     constName,
		ExprType: constant.Type,
	}, nil
}

// HasBuiltinConstant checks if a built-in object has a constant with the given name.
func HasBuiltinConstant(objName, constName string) bool {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return false
	}
	_, ok = obj.Constants[constName]
	return ok
}

// HasBuiltinMethod checks if a built-in object has a method with the given name.
func HasBuiltinMethod(objName, methodName string) bool {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return false
	}
	_, ok = obj.Methods[methodName]
	return ok
}

// ----------------------------------------------------------------------------
// Helper functions for codegen
// ----------------------------------------------------------------------------

// GenerateBuiltinCall generates Go code for a built-in method call.
func GenerateBuiltinCall(objName, methodName string, args []string) (string, error) {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return "", fmt.Errorf("unknown built-in object: %s", objName)
	}

	method, ok := obj.Methods[methodName]
	if !ok {
		return "", fmt.Errorf("unknown method %s.%s", objName, methodName)
	}

	return method.GoCodeGen(args), nil
}

// GenerateBuiltinConstant generates Go code for a built-in constant access.
func GenerateBuiltinConstant(objName, constName string) (string, error) {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return "", fmt.Errorf("unknown built-in object: %s", objName)
	}

	constant, ok := obj.Constants[constName]
	if !ok {
		return "", fmt.Errorf("unknown constant %s.%s", objName, constName)
	}

	return constant.GoCode, nil
}

// GetBuiltinImports returns the Go imports needed for a built-in object.
func GetBuiltinImports(objName string) []string {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return nil
	}
	return obj.Imports
}

// GetAllBuiltinNames returns all registered built-in object names.
func GetAllBuiltinNames() []string {
	names := make([]string, 0, len(BuiltinRegistry))
	for name := range BuiltinRegistry {
		names = append(names, name)
	}
	return names
}

// ----------------------------------------------------------------------------
// Number Built-in Object
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:    "Number",
		Imports: []string{"math", "strconv"},
		Constants: map[string]*BuiltinConstant{
			"MAX_SAFE_INTEGER": {Type: types.NumberType, GoCode: "float64(9007199254740991)"},
			"MIN_SAFE_INTEGER": {Type: types.NumberType, GoCode: "float64(-9007199254740991)"},
			"MAX_VALUE":        {Type: types.NumberType, GoCode: "math.MaxFloat64"},
			"MIN_VALUE":        {Type: types.NumberType, GoCode: "math.SmallestNonzeroFloat64"},
			"POSITIVE_INFINITY": {Type: types.NumberType, GoCode: "math.Inf(1)"},
			"NEGATIVE_INFINITY": {Type: types.NumberType, GoCode: "math.Inf(-1)"},
			"NaN":              {Type: types.NumberType, GoCode: "math.NaN()"},
		},
		Methods: map[string]*BuiltinMethod{
			// Static methods
			"isFinite": {
				Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
				ReturnType: types.BooleanType,
				GoCodeGen:  func(args []string) string { return fmt.Sprintf("(!math.IsInf(%s, 0) && !math.IsNaN(%s))", args[0], args[0]) },
			},
			"isNaN": {
				Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
				ReturnType: types.BooleanType,
				GoCodeGen:  func(args []string) string { return fmt.Sprintf("math.IsNaN(%s)", args[0]) },
			},
			"isInteger": {
				Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
				ReturnType: types.BooleanType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("(math.Trunc(%s) == %s && !math.IsInf(%s, 0))", args[0], args[0], args[0])
				},
			},
			"isSafeInteger": {
				Params:     []*types.Param{{Name: "x", Type: types.NumberType}},
				ReturnType: types.BooleanType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("(math.Trunc(%s) == %s && math.Abs(%s) <= 9007199254740991)", args[0], args[0], args[0])
				},
			},
			"parseFloat": {
				Params:     []*types.Param{{Name: "s", Type: types.StringType}},
				ReturnType: types.NumberType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() float64 { v, _ := strconv.ParseFloat(%s, 64); return v }()", args[0])
				},
			},
			"parseInt": {
				Params:     []*types.Param{{Name: "s", Type: types.StringType}, {Name: "radix", Type: types.IntType}},
				ReturnType: types.IntType,
				Variadic:   true, // radix is optional
				GoCodeGen: func(args []string) string {
					if len(args) == 1 {
						return fmt.Sprintf("func() int { v, _ := strconv.ParseInt(%s, 10, 64); return int(v) }()", args[0])
					}
					return fmt.Sprintf("func() int { v, _ := strconv.ParseInt(%s, %s, 64); return int(v) }()", args[0], args[1])
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------
// JSON Built-in Object
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:      "JSON",
		Imports:   []string{"encoding/json"},
		Constants: map[string]*BuiltinConstant{},
		Methods: map[string]*BuiltinMethod{
			"stringify": {
				Params:     []*types.Param{{Name: "value", Type: types.AnyType}},
				ReturnType: types.StringType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() string { b, _ := json.Marshal(%s); return string(b) }()", args[0])
				},
			},
			"parse": {
				Params:     []*types.Param{{Name: "text", Type: types.StringType}},
				ReturnType: types.AnyType, // Returns any, caller casts to expected type
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() interface{} { var v interface{}; json.Unmarshal([]byte(%s), &v); return v }()", args[0])
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Object Built-in Object
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:      "Object",
		Imports:   []string{},
		Constants: map[string]*BuiltinConstant{},
		Methods: map[string]*BuiltinMethod{
			"keys": {
				Params:     []*types.Param{{Name: "obj", Type: types.AnyType}},
				ReturnType: &types.Array{Element: types.StringType},
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() []string { keys := make([]string, 0); for k := range %s { keys = append(keys, k) }; return keys }()", args[0])
				},
			},
			"values": {
				Params:     []*types.Param{{Name: "obj", Type: types.AnyType}},
				ReturnType: types.AnyType, // Return type depends on map value type
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() interface{} { var vals []interface{}; for _, v := range %s { vals = append(vals, v) }; return vals }()", args[0])
				},
			},
			"assign": {
				Params:     []*types.Param{{Name: "target", Type: types.AnyType}, {Name: "source", Type: types.AnyType}},
				ReturnType: types.AnyType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() interface{} { for k, v := range %s { %s[k] = v }; return %s }()", args[1], args[0], args[0])
				},
			},
			"hasOwn": {
				Params:     []*types.Param{{Name: "obj", Type: types.AnyType}, {Name: "prop", Type: types.StringType}},
				ReturnType: types.BooleanType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() bool { _, ok := %s[%s]; return ok }()", args[0], args[1])
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Array Built-in Object (static methods)
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:      "Array",
		Imports:   []string{"reflect"},
		Constants: map[string]*BuiltinConstant{},
		Methods: map[string]*BuiltinMethod{
			"isArray": {
				Params:     []*types.Param{{Name: "value", Type: types.AnyType}},
				ReturnType: types.BooleanType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("(reflect.TypeOf(%s).Kind() == reflect.Slice)", args[0])
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Date Built-in Object (static methods)
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:      "Date",
		Imports:   []string{"time"},
		Constants: map[string]*BuiltinConstant{},
		Methods: map[string]*BuiltinMethod{
			"now": {
				Params:     []*types.Param{},
				ReturnType: types.NumberType,
				GoCodeGen: func(args []string) string {
					return "float64(time.Now().UnixMilli())"
				},
			},
			"parse": {
				Params:     []*types.Param{{Name: "dateString", Type: types.StringType}},
				ReturnType: types.NumberType,
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("func() float64 { t, err := time.Parse(time.RFC3339, %s); if err != nil { return 0 }; return float64(t.UnixMilli()) }()", args[0])
				},
			},
		},
	})
}

// ----------------------------------------------------------------------------
// Promise Built-in Object (static methods)
// ----------------------------------------------------------------------------

func init() {
	RegisterBuiltin(&BuiltinObject{
		Name:      "Promise",
		Imports:   []string{},
		Constants: map[string]*BuiltinConstant{},
		Methods: map[string]*BuiltinMethod{
			"resolve": {
				Params:     []*types.Param{{Name: "value", Type: types.AnyType}},
				ReturnType: &types.Promise{Value: types.AnyType},
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("GTS_Promise_Resolve[interface{}](%s)", args[0])
				},
			},
			"reject": {
				Params:     []*types.Param{{Name: "reason", Type: types.AnyType}},
				ReturnType: &types.Promise{Value: types.AnyType},
				GoCodeGen: func(args []string) string {
					return fmt.Sprintf("GTS_Promise_Reject[interface{}](fmt.Errorf(\"%%v\", %s))", args[0])
				},
			},
		},
	})
}

// DescribeBuiltin returns a description of a built-in object for documentation.
func DescribeBuiltin(objName string) string {
	obj, ok := GetBuiltin(objName)
	if !ok {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## %s\n\n", objName))

	if len(obj.Constants) > 0 {
		sb.WriteString("### Constants\n\n")
		for name, c := range obj.Constants {
			sb.WriteString(fmt.Sprintf("- `%s.%s`: %s\n", objName, name, c.Type.String()))
		}
		sb.WriteString("\n")
	}

	if len(obj.Methods) > 0 {
		sb.WriteString("### Methods\n\n")
		for name, m := range obj.Methods {
			params := make([]string, len(m.Params))
			for i, p := range m.Params {
				params[i] = fmt.Sprintf("%s: %s", p.Name, p.Type.String())
			}
			paramStr := strings.Join(params, ", ")
			if m.Variadic {
				paramStr = "..." + paramStr
			}
			sb.WriteString(fmt.Sprintf("- `%s.%s(%s)`: %s\n", objName, name, paramStr, m.ReturnType.String()))
		}
	}

	return sb.String()
}
