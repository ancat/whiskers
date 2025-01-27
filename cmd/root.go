package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "whiskers",
	Short: "Whiskers is a CLI tool for cats",
	Long: `A longer description of the Whiskers CLI tool
that can span multiple lines and provide more detailed
information about the application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you can define flags and configuration settings that are
	// global to all commands
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.whiskers.yaml)")
} 
