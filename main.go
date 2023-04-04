package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	gitUrls "github.com/whilp/git-urls"

	"github.com/mojotx/git-goclone/pkg/path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})
	var errCount int

	for _, a := range os.Args[1:] {
		if err := processUrl(a); err != nil {
			errCount++
		}
	}

	os.Exit(errCount)
}

func processUrl(gitUrl string) error {
	log.Info().Str("url", gitUrl).Msg("processing")

	url, err := gitUrls.Parse(gitUrl)
	if err != nil {
		log.Error().Err(err).Str("url", gitUrl).Msg("cannot parse")
		return err
	} else if url == nil {
		log.Panic().Str("url", gitUrl).Msg("never happen, but url is nil")
		return fmt.Errorf("url parse failed")
	}

	// See if the path ends with .git, and if so, remove that before cloning
	clonePath := path.Sanitize(url.Path)

	log.Info().Str("url", gitUrl).Str("clonePath", clonePath).Msg("cloning repo")

	// Clone the repository
	if _, gErr := git.PlainClone(clonePath, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		log.Error().Err(gErr).Str("url", gitUrl).Str("clonePath", clonePath).Msg("cannot clone repo")
		return gErr
	}

	return nil
}
