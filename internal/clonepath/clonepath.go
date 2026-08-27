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

	if err := ValidateWithinCWD(cleaned); err != nil {
		return "", err
	}

	return cleaned, nil
}

// ValidateWithinCWD ensures an already-relative destination still resolves under
// the current working directory after canonicalization and symlink resolution.
func ValidateWithinCWD(dest string) error {
	if dest == "" || dest == "." {
		return fmt.Errorf("empty path after sanitization: %q", dest)
	}
	if strings.Contains(dest, "..") {
		return fmt.Errorf("path traversal detected in: %q", dest)
	}
	if filepath.IsAbs(dest) {
		return fmt.Errorf("absolute paths not allowed: %q", dest)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: %w", err)
	}
	resolvedCWD, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return fmt.Errorf("cannot resolve working directory: %w", err)
	}

	candidate := filepath.Join(cwd, dest)
	resolvedCandidate, err := resolveCanonicalPath(resolvedCWD, candidate)
	if err != nil {
		return fmt.Errorf("path resolves outside working directory: %q", dest)
	}
	if !isWithinDir(resolvedCWD, resolvedCandidate) {
		return fmt.Errorf("path resolves outside working directory: %q", dest)
	}

	return nil
}

func resolveCanonicalPath(base, candidate string) (string, error) {
	candidate = filepath.Clean(candidate)
	candidateRel, err := filepath.Rel(base, candidate)
	if err != nil {
		return "", err
	}
	if candidateRel == ".." || strings.HasPrefix(candidateRel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path resolves outside working directory")
	}

	cur := base
	parts := strings.Split(candidateRel, string(filepath.Separator))
	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if part == ".." {
			return "", fmt.Errorf("path traversal detected")
		}

		next := filepath.Join(cur, part)
		if info, err := os.Lstat(next); err == nil && info.Mode()&os.ModeSymlink != 0 {
			target, err := filepath.EvalSymlinks(next)
			if err != nil {
				return "", err
			}
			cur = target
			continue
		}

		if _, err := os.Stat(next); err == nil {
			target, err := filepath.EvalSymlinks(next)
			if err != nil {
				return "", err
			}
			cur = target
			continue
		}

		cur = next
	}

	return cur, nil
}

func isWithinDir(base, candidate string) bool {
	rel, err := filepath.Rel(base, candidate)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func trimGitSuffix(path string) string {
	return strings.TrimSuffix(path, ".git")
}

func trimLeadingSeparators(path string) string {
	return strings.TrimLeft(path, `/\`)
}
