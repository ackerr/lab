package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"

	"github.com/ackerr/lab/utils"
)

// FuzzyLine : fuzzy finder file line
func FuzzyLine(filePath string) string {
	lines, err := utils.ReadLines(filePath)
	utils.Check(err)
	filtered := FuzzyFinder(lines)
	return filtered
}

// FuzzyLine : fuzzy finder file lines
func FuzzyLines(filePath string) []string {
	lines, err := utils.ReadLines(filePath)
	utils.Check(err)
	filtered := FuzzyMultiFinder(lines)
	return filtered
}

// FuzzyFinder : fuzzy finder content
func FuzzyFinder(lines []string) (filtered string) {
	if checkFZF() {
		filters := withFilter("fzf", func(in io.WriteCloser) {
			for _, line := range lines {
				fmt.Fprintln(in, line)
			}
		})
		if len(filters) == 0 {
			filtered = ""
		} else {
			filtered = filters[0]
		}
	} else {
		index, err := fuzzyfinder.Find(lines, func(i int) string {
			return lines[i]
		})
		if err == fuzzyfinder.ErrAbort {
			return
		}
		utils.Check(err)
		filtered = lines[index]
	}
	return
}

// FuzzyMultiFinder : fuzzy finder multiple content
func FuzzyMultiFinder(lines []string) (filtered []string) {
	if checkFZF() {
		filtered = withFilter("fzf -m", func(in io.WriteCloser) {
			for _, line := range lines {
				fmt.Fprintln(in, line)
			}
		})
	} else {
		index, err := fuzzyfinder.FindMulti(lines, func(i int) string {
			return lines[i]
		})
		if err == fuzzyfinder.ErrAbort {
			return
		}
		utils.Check(err)
		for _, i := range index {
			filtered = append(filtered, lines[i])
		}
	}
	return
}

func checkFZF() bool {
	if !MainConfig.FZF {
		return false
	}
	_, err := exec.LookPath("fzf")
	return err == nil
}

func withFilter(command string, input func(in io.WriteCloser)) []string {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		input(in)
		in.Close()
	}()
	result, err := cmd.Output()
	if _, ok := err.(*exec.ExitError); ok {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(result)), "\n")
}
