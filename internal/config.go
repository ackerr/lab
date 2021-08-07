package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/ackerr/lab/utils"
)

var content = []byte(`[gitlab]
# gitlab domain, like https://gitlab.com
base_url = "$GITLAB_BASE_URL"

# gitlab access token
token = "$GITLAB_TOKEN"

# gitlab projects file
# default $HOME/config/.lab/.projects
projects = ""

# If set, lab clone and lab cs will use this path as target path
# default empty
codespace = ""

# If set, lab clone will auto set user.name in repo gitconfig
# default empty
name = ""

# If set, lab clone will auto set user.email in repo gitconfig
# default empty
email = ""

[main]
# If set 1, it will use fzf as fuzzy finder, default use go-fuzzyfinder
# default 0
fzf = 0

# lab clone extra custom git clone config
# example clone_opts="--origin ackerr --branch fix"
# default empty
clone_opts = ""
`)

var (
	// Config global gitlab config
	Config      *gitlabConfig
	MainConfig  *mainConfig
	LabDir      string
	ConfigPath  string
	ProjectPath string
)

var decodeOpt = func(config *mapstructure.DecoderConfig) { config.TagName = "toml" }

func SetupConfig() {
	home, _ := os.UserHomeDir()
	LabDir = filepath.Join(home, ".config", "lab")
	err := os.MkdirAll(LabDir, utils.DirPerm)
	utils.Check(err)
	if ConfigPath == "" {
		ConfigPath = filepath.Join(LabDir, "config.toml")
	}
	if _, err = os.Stat(ConfigPath); os.IsNotExist(err) {
		err = ioutil.WriteFile(ConfigPath, content, utils.FilePerm)
		utils.Check(err)
	}
	buf, err := envsubst.ReadFile(ConfigPath)
	utils.Check(err)
	viper.AddConfigPath(LabDir)
	err = viper.ReadConfig(bytes.NewReader(buf))
	utils.Check(err)
}

type gitlabConfig struct {
	BaseURL   string `toml:"base_url"`
	Token     string `toml:"token"`
	Codespace string `toml:"codespace"`
	Name      string `toml:"name"`
	Email     string `toml:"email"`
	Projects  string `toml:"projects"`
}

type mainConfig struct {
	ThemeColor     string `toml:"theme_color"`
	CloneOpts      string `toml:"clone_opts"`
	TailLineNumber int64  `toml:"tail_line_number"`
	FZF            bool   `toml:"fzf"`
}

func Setup() {
	// init main config
	MainConfig = &mainConfig{}
	viper.SetDefault("main.theme_color", "79")
	err := viper.Sub("main").Unmarshal(MainConfig, decodeOpt)
	utils.Check(err)
	if len(MainConfig.ThemeColor) == 0 {
		MainConfig.ThemeColor = "79"
	}
	if MainConfig.TailLineNumber == 0 {
		MainConfig.TailLineNumber = 20
	}

	// init gitlab config
	Config = &gitlabConfig{}
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
	Config.BaseURL = strings.TrimSuffix(baseURL, "/")

	home, err := os.UserHomeDir()
	utils.Check(err)
	codespace := Config.Codespace
	if strings.HasPrefix(codespace, "~") {
		codespace = filepath.Join(home, codespace[1:])
	}
	if strings.HasSuffix(codespace, string(os.PathSeparator)) {
		codespace = codespace[:len(codespace)-len(string(os.PathSeparator))]
	}
	Config.Codespace = codespace
	if Config.Projects == "" {
		Config.Projects = filepath.Join(LabDir, ".projects")
	}
	ProjectPath = Config.Projects
}
