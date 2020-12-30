package utils

import (
	"io/ioutil"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func CopyFile(source string, target string) {
	buf, err := ioutil.ReadFile(source)
	print("???")
	Check(err)
	print("??")
	err = ioutil.WriteFile(target, buf, 0644)
	Check(err)
}
