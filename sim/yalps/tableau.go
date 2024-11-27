package yalps

import (
	"math"
)

type Tableau struct {
	Matrix             []float64
	Width              int
	Height             int
	PositionOfVariable []int
	VariableAtPosition []int
}

type TableauModel struct {
	Tableau   *Tableau
	Sign      float64
	Variables [][2]interface{}
	Integers  []int
}

func index(tableau *Tableau, row, col int) float64 {
	return tableau.Matrix[row*tableau.Width+col]
}

func update(tableau *Tableau, row, col int, value float64) {
	tableau.Matrix[row*tableau.Width+col] = value
}

func tableauModel(model Model) TableauModel {
	var direction float64
	if model.Direction == Minimize {
		direction = -1.0
	} else {
		direction = 1.0
	}

	variables := make([][2]interface{}, 0)
	variableIndices := make(map[string]int)
	for idx, key := range sortedVariableKeys(model.Variables) {
		coeffs := model.Variables[key]
		variables = append(variables, [2]interface{}{key, coeffs})
		variableIndices[key] = idx + 1 // Variable indices start from 1
	}

	// Handle integer and binary variables
	integers := make([]int, 0)
	binaryConstraintCols := make([]int, 0) // Initialize binaryConstraintCols
	// Process Binaries
	if model.Binaries != nil {
		for _, varName := range model.Binaries {
			if varIndex, exists := variableIndices[varName]; exists {
				integers = append(integers, varIndex)
				binaryConstraintCols = append(binaryConstraintCols, varIndex)
			}
		}
	} else if model.AllBinaries {
		for i := 1; i <= len(variables); i++ {
			integers = append(integers, i)
			binaryConstraintCols = append(binaryConstraintCols, i)
		}
	}
	// Handle Integers field similarly if necessary

	// Constraints processing
	constraints := make(map[string]struct {
		Row   int
		Lower float64
		Upper float64
	})
	for key, constraint := range model.Constraints {
		bounds := struct {
			Row   int
			Lower float64
			Upper float64
		}{Row: -1, Lower: -math.Inf(1), Upper: math.Inf(1)}

		if constraint.Equal != nil {
			bounds.Lower = *constraint.Equal
			bounds.Upper = *constraint.Equal
		} else {
			if constraint.Min != nil {
				bounds.Lower = math.Max(bounds.Lower, *constraint.Min)
			}
			if constraint.Max != nil {
				bounds.Upper = math.Min(bounds.Upper, *constraint.Max)
			}
		}
		constraints[key] = bounds
	}

	numConstraints := 1
	for key, constraint := range constraints {
		bounds := constraint
		bounds.Row = numConstraints
		numConstraints += 1
		if bounds.Lower > -math.Inf(1) {
			numConstraints += 1
		}
		constraints[key] = bounds
	}

	width := len(variables) + 1
	height := numConstraints + len(binaryConstraintCols)
	numVars := width + height
	matrix := make([]float64, width*height)
	positionOfVariable := make([]int, numVars)
	variableAtPosition := make([]int, numVars)
	tableau := &Tableau{
		Matrix:             matrix,
		Width:              width,
		Height:             height,
		PositionOfVariable: positionOfVariable,
		VariableAtPosition: variableAtPosition,
	}

	for i := 0; i < numVars; i++ {
		positionOfVariable[i] = i
		variableAtPosition[i] = i
	}

	// Build the tableau
	for c := 1; c < width; c++ {
		varPair := variables[c-1]
		varCoeffs := varPair[1].(Coefficients)
		for constraintKey, coef := range varCoeffs {
			if constraintKey == model.Objective {
				update(tableau, 0, c, direction*coef)
			}
			bounds, exists := constraints[constraintKey]
			if exists {
				if bounds.Upper < math.Inf(1) {
					update(tableau, bounds.Row, c, coef)
					if bounds.Lower > -math.Inf(1) {
						update(tableau, bounds.Row+1, c, -coef)
					}
				} else if bounds.Lower > -math.Inf(1) {
					update(tableau, bounds.Row, c, -coef)
				}
			}
		}
	}

	// Set RHS values
	for _, bounds := range constraints {
		if bounds.Upper < math.Inf(1) {
			update(tableau, bounds.Row, 0, bounds.Upper)
			if bounds.Lower > -math.Inf(1) {
				update(tableau, bounds.Row+1, 0, -bounds.Lower)
			}
		} else if bounds.Lower > -math.Inf(1) {
			update(tableau, bounds.Row, 0, -bounds.Lower)
		}
	}

	// Binary constraints
	for b, col := range binaryConstraintCols {
		row := numConstraints + b
		update(tableau, row, 0, 1.0)
		// Since variable indices start from 1, adjust the column index
		update(tableau, row, col, 1.0)
	}

	return TableauModel{
		Tableau:   tableau,
		Sign:      direction,
		Variables: variables,
		Integers:  integers,
	}
}
