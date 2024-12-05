package yalps

import (
	"container/heap"
	"math"
	"time"
)

type Cut struct {
	Sign     float64
	Variable int
	Value    float64
}

type Branch struct {
	Eval float64
	Cuts []Cut
}

type BranchHeap []Branch

func (h BranchHeap) Len() int            { return len(h) }
func (h BranchHeap) Less(i, j int) bool  { return h[i].Eval < h[j].Eval }
func (h BranchHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *BranchHeap) Push(x interface{}) { *h = append(*h, x.(Branch)) }
func (h *BranchHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// Find the most fractional variable
func mostFractionalVar(tableau *Tableau, intVars []int) (int, float64, float64) {
	highestFrac := 0.0
	variable := 0
	value := 0.0
	for _, intVar := range intVars {
		row := tableau.PositionOfVariable[intVar] - tableau.Width
		if row < 0 {
			continue
		}
		val := index(tableau, row, 0)
		frac := math.Abs(val - math.Round(val))
		if frac > highestFrac {
			highestFrac = frac
			variable = intVar
			value = val
		}
	}
	return variable, value, highestFrac
}

// Apply cuts to the tableau
func applyCuts(tableau *Tableau, buf *Buffer, cuts []Cut) *Tableau {
	width := tableau.Width
	height := tableau.Height
	matrix := buf.Matrix
	copy(matrix, tableau.Matrix)
	for i, cut := range cuts {
		sign, variable, value := cut.Sign, cut.Variable, cut.Value
		r := (height + i) * width
		pos := tableau.PositionOfVariable[variable]
		if pos < width {
			matrix[r] = sign * value
			for c := 1; c < width; c++ {
				matrix[r+c] = 0.0
			}
			matrix[r+pos] = sign
		} else {
			row := (pos - width) * width
			matrix[r] = sign * (value - tableau.Matrix[row])
			for c := 1; c < width; c++ {
				matrix[r+c] = -sign * tableau.Matrix[row+c]
			}
		}
	}
	positionOfVariable := buf.PositionOfVariable
	variableAtPosition := buf.VariableAtPosition
	copy(positionOfVariable, tableau.PositionOfVariable)
	copy(variableAtPosition, tableau.VariableAtPosition)
	length := width + height + len(cuts)
	for i := width + height; i < length; i++ {
		positionOfVariable[i] = i
		variableAtPosition[i] = i
	}
	return &Tableau{
		Matrix:             matrix[:len(tableau.Matrix)+len(cuts)*width],
		Width:              width,
		Height:             height + len(cuts),
		PositionOfVariable: positionOfVariable[:length],
		VariableAtPosition: variableAtPosition[:length],
		ColIdxBuffer:       tableau.ColIdxBuffer,
	}
}

// Branch and cut algorithm
func branchAndCut(tabmod TableauModel, initResult float64, options *Options) (TableauModel, SolutionStatus, float64) {
	tableau := tabmod.Tableau
	sign := tabmod.Sign
	integers := tabmod.Integers

	variable, value, frac := mostFractionalVar(tableau, integers)
	if frac <= options.Precision {
		// Initial solution is integer
		return tabmod, StatusOptimal, initResult
	}

	branches := &BranchHeap{}
	heap.Init(branches)
	heap.Push(branches, Branch{Eval: initResult, Cuts: []Cut{{-1, variable, math.Ceil(value)}}})
	heap.Push(branches, Branch{Eval: initResult, Cuts: []Cut{{1, variable, math.Floor(value)}}})

	maxExtraRows := len(integers) * 2
	matrixLength := len(tableau.Matrix) + maxExtraRows*tableau.Width
	posVarLength := len(tableau.PositionOfVariable) + maxExtraRows
	candidateBuffer := newBuffer(matrixLength, posVarLength)
	solutionBuffer := newBuffer(matrixLength, posVarLength)

	optimalThreshold := initResult * (1.0 - sign*options.Tolerance)
	stopTime := time.Now().Add(time.Millisecond * time.Duration(options.TimeoutMs))
	timedOut := time.Now().After(stopTime)
	solutionFound := false
	bestEval := math.Inf(1)
	bestTableau := tableau
	iter := 0

	for iter < options.MaxIterations && branches.Len() > 0 && bestEval >= optimalThreshold && !timedOut {
		branch := heap.Pop(branches).(Branch)
		relaxedEval, cuts := branch.Eval, branch.Cuts
		if relaxedEval > bestEval {
			break
		}

		currentTableau := applyCuts(tableau, candidateBuffer, cuts)
		status, result := simplex(currentTableau, options)

		if status == StatusOptimal && result < bestEval {
			variable, value, frac := mostFractionalVar(currentTableau, integers)
			if frac <= options.Precision {
				// Found integer solution
				solutionFound = true
				bestEval = result
				bestTableau = currentTableau
				candidateBuffer, solutionBuffer = solutionBuffer, candidateBuffer
			} else {
				numCuts := len(cuts)
				cutsUpper := make([]Cut, numCuts + 1)
				cutsLower := make([]Cut, numCuts + 1)
				copy(cutsUpper, cuts)
				copy(cutsLower, cuts)

				cutsLower[numCuts] = Cut{1, variable, math.Floor(value)}
				cutsUpper[numCuts] = Cut{-1, variable, math.Ceil(value)}

				heap.Push(branches, Branch{Eval: result, Cuts: cutsUpper})
				heap.Push(branches, Branch{Eval: result, Cuts: cutsLower})
			}
		}
		timedOut = time.Now().After(stopTime)
		iter++
	}

	unfinished := (timedOut || iter >= options.MaxIterations) && branches.Len() > 0 && bestEval >= optimalThreshold
	var status SolutionStatus
	if unfinished {
		status = StatusTimedOut
	} else if !solutionFound {
		status = StatusInfeasible
	} else {
		status = StatusOptimal
	}

	return TableauModel{Tableau: bestTableau, Sign: sign, Variables: tabmod.Variables, Integers: integers}, status, bestEval
}

// Buffer for storing tableau data
type Buffer struct {
	Matrix             []float64
	PositionOfVariable []int
	VariableAtPosition []int
}

func newBuffer(matrixLength, posVarLength int) *Buffer {
	return &Buffer{
		Matrix:             make([]float64, matrixLength),
		PositionOfVariable: make([]int, posVarLength),
		VariableAtPosition: make([]int, posVarLength),
	}
}
