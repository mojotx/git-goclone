package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
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

// redactURL removes credentials from URLs before logging to prevent password exposure
func redactURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "[invalid URL]"
	}

	// Build the redacted URL manually to avoid URL encoding of the asterisks
	var sb strings.Builder

	if u.Scheme != "" {
		sb.WriteString(u.Scheme)
		sb.WriteString("://")
	}

	if u.User != nil {
		username := u.User.Username()
		sb.WriteString(username)
		if _, hasPass := u.User.Password(); hasPass {
			sb.WriteString(":***")
		}
		sb.WriteString("@")
	}

	sb.WriteString(u.Host)
	sb.WriteString(u.Path)

	if u.RawQuery != "" {
		sb.WriteString("?")
		sb.WriteString(u.RawQuery)
	}

	if u.Fragment != "" {
		sb.WriteString("#")
		sb.WriteString(u.Fragment)
	}

	return sb.String()
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

	// Exit with 1 if any errors occurred, 0 otherwise
	if errCount > 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func processUrl(gitUrl string, parser GitURLsParser, cloner GitCloner) error {
	// Redact credentials before logging
	redactedURL := redactURL(gitUrl)
	log.Info().Str("url", redactedURL).Msg("processing")
	// Use parser and cloner instead of directly calling gitUrls.Parse and git.PlainClone

	url, err := parser.Parse(gitUrl)
	if err != nil {
		log.Error().Err(err).Str("url", redactedURL).Msg("cannot parse")
		return err
	} else if url == nil {
		log.Error().Str("url", redactedURL).Msg("url parse returned nil")
		return fmt.Errorf("url parse failed")
	}

	// See if the path ends with .git, and if so, remove that before cloning
	clonePath, sanitizeErr := path.Sanitize(url.Path)
	if sanitizeErr != nil {
		log.Error().Err(sanitizeErr).Str("url", redactedURL).Str("path", url.Path).Msg("invalid path")
		return sanitizeErr
	}

	// Check if destination already exists
	if _, statErr := os.Stat(clonePath); statErr == nil {
		err := fmt.Errorf("destination already exists: %s", clonePath)
		log.Error().Err(err).Str("url", redactedURL).Str("clonePath", clonePath).Msg("cannot clone")
		return err
	}

	log.Info().Str("url", redactedURL).Str("clonePath", clonePath).Msg("cloning repo")

	// Create context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Clone the repository with shallow clone to limit resource usage
	// Note: Context-aware cloning requires manual implementation or newer go-git version
	// For now, we use Depth:1 to limit resource consumption
	cloneDone := make(chan error, 1)
	go func() {
		_, gErr := cloner.PlainClone(clonePath, false, &git.CloneOptions{
			URL:      gitUrl,
			Progress: os.Stderr,
			Depth:    1,
		})
		cloneDone <- gErr
	}()

	// Wait for clone to complete or timeout
	select {
	case gErr := <-cloneDone:
		if gErr != nil {
			log.Error().Err(gErr).Str("url", redactedURL).Str("clonePath", clonePath).Msg("cannot clone repo")
			return gErr
		}
	case <-ctx.Done():
		err := fmt.Errorf("clone timeout after 5 minutes")
		log.Error().Err(err).Str("url", redactedURL).Str("clonePath", clonePath).Msg("cannot clone repo")
		return err
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
