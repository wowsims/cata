package yalps

import (
	"math"
	"sort"
)

// Main solver function
func Solve(model Model, options *Options) Solution {
	tabmod := tableauModel(model)
	opt := defaultOptions()
	if options != nil {
		opt = *options
	}
	status, result := simplex(tabmod.Tableau, &opt)

	if len(tabmod.Integers) == 0 || status != StatusOptimal {
		// Non-integer problem or no optimal solution
		return buildSolution(tabmod, status, result, &opt)
	} else {
		// Integer problem and optimal non-integer solution found
		intTabmod, intStatus, intResult := branchAndCut(tabmod, result, &opt)
		return buildSolution(intTabmod, intStatus, intResult, &opt)
	}
}

// Build the solution object
func buildSolution(tabmod TableauModel, status SolutionStatus, result float64, opt *Options) Solution {
	tableau := tabmod.Tableau
	sign := tabmod.Sign
	vars := tabmod.Variables

	if status == StatusOptimal || (status == StatusTimedOut && !math.IsNaN(result)) {
		var variables = make(map[string]float64)
		for i, varPair := range vars {
			varName := varPair[0].(string)
			row := tableau.PositionOfVariable[i+1] - tableau.Width
			value := 0.0
			if row >= 0 {
				value = index(tableau, row, 0)
			}
			if value > opt.Precision {
				variables[varName] = roundToPrecision(value, opt.Precision)
			} else if opt.IncludeZeroVars {
				variables[varName] = 0.0
			}
		}
		return Solution{
			Status:    status,
			Result:    -sign * result,
			Variables: variables,
		}
	} else if status == StatusUnbounded {
		var variableName string
		variable := tableau.VariableAtPosition[int(result)] - 1
		if 0 <= variable && variable < len(vars) {
			variableName = vars[variable][0].(string)
		}
		return Solution{
			Status:    StatusUnbounded,
			Result:    sign * math.Inf(1),
			Variables: map[string]float64{variableName: math.Inf(1)},
		}
	} else {
		// Infeasible, cycled, or timed out with NaN result
		return Solution{
			Status:    status,
			Result:    math.NaN(),
			Variables: make(map[string]float64),
		}
	}
}

// Default options
func defaultOptions() Options {
	return Options{
		Precision:       1e-8,
		CheckCycles:     false,
		MaxPivots:       8192,
		Tolerance:       0,
		TimeoutMs:       math.MaxInt32,
		MaxIterations:   32768,
		IncludeZeroVars: false,
	}
}

func sortedVariableKeys(vars map[string]Coefficients) []string {
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
