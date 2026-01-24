package declaration

import (
	"testing"

	"github.com/zhy0216/quickts/gots/pkg/ast"
)

func TestParseDeclareModule(t *testing.T) {
	source := `
declare module "go:strings" {
    function ToUpper(s: string): string
    function Split(s: string, sep: string): string[]
    function Contains(s: string, substr: string): boolean
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(file.Modules) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(file.Modules))
	}

	mod := file.Modules[0]
	if mod.Name != "go:strings" {
		t.Errorf("Expected module name 'go:strings', got '%s'", mod.Name)
	}

	if len(mod.Members) != 3 {
		t.Fatalf("Expected 3 members, got %d", len(mod.Members))
	}
}

func TestParseDeclareFunction(t *testing.T) {
	source := `
declare module "go:fmt" {
    function Println(...args: any[]): void
    function Printf(format: string, ...args: any[]): void
    function Sprintf(format: string, ...args: any[]): string
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	fmtMod, ok := modules["fmt"]
	if !ok {
		t.Fatal("Expected fmt module")
	}

	// Test Println
	println, ok := fmtMod.Functions["Println"]
	if !ok {
		t.Fatal("Expected Println function")
	}
	if !println.Variadic {
		t.Error("Expected Println to be variadic")
	}

	// Test Printf
	printf, ok := fmtMod.Functions["Printf"]
	if !ok {
		t.Fatal("Expected Printf function")
	}
	if len(printf.Params) != 2 {
		t.Errorf("Expected Printf to have 2 params, got %d", len(printf.Params))
	}

	// Test Sprintf
	sprintf, ok := fmtMod.Functions["Sprintf"]
	if !ok {
		t.Fatal("Expected Sprintf function")
	}
	if sprintf.ReturnType == nil {
		t.Error("Expected Sprintf to have return type")
	}
}

func TestParseDeclareInterface(t *testing.T) {
	source := `
declare module "go:io" {
    interface Reader {
        Read(p: byte[]): int
    }
    interface Writer {
        Write(p: byte[]): int
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	ioMod, ok := modules["io"]
	if !ok {
		t.Fatal("Expected io module")
	}

	// Test Reader interface
	reader, ok := ioMod.Interfaces["Reader"]
	if !ok {
		t.Fatal("Expected Reader interface")
	}
	if len(reader.Methods) != 1 {
		t.Errorf("Expected Reader to have 1 method, got %d", len(reader.Methods))
	}

	// Test Writer interface
	writer, ok := ioMod.Interfaces["Writer"]
	if !ok {
		t.Fatal("Expected Writer interface")
	}
	if len(writer.Methods) != 1 {
		t.Errorf("Expected Writer to have 1 method, got %d", len(writer.Methods))
	}
}

func TestParseDeclareType(t *testing.T) {
	source := `
declare module "go:builtin" {
    type Error = {
        Error(): string
    }
    type Duration = int
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod, ok := modules["builtin"]
	if !ok {
		t.Fatal("Expected builtin module")
	}

	if len(mod.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(mod.Types))
	}

	if _, ok := mod.Types["Error"]; !ok {
		t.Error("Expected Error type")
	}
	if _, ok := mod.Types["Duration"]; !ok {
		t.Error("Expected Duration type")
	}
}

func TestParseDeclareConst(t *testing.T) {
	source := `
declare module "go:math" {
    const Pi: float
    const E: float
    const MaxInt: int
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mathMod, ok := modules["math"]
	if !ok {
		t.Fatal("Expected math module")
	}

	if len(mathMod.Constants) != 3 {
		t.Errorf("Expected 3 constants, got %d", len(mathMod.Constants))
	}

	if _, ok := mathMod.Constants["Pi"]; !ok {
		t.Error("Expected Pi constant")
	}
	if _, ok := mathMod.Constants["E"]; !ok {
		t.Error("Expected E constant")
	}
}

func TestParseDeclareClass(t *testing.T) {
	source := `
declare module "go:regexp" {
    class Regexp {
        MatchString(s: string): boolean
        FindString(s: string): string
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod, ok := modules["regexp"]
	if !ok {
		t.Fatal("Expected regexp module")
	}

	regexp, ok := mod.Classes["Regexp"]
	if !ok {
		t.Fatal("Expected Regexp class")
	}

	if len(regexp.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(regexp.Methods))
	}
}

func TestParseNullableType(t *testing.T) {
	source := `
declare module "go:test" {
    function MaybeError(): Error | null
    function GetValue(): string | null
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(file.Modules) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(file.Modules))
	}
}

func TestParseArrayType(t *testing.T) {
	source := `
declare module "go:test" {
    function GetStrings(): string[]
    function GetMatrix(): int[][]
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(file.Modules) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(file.Modules))
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	getStrings, ok := mod.Functions["GetStrings"]
	if !ok {
		t.Fatal("Expected GetStrings function")
	}
	if getStrings.ReturnType == nil {
		t.Error("Expected return type for GetStrings")
	}
}

func TestParseMultipleModules(t *testing.T) {
	source := `
declare module "go:strings" {
    function ToUpper(s: string): string
}

declare module "go:strconv" {
    function Itoa(i: int): string
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(file.Modules) != 2 {
		t.Fatalf("Expected 2 modules, got %d", len(file.Modules))
	}

	modules := ExtractModuleInfo(file)
	if _, ok := modules["strings"]; !ok {
		t.Error("Expected strings module")
	}
	if _, ok := modules["strconv"]; !ok {
		t.Error("Expected strconv module")
	}
}

// ============================================================================
// Additional Unit Tests
// ============================================================================

func TestParseFunctionWithNoParams(t *testing.T) {
	source := `
declare module "go:test" {
    function GetTime(): int
    function Now(): Time
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	getTime, ok := mod.Functions["GetTime"]
	if !ok {
		t.Fatal("Expected GetTime function")
	}
	if len(getTime.Params) != 0 {
		t.Errorf("Expected 0 params, got %d", len(getTime.Params))
	}
}

func TestParseFunctionWithMultipleParams(t *testing.T) {
	source := `
declare module "go:test" {
    function Replace(s: string, old: string, replacement: string, n: int): string
    function Copy(dst: byte[], src: byte[]): int
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	replace, ok := mod.Functions["Replace"]
	if !ok {
		t.Fatal("Expected Replace function")
	}
	if len(replace.Params) != 4 {
		t.Errorf("Expected 4 params, got %d", len(replace.Params))
	}

	// Verify param names
	expectedParams := []string{"s", "old", "replacement", "n"}
	for i, name := range expectedParams {
		if replace.Params[i].Name != name {
			t.Errorf("Expected param %d to be '%s', got '%s'", i, name, replace.Params[i].Name)
		}
	}
}

func TestParseGenericTypes(t *testing.T) {
	source := `
declare module "go:test" {
    function Make(): StringIntMap
    function GetList(): UserList
    function Complex(): NestedMap
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	if len(mod.Functions) != 3 {
		t.Errorf("Expected 3 functions, got %d", len(mod.Functions))
	}
}

func TestParseClassWithFields(t *testing.T) {
	source := `
declare module "go:test" {
    class User {
        name: string
        age: int
        email: string | null
        GetName(): string
        SetAge(age: int): void
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	user, ok := mod.Classes["User"]
	if !ok {
		t.Fatal("Expected User class")
	}

	if len(user.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(user.Fields))
	}

	if len(user.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(user.Methods))
	}

	// Check field names
	fieldNames := make(map[string]bool)
	for _, f := range user.Fields {
		fieldNames[f.Name] = true
	}
	if !fieldNames["name"] {
		t.Error("Expected 'name' field")
	}
	if !fieldNames["age"] {
		t.Error("Expected 'age' field")
	}
}

func TestParseClassWithInheritance(t *testing.T) {
	source := `
declare module "go:test" {
    class Animal {
        name: string
        Speak(): void
    }
    class Dog extends Animal {
        breed: string
        Bark(): void
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	dog, ok := mod.Classes["Dog"]
	if !ok {
		t.Fatal("Expected Dog class")
	}

	if dog.SuperClass != "Animal" {
		t.Errorf("Expected Dog to extend Animal, got '%s'", dog.SuperClass)
	}
}

func TestParseInterfaceWithMultipleMethods(t *testing.T) {
	source := `
declare module "go:test" {
    interface ReadWriter {
        Read(p: byte[]): int
        Write(p: byte[]): int
        Close(): Error | null
        Seek(offset: int, whence: int): int
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	rw, ok := mod.Interfaces["ReadWriter"]
	if !ok {
		t.Fatal("Expected ReadWriter interface")
	}

	if len(rw.Methods) != 4 {
		t.Errorf("Expected 4 methods, got %d", len(rw.Methods))
	}
}

func TestParseObjectType(t *testing.T) {
	source := `
declare module "go:test" {
    type Point = {
        x: int
        y: int
    }
    type Config = {
        host: string
        port: int
        timeout: float
    }
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	if len(mod.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(mod.Types))
	}

	pointType, ok := mod.Types["Point"]
	if !ok {
		t.Fatal("Expected Point type")
	}

	objType, ok := pointType.(*ast.ObjectType)
	if !ok {
		t.Fatal("Expected Point to be an ObjectType")
	}

	if len(objType.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(objType.Properties))
	}
}

func TestParseFunctionType(t *testing.T) {
	source := `
declare module "go:test" {
    type Handler = (req: Request, res: Response) => void
    type Callback = (err: Error | null, data: any) => void
    type Predicate = (value: int) => boolean
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	if len(mod.Types) != 3 {
		t.Errorf("Expected 3 types, got %d", len(mod.Types))
	}
}

func TestParseTupleType(t *testing.T) {
	source := `
declare module "go:test" {
    function Divide(a: int, b: int): (int, Error | null)
    function GetPair(): (string, int)
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	divide, ok := mod.Functions["Divide"]
	if !ok {
		t.Fatal("Expected Divide function")
	}

	// Check that return type is a tuple
	_, isTuple := divide.ReturnType.(*ast.ReturnTupleType)
	if !isTuple {
		t.Error("Expected Divide return type to be a tuple")
	}
}

func TestParseNestedArrayTypes(t *testing.T) {
	source := `
declare module "go:test" {
    function GetMatrix(): int[][]
    function Get3D(): string[][][]
    function GetNullableArray(): (int | null)[]
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	getMatrix, ok := mod.Functions["GetMatrix"]
	if !ok {
		t.Fatal("Expected GetMatrix function")
	}

	// Check nested array type
	arrType, ok := getMatrix.ReturnType.(*ast.ArrayType)
	if !ok {
		t.Fatal("Expected GetMatrix return type to be an array")
	}

	innerArr, ok := arrType.ElementType.(*ast.ArrayType)
	if !ok {
		t.Fatal("Expected inner type to be an array (int[][])")
	}

	if _, ok := innerArr.ElementType.(*ast.PrimitiveType); !ok {
		t.Error("Expected innermost type to be int")
	}
}

func TestParseEmptyModule(t *testing.T) {
	source := `
declare module "go:empty" {
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	if len(file.Modules) != 1 {
		t.Fatalf("Expected 1 module, got %d", len(file.Modules))
	}

	if len(file.Modules[0].Members) != 0 {
		t.Errorf("Expected 0 members, got %d", len(file.Modules[0].Members))
	}
}

func TestParseByteType(t *testing.T) {
	source := `
declare module "go:test" {
    function ReadByte(): byte
    function WriteBytes(data: byte[]): int
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	readByte, ok := mod.Functions["ReadByte"]
	if !ok {
		t.Fatal("Expected ReadByte function")
	}

	_, isByte := readByte.ReturnType.(*ast.ByteType)
	if !isByte {
		t.Error("Expected ReadByte return type to be byte")
	}
}

func TestParseAnyType(t *testing.T) {
	source := `
declare module "go:test" {
    function GetAny(): any
    function SetAny(value: any): void
    function Process(data: any[]): any
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	getAny, ok := mod.Functions["GetAny"]
	if !ok {
		t.Fatal("Expected GetAny function")
	}

	_, isAny := getAny.ReturnType.(*ast.AnyType)
	if !isAny {
		t.Error("Expected GetAny return type to be any")
	}
}

func TestParseAllPrimitiveTypes(t *testing.T) {
	source := `
declare module "go:test" {
    function GetInt(): int
    function GetFloat(): float
    function GetString(): string
    function GetBoolean(): boolean
    function GetVoid(): void
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	mod := modules["test"]

	tests := []struct {
		name     string
		expected ast.PrimitiveKind
	}{
		{"GetInt", ast.TypeInt},
		{"GetFloat", ast.TypeFloat},
		{"GetString", ast.TypeString},
		{"GetBoolean", ast.TypeBoolean},
		{"GetVoid", ast.TypeVoid},
	}

	for _, tt := range tests {
		fn, ok := mod.Functions[tt.name]
		if !ok {
			t.Errorf("Expected %s function", tt.name)
			continue
		}

		prim, ok := fn.ReturnType.(*ast.PrimitiveType)
		if !ok {
			t.Errorf("Expected %s return type to be primitive", tt.name)
			continue
		}

		if prim.Kind != tt.expected {
			t.Errorf("Expected %s return type to be %v, got %v", tt.name, tt.expected, prim.Kind)
		}
	}
}

func TestDeclarationFileString(t *testing.T) {
	source := `
declare module "go:test" {
    function Hello(): string
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	str := file.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
}

func TestDeclareFunctionString(t *testing.T) {
	source := `
declare module "go:test" {
    function Add(a: int, b: int): int
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	fn := file.Modules[0].Members[0].(*ast.DeclareFunction)
	str := fn.String()

	if str != "function Add(a: int, b: int): int" {
		t.Errorf("Unexpected string: %s", str)
	}
}

func TestExtractModuleInfoStripsGoPrefix(t *testing.T) {
	source := `
declare module "go:strings" {
    function ToUpper(s: string): string
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	modules := ExtractModuleInfo(file)

	// Should be accessible without "go:" prefix
	if _, ok := modules["strings"]; !ok {
		t.Error("Expected module to be accessible as 'strings'")
	}
}

func TestParseComplexRealWorldDeclaration(t *testing.T) {
	source := `
declare module "go:net/http" {
    type Header = {
        Get(key: string): string
        Set(key: string, value: string): void
    }

    interface ResponseWriter {
        Header(): Header
        Write(data: byte[]): int
        WriteHeader(statusCode: int): void
    }

    interface Request {
        GetMethod(): string
        GetURL(): URL
        GetHeader(): Header
        GetBody(): Reader
    }

    type HandlerFunc = (w: ResponseWriter, r: Request) => void

    class Server {
        ListenAndServe(): Error | null
        Shutdown(ctx: Context): Error | null
    }

    function ListenAndServe(addr: string, handler: Handler): Error | null
    function Get(url: string): Response
    function Post(url: string, contentType: string, body: Reader): Response
}
`
	parser := NewFromSource(source)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		t.Fatalf("Parser errors: %v", parser.Errors())
	}

	modules := ExtractModuleInfo(file)
	httpMod, ok := modules["net/http"]
	if !ok {
		t.Fatal("Expected net/http module")
	}

	// Check types
	if len(httpMod.Types) < 1 {
		t.Error("Expected at least 1 type (Header)")
	}

	// Check interfaces
	if len(httpMod.Interfaces) < 2 {
		t.Error("Expected at least 2 interfaces")
	}

	// Check classes
	if len(httpMod.Classes) < 1 {
		t.Error("Expected at least 1 class")
	}

	// Check functions
	if len(httpMod.Functions) < 3 {
		t.Error("Expected at least 3 functions")
	}
}
