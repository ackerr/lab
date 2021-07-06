package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ackerr/lab/internal"
	"github.com/ackerr/lab/utils"
	"github.com/spf13/cobra"
)

var maxDepth int

func init() {
	rootCmd.AddCommand(csCmd)
	csCmd.Flags().IntVarP(&maxDepth, "depth", "d", 5, "maximum depth to filepath walk")
}

var csCmd = &cobra.Command{
	Use:     "cs",
	Aliases: []string{"ws"},
	Short:   "Search repo in your codespace",
	Run:     searchCodespace,
}

func searchCodespace(_ *cobra.Command, args []string) {
	internal.Setup()
	codespace := internal.Config.Codespace
	if codespace == "" {
		utils.Err("use <lab config> to set codespace first")
	}
	baseDepth := strings.Count(codespace, string(os.PathSeparator))
	var projects []string
	err := filepath.Walk(codespace, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return filepath.SkipDir
		}
		currentDepth := strings.Count(path, string(os.PathSeparator)) - baseDepth
		if currentDepth > maxDepth {
			return filepath.SkipDir
		}
		if utils.FileExists(filepath.Join(path, ".git")) {
			p := strings.Replace(path, codespace, "", 1)
			if len(p) > 0 {
				projects = append(projects, p[1:])
			}
			return filepath.SkipDir
		}
		return err
	})
	utils.Check(err)
	if len(projects) == 0 {
		utils.Err("no projects in codespace")
	}
	path := internal.FuzzyFinder(projects)
	if path == "" {
		return
	}
	fmt.Println(filepath.Join(codespace, path))
}
