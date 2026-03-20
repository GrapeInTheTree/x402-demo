package quiz

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runner executes quiz solutions in a temporary Go module.
type Runner struct {
	workDir string
}

// NewRunner creates a temp directory for compiling/testing quiz solutions.
func NewRunner() (*Runner, error) {
	dir, err := os.MkdirTemp("", "x402-quiz-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}

	// Initialize a Go module in the temp directory
	modContent := "module x402quiz\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(modContent), 0644); err != nil {
		os.RemoveAll(dir)
		return nil, fmt.Errorf("write go.mod: %w", err)
	}

	return &Runner{workDir: dir}, nil
}

// Run writes the solution and test files, then runs `go test`.
func (r *Runner) Run(solution, testCode string) *Result {
	solPath := filepath.Join(r.workDir, "solution.go")
	testPath := filepath.Join(r.workDir, "solution_test.go")

	if err := os.WriteFile(solPath, []byte(solution), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write solution: %v", err)}
	}
	if err := os.WriteFile(testPath, []byte(testCode), 0644); err != nil {
		return &Result{Error: fmt.Sprintf("write test: %v", err)}
	}

	cmd := exec.Command("go", "test", "-v", "-count=1", "./...")
	cmd.Dir = r.workDir
	out, err := cmd.CombinedOutput()
	output := string(out)

	result := &Result{Output: output}

	if err != nil {
		// Check if it's a compilation error vs test failure
		if strings.Contains(output, "build failed") ||
			strings.Contains(output, "cannot") ||
			strings.Contains(output, "undefined") ||
			strings.Contains(output, "syntax error") {
			result.Error = "Compilation failed"
			return result
		}
		result.Compiled = true
	} else {
		result.Compiled = true
	}

	// Count pass/fail from verbose output
	result.Passed = strings.Count(output, "--- PASS")
	result.Total = result.Passed + strings.Count(output, "--- FAIL")

	if result.Total == 0 && result.Compiled {
		// No test functions found or all skipped
		result.Total = 1
		if err == nil {
			result.Passed = 1
		}
	}

	return result
}

// TemplatePath returns the path where the solution file is written.
func (r *Runner) TemplatePath() string {
	return filepath.Join(r.workDir, "solution.go")
}

// Cleanup removes the temporary directory.
func (r *Runner) Cleanup() {
	os.RemoveAll(r.workDir)
}
