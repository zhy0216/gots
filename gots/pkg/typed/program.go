// Package typed provides type-annotated AST nodes for code generation.
package typed

import (
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Program represents a type-checked program ready for code generation.
type Program struct {
	// GoImports contains Go package imports.
	GoImports []*GoImportDecl

	// ModuleImports contains local module imports.
	ModuleImports []*ModuleImportDecl

	// DefaultImports contains default imports.
	DefaultImports []*DefaultImport

	// NamespaceImports contains namespace imports.
	NamespaceImports []*NamespaceImport

	// TypeAliases contains type alias declarations.
	TypeAliases []*TypeAlias

	// Enums contains enum declarations.
	Enums []*EnumDecl

	// Interfaces contains interface declarations.
	Interfaces []*InterfaceDecl

	// Classes contains class declarations.
	Classes []*ClassDecl

	// Functions contains top-level function declarations.
	Functions []*FuncDecl

	// TopLevel contains top-level statements (executed in main).
	TopLevel []Stmt

	// Exports contains the names of exported declarations.
	Exports []string

	// ReExports contains re-export statements.
	ReExports []*ReExportDecl

	// DefaultExports contains default export statements.
	DefaultExports []*DefaultExport
}

// TypeAlias represents a resolved type alias.
type TypeAlias struct {
	Name     string
	Resolved types.Type
}

// EnumDecl represents a typed enum declaration.
type EnumDecl struct {
	Name    string
	Members []*EnumMember
}

// EnumMember represents a member of an enum.
type EnumMember struct {
	Name  string
	Value int // The numeric value of this member
}

// Capture represents a captured variable in a closure.
type Capture struct {
	Name  string      // Variable name
	Type  types.Type  // Variable type
	Depth int         // Lexical depth where defined (0 = current, 1 = parent, etc.)
	Index int         // Slot index in that scope's frame
}
