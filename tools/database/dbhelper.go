package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type DBHelper struct {
	db *sql.DB
}

var DatabasePath string

func NewDBHelper() (*DBHelper, error) {
	db, err := sql.Open("sqlite", DatabasePath)
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
			fmt.Println(err.Error())
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

func RunOverrides(dbHelper *DBHelper, overridesFolder string) error {
	entries, err := os.ReadDir(overridesFolder)
	if err != nil {
		return fmt.Errorf("error reading overrides folder: %w", err)
	}

	type sqlFile struct {
		path  string
		order int
	}

	var files []sqlFile
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			base := strings.TrimSuffix(entry.Name(), ".sql")
			order, err := strconv.Atoi(base)
			if err != nil {
				continue
			}

			files = append(files, sqlFile{
				path:  filepath.Join(overridesFolder, entry.Name()),
				order: order,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].order < files[j].order
	})

	for _, f := range files {
		content, err := os.ReadFile(f.path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", f.path, err)
		}
		fmt.Println("Running override file", f.path)
		_, err = dbHelper.db.Exec(string(content))
		if err != nil {

			fmt.Println(err.Error())
			return fmt.Errorf("error executing SQL file %s: %w", f.path, err)
		}
	}

	return nil
}
