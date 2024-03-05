package main

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
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
func TestProcessUrl(t *testing.T) {
	// Mock the external dependencies
	mockParser := &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/example"}, nil
		},
	}
	mockCloner := &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			return nil, nil
		},
	}

	// Call processUrl with the mocked dependencies
	err := processUrl("https://example.com", mockParser, mockCloner)

	// Assert that no error was returned
	if err != nil {
		t.Errorf("processUrl returned an error: %v", err)
	}
}
