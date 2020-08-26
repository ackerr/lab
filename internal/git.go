package internal

import (
	"os/exec"

	"github.com/ackerr/lab/utils"
)

// Clone : git clone the gitlab project
func Clone(giturl, path string) {
	cmd := exec.Command("git", "clone", giturl, path)
	err := cmd.Run()
	utils.Check(err)
}
