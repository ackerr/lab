package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ackerr/lab/utils"
)

func GitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	return string(output), err
}

// Clone git clone the gitlab project
func Clone(gitURL, path string) error {
	cmd := exec.Command("git", "clone", gitURL, path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if len(Config.Email) > 0 {
		_ = SetGitConfig("name", Config.Username, path)
		_ = SetGitConfig("email", Config.Email, path)
	}
	return err
}

func SetGitConfig(key, value, path string) error {
	args := []string{}
	if len(path) > 0 {
		args = append(args, "-C", path)
	}
	args = append(args, "config", "user."+key, value)
	_, err := GitCommand(args...)
	return err
}

// CurrentGitRepo return the GitRepo path
func CurrentGitRepo() (string, error) {
	output, err := GitCommand("rev-parse", "-q", "--show-toplevel")
	return string(output), err
}

// CurrentGitRepo return the current branch
func CurrentBranch() string {
	branch, err := SymbolicRef("HEAD", true)
	if branch == "" || err != nil {
		branch = "master"
	}
	return branch
}

func RemoteURL(remote string) string {
	gitURL, err := GitCommand("ls-remote", "--get-url", remote)
	if err != nil {
		utils.Err("git remote is not set for", remote)
	}
	gitURL = firstLine(gitURL)
	if gitURL == remote {
		utils.Err(remote, "is a wrong remote")
	}
	return gitURL
}

// CurrentGitRepo return the current branch ref remote
func CurrentRemote(branch string) string {
	remote, err := GitCommand("config", fmt.Sprintf("branch.%s.remote", branch))
	remote = firstLine(remote)
	if remote == "" || err != nil {
		remote = "origin"
	}
	return remote
}

// SymbolicRef return the ref branch
func SymbolicRef(ref string, short bool) (string, error) {
	args := []string{"symbolic-ref"}
	if short {
		args = append(args, "--short")
	}
	args = append(args, ref)
	output, err := GitCommand(args...)
	return firstLine(output), err
}

// the git command output always has the "/n"
func firstLine(output string) string {
	if i := strings.Index(output, "\n"); i >= 0 {
		return output[0:i]
	}
	return output
}

// TransferGitURLToURL example:
// git@github.com/Ackerr:lab.git     -> https://github.com/Ackerr/lab
// https://github.com/Ackerr/lab.git -> https://github.com/Ackerr/lab
func TransferGitURLToURL(gitURL string) string {
	var url string
	if strings.HasPrefix(gitURL, "https://") {
		url = gitURL[:len(gitURL)-4]
	}
	if strings.HasPrefix(gitURL, "git@") {
		url = gitURL[:len(gitURL)-4]
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}
	return url
}
