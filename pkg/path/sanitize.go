package path

import (
	"regexp"
)

// Sanitize is a wrapper around some helper functions that clean up the url path so that we can clone
// the repository on the file system.
func Sanitize(path string) string {
	return PostTrim(PreTrim(path))
}

// PostTrim checks to see if the string ends with the substring ".git", and if so, removes it.
func PostTrim(path string) string {

	// Match \.git$ on the end
	gitPattern := regexp.MustCompile(`.git$`)

	// Find the start and end indices of the pattern
	if indices := gitPattern.FindStringIndex(path); indices != nil {
		path = path[:indices[0]]
	}

	return path
}

// PreTrim Checks to see if the first character is a path separator, and if so, removes it.
func PreTrim(path string) string {

	// Check to see if the first character is a path separator, and if so, remove it
	pathPattern := regexp.MustCompile(`^[/\\]*`)

	// Find the start and end indices of the pattern
	if indices := pathPattern.FindStringIndex(path); indices != nil {
		path = path[indices[1]:]
	}

	return path
}
