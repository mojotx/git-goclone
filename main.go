package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	// Get list of URLs to clone
	var errCount int
	for _, a := range os.Args[1:] {
		if err := processUrl(a); err != nil {
			errCount++
		}
	}

	os.Exit(errCount)
}

func processUrl(gitUrl string) error {
	log.Info().Msgf("processing %s", gitUrl)

	// Save the current working directory, so we can make sure that we change back here.
	// Do the chdir in a deferred closure
	if savedWD, err := os.Getwd(); err != nil {
		log.Error().Err(err).Msg("cannot determine current directory name")
		return err
	} else {

		// Change back to the saved working directory
		defer func(s string) {

			if cerr := os.Chdir(s); cerr != nil {
				log.Error().Err(cerr).Msgf("cannot change back to %s", s)
			}

		}(savedWD)

	}

	url, err := gitUrls.Parse(gitUrl)
	if err != nil {
		log.Error().Err(err).Msgf("cannot parse %s", gitUrl)
		return err
	} else if url == nil {
		log.Error().Msgf("this should never happen, but url is a nil pointer (%s)", gitUrl)
		return fmt.Errorf("url parse failed")
	}

	dir, file := filepath.Split(url.Path)
	log.Debug().Msgf("path is '%s', dir is '%s', and file is '%s'", url.Path, dir, file)

	if dir[0] == filepath.Separator {
		dir = dir[1:]
	}

	log.Debug().Msgf("----- dir is now '%s' ----------------------", dir)

	// Create the directory or directories
	if merr := os.MkdirAll(dir, 0755); merr != nil {
		log.Error().Err(merr).Msgf("cannot create dirs '%s'", dir)
		return merr
	}

	/*
	// Change to the directory
	if cdErr := os.Chdir(dir); cdErr != nil {
		log.Error().Err(cdErr).Msgf("cannot change to new directory '%s'", dir)
		return cdErr
	}

	 */

	log.Info().Msgf("Cloning repo %s into %s", gitUrl, dir)

	// Clone the repository
	if _, gErr := git.PlainClone(dir, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		log.Error().Err(gErr).Msgf("cannot clone repo %s", gitUrl)
		return gErr
	}

	return nil
}
