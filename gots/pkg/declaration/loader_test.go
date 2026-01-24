package declaration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderNewLoader(t *testing.T) {
	loader := NewLoader()
	if loader == nil {
		t.Fatal("Expected non-nil loader")
	}
	if loader.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestLoaderAddSearchPath(t *testing.T) {
	loader := NewLoader()
	loader.AddSearchPath("/some/path")
	loader.AddSearchPath("/another/path")

	if len(loader.searchPaths) != 2 {
		t.Errorf("Expected 2 search paths, got %d", len(loader.searchPaths))
	}
}

func TestLoaderLoadFromCustomDeclaration(t *testing.T) {
	// Create a temporary directory with a custom declaration
	tmpDir, err := os.MkdirTemp("", "gots-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a test declaration file
	declContent := `
declare module "go:mytest" {
    function Hello(name: string): string
    function Add(a: int, b: int): int
    const Version: string
}
`
	declPath := filepath.Join(tmpDir, "go_mytest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	// Create loader and add search path
	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Load the module
	info, err := loader.Load("go:mytest")
	if err != nil {
		t.Fatalf("Failed to load module: %v", err)
	}

	// Verify the loaded info
	if len(info.Functions) != 2 {
		t.Errorf("Expected 2 functions, got %d", len(info.Functions))
	}

	if _, ok := info.Functions["Hello"]; !ok {
		t.Error("Expected Hello function")
	}
	if _, ok := info.Functions["Add"]; !ok {
		t.Error("Expected Add function")
	}

	if _, ok := info.Constants["Version"]; !ok {
		t.Error("Expected Version constant")
	}
}

func TestLoaderCaching(t *testing.T) {
	// Create a temporary directory with a declaration
	tmpDir, err := os.MkdirTemp("", "gots-test-cache")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:cachetest" {
    function Foo(): int
}
`
	declPath := filepath.Join(tmpDir, "go_cachetest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Load twice
	info1, err := loader.Load("go:cachetest")
	if err != nil {
		t.Fatalf("First load failed: %v", err)
	}

	info2, err := loader.Load("go:cachetest")
	if err != nil {
		t.Fatalf("Second load failed: %v", err)
	}

	// Should be the same cached instance
	if info1 != info2 {
		t.Error("Expected cached instance to be returned")
	}
}

func TestLoaderGetFunction(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-fn")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:functest" {
    function Greet(name: string): string
    function Calculate(a: int, b: int): int
}
`
	declPath := filepath.Join(tmpDir, "go_functest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get existing function
	fn, err := loader.GetFunction("go:functest", "Greet")
	if err != nil {
		t.Fatalf("Failed to get function: %v", err)
	}
	if fn.Name != "Greet" {
		t.Errorf("Expected function name 'Greet', got '%s'", fn.Name)
	}

	// Try to get non-existent function
	_, err = loader.GetFunction("go:functest", "NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent function")
	}
}

func TestLoaderGetConstant(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-const")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:consttest" {
    const Pi: float
    const MaxSize: int
}
`
	declPath := filepath.Join(tmpDir, "go_consttest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get existing constant
	constType, err := loader.GetConstant("go:consttest", "Pi")
	if err != nil {
		t.Fatalf("Failed to get constant: %v", err)
	}
	if constType == nil {
		t.Error("Expected non-nil constant type")
	}

	// Try to get non-existent constant
	_, err = loader.GetConstant("go:consttest", "NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent constant")
	}
}

func TestLoaderGetInterface(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-iface")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:ifacetest" {
    interface Reader {
        Read(p: byte[]): int
    }
    interface Writer {
        Write(p: byte[]): int
    }
}
`
	declPath := filepath.Join(tmpDir, "go_ifacetest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get existing interface
	iface, err := loader.GetInterface("go:ifacetest", "Reader")
	if err != nil {
		t.Fatalf("Failed to get interface: %v", err)
	}
	if iface.Name != "Reader" {
		t.Errorf("Expected interface name 'Reader', got '%s'", iface.Name)
	}
	if len(iface.Methods) != 1 {
		t.Errorf("Expected 1 method, got %d", len(iface.Methods))
	}

	// Try to get non-existent interface
	_, err = loader.GetInterface("go:ifacetest", "NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent interface")
	}
}

func TestLoaderGetClass(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-class")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:classtest" {
    class Buffer {
        data: byte[]
        Write(p: byte[]): int
        Read(p: byte[]): int
    }
}
`
	declPath := filepath.Join(tmpDir, "go_classtest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get existing class
	class, err := loader.GetClass("go:classtest", "Buffer")
	if err != nil {
		t.Fatalf("Failed to get class: %v", err)
	}
	if class.Name != "Buffer" {
		t.Errorf("Expected class name 'Buffer', got '%s'", class.Name)
	}
	if len(class.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(class.Fields))
	}
	if len(class.Methods) != 2 {
		t.Errorf("Expected 2 methods, got %d", len(class.Methods))
	}

	// Try to get non-existent class
	_, err = loader.GetClass("go:classtest", "NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent class")
	}
}

func TestLoaderGetType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-type")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:typetest" {
    type Point = {
        x: int
        y: int
    }
    type Handler = (data: any) => void
}
`
	declPath := filepath.Join(tmpDir, "go_typetest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Get existing type
	pointType, err := loader.GetType("go:typetest", "Point")
	if err != nil {
		t.Fatalf("Failed to get type: %v", err)
	}
	if pointType == nil {
		t.Error("Expected non-nil type")
	}

	// Try to get non-existent type
	_, err = loader.GetType("go:typetest", "NonExistent")
	if err == nil {
		t.Error("Expected error for non-existent type")
	}
}

func TestLoaderHasPackage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-has")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:hastest" {
    function Foo(): int
}
`
	declPath := filepath.Join(tmpDir, "go_hastest.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Should have the package
	if !loader.HasPackage("go:hastest") {
		t.Error("Expected HasPackage to return true for existing package")
	}

	// Should not have non-existent package
	if loader.HasPackage("go:nonexistent") {
		t.Error("Expected HasPackage to return false for non-existent package")
	}
}

func TestLoaderLoadModuleNotFound(t *testing.T) {
	loader := NewLoader()

	_, err := loader.Load("go:definitely-does-not-exist-12345")
	if err == nil {
		t.Error("Expected error for non-existent module")
	}
}

func TestLoaderLoadWithoutGoPrefix(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-noprefix")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:noprefix" {
    function Bar(): string
}
`
	declPath := filepath.Join(tmpDir, "go_noprefix.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Load with prefix
	info1, err := loader.Load("go:noprefix")
	if err != nil {
		t.Fatalf("Failed to load with prefix: %v", err)
	}

	// Load without prefix should also work (from cache)
	info2, err := loader.Load("noprefix")
	if err != nil {
		t.Fatalf("Failed to load without prefix: %v", err)
	}

	if info1 != info2 {
		t.Error("Expected same cached instance")
	}
}

func TestLoaderMultipleModulesInOneFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-multi")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Single file with multiple modules (unusual but valid)
	declContent := `
declare module "go:multi1" {
    function Foo(): int
}

declare module "go:multi2" {
    function Bar(): string
}
`
	declPath := filepath.Join(tmpDir, "go_multi1.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Load first module
	info1, err := loader.Load("go:multi1")
	if err != nil {
		t.Fatalf("Failed to load multi1: %v", err)
	}
	if _, ok := info1.Functions["Foo"]; !ok {
		t.Error("Expected Foo function in multi1")
	}

	// Second module should also be cached
	info2, err := loader.Load("go:multi2")
	if err != nil {
		t.Fatalf("Failed to load multi2: %v", err)
	}
	if _, ok := info2.Functions["Bar"]; !ok {
		t.Error("Expected Bar function in multi2")
	}
}

func TestLoaderParseError(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-parseerr")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Invalid declaration syntax
	declContent := `
declare module "go:invalid" {
    function Broken(
}
`
	declPath := filepath.Join(tmpDir, "go_invalid.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Should fail to parse
	_, err = loader.Load("go:invalid")
	if err == nil {
		t.Error("Expected error for invalid declaration syntax")
	}
}

func TestDefaultLoader(t *testing.T) {
	// Test that default loader exists and works
	if DefaultLoader == nil {
		t.Fatal("DefaultLoader should not be nil")
	}

	// LoadModule should use the default loader
	// This will likely fail since we don't have stdlib files embedded yet
	// but it should at least not panic
	_, _ = LoadModule("go:nonexistent")
}

func TestLoaderConcurrentAccess(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gots-test-concurrent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	declContent := `
declare module "go:concurrent" {
    function Test(): int
}
`
	declPath := filepath.Join(tmpDir, "go_concurrent.d.gts")
	if err := os.WriteFile(declPath, []byte(declContent), 0644); err != nil {
		t.Fatalf("Failed to write declaration file: %v", err)
	}

	loader := NewLoader()
	loader.AddSearchPath(tmpDir)

	// Run concurrent loads
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := loader.Load("go:concurrent")
			if err != nil {
				t.Errorf("Concurrent load failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
