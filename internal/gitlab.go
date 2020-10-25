package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/a8m/envsubst"
	"github.com/ackerr/lab/utils"
	"github.com/goware/prefixer"
	"github.com/mitchellh/mapstructure"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var (
	// Config global gitlab config
	Config     *gitlabConfig
	prePage    = 100
	apiVersion = "v4"
)

type gitlabConfig struct {
	BaseURL   string `toml:"base_url"`
	Token     string `toml:"token"`
	Codespace string `toml:"codespace"`
	Name      string `toml:"name"`
	Email     string `toml:"email"`
}

func Setup() {
	buf, err := envsubst.ReadFile(ConfigPath)
	utils.Check(err)
	viper.AddConfigPath(LabDir)
	err = viper.ReadConfig(bytes.NewReader(buf))
	utils.Check(err)

	Config = &gitlabConfig{}
	decodeOpt := func(config *mapstructure.DecoderConfig) { config.TagName = "toml" }
	err = viper.Sub("gitlab").Unmarshal(Config, decodeOpt)
	utils.Check(err)

	if len(Config.Token) == 0 {
		utils.Err("set Gitlab token first, use `lab config`")
	}

	baseURL := Config.BaseURL
	if len(baseURL) == 0 {
		utils.Err("set Gitlab base url first, use `lab config`")
	}
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}
	Config.BaseURL = baseURL

	home, err := os.UserHomeDir()
	utils.Check(err)
	codespace := Config.Codespace
	if strings.HasPrefix(codespace, "~") {
		codespace = filepath.Join(home, codespace[1:])
	}
	if strings.HasSuffix(codespace, "/") {
		codespace = codespace[:len(codespace)-1]
	}
	Config.Codespace = codespace
}

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
	if err != nil {
		return []*gitlab.Project{}, 0
	}
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

func TraceJobs(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job, isAll bool) {
	if isAll {
		traceAllJobs(client, pid, jobs)
		return
	}
	traceRunningJobs(client, pid, jobs)
}

func traceRunningJobs(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job) {
	wg := sync.WaitGroup{}
	for _, job := range jobs {
		wg.Add(1)
		go func(j *gitlab.Job) {
			_ = doTrace(client, pid, j)
			wg.Done()
		}(job)
	}
	wg.Wait()
	log.Println("\u001b[38;1mall jobs done")
}

func traceAllJobs(client *gitlab.Client, pid interface{}, jobs []*gitlab.Job) {
	m := make(map[string]*gitlab.Job)
	j := make([]string, 0)
	for _, job := range jobs {
		d := fmt.Sprintf("%-20s\t%-10s\t%-10s", job.Name, job.Status, job.Stage)
		m[d] = job
		j = append(j, d)
	}
	names := FuzzyMultiFinder(j)
	if len(names[0]) == 0 {
		return
	}
	// remove the last empty item
	names = names[:len(names)-1]

	wg := sync.WaitGroup{}
	wg.Add(len(names))
	for _, name := range names {
		go func(n string) {
			_ = doTrace(client, pid, m[n])
			wg.Done()
		}(name)
	}
	wg.Wait()

}

func isRunning(status string) bool {
	if status == "success" || status == "failed" || status == "cancelled" {
		return false
	}
	return true
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
		if !isRunning(job.Status) {
			return nil
		}
		job, _, err = client.Jobs.GetJob(pid, job.ID)
		utils.Check(err)
	}
	return nil
}
