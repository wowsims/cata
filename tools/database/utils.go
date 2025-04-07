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
		return nil, err
	}
	if len(arr) != expectedLen {
		fmt.Println("expected array of length %d, got %d", expectedLen, len(arr))
		return nil, fmt.Errorf("expected array of length %d, got %d", expectedLen, len(arr))
	}
	return arr, nil
}

func parseFloatArrayField(jsonStr string, expectedLen int) ([]float64, error) {
	var arr []float64
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return nil, err
	}
	if len(arr) != expectedLen {
		fmt.Println("expected array of length %d, got %d", expectedLen, len(arr))
		return nil, fmt.Errorf("expected array of length %d, got %d", expectedLen, len(arr))
	}
	return arr, nil
}

func ParseRandomSuffixOptions(optionsString sql.NullString) []int32 {
	if !optionsString.Valid || optionsString.String == "" {
		return []int32{}
	}
	parts := strings.Split(optionsString.String, ",")
	var opts []int32
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if num, err := strconv.Atoi(part); err == nil {
			opts = append(opts, int32(num))
		}
	}
	return opts
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
