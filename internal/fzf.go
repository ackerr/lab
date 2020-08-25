package internal

import (
	"bufio"
	"io"
	"os"

	"github.com/ackerr/lab/utils"
	"github.com/ktr0731/go-fuzzyfinder"
)

func readLines(filePath string) (lines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		utils.Err("If no ~/.projects file, you should run `lab sync` first")
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

// Fuzzy : fuzzy finder
func Fuzzy(filePath string) string {
	lines, err := readLines(filePath)
	if err != nil {
		utils.Err(err)
	}
	index, err := fuzzyfinder.Find(lines, func(i int) string {
		return lines[i]
	})
	if err != nil {
		utils.Err(err)
	}
	return lines[index]
}
