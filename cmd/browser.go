package cmd

import (
	"fmt"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

func init() {
	browserCmd.Flags().BoolP("pipelines", "p", false, "open pipeline page")
	browserCmd.Flags().BoolP("merge_requests", "m", false, "open merge_requests page")
	rootCmd.AddCommand(browserCmd)
}

var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "Browser open the gitlab project",
	Run:   openURL,
}

func openURL(cmd *cobra.Command, args []string) {
	internal.Setup()
	var project string
	if len(args) > 0 {
		project = args[0]
	} else {
		project = ""
	}
	if project == "" {
		project = internal.FuzzyLine(internal.ProjectPath)
	}
	// ctrl-c
	if project == "" {
		return
	}

	url := fmt.Sprintf("%s/%s", internal.Config.BaseURL, project)

	subpage := ""
	isPL, _ := cmd.Flags().GetBool("pipelines")
	if isPL {
		subpage = "pipelines"
	}
	isMR, _ := cmd.Flags().GetBool("merge_requests")
	if isMR {
		subpage = "merge_requests"
	}
	if len(subpage) > 0 {
		url = fmt.Sprintf("%s/-/%s", url, subpage)
	}
	err := utils.OpenBrowser(url)
	utils.Check(err)
}
