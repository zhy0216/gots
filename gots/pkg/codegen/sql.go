package codegen

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
