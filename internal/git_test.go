package internal

import (
	"github.com/ackerr/lab/utils"
	"github.com/go-git/go-git/v5"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	tempDir = filepath.Join(utils.GetEnv("ROOT", os.TempDir()), "temp")
)

func TestCloneUseSSH(t *testing.T) {
	path, _ := ioutil.TempDir(tempDir, "ssh")
	repoURL := "git@github.com:Ackerr/lab.git"
	r := Clone(repoURL, path, false)
	checkRepo(t, r)
}

func TestCloneUseHTTPS(t *testing.T) {
	path, _ := ioutil.TempDir(tempDir, "https")
	repoURL := "https://github.com/Ackerr/lab.git"
	r := Clone(repoURL, path, true)
	checkRepo(t, r)
}

func checkRepo(t *testing.T, r *git.Repository) {
	remotes, err := r.Remotes()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(remotes), 1)
	cfg, err := r.Config()
	assert.Equal(t, err, nil)
	assert.Equal(t, cfg.Branches["master"].Name, "master")
}
