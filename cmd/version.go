package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lab",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lab v0.2.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
