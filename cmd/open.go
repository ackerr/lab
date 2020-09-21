package cmd

import (
	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [remote]",
	Short: "Open the current repo in web browser",
	Run:   openCurrentRepo,
}

func openCurrentRepo(cmd *cobra.Command, args []string) {
	// check the path in a git repo
	var remote string
	if len(args) > 0 {
		remote = args[0]
	}
	_, err := internal.CurrentGitRepo()
	if err != nil {
		utils.Err("not a git repository")
	}
	if remote == "" {
		branch := internal.CurrentBranch()
		remote = internal.CurrentRemote(branch)
	}
	gitURL := internal.RemoteURL(remote)
	url := internal.TransferGitURLToURL(gitURL)
	err = utils.OpenBrowser(url)
	utils.Check(err)
}
