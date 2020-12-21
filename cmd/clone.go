package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

func init() {
	cloneCmd.Flags().Bool("https", false, "clone with https, default use ssh")
	cloneCmd.Flags().BoolP("current", "c", false, "clone repo to current directory")
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the gitlab project, like git clone",
	Run:   cloneRepo,
}

func cloneRepo(cmd *cobra.Command, args []string) {
	project := internal.FuzzyLine(internal.ProjectPath)
	if project == "" {
		return
	}

	var path, gitURL string
	if len(args) > 0 {
		path = args[0]
	} else {
		path, _ = os.Getwd()
		dir := strings.Split(project, "/")
		path = strings.Join([]string{path, dir[len(dir)-1]}, "/")
	}

	baseURL := internal.Config.BaseURL
	if strings.HasPrefix(baseURL, "http") {
		baseURL = strings.Split(baseURL, "://")[1]
	}

	useHTTPS, _ := cmd.Flags().GetBool("https")
	if !useHTTPS {
		gitURL = strings.Join([]string{"git@", baseURL, ":", project, ".git"}, "")
	} else {
		gitURL = strings.Join([]string{internal.Config.BaseURL, project}, "/")
	}
	current, _ := cmd.Flags().GetBool("current")
	codespace := internal.Config.Codespace
	if !current && len(codespace) > 0 {
		dirs := []string{codespace, baseURL}
		dirs = append(dirs, strings.Split(project, "/")...)
		path = filepath.Join(dirs...)
		if utils.FileExists(path) {
			_ = internal.Fetch(path)
			return
		}
		err := os.MkdirAll(path, 0755)
		utils.Check(err)
	}

	// if clone abort, remove the empty directory
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sign
		_ = os.Remove(path)
	}()
	_ = internal.Clone(gitURL, path)
}
