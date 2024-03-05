package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	gitUrls "github.com/mojotx/git-urls"

	"github.com/mojotx/git-goclone/pkg/path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})
}

// GitURLsParser is an interface for parsing Git URLs.
type GitURLsParser interface {
	Parse(gitUrl string) (*url.URL, error)
}

// GitCloner is an interface for cloning Git repositories.
type GitCloner interface {
	PlainClone(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error)
}

func main() {
	var errCount int

	for _, a := range os.Args[1:] {
		if err := processUrl(a, &RealGitURLsParser{}, &RealGitCloner{}); err != nil {
			errCount++
		}
	}

	os.Exit(errCount)
}

func processUrl(gitUrl string, parser GitURLsParser, cloner GitCloner) error {
	log.Info().Str("url", gitUrl).Msg("processing")
	// Use parser and cloner instead of directly calling gitUrls.Parse and git.PlainClone

	url, err := parser.Parse(gitUrl)
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
	if _, gErr := cloner.PlainClone(clonePath, false, &git.CloneOptions{URL: gitUrl, Progress: os.Stderr}); gErr != nil {
		log.Error().Err(gErr).Str("url", gitUrl).Str("clonePath", clonePath).Msg("cannot clone repo")
		return gErr
	}

	return nil
}

// MockGitURLsParser is a mock for gitUrls.Parse
type MockGitURLsParser struct {
	ParseFunc func(string) (*url.URL, error)
}

func (m *MockGitURLsParser) Parse(gitUrl string) (*url.URL, error) {
	return m.ParseFunc(gitUrl)
}

// MockGitCloner is a mock for git.PlainClone
type MockGitCloner struct {
	PlainCloneFunc func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error)
}

func (m *MockGitCloner) PlainClone(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
	return m.PlainCloneFunc(path, isBare, options)
}

// RealGitURLsParser is a real implementation of the GitURLsParser interface.
type RealGitURLsParser struct{}

func (r *RealGitURLsParser) Parse(gitUrl string) (*url.URL, error) {
	return gitUrls.Parse(gitUrl)
}

// RealGitCloner is a real implementation of the GitCloner interface.
type RealGitCloner struct{}

func (r *RealGitCloner) PlainClone(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, isBare, options)
}
