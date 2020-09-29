package internal

import (
	"os"
	"path/filepath"

	"github.com/ackerr/lab/utils"
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
	LabDir      string
	ConfigPath  string
	ProjectPath string
)

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
