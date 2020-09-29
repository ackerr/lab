package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ackerr/lab/utils"
	"github.com/magiconair/properties/assert"
)

func TestMain(m *testing.M) {
	Setup()
	os.Exit(m.Run())
}

var tempDir = filepath.Join(utils.GetEnv("ROOT", os.TempDir()), "temp")

func TestCloneUseSSH(t *testing.T) {
	path, _ := ioutil.TempDir(tempDir, "ssh")
	repoURL := "git@github.com:Ackerr/lab.git"
	r := Clone(repoURL, path)
	assert.Equal(t, r, nil)
}

func TestCloneUseHTTPS(t *testing.T) {
	path, _ := ioutil.TempDir(tempDir, "https")
	repoURL := "https://github.com/Ackerr/lab.git"
	r := Clone(repoURL, path)
	assert.Equal(t, r, nil)
}
