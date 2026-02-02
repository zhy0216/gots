package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zhy0216/quickts/gots/pkg/typed"
	"github.com/zhy0216/quickts/gots/pkg/types"
)

const sqlRuntime = `
// SQL runtime helpers
type GTS_DB struct {
	handle *sql.DB
}

type GTS_Tx struct {
	handle *sql.Tx
}

type GTS_ExecResult struct {
	RowsAffected int
	LastInsertId int
}

func gts_connect(path string) *GTS_DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		panic(err)
	}
	return &GTS_DB{handle: db}
}

func gts_db_close(db *GTS_DB) {
	db.handle.Close()
}

func gts_db_begin(db *GTS_DB, callback func(*GTS_Tx)) {
	tx, err := db.handle.Begin()
	if err != nil {
		panic(err)
	}
	gtsTx := &GTS_Tx{handle: tx}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	callback(gtsTx)
	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func gts_sql_exec(queryable interface{ ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) }, query string, args ...any) *GTS_ExecResult {
	result, err := queryable.ExecContext(context.Background(), query, args...)
	if err != nil {
		panic(err)
	}
	rowsAffected, _ := result.RowsAffected()
	lastInsertId, _ := result.LastInsertId()
	return &GTS_ExecResult{RowsAffected: int(rowsAffected), LastInsertId: int(lastInsertId)}
}
`

// genSQLTaggedTemplate generates Go code for a SQL tagged template literal.
// It converts template interpolations to ? placeholders and passes values as query parameters.
func (g *Generator) genSQLTaggedTemplate(expr *typed.TaggedTemplateLit) string {
	// Build the query string with ? placeholders
	query := ""
	for i, part := range expr.Parts {
		query += part
		if i < len(expr.Expressions) {
			query += "?"
		}
	}

	// Generate argument list
	args := make([]string, len(expr.Expressions))
	for i, e := range expr.Expressions {
		args[i] = g.genExpr(e)
	}

	// Get the DB/Tx handle expression
	propExpr := expr.Tag.(*typed.PropertyExpr)
	dbExpr := g.genExpr(propExpr.Object)

	argsStr := ""
	if len(args) > 0 {
		argsStr = ", " + strings.Join(args, ", ")
	}

	// Determine query type from type args
	if len(expr.TypeArgs) == 0 {
		// Exec: no type param
		return fmt.Sprintf("gts_sql_exec(%s.handle, %q%s)", dbExpr, query, argsStr)
	}

	ta := expr.TypeArgs[0]

	if arrType, ok := ta.(*types.Array); ok {
		// Multi-row query: T[]
		return g.genSQLMultiRowQuery(dbExpr, query, argsStr, arrType.Element)
	}

	// Single-row query: T
	return g.genSQLSingleRowQuery(dbExpr, query, argsStr, ta)
}

// genSQLMultiRowQuery generates an IIFE that calls Query, iterates rows, and scans into structs.
func (g *Generator) genSQLMultiRowQuery(dbExpr, query, argsStr string, elemType types.Type) string {
	iface, ok := elemType.(*types.Interface)
	if !ok {
		return fmt.Sprintf("nil /* unsupported SQL row type: %s */", elemType.String())
	}

	typeName := exportName(iface.Name)
	scanFields := g.getSQLScanFields(iface)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("func() []*%s {\n", typeName))
	buf.WriteString(fmt.Sprintf("\trows, err := %s.handle.Query(%q%s)\n", dbExpr, query, argsStr))
	buf.WriteString("\tif err != nil { panic(err) }\n")
	buf.WriteString("\tdefer rows.Close()\n")
	buf.WriteString(fmt.Sprintf("\tvar result []*%s\n", typeName))
	buf.WriteString("\tfor rows.Next() {\n")
	buf.WriteString(fmt.Sprintf("\t\titem := &%s{}\n", typeName))
	buf.WriteString(fmt.Sprintf("\t\terr := rows.Scan(%s)\n", scanFields))
	buf.WriteString("\t\tif err != nil { panic(err) }\n")
	buf.WriteString("\t\tresult = append(result, item)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn result\n")
	buf.WriteString("}()")

	return buf.String()
}

// genSQLSingleRowQuery generates an IIFE that calls QueryRow and scans into a single struct.
// Returns nil if no row is found (sql.ErrNoRows).
func (g *Generator) genSQLSingleRowQuery(dbExpr, query, argsStr string, rowType types.Type) string {
	iface, ok := rowType.(*types.Interface)
	if !ok {
		return fmt.Sprintf("nil /* unsupported SQL row type: %s */", rowType.String())
	}

	typeName := exportName(iface.Name)
	scanFields := g.getSQLScanFields(iface)

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("func() *%s {\n", typeName))
	buf.WriteString(fmt.Sprintf("\trow := %s.handle.QueryRow(%q%s)\n", dbExpr, query, argsStr))
	buf.WriteString(fmt.Sprintf("\titem := &%s{}\n", typeName))
	buf.WriteString(fmt.Sprintf("\terr := row.Scan(%s)\n", scanFields))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\tif err == sql.ErrNoRows { return nil }\n")
	buf.WriteString("\t\tpanic(err)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn item\n")
	buf.WriteString("}()")

	return buf.String()
}

// getSQLScanFields generates "&item.Field1, &item.Field2, ..." for sql.Scan.
// It uses FieldOrder from the interface to preserve declaration order,
// which is critical since maps in Go have random iteration order.
func (g *Generator) getSQLScanFields(iface *types.Interface) string {
	scanFields := []string{}
	for _, fieldName := range iface.FieldOrder {
		scanFields = append(scanFields, fmt.Sprintf("&item.%s", exportName(fieldName)))
	}
	return strings.Join(scanFields, ", ")
}
