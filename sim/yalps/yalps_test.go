package yalps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"testing"
)

// ExpectedSolution represents the expected outcome of the optimization.
type ExpectedSolution struct {
	Status    SolutionStatus     `json:"status"`
	Result    float64            `json:"result"`
	Variables map[string]float64 `json:"variables"`
}

// TestCase represents a single test case with a model and its expected solution.
type TestCase struct {
	Name     string
	Model    Model            `json:"model"`
	Expected ExpectedSolution `json:"expected"`
	Options  Options          `json:"options"`
}

func mergeOptions(defaults, overrides Options) Options {
	if overrides.Precision != 0 {
		defaults.Precision = overrides.Precision
	}
	if overrides.CheckCycles {
		defaults.CheckCycles = overrides.CheckCycles
	}
	if overrides.MaxPivots != 0 {
		defaults.MaxPivots = overrides.MaxPivots
	}
	if overrides.Tolerance != 0 {
		defaults.Tolerance = overrides.Tolerance
	}
	if overrides.Timeout != 0 {
		defaults.Timeout = overrides.Timeout
	}
	if overrides.MaxIterations != 0 {
		defaults.MaxIterations = overrides.MaxIterations
	}
	if overrides.IncludeZeroVars {
		defaults.IncludeZeroVars = overrides.IncludeZeroVars
	}
	return defaults
}

func TestSolver(t *testing.T) {
	// Define the directory containing your JSON test cases.
	testCasesDir := "tests/cases"

	// Read all test cases from the directory.
	testCases, err := readTestCases(testCasesDir)
	if err != nil {
		t.Fatalf("Error reading test cases: %v", err)
	}

	if len(testCases) == 0 {
		t.Fatalf("No test cases found in directory: %s", testCasesDir)
	}

	// Iterate over each test case.
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase_%s", tc.Name), func(t *testing.T) {
			// Set up constraints with Float64Ptr.
			constraints := make(map[string]Constraint)
			for name, c := range tc.Model.Constraints {
				constraints[name] = Constraint{
					Equal: c.Equal,
					Min:   c.Min,
					Max:   c.Max,
				}
			}

			// Set up variables.
			variables := make(map[string]Coefficients)
			for name, coeffs := range tc.Model.Variables {
				variables[name] = Coefficients(coeffs)
			}

			// Set up the model.
			model := Model{
				Direction:   tc.Model.Direction,
				Objective:   tc.Model.Objective,
				Constraints: constraints,
				Variables:   variables,
				AllBinaries: tc.Model.AllBinaries,
				AllIntegers: tc.Model.AllIntegers,
				Integers:    tc.Model.Integers,
				Binaries:    tc.Model.Binaries,
			}

			// Set up solver options.

			options := mergeOptions(defaultOptions(), tc.Options)
			t.Logf("Running %s: Model Direction=%s, Objective=%s", tc.Name, model.Direction, model.Objective)

			// Execute the solver.
			solution := Solve(model, &options)
			t.Logf("%s Solution: %+v", tc.Name, solution)
			// Assert the status.
			if solution.Status != tc.Expected.Status {
				t.Errorf("Expected status '%s', got '%s'", tc.Expected.Status, solution.Status)
			}

			// Assert the result with a tolerance.
			resultDiff := math.Abs(solution.Result - tc.Expected.Result)
			if resultDiff > 1e-6 {
				t.Errorf("Expected result '%f', got '%f' (diff: %f)", tc.Expected.Result, solution.Result, resultDiff)
			}

			// Assert the variables.
			for varName, expectedValue := range tc.Expected.Variables {
				actualValue, exists := solution.Variables[varName]
				if !exists {
					t.Errorf("Expected variable '%s' not found in solution", varName)
					continue
				}

				valueDiff := math.Abs(actualValue - expectedValue)
				if valueDiff > 1e-6 {
					t.Errorf("Variable '%s': expected '%f', got '%f' (diff: %f)", varName, expectedValue, actualValue, valueDiff)
				}
			}

			// Optionally, assert that no unexpected variables are present.
			if len(tc.Expected.Variables) != len(solution.Variables) {
				t.Errorf("Expected %d variables, but got %d", len(tc.Expected.Variables), len(solution.Variables))
			}

			if !t.Failed() {
				t.Logf("Test case '%s' passed successfully!", tc.Name)
			}
		})
	}
}

// readTestCases reads all JSON test case files from the specified directory.
func readTestCases(dir string) ([]TestCase, error) {
	var testCases []TestCase

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read test cases directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue // Skip non-JSON files
		}

		filePath := filepath.Join(dir, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
		}

		var testCase TestCase
		if err := json.Unmarshal(data, &testCase); err != nil {
			return nil, fmt.Errorf("failed to parse JSON in file %s: %v", filePath, err)
		}
		testCase.Name = file.Name()

		testCases = append(testCases, testCase)
	}

	return testCases, nil
}
