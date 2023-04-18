package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteBackend struct {
	db *sql.DB
}

func (sqlb *SQLiteBackend) DB() *sql.DB {
	return sqlb.db
}

func (sqlb *SQLiteBackend) Open(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Could not open sqlite database: %s", err)
		return err
	}

	sqlb.db = db

	return err
}

func (sqlb *SQLiteBackend) Close() error {
	err := sqlb.Close()
	if err != nil {
		log.Fatalf("Could not close sqlite database: %s", err)
	}

	return err
}
