package database

import (
	"encoding/json"
	"fmt"
)

func parseIntArrayField(jsonStr string, expectedLen int) ([]int, error) {
	var arr []int
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, err
	}
	if len(arr) != expectedLen {
		return nil, fmt.Errorf("expected array of length %d, got %d", expectedLen, len(arr))
	}
	return arr, nil
}

func decodeJSONIntArray(raw string) []int {
	var arr []int
	if raw == "" {
		return arr
	}
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil
	}
	return arr
}

func decodeJSONFloatArray(raw string) []float64 {
	var arr []float64
	if raw == "" {
		return arr
	}
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil
	}
	return arr
}
