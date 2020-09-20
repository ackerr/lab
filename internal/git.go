package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ackerr/lab/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func GitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	return string(output), err
}

// Clone : git clone the gitlab project
func Clone(gitURL, path string, useHTTPS bool) *git.Repository {
	var auth transport.AuthMethod
	if !useHTTPS {
		sshPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
		keyFile := utils.GetEnv("PRIVATE_KEY_PATH", sshPath)
		sshKey, err := ioutil.ReadFile(keyFile)
		utils.Check(err)
		auth, err = ssh.NewPublicKeys("git", sshKey, "")
		utils.Check(err)
	}
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		Auth:     auth,
		URL:      gitURL,
		Progress: os.Stdout,
	})
	utils.Check(err)
	return r
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
