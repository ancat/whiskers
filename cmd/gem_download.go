package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"whiskers/gem"

	"github.com/spf13/cobra"
)

var (
	sourceURL string
)

var gemDownloadCmd = &cobra.Command{
	Use:   "gem-download [gem-name] [version]",
	Short: "Download and extract a Ruby gem",
	Long: `Download and extract a Ruby gem from RubyGems.org or a specified source.
For example:
  whiskers gem-download rails 7.0.8.5
  whiskers gem-download rails 7.0.8.5 --source https://custom-gems.org`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		version := args[1]

		// Create gem instance with appropriate source
		source := gem.DefaultSource()
		if sourceURL != "" {
			source = gem.Source{
				Type: "rubygems",
				URL:  sourceURL,
			}
		}

		g := gem.NewGem(name, version, source)

		// Create the base directory for gems
		baseDir := "/tmp/gems"
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return fmt.Errorf("failed to create gems directory: %w", err)
		}

		fmt.Printf("Downloading %s (%s) from %s...\n", g.Name, g.Version, g.Source.URL)
		
		if err := g.DownloadAndExtract(baseDir); err != nil {
			return fmt.Errorf("failed to download and extract gem: %w", err)
		}

		targetDir := filepath.Join(baseDir, fmt.Sprintf("%s-%s", g.Name, g.Version))
		fmt.Printf("Successfully downloaded and extracted to %s\n", targetDir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(gemDownloadCmd)
	gemDownloadCmd.Flags().StringVarP(&sourceURL, "source", "s", "", "gem source URL (default is RubyGems.org)")
} 
