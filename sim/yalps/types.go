package yalps

import (
	"fmt"
	"strings"
)

type Constraint struct {
	Equal *float64 `json:"equal,omitempty"`
	Min   *float64 `json:"min,omitempty"`
	Max   *float64 `json:"max,omitempty"`
}

func (constraint Constraint) String() string {
	components := make([]string, 0, 3)

	if constraint.Equal != nil {
		components = append(components, fmt.Sprintf("Equal: %f", *constraint.Equal))
	}

	if constraint.Min != nil {
		components = append(components, fmt.Sprintf("Min: %f", *constraint.Min))
	}

	if constraint.Max != nil {
		components = append(components, fmt.Sprintf("Max: %f", *constraint.Max))
	}

	return fmt.Sprintf("{%s}", strings.Join(components, ", "))
}

type Coefficients map[string]float64

type OptimizationDirection string

const (
	Maximize OptimizationDirection = "maximize"
	Minimize                       = "minimize"
)

type Model struct {
	Direction   OptimizationDirection   `json:"direction"`
	Objective   string                  `json:"objective"`
	Constraints map[string]Constraint   `json:"constraints"`
	Variables   map[string]Coefficients `json:"variables"`
	AllIntegers bool                    `json:"allIntegers"`
	Integers    []string                `json:"integers,omitempty"`
	AllBinaries bool                    `json:"allBinaries"`
	Binaries    []string                `json:"binaries,omitempty"`
}

type SolutionStatus string

const (
	StatusOptimal    SolutionStatus = "optimal"
	StatusInfeasible                = "infeasible"
	StatusUnbounded                 = "unbounded"
	StatusTimedOut                  = "timedout"
	StatusCycled                    = "cycled"
)

type Options struct {
	Precision       float64
	CheckCycles     bool          `json:"checkCycles"`
	MaxPivots       int
	Tolerance       float64       `json:"tolerance"`
	TimeoutMs       int32         `json:"timeout"`
	MaxIterations   int           `json:"maxIterations"`
	IncludeZeroVars bool
}

type Solution struct {
	Status    SolutionStatus     `json:"status"`
	Result    float64            `json:"result"`
	Variables map[string]float64 `json:"variables"`
}
