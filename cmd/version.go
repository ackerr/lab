package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.2.7"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lab",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("lab %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
