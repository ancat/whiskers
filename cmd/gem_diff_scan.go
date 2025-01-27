package cmd

import (
	"fmt"
	"path/filepath"
	"whiskers/gem"
	"whiskers/semgrep"
	"whiskers/utils"

	"github.com/spf13/cobra"
)

var (
	gemDiffScanSourceURL string
	rulesPath           string
)

var gemDiffScanCmd = &cobra.Command{
	Use:   "gem-diff-scan [gem-name] [version1] [version2]",
	Short: "Compare two versions of a gem and scan for new issues",
	Long: `Download and compare two versions of a Ruby gem, then run semgrep on the changes to find new issues.
For example:
  whiskers gem-diff-scan rails 7.0.0 7.0.8.5
  whiskers gem-diff-scan rails 7.0.0 7.0.8.5 --rules ./my-rules`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version1 := args[1]
		version2 := args[2]

		// Create gem instances
		source := gem.DefaultSource()
		if gemDiffScanSourceURL != "" {
			source = gem.Source{
				Type: "rubygems",
				URL:  gemDiffScanSourceURL,
			}
		}

		gem1 := gem.NewGem(name, version1, source)
		gem2 := gem.NewGem(name, version2, source)

		// Create base directory for downloads
		baseDir := "/tmp/gems"

		// Download and extract both versions
		fmt.Printf("Downloading %s (%s)...\n", name, version1)
		if err := gem1.DownloadAndExtract(baseDir); err != nil {
			return fmt.Errorf("failed to download %s version %s: %w", name, version1, err)
		}

		fmt.Printf("Downloading %s (%s)...\n", name, version2)
		if err := gem2.DownloadAndExtract(baseDir); err != nil {
			return fmt.Errorf("failed to download %s version %s: %w", name, version2, err)
		}

		// Get paths to the extracted gems
		path1 := filepath.Join(baseDir, fmt.Sprintf("%s-%s", name, version1))
		path2 := filepath.Join(baseDir, fmt.Sprintf("%s-%s", name, version2))

		// Compare the directories
		ignoreFiles := []string{
			"Gemfile.lock",
			".gitignore",
			"gem.deps.rb",
		}

		diff, err := utils.ComparePaths(path1, path2, ignoreFiles)
		if err != nil {
			return fmt.Errorf("failed to compare gem versions: %w", err)
		}

		if !diff.HasChanges() {
			fmt.Printf("\nNo changes found between %s and %s\n", version1, version2)
			return nil
		}

		// Create semgrep runner
		runner := semgrep.NewRunner(rulesPath)

		// Get files to scan (changed + added)
		var filesToScan1 []string
		var filesToScan2 []string

		// Add changed files
		for _, file := range diff.Changed {
			filesToScan1 = append(filesToScan1, filepath.Join(path1, file))
			filesToScan2 = append(filesToScan2, filepath.Join(path2, file))
		}

		// Add new files (only in version 2)
		for _, file := range diff.Added {
			filesToScan2 = append(filesToScan2, filepath.Join(path2, file))
		}

		// Run semgrep on both versions
		fmt.Printf("\nScanning files in version %s...\n", version1)
		findings1, err := runner.Scan(filesToScan1)
		if err != nil {
			return fmt.Errorf("failed to scan version %s: %w", version1, err)
		}

		fmt.Printf("Scanning files in version %s...\n", version2)
		findings2, err := runner.Scan(filesToScan2)
		if err != nil {
			return fmt.Errorf("failed to scan version %s: %w", version2, err)
		}

		// Create map of findings in version 1 for easy comparison
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
				newFindings = append(newFindings, f)
			}
		}

		// Print results
		if len(newFindings) == 0 {
			fmt.Println("\nNo new issues found!")
			return nil
		}

		fmt.Printf("\nFound %d new issues:\n", len(newFindings))
		for _, f := range newFindings {
			// Make the path relative to the gem root
			if err := f.Rebase(path2); err != nil {
				return fmt.Errorf("failed to rebase path: %w", err)
			}
			fmt.Println(f.Display())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemDiffScanCmd)
	gemDiffScanCmd.Flags().StringVarP(&gemDiffScanSourceURL, "source", "s", "", "gem source URL (default is RubyGems.org)")
	gemDiffScanCmd.Flags().StringVarP(&rulesPath, "rules", "r", "./semgrep-rules", "path to semgrep rules")
} 
