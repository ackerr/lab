package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/spf13/cobra"
)

func init() {
	cloneCmd.Flags().Bool("https", false, "clone with https")
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the gitlab project, like git clone",
	Run: func(cmd *cobra.Command, args []string) {
		project := internal.Fuzzy(internal.Config.ProjectsPath)

		var path, giturl string
		if len(args) > 0 {
			path = args[0]
		} else {
			path, _ = os.Getwd()
			dir := strings.Split(project, "/")
			path = strings.Join([]string{path, dir[len(dir)-1]}, "/")
		}

		useHTTPS, _ := cmd.Flags().GetBool("https")
		if !useHTTPS {
			baseURL := internal.Config.BaseURL
			if strings.HasPrefix(baseURL, "http") {
				baseURL = strings.Split(baseURL, "://")[1]
			}
			giturl = strings.Join([]string{"git@", baseURL, ":", project, ".git"}, "")
		} else {
			giturl = strings.Join([]string{internal.Config.BaseURL, project}, "/")
		}
		fmt.Println("Cloning", project)
		internal.Clone(giturl, path)
	},
}
