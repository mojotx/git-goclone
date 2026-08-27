// Package clonepath sanitizes URL paths into safe on-disk clone destinations.
package clonepath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Sanitize converts a URL path into a relative filesystem path suitable for
// cloning into the current working directory. It strips leading separators
// and a trailing ".git" suffix, then rejects any input that would escape the
// working directory (absolute paths, "..", etc.).
func Sanitize(path string) (string, error) {
	cleaned := trimGitSuffix(trimLeadingSeparators(path))

	cleaned = filepath.Clean(cleaned)
	cleaned = strings.TrimPrefix(cleaned, string(filepath.Separator))

	if cleaned == "" || cleaned == "." {
		return "", fmt.Errorf("empty path after sanitization: %q", path)
	}
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal detected in: %q", path)
	}
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute paths not allowed: %q", path)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %w", err)
	}

	// Guard against a prefix collision like cwd="/tmp/work" vs "/tmp/workstuff".
	absPath := filepath.Join(cwd, cleaned)
	if !strings.HasPrefix(absPath, cwd+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes working directory: %q", path)
	}

	return cleaned, nil
}

func trimGitSuffix(path string) string {
	return strings.TrimSuffix(path, ".git")
}

func trimLeadingSeparators(path string) string {
	return strings.TrimLeft(path, `/\`)
}
