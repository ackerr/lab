package cmd

import (
	"fmt"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func init() {
	ciCmd.Flags().StringP("remote", "r", "", "the remote's pipeline")
	ciCmd.Flags().StringP("branch", "b", "", "the branch's pipeline")
	ciCmd.Flags().BoolP("all", "l", false, "view all jobs, default view running job")
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "View the pipeline jobs trace log",
	Run:   currentJobs,
}

func currentJobs(cmd *cobra.Command, args []string) {
	internal.Setup()
	remote, _ := cmd.Flags().GetString("remote")
	branch, _ := cmd.Flags().GetString("branch")

	if _, err := internal.CurrentGitRepo(); err != nil {
		return
	}
	if branch == "" {
		branch = internal.CurrentBranch()
	}
	if remote == "" {
		remote = internal.CurrentRemote(branch)
	}
	gitURL := internal.RemoteURL(remote)
	if !strings.Contains(gitURL, internal.Config.BaseURL[8:]) {
		utils.Err("not a gitlab repo")
	}
	project := internal.TransferGitURLToProject(gitURL)

	client := internal.NewClient()
	p, _, err := client.Projects.GetProject(project, &gitlab.GetProjectOptions{})
	utils.Check(err)
	c, resp, _ := client.Commits.GetCommit(p.ID, branch)
	if resp.StatusCode == 404 {
		utils.Err(fmt.Sprintf("branch %s/%s not exist!", remote, branch))
	}
	opt := &gitlab.ListProjectPipelinesOptions{SHA: &c.ID}
	pipelines, _, err := client.Pipelines.ListProjectPipelines(p.ID, opt)
	utils.Check(err)
	jobs, _, err := client.Jobs.ListPipelineJobs(p.ID, pipelines[0].ID, nil)
	utils.Check(err)

	isAll, _ := cmd.Flags().GetBool("all")
	internal.TraceJobs(client, p.ID, jobs, isAll)
}
