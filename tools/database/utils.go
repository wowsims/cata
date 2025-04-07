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

// func GetProfession(id int) proto.Profession {
// 	if profession, ok := dbc.MapProfessionIdToProfession[id]; ok {
// 		return profession
// 	}
// 	return 0
// }

// func GetSubclasses(mask int) []string {
// 	var result []string
// 	for flag, name := range dbc.MapItemSubclassNames {
// 		if mask&int(flag) != 0 {
// 			result = append(result, name)
// 		}
// 	}
// 	return result
// }

// func GetClassesFromClassMask(mask int) []proto.Class {
// 	var result []proto.Class
// 	for _, class := range classes {
// 		// Calculate the bit flag using 1 << (ID - 1)
// 		if mask&(1<<(class.ID-1)) != 0 {
// 			result = append(result, class.protoClass)
// 		}
// 	}
// 	return result
// }
