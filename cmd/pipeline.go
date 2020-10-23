package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/goware/prefixer"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

func init() {
	ciCmd.Flags().StringP("remote", "r", "", "the remote's pipeline")
	ciCmd.Flags().StringP("branch", "b", "", "the branch's pipeline")
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "View the pipeline jobs trace log",
	Run:   currentJobs,
}

func currentJobs(cmd *cobra.Command, args []string) {
	remote, _ := cmd.Flags().GetString("remote")
	branch, _ := cmd.Flags().GetString("branch")

	_, _ = internal.CurrentGitRepo()
	if branch == "" {
		branch = internal.CurrentBranch()
	}
	if remote == "" {
		remote = internal.CurrentRemote(branch)
	}
	gitURL := internal.RemoteURL(remote)
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

	m := make(map[string]*gitlab.Job)
	j := make([]string, 0)
	for _, job := range jobs {
		d := fmt.Sprintf("%-20s\t%-10s\t%-10s", job.Name, job.Status, job.Stage)
		m[d] = job
		j = append(j, d)
	}
	names := internal.FuzzyMultiFinder(j)
	if len(names[0]) == 0 {
		return
	}
	// remove the last empty item
	names = names[:len(names)-1]

	wg := sync.WaitGroup{}
	wg.Add(len(names))
	for _, name := range names {
		go func(n string) {
			_ = doTrace(client, p.ID, m[n])
			wg.Done()
		}(name)
	}
	wg.Wait()
}

func doTrace(client *gitlab.Client, pid interface{}, job *gitlab.Job) error {
	var offset int64
	prefix := utils.RandomColor(fmt.Sprintf("[%s] \u001b[0m", job.Name))
	for range time.NewTicker(time.Second * 3).C {
		trace, _, err := client.Jobs.GetTraceFile(pid, job.ID)
		utils.Check(err)
		prefixReader := prefixer.New(trace, prefix)
		_, err = io.CopyN(ioutil.Discard, prefixReader, offset)
		utils.Check(err)
		lenT, err := io.Copy(os.Stdout, prefixReader)
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		atomic.AddInt64(&offset, lenT)

		// job finshed
		if job.Status == "success" || job.Status == "failed" || job.Status == "cancelled" {
			return nil
		}
		job, _, err = client.Jobs.GetJob(pid, job.ID)
		utils.Check(err)
	}
	return nil
}
