package internal

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ackerr/lab/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/xanzy/go-gitlab"
)

var (
	prePage    = 100
	apiVersion = "v4"
)

const interval = 3 * time.Second

func NewClient() *gitlab.Client {
	Setup()
	path := gitlab.WithBaseURL(strings.Join([]string{Config.BaseURL, "api", apiVersion}, "/"))
	client, err := gitlab.NewClient(Config.Token, path)
	if err != nil {
		utils.Err(err)
	}
	return client
}

// Projects will return all projects path with namespace
func Projects(syncAll bool) []string {
	client := NewClient()
	items, totalPage := projects(client, 1, syncAll)
	bar := progressbar.Default(int64(totalPage), "Syncing")

	allProjects := make([]string, totalPage*prePage)
	ns := projectNameSpaces(items)
	copy(allProjects[:len(ns)], ns)
	_ = bar.Add(1)
	var wg sync.WaitGroup
	for curPage := 2; curPage <= totalPage; curPage++ {
		wg.Add(1)
		go func(cur int) {
			defer wg.Done()
			p, _ := projects(client, cur, syncAll)
			ns := projectNameSpaces(p)
			start := prePage * (cur - 1)
			copy(allProjects[start:start+len(ns)], ns)
			_ = bar.Add(1)
		}(curPage)
	}
	wg.Wait()
	_ = bar.Finish()
	return allProjects
}

func projectNameSpaces(projects []*gitlab.Project) []string {
	ns := make([]string, len(projects))
	for i, p := range projects {
		ns[i] = p.PathWithNamespace
	}
	return ns
}

func projects(client *gitlab.Client, page int, syncAll bool) ([]*gitlab.Project, int) {
	listOpt := gitlab.ListOptions{PerPage: prePage, Page: page}
	projectsOpt := gitlab.ListProjectsOptions{Simple: gitlab.Bool(true), Membership: gitlab.Bool(!syncAll), ListOptions: listOpt}
	projects, res, err := client.Projects.ListProjects(&projectsOpt)
	utils.Check(err)
	return projects, res.TotalPages
}

// TransferGitURLToProject example:
// git@gitlab.com/Ackerr:lab.git     -> Ackerr/lab
// https://gitlab.com/Ackerr/lab.git -> Ackerr/lab
func TransferGitURLToProject(gitURL string) string {
	var url string
	if strings.HasPrefix(gitURL, "https://") {
		url = gitURL[len(Config.BaseURL) : len(gitURL)-4]
	}
	if strings.HasPrefix(gitURL, "git@") {
		url = gitURL[:len(gitURL)-4]
		url = strings.Split(url, ":")[1]
	}
	return url
}

func TraceRunningJobs(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job, tailLine int64) bool {
	wg := sync.WaitGroup{}
	allDone := true
	for _, job := range jobs {
		if !IsRunning(job.Status) {
			continue
		}
		allDone = false
		wg.Add(1)
		go func(j *gitlab.Job) {
			_ = DoTrace(client, pid, j, tailLine)
			wg.Done()
		}(job)
	}
	wg.Wait()
	return allDone
}

// IsRunning check job status
func IsRunning(status string) bool {
	if status == "created" || status == "pending" || status == "running" {
		return true
	}
	return false
}

// gitlab trace log has some hidden content like this ^[[0m^[[0K^[[36;1mStart^[[0;m`
// It will make prefix failure, so replace it. PS: \x1b == ^[
var re = regexp.MustCompile(`\x1b\[0m.*\[0K`)

func DoTrace(client *gitlab.Client, pid interface{}, job *gitlab.Job, tailLine int64) error {
	var offset int64 = 0
	firstTail := true
	prefix := utils.RandomColor(fmt.Sprintf("[%s] ", job.Name))
	for range time.NewTicker(interval).C {
		trace, _, err := client.Jobs.GetTraceFile(pid, job.ID)
		utils.Check(err)
		buffer, err := ioutil.ReadAll(trace)
		utils.Check(err)
		lines := strings.Split(string(buffer), "\n")
		length := len(lines)
		if firstTail {
			begin := int64(length) - tailLine
			if begin < 0 {
				begin = 0
			}
			lines = lines[begin : len(lines)-1]
			firstTail = false
		}
		for _, line := range lines[offset:] {
			println(re.ReplaceAllString(prefix+line, ``))
		}
		offset = int64(length)
		if !IsRunning(job.Status) {
			return nil
		}
		job, _, err = client.Jobs.GetJob(pid, job.ID)
		utils.Check(err)
	}
	return nil
}
