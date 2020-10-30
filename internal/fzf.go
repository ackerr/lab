package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ackerr/lab/utils"
	"github.com/ktr0731/go-fuzzyfinder"
)

func readLines(filePath string) (lines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Sprintf("file %s doesn't exist, please run `lab sync` first", filePath)
		utils.Err(msg)
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

// FuzzyLine : fuzzy finder file line
func FuzzyLine(filePath string) string {
	lines, err := readLines(filePath)
	utils.Check(err)
	filtered := FuzzyFinder(lines)
	return filtered
}

// FuzzyFinder : fuzzy finder a content, if enter ctrl-c will return ""
func FuzzyFinder(lines []string) (filtered string) {
	if checkFZF() {
		filters := withFilter("fzf", func(in io.WriteCloser) {
			for _, line := range lines {
				fmt.Fprintln(in, line)
			}
		})
		filtered = filters[0]
	} else {
		index, err := fuzzyfinder.Find(lines, func(i int) string {
			return lines[i]
		})
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
	term := os.Getenv("TERM")
	if runtime.GOOS == "windows" && term == "srceen-256color" {
		os.Setenv("TERM", "")
		defer os.Setenv("TERM", term)
	}
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
	result, _ := cmd.Output()
	return strings.Split(string(result), "\n")
}
