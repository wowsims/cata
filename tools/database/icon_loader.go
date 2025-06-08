package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LoadArtTexturePaths(filePath string) (map[int]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'
	r.TrimLeadingSpace = true

	paths := make(map[int]string)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}
		if len(record) < 2 {
			continue
		}

		key, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid key %q on line %v: %w", record[0], record, err)
		}
		paths[key] = record[1]
	}

	return paths, nil
}

func GetIconName(artPaths map[int]string, fdid int) string {
	path, ok := artPaths[fdid]
	if !ok {
		return ""
	}

	fileName := filepath.Base(path)
	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	return strings.ToLower(fileName)
}
