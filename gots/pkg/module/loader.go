// Package module provides module loading and resolution for goTS.
package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/ast"
	"github.com/zhy0216/quickts/gots/pkg/lexer"
	"github.com/zhy0216/quickts/gots/pkg/parser"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Module represents a loaded goTS module.
type Module struct {
	Path    string                 // Absolute file path
	Program *ast.Program           // Parsed AST
	Exports map[string]ExportInfo  // Exported symbols
}

// ExportInfo contains information about an exported symbol.
type ExportInfo struct {
	Name string      // Symbol name
	Type types.Type  // Type of the symbol (for type checking)
	Decl ast.Statement // The declaration
}

// Loader handles loading and resolving modules.
type Loader struct {
	baseDir  string              // Base directory for resolving relative paths
	cache    map[string]*Module  // Cache of loaded modules
	loading  map[string]bool     // Tracks modules currently being loaded (for circular detection)
}

// NewLoader creates a new module loader.
func NewLoader(baseDir string) *Loader {
	return &Loader{
		baseDir: baseDir,
		cache:   make(map[string]*Module),
		loading: make(map[string]bool),
	}
}

// Load loads a module from the given path.
// The path can be relative (starting with "./" or "../") or absolute.
func (l *Loader) Load(importPath string, fromPath string) (*Module, error) {
	// Resolve the absolute path
	absPath, err := l.ResolvePath(importPath, fromPath)
	if err != nil {
		return nil, err
	}

	// Check cache
	if mod, ok := l.cache[absPath]; ok {
		return mod, nil
	}

	// Check for circular dependency
	if l.loading[absPath] {
		return nil, fmt.Errorf("circular dependency detected: %s", absPath)
	}

	// Mark as loading
	l.loading[absPath] = true
	defer func() { delete(l.loading, absPath) }()

	// Read the file
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read module %s: %v", absPath, err)
	}

	// Parse the module
	lex := lexer.New(string(content))
	p := parser.New(lex)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parse errors in module %s: %s", absPath, strings.Join(p.Errors(), "; "))
	}

	// Collect exports
	exports := make(map[string]ExportInfo)
	for _, stmt := range program.Statements {
		if export, ok := stmt.(*ast.ExportModifier); ok {
			info := l.extractExportInfo(export.Decl)
			if info.Name != "" {
				exports[info.Name] = info
			}
		}
	}

	mod := &Module{
		Path:    absPath,
		Program: program,
		Exports: exports,
	}

	// Cache the module
	l.cache[absPath] = mod

	return mod, nil
}

// ResolvePath resolves an import path to an absolute file path.
func (l *Loader) ResolvePath(importPath string, fromPath string) (string, error) {
	var basePath string
	if fromPath != "" {
		basePath = filepath.Dir(fromPath)
	} else {
		basePath = l.baseDir
	}

	// Handle relative paths
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		resolved := filepath.Join(basePath, importPath)
		// Add .gts extension if not present
		if !strings.HasSuffix(resolved, ".gts") {
			resolved += ".gts"
		}
		absPath, err := filepath.Abs(resolved)
		if err != nil {
			return "", fmt.Errorf("failed to resolve path %s: %v", importPath, err)
		}
		return absPath, nil
	}

	return "", fmt.Errorf("invalid import path: %s (must start with ./ or ../)", importPath)
}

// extractExportInfo extracts export information from a declaration.
func (l *Loader) extractExportInfo(decl ast.Statement) ExportInfo {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		// Build function type
		params := make([]*types.Param, len(d.Params))
		for i, p := range d.Params {
			params[i] = &types.Param{
				Name: p.Name,
				Type: types.AnyType, // TODO: Resolve actual type
			}
		}
		return ExportInfo{
			Name: d.Name,
			Type: &types.Function{
				Params:     params,
				ReturnType: types.AnyType, // TODO: Resolve actual type
			},
			Decl: decl,
		}

	case *ast.ClassDecl:
		return ExportInfo{
			Name: d.Name,
			Type: &types.Class{Name: d.Name},
			Decl: decl,
		}

	case *ast.VarDecl:
		return ExportInfo{
			Name: d.Name,
			Type: types.AnyType, // TODO: Resolve actual type
			Decl: decl,
		}

	case *ast.TypeAliasDecl:
		return ExportInfo{
			Name: d.Name,
			Type: types.AnyType, // TODO: Resolve actual type
			Decl: decl,
		}

	case *ast.InterfaceDecl:
		return ExportInfo{
			Name: d.Name,
			Type: &types.Interface{Name: d.Name},
			Decl: decl,
		}
	}

	return ExportInfo{}
}

// GetExport returns the export info for a given name from a module.
func (m *Module) GetExport(name string) (ExportInfo, bool) {
	info, ok := m.Exports[name]
	return info, ok
}
