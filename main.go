package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gitUrls "github.com/whilp/git-urls"
)

func main() {
	// Set local time to UTC
	go func() {
		utc, err := time.LoadLocation("UTC")
		if err != nil {
			panic(err.Error())
		} else {
			time.Local = utc
		}
	}()

	// Initialize zerolog
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Get list of URLs to clone
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

	if m, err := regexp.MatchString( pattern, path); m {
		return path[0 : len(path)-4]
	} else if err != nil {
		log.Error().Err(err).Msgf("weird error trying to match '%s' against regexp '%s'", path, pattern)
	}
	return path
}

func processUrl(gitUrl string) error {
	log.Info().Msgf("processing %s", gitUrl)

	url, err := gitUrls.Parse(gitUrl)
	if err != nil {
		log.Error().Err(err).Msgf("cannot parse %s", gitUrl)
		return err
	} else if url == nil {
		log.Error().Msgf("this should never happen, but url is a nil pointer (%s)", gitUrl)
		return fmt.Errorf("url parse failed")
	}

	dir, file := filepath.Split(url.Path)
	log.Debug().Msgf("url is %+v", *url)
	log.Debug().Msgf("path is '%s', dir is '%s', and file is '%s'", url.Path, dir, file)

	if dir[0] == filepath.Separator {
		dir = dir[1:]
	}

	log.Debug().Msgf("----- dir is now '%s' ----------------------", dir)


	// See if the path ends with .git, and if so, remove that before cloning
	clonePath := sanitize( url.Path )


	log.Info().Msgf("Cloning repo %s into %s", gitUrl, clonePath)

	// Clone the repository
	if _, gErr := git.PlainClone(clonePath, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		log.Error().Err(gErr).Msgf("cannot clone repo %s", gitUrl)
		return gErr
	}

	return nil
}
