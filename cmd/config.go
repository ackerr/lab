package cmd

import (
	"os"
	"os/exec"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Use $EDITOR open config file",
	Run:   editConfig,
}

func editConfig(cmd *cobra.Command, args []string) {
	configPath := internal.ConfigPath
	editor := utils.GetEnv("EDITOR", "vim")
	command := exec.Command(editor, configPath)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	_ = command.Run()
}
