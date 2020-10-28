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
	var err error
	if checkCommandExists("fzf") {
		var filters []string
		filters, err = withFilter("fzf", func(in io.WriteCloser) {
			for _, line := range lines {
				fmt.Fprintln(in, line)
			}
		})
		filtered = filters[0]
	}
	if err != nil {
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

func withFilter(command string, input func(in io.WriteCloser)) ([]string, error) {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", command)
	cmd.Stderr = os.Stderr
	in, err := cmd.StdinPipe()
	if err != nil {
		return []string{}, err
	}
	go func() {
		input(in)
		in.Close()
	}()
	result, err := cmd.Output()
	return strings.Split(string(result), "\n"), err
}

func FuzzyMultiFinder(lines []string) (filtered []string) {
	var err error
	if checkCommandExists("fzf") {
		filtered, err = withFilter("fzf -m", func(in io.WriteCloser) {
			for _, line := range lines {
				fmt.Fprintln(in, line)
			}
		})
	}
	if err != nil {
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
