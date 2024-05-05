package main

import (
	"fmt"
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
	// Test case when ParseFunc returns an error
	mockParser := &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return nil, fmt.Errorf("mock error")
		},
	}
	mockCloner := &MockGitCloner{}
	err := processUrl("https://example.com", mockParser, mockCloner)
	if err == nil || err.Error() != "mock error" {
		t.Errorf("Expected 'mock error', got %v", err)
	}

	// Test case when PlainCloneFunc returns an error
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/example"}, nil
		},
	}
	mockCloner = &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			return nil, fmt.Errorf("mock error")
		},
	}
	err = processUrl("https://example.com", mockParser, mockCloner)
	if err == nil || err.Error() != "mock error" {
		t.Errorf("Expected 'mock error', got %v", err)
	}

	// Test case when ParseFunc returns a different URL
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/different"}, nil
		},
	}
	mockCloner = &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			return nil, nil
		},
	}
	err = processUrl("https://example.com", mockParser, mockCloner)
	if err != nil {
		t.Errorf("processUrl returned an error: %v", err)
	}
}

func TestParse(t *testing.T) {
	// Mock the external dependencies
	mockParser := &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/example"}, nil
		},
	}

	// Call Parse with the mocked dependencies
	url, err := mockParser.Parse("https://example.com")

	// Assert that the URL was parsed correctly
	if url.Path != "/example" {
		t.Errorf("Parse returned an unexpected URL: %v", url)
	}

	// Assert that no error was returned
	if err != nil {
		t.Errorf("Parse returned an error: %v", err)
	}
}

func TestPlainClone(t *testing.T) {
	// Mock the external dependencies
	mockCloner := &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			return nil, nil
		},
	}

	// Call PlainClone with the mocked dependencies
	repo, err := mockCloner.PlainClone("/path/to/repo", false, &git.CloneOptions{})

	// Assert that no error was returned
	if err != nil {
		t.Errorf("PlainClone returned an error: %v", err)
	}

	// Assert that the repository was cloned correctly
	if repo != nil {
		t.Errorf("PlainClone returned an unexpected repository: %v", repo)
	}
}
