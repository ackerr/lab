package cmd

import (
	"github.com/spf13/cobra"
)

const version = "v0.2.20"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lab",
	Run: func(cmd *cobra.Command, args []string) {
		println("lab", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
