package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func parseIntArrayField(jsonStr string, expectedLen int) ([]int, error) {
	var arr []int
	if jsonStr == "" {
		return arr, nil
	}
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON array: %w", err)
	}
	if len(arr) != expectedLen {
		return nil, fmt.Errorf("invalid array length: expected %d, got %d", expectedLen, len(arr))
	}
	return arr, nil
}

func parseFloatArrayField(jsonStr string, expectedLen int) ([]float64, error) {
	var arr []float64
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, fmt.Errorf("unmarshaling JSON array: %w", err)
	}
	if len(arr) != expectedLen {
		return nil, fmt.Errorf("invalid array length: expected %d, got %d", expectedLen, len(arr))
	}
	return arr, nil
}

func ParseRandomSuffixOptions(optionsString sql.NullString) ([]int32, error) {
	if !optionsString.Valid || optionsString.String == "" {
		return []int32{}, nil
	}

	parts := strings.Split(optionsString.String, ",")
	var opts []int32
	var parseErrors []string

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		num, err := strconv.Atoi(part)
		if err != nil {
			parseErrors = append(parseErrors, fmt.Sprintf("part %d (%s): %v", i, part, err))
			continue
		}
		opts = append(opts, int32(num))
	}

	if len(parseErrors) > 0 {
		return opts, fmt.Errorf("some values couldn't be parsed: %s", strings.Join(parseErrors, "; "))
	}

	return opts, nil
}

// Formats the input string so that it does not use more than maxLength characters
// as soon a whole word exceeds the character limit a new line will be created
func formatStrings(maxLength int, input []string) []string {
	result := []string{}
	for _, line := range input {
		words := strings.Split(strings.Trim(line, "\n\r "), " ")
		currentLine := ""
		for _, word := range words {
			if len(currentLine) > maxLength {
				result = append(result, currentLine)
				currentLine = ""
			}

			if len(currentLine) > 0 {
				currentLine += " "
			}

			currentLine += word
		}

		if len(result) > 0 || len(currentLine) > 0 {
			result = append(result, currentLine)
		}
	}

	return result
}
