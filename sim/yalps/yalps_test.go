package yalps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"slices"
	"strings"
	"testing"
	"time"
)

const maxDiff = 1e-5
const includeProfileRun = false
const profileTestCase = "mage"

// ExpectedSolution represents the expected outcome of the optimization.
type ExpectedSolution struct {
	Status    SolutionStatus     `json:"status"`
	Result    *float64           `json:"result"`
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
	if overrides.TimeoutMs != 0 {
		defaults.TimeoutMs = overrides.TimeoutMs
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
			// CPU Profile if requestd.
			if includeProfileRun && (tc.Name == profileTestCase) {
				f, err := os.Create(fmt.Sprintf("%s_cpu.pprof", tc.Name))
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()

				pprof.StartCPUProfile(f)
				defer pprof.StopCPUProfile()
			}

			// Set up constraints.
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
			startTime := time.Now()
			solution := Solve(model, &options)
			elapsedTime := time.Since(startTime)
			t.Logf("%s Solution: %+v", tc.Name, solution)
			t.Logf("Execution time: %v", elapsedTime)

			// Heap Profile if requestd.
			if includeProfileRun && (tc.Name == profileTestCase) {
				f2, err := os.Create(fmt.Sprintf("%s_heap.pprof", tc.Name))
				if err != nil {
					t.Fatal(err)
				}
				defer f2.Close()

				err = pprof.WriteHeapProfile(f2)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Assert the status.
			if solution.Status != tc.Expected.Status {
				t.Errorf("Expected status '%s', got '%s'", tc.Expected.Status, solution.Status)
			}

			// Assert the result with a tolerance.
			expectedResult := parseExpectedResult(&tc)

			if !resultIsOptimal(solution.Result, expectedResult, options) {
				t.Errorf("Expected result '%f', got '%f'", expectedResult, solution.Result)
			}

			// Assert validity of the variables.
			invalidVariables := checkVariables(solution, model, options.Precision)

			if invalidVariables != nil {
				t.Errorf("The following variables have invalid values: %s", strings.Join(invalidVariables, ", "))
			}

			// Assert that constraints were satisfied.
			violatedConstraints := checkConstraints(solution, model, options.Precision)

			if (len(violatedConstraints) > 0) && IsFinite(expectedResult) {
				t.Errorf("The following constraints were violated: %v", violatedConstraints)
			}

			if !t.Failed() {
				t.Logf("Test case '%s' passed successfully!", tc.Name)
			} else {
				t.Logf("Model struct for debugging:\n%+v", model)
				t.Logf("Options struct for debugging:\n%+v", options)
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
		testCase.Name = strings.Join(strings.Split(strings.Split(file.Name(), ".")[0], " "), "_")

		if !includeProfileRun || (testCase.Name == profileTestCase) {
			testCases = append(testCases, testCase)
		}
	}

	return testCases, nil
}

func parseExpectedResult(tc *TestCase) float64 {
	if tc.Expected.Status == StatusOptimal {
		return *tc.Expected.Result
	} else if tc.Expected.Status == StatusUnbounded {
		if tc.Model.Direction == Minimize {
			return math.Inf(-1)
		} else {
			return math.Inf(1)
		}
	} else {
		return math.NaN()
	}
}

func resultIsOptimal(result float64, expected float64, options Options) bool {
	if math.IsNaN(expected) {
		return math.IsNaN(result)
	} else if math.IsInf(expected, 0) {
		return result == expected
	} else {
		return relativeDifference(result, expected, options.Precision) <= max(options.Tolerance, maxDiff)
	}
}

func checkVariables(solution Solution, model Model, precision float64) []string {
	var invalidVariables []string

	for variableName, value := range solution.Variables {
		isBinary := model.AllBinaries || slices.Contains(model.Binaries, variableName)
		isInteger := isBinary || model.AllIntegers || slices.Contains(model.Integers, variableName)

		if (value < -precision) || (isInteger && (math.Abs(value - math.Round(value)) > precision)) || (isBinary && (value > 1.0 + precision)) {
			invalidVariables = append(invalidVariables, variableName)
		}
	}

	return invalidVariables
}

func getSummedCoefficients(solution Solution, model Model) Coefficients {
	sums := make(Coefficients)

	for variableKey, value := range solution.Variables {
		for coefficientKey, weight := range model.Variables[variableKey] {
			sums[coefficientKey] += value * weight
		}
	}

	return sums
}

func checkConstraints(solution Solution, model Model, precision float64) map[string]Constraint {
	violatedConstraints := map[string]Constraint{}
	summedCoefficients := getSummedCoefficients(solution, model)

	for key, constraint := range model.Constraints {
		sum := summedCoefficients[key]

		if ((constraint.Equal != nil) && (relativeDifference(sum, *constraint.Equal, precision) > maxDiff)) || ((constraint.Min != nil) && (relativeDifferenceFrom(*constraint.Min - sum, *constraint.Min, precision) > maxDiff)) || ((constraint.Max != nil) && (relativeDifferenceFrom(sum - *constraint.Max, *constraint.Max, precision) > maxDiff)) {
			violatedConstraints[key] = constraint
		}
	}

	return violatedConstraints
}
