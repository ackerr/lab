package cmd

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/internal/ui"
	"github.com/ackerr/lab/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var tailLine int64

func init() {
	ciCmd.Flags().StringP("remote", "r", "", "the remote's pipeline")
	ciCmd.Flags().StringP("branch", "b", "", "the branch's pipeline")
	ciCmd.Flags().BoolP("all", "l", false, "view all jobs, use <tab> to view multiple jobs")
	ciCmd.Flags().Int64P("lines", "n", 0, "display the last part of a job log")
	rootCmd.AddCommand(ciCmd)
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "View the pipeline jobs trace log, default view running job,",
	Run:   currentJobs,
}

func currentJobs(cmd *cobra.Command, args []string) {
	internal.Setup()
	tailLine, err := cmd.Flags().GetInt64("lines")
	utils.Check(err)
	if tailLine == 0 {
		tailLine = internal.MainConfig.TailLineNumber
	}
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
	if resp.StatusCode == http.StatusNotFound {
		utils.Err(fmt.Sprintf("branch %s/%s not exist!", remote, branch))
	}
	opt := &gitlab.ListProjectPipelinesOptions{SHA: &c.ID}
	pipelines, _, err := client.Pipelines.ListProjectPipelines(p.ID, opt)
	utils.Check(err)
	jobs, _, err := client.Jobs.ListPipelineJobs(p.ID, pipelines[0].ID, &gitlab.ListJobsOptions{})
	utils.Check(err)

	isAll, _ := cmd.Flags().GetBool("all")
	if isAll {
		jobUI(client, p.ID, jobs)
		return
	}
	allDone := internal.TraceRunningJobs(client, p.ID, jobs, tailLine)
	if allDone {
		jobUI(client, p.ID, jobs)
	}
}

func jobUI(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job) {
	model := ui.NewJobUI(client, pid, jobs, tailLine)
	program := tea.NewProgram(model)
	err := program.Start()
	utils.Check(err)
}
