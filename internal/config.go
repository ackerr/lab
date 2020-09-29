package internal

import (
	"os"

	"github.com/ackerr/lab/utils"
)

var template = []byte(`[gitlab]
# Gitlab domain, like https://gitlab.com
base_url = ""

# Gitlab api token
token = ""

# Optional
# A directory, If you want $GOPATH/src style manager your repo，set codesapce
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
	LabDir = home + "/.config/lab/"
	ConfigPath = LabDir + "config.toml"
	ProjectPath = LabDir + ".projects"
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		file, err := os.OpenFile(ConfigPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		utils.Check(err)
		_, err = file.Write(template)
		utils.Check(err)
		defer file.Close()
	}
}
