package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/ackerr/lab/internal"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync gitlab projects",
	Run: func(cmd *cobra.Command, args []string) {
		syncProjects()
	},
}

// 同步项目, 顺便按字母排个序
func syncProjects() {
	projects := internal.Projects()
	file, err := os.Create(internal.Config.ProjectsPath)
	if err != nil {
		internal.Err(err)
	}
	defer file.Close()
	var allNameSpace []string
	for _, p := range projects {
		allNameSpace = append(allNameSpace, p.PathWithNamespace)
	}
	sort.Strings(allNameSpace)
	for _, ns := range allNameSpace {
		fmt.Fprintln(file, ns)
	}
}
