package sql

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	// sqlite will create this if it doesn't exist
	db, err := sql.Open("sqlite3", "./notey.db")
	if err != nil {
		slog.Error("fatal error database init failure", slog.Any("err", err))
		os.Exit(1)
	}
	return db
}
