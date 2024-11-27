package yalps

import (
	"math"
	"slices"
)

func Float64Ptr(v float64) *float64 {
	return &v
}

func roundToPrecision(num, precision float64) float64 {
	rounding := math.Round(1.0 / precision)
	return math.Round((num+math.SmallestNonzeroFloat64)*rounding) / rounding
}

func SetToSortedSlice(idxSet map[int]bool) []int {
	keys := make([]int, 0, len(idxSet))

	for k := range idxSet {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	return keys
}
