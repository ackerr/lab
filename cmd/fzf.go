package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ackerr/lab/internal"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
)

func init() {
	fzfCmd.Flags().Bool("projects", false, "Return the gitlab projects")
}

var fzfCmd = &cobra.Command{
	Use:   "fzf",
	Short: "Use fzf to show file",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetBool("projects")
		if name {
			lines, _ := readLines(internal.Config.ProjectsPath)
			for _, line := range lines {
				fmt.Println(line)
			}
			return
		}

		var filePath string
		if len(args) > 0 {
			filePath = args[0]
			fmt.Println(fuzzy(filePath))
		}
	},
}

func readLines(filePath string) (lines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		internal.Err("If no ~/.projects file, you should run `lab sync` first")
	}
	defer file.Close()
	buffer := bufio.NewReader(file)
	for {
		value, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		}
		lines = append(lines, string(value))
	}
	return lines, err
}

func fuzzy(filePath string) string {
	lines, err := readLines(filePath)
	if err != nil {
		internal.Err(err)
	}
	index, err := fuzzyfinder.Find(lines, func(i int) string {
		return lines[i]
	})
	if err != nil {
		internal.Err(err)
	}
	return lines[index]
}
