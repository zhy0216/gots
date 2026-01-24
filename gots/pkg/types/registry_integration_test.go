package types

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/declaration"
)

func TestGetGoPackageFunctionFromDeclaration(t *testing.T) {
	// Create a temporary directory with a custom declaration
	tmpDir, err := os.MkdirTemp("", "gots-registry-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a test declaration file
	declContent := `
declare module "go:customtest" {
    function CustomFunc(s: string, n: int): boolean
    function AnotherFunc(): float
    const CustomConst: int
}
`
	declPath := filepath.Join(tmpDir, "go_customtest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	// Add search path to default loader
	declaration.DefaultLoader.AddSearchPath(tmpDir)

	// Test loading function from declaration
	fnType := GetGoPackageFunction("customtest", "CustomFunc")
	if fnType == nil {
		t.Fatal("Expected to get CustomFunc from declaration")
	}

	fn, ok := fnType.(*Function)
	if !ok {
		t.Fatal("Expected Function type")
	}

	if len(fn.Params) != 2 {
		t.Errorf("Expected 2 params, got %d", len(fn.Params))
	}

	// Check return type is boolean
	if !fn.ReturnType.Equals(BooleanType) {
		t.Errorf("Expected boolean return type, got %s", fn.ReturnType.String())
	}
}

func TestGetGoPackageConstantFromDeclaration(t *testing.T) {
	// This test uses the declaration from TestGetGoPackageFunctionFromDeclaration
	// which should still be cached

	constType := GetGoPackageConstant("customtest", "CustomConst")
	if constType == nil {
		t.Fatal("Expected to get CustomConst from declaration")
	}

	if !constType.Equals(IntType) {
		t.Errorf("Expected int type, got %s", constType.String())
	}
}

func TestGetGoPackageFunctionFallbackToRegistry(t *testing.T) {
	// Test that hardcoded registry is used when no declaration is found
	// The "fmt" package should be in the hardcoded registry
	fnType := GetGoPackageFunction("fmt", "Println")
	if fnType == nil {
		t.Fatal("Expected to get Println from registry")
	}

	fn, ok := fnType.(*Function)
	if !ok {
		t.Fatal("Expected Function type")
	}

	if fn.ReturnType == nil || !fn.ReturnType.Equals(VoidType) {
		t.Error("Expected void return type for Println")
	}
}

func TestGetGoPackageFunctionNotFound(t *testing.T) {
	// Test that nil is returned for non-existent functions
	fnType := GetGoPackageFunction("fmt", "NonExistentFunction12345")
	if fnType != nil {
		t.Error("Expected nil for non-existent function")
	}
}

func TestConvertAstTypeToTypePrimitives(t *testing.T) {
	tests := []struct {
		name     string
		declSrc  string
		expected Type
	}{
		{
			name:     "int return",
			declSrc:  `declare module "go:test1" { function Foo(): int }`,
			expected: IntType,
		},
		{
			name:     "float return",
			declSrc:  `declare module "go:test2" { function Foo(): float }`,
			expected: FloatType,
		},
		{
			name:     "string return",
			declSrc:  `declare module "go:test3" { function Foo(): string }`,
			expected: StringType,
		},
		{
			name:     "boolean return",
			declSrc:  `declare module "go:test4" { function Foo(): boolean }`,
			expected: BooleanType,
		},
		{
			name:     "void return",
			declSrc:  `declare module "go:test5" { function Foo(): void }`,
			expected: VoidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file for this test
			tmpDir, err := os.MkdirTemp("", "gots-type-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Extract module name from declaration
			parser := declaration.NewFromSource(tt.declSrc)
			file := parser.Parse()
			if len(file.Modules) == 0 {
				t.Fatal("No modules parsed")
			}
			modName := file.Modules[0].Name
			pkgName := modName[3:] // strip "go:"

			// Write file
			declPath := filepath.Join(tmpDir, "go_"+pkgName+".d.gts")
			if err := os.WriteFile(declPath, []byte(tt.declSrc), 0644); err != nil {
				t.Fatalf("Failed to write declaration file: %v", err)
			}

			// Create new loader for this test
			loader := declaration.NewLoader()
			loader.AddSearchPath(tmpDir)

			// Load and convert
			fn, err := loader.GetFunction(modName, "Foo")
			if err != nil {
				t.Fatalf("Failed to get function: %v", err)
			}

			converted := convertAstTypeToType(fn.ReturnType)
			if !converted.Equals(tt.expected) {
				t.Errorf("Expected %s, got %s", tt.expected.String(), converted.String())
			}
		})
	}
}

func TestConvertAstTypeToTypeArrays(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-array-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:arraytest" {
    function GetInts(): int[]
    function GetStrings(): string[]
}
`
	declPath := filepath.Join(tmpDir, "go_arraytest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := declaration.NewLoader()
	loader.AddSearchPath(tmpDir)

	// Test int[]
	fn, err := loader.GetFunction("go:arraytest", "GetInts")
	if err != nil {
		t.Fatalf("Failed to get GetInts: %v", err)
	}

	converted := convertAstTypeToType(fn.ReturnType)
	arr, ok := converted.(*Array)
	if !ok {
		t.Fatalf("Expected Array type, got %T", converted)
	}

	if !arr.Element.Equals(IntType) {
		t.Errorf("Expected int element type, got %s", arr.Element.String())
	}

	// Test string[]
	fn, err = loader.GetFunction("go:arraytest", "GetStrings")
	if err != nil {
		t.Fatalf("Failed to get GetStrings: %v", err)
	}

	converted = convertAstTypeToType(fn.ReturnType)
	arr, ok = converted.(*Array)
	if !ok {
		t.Fatalf("Expected Array type, got %T", converted)
	}

	if !arr.Element.Equals(StringType) {
		t.Errorf("Expected string element type, got %s", arr.Element.String())
	}
}

func TestConvertAstTypeToTypeNullable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-nullable-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:nullabletest" {
    function MaybeString(): string | null
    function MaybeInt(): int | null
}
`
	declPath := filepath.Join(tmpDir, "go_nullabletest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := declaration.NewLoader()
	loader.AddSearchPath(tmpDir)

	// Test string | null
	fn, err := loader.GetFunction("go:nullabletest", "MaybeString")
	if err != nil {
		t.Fatalf("Failed to get MaybeString: %v", err)
	}

	converted := convertAstTypeToType(fn.ReturnType)
	nullable, ok := converted.(*Nullable)
	if !ok {
		t.Fatalf("Expected Nullable type, got %T", converted)
	}

	if !nullable.Inner.Equals(StringType) {
		t.Errorf("Expected string inner type, got %s", nullable.Inner.String())
	}
}

func TestConvertAstTypeToTypeFunction(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-functype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:functypetest" {
    type Handler = (s: string, n: int) => void
}
`
	declPath := filepath.Join(tmpDir, "go_functypetest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := declaration.NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get the type
	handlerType, err := loader.GetType("go:functypetest", "Handler")
	if err != nil {
		t.Fatalf("Failed to get Handler type: %v", err)
	}

	converted := convertAstTypeToType(handlerType)
	fn, ok := converted.(*Function)
	if !ok {
		t.Fatalf("Expected Function type, got %T", converted)
	}

	if len(fn.Params) != 2 {
		t.Errorf("Expected 2 params, got %d", len(fn.Params))
	}

	if !fn.ReturnType.Equals(VoidType) {
		t.Errorf("Expected void return type, got %s", fn.ReturnType.String())
	}
}

func TestConvertAstTypeToTypeObject(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-objtype-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:objtypetest" {
    type Point = {
        x: int
        y: int
    }
}
`
	declPath := filepath.Join(tmpDir, "go_objtypetest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := declaration.NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get the type
	pointType, err := loader.GetType("go:objtypetest", "Point")
	if err != nil {
		t.Fatalf("Failed to get Point type: %v", err)
	}

	converted := convertAstTypeToType(pointType)
	obj, ok := converted.(*Object)
	if !ok {
		t.Fatalf("Expected Object type, got %T", converted)
	}

	if len(obj.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(obj.Properties))
	}

	xProp := obj.GetProperty("x")
	if xProp == nil {
		t.Error("Expected 'x' property")
	} else if !xProp.Type.Equals(IntType) {
		t.Errorf("Expected x to be int, got %s", xProp.Type.String())
	}

	yProp := obj.GetProperty("y")
	if yProp == nil {
		t.Error("Expected 'y' property")
	} else if !yProp.Type.Equals(IntType) {
		t.Errorf("Expected y to be int, got %s", yProp.Type.String())
	}
}

func TestConvertDeclFunctionToType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-declfn-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:declfntest" {
    function Process(data: string, count: int, flag: boolean): string[]
}
`
	declPath := filepath.Join(tmpDir, "go_declfntest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := declaration.NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get the function info
	fnInfo, err := loader.GetFunction("go:declfntest", "Process")
	if err != nil {
		t.Fatalf("Failed to get Process function: %v", err)
	}

	// Convert to types.Type
	converted := convertDeclFunctionToType(fnInfo)
	fn, ok := converted.(*Function)
	if !ok {
		t.Fatalf("Expected Function type, got %T", converted)
	}

	// Check params
	if len(fn.Params) != 3 {
		t.Fatalf("Expected 3 params, got %d", len(fn.Params))
	}

	expectedParamTypes := []Type{StringType, IntType, BooleanType}
	for i, expected := range expectedParamTypes {
		if !fn.Params[i].Type.Equals(expected) {
			t.Errorf("Param %d: expected %s, got %s", i, expected.String(), fn.Params[i].Type.String())
		}
	}

	// Check return type is string[]
	arr, ok := fn.ReturnType.(*Array)
	if !ok {
		t.Fatalf("Expected Array return type, got %T", fn.ReturnType)
	}
	if !arr.Element.Equals(StringType) {
		t.Errorf("Expected string[] return type, got %s", fn.ReturnType.String())
	}
}
