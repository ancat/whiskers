package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"whiskers/gem"
	"whiskers/semgrep"
	"whiskers/utils"

	"github.com/spf13/cobra"
)

var (
	gemfileDiffScanRulesPath string
)

var gemfileDiffScanCmd = &cobra.Command{
	Use:   "gemfile-diff-scan [diff.json]",
	Short: "Load a Gemfile diff and scan changed gems for new issues",
	Long: `Load a Gemfile diff from a JSON file, download changed gems, and scan for new security issues.
For example:
  whiskers gemfile-diff-scan diff.json
  whiskers gemfile-diff-scan diff.json --rules ./my-rules`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		diffPath := args[0]

		// Check if file exists
		if _, err := os.Stat(diffPath); os.IsNotExist(err) {
			return fmt.Errorf("diff file not found: %s", diffPath)
		}

		// Load the diff from JSON
		diff, err := gem.LoadFromJSON(diffPath)
		if err != nil {
			return fmt.Errorf("failed to load diff from JSON: %w", err)
		}

		// Print the diff summary
		printDiffSummary(diff)

		// Create semgrep runner
		runner := semgrep.NewRunner(gemfileDiffScanRulesPath)

		// Process version changes
		changes := diff.GetVersionChanges()
		if len(changes) == 0 {
			fmt.Println("\nNo version changes to scan")
			return nil
		}

		fmt.Printf("\nScanning %d gems for security changes...\n", len(changes))

		// Map to store findings by gem
		newFindingsByGem := make(map[string][]*semgrep.Finding)

		// Create base directory for downloads
		baseDir := "/tmp/gems"

		// Process each changed gem
		for _, change := range changes {
			fmt.Printf("\nAnalyzing %s (%s → %s)...\n", change.Name, change.Before.Version, change.After.Version)

			// Download and extract both versions
			fmt.Printf("  Downloading version %s...\n", change.Before.Version)
			if err := change.Before.DownloadAndExtract(baseDir); err != nil {
				fmt.Printf("  Warning: failed to download version %s: %v\n", change.Before.Version, err)
				continue
			}

			fmt.Printf("  Downloading version %s...\n", change.After.Version)
			if err := change.After.DownloadAndExtract(baseDir); err != nil {
				fmt.Printf("  Warning: failed to download version %s: %v\n", change.After.Version, err)
				continue
			}

			// Get paths to the extracted gems
			beforePath := filepath.Join(baseDir, fmt.Sprintf("%s-%s", change.Name, change.Before.Version))
			afterPath := filepath.Join(baseDir, fmt.Sprintf("%s-%s", change.Name, change.After.Version))

			// Compare the directories
			diff, err := utils.ComparePaths(beforePath, afterPath, []string{
				"Gemfile.lock",
				".gitignore",
				"gem.deps.rb",
			})
			if err != nil {
				fmt.Printf("  Warning: failed to compare versions: %v\n", err)
				continue
			}

			if !diff.HasChanges() {
				fmt.Println("  No file changes found")
				continue
			}

			// Get files to scan
			var filesToScan1 []string
			var filesToScan2 []string

			// Add changed files
			for _, file := range diff.Changed {
				filesToScan1 = append(filesToScan1, filepath.Join(beforePath, file))
				filesToScan2 = append(filesToScan2, filepath.Join(afterPath, file))
			}

			// Add new files (only in version 2)
			for _, file := range diff.Added {
				filesToScan2 = append(filesToScan2, filepath.Join(afterPath, file))
			}

			// Run semgrep on both versions
			fmt.Printf("  Scanning files in version %s...\n", change.Before.Version)
			findings1, err := runner.Scan(filesToScan1)
			if err != nil {
				fmt.Printf("  Warning: failed to scan version %s: %v\n", change.Before.Version, err)
				continue
			}

			fmt.Printf("  Scanning files in version %s...\n", change.After.Version)
			findings2, err := runner.Scan(filesToScan2)
			if err != nil {
				fmt.Printf("  Warning: failed to scan version %s: %v\n", change.After.Version, err)
				continue
			}

			// Create map of findings in version 1
			oldFindings := make(map[string]bool)
			for _, f := range findings1 {
				key := fmt.Sprintf("%s:%s:%s", f.Path, f.RuleID, f.Lines)
				oldFindings[key] = true
			}

			// Find new findings
			var newFindings []*semgrep.Finding
			for _, f := range findings2 {
				key := fmt.Sprintf("%s:%s:%s", f.Path, f.RuleID, f.Lines)
				if !oldFindings[key] {
					if err := f.Rebase(afterPath); err != nil {
						fmt.Printf("  Warning: failed to rebase path: %v\n", err)
						continue
					}
					newFindings = append(newFindings, f)
				}
			}

			if len(newFindings) > 0 {
				newFindingsByGem[change.Name] = newFindings
			}
		}

		// Print results
		if len(newFindingsByGem) == 0 {
			fmt.Println("\nNo new security issues found!")
			return nil
		}

		fmt.Println("\nNew security issues found:")
		for gemName, findings := range newFindingsByGem {
			fmt.Printf("\n%s:\n", gemName)
			for _, f := range findings {
				fmt.Println(f.Display())
			}
		}

		return nil
	},
}

func printDiffSummary(diff *gem.GemfileDiff) {
	// Print added gems
	if added := diff.GetAddedGems(); len(added) > 0 {
		fmt.Println("\nAdded gems:")
		for _, gem := range added {
			fmt.Printf("  + %s (%s)\n", gem.Name, gem.Version)
			if !gem.IsFromRubyGems() {
				fmt.Printf("    source: %s (%s)\n", gem.Source.URL, gem.Source.Type)
			}
		}
	}

	// Print removed gems
	if removed := diff.GetRemovedGems(); len(removed) > 0 {
		fmt.Println("\nRemoved gems:")
		for _, gem := range removed {
			fmt.Printf("  - %s (%s)\n", gem.Name, gem.Version)
			if !gem.IsFromRubyGems() {
				fmt.Printf("    source: %s (%s)\n", gem.Source.URL, gem.Source.Type)
			}
		}
	}

	// Print version changes
	if changes := diff.GetVersionChanges(); len(changes) > 0 {
		fmt.Println("\nVersion changes:")
		for _, change := range changes {
			fmt.Printf("  ~ %s: %s → %s\n",
				change.Name,
				change.Before.Version,
				change.After.Version)

			if !change.Before.IsFromRubyGems() || !change.After.IsFromRubyGems() {
				if change.Before.Source != change.After.Source {
					fmt.Printf("    source changed: %s (%s) → %s (%s)\n",
						change.Before.Source.URL, change.Before.Source.Type,
						change.After.Source.URL, change.After.Source.Type)
				}
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(gemfileDiffScanCmd)
	gemfileDiffScanCmd.Flags().StringVarP(&gemfileDiffScanRulesPath, "rules", "r", "./semgrep-rules", "path to semgrep rules")
}
