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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})

}

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with username and password",
			input:    "https://user:password@github.com/repo.git",
			expected: "https://user:***@github.com/repo.git",
		},
		{
			name:     "URL with only username",
			input:    "https://user@github.com/repo.git",
			expected: "https://user@github.com/repo.git",
		},
		{
			name:     "URL without credentials",
			input:    "https://github.com/repo.git",
			expected: "https://github.com/repo.git",
		},
		{
			name:     "SSH URL with credentials",
			input:    "ssh://user:pass@git.example.com/repo.git",
			expected: "ssh://user:***@git.example.com/repo.git",
		},
		{
			name:     "Invalid URL with scheme",
			input:    "ht!tp://invalid",
			expected: "[invalid URL]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactURL(tt.input)
			assert.Equal(t, tt.expected, result, "redactURL output mismatch")
		})
	}
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
	assert.Error(t, err, "Expected error when parser fails")
	assert.EqualError(t, err, "mock error")

	// Test case when parser returns nil URL
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return nil, nil
		},
	}
	err = processUrl("https://example.com", mockParser, mockCloner)
	assert.Error(t, err, "Expected error when URL is nil")
	assert.EqualError(t, err, "url parse failed")

	// Test case when sanitize returns an error (path traversal)
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/../../../etc/passwd"}, nil
		},
	}
	err = processUrl("https://evil.com/../../../etc/passwd", mockParser, mockCloner)
	assert.Error(t, err, "Expected error for path traversal attack")
	assert.Contains(t, err.Error(), "path traversal detected")

	// Test case when destination already exists
	// First create a temporary directory
	tempDir := t.TempDir()
	existingPath := "test-repo"
	fullPath := tempDir + "/" + existingPath
	err = os.MkdirAll(fullPath, 0755)
	require.NoError(t, err, "Failed to create test directory")

	// Change to temp directory
	oldWd, err := os.Getwd()
	require.NoError(t, err, "Failed to get working directory")
	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change directory")
	t.Cleanup(func() {
		assert.NoError(t, os.Chdir(oldWd), "Failed to restore working directory")
	})

	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/" + existingPath}, nil
		},
	}
	err = processUrl("https://example.com/test-repo", mockParser, mockCloner)
	assert.Error(t, err, "Expected error when destination exists")
	assert.Contains(t, err.Error(), "destination already exists")

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
	assert.Error(t, err, "Expected error when clone fails")
	assert.EqualError(t, err, "mock error")

	// Test case for successful clone
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/different"}, nil
		},
	}
	mockCloner = &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			// Verify shallow clone is configured
			assert.Equal(t, 1, options.Depth, "Expected shallow clone with Depth=1")
			return nil, nil
		},
	}
	err = processUrl("https://example.com", mockParser, mockCloner)
	assert.NoError(t, err, "processUrl should succeed with valid inputs")

	// Test case for timeout (clone taking too long)
	// Note: This is hard to test comprehensively without actually timing out
	// but we can at least verify the goroutine and channel setup works
	mockParser = &MockGitURLsParser{
		ParseFunc: func(gitUrl string) (*url.URL, error) {
			return &url.URL{Path: "/timeout-test"}, nil
		},
	}
	mockCloner = &MockGitCloner{
		PlainCloneFunc: func(path string, isBare bool, options *git.CloneOptions) (*git.Repository, error) {
			// Simulate a quick successful clone (not an actual timeout)
			return nil, nil
		},
	}
	err = processUrl("https://example.com", mockParser, mockCloner)
	assert.NoError(t, err, "processUrl should succeed with valid inputs")
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
	assert.NoError(t, err, "Parse should not return error")
	require.NotNil(t, url, "Parse should return non-nil URL")
	assert.Equal(t, "/example", url.Path, "Parse should return correct path")
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
	assert.NoError(t, err, "PlainClone should not return error")
	// Assert that the repository was cloned correctly
	assert.Nil(t, repo, "PlainClone should return nil repository in this mock")
}

func TestRealImplementations(t *testing.T) {
	// Test RealGitURLsParser
	parser := &RealGitURLsParser{}
	url, err := parser.Parse("https://github.com/mojotx/git-goclone.git")

	assert.NoError(t, err, "RealGitURLsParser.Parse should not return error")
	require.NotNil(t, url, "RealGitURLsParser.Parse should return non-nil URL")
	assert.Equal(t, "github.com", url.Host, "URL should have correct host")
}

func TestRedactURLEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "URL with query parameters",
			input: "https://user:pass@example.com/path?key=value",
		},
		{
			name:  "URL with fragment",
			input: "https://user:pass@example.com/path#section",
		},
		{
			name:  "URL with port",
			input: "https://user:pass@example.com:8080/repo.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactURL(tt.input)
			// Verify password is redacted
			assert.NotContains(t, result, "pass", "Password should be redacted")
			// Verify username is preserved
			assert.Contains(t, result, "user", "Username should be preserved")
			// Verify *** is present (password redaction marker)
			assert.Contains(t, result, "***", "Redaction marker should be present")
		})
	}
}
