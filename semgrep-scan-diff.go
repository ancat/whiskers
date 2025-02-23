package main

import (
	"encoding/json"
	"fmt"
	"os"
	"whiskers/semgrep"

	"github.com/spf13/cobra"
)

// OutputData represents a sample data structure to be serialized as JSON.
type OutputData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	var beforeDir string
	var afterDir string
	var rulesDir string
	var outputFile string

	rootCmd := &cobra.Command{
		Use:   "semgrep-scan-diff",
		Short: "A CLI tool that runs semgrep on two directories and returns the diff of the findings",
		Run: func(cmd *cobra.Command, args []string) {
			if !directoryExists(beforeDir) {
				fmt.Fprintf(os.Stderr, "Error: 'before' directory does not exist: %s\n", beforeDir)
				os.Exit(1)
			}

			if !directoryExists(afterDir) {
				fmt.Fprintf(os.Stderr, "Error: 'after' directory does not exist: %s\n", afterDir)
				os.Exit(1)
			}

			if !directoryExists(rulesDir) {
				fmt.Fprintf(os.Stderr, "Error: 'rules' directory does not exist: %s\n", rulesDir)
				os.Exit(1)
			}

			// Sample output data.
			outputData := OutputData{
				Status:  "success",
				Message: "All required directories exist. Proceeding with processing.",
			}

			runner := semgrep.NewRunner(rulesDir)
			findings_before, err := runner.Scan([]string{beforeDir})
			findings_after, err2 := runner.Scan([]string{afterDir})

			if err != nil || err2 != nil {
				fmt.Printf("Before: %v; After: %v\n", err, err2)
				os.Exit(1)
			}


			// Create map of findings in version 1
			oldFindings := make(map[string]bool)
			for _, f := range findings_before {
				if err := f.Rebase(beforeDir); err != nil {
						fmt.Printf("  Warning: failed to rebase path: %v\n", err)
						continue
				}

				key := fmt.Sprintf("%s:%s:%s", f.Path, f.RuleID, f.Lines)
				fmt.Printf("old|%s\n", key)
				oldFindings[key] = true
			}

			// Find new findings
			var newFindings []*semgrep.Finding
			for _, f := range findings_after {
				if err := f.Rebase(afterDir); err != nil {
					fmt.Printf("  Warning: failed to rebase path: %v\n", err)
					continue
				}

				key := fmt.Sprintf("%s:%s:%s", f.Path, f.RuleID, f.Lines)
				fmt.Printf("new|%s\n", key)
				if !oldFindings[key] {
					fmt.Printf("comparing %s vs %s\n", f.Path, key)
					newFindings = append(newFindings, f)
				}
			}

			for _, f := range newFindings {
				fmt.Println(f.Display())
			}


			// If an output file is specified, write JSON data to it.
			if outputFile != "" {
				file, err := os.Create(outputFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
					os.Exit(1)
				}
				defer file.Close()

				encoder := json.NewEncoder(file)
				encoder.SetIndent("", "  ") // For pretty-printing.
				if err := encoder.Encode(outputData); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing JSON to file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("JSON output written to file: %s\n", outputFile)
			} else {
				// Otherwise, print a standard message.
				fmt.Println("All required directories exist. Proceeding with standard output...")
			}
		},
	}

	// Define the flags.
	rootCmd.Flags().StringVar(&beforeDir, "before", "", "Path to the 'before' directory")
	rootCmd.Flags().StringVar(&afterDir, "after", "", "Path to the 'after' directory")
	rootCmd.Flags().StringVar(&rulesDir, "rules", "", "Path to the 'rules' directory")
	rootCmd.Flags().StringVar(&outputFile, "output", "", "File path to write JSON output (optional)")

	// Mark the directory flags as required.
	rootCmd.MarkFlagRequired("before")
	rootCmd.MarkFlagRequired("after")
	rootCmd.MarkFlagRequired("rules")

	// Execute the command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// directoryExists checks if a given path exists and is a directory.
func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
