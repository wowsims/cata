package yalps

import (
	"math"
	"slices"
)

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

func relativeDifferenceFrom(delta float64, expected float64, precision float64) float64 {
	return (delta - precision) / max(math.Abs(expected), 1.0)
}

func relativeDifference(result float64, expected float64, precision float64) float64 {
	return relativeDifferenceFrom(math.Abs(result - expected), expected, precision)
}

func IsFinite(result float64) bool {
	return !math.IsNaN(result) && !math.IsInf(result, 0)
}
