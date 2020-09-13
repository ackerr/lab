package internal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/ackerr/lab/utils"
	"github.com/ktr0731/go-fuzzyfinder"
)

func readLines(filePath string) (lines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Sprintf("If thr file %s doesn't exist, you should run `lab sync` first", filePath)
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
	if checkCommandExists("fzf") {
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

func checkCommandExists(command string) bool {
	_, err := exec.LookPath(command)
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
	result, _ := cmd.Output()
	return strings.Split(string(result), "\n")
}
