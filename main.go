package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	gitUrls "github.com/whilp/git-urls"
)

const (
	red = "\x1b[1;31m"
	green = "\x1b[32m"
	reset = "\x1b[0m"
)

func main() {
	var errCount int

	for _, a := range os.Args[1:] {
		if err := processUrl(a); err != nil {
			errCount++
		}
	}

	os.Exit(errCount)
}

func sanitize(path string) string {

	pattern := ".git$"

	if m, err := regexp.MatchString(pattern, path); m {
		return path[0 : len(path)-4]
	} else if err != nil {
		fmt.Printf("weird error trying to match '%s' against regexp '%s': %s", path, pattern, err.Error())
		os.Exit(-1)
	}
	return path
}

func processUrl(gitUrl string) error {
	fmt.Printf("processing %s...\n", gitUrl)

	url, err := gitUrls.Parse(gitUrl)
	if err != nil {
		fmt.Printf(red + "cannot parse %s: %s\n" + reset, gitUrl, err.Error())
		return err
	} else if url == nil {
		fmt.Printf(red + "this should never happen, but url is a nil pointer (%s)\n" + reset, gitUrl)
		return fmt.Errorf("url parse failed")
	}

	// See if the path ends with .git, and if so, remove that before cloning
	clonePath := sanitize(url.Path)

	fmt.Printf(green + "Cloning repo %s into %s...\n" + reset, gitUrl, clonePath)

	// Clone the repository
	if _, gErr := git.PlainClone(clonePath, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		fmt.Printf(red + "cannot clone repo %s: %s\n" + reset, gitUrl, gErr.Error())
		return gErr
	}

	return nil
}
