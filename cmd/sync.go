package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/ackerr/lab/internal"
	"github.com/spf13/cobra"
)

func init() {
	syncCmd.Flags().Bool("all", false, "sync all projects, defalut sync project if you are the membership")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync gitlab projects",
	Run: func(cmd *cobra.Command, args []string) {
		syncAll, _ := cmd.Flags().GetBool("all")
		syncProjects(syncAll)
	},
}

// 同步项目, 顺便按字母排个序
func syncProjects(syncAll bool) {
	file, err := os.Create(internal.Config.ProjectsPath)
	if err != nil {
		internal.Err(err)
	}

	defer file.Close()
	ns := internal.Projects(syncAll)
	sort.Strings(ns)
	for _, n := range ns {
		if n != "" {
			fmt.Fprintln(file, n)
		}
	}
}
