package api

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func OpenSQLite(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath)
}
