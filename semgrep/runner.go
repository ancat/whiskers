package semgrep

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Runner handles executing semgrep and parsing its results
type Runner struct {
	rulesPath string
}

// SemgrepOutput represents the JSON structure of semgrep's output
type SemgrepOutput struct {
	Results []map[string]interface{} `json:"results"`
}

// NewRunner creates a new Runner instance
func NewRunner(rulesPath string) *Runner {
	if rulesPath == "" {
		rulesPath = "./semgrep-rules"
	}
	return &Runner{
		rulesPath: rulesPath,
	}
}

// Scan runs semgrep on the given files and returns the findings
func (r *Runner) Scan(files []string) ([]*Finding, error) {
	if len(files) == 0 {
		return []*Finding{}, nil
	}

	// Build the command
	args := []string{
		"--config", r.rulesPath,
		"--json",
		"--quiet",
	}
	args = append(args, files...)

	// Run semgrep
	cmd := exec.Command("semgrep", args...)
	output, err := cmd.Output()

	// Handle command execution errors
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("semgrep failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run semgrep: %w", err)
	}

	// Parse the JSON output
	var semgrepOutput SemgrepOutput
	if err := json.Unmarshal(output, &semgrepOutput); err != nil {
		return nil, fmt.Errorf("failed to parse semgrep output: %w", err)
	}

	// Convert results to findings
	findings := make([]*Finding, 0, len(semgrepOutput.Results))
	for _, result := range semgrepOutput.Results {
		findings = append(findings, NewFinding(result))
	}

	return findings, nil
}
