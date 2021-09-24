package path

import (
	"os"
	"regexp"

	"github.com/mojotx/git-goclone/pkg/msg"
)

// Sanitize is a wrapper around some helper functions that clean up the url path so that we can clone
// the repository on the file system.
func Sanitize(path string) string {
	return PostTrim(PreTrim(path))
}

// PostTrim checks to see if the string ends with the substring ".git", and if so, removes it.
func PostTrim(path string) string {

	gitPattern := ".git$"
	if m, err := regexp.MatchString(gitPattern, path); m {
		return path[0 : len(path)-4]
	} else if err != nil {
		msg.Err("weird error trying to match '%s' against regexp '%s': %s", path, gitPattern, err.Error())
		os.Exit(-1)
	}
	return path
}

// PreTrim Checks to see if the first character is a path separator, and if so, removes it.
func PreTrim(path string) string {
	// Check to see if the first character is a path separator, and if so, remove it
	pathPattern := `^[/\\]`
	if m, err := regexp.MatchString(pathPattern, path); m {
		path = path[1:]
	} else if err != nil {
		msg.Err("weird error trying to match '%s' against regexp '%s': %s", path, pathPattern, err.Error())
		os.Exit(-1)
	}
	return path
}
