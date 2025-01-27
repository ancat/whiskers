package cmd

import (
	"fmt"
	"whiskers/gem"

	"github.com/spf13/cobra"
)

var (
	outputPath string
)

var gemfileDiffCmd = &cobra.Command{
	Use:   "gemfile-diff [before-gemfile] [after-gemfile]",
	Short: "Compare two Gemfile.lock files",
	Long: `Compare two Gemfile.lock files and show what gems were added, removed, or changed.
For example:
  whiskers gemfile-diff Gemfile.lock.before Gemfile.lock.after
  whiskers gemfile-diff Gemfile.lock.before Gemfile.lock.after --output diff.json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		beforePath := args[0]
		afterPath := args[1]

		diff, err := gem.NewGemfileDiff(beforePath, afterPath)
		if err != nil {
			return fmt.Errorf("failed to compare Gemfile.lock files: %w", err)
		}

		if !diff.HasChanges() {
			fmt.Println("No changes found between the Gemfile.lock files")
			return nil
		}

		// Print added gems
		if added := diff.GetAddedGems(); len(added) > 0 {
			fmt.Println("\nAdded gems:")
			for _, gem := range added {
				fmt.Printf("  + %s (%s)\n", gem.Name, gem.Version)
			}
		}

		// Print removed gems
		if removed := diff.GetRemovedGems(); len(removed) > 0 {
			fmt.Println("\nRemoved gems:")
			for _, gem := range removed {
				fmt.Printf("  - %s (%s)\n", gem.Name, gem.Version)
			}
		}

		// Print version changes
		if changes := diff.GetVersionChanges(); len(changes) > 0 {
			fmt.Println("\nVersion changes:")
			for _, change := range changes {
				fmt.Printf("  ~ %s: %s → %s\n", 
					change.Name, change.Before.Version, change.After.Version)
			}
		}

		// Save to JSON file if output path is specified
		if outputPath != "" {
			if err := diff.SaveToJSON(outputPath); err != nil {
				return fmt.Errorf("failed to save diff to JSON: %w", err)
			}
			fmt.Printf("\nDiff saved to %s\n", outputPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemfileDiffCmd)
	gemfileDiffCmd.Flags().StringVarP(&outputPath, "output", "o", "", "save diff to JSON file")
} 
