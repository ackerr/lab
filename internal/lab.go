package internal

import (
	"os"
	"strings"
	"sync"

	"github.com/ackerr/lab/utils"
	"github.com/schollz/progressbar/v3"
	"github.com/xanzy/go-gitlab"
)

var (
	// Config is global gitlab config
	Config  *gitlabConfig
	prePage int = 100
)

type gitlabConfig struct {
	BaseURL      string
	Version      string
	Token        string
	ProjectsPath string
}

func init() {
	token := getenv("GITLAB_TOKEN", "")
	if len(token) == 0 {
		utils.Err("set the gitlab token first, like\nexport GITLAB_TOKEN='GITLAB_TOKEN'")
	}

	baseURL := getenv("GITLAB_BASE_URL", "")
	if len(baseURL) == 0 {
		utils.Err("set the gitlab base url first, like\nexport GITLAB_BASE_URL='GITLAB_BASE_URL'")
	}
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}

	home, err := os.UserHomeDir()
	if err != nil {
		utils.Err(err)
	}

	Config = &gitlabConfig{
		BaseURL:      baseURL,
		Version:      getenv("GITLAB_API_VERSION", "v4"),
		ProjectsPath: getenv("GITLAB_PROJECT_PATH", home+"/.projects"),
		Token:        token,
	}
}

//Projects will return all projects's path with namespace
func Projects(syncAll bool) []string {
	path := gitlab.WithBaseURL(strings.Join([]string{Config.BaseURL, "api", Config.Version}, "/"))
	client, err := gitlab.NewClient(Config.Token, path)
	if err != nil {
		utils.Err(err)
	}
	items, totalPage := projects(client, 1, syncAll)
	bar := progressbar.Default(int64(totalPage), "Syncing")

	allProjects := make([]string, totalPage*prePage)
	ns := projectNameSpaces(items)
	for i, n := range ns {
		allProjects[i] = n
	}
	bar.Add(1)
	var wg sync.WaitGroup
	for curPage := 2; curPage <= totalPage; curPage++ {
		wg.Add(1)
		go func(cur int) {
			defer wg.Done()
			defer bar.Add(1)
			p, _ := projects(client, cur, syncAll)
			ns := projectNameSpaces(p)
			for i, n := range ns {
				allProjects[i+prePage*(cur-1)] = n
			}
		}(curPage)
	}
	wg.Wait()
	bar.Finish()
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

func getenv(key, defalut string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defalut
	}
	return value
}
