package database

type ExecResult struct {
	Error        error
	RowsAffected int64
}

// Exec a query in the database.
func (db *Database) Exec(query string, values ...any) *ExecResult {
	result := db.session.Exec(query, values...)

	return &ExecResult{
		RowsAffected: result.RowsAffected,
		Error:        result.Error,
	}
}
