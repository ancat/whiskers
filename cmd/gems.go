package cmd

import (
	"fmt"
	"whiskers/gem"

	"github.com/spf13/cobra"
)

var (
	showSource bool
)

var gemsCmd = &cobra.Command{
	Use:   "gems [path/to/Gemfile.lock]",
	Short: "List all gems in a Gemfile.lock",
	Long: `Parse a Gemfile.lock and display all gem dependencies with their versions.
For example:
  whiskers gems Gemfile.lock
  whiskers gems Gemfile.lock --show-source`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gemfileLock, err := gem.FromFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read Gemfile.lock: %w", err)
		}

		deps := gemfileLock.GetAllDependencies()
		fmt.Printf("Found %d dependencies\n", len(deps))

		for _, gem := range deps {
			if showSource {
				fmt.Printf("%s (%s) from %s\n", gem.Name, gem.Version, gem.Source.URL)
			} else {
				fmt.Printf("%s (%s)\n", gem.Name, gem.Version)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemsCmd)
	gemsCmd.Flags().BoolVarP(&showSource, "show-source", "s", false, "show the source URL for each gem")
} 
