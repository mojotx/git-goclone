// Package clone contains the core "clone one URL into a nested directory" logic.
package clone

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5"
	gitUrls "github.com/mojotx/git-urls"
	"github.com/rs/zerolog"

	"github.com/mojotx/git-goclone/internal/clonepath"
)

// Options controls a single clone operation.
type Options struct {
	// Depth is the clone depth; 0 means a full clone.
	Depth int
	// Progress receives git's progress output; nil discards it.
	Progress io.Writer
}

// URLParser parses a raw Git URL string into a *url.URL. Injected for testing.
type URLParser func(string) (*url.URL, error)

// Cloner performs the actual clone. Injected for testing.
type Cloner func(ctx context.Context, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)

// DefaultURLParser is the production URL parser.
var DefaultURLParser URLParser = gitUrls.Parse

// DefaultCloner is the production cloner.
var DefaultCloner Cloner = git.PlainCloneContext

// URL clones a single repository into a directory that mirrors the URL path,
// e.g. "https://gitlab.com/org/team/proj.git" -> "./org/team/proj".
func URL(ctx context.Context, logger zerolog.Logger, rawURL string, opts Options) error {
	return URLWith(ctx, logger, rawURL, opts, DefaultURLParser, DefaultCloner)
}

// URLWith is URL with injectable dependencies for testing.
func URLWith(
	ctx context.Context,
	logger zerolog.Logger,
	rawURL string,
	opts Options,
	parse URLParser,
	clone Cloner,
) error {
	redacted := redact(rawURL)
	logger.Info().Str("url", redacted).Msg("processing")

	parsed, err := parse(rawURL)
	if err != nil {
		logger.Error().Err(err).Str("url", redacted).Msg("cannot parse URL")
		return fmt.Errorf("parse URL: %w", err)
	}
	if parsed == nil {
		logger.Error().Str("url", redacted).Msg("URL parser returned nil")
		return errors.New("URL parser returned nil")
	}

	dest, err := clonepath.Sanitize(parsed.Path)
	if err != nil {
		logger.Error().Err(err).Str("url", redacted).Str("path", parsed.Path).Msg("invalid path")
		return err
	}

	if _, statErr := os.Stat(dest); statErr == nil {
		logger.Error().Str("url", redacted).Str("dest", dest).Msg("destination already exists")
		return fmt.Errorf("destination already exists: %s", dest)
	}

	logger.Info().Str("url", redacted).Str("dest", dest).Int("depth", opts.Depth).Msg("cloning")

	progress := opts.Progress
	if progress == nil {
		progress = io.Discard
	}

	_, cloneErr := clone(ctx, dest, false, &git.CloneOptions{
		URL:      rawURL,
		Progress: progress,
		Depth:    opts.Depth,
	})
	if cloneErr != nil {
		// If the context was cancelled or timed out, tidy up any partial clone.
		if ctx.Err() != nil {
			if removeErr := os.RemoveAll(dest); removeErr != nil {
				logger.Warn().Err(removeErr).Str("dest", dest).Msg("failed to remove partial clone")
			}
		}
		logger.Error().Err(cloneErr).Str("url", redacted).Str("dest", dest).Msg("clone failed")
		return fmt.Errorf("clone %s: %w", redacted, cloneErr)
	}

	return nil
}

// redact returns rawURL with any embedded password replaced. Falls back to a
// placeholder if the URL cannot be parsed at all.
func redact(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "[unparseable URL]"
	}
	return u.Redacted()
}
