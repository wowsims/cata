package yalps

import "math"

func Float64Ptr(v float64) *float64 {
	return &v
}

func roundToPrecision(num, precision float64) float64 {
	rounding := math.Round(1.0 / precision)
	return math.Round((num+math.SmallestNonzeroFloat64)*rounding) / rounding
}
