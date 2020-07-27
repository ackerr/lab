package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type gitlabConfig struct {
	BaseURL      string
	Version      string
	Token        string
	ProjectsPath string
}

var (
	// Config is global gitlab config
	Config *gitlabConfig
	client *gitlab.Client
	t      = gitlab.Bool(true)
	f      = gitlab.Bool(false)
)

func init() {
	token := getenv("GITLAB_TOKEN", "")
	if len(token) == 0 {
		Err("set the gitlab token first, like\nexport GITLAB_TOKEN='GITLAB_TOKEN'")
	}

	baseURL := getenv("GITLAB_BASE_URL", "")
	if len(baseURL) == 0 {
		Err("set the gitlab base url first, like\nexport GITLAB_BASE_URL='GITLAB_BASE_URL'")
	}
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "https://" + baseURL
	}
	if strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL[:len(baseURL)-1]
	}

	home, err := os.UserHomeDir()
	if err != nil {
		Err(err)
	}

	Config = &gitlabConfig{
		BaseURL:      baseURL,
		Version:      getenv("GITLAB_API_VERSION", "v4"),
		ProjectsPath: getenv("GITLAB_PROJECT_PATH", home+"/.projects"),
		Token:        token,
	}
}

//Projects will return all projects
func Projects() []*gitlab.Project {
	path := gitlab.WithBaseURL(strings.Join([]string{Config.BaseURL, "api", Config.Version}, "/"))
	client, err := gitlab.NewClient(Config.Token, path)
	if err != nil {
		Err(err)
	}
	prePage := 100
	totalPages := 2
	var allProjects []*gitlab.Project

	for curPage := 1; curPage <= totalPages; curPage++ {
		listOpt := gitlab.ListOptions{PerPage: prePage, Page: curPage}
		projectsOpt := gitlab.ListProjectsOptions{Simple: t, Membership: t, OrderBy: gitlab.String("name"), ListOptions: listOpt}
		projects, res, err := client.Projects.ListProjects(&projectsOpt)
		if err != nil {
			Err(err)
		}
		allProjects = append(allProjects, projects...)
		totalPages = res.TotalPages
	}
	return allProjects
}

func getenv(key, defalut string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defalut
	}
	return value
}

// Err will return the error message
func Err(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
