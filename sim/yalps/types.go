package yalps

import "time"

type Constraint struct {
	Equal *float64 `json:"equal,omitempty"`
	Min   *float64 `json:"min,omitempty"`
	Max   *float64 `json:"max,omitempty"`
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
	Integers    interface{}             `json:"integers,omitempty"`
	Binaries    interface{}             `json:"binaries"`
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
	CheckCycles     bool
	MaxPivots       int
	Tolerance       float64
	Timeout         time.Duration
	MaxIterations   int
	IncludeZeroVars bool
}

type Solution struct {
	Status    SolutionStatus     `json:"status"`
	Result    float64            `json:"result"`
	Variables map[string]float64 `json:"variables"`
}
