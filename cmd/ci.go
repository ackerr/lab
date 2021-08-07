package cmd

import (
	"fmt"
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/internal/ui"
	"github.com/ackerr/lab/utils"
)

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

	if _, err = internal.CurrentGitRepo(); err != nil {
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
		utils.Err("not", internal.Config.BaseURL, "repo")
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
	jobs = sortJob(jobs)

	isAll, _ := cmd.Flags().GetBool("all")
	if isAll {
		jobUI(client, p.ID, jobs, tailLine)
		return
	}
	allDone := internal.TraceRunningJobs(client, p.ID, jobs, tailLine)
	if allDone {
		jobUI(client, p.ID, jobs, tailLine)
	}
}

func jobUI(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job, tailLine int64) {
	model := ui.NewJobUI(client, pid, jobs, tailLine)
	program := tea.NewProgram(model)
	err := program.Start()
	utils.Check(err)
}

// list jobs api only sorted in descending order of their IDs.
func sortJob(jobs []*gitlab.Job) []*gitlab.Job {
	for i, j := 0, len(jobs)-1; i < j; i, j = i+1, j-1 {
		jobs[i], jobs[j] = jobs[j], jobs[i]
	}
	return jobs
}
