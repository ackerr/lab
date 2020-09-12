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

var pageMap = map[string]string{
	"merge_requests": "merge_requests",
	"pipelines":      "pipelines",
	"overview":       "",
}

func openURL(project string) {
	if project == "" {
		project = internal.Fuzzy(internal.Config.ProjectsPath)
	}
	url := internal.Config.BaseURL
	pages := make([]string, 0, len(pageMap))
	for k, _ := range pageMap {
		pages = append(pages, k)
	}
	page := pageMap[internal.FuzzyFinder(pages)]
	url = strings.Join([]string{url, project, page}, "/")
	launcher, err := utils.BrowserLauncher()
	utils.Check(err)
	args := append(launcher, url)
	cmd := exec.Command(args[0], args[1:]...)
	err = cmd.Run()
	utils.Check(err)
}
