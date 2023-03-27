package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ackerr/lab/internal"
)

func init() {
	// init config after cobra command called
	cobra.OnInitialize(internal.SetupConfig)
	rootCmd.PersistentFlags().StringVar(&internal.ConfigPath, "config", "", "target config file (default is $HOME/.config/lab/config.toml)")
}

var rootCmd = &cobra.Command{
	Use:   "lab",
	Short: "Lab is a cli tool, include some shortcut for gitlab",
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
	},
}

// Execute is the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
