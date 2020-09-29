package cmd

import (
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
	internal.Setup()
	if project == "" {
		project = internal.FuzzyLine(internal.ProjectPath)
	}
	// ctrl-c
	if project == "" {
		return
	}
	url := internal.Config.BaseURL
	pages := make([]string, 0, len(pageMap))
	for k := range pageMap {
		pages = append(pages, k)
	}
	key := internal.FuzzyFinder(pages)
	// ctrl-c
	if key == "" {
		return
	}
	page := pageMap[key]
	url = strings.Join([]string{url, project, page}, "/")
	err := utils.OpenBrowser(url)
	utils.Check(err)
}
