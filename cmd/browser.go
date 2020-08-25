package cmd

import (
	"os/exec"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "Browser open the gitlab project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			project := args[0]
			openURL(project)
		} else {
			openURL("")
		}
	},
}

func openURL(project string) {
	if project == "" {
		project = internal.Fuzzy(internal.Config.ProjectsPath)
	}
	url := internal.Config.BaseURL
	url = strings.Join([]string{url, project, "merge_requests"}, "/")
	launcher, err := utils.BrowserLauncher()
	if err != nil {
		utils.Err(err)
	}
	cmd := exec.Command(launcher, url)
	err = cmd.Run()
	if err != nil {
		utils.Err(err)
	}
}
