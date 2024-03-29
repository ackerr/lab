package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v0.4.3"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lab",
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Println("lab", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
