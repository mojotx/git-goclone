package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	gitUrls "github.com/whilp/git-urls"

	"github.com/mojotx/git-goclone/pkg/msg"
	"github.com/mojotx/git-goclone/pkg/path"
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

func processUrl(gitUrl string) error {
	fmt.Printf("processing %s...\n", gitUrl)

	url, err := gitUrls.Parse(gitUrl)
	if err != nil {
		msg.Err("cannot parse %s: %s\n", gitUrl, err.Error())
		return err
	} else if url == nil {
		msg.Err("this should never happen, but url is a nil pointer (%s)\n", gitUrl)
		return fmt.Errorf("url parse failed")
	}

	// See if the path ends with .git, and if so, remove that before cloning
	clonePath := path.Sanitize(url.Path)

	msg.Info("Cloning repo %s into %s...\n", gitUrl, clonePath)

	// Clone the repository
	if _, gErr := git.PlainClone(clonePath, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		msg.Err("cannot clone repo %s: %s\n", gitUrl, gErr.Error())
		return gErr
	}

	return nil
}
