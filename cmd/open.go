package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
)

var subpage string

func init() {
	openCmd.Flags().StringP("remote", "r", "", "the remote's pipeline")
	openCmd.Flags().BoolP("pipelines", "p", false, "open pipeline page, only support gitlab")
	openCmd.Flags().BoolP("merge_requests", "m", false, "open merge_requests page, only support gitlab")
	openCmd.Flags().StringVar(&subpage, "subpage", "", "open the input subpage")
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the current repo in web browser",
	Run:   openCurrentRepo,
}

func openCurrentRepo(cmd *cobra.Command, _ []string) {
	internal.Setup()
	if _, err := internal.CurrentGitRepo(); err != nil {
		return
	}
	remote, _ := cmd.Flags().GetString("remote")
	if remote == "" {
		branch := internal.CurrentBranch()
		remote = internal.CurrentRemote(branch)
	}
	gitURL := internal.RemoteURL(remote)
	url := internal.TransferGitURLToURL(gitURL)

	if len(subpage) == 0 {
		isPL, _ := cmd.Flags().GetBool("pipelines")
		if isPL {
			subpage = "pipelines"
		}
	}
	if len(subpage) == 0 {
		isMR, _ := cmd.Flags().GetBool("merge_requests")
		if isMR {
			subpage = "merge_requests"
		}
	}
	if url[len(url)-5:] == ".wiki" {
		url = fmt.Sprintf("%s/wikis", url[:len(url)-5])
	} else {
		url = fmt.Sprintf("%s/%s", url, subpage)
	}
	err := utils.OpenBrowser(url)
	utils.Check(err)
}
