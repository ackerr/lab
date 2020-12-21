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
	err := os.MkdirAll(LabDir, 0755)
	utils.Check(err)
	if ConfigPath == "" {
		ConfigPath = filepath.Join(LabDir, "config.toml")
	}
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		utils.CopyFile("config.toml", ConfigPath)
	}
	buf, err := envsubst.ReadFile(ConfigPath)
	utils.Check(err)
	viper.AddConfigPath(LabDir)
	err = viper.ReadConfig(bytes.NewReader(buf))
	utils.Check(err)
	setup()
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
	TailLineNumber int64  `toml:"tail_line_number"`
	FZF            bool   `toml:"fzf"`
	CloneOpts      string `toml:"clone_opts"`
}

func setup() {
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
	if Config.Projects == "" {
		Config.Projects = filepath.Join(LabDir, ".projects")
	}
	ProjectPath = Config.Projects
}
