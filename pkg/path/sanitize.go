package path

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// gitPattern matches .git at the end of a path
	gitPattern = regexp.MustCompile(`.git$`)
	// pathPattern matches leading path separators
	pathPattern = regexp.MustCompile(`^[/\\]*`)
)

// Sanitize is a wrapper around some helper functions that clean up the url path so that we can clone
// the repository on the file system. It validates against path traversal attacks and ensures the
// destination is within the current working directory.
func Sanitize(path string) (string, error) {
	// Remove leading slashes and .git suffix
	cleaned := PostTrim(PreTrim(path))

	// Use filepath.Clean to normalize the path and resolve . and ..
	cleaned = filepath.Clean(cleaned)

	// Strip any leading path separator again after cleaning
	cleaned = strings.TrimPrefix(cleaned, string(filepath.Separator))

	// Explicitly check for path traversal attempts
	if strings.Contains(cleaned, "..") {
		return "", fmt.Errorf("path traversal detected in: %s", path)
	}

	// Prevent absolute paths
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute paths not allowed: %s", path)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %w", err)
	}

	// Build the absolute path to verify it's under CWD
	absPath := filepath.Join(cwd, cleaned)

	// Ensure the resolved path is still under the working directory
	if !strings.HasPrefix(absPath, cwd) {
		return "", fmt.Errorf("path escapes working directory: %s", path)
	}

	return cleaned, nil
}

// PostTrim checks to see if the string ends with the substring ".git", and if so, removes it.
func PostTrim(path string) string {
	// Find the start and end indices of the pattern
	if indices := gitPattern.FindStringIndex(path); indices != nil {
		path = path[:indices[0]]
	}

	return path
}

// PreTrim Checks to see if the first character is a path separator, and if so, removes it.
func PreTrim(path string) string {
	// Find the start and end indices of the pattern
	if indices := pathPattern.FindStringIndex(path); indices != nil {
		path = path[indices[1]:]
	}

	return path
}
