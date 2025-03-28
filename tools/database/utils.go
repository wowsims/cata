package database

import (
	"encoding/json"
	"fmt"

	"github.com/wowsims/cata/sim/core/proto"
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

func (f ItemFlags) Has(flag ItemFlags) bool {
	return f&flag != 0
}

func (f InventoryType) Has(flag InventoryType) bool {
	return f&flag != 0
}

func GetSubclasses(mask int) []string {
	var result []string
	for flag, name := range SubclassNames {
		if mask&int(flag) != 0 {
			result = append(result, name)
		}
	}
	return result
}

func GetClassesFromClassMask(mask int) []proto.Class {
	var result []proto.Class
	for _, class := range classes {
		// Calculate the bit flag using 1 << (ID - 1)
		if mask&(1<<(class.ID-1)) != 0 {
			result = append(result, class.protoClass)
		}
	}
	return result
}
func ConvertGemTypeToProto(input int) proto.GemColor {
	switch GemType(input) {
	case Meta:
		return proto.GemColor_GemColorMeta
	case Red:
		return proto.GemColor_GemColorRed
	case Yellow:
		return proto.GemColor_GemColorYellow
	case Blue:
		return proto.GemColor_GemColorBlue
	case Orange:
		return proto.GemColor_GemColorOrange
	case Purple:
		return proto.GemColor_GemColorPurple
	case Green:
		return proto.GemColor_GemColorGreen
	case Prismatic:
		return proto.GemColor_GemColorPrismatic
	case Cogwheel:
		return proto.GemColor_GemColorCogwheel
	default:
		return proto.GemColor_GemColorUnknown
	}
}
