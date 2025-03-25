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

func LoadRows[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error), args ...interface{}) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		results = append(results, item)
	}
	return results, nil
}

func CacheBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T {
	cache := make(map[K]T)
	for _, item := range items {
		cache[keyFunc(item)] = item
	}
	return cache
}
