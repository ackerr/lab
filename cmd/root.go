package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(browserCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cloneCmd)
}

var rootCmd = &cobra.Command{
	Use:   "lab",
	Short: "Lab is a cli tool, include some shortcut for gitlab",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute is the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
