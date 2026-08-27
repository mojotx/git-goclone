package clone

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func silentLogger() zerolog.Logger {
	return zerolog.New(io.Discard)
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { _ = os.Chdir(orig) })
	return dir
}

func TestURLWith_ParseError(t *testing.T) {
	parse := func(string) (*url.URL, error) { return nil, errors.New("boom") }
	clone := func(context.Context, string, bool, *git.CloneOptions) (*git.Repository, error) {
		t.Fatal("cloner must not be called")
		return nil, nil
	}
	err := URLWith(context.Background(), silentLogger(), "https://x", Options{}, parse, clone)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestURLWith_NilParsedURL(t *testing.T) {
	parse := func(string) (*url.URL, error) { return nil, nil }
	clone := func(context.Context, string, bool, *git.CloneOptions) (*git.Repository, error) {
		t.Fatal("cloner must not be called")
		return nil, nil
	}
	err := URLWith(context.Background(), silentLogger(), "https://x", Options{}, parse, clone)
	require.Error(t, err)
}

func TestURLWith_SanitizeError(t *testing.T) {
	parse := func(string) (*url.URL, error) { return &url.URL{Path: "/../../etc/passwd"}, nil }
	clone := func(context.Context, string, bool, *git.CloneOptions) (*git.Repository, error) {
		t.Fatal("cloner must not be called")
		return nil, nil
	}
	err := URLWith(context.Background(), silentLogger(), "https://x", Options{}, parse, clone)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestURLWith_DestinationExists(t *testing.T) {
	chdirTemp(t)
	require.NoError(t, os.MkdirAll(filepath.Join("org", "proj"), 0o755))

	parse := func(string) (*url.URL, error) { return &url.URL{Path: "/org/proj.git"}, nil }
	clone := func(context.Context, string, bool, *git.CloneOptions) (*git.Repository, error) {
		t.Fatal("cloner must not be called when destination exists")
		return nil, nil
	}
	err := URLWith(context.Background(), silentLogger(), "https://x/org/proj.git", Options{}, parse, clone)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestURLWith_CloneFailure(t *testing.T) {
	chdirTemp(t)

	parse := func(string) (*url.URL, error) { return &url.URL{Path: "/org/proj.git"}, nil }
	clone := func(context.Context, string, bool, *git.CloneOptions) (*git.Repository, error) {
		return nil, errors.New("network fell over")
	}
	err := URLWith(context.Background(), silentLogger(), "https://x/org/proj.git", Options{Depth: 1}, parse, clone)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "network fell over")
}

func TestURLWith_Success(t *testing.T) {
	chdirTemp(t)

	var seenPath string
	var seenDepth int
	parse := func(string) (*url.URL, error) { return &url.URL{Path: "/org/team/proj.git"}, nil }
	clone := func(_ context.Context, path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
		seenPath = path
		seenDepth = o.Depth
		assert.False(t, isBare)
		return nil, nil
	}

	err := URLWith(context.Background(), silentLogger(), "https://x/org/team/proj.git", Options{Depth: 1}, parse, clone)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join("org", "team", "proj"), seenPath)
	assert.Equal(t, 1, seenDepth)
}

func TestURLWith_CancelledContextRemovesPartial(t *testing.T) {
	dir := chdirTemp(t)

	parse := func(string) (*url.URL, error) { return &url.URL{Path: "/org/proj.git"}, nil }
	clone := func(_ context.Context, path string, _ bool, _ *git.CloneOptions) (*git.Repository, error) {
		// Simulate go-git leaving a partial directory behind on cancel.
		require.NoError(t, os.MkdirAll(path, 0o755))
		return nil, context.DeadlineExceeded
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	time.Sleep(5 * time.Millisecond) // ensure ctx is already expired

	err := URLWith(ctx, silentLogger(), "https://x/org/proj.git", Options{}, parse, clone)
	require.Error(t, err)
	_, statErr := os.Stat(filepath.Join(dir, "org", "proj"))
	assert.True(t, os.IsNotExist(statErr), "partial clone should have been cleaned up")
}

func TestRedact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		in              string
		wantContains    string
		wantNotContains string
	}{
		{"with password", "https://user:secret@example.com/x", "xxxxx", "secret"},
		{"only username", "https://user@example.com/x", "user@example.com", ""},
		{"no credentials", "https://example.com/x", "example.com", ""},
		{"scp-style ssh", "git@github.com:mojotx/git-goclone.git", "mojotx/git-goclone.git", ""},
		{"ssh scheme", "ssh://git@github.com/mojotx/git-goclone.git", "git@github.com", ""},
		{"ssh with password", "ssh://user:secret@host/repo.git", "xxxxx", "secret"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := redact(tt.in)
			assert.Contains(t, got, tt.wantContains)
			if tt.wantNotContains != "" {
				assert.NotContains(t, got, tt.wantNotContains)
			}
		})
	}
}
