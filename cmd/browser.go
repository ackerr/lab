package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

var openBrowerCmd = &cobra.Command{
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
		project = fuzzy(internal.Config.ProjectsPath)
	}
	url := internal.Config.BaseURL
	url = strings.Join([]string{url, project, "merge_requests"}, "/")
	browserCmd, err := utils.BrowserLauncher()
	if err != nil {
		internal.Err(err)
	}
	fmt.Println(browserCmd)
	cmd := exec.Command(browserCmd, url)
	err = cmd.Run()
	if err != nil {
		internal.Err(err)
	}
}
