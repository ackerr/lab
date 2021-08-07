package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
)

func init() {
	syncCmd.Flags().Bool("all", false, "sync all projects, default sync project if you are the membership")
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync [--all]",
	Short: "Sync gitlab projects",
	Run: func(cmd *cobra.Command, args []string) {
		syncAll, _ := cmd.Flags().GetBool("all")
		syncProjects(syncAll)
	},
}

// 同步项目, 顺便按字母排个序
func syncProjects(syncAll bool) {
	internal.Setup()
	file, err := os.Create(internal.ProjectPath)
	utils.Check(err)

	defer file.Close()
	ns := internal.Projects(syncAll)
	sort.Strings(ns)
	for _, n := range ns {
		if n != "" {
			fmt.Fprintln(file, n)
		}
	}
	println("Done.")
}
