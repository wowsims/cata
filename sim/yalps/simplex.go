package yalps

import (
	"math"
)

// Pivot operation
func pivot(tableau *Tableau, row, col int, nonZeroColumns []int) {
	quotient := tableau.Matrix[row*tableau.Width+col]
	leaving := tableau.VariableAtPosition[tableau.Width+row]
	entering := tableau.VariableAtPosition[col]
	tableau.VariableAtPosition[tableau.Width+row] = entering
	tableau.VariableAtPosition[col] = leaving
	tableau.PositionOfVariable[leaving] = col
	tableau.PositionOfVariable[entering] = tableau.Width + row

	// Reset nonZeroColumns without reallocating
	nzCount := 0

	// (1 / quotient) * R_pivot -> R_pivot
	rowStart := row * tableau.Width
	rowSlice := tableau.Matrix[rowStart : rowStart + tableau.Width]
	for c := 0; c < tableau.Width; c++ {
		value := rowSlice[c]
		if value > 1e-16 || value < -1e-16 {
			rowSlice[c] = value / quotient
			nonZeroColumns[nzCount] = c
			nzCount++
		} else {
			rowSlice[c] = 0.0
		}
	}
	tableau.Matrix[rowStart+col] = 1.0 / quotient

	// -M[r, col] * R_pivot + R_r -> R_r
	for r := 0; r < tableau.Height; r++ {
		if r == row {
			continue
		}
		rowRStart := r * tableau.Width
		coef := tableau.Matrix[rowRStart+col]
		if coef > 1e-16 || coef < -1e-16 {
			rowRSlice := tableau.Matrix[rowRStart : rowRStart + tableau.Width]
			for i := 0; i < nzCount; i++ {
				c := nonZeroColumns[i]
				rowRSlice[c] -= coef * rowSlice[c]
			}
			rowRSlice[col] = -coef / quotient
		}
	}
}

func phase2(tableau *Tableau, options *Options) (SolutionStatus, float64) {
	pivotHistory := make([][2]int, 0)
	nonZeroColumns := make([]int, tableau.Width)
	for iter := 0; iter < options.MaxPivots; iter++ {
		// Find entering column
		col := -1
		value := options.Precision
		for c := 1; c < tableau.Width; c++ {
			reducedCost := index(tableau, 0, c)
			if reducedCost > value {
				value = reducedCost
				col = c
			}
		}
		if col == -1 {
			return StatusOptimal, roundToPrecision(index(tableau, 0, 0), options.Precision)
		}

		// Find leaving row
		row := -1
		minRatio := math.Inf(1)
		for r := 1; r < tableau.Height; r++ {
			value := index(tableau, r, col)
			if value <= options.Precision {
				continue
			}
			rhs := index(tableau, r, 0)
			ratio := rhs / value
			if ratio < minRatio {
				row = r
				minRatio = ratio
				if ratio <= options.Precision {
					break
				}
			}
		}
		if row == -1 {
			return StatusUnbounded, float64(col)
		}

		if options.CheckCycles && hasCycle(pivotHistory, tableau, row, col) {
			return StatusCycled, math.NaN()
		}

		pivot(tableau, row, col, nonZeroColumns)
	}
	return StatusCycled, math.NaN()
}

func phase1(tableau *Tableau, options *Options) (SolutionStatus, float64) {
	pivotHistory := make([][2]int, 0)
	nonZeroColumns := make([]int, tableau.Width)
	for iter := 0; iter < options.MaxPivots; iter++ {

		// Find leaving row
		row := -1
		rhs := -options.Precision
		for r := 1; r < tableau.Height; r++ {
			value := tableau.Matrix[r*tableau.Width]
			if value < rhs {
				rhs = value
				row = r
			}
		}
		if row == -1 {
			return phase2(tableau, options)
		}

		// Find entering column
		col := -1
		maxRatio := -math.Inf(1)
		rowStart := row * tableau.Width
		for c := 1; c < tableau.Width; c++ {
			coef := tableau.Matrix[rowStart+c]
			if coef < -options.Precision {
				ratio := -tableau.Matrix[c] / coef
				if ratio > maxRatio {
					maxRatio = ratio
					col = c
				}
			}
		}
		if col == -1 {
			return StatusInfeasible, math.NaN()
		}

		if options.CheckCycles && hasCycle(pivotHistory, tableau, row, col) {
			return StatusCycled, math.NaN()
		}

		pivot(tableau, row, col, nonZeroColumns)
	}
	return StatusCycled, math.NaN()
}

// Cycle detection
func hasCycle(history [][2]int, tableau *Tableau, row, col int) bool {
	history = append(history, [2]int{tableau.VariableAtPosition[tableau.Width+row], tableau.VariableAtPosition[col]})
	for length := 6; length <= len(history)/2; length++ {
		cycle := true
		for i := 0; i < length; i++ {
			item := len(history) - 1 - i
			row1, col1 := history[item][0], history[item][1]
			row2, col2 := history[item-length][0], history[item-length][1]
			if row1 != row2 || col1 != col2 {
				cycle = false
				break
			}
		}
		if cycle {
			return true
		}
	}
	return false
}

func simplex(tableau *Tableau, options *Options) (SolutionStatus, float64) {
	status, result := phase1(tableau, options)
	return status, result
}
