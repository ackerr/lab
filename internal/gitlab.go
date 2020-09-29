package internal

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ackerr/lab/utils"
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
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(LabDir)

	err := viper.ReadInConfig()
	utils.Check(err)

	Config = &gitlabConfig{}
	decodeOpt := func(config *mapstructure.DecoderConfig) { config.TagName = "toml" }
	err = viper.Sub("gitlab").Unmarshal(Config, decodeOpt)
	utils.Check(err)

	if len(Config.Token) == 0 {
		token := utils.GetEnv("GITLAB_TOKEN", "")
		if len(token) == 0 {
			utils.Err("set Gitlab token first, use `lab config`")
		}
		Config.Token = token
	}

	baseURL := Config.BaseURL
	if len(baseURL) == 0 {
		baseURL = utils.GetEnv("GITLAB_BASE_URL", "")
		if len(baseURL) == 0 {
			utils.Err("set Gitlab base url first, use `lab config`")
		}
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

// Projects will return all projects path with namespace
func Projects(syncAll bool) []string {
	path := gitlab.WithBaseURL(strings.Join([]string{Config.BaseURL, "api", apiVersion}, "/"))
	client, err := gitlab.NewClient(Config.Token, path)
	if err != nil {
		utils.Err(err)
	}
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
