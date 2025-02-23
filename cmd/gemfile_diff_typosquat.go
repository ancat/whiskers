package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"whiskers/gem"
	"whiskers/utils"

	"github.com/spf13/cobra"
)

var (
	typosquatOutputPath string
)

var gemfileDiffTyposquatCmd = &cobra.Command{
	Use:   "gemfile-diff-typosquat [diff.json]",
	Short: "Check new gems in a Gemfile diff for potential typosquatting",
	Long: `Load a Gemfile diff from a JSON file and check new gems for potential typosquatting.
For example:
  whiskers gemfile-diff-typosquat diff.json`,
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

		potential_typosquats := utils.TyposquatsJSON{}
		potential_typosquats.Candidates = make(map[string][]string)


		// Check new gems for potential typosquatting
		if added := diff.GetAddedGems(); len(added) > 0 {
			fmt.Println("\nChecking new gems for potential typosquatting...")
			foundIssues := false
			for _, gem := range added {
				matches := utils.CheckForTyposquats(gem.Name)
				if len(matches) > 0 {
					if potential_typosquats.Candidates[gem.Name] == nil {
						potential_typosquats.Candidates[gem.Name] = matches
					}

					foundIssues = true
					fmt.Printf("\nWarning: %s might be a typosquat of these popular packages:\n", gem.Name)
					for _, match := range matches {
						fmt.Printf("  - %s\n", match)
					}
				}
			}
			if !foundIssues {
				fmt.Println("\nNo potential typosquatting issues found!")
			}
		} else {
			fmt.Println("\nNo new gems to check for typosquatting")
		}

		data, err := json.MarshalIndent(potential_typosquats, "", "  ")
		if err != nil {
			return err
		}

		os.WriteFile(typosquatOutputPath, data, 0644)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemfileDiffTyposquatCmd)
	gemfileDiffTyposquatCmd.Flags().StringVarP(&typosquatOutputPath, "output", "o", "", "save diff to JSON file")
}
