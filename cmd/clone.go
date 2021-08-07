package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
)

func init() {
	cloneCmd.Flags().Bool("https", false, "clone with https, default use ssh")
	cloneCmd.Flags().BoolP("current", "c", false, "clone repo to current directory")
	rootCmd.AddCommand(cloneCmd)
}

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the gitlab projects, like git clone",
	Run:   cloneRepo,
}

func cloneRepo(cmd *cobra.Command, args []string) {
	internal.Setup()
	projects := internal.FuzzyLines(internal.ProjectPath)
	if len(projects) == 0 {
		return
	}

	isHTTPS, _ := cmd.Flags().GetBool("https")
	isCurrent, _ := cmd.Flags().GetBool("current")
	for _, p := range projects {
		fmt.Println(utils.ColorFg("Cloning "+p, internal.MainConfig.ThemeColor))
		clone(p, isHTTPS, isCurrent)
	}
}

func clone(project string, isHTTPS, isCurrent bool) {
	var gitURL, path string
	baseURL := internal.Config.BaseURL
	if strings.HasPrefix(baseURL, "http") {
		baseURL = strings.Split(baseURL, "://")[1]
	}

	if !isHTTPS {
		gitURL = strings.Join([]string{"git@", baseURL, ":", project, ".git"}, "")
	} else {
		gitURL = strings.Join([]string{internal.Config.BaseURL, project}, "/")
	}
	codespace := internal.Config.Codespace
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	if !isCurrent && len(codespace) > 0 {
		dirs := []string{codespace, baseURL}
		dirs = append(dirs, strings.Split(project, "/")...)
		path = filepath.Join(dirs...)
		if utils.FileExists(path) {
			_ = internal.Fetch(path)
			return
		}
		err := os.MkdirAll(path, utils.DirPerm)
		go func() {
			<-sign
			_ = os.Remove(path)
		}()
		utils.Check(err)
	} else {
		// current path
		path, _ = os.Getwd()
		dir := strings.Split(project, "/")
		path = strings.Join([]string{path, dir[len(dir)-1]}, "/")
	}
	_ = internal.Clone(gitURL, path)
}
