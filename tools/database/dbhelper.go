package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DBHelper struct {
	db *sql.DB
}

func NewDBHelper(dbPath string) (*DBHelper, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return &DBHelper{db: db}, nil
}

func (d *DBHelper) Close() error {
	return d.db.Close()
}

func (d *DBHelper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

func (d *DBHelper) QueryAndProcess(query string, process func(*sql.Rows) error, args ...interface{}) {
	rows, err := d.Query(query, args...)
	if err != nil {
		log.Fatalf("query error: %v", err)
	}
	defer rows.Close()

	if err := process(rows); err != nil {
		log.Fatalf("processing error: %v", err)
	}
}
