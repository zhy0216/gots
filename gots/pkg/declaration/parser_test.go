package declaration

import (
	"testing"
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
