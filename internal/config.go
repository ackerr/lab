package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/a8m/envsubst"
	"github.com/ackerr/lab/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var template = []byte(`[gitlab]
# Gitlab domain, like https://gitlab.com
base_url = "$GITLAB_BASE_URL"

# Gitlab api token
token = "$GITLAB_TOKEN"

# Optional
# A directory path, If you want $GOPATH/src style manager your repo，set codesapce like "~/workdir"
codespace = ""

# Default git config user.name
name = ""

# Default git config user.name
email = ""
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

func init() {
	home, _ := os.UserHomeDir()
	LabDir = filepath.Join(home, ".config", "lab")
	ConfigPath = filepath.Join(LabDir, "config.toml")
	ProjectPath = filepath.Join(LabDir, ".projects")
	err := os.MkdirAll(LabDir, 0755)
	utils.Check(err)
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		file, err := os.OpenFile(ConfigPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		utils.Check(err)
		_, err = file.Write(template)
		utils.Check(err)
		defer file.Close()
	}
}

type gitlabConfig struct {
	BaseURL   string `toml:"base_url"`
	Token     string `toml:"token"`
	Codespace string `toml:"codespace"`
	Name      string `toml:"name"`
	Email     string `toml:"email"`
}

type mainConfig struct {
	ThemeColor     string `toml:"theme_color"`
	TailLineNumber int64  `toml:"tail_line_number"`
	FZF            bool   `toml:"fzf"`
}

func Setup() {
	buf, err := envsubst.ReadFile(ConfigPath)
	utils.Check(err)
	viper.AddConfigPath(LabDir)
	err = viper.ReadConfig(bytes.NewReader(buf))
	utils.Check(err)

	Config = &gitlabConfig{}
	err = viper.Sub("gitlab").Unmarshal(Config, decodeOpt)
	utils.Check(err)

	MainConfig = &mainConfig{}
	viper.SetDefault("main.theme_color", "79")
	err = viper.Sub("main").Unmarshal(MainConfig, decodeOpt)
	utils.Check(err)
	if len(MainConfig.ThemeColor) == 0 {
		MainConfig.ThemeColor = "79"
	}
	if MainConfig.TailLineNumber == 0 {
		MainConfig.TailLineNumber = 20
	}

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
