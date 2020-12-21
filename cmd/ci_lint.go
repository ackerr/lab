package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lintCmd)
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Check .gitlab-ci.yml syntax",
	Run:   lint,
}

func lint(cmd *cobra.Command, args []string) {
	var path string
	var err error
	if len(args) > 0 {
		path = args[0]
		path, err = filepath.Abs(path)
	} else {
		path, err = os.Getwd()
	}
	utils.Check(err)
	if !strings.HasSuffix(path, ".gitlab-ci.yml") {
		path = filepath.Join(path, ".gitlab-ci.yml")
	}
	if !utils.FileExists(path) {
		utils.Err(path, "not exist")
	}
	content, err := ioutil.ReadFile(path)
	utils.Check(err)
	if len(content) == 0 {
		utils.Err("empty .gitlab-ci.yml")
	}
	client := internal.NewClient()
	result, _, err := client.Validate.Lint(string(content))
	utils.Check(err)
	if result.Status != "valid" {
		for _, e := range result.Errors {
			fmt.Println(utils.ColorBg("ERROR", "#F08080)"), e)
		}
	}
}
