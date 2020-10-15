package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(browserCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(ciCmd)
}

var rootCmd = &cobra.Command{
	Use:   "lab",
	Short: "Lab is a cli tool, include some shortcut for gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute is the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
