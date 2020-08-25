package utils

import (
	"errors"
	"github.com/kballard/go-shellquote"
	"os"
	"runtime"
)

// BrowserLauncher : return the browser launcher command, copy from github.com/github/hub
func BrowserLauncher() ([]string, error) {
	browser := os.Getenv("BROWSER")
	if browser == "" {
		browser = searchBrowserLauncher(runtime.GOOS)
	} else {
		browser = os.ExpandEnv(browser)
	}

	if browser == "" {
		return nil, errors.New("please set $BROWSER to a web launcher")
	}
	return shellquote.Split(browser)
}

func searchBrowserLauncher(goos string) (browser string) {
	switch goos {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	case "linux":
		browser = "xdg-open"
	default:
		browser = ""
	}
	return browser
}
