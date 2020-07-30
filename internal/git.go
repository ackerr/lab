package internal

import (
	"fmt"
	"os/exec"
)

// Clone : git clone the gitlab project
func Clone(giturl, path string) {
	fmt.Println(giturl, path)
	cmd := exec.Command("git", "clone", giturl, path)
	err := cmd.Run()
	if err != nil {
		Err(err)
	}
	// _, err := git.PlainClone(path, false, &git.CloneOptions{
	// 	URL:      url,
	// 	Progress: os.Stdout,
	// })
	// if err != nil {
	// 	Err(err)
	// }
}
