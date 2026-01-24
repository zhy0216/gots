// Package typed provides type-annotated AST nodes for code generation.
package typed

import (
	"github.com/zhy0216/quickts/gots/pkg/types"
)

// Program represents a type-checked program ready for code generation.
type Program struct {
	// TypeAliases contains type alias declarations.
	TypeAliases []*TypeAlias

	// Classes contains class declarations.
	Classes []*ClassDecl

	// Functions contains top-level function declarations.
	Functions []*FuncDecl

	// TopLevel contains top-level statements (executed in main).
	TopLevel []Stmt
}

// TypeAlias represents a resolved type alias.
type TypeAlias struct {
	Name     string
	Resolved types.Type
}

// Capture represents a captured variable in a closure.
type Capture struct {
	Name  string      // Variable name
	Type  types.Type  // Variable type
	Depth int         // Lexical depth where defined (0 = current, 1 = parent, etc.)
	Index int         // Slot index in that scope's frame
}
