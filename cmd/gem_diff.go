package cmd

import (
	"fmt"
	"path/filepath"
	"whiskers/gem"
	"whiskers/utils"

	"github.com/spf13/cobra"
)

var (
	gemDiffSourceURL string
)

var gemDiffCmd = &cobra.Command{
	Use:   "gem-diff [gem-name] [version1] [version2]",
	Short: "Compare two versions of a gem",
	Long: `Download and compare two versions of a Ruby gem to see what files changed.
For example:
  whiskers gem-diff rails 7.0.0 7.0.8.5
  whiskers gem-diff rails 7.0.0 7.0.8.5 --source https://custom-gems.org`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version1 := args[1]
		version2 := args[2]

		// Create gem instances
		source := gem.DefaultSource()
		if gemDiffSourceURL != "" {
			source = gem.Source{
				Type: "rubygems",
				URL:  gemDiffSourceURL,
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

		// Print the differences
		if !diff.HasChanges() {
			fmt.Printf("\nNo changes found between %s and %s\n", version1, version2)
			return nil
		}

		fmt.Printf("\nChanges from %s to %s:\n", version1, version2)

		if len(diff.Added) > 0 {
			fmt.Println("\nAdded files:")
			for _, file := range diff.Added {
				fmt.Printf("  + %s\n", file)
			}
		}

		if len(diff.Removed) > 0 {
			fmt.Println("\nRemoved files:")
			for _, file := range diff.Removed {
				fmt.Printf("  - %s\n", file)
			}
		}

		if len(diff.Changed) > 0 {
			fmt.Println("\nModified files:")
			for _, file := range diff.Changed {
				fmt.Printf("  ~ %s\n", file)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemDiffCmd)
	gemDiffCmd.Flags().StringVarP(&gemDiffSourceURL, "source", "s", "", "gem source URL (default is RubyGems.org)")
} 
