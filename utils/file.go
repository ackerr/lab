package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	FilePerm = 0644
	DirPerm  = 0755
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func ReadLines(filePath string) (lines []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Sprintf("file %s doesn't exist, please run `lab sync` first", filePath)
		Err(msg)
	}
	defer file.Close()
	buffer := bufio.NewReader(file)
	for {
		value, _, e := buffer.ReadLine()
		if e == io.EOF {
			break
		}
		lines = append(lines, string(value))
	}
	return lines, err
}
