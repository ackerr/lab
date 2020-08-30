package internal

import (
	"os/exec"

	"github.com/ackerr/lab/utils"
)

// Clone : git clone the gitlab project
func Clone(gitURL, path string) {
	cmd := exec.Command("git", "clone", gitURL, path)
	err := cmd.Run()
	utils.Check(err)
}
