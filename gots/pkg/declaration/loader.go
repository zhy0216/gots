package declaration

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zhy0216/quickts/gots/pkg/ast"
)

//go:embed stdlib/*.d.gts
var stdlibFS embed.FS

// Loader loads and caches declaration files
type Loader struct {
	cache      map[string]*ModuleInfo
	cacheMutex sync.RWMutex
	searchPaths []string
}

// NewLoader creates a new declaration loader
func NewLoader() *Loader {
	return &Loader{
		cache:       make(map[string]*ModuleInfo),
		searchPaths: []string{},
	}
}

// AddSearchPath adds a path to search for declaration files
func (l *Loader) AddSearchPath(path string) {
	l.searchPaths = append(l.searchPaths, path)
}

// Load loads type information for a module
// modulePath should be like "go:strings" or "go:github.com/gin-gonic/gin"
func (l *Loader) Load(modulePath string) (*ModuleInfo, error) {
	// Check cache first
	l.cacheMutex.RLock()
	if info, ok := l.cache[modulePath]; ok {
		l.cacheMutex.RUnlock()
		return info, nil
	}
	l.cacheMutex.RUnlock()

	// Parse the module path
	var pkgPath string
	if strings.HasPrefix(modulePath, "go:") {
		pkgPath = strings.TrimPrefix(modulePath, "go:")
	} else {
		pkgPath = modulePath
	}

	// Try to find the declaration file
	content, err := l.findDeclaration(pkgPath)
	if err != nil {
		return nil, err
	}

	// Parse the declaration file
	parser := NewFromSource(content)
	file := parser.Parse()

	if len(parser.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors in declaration file: %v", parser.Errors())
	}

	// Extract module info
	modules := ExtractModuleInfo(file)

	// Cache all modules from the file
	l.cacheMutex.Lock()
	for name, info := range modules {
		l.cache["go:"+name] = info
		l.cache[name] = info
	}
	l.cacheMutex.Unlock()

	// Return the requested module
	if info, ok := modules[pkgPath]; ok {
		return info, nil
	}

	return nil, fmt.Errorf("module %s not found in declaration file", modulePath)
}

// findDeclaration searches for a declaration file
func (l *Loader) findDeclaration(pkgPath string) (string, error) {
	// Convert package path to filename
	// e.g., "strings" -> "go_strings.d.gts"
	// e.g., "github.com/gin-gonic/gin" -> "go_github_com_gin_gonic_gin.d.gts"
	filename := "go_" + strings.ReplaceAll(pkgPath, "/", "_") + ".d.gts"
	filename = strings.ReplaceAll(filename, ".", "_") + ".d.gts"
	filename = strings.TrimSuffix(filename, ".d.gts.d.gts") + ".d.gts"

	// Simpler filename for stdlib
	simpleFilename := "go_" + strings.ReplaceAll(pkgPath, "/", "_") + ".d.gts"

	// 1. Try embedded stdlib
	content, err := stdlibFS.ReadFile("stdlib/" + simpleFilename)
	if err == nil {
		return string(content), nil
	}

	// 2. Try search paths
	for _, searchPath := range l.searchPaths {
		fullPath := filepath.Join(searchPath, simpleFilename)
		content, err := os.ReadFile(fullPath)
		if err == nil {
			return string(content), nil
		}
	}

	// 3. Try current directory's declarations folder
	content, err = os.ReadFile(filepath.Join("declarations", simpleFilename))
	if err == nil {
		return string(content), nil
	}

	// 4. Try home directory
	home, _ := os.UserHomeDir()
	if home != "" {
		content, err = os.ReadFile(filepath.Join(home, ".gots", "declarations", simpleFilename))
		if err == nil {
			return string(content), nil
		}
	}

	return "", fmt.Errorf("declaration file not found for package: %s (looked for %s)", pkgPath, simpleFilename)
}

// GetFunction returns function info for a Go package function
func (l *Loader) GetFunction(pkgPath, funcName string) (*FunctionInfo, error) {
	info, err := l.Load(pkgPath)
	if err != nil {
		return nil, err
	}

	if fn, ok := info.Functions[funcName]; ok {
		return fn, nil
	}

	return nil, fmt.Errorf("function %s not found in package %s", funcName, pkgPath)
}

// GetConstant returns constant type for a Go package constant
func (l *Loader) GetConstant(pkgPath, constName string) (ast.Type, error) {
	info, err := l.Load(pkgPath)
	if err != nil {
		return nil, err
	}

	if t, ok := info.Constants[constName]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("constant %s not found in package %s", constName, pkgPath)
}

// GetType returns type info for a Go package type
func (l *Loader) GetType(pkgPath, typeName string) (ast.Type, error) {
	info, err := l.Load(pkgPath)
	if err != nil {
		return nil, err
	}

	if t, ok := info.Types[typeName]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("type %s not found in package %s", typeName, pkgPath)
}

// GetInterface returns interface info for a Go package interface
func (l *Loader) GetInterface(pkgPath, ifaceName string) (*InterfaceInfo, error) {
	info, err := l.Load(pkgPath)
	if err != nil {
		return nil, err
	}

	if iface, ok := info.Interfaces[ifaceName]; ok {
		return iface, nil
	}

	return nil, fmt.Errorf("interface %s not found in package %s", ifaceName, pkgPath)
}

// GetClass returns class info for a Go package struct/class
func (l *Loader) GetClass(pkgPath, className string) (*ClassInfo, error) {
	info, err := l.Load(pkgPath)
	if err != nil {
		return nil, err
	}

	if class, ok := info.Classes[className]; ok {
		return class, nil
	}

	return nil, fmt.Errorf("class %s not found in package %s", className, pkgPath)
}

// HasPackage checks if a package has declarations available
func (l *Loader) HasPackage(pkgPath string) bool {
	_, err := l.Load(pkgPath)
	return err == nil
}

// DefaultLoader is the global declaration loader
var DefaultLoader = NewLoader()

// LoadModule loads a module using the default loader
func LoadModule(modulePath string) (*ModuleInfo, error) {
	return DefaultLoader.Load(modulePath)
}
